package linkcontrol

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/vc3000/record"
	"github.com/weareplanet/ifcv5-drivers/vc3000/template"
)

const (
	maxIdleTime     = 15 * time.Minute
	minSlotsBacklog = 3
)

func (p *Plugin) handleTimeout(addr string, state automatestate.State, t *ticker.ResetTicker) {

	automate := p.automate
	name := automate.Name()

	var err error
	action := dispatcher.StateAction{NextTimeout: retryDelay}

	switch state {

	case packetSent, commandSent:
		automate.ChangeState(name, addr, idle)

	case linkDown:
		licenseCode, ok := automate.GetSetting(addr, "LicenseCode", "").(string)
		if !ok || len(licenseCode) == 0 {
			log.Warn().Msgf("%s addr '%s' no license code configured", name, addr)
		}

		record := record.Register{
			LicenseCode: licenseCode,
		}
		err = p.sendRegister(addr, nil, &action, record)

	case linkUp:
		// cleanup unused slot
		p.cleanupSlots(addr)

		action = dispatcher.StateAction{NextTimeout: aliveTimeout}
	}

	if err != nil {

		if err == dispatcher.ErrPeerDisconnected || err == dispatcher.ErrDeviceDisconnected { // IFCDEV-65 unnötige Erweiterung
			// connection lost, handler shutdown
			automate.ChangeState(name, addr, shutdown)
			t.Tick()
			return
		}

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

func (p *Plugin) cleanupSlots(addr string) {
	automate := p.automate
	driver := p.driver
	name := automate.Name()

	slot, exist := driver.GetSlotInfo(addr)
	if !exist {
		return
	}

	timestamp := time.Unix(0, slot.Timestamp)
	lastActive := time.Since(timestamp)
	log.Debug().Msgf("%s addr '%s' slot currently in use: %t, last assignment %s before", name, addr, slot.Used, lastActive)

	if !slot.Used && lastActive >= maxIdleTime {
		slotsCount := driver.GetSlotsCount(addr)
		if slotsCount > minSlotsBacklog && driver.RemoveSlotIfUnsed(addr) {
			log.Info().Msgf("%s addr '%s' close slot, because of longer idle time, slots backlog: %d", name, addr, slotsCount-1)
			p.closeConnection(addr)
		}
	}

}

func (p *Plugin) handleRegister(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	success := false

	if data, ok := p.driver.GetPayload(in); ok && len(data) == 44 {
		ret := p.driver.GetUint(data[40:44])
		success = (ret == 0)
	} else {
		log.Warn().Msgf("%s addr '%s' packet '%s' expected payload length: 44, got: %d", name, addr, in.Name, len(data))
	}

	if success {
		action.NextTimeout = aliveTimeout
		dispatcher.ChangeLinkState(addr, action.NextState)
		p.addSlot(addr)
		return nil
	}

	action.NextTimeout = pendingDelay
	automate.ChangeState(name, addr, idle)
	return nil
}

func (p *Plugin) handleRegisterAck(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	action.NextTimeout = aliveTimeout
	dispatcher.ChangeLinkState(addr, action.NextState)
	p.addSlot(addr)
	return nil
}

func (p *Plugin) handleRegisterNak(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	name := automate.Name()
	action.NextTimeout = pendingDelay
	automate.ChangeState(name, addr, idle)
	return nil
}

func (p *Plugin) handleUnregister(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	p.removeSlot(addr)
	p.closeConnection(addr)
	dispatcher.ChangeLinkState(addr, action.NextState)
	return nil
}

func (p *Plugin) sendRegister(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, context interface{}) error {
	packet, _ := p.driver.ConstructPacket(addr, template.PacketRegister, false, "", context)
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

func (p *Plugin) addSlot(addr string) {
	driver := p.driver
	if driver != nil {
		driver.AddSlot(addr)
	}
}

func (p *Plugin) removeSlot(addr string) {
	driver := p.driver
	if driver != nil {
		driver.RemoveSlot(addr)
	}
}
