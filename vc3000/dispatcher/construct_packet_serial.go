package vc3000

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/vc3000/template"

	"github.com/pkg/errors"
)

func (d *Dispatcher) constructPacketSerial(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketEnq:
		return packet, nil

	case template.PacketCodeCard, template.PacketCodeCardModify, template.PacketCheckout:

		command := d.getCodeCardCommand(addr, packetName, context)
		if command == 0 {
			return nil, errors.New("command not supported")
		}

		station, _ := d.GetStationAddr(addr)

		packet.Add("FF", d.formatString(string(command), 1))

		dta := make([]byte, 0)
		if err := d.constructData(addr, string(command), station, context, &dta); err != nil {
			return nil, err
		}
		if dta != nil && len(dta) > 0 {
			packet.Add("Data", dta)
		}

		destination := d.getDestination(context)
		packet.Add("Destination", d.formatString(destination, 2))
		packet.Add("Source", d.formatString("00", 2))

	case template.PacketReadKey:

		command := d.getCodeCardCommand(addr, packetName, context)
		if command == 0 {
			return nil, errors.New("command not supported")
		}

		packet.Add("FF", d.formatString(string(command), 1))

		destination := d.getDestination(context)
		packet.Add("Destination", d.formatString(destination, 2))
		packet.Add("Source", d.formatString("00", 2))

	default:

		log.Warn().Msgf("%T template '%s' handler not defined", d, packetName)
	}

	return packet, nil
}
