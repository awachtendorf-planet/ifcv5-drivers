package guestlink

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"

	"github.com/spf13/cast"
)

// ProcessIncomingPacket callback, ack/nak
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket) {

	if in == nil {
		return
	}

	if d.IsAcknowledgement(in.Addr, in.Name) {

		switch in.Name {

		case template.PacketAck, template.PacketNak, template.PacketEnq:
			break

		case template.PacketGarbage:
			d.sendRefusion(in.Addr)

		default:
			d.sendAcknowledge(in.Addr)

		}
	}

}

func (d *Dispatcher) sendAcknowledge(addr string) {
	driver := d.Network
	data := ifc.NewLogicalPacket(template.PacketAck, addr, "")
	driver.Send(data, 0)
}

func (d *Dispatcher) sendRefusion(addr string) {
	driver := d.Network
	data := ifc.NewLogicalPacket(template.PacketNak, addr, "")
	driver.Send(data, 0)
}

func (d *Dispatcher) GetErrorCode(in *ifc.LogicalPacket) int {

	errorCode := int(0)

	if in == nil {
		return errorCode
	}

	data := in.Data()

	if value, exist := data[guestlink_error]; exist {
		errorCode = cast.ToInt(value)
	}

	return errorCode
}
