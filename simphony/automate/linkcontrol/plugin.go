package linkcontrol

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	lc "github.com/weareplanet/ifcv5-main/ifc/automate/linkcontrol"

	simphony "github.com/weareplanet/ifcv5-drivers/simphony/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/simphony/template"
)

// Plugin ...
type Plugin struct {
	*lc.Plugin
	driver *simphony.Dispatcher
}

const (
	aliveTimeout = simphony.AliveTimeout
	retryTimeout = simphony.RetryDelay
)

// New return a new plugin
func New(parent *simphony.Dispatcher) *Plugin {

	p := &Plugin{
		lc.New(),
		parent,
	}

	p.Setup(lc.Config{
		InitHandler: p.init,
	})

	return p
}

func (p *Plugin) init() {
	// state maschine rules

	p.RegisterRule(template.PacketPing, lc.LinkUp, p.handlePing, dispatcher.StateAction{NextState: lc.LinkUp, NextTimeout: aliveTimeout})
	p.RegisterRule(template.PacketPing, lc.LinkDown, p.handlePing, dispatcher.StateAction{NextState: lc.LinkUp, NextTimeout: aliveTimeout})
	p.RegisterRule(template.PacketHandshake, lc.LinkUp, p.handlePing, dispatcher.StateAction{})

	p.RegisterHandler(lc.LinkDown, p.handleLinkDown)
	p.RegisterHandler(lc.LinkUp, p.handleLinkUp)
}

func (p *Plugin) handleLinkDown(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = retryTimeout

	if !p.IsSerialDevice(addr) {
		p.ChangeLinkState(addr, lc.LinkUp)

	}

	return nil
}

func (p *Plugin) handleLinkUp(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = aliveTimeout

	return nil
}

func (p *Plugin) handlePing(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
	p.ChangeLinkState(addr, lc.LinkUp)

	err := p.SendPacket(addr, in, action)

	return err
}
