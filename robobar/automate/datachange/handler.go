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

	"github.com/weareplanet/ifcv5-drivers/robobar/template"

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

		if p.waitForReplyPacket(addr) {
			action.NextState = answer
			action.NextTimeout = answerTimeout
		}

		switch job.Action {

		case order.Checkin:
			err = p.sendPacket(addr, &action, template.PacketCheckIn, job.Initiator, job.Context)

		case order.Checkout:
			err = p.sendPacket(addr, &action, template.PacketCheckOut, job.Initiator, job.Context)

		case order.RoomStatus:

			driver := p.driver

			mb := driver.GetMinibarRight(job.Context)

			switch mb {

			case 0:
				err = p.sendPacket(addr, &action, template.PacketLockBar, job.Initiator, job.Context)

			case 2:
				err = p.sendPacket(addr, &action, template.PacketUnlockBar, job.Initiator, job.Context)

			default:
				log.Warn().Msgf("%s addr '%s' cancel job, because unknown minibar right: %d", name, addr, mb)
				p.handleError(addr, job, "unknown minibar right (expected 0 or 2)")
				automate.NextAction(name, addr, shutdown, t)
				return
			}

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

func (p *Plugin) waitForReplyPacket(addr string) bool {

	dispatcher := p.automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	state := dispatcher.GetConfig(station, "WaitForReplyPacket", "true")
	return cast.ToBool(state)

}

func (p *Plugin) handlePacketAcknowledge(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	if !p.waitForReplyPacket(addr) {
		p.handleAnswerPositiv(addr, job)
	}

	return p.automate.HandlePacketAcknowledge(addr, packet, action, job)
}

func (p *Plugin) handleReply(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	driver := p.driver
	name := automate.Name()

	data := packet.Data()["Data"]

	//room := driver.DecodeString(addr, data[:len(data)-2])
	room := string(data[:len(data)-2])
	room = strings.TrimLeft(room, "0")

	pmsRoom := driver.GetRoom(job.Context)

	if room == pmsRoom && (len(data)-2 >= 2) {

		result := data[len(data)-2:]
		command := int(result[0])

		if (job.Action == order.Checkin && command == 'I') || (job.Action == order.Checkout && command == 'O') {
			if result[1] == 'S' {
				p.handleAnswerPositiv(addr, job)
			} else {
				p.handleAnswerNegativ(addr, job, result[1])
			}
			action.NextTimeout = nextActionDelay
			automate.ChangeState(name, addr, success)
			return nil
		}

		mb := driver.GetMinibarRight(job.Context)

		if (job.Action == order.RoomStatus && command == 'U' && mb == 2) || (job.Action == order.RoomStatus && command == 'L' && mb == 0) {
			if result[1] == 'S' {
				p.handleAnswerPositiv(addr, job)
			} else {
				p.handleAnswerNegativ(addr, job, result[1])
			}
			action.NextTimeout = nextActionDelay
			automate.ChangeState(name, addr, success)
			return nil
		}

		log.Warn().Msgf("%s addr '%s' received unexpected reply command '%c' for action '%s'", name, addr, result[0], job.Action)

	} else {

		log.Warn().Msgf("%s addr '%s' received reply for room '%s', expected '%s' ", name, addr, room, pmsRoom)

	}

	p.setNextTimeout(action, job)
	return nil
}

func (p *Plugin) setNextTimeout(action *dispatcher.StateAction, job *order.Job) {
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

func (p *Plugin) sendPacket(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		p.automate.AddInternalErrorCounter(addr)
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}
