package xmlpos

import (
	"github.com/weareplanet/ifcv5-drivers/xmlpos/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"
)

// ProcessIncomingPacket callback, ack/nak, encoding decode
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket) {

	if in == nil {
		return
	}

	switch in.Name {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return

	case template.PacketGarbage:
		return

	case template.PacketFramed:

		err := d.unmarshal(in)
		if err != nil {
			log.Warn().Msgf("unmarshal incoming packet '%s' from addr '%s' failed, err:=%s", in.Name, in.Addr, err)
		}

		if d.IsAcknowledgement(in.Addr, in.Name) {
			if err == nil {
				d.sendAcknowledge(in.Addr)
			} else {
				break
			}
		}

		return

	}

	if d.IsAcknowledgement(in.Addr, in.Name) {
		d.sendRefusion(in.Addr)
	}

}

// ProcessOutgoingPacket callback, encoding encode
func (d *Dispatcher) ProcessOutgoingPacket(out *ifc.LogicalPacket) {

	if out == nil {
		return
	}

	if err := d.marshal(out); err != nil {
		log.Warn().Msgf("marshal outgoing packet '%s' from addr '%s' failed, err:=%s", out.Name, out.Addr, err)
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
