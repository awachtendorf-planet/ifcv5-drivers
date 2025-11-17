package detewe

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/detewe/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	// "github.com/weareplanet/ifcv5-main/log"
)

// ProcessIncomingPacket callback, ack/nak, encoding decode
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

// ProcessOutgoingPacket callback, encoding encode
func (d *Dispatcher) ProcessOutgoingPacket(out *ifc.LogicalPacket) {
	if out == nil {
		return
	}

}

// PostProcessOutgoing callback, frame the telegram
func (d *Dispatcher) PostProcessOutgoing(addr string, _ string, data *[]byte) {
	if data != nil && !strings.Contains(string(*data), "Admin") && len(*data) > 1 {
		if d.IsSerialDevice(addr) {

			high, low := d.CalcLRCHighLow(data)

			payload := *data

			newSlice := append([]byte{0x02}, payload...)
			newSlice = append(newSlice, []byte{high, low, 0x03}...)

			*data = newSlice

		} else {

			payload := *data
			if len(payload) > 3 {

				newLoad := append([]byte("DeTeWeCItoPBX:"), payload...)
				lenload := len(newLoad)
				high := (lenload & 0xFF00) >> 8
				low := lenload & 0x00FF

				newSlice := []byte{byte(high), byte(low)}

				newSlice = append(newSlice, newLoad...)

				*data = newSlice

			}
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
