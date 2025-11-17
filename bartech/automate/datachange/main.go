package datachange

import (
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) main() {

	p.setWorkflow()

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

	// database swap running
	if p.driver.GetSwapState(addr) {
		p.handleError(addr, job, "database swap running")
		p.finaliseJob(addr, job, "database swap running")
		return
	}

	// pre checks

	// check sharer flag, ignore driver packet, send pms ok
	if job != nil && p.driver.IsSharer(job.Context) {
		log.Debug().Msgf("%s addr '%s' ignore packet, because of sharer flag is set", name, addr)
		p.handleSuccess(addr, job)
		p.finaliseJob(addr, job, "")
		return
	}

	// is workflow defined, job starts with 1
	if _, exist := p.workflow.Get(int(job.Action), 1); !exist {
		log.Error().Msgf("%s addr '%s' no workflow for '%s' defined", name, addr, job.Action)
		p.handleError(addr, job, "no workflow defined")
		p.finaliseJob(addr, job, "no workflow defined")
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	// run state maschine
	go func(addr string, job *order.Job) {

		automate := p.automate
		name := automate.Name()

		err := automate.StateMaschine(
			// device address
			addr,
			// Job
			job,
			// preProcess
			func() {
				automate.ChangeState(name, addr, nextAction)
			},
			// postProcess
			nil,
			// shutdown chan
			p.kill,
		)

		p.state.Remove(addr)

		if err != nil {
			log.Error().Msgf("%s addr '%s' canceled, err=%s", name, addr, err)
			p.handleError(addr, job, err.Error())
			p.finaliseJob(addr, job, err.Error())
		} else {
			p.finaliseJob(addr, job, "")
		}

		p.waitGroup.Done()
	}(addr, job)

}
