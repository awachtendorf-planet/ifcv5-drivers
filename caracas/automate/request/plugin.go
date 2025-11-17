package request

import (
	"fmt"
	"strings"

	caracas "github.com/weareplanet/ifcv5-drivers/caracas/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	//"github.com/weareplanet/ifcv5-main/ifc/record"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/caracas/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/caracas/template"

	// "github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *caracas.Dispatcher
}

// New return a new plugin
func New(parent *caracas.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		// register Ack/Nak packets if needed
		// TemplateAck: template.PacketAck,
		// TemplateNak: template.PacketNak,

		// register init handler if needed
		// InitHandler: p.init,

		// register pre/post handler
		// PreHandler:    p.preHandler,
		// PostHandler:   p.postHandler,

		// register incoming packet handler
		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(
		template.PacketCallCharge,
		template.PacketMinibarCharge,
		template.PacketMinibarChargeAdvanced,
		template.PacketMinibarChargeComplete,
		template.PacketVoicemail,
		template.PacketVoicemailAdvanced,
		template.PacketRoomstatus,
		template.PacketRoomstatusAdvanced,
		template.PacketWakeupOrder,
		template.PacketWakeupResult,
		template.PacketDBSyncRequest,
	)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketCallCharge:
		return p.handleCallPacket(addr, packet, action)

	case template.PacketMinibarCharge, template.PacketMinibarChargeAdvanced, template.PacketMinibarChargeComplete:
		return p.handleMinibarCharge(addr, packet, action)

	case template.PacketRoomstatus, template.PacketRoomstatusAdvanced:
		return p.handleRoomStatus(addr, packet)

	case template.PacketVoicemail, template.PacketVoicemailAdvanced:
		return p.handleVoicemail(addr, packet)

	case template.PacketWakeupOrderIncoming:
		return p.handleWakeupOrder(addr, packet)

	case template.PacketWakeupResult:
		return p.handleWakeupResult(addr, packet)

	case template.PacketDBSyncRequest:
		return p.handleDBSyncReq(addr, packet)

	default:
		return fmt.Errorf("no handler defined to process packet '%s'", packet.Name)

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
