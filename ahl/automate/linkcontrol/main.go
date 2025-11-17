package linkcontrol

import (
	deviceevent "github.com/weareplanet/ifcv5-main/device/event"
	devicestate "github.com/weareplanet/ifcv5-main/device/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/ahl/template"
)

func (p *Plugin) main() {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	subscriber, err := dispatcher.RegisterEvents(
		deviceevent.Connected.String(),
		deviceevent.StateChanged.String(),
		deviceevent.Disconnected.String(),
	)
	if err != nil {
		log.Error().Msgf("%s", err)
		return
	}

	defer func() {
		dispatcher.DeregisterEvents(subscriber)
	}()

	devState := make(map[string][]string)

	for {
		select {

		// shutdown
		case <-p.kill:
			return

		// event broker
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

				isSerialLayer := p.isSerialDevice(addr)
				isSocketLayer := !isSerialLayer

				switch event.Type {

				case deviceevent.Connected:

					if isSocketLayer {

						p.handleConnection(addr)

					} else if isSerialLayer {

						alreadyExist := false
						devAddr, _ := dispatcher.GetCustomerAddr(addr)
						if connections, exist := devState[devAddr]; exist {
							for i := range connections {
								if connections[i] == addr {
									alreadyExist = true
									break
								}
							}
						}
						if !alreadyExist {
							devState[devAddr] = append(devState[devAddr], addr)
						}

						dispatcher.ChangeLinkState(addr, linkUp)

					}

				case deviceevent.Disconnected:

					if isSerialLayer {

						devAddr, _ := dispatcher.GetCustomerAddr(addr)
						if connections, exist := devState[devAddr]; exist {
							for i := range connections {
								if connections[i] == addr {
									connections = append(connections[:i], connections[i+1:]...)
									break
								}
							}
							devState[devAddr] = connections
						}

						dispatcher.ChangeLinkState(addr, linkDown)
						dispatcher.SetLinkState(addr, initial)
						automate.ClearState(name, addr)

					}

				case deviceevent.StateChanged:

					if isSerialLayer {

						devAddr := addr
						if connections, exist := devState[devAddr]; exist {
							for i := range connections {
								addr := connections[i]
								if dispatcher.GetLinkState(addr) == linkUp {
									state := devicestate.DeviceState(event.ID)
									if state == devicestate.Shutdown || state == devicestate.WaitForPeer || state == devicestate.Initial {
										dispatcher.ChangeLinkState(addr, linkDown)
										dispatcher.SetLinkState(addr, initial)
										automate.ClearState(name, addr)
									}
								}
							}
							delete(devState, devAddr)
						}

					}

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

		route, err := dispatcher.RegisterPacketRoute(name, addr, template.PacketLinkAlive)

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

				// execute state rule
				err = rule.Execute(packet, nil)

				// set next timeout
				automate.SetNextTimeoutSave(packet.Addr, rule.Action, err, t)

				// put rule to pool
				automate.PutRule(rule)
			}

		}

	}(addr)

}

func (p *Plugin) isSerialDevice(addr string) bool {

	// normalise addr
	// device state changed addr -> pcon-local:1234567890:1
	// device connect/disconnect addr -> pcon-local:1234567890:1:COM1
	// in case of device state addr create a valid addr e.g. pcon-local:1234567890:1:0

	automate := p.automate
	dispatcher := automate.Dispatcher()
	addr = ifc.NormaliseAddress(addr, 4)
	isSerialLayer := dispatcher.IsSerialDevice(addr)
	return isSerialLayer
}
