package linkcontrol

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
)

func (p *Plugin) getState(addr string) automatestate.State {
	state := p.automate.GetState(addr)
	if state == idle {
		dispatcher := p.automate.Dispatcher()
		state = dispatcher.GetLinkState(addr)
	}
	return state
}

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

		linkState := dispatcher.GetLinkState(addr)

		if linkState == linkUp {

			lastAlive := dispatcher.GetAlive(addr)
			if lastAlive > 0 {
				elapsed := time.Since(time.Unix(0, lastAlive))
				if elapsed.Seconds() < 60 {
					break
				}
			}

			if state == packetSent {
				// ignore first send error because LS triggers an DB swap
				if failed := automate.AddErrorCounter(addr); failed > 1 {
					dispatcher.ChangeLinkState(addr, linkDown)
				}
			} else {
				dispatcher.ChangeLinkState(addr, linkDown)
			}

			// dispatcher.ChangeLinkState(addr, linkDown)

		}

		if dispatcher.IsAcknowledgement(addr, "") {
			if packetTimeout >= retryDelay { // if protocol acknowledgement timeout > retryTimeout then retry immediately
				action.NextTimeout = nextActionDelay
			} else { // correct next timeout
				action.NextTimeout = (retryDelay - packetTimeout)
			}
		}

	case linkDown:
		err = p.sendLinkStart(addr, nil, &action, nil)
		if err == nil {
			automate.ClearErrorCounter(addr)
		}

	case linkUp:
		action = dispatcher.StateAction{NextTimeout: aliveTimeout}
		dispatcher := automate.Dispatcher()

		lastAlive := dispatcher.GetAlive(addr)
		if lastAlive > 0 {
			elapsed := time.Since(time.Unix(0, lastAlive))
			if elapsed < aliveTimeout {
				diff := aliveTimeout - elapsed
				if diff.Seconds() > 5 {
					action.NextTimeout = diff
					log.Debug().Msgf("%s addr '%s' last seen before %s, calculate new timestamp", name, addr, elapsed)
					break
				}
			}
		}

		err = p.sendLinkAlive(addr, nil, &action, nil)
		if err == nil {
			automate.ClearErrorCounter(addr)
		}

	}

	if err != nil {

		if err == dispatcher.ErrPeerDisconnected || err == dispatcher.ErrDeviceDisconnected {
			// connection lost, handler shutdown
			automate.ChangeState(name, addr, shutdown)
			t.Tick()
			return
		}

		currentState := p.automate.GetState(addr)
		if currentState == shutdown {
			// prepare connection handler shutdown
			t.Tick()
			return
		}

	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handlePacketAcknowledge(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	automate.ChangeState(name, addr, idle)
	if pendingAction, exist := automate.GetPendingAction(addr); exist {
		automate.ClearPendingAction(addr)
		if pendingAction.NextState != initial {
			dispatcher.ChangeLinkState(addr, pendingAction.NextState)
		}
		action.NextTimeout = pendingAction.NextTimeout
	} else {
		action.NextTimeout = retryDelay
	}
	return nil
}

func (p *Plugin) handlePacketRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	automate.ChangeState(name, addr, idle)
	automate.ClearPendingAction(addr)
	dispatcher.ChangeLinkState(addr, linkDown)
	return nil
}

func (p *Plugin) handleLinkDescription(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	if record, found := p.driver.GetPayload(in); found {
		log.Info().Msgf("%s.%s received link description '%s'", name, addr, record[:])
		p.driver.ClearLinkDescription(addr)
		// vorbereitung
		// if family, exist := p.driver.GetField(in, "IF"); exist {

		// }
	}
	dispatcher.ChangeLinkState(addr, action.NextState)
	return nil
}

func (p *Plugin) handleLinkRecord(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	if record, found := p.driver.GetPayload(in); found {
		if len(record) > 7 && record[:2] == "RI" && record[5:7] == "FL" {
			log.Info().Msgf("%s.%s received link record '%s' fields '%s'", name, addr, record[2:4], record[7:])
			p.driver.SetLinkRecord(addr, record[2:4], record[7:])
		}
	}
	dispatcher.ChangeLinkState(addr, action.NextState)
	return nil
}

func (p *Plugin) sendLinkStart(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	packet := p.driver.ConstructPacket(addr, template.PacketLinkStart, "", "LS", nil)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendLinkAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	packet := p.driver.ConstructPacket(addr, template.PacketLinkAlive, "", "LA", nil)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

// func (p *Plugin) sendLinkEnd(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
// 	packet := p.driver.ConstructPacket(addr, template.PacketLinkEnd, "", "LE", nil)
// 	err := p.automate.SendPacket(addr, packet, action)
// 	return err
// }
