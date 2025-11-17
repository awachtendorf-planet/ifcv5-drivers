package linkcontrol

import (
	syncstate "github.com/weareplanet/ifcv5-drivers/bartech/state"
	deviceevent "github.com/weareplanet/ifcv5-main/device/event"
	devicestate "github.com/weareplanet/ifcv5-main/device/state"
	"github.com/weareplanet/ifcv5-main/log"

	// "github.com/weareplanet/ifcv5-drivers/bartech/template"

	"github.com/pkg/errors"
)

const (
	wildcardAddress = "*"
)

func (p *Plugin) main() {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	subscriber, err := dispatcher.RegisterEvents(
		deviceevent.Connected.String(),
		deviceevent.StateChanged.String(),
		deviceevent.Disconnected.String(),
		syncstate.Start.String(),
		syncstate.End.String(),
		syncstate.Error.String(),
	)
	if err != nil {
		log.Error().Msgf("%s", err)
		return
	}

	defer func() {
		dispatcher.DeregisterEvents(subscriber)
	}()

	// route, err := dispatcher.RegisterPacketRoute(name, wildcardAddress,
	// 	template.PacketSwapEnquiry,
	// )
	// if err != nil {
	// 	log.Error().Msgf("%s", err)
	// 	return
	// }

	dbRunning := make(map[string]bool)
	devState := make(map[string][]string)

	defer func() {
		dispatcher.DeregisterPacketRoute(name, wildcardAddress)
	}()

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

					p.setLinkUp(addr)

				case deviceevent.Disconnected:

					if dbRunning[addr] {
						p.triggerEvent(addr, syncstate.Cancel, errors.New("vendor disconnected"))
					}

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

					p.setLinkDown(addr)

				case deviceevent.StateChanged:

					devAddr := addr
					if connections, exist := devState[devAddr]; exist {
						for i := range connections {
							addr := connections[i]
							if dispatcher.GetLinkState(addr) == linkUp {
								state := devicestate.DeviceState(event.ID)
								if state == devicestate.Shutdown || state == devicestate.WaitForPeer || state == devicestate.Initial {
									if dbRunning[addr] {
										p.triggerEvent(addr, syncstate.Cancel, errors.New("vendor disconnected"))
									}
									p.setLinkDown(addr)
								}
							}
						}
						delete(devState, devAddr)
					}

				}

			case *syncstate.Event:

				addr := event.Addr

				switch event.State {

				case syncstate.Start:
					p.driver.SetSwapState(addr, true)
					dbRunning[addr] = true

				case syncstate.End, syncstate.Error:
					p.driver.SetSwapState(addr, false)
					dbRunning[addr] = false

				}

			}

			// case packet := <-route.C:

			// 	if packet == nil {
			// 		continue
			// 	}

			// 	if packet.Name == template.PacketSwapEnquiry {

			// 		addr := packet.Addr

			// 		if !dbRunning[addr] && p.driver.SwapRequest(addr) {

			// 			dbRunning[addr] = true
			// 			p.triggerEvent(addr, syncstate.Start, nil)

			// 		}

			// 	}

		}
	}

}

func (p *Plugin) setLinkUp(addr string) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	dispatcher.ChangeLinkState(addr, linkUp)
}

func (p *Plugin) setLinkDown(addr string) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	dispatcher.ChangeLinkState(addr, linkDown)
	dispatcher.SetLinkState(addr, initial)
	automate.ClearState(name, addr)
}

func (p *Plugin) triggerEvent(addr string, event syncstate.State, err error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	broker := dispatcher.Broker()

	if broker != nil {
		broker.Broadcast(syncstate.NewEvent(addr, event, err), event.String())
	}
}
