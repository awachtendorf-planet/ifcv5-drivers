package keyservice

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/vc3000/template"
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

	case busy, nextAction:

		if job == nil || !job.InProcess() { // order was canceled externally
			err = defines.ErrorJobFailed
			break
		}

		if job.Context == nil {
			err = defines.ErrorJobContext
			break
		}

		dispatcher := p.automate.Dispatcher()
		isSerialLayer := dispatcher.IsSerialDevice(addr)

		if isSerialLayer {
			if state == busy {
				action.NextState = nextAction
				err = p.sendEnquiry(addr, &action, "", nil)
				break
			}

			action.CurrentState = busy // change back to busy state on error to start again with enquiry
		}

		action.NextState = answer
		action.NextTimeout = p.getKeyAnswerTimeout(addr)

		// Serial Layer:
		// Wenn das Enquiry Packet mit ACK beantwortet wird und das eigentliche Payload Packet mit NAK
		// dann würde der Automate nie aus dem Versandversuch raus kommen, da der ErrorCounter mit dem ACK für
		// das Enquiy Packet resetet wird.
		// Wir benutzen den job.Tries Zähler um den ErrorCounter "künstlich" zu setzen

		if isSerialLayer {
			p.automate.SetErrorCounter(addr, uint16(job.Tries))
		}

		switch job.Action {

		case order.KeyRequest:
			err = p.sendKeyRequest(addr, &action, job.Initiator, job.Context, isSerialLayer)

		case order.KeyDelete:
			err = p.sendCheckout(addr, &action, job.Initiator, job.Context, isSerialLayer)

		case order.KeyChange:
			err = p.sendKeyRequestModify(addr, &action, job.Initiator, job.Context, isSerialLayer)

		case order.KeyRead:
			err = p.sendReadKey(addr, &action, job.Initiator, job.Context, isSerialLayer)

		}

		if err != nil {
			p.handleFailed(addr, job, err.Error())
			automate.ChangeState(name, addr, shutdown)
			t.Tick()
			return
		}

		job.Tries++
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

	var answerCode byte
	keyAnswer := record.KeyAnswer{
		EncoderNumber: p.getEncoder(addr),
		Station:       job.Station,
	}

	if job.Action == order.KeyRead || job.Action == order.KeyRequest || job.Action == order.KeyChange {

		if err := driver.UnmarshalPacket(packet, &keyAnswer); err != nil {
			log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
		}

		if job.Action == order.KeyRead {

			guest := &record.Guest{}
			if err := driver.UnmarshalPacket(packet, guest); err != nil {
				log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
			}

			if date, exist := driver.GetPayLoadField(packet, "D"); exist {
				if t, err := time.Parse("200601021504", date); err == nil {
					guest.Reservation.ArrivalDate = t
				}
			}

			if date, exist := driver.GetPayLoadField(packet, "O"); exist {
				if t, err := time.Parse("200601021504", date); err == nil {
					guest.Reservation.DepartureDate = t
				}
			}

			// freaky R/L logic
			if len(keyAnswer.RoomNumber) == 0 && len(keyAnswer.AdditionalRooms) > 0 {
				rooms := strings.Split(keyAnswer.AdditionalRooms, ",")
				if len(rooms) > 0 {
					keyAnswer.RoomNumber = rooms[0]
				}
				if len(rooms) > 1 {
					keyAnswer.AdditionalRooms = strings.Join(rooms[1:], ",")
				}
			}

			// radisson usertype and usergroup from key options and vice versa
			keyAnswer.KeyOptions = p.reconstructKeyOptions(packet, job.Station)

			keyAnswer.Guest = guest

		}

		if len(keyAnswer.Track2) == 0 && len(keyAnswer.SerialNumber) > 0 {
			keyAnswer.Track2 = keyAnswer.SerialNumber
			keyAnswer.SerialNumber = ""
		}

	}

	if answer, exist := driver.GetField(packet, "FF"); exist && len(answer) > 0 {
		answerCode = answer[0]
	}

	keyAnswer.Message = driver.GetAnswerText(answerCode)

	switch answerCode {

	case '0':

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

func (p *Plugin) sendEnquiry(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketEnq, true, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyRequest(addr string, action *dispatcher.StateAction, tracking string, context interface{}, isSerialLayer bool) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketCodeCard, isSerialLayer, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyRequestModify(addr string, action *dispatcher.StateAction, tracking string, context interface{}, isSerialLayer bool) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketCodeCardModify, isSerialLayer, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendCheckout(addr string, action *dispatcher.StateAction, tracking string, context interface{}, isSerialLayer bool) error {
	packetName := template.PacketCheckout
	if !isSerialLayer && p.driver.UseRmtCommandForKeyDelete(addr) {
		packetName = template.PacketCheckoutRmt
	}
	packet, err := p.driver.ConstructPacket(addr, packetName, isSerialLayer, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendReadKey(addr string, action *dispatcher.StateAction, tracking string, context interface{}, isSerialLayer bool) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketReadKey, isSerialLayer, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}
