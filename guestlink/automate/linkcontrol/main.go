package linkcontrol

import (
	"time"

	deviceevent "github.com/weareplanet/ifcv5-main/device/event"
	devicestate "github.com/weareplanet/ifcv5-main/device/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
)

func (p *Plugin) main() {

	dispatcher := p.automate.Dispatcher()

	subscriber, err := dispatcher.RegisterEvents(
		deviceevent.Connected.String(),
	)
	if err != nil {
		log.Error().Msgf("%s", err)
		return
	}

	defer func() {
		dispatcher.DeregisterEvents(subscriber)
	}()

	pending := map[string]bool{}
	retry := make(chan string, 64)

	for {

		select {

		// shutdown
		case <-p.kill:
			return

		// event broker client connected
		case sub := <-subscriber.GetMessages():

			if sub == nil {
				continue
			}

			event := sub.GetPayload()

			switch event.(type) {

			case *deviceevent.Event:

				event := event.(*deviceevent.Event)
				p.automate.LogDeviceEvent(event, sub.GetRuntime())

				switch event.Type {

				case deviceevent.Connected:

					if !p.handleConnection(event.Addr) {

						// handler still exist, because handler is still finishing, or stays valid

						if pending[event.Addr] { // one retry is enough
							continue
						}

						pending[event.Addr] = true

						go func(addr string) {
							select {
							case <-p.kill:
							case <-time.After(2 * time.Second):
								retry <- addr // retry handleConnection
							}
						}(event.Addr)

					}

				}

			}

		case addr := <-retry:

			delete(pending, addr)

			// is the addr still valid, because we did not catch the device disconnect or peer state changed event here
			if dispatcher.Network.Exist(addr) {
				p.handleConnection(addr)
			}

		}

	}

}

func (p *Plugin) handleConnection(addr string) bool {

	automate := p.automate
	name := automate.Name()

	if p.state.Exist(addr) {
		log.Warn().Msgf("%s addr '%s' handler still exist", name, addr)
		return false
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	go func(addr string) {

		automate := p.automate
		dispatcher := automate.Dispatcher()
		name := automate.Name()

		route, err := dispatcher.RegisterPacketRoute(name, addr,
			//template.PacketVerify, // -> via router automate
			//template.PacketError, // -> via router automate
			template.PacketStart,
			template.PacketTest,
		)
		if err != nil {
			log.Error().Msgf("%s", err)
			p.state.Remove(addr)
			p.waitGroup.Done()
			return
		}

		subscriber, err := dispatcher.RegisterEvents(
			deviceevent.StateChanged.String(),
			deviceevent.Disconnected.String(),
		)
		if err != nil {
			log.Error().Msgf("%s", err)
			dispatcher.DeregisterPacketRoute(name, addr)
			p.state.Remove(addr)
			p.waitGroup.Done()
			return
		}

		automate.ChangeState(name, addr, idle)
		dispatcher.SetLinkState(addr, linkDown)
		dispatcher.SetAlive(addr)

		p.driver.RegisterRoute(p, addr, route)

		t := ticker.New(retryDelay, false)

		defer func() {
			t.Stop()
			// driver cleanup
			p.driver.UnregisterRoute(p, addr)
			p.driver.FreeTransaction(p, addr)
			//
			dispatcher.DeregisterEvents(subscriber)
			dispatcher.DeregisterPacketRoute(name, addr)
			// broadcast linkdown state
			dispatcher.ChangeLinkState(addr, linkDown)
			// reset automate state
			dispatcher.SetLinkState(addr, initial)
			automate.ClearState(name, addr)
			// cleanup
			p.clearAddress(addr)
			p.state.Remove(addr)
			p.waitGroup.Done()
		}()

		t.Tick()

		for {

			select {

			// shutdown
			case <-p.kill:
				return

			// automate timeout
			case tick := <-t.C:
				t.Pause()
				state := p.getState(addr)
				if state == success || state == shutdown {
					return
				}

				// intentional state transition, do not log
				if tick > 0 && t.Interval() != nextActionDelay {
					log.Info().Msgf("%s addr '%s' timeout at state '%s'", name, addr, state)
				}

				p.handleTimeout(addr, state, t)
				t.Resume()

			// event broker client disconnected, device state changed
			case sub := <-subscriber.GetMessages():
				if sub == nil {
					continue
				}
				event := sub.GetPayload()

				switch event.(type) {

				case *deviceevent.Event:
					event := event.(*deviceevent.Event)

					switch event.Type {

					case deviceevent.Disconnected:
						if event.Addr == addr {
							automate.LogDeviceEvent(event, sub.GetRuntime())
							automate.ChangeState(name, addr, shutdown)
							return
						}

					case deviceevent.StateChanged:
						if ifc.DeviceMask(addr) == ifc.DeviceMask(ifc.NormaliseAddress(event.Addr, 4)) {
							automate.LogDeviceEvent(event, sub.GetRuntime())
							dstate := devicestate.DeviceState(event.ID)
							if dstate == devicestate.Shutdown || dstate == devicestate.WaitForPeer || dstate == devicestate.Initial {
								automate.ChangeState(name, addr, shutdown)
								return
							}
						}
					}
				}

			// incoming logical packet
			case packet := <-route.C:
				if packet == nil {
					return
				}

				if packet.Addr != addr {
					continue
				}

				log.Info().Msgf("%s addr '%s' received packet '%s'", name, packet.Addr, packet.Name)

				state := p.getState(packet.Addr)

				rule, err := automate.GetRule(packet.Name, state)
				if err != nil {
					log.Error().Msgf("%s", err)
					continue
				}

				// set connetion active timestamp
				linkState := dispatcher.GetLinkState(packet.Addr)
				if linkState == linkUp {
					dispatcher.SetAlive(packet.Addr)
				}

				// execute state rule
				err = rule.Execute(packet, nil)

				// set next timeout
				automate.SetNextTimeout(packet.Addr, rule.Action, err, t)

				// put rule to pool
				automate.PutRule(rule)
			}

		}

	}(addr)

	return true

}
