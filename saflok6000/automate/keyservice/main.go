package keyservice

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	keyLifetime          = 120 * time.Second
	virtualEncoderNumber = 254
)

func (p *Plugin) main() {

	for {

		select {

		// shutdown
		case <-p.kill:
			return

		// action
		case order := <-p.job:
			if guest, ok := order.Job.Context.(*record.Guest); ok {
				encoder := 0
				if coder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
					encoder = cast.ToInt(coder)
					// encoder mapped to PMS Terminal (001-254)
					if encoder == 0 {
						encoder = virtualEncoderNumber
						guest.SetGeneric(defines.EncoderNumber, encoder)
					}
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

func (p *Plugin) preCheck(addr string, job *order.Job) error {

	if job == nil {
		return errors.New("job nil object")
	}

	if job.Timestamp > 0 {
		timestamp := time.Unix(0, job.Timestamp)
		since := time.Since(timestamp)
		if since > keyLifetime {
			return errors.Errorf("key service deadline exceeded (%s)", since)
		}
	}

	// pre checks
	switch job.Action {

	case order.KeyRequest, order.KeyDelete:
		// OK

	default:
		return errors.Errorf("key service '%s' not supported", job.Action.String())

	}

	return nil
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

func (p *Plugin) isSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Reservation.SharedInd
	}
	return false
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

	// CO and Sharer -> ignore driver packet, send pms ok
	if job != nil && job.Action == order.KeyDelete {

		data := automate.GetSetting(addr, "KeyDelete", true)
		if state, ok := data.(bool); ok && state == false {
			log.Debug().Msgf("%s addr '%s' ignore key delete, because of config parameter", name, addr)
			keyAnswer := record.KeyAnswer{
				Success: true,
			}
			p.handleAnswer(addr, job, keyAnswer)
			p.finaliseJob(addr, job, "")
			return
		}

		if p.isSharer(job.Context) {
			log.Debug().Msgf("%s addr '%s' ignore key delete, because of sharer flag is set", name, addr)
			keyAnswer := record.KeyAnswer{
				Success: true,
			}
			p.handleAnswer(addr, job, keyAnswer)
			p.finaliseJob(addr, job, "")
			return
		}
	}

	// requested records
	if err := p.preCheck(addr, job); err != nil {
		reason := err.Error()
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

			process := func(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
				err := p.handleKeyAnswer(addr, in, action, job)
				return err
			}

			driver := p.driver
			encoder := p.getEncoder(addr)
			encoderAddr, _ := driver.GetTerminal(in)

			if encoder == encoderAddr {
				err := process(addr, in, action, job)
				return err
			} else {
				log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' for encoder: %d, expect: %d", name, addr, in.Name, encoderAddr, encoder)
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

		automate.RegisterRule(template.PacketSuccess, answer, handleKeyAnswer, dispatcher.StateAction{})
		automate.RegisterRule(template.PacketError, answer, handleKeyAnswer, dispatcher.StateAction{})

		err := automate.StateMaschine(
			addr,
			job,
			nil,
			nil,
			p.kill,
			template.PacketSuccess,
			template.PacketError,
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
