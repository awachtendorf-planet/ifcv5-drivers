package vc3000

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"

	driverRecord "github.com/weareplanet/ifcv5-drivers/vc3000/record"
	"github.com/weareplanet/ifcv5-drivers/vc3000/template"

	"github.com/pkg/errors"
)

func (d *Dispatcher) constructPacketSocket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)
	station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketRegister:

		switch context.(type) {

		case driverRecord.Register:

			p := context.(driverRecord.Register)
			packet.Add("License", d.formatString(p.LicenseCode, 19))
			packet.Add("AppName", d.formatString(p.AppName, 19))
			packet.Add("Len", d.formatUint32(uint32(62-18))) // total packet size - header size

		default:

			log.Warn().Msgf("%T template '%s' with context '%T' not supported", d, packetName, context)
		}

	case template.PacketCodeCard, template.PacketCodeCardModify, template.PacketReadKey:

		command := d.getCodeCardCommand(addr, packetName, context)
		if command == 0 {
			return nil, errors.New("command not supported")
		}

		packet.Add("FF", d.formatString(string(command), 1))

		dta := make([]byte, 511, 511)
		if err := d.constructData(addr, string(command), station, context, &dta); err != nil {
			return nil, err
		}
		packet.Add("Data", dta)

		destination := d.getDestination(context)
		packet.Add("Destination", d.formatString(destination, 2))
		packet.Add("Source", d.formatString("00", 2))
		packet.Add("Len", d.formatUint32(uint32(583-18))) // total packet size - header size

	case template.PacketCheckout, template.PacketCheckoutRmt:

		var command = 'B'
		packet.Add("FF", d.formatString(string(command), 1))

		dta := make([]byte, 511, 511)
		if err := d.constructData(addr, string(command), station, context, &dta); err != nil {
			return nil, err
		}
		packet.Add("Data", dta)

		if packetName == template.PacketCheckoutRmt {
			destination := d.getDestination(context)
			packet.Add("Destination", d.formatString(destination, 2))
			packet.Add("Source", d.formatString("00", 2))
		}
		packet.Add("Len", d.formatUint32(uint32(577-18))) // total packet size - header size
	}

	if packetName != template.PacketRegister {
		packet.Add("OpID", d.formatString(appName, 9))
		packet.Add("OpFirst", d.formatString("John", 15))
		packet.Add("OpLast", d.formatString("Doe", 15))
	}

	packet.Add("Version", []byte{0x01, 0x0})
	return packet, nil
}
