package dummy

import (
	"github.com/weareplanet/ifcv5-drivers/dummy/template"
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
			//d.sendRefusion(in.Addr)
			break

		case template.PacketUnknown:
			d.sendRefusion(in.Addr)

		default:
			d.sendAcknowledge(in.Addr)
		}
	}

	// uncomment if we can decode the entire framed packet
	// usually works with one data blob, eg. fidserv, visionline, vc3000
	// otherwise the function should be called in driver automate's

	/*

		// do we have a encoding configured
		encoding := d.GetEncoding(in.Addr)
		if len(encoding) == 0 {
			return
		}

		// do we have a payload packet
		raw := in.Data()
		data, exist := raw[payload]
		if !exist {
			return
		}

		if dec, err := d.Decode(data, encoding); err == nil && len(dec) > 0 {
			in.Add(payload, dec)
		} else if err != nil && len(encoding) > 0 {
			log.Warn().Msgf("%s", err)
		}

	*/

}

// ProcessOutgoingPacket callback, encoding encode
func (d *Dispatcher) ProcessOutgoingPacket(out *ifc.LogicalPacket) {
	if out == nil {
		return
	}

	// uncomment if we can encode the entire framed packet
	// usually works with one data blob, eg. fidserv, visionline, vc3000
	// otherwise the function should be called in dispatcher construct packet

	/*

		// do we have a encoding configured
		encoding := d.GetEncoding(out.Addr)
		if len(encoding) == 0 {
			return
		}

		// do we have a payload packet
		raw := out.Data()
		data, exist := raw[payload]
		if !exist {
			return
		}

		if enc, err := d.Encode(data, encoding); err == nil && len(enc) > 0 {
			out.Add(payload, enc)
		} else if err != nil && len(encoding) > 0 {
			log.Warn().Msgf("%s", err)
		}

	*/
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
