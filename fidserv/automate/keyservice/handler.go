package keyservice

import (
	// "fmt"

	"time"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	if !p.linkState(addr) {
		automate.ChangeState(name, addr, shutdown)
		t.Tick()
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

		if job.Context == nil {
			err = defines.ErrorJobContext
			break
		}

		action.NextState = answer
		action.NextTimeout = p.getKeyAnswerTimeout(addr)

		switch job.Action {

		case order.KeyRequest:
			err = p.sendKeyRequest(addr, &action, job.Initiator, job.Context)

		case order.KeyDelete:
			err = p.sendKeyDelete(addr, &action, job.Initiator, job.Context)

		case order.KeyChange:
			err = p.sendKeyChange(addr, &action, job.Initiator, job.Context)

		case order.KeyRead:
			err = p.sendKeyRead(addr, &action, job.Initiator, job.Context)

		}

		if err != nil {
			p.handleFailed(addr, job, err.Error())
			automate.ChangeState(name, addr, shutdown)
			t.Tick()
			return
		}
		job.Timestamp = time.Now().UnixNano()

	case answer:
		log.Warn().Msgf("%s addr '%s' key service failed, key system did not answer", name, addr)
		p.handleFailed(addr, job, "key system did not answer")
		automate.ChangeState(name, addr, timeout)
		t.Tick()
		return

	case success, timeout, shutdown:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {

		log.Error().Msgf("%s %s", name, err)
		if err == defines.ErrorJobContext || err == defines.ErrorJobFailed {
			p.handleFailed(addr, job, err.Error())
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		if job != nil {
			job.Tries++
			if job.Tries >= maxError {
				p.handleFailed(addr, job, err.Error())
				automate.NextAction(name, addr, shutdown, t)
				return
			}
		}

		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handleKeyAnswer(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	keyAnswer := record.KeyAnswer{}
	if err := driver.UnmarshalPacket(packet, &keyAnswer); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	answer, _ := driver.GetField(packet, "AS")

	if len(keyAnswer.Message) == 0 {
		keyAnswer.Message = answer
	}

	if answer == "DE" && job != nil && job.Action == order.KeyDelete {
		answer = "OK"
	}

	switch answer {

	case "OK":

		if job != nil && job.Action == order.KeyRead {

			guest := &record.Guest{}
			if err := driver.UnmarshalPacket(packet, guest); err != nil {
				log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
			}

			departureDate, _ := driver.GetField(packet, "GD")
			departureTime, _ := driver.GetField(packet, "DT")

			// YYMMDDhhmm
			if t, err := time.Parse("0601021504", departureDate+departureTime); err == nil {
				guest.Reservation.DepartureDate = t
			}

			keyAnswer.Guest = guest

		}

		keyAnswer.Success = true
		p.handleAnswer(addr, job, keyAnswer)
		action.NextState = success
		action.NextTimeout = nextActionDelay

	default:

		keyAnswer.Success = false
		p.handleAnswer(addr, job, keyAnswer)
		action.NextState = shutdown
		action.NextTimeout = nextActionDelay
	}

	automate.ChangeState(name, addr, action.NextState)

	return nil
}

func (p *Plugin) logOutgoingPacket(packet *ifc.LogicalPacket) {
	name := p.automate.Name()
	p.driver.LogOutgoingPacket(name, packet)
}

func (p *Plugin) sendKeyRequest(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketKeyRequest, tracking, "KR", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyDelete(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketKeyDelete, tracking, "KD", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyChange(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketKeyChange, tracking, "KM", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyRead(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketKeyRead, tracking, "KZ", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
