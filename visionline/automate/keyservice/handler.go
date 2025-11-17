package keyservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/visionline/template"
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
			err = p.sendKeyUpdate(addr, &action, job.Initiator, job.Context)

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
	dispatcher := automate.Dispatcher()

	name := automate.Name()

	// normalise address, remove :encoder
	// pcon-4711-local:47110815:1:192.168.0.205#65231:1 -> pcon-4711-local:47110815:1:192.168.0.205#65231
	if conAddr, err := dispatcher.NormalizePhysicalAddr(addr); err == nil {
		dispatcher.SetAlive(conAddr)
	}

	keyAnswer := record.KeyAnswer{}
	if err := driver.UnmarshalPacket(packet, &keyAnswer); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	resultCode, _ := driver.GetResultCode(packet)
	if resultCode != 0 && resultCode != 13 {
		if len(keyAnswer.Message) == 0 {
			keyAnswer.Message = fmt.Sprintf("visionline error code: %d", resultCode)
		}
	}

	switch resultCode {

	case 0, 13: // success, card is valid

		if job != nil && job.Action == order.KeyRead {

			guest := &record.Guest{}
			if err := driver.UnmarshalPacket(packet, guest); err != nil {
				log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
			}

			if date, exist := driver.GetField(packet, "CI"); exist {
				if t, err := time.Parse("200601021504", date); err == nil {
					guest.Reservation.ArrivalDate = t
				}
			}

			if date, exist := driver.GetField(packet, "CO"); exist {
				if t, err := time.Parse("200601021504", date); err == nil {
					guest.Reservation.DepartureDate = t
				}
			}

			// re-construct keyoptions CRGarage,Spa -> 003

			keyAnswer.KeyOptions = p.reconstructKeyOptions(keyAnswer.AdditionalRooms, job.Station)

			// normalize GR101-103,110
			// main room -> 101
			// additional rooms -> 102,103,110

			rooms := p.normalizeRoom(keyAnswer.RoomNumber)
			if len(rooms) > 0 {
				keyAnswer.RoomNumber = rooms[0]
			}

			var additionalRooms string
			if len(rooms) > 1 {
				additionalRooms = strings.Join(rooms[1:], ",")
			}
			keyAnswer.AdditionalRooms = additionalRooms

			keyAnswer.Guest = guest

		}

		if len(keyAnswer.Track2) == 0 && len(keyAnswer.SerialNumber) > 0 {
			keyAnswer.Track2 = keyAnswer.SerialNumber
			keyAnswer.SerialNumber = ""
		}

		if len(keyAnswer.Track2) > 0 {
			trim := automate.GetSetting(addr, "Track2Trim", false)
			if t, ok := trim.(bool); ok && t {
				keyAnswer.Track2 = p.trimTrack(keyAnswer.Track2)
			}
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
	packet, err := p.driver.ConstructPacket(addr, template.PacketCodeCard, tracking, context)
	if err != nil {
		return err
	}
	p.logOutgoingPacket(packet)
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyDelete(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketCheckout, tracking, context)
	if err != nil {
		return err
	}
	p.logOutgoingPacket(packet)
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyUpdate(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketCardUpdate, tracking, context)
	if err != nil {
		return err
	}
	p.logOutgoingPacket(packet)
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyRead(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketReadKey, tracking, context)
	if err != nil {
		return err
	}
	p.logOutgoingPacket(packet)
	err = p.automate.SendPacket(addr, packet, action)
	return err
}
