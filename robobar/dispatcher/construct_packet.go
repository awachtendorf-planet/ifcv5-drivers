package robobar

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/robobar/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	case template.PacketRestart:
		return packet, nil

	case template.PacketStartup:
		sequence := "     "
		if stored := d.GetSequenceNumber(station); stored >= 0 {
			sequence = pad.Left(cast.ToString(stored), 5, "0")
		}
		packet.Add("SequenceNumber", []byte(sequence))
		packet.Add("Date", []byte(time.Now().Format("0601021504")))
		return packet, nil

	}

	room := d.GetRoom(context)
	//room = d.EncodeString(addr, room)
	if len(room) == 0 {
		return packet, errors.New("empty room")
	}

	roomLength := d.GetRoomLength(station)
	if roomLength > 0 && len(room) > roomLength {
		return packet, errors.Errorf("room '%s' too long (max: %d)", room, roomLength)
	}

	room = pad.Left(room, roomLength, "0")

	packet.Add("Room", []byte(room))

	switch packetName {

	case template.PacketCheckIn, template.PacketUpdateCheckIn:
		packet.Add("Command", []byte("I"))

	case template.PacketCheckOut, template.PacketUpdateCheckOut:
		packet.Add("Command", []byte("O"))

	case template.PacketLockBar, template.PacketUpdateLockBar:
		packet.Add("Command", []byte("L"))

	case template.PacketUnlockBar, template.PacketUpdateUnlockBar:
		packet.Add("Command", []byte("U"))

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) GetRoom(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		return strings.TrimLeft(guest.Reservation.RoomNumber, "0")
	}
	return ""
}

func (d *Dispatcher) GetMinibarRight(context interface{}) int {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Rights.Minibar
	}
	return -1
}

func (d *Dispatcher) IsSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}
	return false
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
