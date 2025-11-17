package request

import (
	"fmt"
	"strings"

	informatel "github.com/weareplanet/ifcv5-drivers/informatel/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"

	//"github.com/weareplanet/ifcv5-main/ifc/record"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/informatel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/informatel/template"

	// "github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *informatel.Dispatcher
}

// New return a new plugin
func New(parent *informatel.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,

		// register init handler if needed
		// InitHandler: p.init,

		// register pre/post handler
		// PreHandler:    p.preHandler,
		// PostHandler:   p.postHandler,

		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(template.PacketChargeRecord,
		template.PacketMinibarCharge,
		template.PacketRoomStatus,
		template.PacketVoicemail,
		template.PacketWakeupIncomingCall,
		template.PacketDBSwap,
	)

	return p
}

// func (p *Plugin) init() {
// 	// register some state handler if needed
// 	// p.RegisterHandler(rq.NextAction, p.handleNextAction)
// 	// p.RegisterHandler(rq.NextRecord, p.handleNextRecord)
// }

// func (p *Plugin) preHandler(addr string) {
// 	// called before the state maschine start
// }

// func (p *Plugin) postHandler(addr string) {
// 	// called after the state maschine finished
// }

// func (p *Plugin) handleNextAction(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {
// 	action.NextState = rq.NextRecord
// 	return nil
// }

// func (p *Plugin) handleNextRecord(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {
// 	action.NextState = rq.Success
// 	return nil
// }

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketChargeRecord:
		return p.handleCallPacket(addr, packet, action)

	case template.PacketMinibarCharge:
		return p.handleMinibarCharge(addr, packet, action)

	case template.PacketRoomStatus:
		return p.handleRoomStatus(addr, packet)

	case template.PacketVoicemail:
		return p.handleVoicemail(addr, packet)

	case template.PacketWakeupIncomingCall:
		return p.handleWakeupEvent(addr, packet)

	case template.PacketDBSwap:
		return p.handleDBSyncReq(addr, packet)

	default:
		return fmt.Errorf("no handler defined to process packet '%s'", packet.Name)

	}

	return nil
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
