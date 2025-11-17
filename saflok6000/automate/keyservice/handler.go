package keyservice

import (
	"time"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"
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

	encoder := p.getEncoder(addr)
	if encoder == virtualEncoderNumber {
		encoder = 0
	}

	keyAnswer := record.KeyAnswer{
		EncoderNumber: encoder,
	}

	statusCode, _ := driver.GetStatusCode(packet)

	switch statusCode {

	case 0: // success

		keyAnswer.Success = true
		p.handleAnswer(addr, job, keyAnswer)

		action.NextState = success
		action.NextTimeout = nextActionDelay

	default:

		keyAnswer.Success = false

		if statusCode > 0 {

			txc := 0

			switch job.Action {

			case order.KeyRequest:
				txc = driver.GetKeyType(job.Context) // 1 = new key, 3 = duplicate key

			case order.KeyDelete:
				txc = 15

			}

			responseCode, _ := driver.GetResponseCode(packet)
			keyAnswer.Message = driver.GetAnswerText(statusCode, responseCode, txc)
			responseCode = (statusCode * 10000) + responseCode
			keyAnswer.ResponseCode = responseCode
		}

		p.handleAnswer(addr, job, keyAnswer)

		action.NextState = shutdown
		action.NextTimeout = nextActionDelay
	}

	automate.ChangeState(name, addr, action.NextState)

	return nil
}

func (p *Plugin) sendKeyRequest(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketKeyRequest, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendKeyDelete(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketKeyDelete, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}
