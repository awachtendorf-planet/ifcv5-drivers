package linkcontrol

import (
	"time"

	detewe "github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	lc "github.com/weareplanet/ifcv5-main/ifc/automate/linkcontrol"

	"github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/detewe/template"
)

// Plugin ...
type Plugin struct {
	*lc.Plugin
	driver *detewe.Dispatcher
}

const (
	aliveTimeout = 30 * time.Second
)

// New return a new plugin
func New(parent *detewe.Dispatcher) *Plugin {

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
	p.RegisterRule(template.PacketLoginAnswer, lc.LinkDown, p.handleLoginAnswer, dispatcher.StateAction{NextState: lc.LinkUp})
	p.RegisterRule(template.PacketAnswerTG00, lc.LinkUp, p.handleLinkAlive, dispatcher.StateAction{})

	// state maschine callbacks
	p.RegisterHandler(lc.LinkDown, p.handleLinkDown)
	p.RegisterHandler(lc.LinkUp, p.handleLinkUp)
	p.RegisterHandler(lc.Busy, p.handleLinkDown) // called after packet refusion without linkstate changed. do whatever is necessary.

}

func (p *Plugin) handleLinkDown(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = aliveTimeout

	if p.IsSerialDevice(addr) {
		p.ChangeLinkState(addr, lc.LinkUp)

	} else {
		return p.send(addr, action, template.PacketTCPLogin)
	}

	return nil
}

func (p *Plugin) handleLinkUp(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = aliveTimeout

	return nil
}

func (p *Plugin) handleLinkAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	p.SetAlive(addr)

	p.send(addr, action, template.PacketTG00)
	return nil
}

func (p *Plugin) handleLoginAnswer(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	p.ChangeLinkState(addr, lc.LinkUp)

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, name string) error {

	packet, _ := p.driver.ConstructPacket(addr, name, "", order.Internal, nil)
	err := p.SendPacket(addr, packet, action)

	return err
}
