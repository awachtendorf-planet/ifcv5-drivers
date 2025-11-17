package request

import (
	"fmt"
	"strings"

	detewe "github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	//"github.com/weareplanet/ifcv5-main/ifc/record"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/detewe/template"
	"github.com/weareplanet/ifcv5-main/log"

	// "github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *detewe.Dispatcher
}

// New return a new plugin
func New(parent *detewe.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(template.PacketTG10,
		template.PacketTG20,
		template.PacketTG40,
		template.PacketTG72,
	)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketTG20:
		p.handleRoomStatus(addr, packet)

	case template.PacketTG72:
		p.handleWakeupEvent(addr, packet)

	case template.PacketTG40:
		p.handleDataTransfer(addr, packet)

	case template.PacketTG10:
		p.handleCallPacket(addr, packet)

	default:
		name := p.GetName()
		log.Warn().Msgf("%s addr '%s' unknown packet '%s'", name, addr, packet.Name)
	}

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, job.Action, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}

func (p *Plugin) getField(packet *ifc.LogicalPacket, field string, trim bool) string {

	if packet == nil {
		return ""
	}

	data := packet.Data()
	value := cast.ToString(data[field])
	if trim {
		value = strings.Trim(value, " ")
	}
	return value
}
