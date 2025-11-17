package datachange

import (
	"fmt"
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/centigram/template"

	"github.com/spf13/cast"
)

const (
	maxRefusion = 3
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// link state
	if !p.linkState(addr) {
		p.handleError(addr, job, "vendor disconnected")
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	// database swap running
	if p.driver.GetSwapState(addr) {
		p.handleError(addr, job, "database swap running")
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: nextAction}

	switch state {

	case busy:

		if job == nil { // should not happend, safty check
			automate.NextAction(name, addr, success, t)
			return
		}

		if job != nil && job.Tries >= maxRefusion { // timeout: check record error counter
			p.handleError(addr, job, "vendor did not answer")
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		// send packets defined at workflow
		if packet, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {
			err = p.sendPacket(addr, &action, packet, job.Initiator, job.Context)
		} else {
			// should not happend, safty check
			automate.NextAction(name, addr, nextAction, t)
			return
		}

		if err == nil {
			job.Tries++
		}

	case nextAction:

		if job != nil {

			job.Tries = 0
			job.Task++

			// exist more packets to send
			if packetName, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {

				switch packetName {

				case template.PacketRoomChange:

					if !p.driver.IsMove(job.Context) {
						automate.NextAction(name, addr, nextAction, t)
						return
					}

				case template.PacketChangeName:

					if !p.driver.SendGuestName(addr) {
						automate.NextAction(name, addr, nextAction, t)
						return
					}

				case template.PacketChangeFCOS:

					if !p.driver.SendGuestLanguage(addr) {
						automate.NextAction(name, addr, nextAction, t)
						return
					}

				case template.PacketMessageLampOn:

					messageLampOn := false
					if guest, ok := job.Context.(*record.Guest); ok {
						messageLampOn = guest.Reservation.MessageLightStatus
					}

					if !messageLampOn {
						automate.NextAction(name, addr, nextAction, t)
						return
					}
				}

				automate.NextAction(name, addr, busy, t) // send packet
				return
			}
		}

		// done
		p.handleSuccess(addr, job)
		automate.NextAction(name, addr, success, t)
		return

	case success, shutdown:

		t.Tick()
		return

	default:

		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)

		if err == defines.ErrorJobContext || err == defines.ErrorJobFailed {
			p.handleError(addr, job, "job was aborted")
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		internalErrorCounter := automate.GetInternalErrorCounter(addr)
		if internalErrorCounter > 0 {
			p.handleError(addr, job, err.Error())
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handlePacketRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	name := automate.Name()

	if in != nil {

		switch in.Name {

		case template.PacketRefusion:

			extension := p.getExtension(in)
			if len(extension) > 0 {
				p.handleError(addr, job, fmt.Sprintf("vendor refused the record, bad mailbox address: %s", extension))
			} else {
				p.handleError(addr, job, "vendor refused the record, bad mailbox address")
			}
			if action != nil {
				action.NextTimeout = nextActionDelay
			}
			automate.ChangeState(name, addr, shutdown)
			return nil

		default:

			if job != nil && job.Tries >= maxRefusion {
				p.handleError(addr, job, fmt.Sprintf("vendor refused the record %d times", maxRefusion))
				if action != nil {
					action.NextTimeout = nextActionDelay
				}
				automate.ChangeState(name, addr, shutdown)
				return nil
			}

		}

	}

	return p.automate.HandlePacketRefusion(addr, in, action, job)
}

func (p *Plugin) getExtension(in *ifc.LogicalPacket) string {
	if in == nil {
		return ""
	}
	data := in.Data()["Extension"]
	extension := cast.ToString(data)
	extension = strings.Trim(extension, " ")
	return extension
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
