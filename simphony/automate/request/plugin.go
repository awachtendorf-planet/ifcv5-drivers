package request

import (
	"fmt"
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	simphony "github.com/weareplanet/ifcv5-drivers/simphony/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/simphony/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *simphony.Dispatcher
}

// New return a new plugin
func New(parent *simphony.Dispatcher) *Plugin {

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

	p.RegisterPacket(template.PacketPostingInquiry,
		template.PacketChargePostingSelectedGuest,
		template.PacketNonFrontOfficeCharge,
		template.PacketCheckFacsimileRequest,
	)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketCheckFacsimileRequest:

		err := p.handleFascimile(addr, packet, action)
		return err

	case template.PacketPing:

		err := p.SendPacket(addr, packet, action)
		return err

	case template.PacketNonFrontOfficeCharge:

		return p.handleNonFrontOfficePostingCharge(addr, packet, action)

	case template.PacketChargePostingSelectedGuest:

		inquiry := p.getField(packet, "GuestID", true)
		selectionNumber := p.getField(packet, "SelectionNumber", true)

		if len(inquiry) != 0 && selectionNumber == "0" {
			return p.handlePostingInquiry(addr, packet, action)
		}

		return p.handlePostingCharge(addr, packet, action)

	case template.PacketPostingInquiry:

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
