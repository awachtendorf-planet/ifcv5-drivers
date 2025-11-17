package linkcontrol

import (
	deviceevent "github.com/weareplanet/ifcv5-main/device/event"
	devicestate "github.com/weareplanet/ifcv5-main/device/state"
	"github.com/weareplanet/ifcv5-main/log"
)

const (
	wildcardAddress = "*"
)

func (p *Plugin) main() {

	automate := p.automate
	dispatcher := automate.Dispatcher()

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
									p.setLinkDown(addr)
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

func (p *Plugin) setLinkUp(addr string) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	dispatcher.SetAlive(addr)
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
