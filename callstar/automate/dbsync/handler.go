package dbsync

import (
	"fmt"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/callstar/template"
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

		if job != nil && job.Tries >= maxRefusion { // timeout: check record error counter
			job.Task = 9999 // mark job as done
			log.Warn().Msgf("%s addr '%s' drop job '%s', err=%s", name, addr, job.Action.String(), fmt.Sprintf("vendor did not answer %d times", maxRefusion))
			automate.NextAction(name, addr, nextRecord, t)
			return
		}

		action.NextState = nextAction
		err = p.sendPacket(addr, &action, template.PacketEnq, "", nil)

	case nextAction:

		dispatcher := automate.Dispatcher()
		job := dispatcher.GetSyncRecord(addr)

		if job == nil { // should not happend, safty check
			automate.NextAction(name, addr, success, t)
			return
		}

		action.CurrentState = busy // on error fallback to send enquiry

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
			if _, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {
				automate.NextAction(name, addr, busy, t) // send enquiry
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

func (p *Plugin) handlePacketRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	job := dispatcher.GetSyncRecord(addr)

	if job != nil && job.Tries >= maxRefusion {
		log.Warn().Msgf("%s addr '%s' drop job '%s', err=%s", name, addr, job.Action.String(), fmt.Sprintf("vendor refused the record %d times", maxRefusion))
		if action != nil {
			action.NextTimeout = nextActionDelay
		}
		automate.ChangeState(name, addr, nextRecord)
		return nil
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
