package xmlpos

import (
	"github.com/weareplanet/ifcv5-drivers/xmlpos/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"

	"github.com/pkg/errors"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	switch packetName {

	case template.PacketFramed:
		packet.Context = context

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}
