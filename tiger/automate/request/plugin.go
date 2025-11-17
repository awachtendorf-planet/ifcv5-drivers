package request

import (
	"fmt"
	"strings"

	tiger "github.com/weareplanet/ifcv5-drivers/tiger/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	//"github.com/weareplanet/ifcv5-main/ifc/record"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/tiger/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/tiger/template"
	"github.com/weareplanet/ifcv5-main/log"

	// "github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *tiger.Dispatcher
}

// New return a new plugin
func New(parent *tiger.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(template.PacketRoomBasedCall,
		template.PacketMiniBarBilling,
		template.PacketOtherChargePosting,
		template.PacketRoomstatus,
		// template.PacketMessageWaitingGuest,
		// template.PacketMessageWaitingReservation,
	)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketRoomstatus:
		p.handleRoomStatus(addr, packet)

	case template.PacketMiniBarBilling, template.PacketOtherChargePosting:
		p.handleCharges(addr, packet)

	case template.PacketRoomBasedCall:
		p.handleCallPacket(addr, packet)

	default:
		name := p.GetName()
		log.Warn().Msgf("%s addr '%s' unknown packet '%s'", name, addr, packet.Name)
	}

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
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
