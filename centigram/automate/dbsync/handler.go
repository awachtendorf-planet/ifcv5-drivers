package dbsync

import (
	"fmt"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/centigram/template"
)

const (
	maxRefusion = 3
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, _ *order.Job) {

	automate := p.automate
	name := automate.Name()

	if !p.linkState(addr) {
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: nextRecord}

	switch state {

	case busy:

		dispatcher := automate.Dispatcher()
		job := dispatcher.GetSyncRecord(addr)

		if job == nil { // should not happend, safty check
			automate.NextAction(name, addr, success, t)
			return
		}

		if job != nil && job.Tries >= maxRefusion { // check record error counter
			job.Task = 9999 // mark job as done
			log.Warn().Msgf("%s addr '%s' drop job '%s', err=%s", name, addr, job.Action.String(), fmt.Sprintf("vendor refused the record %d times", maxRefusion))
			automate.NextAction(name, addr, nextRecord, t)
			return
		}

		// send packets defined at workflow
		if packet, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {
			err = p.sendPacket(addr, &action, packet, job.Initiator, job.Context)
		} else {
			// should not happend, safty check
			automate.NextAction(name, addr, nextRecord, t)
			return
		}

		if err == nil {
			job.Tries++
		}

	case nextRecord:

		dispatcher := automate.Dispatcher()
		job := dispatcher.GetSyncRecord(addr)

		if job != nil {

			if job.Task == 0 { // first call, driver pre check

				if err := p.driver.PreCheckDataSwap(job); err != nil {
					log.Warn().Msgf("%s addr '%s' ignore job '%s', err=%s", name, addr, job.Action.String(), err)
					job.Task = 9999 // mark job as done
					automate.NextAction(name, addr, nextRecord, t)
					return
				}

			}

			job.Tries = 0
			job.Task++

			// exist more packets to send
			if packetName, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {

				switch packetName {

				case template.PacketChangeName:

					if !p.driver.SendGuestName(addr) {
						automate.NextAction(name, addr, nextRecord, t)
						return
					}

				case template.PacketChangeFCOS:

					if !p.driver.SendGuestLanguage(addr) {
						automate.NextAction(name, addr, nextRecord, t)
						return
					}

				case template.PacketMessageLampOn:

					messageLampOn := false
					if guest, ok := job.Context.(*record.Guest); ok {
						messageLampOn = guest.Reservation.MessageLightStatus
					}

					if !messageLampOn {
						automate.NextAction(name, addr, nextRecord, t)
						return
					}
				}

				automate.NextAction(name, addr, busy, t)
				return
			}
		}

		// exist more jobs to process
		if dispatcher.GetNextSyncRecord(addr) {
			automate.NextAction(name, addr, nextRecord, t)
			return
		}

		// done
		automate.NextAction(name, addr, success, t)
		return

	case success, shutdown:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handlePacket(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, j *order.Job) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	job := dispatcher.GetSyncRecord(addr)

	if packetName, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {

		switch packetName {

		case template.PacketMessageLampOn:
			if in != nil && (in.Name == template.PacketMessageWaitingStatus || in.Name == template.PacketMessageWaitingStatusOnSwap) {
				return p.automate.HandlePacketAcknowledge(addr, in, action, j)
			}
		}
	}

	return p.automate.HandlePacketRefusion(addr, in, action, job)
}

func (p *Plugin) sendPacket(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		p.automate.AddInternalErrorCounter(addr)
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}
