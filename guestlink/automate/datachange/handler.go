package datachange

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"

	"github.com/pkg/errors"
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

		action.NextState = answer
		action.NextTimeout = answerTimeout

		switch job.Action {

		case order.Checkin, order.DataChange:
			err = p.sendPacket(addr, &action, template.PacketCheckIn, job.Initiator, job.Context)

		case order.Checkout:
			err = p.sendPacket(addr, &action, template.PacketCheckOut, job.Initiator, job.Context)

		case order.RoomStatus:
			err = p.sendPacket(addr, &action, template.PacketGuestMessageStatus, job.Initiator, job.Context)

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
		p.clearTransaction(addr)
		log.Warn().Msgf("%s addr '%s' finalise job '%s', vendor did not answer", name, addr, job.Action.String())
		p.handleError(addr, job, "vendor did not answer")
		automate.NextAction(name, addr, shutdown, t)
		return

	case success, shutdown:
		p.clearTransaction(addr)
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)

		// reset transaction on error
		p.driver.FreeTransaction(p, addr)

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

func (p *Plugin) handleCommandAcknowledge(addr string, _ *ifc.LogicalPacket, _ *dispatcher.StateAction, job *order.Job) error {
	p.clearTransaction(addr)

	automate := p.automate
	name := automate.Name()

	p.handleAnswerPositiv(addr, job)

	automate.ChangeState(name, addr, success)
	return nil
}

func (p *Plugin) handleCommandRefusion(addr string, packet *ifc.LogicalPacket, _ *dispatcher.StateAction, job *order.Job) error {
	p.clearTransaction(addr)

	automate := p.automate
	name := automate.Name()

	errorCode := p.driver.GetErrorCode(packet)
	p.handleAnswerNegativ(addr, job, errorCode)

	automate.ChangeState(name, addr, success)
	return nil
}

func (p *Plugin) sendPacket(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}) error {

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	transaction, err := driver.NewTransaction(addr)
	if err != nil {
		log.Error().Msgf("%s new transaction failed, err=", name, err)
		return errors.New("transaction handling failed")
	}

	if err = driver.RegisterTransaction(p, addr, transaction); err != nil {
		log.Error().Msgf("%s register transaction '%s' failed, err=%s", name, transaction.Identifier, err)
		return errors.New("transaction handling failed")
	}

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context, transaction)
	if err == nil {
		if err = p.automate.SendPacket(addr, packet, action); err == nil {
			return nil
		}
	}

	driver.UnregisterTransaction(p, addr, transaction)
	return err
}

func (p *Plugin) clearTransaction(addr string) {

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	if transaction, exist := driver.LastTransaction(p, addr); exist {
		if err := driver.UnregisterTransaction(p, addr, transaction); err != nil {
			log.Error().Msgf("%s unregister transaction '%s' failed, err=", name, transaction.Identifier, err)
		}
	}
}
