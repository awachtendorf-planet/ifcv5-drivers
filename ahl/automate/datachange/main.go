package datachange

import (
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/ahl/template"
)

func (p *Plugin) main() {

	for {

		select {

		// shutdown
		case <-p.kill:
			return

		// action
		case order := <-p.job:
			p.handleJob(order.Addr, order.Job)

		}

	}

}

func (p *Plugin) finaliseJob(addr string, job *order.Job, reason string) {
	if job == nil {
		return
	}
	automate := p.automate
	dispatcher := automate.Dispatcher()
	automate.LogFinaliseJob(addr, job, "")
	dispatcher.RemoveJob(job)
}

func (p *Plugin) linkState(addr string) bool {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	linkState := dispatcher.GetLinkState(addr)
	if linkState != automatestate.LinkUp {
		log.Info().Msgf("%s addr '%s' canceled, because of link state '%s'", name, addr, linkState.String())
		return false
	}
	return true
}

func (p *Plugin) handleJob(addr string, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// handler already running
	if p.state.Exist(addr) {
		log.Error().Msgf("%s addr '%s' handler still exist", name, addr)
		p.handleError(addr, job, "handler still busy")
		p.finaliseJob(addr, job, "handler still exist")
		return
	}

	// link state
	if !p.linkState(addr) {
		p.handleError(addr, job, "vendor disconnected")
		p.finaliseJob(addr, job, "link state")
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	// run state maschine
	go func(addr string, job *order.Job) {

		automate := p.automate
		name := automate.Name()

		automate.RegisterRule(template.PacketReply, answer, p.handleReply, dispatcher.StateAction{})
		automate.RegisterRule(template.PacketReply, commandSent, p.handleReply, dispatcher.StateAction{})

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
			template.PacketReply,
		)

		p.state.Remove(addr)

		if err != nil {
			log.Error().Msgf("%s addr '%s' canceled, err=%s", name, addr, err)
			p.finaliseJob(addr, job, err.Error())
		} else {
			p.finaliseJob(addr, job, "")
		}

		p.waitGroup.Done()
	}(addr, job)

}
