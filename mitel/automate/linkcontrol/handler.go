package linkcontrol

import (
	"time"

	syncstate "github.com/weareplanet/ifcv5-drivers/mitel/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/mitel/template"
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

		p.freeCommunication(addr)

		linkState := dispatcher.GetLinkState(addr)
		automate.ChangeState(name, addr, linkState)

		automate.ClearPendingAction(addr)

		if p.checkErrorCounter(addr) {
			break
		}

		// correct next timeout
		if dispatcher.IsAcknowledgement(addr, "") {
			if packetTimeout >= retryDelay {
				// if protocol acknowledgement timeout > retryTimeout then retry immediately
				action.NextTimeout = nextActionDelay
			} else {
				// correct next timeout
				action.NextTimeout = (retryDelay - packetTimeout)
			}
		}

	case busy, linkDown:

		// send initial ENQ
		action.NextTimeout = packetTimeout
		err = p.sendEnquiry(addr, &action, "", nil)

	case linkUp:

		driver := p.driver
		if !driver.SendAliveRecord(addr) {
			action.NextTimeout = 60 * time.Second
			break
		}

		dispatcher := automate.Dispatcher()
		lastAlive := dispatcher.GetAlive(addr)
		if lastAlive > 0 {
			elapsed := time.Since(time.Unix(0, lastAlive))
			if elapsed < aliveTimeout {
				diff := aliveTimeout - elapsed
				if diff.Seconds() > 1 {
					action.NextTimeout = diff
					log.Debug().Msgf("%s addr '%s' last seen before %s, calculate new timestamp", name, addr, elapsed)
					break
				}
			}
		}

		// synchronization between datachange automate and linkcontrol automate
		if !p.blockCommunication(addr) {
			action.NextTimeout = retryDelay
			log.Debug().Msgf("%s addr '%s' communication blocked, retry in %s", name, addr, action.NextTimeout)
			break
		}

		if p.checkErrorCounter(addr) {
			break
		}

		automate.AddErrorCounter(addr) // artificially set, ENQ -> ACK, ALIVE -> NAK handling

		action.NextTimeout = nextActionDelay
		action.NextState = nextAction // on success send AREYUTHERE
		err = p.sendEnquiry(addr, &action, "", nil)

	case nextAction:

		automate.AddErrorCounter(addr) // artificially set, ENQ -> ACK, ALIVE -> NAK handling

		action.NextTimeout = nextActionDelay
		action.NextState = nextRecord // on success reset error counter, ENQ -> ACK, ALIVE -> NAK handling
		action.CurrentState = linkUp  // on error fallback to sendEnquiry
		err = p.sendAlive(addr, &action, "", nil)

		p.freeCommunication(addr) // ok here, because ack slot is still registered for linkcontrol automate and datachange automate can not register a ack slot

	case nextRecord:

		action.NextTimeout = aliveTimeout
		automate.ClearErrorCounter(addr) // clear error, ENQ -> ACK, ALIVE -> NAK handling
		automate.ChangeState(name, addr, linkUp)

	case resyncStart:

		dispatcher := automate.Dispatcher()
		dispatcher.UnsetAcknowledgement(name, addr)

		driver := p.driver
		if !driver.SwapRequest(addr) {
			action.NextState = linkUp
			action.NextTimeout = aliveTimeout
			automate.ChangeState(name, addr, linkUp)
			dispatcher.ChangeLinkState(addr, linkUp)
			break
		}

		dispatcher.ChangeLinkState(addr, linkDown)

		// make sure we first t.Tick the automate
		defer func() {
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

func (p *Plugin) checkErrorCounter(addr string) bool {

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	if errorCounter := automate.GetErrorCounter(addr); errorCounter >= maxError {
		automate.ClearErrorCounter(addr)
		if dispatcher.IsSerialDevice(addr) {
			linkState := dispatcher.GetLinkState(addr)
			if linkState == linkUp {
				log.Warn().Msgf("%s addr '%s' reset connection, err=%s", name, addr, "mitel did not answer")
			}
			automate.ChangeState(name, addr, linkDown)
			dispatcher.ChangeLinkState(addr, linkDown)
		} else {
			p.closeConnection(addr)
			log.Warn().Msgf("%s addr '%s' close connection, err=%s", name, addr, "mitel did not answer")
			return true
		}
	}

	return false
}

func (p *Plugin) handlePacketAcknowledge(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	if pendingAction, exist := automate.GetPendingAction(addr); exist {

		automate.ClearPendingAction(addr)
		action.NextTimeout = pendingAction.NextTimeout

		if in.Name == template.PacketAck {

			linkState := dispatcher.GetLinkState(addr)
			if linkState == linkDown || linkState == busy {
				automate.ChangeState(name, addr, linkUp)
				dispatcher.ChangeLinkState(addr, linkUp)
				return nil
			}

		}

		automate.ChangeState(name, addr, pendingAction.NextState)
	} else {
		automate.ChangeState(name, addr, idle)
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
	} else {
		automate.ChangeState(name, addr, busy)
	}

	errorCounter := automate.GetErrorCounter(addr)
	if errorCounter >= maxError {
		action.NextTimeout = nextActionDelay
	} else {
		action.NextTimeout = retryDelay
	}
	return nil
}

func (p *Plugin) handlePacketEnquire(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	owner := dispatcher.UnsetAcknowledgement(name, addr)

	if owner {
		linkState := dispatcher.GetLinkState(addr)
		if linkState == linkDown {
			p.handlePacketAcknowledge(addr, in, action, job)
		} else {
			p.handlePacketRefusion(addr, in, action, job)
		}
	}

	return nil
}

func (p *Plugin) handleSwapEnquiry(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	name := automate.Name()
	automate.ChangeState(name, addr, action.NextState)
	return nil
}

func (p *Plugin) setAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	dispatcher.SetAlive(addr)
	return nil
}

func (p *Plugin) sendEnquiry(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketEnq, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendAlive(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketAlive, tracking, context)
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
