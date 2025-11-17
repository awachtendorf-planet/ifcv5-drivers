package linkcontrol

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"

	"github.com/pkg/errors"
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

		p.clearTransaction(addr)

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
			dispatcher.ChangeLinkState(addr, linkDown)
		}

		if dispatcher.IsAcknowledgement(addr, "") {
			if packetTimeout >= retryDelay { // if protocol acknowledgement timeout > retryTimeout then retry immediately
				action.NextTimeout = nextActionDelay
			} else { // correct next timeout
				action.NextTimeout = (retryDelay - packetTimeout)
			}
		}

	case linkDown:

		p.clearTransaction(addr)

		action.NextTimeout = answerTimeout
		err = p.sendPacket(addr, &action, template.PacketStart)

	case hello:

		p.clearTransaction(addr)
		dispatcher := automate.Dispatcher()

		// wir betrachten das HELO Paket als optional
		// wenn das unbeantwortet oder abgelehnt wird, gehen wir trotzdem in LinpUp
		// da das initiale STRT Paket erfolgreich versendet wurde

		if !p.isAnswerExpected(addr) {
			station, _ := dispatcher.GetStationAddr(addr)
			if p.driver.SendHeloPacket(station) {
				action.NextTimeout = answerTimeout
				err = p.sendPacket(addr, &action, template.PacketHelo)
				break
			}
		}

		p.receivedAnswer(addr)
		dispatcher.ChangeLinkState(addr, linkUp)
		t.Tick()

		return

	case linkUp:

		if p.isAnswerExpected(addr) {
			p.clearTransaction(addr)
			dispatcher := automate.Dispatcher()
			dispatcher.ChangeLinkState(addr, linkDown)
			t.Tick()
			return
		}

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

		// reset transaction on error
		p.driver.FreeTransaction(p, addr)

		log.Error().Msgf("%s %s", name, err)

	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) handlePacketAcknowledge(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {
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

func (p *Plugin) handlePacketRefusion(addr string, _ *ifc.LogicalPacket, _ *dispatcher.StateAction, _ *order.Job) error {
	p.clearTransaction(addr)

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	automate.ChangeState(name, addr, idle)
	automate.ClearPendingAction(addr)
	dispatcher.ChangeLinkState(addr, linkDown)

	return nil
}

func (p *Plugin) handleCommandAcknowledge(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {
	p.clearTransaction(addr)

	automate := p.automate
	dispatcher := automate.Dispatcher()

	if !p.isAnswerExpected(addr) {
		return nil
	}
	p.receivedAnswer(addr)

	dispatcher.ChangeLinkState(addr, action.NextState)
	return nil
}

func (p *Plugin) handleCommandRefusion(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {
	p.clearTransaction(addr)

	automate := p.automate
	dispatcher := automate.Dispatcher()

	if !p.isAnswerExpected(addr) {
		return nil
	}
	p.receivedAnswer(addr)

	dispatcher.ChangeLinkState(addr, action.NextState)
	return nil
}

func (p *Plugin) sendCommandAcknowledge(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	// wenn wir ein STRT oder TEST Paket erhalten, dann gehen wir ggf. gleich in LinkUp über ohne auf das VER/ERR Packet zu warten
	p.receivedAnswer(addr)

	driver := p.driver
	transaction := driver.GetTransaction(in)
	if transaction.IsLastSequence() {
		if packet, err := driver.ConstructPacket(addr, template.PacketVerify, "", nil, transaction); err == nil {
			err = p.automate.SendPacket(addr, packet, action)
		} else {
			return err
		}
	}
	return nil
}

func (p *Plugin) sendPacket(addr string, action *dispatcher.StateAction, packetName string) error {

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	transaction, err := driver.NewTransaction(addr)
	if err != nil {
		log.Error().Msgf("%s new transaction failed, err=", name, err)
		return errors.New("transaction handling failed")
	}

	if err = driver.RegisterTransaction(p, addr, transaction); err != nil {
		log.Error().Msgf("%s register transaction '%s' failed, err=%s", name, transaction.Identifier, err)
		return errors.New("transaction handling failed")
	}

	packet, err := p.driver.ConstructPacket(addr, packetName, "", nil, transaction)
	if err == nil {
		if err = p.automate.SendPacket(addr, packet, action); err == nil {
			p.expectAnswer(addr)
			return nil
		}
	}

	driver.UnregisterTransaction(p, addr, transaction)
	return err
}

func (p *Plugin) clearTransaction(addr string) {

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	if transaction, exist := driver.LastTransaction(p, addr); exist {
		if err := driver.UnregisterTransaction(p, addr, transaction); err != nil {
			log.Error().Msgf("%s unregister transaction '%s' failed, err=", name, transaction.Identifier, err)
		}
	}
}
