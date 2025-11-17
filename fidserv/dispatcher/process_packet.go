package fidserv

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"
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
}

// ProcessOutgoingPacket callback, encoding encode
func (d *Dispatcher) ProcessOutgoingPacket(out *ifc.LogicalPacket) {
	if out == nil {
		return
	}

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
}

// GetPayload return the Data object from a logical packet
func (d *Dispatcher) GetPayload(in *ifc.LogicalPacket) (string, bool) {
	if in == nil {
		return "", false
	}
	data := in.Data()["Data"]
	record := string(data[:])
	record = strings.Trim(record, "|")
	return record, len(record) > 0
}

// GetField return the field record from a logical packet
func (d *Dispatcher) GetField(in *ifc.LogicalPacket, key string) (string, bool) {
	if in == nil {
		return "", false
	}
	data := in.Data()["Data"]
	record := string(data[:])
	offset := strings.Index(record, "|"+key)
	if offset > -1 {
		offset = offset + 1 + len(key)
		if offset >= len(record) {
			return "", true
		}
		record = string(data[offset:])
		if index := strings.Index(record, "|"); index > -1 {
			record = string(data[offset : offset+index])
		}
		return record, true
	}
	return "", false
}

// EachField call a function for each field
func (d *Dispatcher) EachField(in *ifc.LogicalPacket, fn func(key string, value string)) {
	if in == nil || fn == nil {
		return
	}
	data := in.Data()["Data"]
	record := string(data[:])

	records := strings.Split(record, "|")

	for i := range records {
		data := records[i]
		var key, value string
		if len(data) > 1 {
			key = data[:2]
			if len(data) > 2 {
				value = data[2:]
			}
		}
		if len(key) > 0 {
			fn(key, value)
		}
	}
}

// ExistField return true if the field exist in the logical packet
func (d *Dispatcher) ExistField(in *ifc.LogicalPacket, key string) bool {
	_, exist := d.GetField(in, key)
	return exist
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
