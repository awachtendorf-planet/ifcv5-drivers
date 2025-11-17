package datachange

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
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

	dispatcher := p.automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)
	accountNumberLength := p.driver.GetAccountNumberLength(station)

	switch job.Context.(type) {

	case *record.Guest:
		guest := job.Context.(*record.Guest)
		if len(guest.Reservation.RoomNumber) > 6 {
			return "roomnumber is too long", false
		} else if len(guest.Reservation.ReservationID) > accountNumberLength {
			return "reservation-id is too long", false
		}

	case *record.Generic:
		generic := job.Context.(*record.Generic)
		if data, exist := generic.Get(defines.WakeExtension); exist {
			extension := cast.ToString(data)
			if len(extension) > 6 {
				return "wakeup extension is too long", false
			}
		}
	}
	return "", true
}

func (p *Plugin) handleJob(addr string, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// handler already running
	if p.state.Exist(addr) {
		log.Error().Msgf("%s addr '%s' handler still exist", name, addr)
		p.handleError(addr, job, "handler still busy")
		p.finaliseJob(addr, job)
		return
	}

	// link state
	if !p.linkState(addr) {
		p.handleError(addr, job, "vendor disconnected")
		p.finaliseJob(addr, job)
		return
	}

	// CO and Sharer -> ignore driver packet, send pms ok
	if job != nil && job.Action == order.Checkout && p.isSharer(job.Context) {
		log.Debug().Msgf("%s addr '%s' ignore checkout, because of sharer flag is set", name, addr)
		p.handleAnswerPositiv(addr, job)
		p.finaliseJob(addr, job)
		return
	}

	// precheck
	if reason, success := p.preCheck(addr, job); !success {
		if job == nil {
			log.Info().Msgf("%s addr '%s' canceled, because of %s", name, addr, reason)
		} else {
			log.Info().Msgf("%s addr '%s' job-id: %d canceled, because of %s", name, addr, job.GetQueue().Id, reason)
		}
		p.handleError(addr, job, reason)
		p.finaliseJob(addr, job)
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	// run state maschine
	go func(addr string, job *order.Job) {

		automate := p.automate
		name := automate.Name()

		cleanup := func() {
			// driver cleanup
			p.driver.FreeTransaction(p, addr) // make sure that the router automate does not find/access the route anymore
			p.driver.UnregisterAndCloseRoute(p, addr)
		}

		err := automate.StateMaschine(
			// device address
			addr,
			// Job
			job,
			// preProcess
			nil,
			// postProcess
			cleanup,
			// shutdown chan
			p.kill,
			// templates
			// template.PacketVerify, // -> via router automate
			// template.PacketError, // -> via router automate
		)

		p.state.Remove(addr)

		if err != nil {
			log.Error().Msgf("%s addr '%s' canceled, err=%s", name, addr, err)
		}

		p.finaliseJob(addr, job)

		p.waitGroup.Done()

	}(addr, job)

}
