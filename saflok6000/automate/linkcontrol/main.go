package linkcontrol

import (
	deviceevent "github.com/weareplanet/ifcv5-main/device/event"
	devicestate "github.com/weareplanet/ifcv5-main/device/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"
)

func (p *Plugin) main() {

	automate := p.automate
	dispatcher := automate.Dispatcher()

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
				automate.LogDeviceEvent(event, sub.GetRuntime())

				addr := event.Addr

				switch event.Type {

				case deviceevent.Connected:
					p.handleConnection(addr)

				}

			}

		}
	}

}

func (p *Plugin) handleConnection(addr string) {

	automate := p.automate
	name := automate.Name()

	if p.state.Exist(addr) {
		log.Warn().Msgf("%s addr '%s' handler still exist", name, addr)
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	go func(addr string) {

		automate := p.automate
		dispatcher := automate.Dispatcher()
		name := automate.Name()

		route, err := dispatcher.RegisterPacketRoute(name, addr, template.PacketAlive, template.PacketBeaconRequest)

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
		dispatcher.SetLinkState(addr, linkDown) // change linkstate from initial to linkdown
		dispatcher.SetAlive(addr)

		t := ticker.New(retryDelay, false)

		defer func() {
			t.Stop()
			dispatcher.DeregisterEvents(subscriber)
			dispatcher.DeregisterPacketRoute(name, addr)
			// broadcast linkdown state
			dispatcher.ChangeLinkState(addr, linkDown)
			// reset automate state
			dispatcher.SetLinkState(addr, initial)
			automate.ClearState(name, addr)
			// cleanup
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

}
