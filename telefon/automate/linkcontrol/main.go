package linkcontrol

import (
	templatestate "github.com/weareplanet/ifcv5-drivers/telefon/state"
	deviceevent "github.com/weareplanet/ifcv5-main/device/event"
	devicestate "github.com/weareplanet/ifcv5-main/device/state"
	configstate "github.com/weareplanet/ifcv5-main/ifc/generic/state/config"
	"github.com/weareplanet/ifcv5-main/log"
)

const (
	wildcardAddress = "*"
)

func (p *Plugin) main() {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := p.driver
	name := automate.Name()

	subscriber, err := dispatcher.RegisterEvents(
		deviceevent.Connected.String(),
		deviceevent.StateChanged.String(),
		deviceevent.Disconnected.String(),
		templatestate.Changed.String(),
		configstate.ConfigEvent,
	)
	if err != nil {
		log.Error().Msgf("%s", err)
		return
	}

	p.waitGroup.Add(1)

	defer func() {
		dispatcher.DeregisterEvents(subscriber)
		dispatcher.DeregisterPacketRoute(name, wildcardAddress)
		p.waitGroup.Done()
	}()

	devState := make(map[string][]string) // connections (pcon-4711-local:47110815:1 -> pcon-4711-local:47110815:1:192.168.0.205#65231, ...)
	protocol := make(map[uint64]string)   // station -> protocol
	polling := make(map[string]chan bool) // running polling handler -> close channel

	// start polling if needed and not running
	startPolling := func(addr string) {

		if driver == nil {
			return
		}

		if state := driver.NeedPolling(addr); !state {
			return
		}

		if _, exist := polling[addr]; exist {
			return
		}

		kill := make(chan bool)
		polling[addr] = kill

		go p.polling(addr, kill)

	}

	// stop polling if running
	stopPolling := func(addr string) {

		if kill, exist := polling[addr]; exist {
			close(kill)
			delete(polling, addr)
		}
	}

	// reconfigure polling
	configure := func(station uint64, protocolName string) {

		if station == 0 || len(protocolName) == 0 {
			return
		}

		for devAddr, connections := range devState {
			if s, err := dispatcher.GetStationAddr(devAddr + ":0"); err == nil && s == station && driver != nil {
				for i := range connections {
					addr := connections[i]
					if state := driver.NeedPolling(addr); state {
						startPolling(addr)
					} else {
						stopPolling(addr)
					}
				}
			}

		}
	}

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

			switch event := event.(type) {

			case *deviceevent.Event:

				automate.LogDeviceEvent(event, sub.GetRuntime())

				addr := event.Addr

				switch event.Type {

				case deviceevent.Connected:

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
					//

					if driver != nil {
						if protocolName := driver.GetProtocol(addr); len(protocolName) > 0 {
							if station, err := dispatcher.GetStationAddr(addr); err == nil {
								protocol[station] = protocolName
							}
						}
					}

					p.setLinkUp(addr)
					startPolling(addr)

				case deviceevent.Disconnected:

					devAddr, _ := dispatcher.GetCustomerAddr(addr)
					if connections, exist := devState[devAddr]; exist {

						for i := range connections {
							if connections[i] == addr {
								connections = append(connections[:i], connections[i+1:]...)
								break
							}
						}
						devState[devAddr] = connections

						if len(connections) == 0 { // no more active connections
							if station, err := dispatcher.GetStationAddr(addr); err == nil {
								delete(protocol, station)
							}
						}

					}

					stopPolling(addr)
					p.setLinkDown(addr)

				case deviceevent.StateChanged:

					devAddr := addr
					if connections, exist := devState[devAddr]; exist {
						for i := range connections {

							addr := connections[i]

							if station, err := dispatcher.GetStationAddr(addr); err == nil {
								delete(protocol, station)
							}

							if dispatcher.GetLinkState(addr) == linkUp {
								state := devicestate.DeviceState(event.ID)
								if state == devicestate.Shutdown || state == devicestate.WaitForPeer || state == devicestate.Initial {
									stopPolling(addr)
									p.setLinkDown(addr)
								}
							}

						}

						delete(devState, devAddr)
					}

				}

			case *templatestate.Event:

				log.Debug().Msgf("- received event 'Template%s', slot: %d, event took: %s", event.State, event.Slot, sub.GetRuntime())

				slot := event.Slot

				if driver != nil {
					if protocolName := driver.GetProtocolFromSlot(slot); len(protocolName) > 0 {
						log.Info().Msgf("%s template slot: %d, protocol type '%s' changed", name, slot, protocolName)
						for k, v := range protocol {
							if v == protocolName {
								configure(k, protocolName)
							}
						}
					}
				}

			case *configstate.Event:

				if event.Station == 0 || event.Module != "config" {
					continue

				}

				log.Debug().Msgf("- received event '%s', module: '%s', station: %d, event took: %s", configstate.ConfigEvent, event.Module, event.Station, sub.GetRuntime())

				station := event.Station

				if currentProtocol, exist := protocol[station]; exist && driver != nil {
					protocolName := driver.GetProtocolByStation(station)
					if currentProtocol != protocolName {
						log.Info().Msgf("%s station: %d, protocol type '%s' changed to '%s'", name, station, currentProtocol, protocolName)
						protocol[station] = protocolName
					}
					configure(station, protocolName) // some station config parameters has been changed (maybe 'Protocol' or 'IgnorePolling')
				}

			}

		}
	}

}

func (p *Plugin) setLinkUp(addr string) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	dispatcher.ChangeLinkState(addr, linkUp)
	dispatcher.SetAlive(addr)
}

func (p *Plugin) setLinkDown(addr string) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	dispatcher.ChangeLinkState(addr, linkDown)
	dispatcher.SetLinkState(addr, initial)
	automate.ClearState(name, addr)
}
