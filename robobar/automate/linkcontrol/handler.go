package linkcontrol

import (
	"time"

	syncstate "github.com/weareplanet/ifcv5-drivers/robobar/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/robobar/template"
)

func (p *Plugin) handleTimeout(addr string, state automatestate.State, t *ticker.ResetTicker) {

	automate := p.automate
	name := automate.Name()

	var err error
	action := dispatcher.StateAction{NextTimeout: retryDelay}

	switch state {

	case packetSent, commandSent:

		dispatcher := automate.Dispatcher()
		dispatcher.UnsetAcknowledgement(name, addr)

		automate.ChangeState(name, addr, linkDown)
		dispatcher.ChangeLinkState(addr, linkDown)
		automate.ClearPendingAction(addr)

		station, _ := dispatcher.GetStationAddr(addr)
		if !dispatcher.IsDebugMode(station) {
			if errorCounter := automate.AddErrorCounter(addr); errorCounter >= maxError {
				if dispatcher.IsSerialDevice(addr) {
					dispatcher.ChangeLinkState(addr, linkDown)
				} else {
					p.closeConnection(addr)
					log.Warn().Msgf("%s addr '%s' close connection, err=%s", name, addr, "robobar did not answer")
					break
				}
			}
		}

		if dispatcher.IsAcknowledgement(addr, "") {
			if packetTimeout >= retryDelay { // if protocol acknowledgement timeout > retryTimeout then retry immediately
				action.NextTimeout = nextActionDelay
			} else { // correct next timeout
				action.NextTimeout = (retryDelay - packetTimeout)
			}
		}

	case busy, linkDown:

		last := p.getLastRestartMessage(addr)
		elapsed := time.Since(last)

		// we run into the state if db swap is triggered too often within a period of time (dispatcher.isDatabaseSyncAllowed)
		if elapsed < 3*time.Second {
			log.Debug().Msgf("%s addr '%s' restart message delayed, because sent %s ago", name, addr, elapsed)
			break
		}

		action.NextTimeout = packetTimeout
		err = p.sendRestart(addr, &action)

		if err == nil {
			p.setLastRestartMessage(addr)
		}

	case linkUp:

		action.NextTimeout = aliveTimeout

	case resyncStart:

		// make sure we first t.Tick the automate
		defer func() {
			dispatcher := automate.Dispatcher()
			broker := dispatcher.Broker()
			if broker != nil {
				broker.Broadcast(syncstate.NewEvent(addr, syncstate.Start, nil), syncstate.Start.String())
			}
		}()

		automate.ChangeState(name, addr, resyncRecord)
		t.Tick()
		return

	case resyncRecord:

		// waiting for dbsync trigger
		action.NextTimeout = pmsSyncTimeout
		log.Debug().Msgf("%s addr '%s' pending sync records, waiting for dbsync automate", name, addr)

	}

	if err != nil {

		if err == dispatcher.ErrPeerDisconnected || err == dispatcher.ErrDeviceDisconnected {
			// connection lost, handler shutdown
			automate.ChangeState(name, addr, shutdown)
			t.Tick()
			return
		}

		currentState := automate.GetState(addr)
		if currentState == shutdown { // prepare connection handler shutdown
			t.Tick()
			return
		}

		log.Error().Msgf("%s %s", name, err)

	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handlePacketAcknowledge(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	name := automate.Name()

	automate.ClearErrorCounter(addr)

	automate.ChangeState(name, addr, idle)
	if pendingAction, exist := automate.GetPendingAction(addr); exist {
		automate.ClearPendingAction(addr)
		action.NextTimeout = pendingAction.NextTimeout

		dispatcher := p.automate.Dispatcher()
		dispatcher.ChangeLinkState(addr, pendingAction.NextState)

	} else {
		action.NextTimeout = retryDelay
	}

	return nil
}

func (p *Plugin) handlePacketRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	name := automate.Name()

	if pendingAction, exist := automate.GetPendingAction(addr); exist {
		automate.ClearPendingAction(addr)
		automate.ChangeState(name, addr, pendingAction.CurrentState) // current state
		action.NextTimeout = retryDelay
	} else {
		automate.ChangeState(name, addr, busy)
		action.NextTimeout = retryDelay
	}
	return nil
}

func (p *Plugin) sendRestart(addr string, action *dispatcher.StateAction) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketRestart, "", nil)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendStartup(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	dispatcher := p.automate.Dispatcher()
	dispatcher.ChangeLinkState(addr, linkDown)

	packet, err := p.driver.ConstructPacket(addr, template.PacketStartup, "", nil)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) getState(addr string) automatestate.State {
	state := p.automate.GetState(addr)
	if state == idle {
		dispatcher := p.automate.Dispatcher()
		state = dispatcher.GetLinkState(addr)
	}
	return state
}

func (p *Plugin) handleError(addr string, job *order.Job, reason string) {
	p.closeConnection(addr)
}

func (p *Plugin) closeConnection(addr string) {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := dispatcher.Network
	if driver != nil {
		driver.Disconnect(addr)
	}
}
