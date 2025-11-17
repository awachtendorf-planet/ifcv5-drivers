package linkcontrol

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/ahl/template"
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

		automate.ChangeState(name, addr, idle)
		automate.ClearPendingAction(addr)

		station, _ := dispatcher.GetStationAddr(addr)
		if !dispatcher.IsDebugMode(station) {
			if errorCounter := automate.AddErrorCounter(addr); errorCounter >= maxError {
				if dispatcher.IsSerialDevice(addr) {
					dispatcher.ChangeLinkState(addr, linkDown)
				} else {
					p.closeConnection(addr)
					log.Warn().Msgf("%s addr '%s' close connection, err=%s", name, addr, "ahl did not answer")
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

		action.NextTimeout = aliveTimeout
		dispatcher := automate.Dispatcher()

		if dispatcher.IsSerialDevice(addr) {
			err = p.sendLinkAlive(addr, &action)
		} else {
			err = p.sendLinkStart(addr, &action)
		}

	case linkUp:

		action.NextTimeout = aliveTimeout
		err = p.sendLinkAlive(addr, &action)
	}

	if err != nil {

		if err == dispatcher.ErrPeerDisconnected || err == dispatcher.ErrDeviceDisconnected { // IFCDEV-65 unnötige Erweiterung
			// connection lost, handler shutdown
			automate.ChangeState(name, addr, shutdown)
			t.Tick()
			return
		}

		currentState := automate.GetState(addr)

		// prepare connection handler shutdown
		if currentState == shutdown {
			t.Tick()
			return
		}
		log.Error().Msgf("%s %s", name, err)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handlePacketAcknowledge(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	automate.ClearErrorCounter(addr)
	dispatcher.SetAlive(addr)

	automate.ChangeState(name, addr, idle)
	if pendingAction, exist := automate.GetPendingAction(addr); exist {
		automate.ClearPendingAction(addr)
		dispatcher.ChangeLinkState(addr, linkUp)
		action.NextTimeout = pendingAction.NextTimeout
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

func (p *Plugin) handleLinkAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	// re calculate next timeout, because we need every 30 second a link alive packet
	lastAlive := dispatcher.GetAlive(addr)
	if lastAlive > 0 {
		elapsed := time.Since(time.Unix(0, lastAlive))
		if elapsed < aliveTimeout {
			diff := aliveTimeout - elapsed
			action.NextTimeout = diff
		}
	}

	return nil
}

func (p *Plugin) handleError(addr string, job *order.Job, reason string) {
	p.closeConnection(addr)
}

func (p *Plugin) sendLinkStart(addr string, action *dispatcher.StateAction) error {
	packet, _ := p.driver.ConstructPacket(addr, template.PacketLinkStart, "", nil)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendLinkAlive(addr string, action *dispatcher.StateAction) error {
	packet, _ := p.driver.ConstructPacket(addr, template.PacketLinkAlive, "", nil)
	err := p.automate.SendPacket(addr, packet, action)
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

func (p *Plugin) closeConnection(addr string) {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := dispatcher.Network
	if driver != nil {
		driver.Disconnect(addr)
	}
}
