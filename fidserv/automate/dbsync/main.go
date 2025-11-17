package dbsync

import (
	"time"

	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	databasesyncstate "github.com/weareplanet/ifcv5-main/ifc/generic/state/databasesync"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"

	"github.com/spf13/cast"
)

const (
	wildcardAddress  = "*"
	respectLinkState = true
)

func (p *Plugin) main() error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	route, err := dispatcher.RegisterPacketRoute(name, wildcardAddress,
		template.PacketResyncRequest,
	)
	if err != nil {
		return err
	}

	defer func() {
		dispatcher.DeregisterPacketRoute(name, wildcardAddress)
	}()

	for {

		select {

		// shutdown
		case <-p.kill:
			return nil

		// incoming logical packet
		case packet := <-route.C:
			if packet == nil {
				return nil
			}

			if packet.Name != template.PacketResyncRequest {
				continue
			}

			dispatcher.SetAlive(packet.Addr)
			p.handleRequest(packet.Addr)

		}

	}

	return nil
}

func (p *Plugin) handleRequest(addr string) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	if p.state.Exist(addr) {
		log.Warn().Msgf("%s addr '%s' connection handler still exist", name, addr)
		return
	}

	if respectLinkState {
		linkState := dispatcher.GetLinkState(addr)
		if linkState != automatestate.LinkUp {
			log.Info().Msgf("%s addr '%s' canceled, because of link state '%s'", name, addr, linkState.String())
			return
		}
	}

	driver := p.driver
	if !driver.IsRequested(addr, "GI") && !driver.IsRequested(addr, "GO") {
		log.Info().Msgf("%s addr '%s' done, because no GI or GO requested", name, addr)
		return
	}

	p.state.Register(addr)

	go func(addr string) {

		defer p.state.Remove(addr)

		automate := p.automate
		dispatcher := automate.Dispatcher()
		name := automate.Name()

		log.Info().Msgf("%s addr '%s' prepare database swap", name, addr)

		customerAddr, err := dispatcher.GetCustomerAddr(addr)
		if err != nil {
			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
			return
		}

		station, err := dispatcher.GetStationAddr(addr)
		if err != nil {
			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
			return
		}

		// prepare cache
		if err = dispatcher.DatabaseSyncPrepare(addr); err != nil {
			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
			return
		}

		// register ready event
		subscriber, err := dispatcher.RegisterEvents(
			databasesyncstate.DatabaseSyncEvent,
		)
		if err != nil {
			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
			return
		}

		// cleanup
		defer func() {
			dispatcher.DeregisterEvents(subscriber)
			dispatcher.DatabaseSyncCleanup(addr)
		}()

		// request databasesync from pms

		requestedMessages := cast.ToInt(driver.IsRequested(addr, "GI"))
		if driver.IsRequested(addr, "GO") {
			requestedMessages += 2
		}

		databaseSync := record.DatabaseSync{
			Station:           station,
			RequestedMessages: requestedMessages,
		}

		if _, pmsErr, sendErr := automate.PmsRequest(station, databaseSync, pmsTimeOut, "", ""); pmsErr != nil || sendErr != nil {
			if pmsErr == nil {
				pmsErr = sendErr
			}

			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, pmsErr)
			return
		}

		start := time.Now()

		// pre/post trigger
		triggerEvent := func(event automatestate.State) {
			broker := dispatcher.Broker()
			if broker != nil {
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

				switch event.(type) {

				case *databasesyncstate.Event:
					event := event.(*databasesyncstate.Event)
					if event.Station == station { // databasesync was fulfilled by pms
						ready = true
					}
				}

			case <-p.kill: // shutdown
				return

			case <-time.After(timeout): // timeout
				log.Error().Msgf("%s addr '%s' timeout, err=%s", name, addr, "pms did not response")
				return

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
				automate.ChangeState(name, addr, resyncStart)
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

	}(addr)

}
