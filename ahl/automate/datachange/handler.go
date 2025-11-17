package datachange

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/ahl/template"

	"github.com/spf13/cast"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	if !p.linkState(addr) {
		p.handleError(addr, job, "vendor disconnected")
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case busy:

		if job == nil || !job.InProcess() { // order was canceled externally
			err = defines.ErrorJobFailed
			break
		}

		if job.Context == nil && job.Action != order.NightAuditStart && job.Action != order.NightAuditEnd {
			err = defines.ErrorJobContext
			break
		}

		if p.waitForReplyPacket(addr, job) {
			action.NextState = answer
			action.NextTimeout = answerTimeout
		}

		switch job.Action {

		case order.Checkin:
			err = p.sendPacket(addr, &action, template.PacketCheckIn, job.Initiator, job.Context)

		case order.DataChange:
			err = p.sendPacket(addr, &action, template.PacketDataChange, job.Initiator, job.Context)

		case order.Checkout:
			err = p.sendPacket(addr, &action, template.PacketCheckOut, job.Initiator, job.Context)

		case order.RoomStatus:
			err = p.sendPacket(addr, &action, template.PacketRoomStatus, job.Initiator, job.Context)

		case order.WakeupRequest:
			err = p.sendPacket(addr, &action, template.PacketWakeupSet, job.Initiator, job.Context)

		case order.WakeupClear:
			err = p.sendPacket(addr, &action, template.PacketWakeupClear, job.Initiator, job.Context)

		default:
			log.Warn().Msgf("%s addr '%s' cancel job, because no handler for action '%s' defined", name, addr, job.Action.String())
			p.handleError(addr, job, "internal error, missing handler")
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		job.Timestamp = time.Now().UnixNano()

	case answer:
		log.Warn().Msgf("%s addr '%s' finalise job '%s', vendor did not answer", name, addr, job.Action.String())
		p.handleError(addr, job, "vendor did not answer")
		automate.NextAction(name, addr, shutdown, t)
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

func (p *Plugin) waitForReplyPacket(addr string, _ *order.Job) bool {

	dispatcher := p.automate.Dispatcher()
	if !dispatcher.IsSerialDevice(addr) {
		return true
	}

	/* IFC-658
	if job != nil && (job.Action == order.Checkin || job.Action == order.Checkout) { // IFC-421
		return true
	}
	*/

	station, _ := dispatcher.GetStationAddr(addr)
	state := dispatcher.GetConfig(station, "WaitForReplyPacket", "true")
	return cast.ToBool(state)

}

func (p *Plugin) handlePacketAcknowledge(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	if !p.waitForReplyPacket(addr, job) {
		p.handleAnswerPositiv(addr, job)
	}

	return p.automate.HandlePacketAcknowledge(addr, packet, action, job)
}

func (p *Plugin) handleReply(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	name := automate.Name()

	data := packet.Data()
	extension := cast.ToString(data["Extension"])
	extension = strings.TrimLeft(extension, " ")

	expectedExtension := p.getExtension(job)
	if extension == expectedExtension {

		status := data["Status"]

		if len(status) == 2 {

			jobAction := job.Action
			if jobAction == order.WakeupRequest || jobAction == order.WakeupClear {
				jobAction = order.DataChange
			}

			answerAction := status[0]

			if (jobAction == order.Checkin && answerAction == 'I') || (jobAction == order.Checkout && answerAction == 'O') || ((jobAction == order.DataChange || jobAction == order.RoomStatus) && answerAction == 'M') {
				p.handleAnswerPositiv(addr, job)
				automate.ChangeState(name, addr, success)
				return nil
			}

			if (jobAction == order.Checkin && answerAction == 'J') || (jobAction == order.Checkout && answerAction == 'P') || ((jobAction == order.DataChange || jobAction == order.RoomStatus) && answerAction == 'N') {
				answerReason := status[1]
				p.handleAnswerNegativ(addr, job, answerReason)
				automate.ChangeState(name, addr, success)
				return nil
			}

			log.Warn().Msgf("%s addr '%s' received unexpected reply status '%s' for action '%s'", name, addr, status, job.Action)

		}

	} else {
		log.Warn().Msgf("%s addr '%s' received reply for extension '%s', expected '%s' ", name, addr, extension, expectedExtension)
	}

	p.calculateNextTimeout(action, job)
	return nil
}

func (p *Plugin) calculateNextTimeout(action *dispatcher.StateAction, job *order.Job) {
	if action == nil || job == nil {
		return
	}

	nextTimeout := answerTimeout
	if job.Timestamp != 0 {
		diff := time.Since(time.Unix(0, job.Timestamp))
		if diff < nextTimeout {
			nextTimeout = nextTimeout - diff
		}
	}
	action.NextTimeout = nextTimeout
	action.NextState = answer
}

func (p *Plugin) getExtension(job *order.Job) string {
	if job == nil || job.Context == nil {
		return ""
	}
	return p.driver.GetExtension(job.Context)
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
