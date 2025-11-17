package nortel

import (
	"github.com/weareplanet/ifcv5-drivers/nortel/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

// ProcessIncomingPacket callback, ack/nak
func (d *Dispatcher) ProcessIncomingPacket(in *ifc.LogicalPacket) {

	if in == nil {
		return
	}

	switch in.Name {

	case template.PacketAck, template.PacketNak, template.PacketEnq, template.PacketTermination, template.PacketGarbage, template.PacketError:
		// early return, IsAcknowledgment is a expensive operation
		return

	}

	if d.IsAcknowledgement(in.Addr, in.Name) {

		switch in.Name {

		case template.PacketUnknown:
			d.sendRefusion(in.Addr)

		default:
			d.sendAcknowledge(in.Addr)
		}
	}

}

// ProcessIncomingBytes remove MSB if set
func (d *Dispatcher) ProcessIncomingBytes(_ string, data *[]byte, size int) {

	if data == nil || size == 0 {
		return
	}

	raw := *data
	for i := 0; i < size; i++ {
		if raw[i] >= 0x80 {
			raw[i] ^= 0x80
		}
	}
}

// ProcessOutgoingPacket callback
func (d *Dispatcher) ProcessOutgoingBytes(addr string, data *[]byte, size int) {

	if data == nil || size < 3 {
		return
	}

	if d.IsBackgroundTerminalMode(addr) {

		// change "STX Payload ETX" to "Payload CRLF"

		raw := *data
		raw = append(raw[:0], raw[1:size-1]...)  // remove first and last byte (STX,ETX)
		raw = append(raw, []byte{0x0d, 0x0a}...) // append CR and LF

		// the size of the object does not change in total
		// so we dont need to resize the object

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
