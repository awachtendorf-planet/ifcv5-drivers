package visionline

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/visionline/template"

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
	data := in.Data()[payload]
	record := string(data[:])
	//record = strings.Trim(record, ";")
	return record, len(record) > 0
}

// GetField return the field record from a logical packet
func (d *Dispatcher) GetField(in *ifc.LogicalPacket, key string) (string, bool) {
	if in == nil {
		return "", false
	}
	data := in.Data()[payload]
	records := strings.Split(string(data[:]), ";")
	keyLen := len(key)
	for i := range records {
		record := records[i]
		if len(record) >= keyLen && record[:keyLen] == key {
			return record[keyLen:], true
		}
	}
	return "", false
}

// GetResultCode ...
func (d *Dispatcher) GetResultCode(in *ifc.LogicalPacket) (int, bool) {
	if result, exist := d.GetField(in, "RC"); exist {
		code := cast.ToInt(result)
		return code, true
	}
	return -1, false
}

// GetEncoder ...
func (d *Dispatcher) GetEncoder(in *ifc.LogicalPacket) (int, bool) {
	if result, exist := d.GetField(in, "EA"); exist {
		code := cast.ToInt(result)
		return code, true
	}
	return -1, false
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
