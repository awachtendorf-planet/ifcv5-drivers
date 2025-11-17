package messerschmitt

import (
	"github.com/weareplanet/ifcv5-drivers/messerschmitt/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

// ProcessIncomingPacket callback, ack/nak
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket) {

	if in == nil {
		return
	}

	if d.IsAcknowledgement(in.Addr, in.Name) {

		switch in.Name {

		case template.PacketAck, template.PacketNak, template.PacketEOT, template.PacketGarbage:

		case template.PacketEnq:
			d.sendAcknowledge0(in.Addr)

		case template.PacketUnknown, template.PacketWACK, template.PacketTTD:
			d.sendRefusion(in.Addr)

		default:
			d.sendAcknowledge1(in.Addr)
		}
	}

}

func (d *Dispatcher) sendAcknowledge0(addr string) {
	driver := d.Network
	data := ifc.NewLogicalPacket(template.PacketAck0, addr, "")
	driver.Send(data, 0)
}

func (d *Dispatcher) sendAcknowledge1(addr string) {
	driver := d.Network
	data := ifc.NewLogicalPacket(template.PacketAck1, addr, "")
	driver.Send(data, 0)
}

func (d *Dispatcher) sendRefusion(addr string) {
	driver := d.Network
	data := ifc.NewLogicalPacket(template.PacketNak, addr, "")
	driver.Send(data, 0)
}
