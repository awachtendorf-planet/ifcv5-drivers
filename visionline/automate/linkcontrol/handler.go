package linkcontrol

import (
	// "fmt"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/visionline/template"
)

const (
	warnMissingHeartbeat    = 2 * time.Minute
	closeOnMissingHeartbeat = false
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

		if dispatcher.IsAcknowledgement(addr, "") {
			if packetTimeout >= retryDelay { // if protocol acknowledgement timeout > retryTimeout then retry immediately
				action.NextTimeout = nextActionDelay
			} else { // correct next timeout
				action.NextTimeout = (retryDelay - packetTimeout)
			}
		}

	case busy: // packet NAK, retry after retryDelay
		action.NextTimeout = aliveTimeout
		err = p.sendAlive(addr, &action)

	case linkDown:
		action.NextTimeout = aliveTimeout
		err = p.sendAlive(addr, &action)

	case linkUp:

		action.NextTimeout = aliveTimeout
		dispatcher := automate.Dispatcher()

		lastAlive := dispatcher.GetAlive(addr)
		if lastAlive > 0 {

			elapsed := time.Since(time.Unix(0, lastAlive))

			// missing heartbeat
			if !dispatcher.IsSerialDevice(addr) {
				if elapsed >= warnMissingHeartbeat {
					log.Warn().Msgf("%s addr '%s' last seen before %s, missing heartbeat", name, addr, elapsed)
					if closeOnMissingHeartbeat && !dispatcher.IsSerialDevice(addr) {
						p.closeConnection(addr)
						break
					}
				}
			}

			// re-calculate timestamp
			if elapsed < aliveTimeout {
				diff := aliveTimeout - elapsed
				if diff.Seconds() > 5 {
					action.NextTimeout = diff
					log.Debug().Msgf("%s addr '%s' last seen before %s, calculate new timestamp", name, addr, elapsed)
					break
				}
			}
		}

		err = p.sendAlive(addr, &action)
	}

	if err != nil {
		currentState := p.automate.GetState(addr)

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
	name := automate.Name()

	// use busy state, because we did not change to link-down, retry after retryDelay
	automate.ChangeState(name, addr, busy)
	automate.ClearPendingAction(addr)
	return nil
}

func (p *Plugin) handleAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	driver := p.driver

	if result, exist := driver.GetResultCode(in); !exist || result != 0 {
		return nil
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()

	dispatcher.SetAlive(addr)

	action.NextTimeout = aliveTimeout
	dispatcher.ChangeLinkState(addr, action.NextState)
	return nil
}

func (p *Plugin) sendAlive(addr string, action *dispatcher.StateAction) error {
	packet, _ := p.driver.ConstructPacket(addr, template.PacketAlive, "", nil)
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
