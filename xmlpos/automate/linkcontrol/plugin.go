package linkcontrol

import (
	"time"

	xmlpos "github.com/weareplanet/ifcv5-drivers/xmlpos/dispatcher"
	vendor "github.com/weareplanet/ifcv5-drivers/xmlpos/record"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	lc "github.com/weareplanet/ifcv5-main/ifc/automate/linkcontrol"

	"github.com/weareplanet/ifcv5-drivers/xmlpos/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/xmlpos/template"
)

// Plugin ...
type Plugin struct {
	*lc.Plugin
	driver *xmlpos.Dispatcher
}

const (
	aliveTimeout = xmlpos.AliveTimeout
)

// New return a new plugin
func New(parent *xmlpos.Dispatcher) *Plugin {

	p := &Plugin{
		lc.New(),
		parent,
	}

	p.Setup(lc.Config{
		InitHandler:         p.init,
		GetInitialLinkState: func(addr string) automatestate.State { return automatestate.LinkUp },
		TemplateAck:         template.PacketAck,
		TemplateNak:         template.PacketNak,
	})

	return p
}

func (p *Plugin) init() {

	// state maschine rules
	p.RegisterRule(template.PacketLinkDescription, lc.LinkUp, p.sendLinkAlive, dispatcher.StateAction{})
	p.RegisterRule(template.PacketLinkStart, lc.LinkUp, p.sendLinkAlive, dispatcher.StateAction{})
	p.RegisterRule(template.PacketLinkAlive, lc.LinkUp, p.handleLinkAlive, dispatcher.StateAction{})

	// state maschine callbacks
	p.RegisterHandler(lc.LinkUp, p.handleLinkUp)

}

func (p *Plugin) handleLinkUp(addr string, action *dispatcher.StateAction) error {

	if lastAlive := p.GetAlive(addr); lastAlive > 0 {

		elapsed := time.Since(time.Unix(0, lastAlive))

		if elapsed < aliveTimeout {

			if elapsed.Seconds() < 1 {
				action.NextTimeout = aliveTimeout
				return nil
			}

			diff := aliveTimeout - elapsed
			if diff.Seconds() > 5 {

				action.NextTimeout = diff
				name := p.GetName()
				log.Debug().Msgf("%s addr '%s' last seen before %s, calculate new timestamp", name, addr, elapsed)
				return nil

			}

		}

	}

	action.NextTimeout = aliveTimeout

	err := p.send(addr, action, &vendor.LinkAlive{
		Date: time.Now(),
		Time: time.Now(),
	})

	return err

}

func (p *Plugin) handleLinkAlive(addr string, _ *ifc.LogicalPacket, _ *dispatcher.StateAction, _ *order.Job) error {

	p.SetAlive(addr)

	return nil
}

func (p *Plugin) sendLinkAlive(addr string, _ *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	err := p.send(addr, action, &vendor.LinkAlive{
		Date: time.Now(),
		Time: time.Now(),
	})

	if err == nil {
		p.SetAlive(addr)
	}

	return err
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, context interface{}) error {

	packet, _ := p.driver.ConstructPacket(addr, template.PacketFramed, "", context)
	err := p.SendPacket(addr, packet, action)

	return err
}
