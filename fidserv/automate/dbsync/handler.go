package dbsync

import (
	// "math/rand"
	// "time"

	"github.com/pkg/errors"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
)

/*
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	// rand.Int63n(300) + 100)
}
*/

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, _ *order.Job) {

	automate := p.automate
	name := automate.Name()

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	if respectLinkState {
		linkState := automate.Dispatcher().GetLinkState(addr)
		if linkState != automatestate.LinkUp {
			log.Info().Msgf("%s addr '%s' canceled, because of link state '%s'", name, addr, linkState.String())
			automate.ChangeState(name, addr, success)
			t.Tick()
			return
		}
	}

	switch state {

	case resyncStart:
		// send DS
		action.NextState = resyncRecord
		err = p.sendDBResyncStart(addr, &action)

	case resyncRecord:
		// send guest record GI/GO
		action.NextState = resyncNextRecord
		err = p.sendRecord(addr, &action)

	case resyncNextRecord:
		// get next record
		// next state -> resyncRecord or resyncEnd if no more records
		action.NextState = resyncRecord
		err = p.getNextRecord(addr, &action)

		/*
			if err == nil && action.NextState == resyncRecord  {
				// if multiple peers are running a database swap then throttle things down, sleep for 100-300 ms
				if automate.Dispatcher().DatabaseSyncRunning() > x {
					time.Sleep(time.Duration(rand.Int63n(300)+100) * time.Millisecond)
				}
			}
		*/

	case resyncEnd:
		// send DE
		action.NextState = success
		err = p.sendDBResyncEnd(addr, &action)

	case busy:
		automate.ChangeState(name, addr, success)
		return

	case success:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) sendRecord(addr string, action *dispatcher.StateAction) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := p.driver
	name := automate.Name()

	job := dispatcher.GetSyncRecord(addr)
	if job != nil {

		if job.Action == order.Checkin && driver.IsRequested(addr, "GI") {
			err := p.sendCheckIn(addr, action, job.Initiator, job.Context)
			return err
		}

		if job.Action == order.Checkout && driver.IsRequested(addr, "GO") {
			err := p.sendCheckOut(addr, action, job.Initiator, job.Context)
			return err
		}

		log.Debug().Msgf("%s addr '%s' drop record '%s'", name, addr, job.Action.String())

		if action != nil {
			action.NextState = resyncNextRecord
			automate.ChangeState(name, addr, action.NextState)
		}
		return nil
	}

	if action != nil {
		action.NextState = resyncEnd
		automate.ChangeState(name, addr, action.NextState)
	}
	return nil
}

func (p *Plugin) getNextRecord(addr string, action *dispatcher.StateAction) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	if dispatcher.GetNextSyncRecord(addr) {
		if action != nil {
			automate.ChangeState(name, addr, action.NextState)
		}
		return nil
	}

	if action != nil {
		action.NextState = resyncEnd
		automate.ChangeState(name, addr, action.NextState)
	}
	return nil
}

func (p *Plugin) sendDBResyncStart(addr string, action *dispatcher.StateAction) error {
	packet := p.driver.ConstructPacket(addr, template.PacketResyncStart, "", "DS", nil)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendDBResyncEnd(addr string, action *dispatcher.StateAction) error {
	packet := p.driver.ConstructPacket(addr, template.PacketResyncEnd, "", "DE", nil)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendCheckIn(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	if context == nil {
		return errors.New("context object missing")
	}
	packet := p.driver.ConstructPacket(addr, template.PacketCheckInSwap, tracking, "GI", context)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendCheckOut(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	if context == nil {
		return errors.New("context object missing")
	}
	packet := p.driver.ConstructPacket(addr, template.PacketCheckOutSwap, tracking, "GO", context)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
