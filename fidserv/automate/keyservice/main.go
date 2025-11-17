package keyservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

const (
	keyLifetime = 120 * time.Second
)

func (p *Plugin) main() {

	for {

		select {

		// shutdown
		case <-p.kill:
			return

		// action
		case order := <-p.job:
			encoder := 0
			if guest, ok := order.Job.Context.(*record.Guest); ok {
				if coder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
					encoder = cast.ToInt(coder)
				}
				p.handleJob(order.Addr, encoder, order.Job)
			} else {
				errStr := fmt.Sprintf("job context is not a guest record (%T)", order.Job.Context)
				p.finaliseJob(order.Addr, order.Job, errStr)
			}

		}

	}

}

func (p *Plugin) handleAnswer(addr string, job *order.Job, keyAnswer record.KeyAnswer) {
	if job == nil || job.IsDone() {
		return
	}
	job.Done()

	automate := p.automate
	dispatcher := automate.Dispatcher()

	keyAnswer.Station = job.Station
	dispatcher.CreatePmsJob(addr, job, keyAnswer)
}

func (p *Plugin) handleFailed(addr string, job *order.Job, reason string) {
	if job == nil || job.IsDone() {
		return
	}
	job.Done()

	automate := p.automate
	dispatcher := automate.Dispatcher()

	keyAnswer := record.KeyAnswer{
		Success:       false,
		Station:       job.Station,
		Message:       reason,
		ResponseCode:  0,
		EncoderNumber: p.getEncoder(addr),
	}

	dispatcher.CreatePmsJob(addr, job, keyAnswer)
}

func (p *Plugin) cancelJob(addr string, job *order.Job, reason string) {
	if job == nil {
		return
	}
	if !job.IsDone() {
		p.handleFailed(addr, job, reason)
	}
	p.finaliseJob(addr, job, reason)
}

func (p *Plugin) finaliseJob(addr string, job *order.Job, reason string) {
	if job == nil || job.IsRemoved() {
		return
	}
	automate := p.automate
	dispatcher := automate.Dispatcher()

	automate.LogFinaliseJob(addr, job, reason)
	dispatcher.RemoveJob(job)
}

func (p *Plugin) linkState(addr string) bool {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	// normalise device address
	deviceAddr := ifc.DeviceAddress(addr)
	linkState := dispatcher.GetLinkState(deviceAddr)

	if linkState != automatestate.LinkUp {
		log.Info().Msgf("%s addr '%s' canceled, because of link state '%s'", name, addr, linkState.String())
		return false
	}
	return true
}

func (p *Plugin) preCheck(addr string, job *order.Job) (string, bool) {

	if job == nil {
		return "job nil object", false
	}

	if job.Timestamp > 0 {
		timestamp := time.Unix(0, job.Timestamp)
		since := time.Since(timestamp)
		if since > keyLifetime {
			return fmt.Sprintf("key service deadline exceeded (%s)", since), false
		}
	}

	driver := p.driver
	switch job.Action {

	case order.KeyRequest:
		if !driver.IsRequested(addr, "KR") {
			return "KR not requested", false
		}

	case order.KeyDelete:
		if !driver.IsRequested(addr, "KD") {
			return "KD not requested", false
		}

	case order.KeyChange:
		if !driver.IsRequested(addr, "KM") {
			return "KM not requested", false
		}

	case order.KeyRead:
		if !driver.IsRequested(addr, "KZ") {
			return "KZ not requested", false
		}

	default:
		return "unkown action " + job.Action.String(), false
	}
	return "", true
}

func (p *Plugin) getEncoder(addr string) int {
	s := strings.Split(addr, ":")
	if len(s) < 5 {
		return -1
	}
	encoder := cast.ToInt(s[4])
	return encoder
}

func (p *Plugin) getKeyAnswerTimeout(addr string) time.Duration {
	automate := p.automate
	if timeout, ok := automate.GetSetting(addr, defines.KeyAnswerTimeout, uint(0)).(uint64); ok && timeout > 0 {
		if duration := cast.ToDuration(time.Duration(timeout) * time.Second); duration > 0 {
			return duration
		}
	}

	return keyAnswerTimeout
}

func (p *Plugin) handleJob(addr string, encoder int, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// addr:encoder
	addr = addr + ":" + cast.ToString(encoder)

	// handler already running
	if p.state.Exist(addr) {
		log.Error().Msgf("%s addr '%s' handler still exist", name, addr)
		p.cancelJob(addr, job, "handler still exist")
		return
	}

	// link state
	if !p.linkState(addr) {
		p.cancelJob(addr, job, "link state")
		return
	}

	// requested records
	if reason, success := p.preCheck(addr, job); !success {
		if job == nil {
			log.Info().Msgf("%s addr '%s' canceled, because of %s", name, addr, reason)
		} else {
			log.Info().Msgf("%s addr '%s' job-id: %d canceled, because of %s", name, addr, job.GetQueue().Id, reason)
		}
		p.cancelJob(addr, job, reason)
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	// run state maschine
	go func(addr string, job *order.Job) {

		automate := p.automate
		name := automate.Name()

		// automate pro encoder
		handleKeyAnswer := func(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
			driver := p.driver
			payload, _ := driver.GetPayload(in)
			log.Info().Msgf("%s addr '%s' received packet '%s' %s", name, addr, in.Name, payload)

			if coder, exist := driver.GetField(in, "KC"); exist {
				if (in.Name == template.PacketKeyRead && job.Action == order.KeyRead) ||
					(in.Name == template.PacketKeyAnswer && job.Action != order.KeyRead) {
					encoder := p.getEncoder(addr)
					if encoder == cast.ToInt(coder) {
						err := p.handleKeyAnswer(addr, in, action, job)
						return err
					}
					log.Debug().Msgf("%s addr '%s' encoder number: %s not processed by this handler (encoder: %d)", name, addr, coder, encoder)
				}
			} else {
				log.Debug().Msgf("%s addr '%s' no encoder number provided", name, addr)
			}

			// diff bilden vom letzten packet versand, und verbleibendes coder timeout neu setzen
			nextTimeout := p.getKeyAnswerTimeout(addr)
			if job.Timestamp != 0 {
				diff := time.Since(time.Unix(0, job.Timestamp))
				if diff < nextTimeout {
					nextTimeout = nextTimeout - diff
				}
			}

			action.NextTimeout = nextTimeout

			return nil

		}

		automate.RegisterRule(template.PacketKeyAnswer, answer, handleKeyAnswer, dispatcher.StateAction{})
		automate.RegisterRule(template.PacketKeyRead, answer, handleKeyAnswer, dispatcher.StateAction{})

		err := automate.StateMaschine(
			// device address
			addr,
			// Job
			job,
			// preProcess
			nil,
			// postProcess
			nil,
			// shutdown chan
			p.kill,
			// templates
			template.PacketKeyAnswer,
			template.PacketKeyRead,
		)

		p.state.Remove(addr)

		if err != nil {
			p.cancelJob(addr, job, err.Error())
		} else {
			dispatcher := automate.Dispatcher()
			dispatcher.SetAlive(addr)
			p.finaliseJob(addr, job, "")
		}

		p.waitGroup.Done()
	}(addr, job)

}
