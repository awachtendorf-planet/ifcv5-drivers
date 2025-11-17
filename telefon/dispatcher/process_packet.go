package telefon

import (
	"github.com/weareplanet/ifcv5-drivers/telefon/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"

	"github.com/pkg/errors"
)

// ProcessIncomingPacket callback, ack/nak
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket, slot uint, garbage bool) {

	if in == nil {
		return
	}

	if d.IsAcknowledgement(in.Addr, in.Name) {

		switch in.Name {

		case template.Ack, template.Nak:
			break

		case template.Enq:

			if d.sendENQReply(in.Addr, in.Name, slot) {
				d.sendPacket(in.Addr, template.EnqReply, slot)
			}

		default:

			if garbage {
				d.sendRefusion(in.Addr, slot)
				return
			}

			d.sendAcknowledge(in.Addr, slot)

		}
	}

}

func (d *Dispatcher) sendAcknowledge(addr string, slot uint) {
	d.sendPacket(addr, template.Ack, slot)
}

func (d *Dispatcher) sendRefusion(addr string, slot uint) {
	d.sendPacket(addr, template.Nak, slot)
}

func (d *Dispatcher) sendPacket(addr string, name string, slot uint) error {

	if !d.parser.ExistOutgoingTemplate(slot, name) {
		return errors.Errorf("outgoing template '%s' slot: %d does not exist", name, slot)
	}

	driver := d.Network
	data := ifc.NewLogicalPacket(name, addr, "")
	err := driver.Send(data, 0)

	return err
}

func (d *Dispatcher) SendPacket(addr string, name string) error {
	slot := d.GetParserSlot(addr)
	err := d.sendPacket(addr, name, slot)
	return err
}
