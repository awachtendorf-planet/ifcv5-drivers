package datachange

import (
	"time"

	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
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

func (p *Plugin) cancelJob(addr string, job *order.Job, reason string) {
	if job == nil {
		return
	}
	automate := p.automate
	dispatcher := automate.Dispatcher()

	automate.LogFinaliseJob(addr, job, reason)
	dispatcher.PauseJob(job, reason)

	// possibly we already have a linkup again, resume job queue

	if dispatcher.HasStationLinkUp(job.Station) {

		go func(station uint64) {

			select {
			case <-p.kill: // shutdown
			case <-time.After(1 * time.Second): // resume job queue if linkup
				if dispatcher.HasStationLinkUp(station) {
					dispatcher.ResumeJobs(station)
				}
			}

		}(job.Station)

	}

}

func (p *Plugin) finaliseJob(addr string, job *order.Job) {
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

func (p *Plugin) preCheck(addr string, job *order.Job) (string, bool) {

	if job == nil {
		return "job nil object", false
	}

	driver := p.driver
	checked := true

	switch job.Action {

	case order.NightAuditStart:
		if !driver.IsRequested(addr, "NS") {
			return "NS not requested", false
		}

	case order.NightAuditEnd:
		if !driver.IsRequested(addr, "NE") {
			return "NE not requested", false
		}

	case order.WakeupRequest:
		if !driver.IsRequested(addr, "WR") {
			return "WR not requested", false
		}

	case order.WakeupClear:
		if !driver.IsRequested(addr, "WC") {
			return "WC not requested", false
		}

	case order.GuestMessageOnline:
		if !driver.IsRequested(addr, "XL") {
			return "XL not requested", false
		}

	case order.GuestMessageDelete:
		if !driver.IsRequested(addr, "XD") {
			return "XD not requested", false
		}

	default:
		checked = false

	}

	if checked {
		return "", true
	}

	if driver.IsRequested(addr, "RE") {
		return "", true
	}

	switch job.Action {

	case order.Checkin:
		if !driver.IsRequested(addr, "GI") {
			return "GI not requested", false
		}

	case order.Checkout:
		if !driver.IsRequested(addr, "GO") {
			return "GO not requested", false
		}

	case order.DataChange:
		if !driver.IsRequested(addr, "GC") {
			return "GC not requested", false
		}

	default:
		return "unkown action '" + job.Action.String() + "'", false
	}

	return "", true
}

func (p *Plugin) handleJob(addr string, job *order.Job) {

	automate := p.automate
	name := automate.Name()

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
		p.finaliseJob(addr, job)
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	// run state maschine
	go func(addr string, job *order.Job) {

		automate := p.automate
		dispatcher := automate.Dispatcher()
		name := automate.Name()

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
			// none for this automate
		)

		p.state.Remove(addr)

		if err != nil {
			log.Error().Msgf("%s addr '%s' canceled, err=%s", name, addr, err)
			p.cancelJob(addr, job, err.Error())
		} else {
			dispatcher.SetAlive(addr)
			p.finaliseJob(addr, job)
		}

		p.waitGroup.Done()
	}(addr, job)

}
