package dummy

import (
	// "strings"
	"time"

	// "github.com/weareplanet/ifcv5-main/ifc/record"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/utils"

	// "github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/dummy/template"

	"github.com/pkg/errors"
	// "github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	//station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	switch packetName {

	case template.PacketCheckIn, template.PacketCheckOut, template.PacketDataChange:
		packet.Add("Data", []byte(packetName))

	case template.PacketKeyRequest, template.PacketKeyDelete:
		packet.Add("Data", []byte(packetName))

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

// func (d *Dispatcher) EncodeString(addr string, data string) string {
// 	encoding := d.GetEncoding(addr)
// 	if len(encoding) == 0 {
// 		return data
// 	}
// 	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
// 		return string(enc)
// 	} else if err != nil && len(encoding) > 0 {
// 		log.Warn().Msgf("%s", err)
// 	}
// 	return data
// }

// func (d *Dispatcher) DecodeString(addr string, data []byte) string {
// 	encoding := d.GetEncoding(addr)
// 	if len(encoding) == 0 {
// 		return string(data)
// 	}
// 	if dec, err := d.Decode(data, encoding); err == nil && len(dec) > 0 {
// 		return string(dec)
// 	} else if err != nil && len(encoding) > 0 {
// 		log.Warn().Msgf("%s", err)
// 	}
// 	return string(data)
// }
