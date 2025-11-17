package centigram

import (
	"github.com/weareplanet/ifcv5-drivers/centigram/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

// ProcessIncomingPacket callback, ack/nak
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket) {

	if in == nil {
		return
	}

	// hier müssen wir etwas tricksen
	// MessageWaitingStatus kann ein spontan Paket für den Request Automaten sein,
	// oder wenn ein Database Swap aktiv ist, auch ein Bestätigungspaket für MessageLampOn sein.

	if in.Name == template.PacketMessageWaitingStatus && d.GetSwapState(in.Addr) {
		in.Name = template.PacketMessageWaitingStatusOnSwap
	}

	if d.IsAcknowledgement(in.Addr, in.Name) {

		switch in.Name {

		case template.PacketAck, template.PacketNak:
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
