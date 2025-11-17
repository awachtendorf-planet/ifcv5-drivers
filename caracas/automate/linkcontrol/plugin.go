package linkcontrol

import (
	"fmt"
	"time"

	caracas "github.com/weareplanet/ifcv5-drivers/caracas/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	lc "github.com/weareplanet/ifcv5-main/ifc/automate/linkcontrol"

	"github.com/weareplanet/ifcv5-drivers/caracas/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/caracas/template"
)

// Plugin ...
type Plugin struct {
	*lc.Plugin
	driver *caracas.Dispatcher
}

const (
	aliveTimeout = 30 * time.Second
)

// New return a new plugin
func New(parent *caracas.Dispatcher) *Plugin {

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
	p.RegisterRule(template.PacketEnq, lc.LinkDown, p.handleEnq, dispatcher.StateAction{})
	p.RegisterRule(template.PacketSyn, lc.LinkDown, p.handleEnq, dispatcher.StateAction{})
	p.RegisterRule(template.PacketBind, lc.LinkDown, p.handleBind, dispatcher.StateAction{})

	p.RegisterRule(template.PacketEnq, lc.LinkUp, p.handleAlive, dispatcher.StateAction{})
	p.RegisterRule(template.PacketSyn, lc.LinkUp, p.handleAlive, dispatcher.StateAction{})

	// p.RegisterRule(template.PacketAnswerTG00, lc.LinkUp, p.handleLinkAlive, dispatcher.StateAction{})

	// state maschine callbacks
	p.RegisterHandler(lc.LinkDown, p.handleLinkDown)
	p.RegisterHandler(lc.LinkUp, p.handleLinkUp)
	p.RegisterHandler(lc.Busy, p.handleLinkDown) // called after packet refusion without linkstate changed. do whatever is necessary.

}

func (p *Plugin) handleEnq(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	action.NextTimeout = 1 * time.Second

	err := p.send(addr, action, template.PacketEot)
	if err == nil && p.GetState(addr) == lc.LinkDown {

		station, err := p.driver.GetStationAddr(addr)
		if err != nil {
			return err
		}
		fmt.Println(p.driver.GetBindMode(station))
		if !p.driver.GetBindMode(station) {
			p.ChangeLinkState(addr, lc.LinkUp)
		}
	}
	return err
}

func (p *Plugin) handleAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	action.NextTimeout = 1 * time.Second

	err := p.send(addr, action, template.PacketEot)
	return err
}

func (p *Plugin) handleBind(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	action.NextTimeout = 1 * time.Second

	err := p.send(addr, action, template.PacketBind)
	if err == nil {
		p.ChangeLinkState(addr, lc.LinkUp)
	}
	return err
}

func (p *Plugin) handleLinkDown(addr string, action *dispatcher.StateAction) error {

	action.NextTimeout = 15 * time.Second

	if !p.IsSerialDevice(addr) {
		p.ChangeLinkState(addr, lc.LinkUp)

	}

	return nil
}

func (p *Plugin) handleLinkUp(addr string, action *dispatcher.StateAction) error {

	station, err := p.driver.GetStationAddr(addr)
	if err != nil {
		return err
	}

	sendTimeMS := p.driver.GetSendTimeMS(station)

	if sendTimeMS > 0 {
		action.NextTimeout = time.Duration(sendTimeMS) * time.Millisecond
		err = p.send(addr, action, template.PacketDatetime)

	} else {
		action.NextTimeout = aliveTimeout
	}

	return err
}

// func (p *Plugin) handleLinkAlive(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

// 	p.SetAlive(addr)

// 	p.send(addr, action, template.PacketTG00)
// 	return nil
// }

// func (p *Plugin) handleLoginAnswer(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

// 	p.ChangeLinkState(addr, lc.LinkUp)

// 	if p.IsSerialDevice(addr) {
// 		p.send(addr, action, template.PacketTG00)
// 	}

// 	return nil
// }

func (p *Plugin) send(addr string, action *dispatcher.StateAction, name string) error {

	packet, _ := p.driver.ConstructPacket(addr, name, "", order.Internal, nil)
	err := p.SendPacket(addr, packet, action)

	return err
}
