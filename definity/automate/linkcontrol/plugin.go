package linkcontrol

import (
	"time"

	definity "github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	syncstate "github.com/weareplanet/ifcv5-drivers/definity/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	lc "github.com/weareplanet/ifcv5-main/ifc/automate/linkcontrol"

	"github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/definity/template"
)

// Plugin ...
type Plugin struct {
	*lc.Plugin
	driver *definity.Dispatcher
}

const (
	aliveTimeout = 15 * time.Second // 5 - 20 seconds
	retryTimeout = 5 * time.Second
)

// New return a new plugin
func New(parent *definity.Dispatcher) *Plugin {

	p := &Plugin{
		lc.New(),
		parent,
	}

	p.Setup(lc.Config{
		InitHandler:     p.init,
		TemplateAck:     template.PacketAck,
		TemplateNak:     template.PacketNak,
		HandleEvent:     p.handleEvent,
		DisableSetAlive: true, // set explicit, disable automatic
	})

	return p
}

func (p *Plugin) init() {

	// state maschine rules
	p.RegisterRule(template.PacketHeartbeatReply, lc.LinkDown, p.handleHeartbeat, dispatcher.StateAction{NextState: lc.LinkUp})
	p.RegisterRule(template.PacketHeartbeatReply, lc.LinkUp, p.handleHeartbeat, dispatcher.StateAction{NextState: lc.LinkUp})

	p.RegisterRule(template.PacketSyncRequest, lc.LinkDown, p.handleSyncRequest, dispatcher.StateAction{NextState: lc.DBSync})
	p.RegisterRule(template.PacketSyncRequest, lc.LinkUp, p.handleSyncRequest, dispatcher.StateAction{NextState: lc.DBSync})

	p.RegisterRule(template.PacketLinkEndRequest, lc.LinkUp, p.handleLinkEnd, dispatcher.StateAction{NextState: lc.LinkDown})
	p.RegisterRule(template.PacketLinkEndRequest, lc.LinkDown, p.handleLinkEnd, dispatcher.StateAction{NextState: lc.LinkDown})

	// state maschine callbacks
	p.RegisterHandler(lc.LinkDown, p.handleLinkDown)
	p.RegisterHandler(lc.LinkUp, p.handleLinkUp)
	p.RegisterHandler(lc.Busy, p.handleLinkDown)
	p.RegisterHandler(lc.DBSync, p.handleSyncState)

	// additional events
	p.RegisterEvents(syncstate.End.String(), syncstate.Error.String())

}

func (p *Plugin) handleHeartbeat(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	dispatcher := p.GetDispatcher()
	dispatcher.ChangeLinkState(addr, action.NextState)

	// set explicit, to detect missing heartbeat reply
	dispatcher.SetAlive(addr)

	action.NextTimeout = aliveTimeout

	return nil
}

func (p *Plugin) handleSyncRequest(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	dispatcher := p.GetDispatcher()
	dispatcher.ChangeLinkState(addr, action.NextState)

	action.NextTimeout = definity.PmsSyncTimeout

	// trigger dbsync automate
	broker := dispatcher.Broker()
	if broker != nil {
		broker.Broadcast(syncstate.NewEvent(addr, syncstate.Start, nil), syncstate.Start.String())
	}

	return nil
}

func (p *Plugin) handleLinkEnd(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	dispatcher := p.GetDispatcher()
	dispatcher.ChangeLinkState(addr, action.NextState)

	action.NextTimeout = retryTimeout

	err := p.send(addr, action, template.PacketLinkEndConfirmed)

	return err
}

func (p *Plugin) handleLinkDown(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = retryTimeout

	err := p.send(addr, action, template.PacketHeartbeat)

	return err
}

func (p *Plugin) handleLinkUp(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = retryTimeout

	err := p.send(addr, action, template.PacketHeartbeat)

	if lastAlive := p.GetAlive(addr); lastAlive > 0 {

		elapsed := time.Since(time.Unix(0, lastAlive))

		if elapsed > aliveTimeout+(retryTimeout*2) { // missing heartbeat reply, change to linkdown
			dispatcher := p.GetDispatcher()
			dispatcher.ChangeLinkState(addr, lc.LinkDown)
		}
	}

	return err
}

func (p *Plugin) handleEvent(addr string, event interface{}, action *dispatcher.StateAction) bool {

	switch event.(type) {

	case *syncstate.Event:

		event := event.(*syncstate.Event)

		if addr == event.Addr {

			switch event.State {

			case syncstate.End: // sync success

				p.ChangeLinkState(addr, lc.LinkDown)

				action.NextState = lc.LinkDown
				action.NextTimeout = 1 * time.Second

				return true

			case syncstate.Error: // sync failed

				p.ChangeLinkState(addr, lc.LinkDown)

				action.NextState = lc.LinkDown
				action.NextTimeout = 5 * time.Second

				return true

			}

		}

	}

	return false
}

func (p *Plugin) handleSyncState(addr string, action *dispatcher.StateAction) error {

	name := p.GetName()
	log.Debug().Msgf("%s addr '%s' pending sync records, waiting for dbsync automate", name, addr)

	action.NextTimeout = definity.PmsSyncTimeout

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, "", nil)

	if err != nil {

		name := p.GetName()
		log.Error().Msgf("%s addr '%s' construct packet '%s' failed, err=%s", name, addr, packetName, err)

		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
