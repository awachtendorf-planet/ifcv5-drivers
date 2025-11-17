package saflok6000

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"

	"github.com/spf13/cast"
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

// GetEncoder ...
func (d *Dispatcher) GetTerminal(in *ifc.LogicalPacket) (int, bool) {
	if in != nil {
		n := d.getNumeric(in, "Terminal")
		return n, true
	}
	return -1, false
}

// GetStatusCode ...
func (d *Dispatcher) GetStatusCode(in *ifc.LogicalPacket) (int, bool) {
	if in != nil {
		n := d.getNumeric(in, "Status")
		return n, true
	}
	return -1, false
}

// GetResponseCode ...
func (d *Dispatcher) GetResponseCode(in *ifc.LogicalPacket) (int, bool) {
	if in != nil {
		n := d.getNumeric(in, "ResponseCode")
		return n, true
	}
	return -1, false
}

func (d *Dispatcher) getNumeric(in *ifc.LogicalPacket, name string) int {

	if in == nil {
		return -1
	}

	data := in.Data()[name]
	value := string(data)
	value = strings.TrimLeft(value, " ")
	value = strings.TrimLeft(value, "0")

	n := cast.ToInt(value)
	return n
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
