package request

import (
	"fmt"
	"strings"

	micros "github.com/weareplanet/ifcv5-drivers/micros/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/micros/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/micros/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *micros.Dispatcher
}

// New return a new plugin
func New(parent *micros.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,

		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(template.PacketInquiryRequest,
		template.PacketOutletChargeRequest,
		template.PacketPing,
	)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketPing:

		err := p.SendPacket(addr, packet, action)
		return err

	case template.PacketOutletChargeRequest:

		inquiry := p.getField(packet, "AccountID", true)
		selectionNumber := p.getField(packet, "SelectionNumber", true)

		if len(inquiry) != 0 && selectionNumber == "0" {
			return p.handlePostingInquiry(addr, packet, action)
		}

		return p.handlePostingCharge(addr, packet, action)

	case template.PacketInquiryRequest:

		return p.handlePostingInquiry(addr, packet, action)

	default:
		return errors.Errorf("no handler defined to process packet '%s'", packet.Name)

	}
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
