package vc3000

import (
	"encoding/binary"
	"strings"

	"github.com/weareplanet/ifcv5-drivers/vc3000/template"
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
		case template.PacketAck, template.PacketNak:
			break
		case template.PacketGarbage, template.PacketUnknown:
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
	data, exist := raw["Data"]
	if !exist {
		return
	}

	if dec, err := d.Decode(data, encoding); err == nil && len(dec) > 0 {
		in.Add("Data", dec)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
}

// GetPayload return the Data object from a logical packet
func (d *Dispatcher) GetPayload(in *ifc.LogicalPacket) ([]byte, bool) {
	if in == nil {
		return []byte{}, false
	}
	data := in.Data()["Data"]
	return data, len(data) > 0
}

// GetUint ...
func (d *Dispatcher) GetUint(data []byte) uint32 {
	value := binary.LittleEndian.Uint32(data)
	return value
}

// GetField ...
func (d *Dispatcher) GetField(in *ifc.LogicalPacket, key string) ([]byte, bool) {
	var nothing []byte
	if in == nil {
		return nothing, false
	}
	if data, exist := in.Data()[key]; exist {
		if len(data) > 0 {
			return data, true
		}
	}
	return nothing, false
}

func (d *Dispatcher) GetPayLoadField(in *ifc.LogicalPacket, key string) (string, bool) {
	data, exist := d.GetPayload(in)
	if !exist {
		return "", false
	}
	record := string(data[:])
	offset := strings.Index(record, string(0x1e)+key)
	if offset > -1 {
		offset = offset + 1 + len(key)
		if offset >= len(record) {
			return "", true
		}
		record = string(data[offset:])
		if index := strings.Index(record, string(0x1e)); index > -1 {
			record = string(data[offset : offset+index])
		} else {
			record = strings.TrimRight(record, string(0x0))
		}
		return record, true
	}
	return "", false
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

func (d *Dispatcher) sendEnquiry(addr string) {
	driver := d.Network
	data := ifc.NewLogicalPacket(template.PacketEnq, addr, "")
	driver.Send(data, 0)
}
