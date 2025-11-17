package dbsync

import (
	"time"

	syncstate "github.com/weareplanet/ifcv5-drivers/callstar/state"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	databasesyncstate "github.com/weareplanet/ifcv5-main/ifc/generic/state/databasesync"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/pkg/errors"
)

func (p *Plugin) main() error {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	p.setWorkflow()

	subscriber, err := dispatcher.RegisterEvents(
		syncstate.Start.String(),
	)
	if err != nil {
		return err
	}

	defer func() {
		dispatcher.DeregisterEvents(subscriber)
	}()

	for {
		select {

		// shutdown
		case <-p.kill:
			return nil

		// event broker
		case sub := <-subscriber.GetMessages():
			if sub == nil {
				continue
			}
			event := sub.GetPayload()

			switch event := event.(type) {

			case *syncstate.Event:

				addr := event.Addr

				switch event.State {

				case syncstate.Start:
					p.handleRequest(addr)

				}

			}

		}
	}

}

func (p *Plugin) handleRequest(addr string) {

	automate := p.automate
	name := automate.Name()

	if p.state.Exist(addr) {
		log.Warn().Msgf("%s addr '%s' connection handler still exist", name, addr)
		return
	}

	p.state.Register(addr)

	do := func(addr string) error {

		defer p.state.Remove(addr)

		automate := p.automate
		dispatcher := automate.Dispatcher()
		name := automate.Name()

		log.Info().Msgf("%s addr '%s' prepare database swap", name, addr)

		customerAddr, err := dispatcher.GetCustomerAddr(addr)
		if err != nil {
			return err
		}

		station, err := dispatcher.GetStationAddr(addr)
		if err != nil {
			return err
		}

		// prepare cache
		if err = dispatcher.DatabaseSyncPrepare(addr); err != nil {
			return err
		}

		// register ready event
		subscriber, err := dispatcher.RegisterEvents(
			databasesyncstate.DatabaseSyncEvent,
			syncstate.Cancel.String(),
		)
		if err != nil {
			return err
		}

		// cleanup
		defer func() {
			dispatcher.DeregisterEvents(subscriber)
			dispatcher.DatabaseSyncCleanup(addr)
		}()

		// request databasesync from pms
		databaseSync := record.DatabaseSync{
			Station:           station,
			RequestedMessages: 15, // CI/CO/RE/WC
		}

		if _, pmsErr, sendErr := automate.PmsRequest(station, databaseSync, pmsTimeOut, "", ""); pmsErr != nil || sendErr != nil {
			if pmsErr == nil {
				pmsErr = sendErr
			}
			return pmsErr
		}

		start := time.Now()

		// pre/post trigger
		triggerEvent := func(event automatestate.State) {
			broker := dispatcher.Broker()
			if broker != nil {
				//broker.Broadcast(automatestate.NewEvent(addr, event), customerAddr, automatestate.ResyncStart.String())
				broker.Broadcast(automatestate.NewEvent(addr, event), customerAddr, event.String())
				utils.TimeTrack(start, "database swap")
			}
		}

		ready := false

		for {

			// recalculate the remaining time for each run
			offset := time.Since(start)
			timeout := pmsSyncTimeout - offset

			select {

			case sub := <-subscriber.GetMessages(): // event broker, ready event
				if sub == nil {
					continue
				}
				event := sub.GetPayload()

				switch event := event.(type) {

				case *databasesyncstate.Event:

					if event.Station == station { // databasesync was fulfilled by pms
						ready = true
					}

				case *syncstate.Event:

					switch event.State {

					case syncstate.Cancel:

						if event.Addr == addr {

							if event.Error != nil {
								return event.Error
							}

							return errors.New("swap request was aborted")

						}

					}

				}

			case <-p.kill: // shutdown
				return defines.ErrorIfcShutdown

			case <-time.After(timeout): // timeout
				return errors.New("pms did not response")

			}

			if ready {
				break
			}
		}

		err = automate.StateMaschine(
			// device address
			addr,
			// Job
			nil,
			// preProcess
			func() {
				// start with get first queued job
				automate.ChangeState(name, addr, nextRecord)
				// broadcast database swap start
				triggerEvent(automatestate.ResyncStart)
				start = time.Now()
			},
			// postProcess
			func() {
				// broadcast database swap end
				triggerEvent(automatestate.ResyncEnd)
			},
			// shutdown chan
			p.kill,
		)

		if err != nil {
			log.Error().Msgf("%s addr '%s' canceled, err=%s", name, addr, err)
		}

		return err

	}

	go func(addr string) {

		if err := do(addr); err != nil {

			name := automate.Name()
			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)

			// trigger err -> linkcontrol
			dispatcher := automate.Dispatcher()
			broker := dispatcher.Broker()
			if broker != nil {
				broker.Broadcast(syncstate.NewEvent(addr, syncstate.Error, err), syncstate.Error.String())
			}

		} else {

			// trigger success -> linkcontrol
			dispatcher := automate.Dispatcher()
			broker := dispatcher.Broker()
			if broker != nil {
				broker.Broadcast(syncstate.NewEvent(addr, syncstate.End, nil), syncstate.End.String())
			}

		}

	}(addr)
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
