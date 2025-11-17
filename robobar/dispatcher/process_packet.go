package robobar

import (
	"github.com/weareplanet/ifcv5-drivers/robobar/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

// ProcessIncomingPacket callback, ack/nak, encoding decode
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket) {
	if in == nil {
		return
	}

	if d.IsAcknowledgement(in.Addr, in.Name) {
		switch in.Name {

		case template.PacketAck, template.PacketNak:
			break

		case template.PacketEnq:
			d.sendAcknowledge(in.Addr)

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
