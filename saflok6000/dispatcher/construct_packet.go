package saflok6000

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"

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

	case template.PacketAlive:
		return packet, nil
	}

	guest, ok := context.(*record.Guest)
	if !ok {
		return packet, errors.Errorf("context '%T' not supported", context)
	}

	encoder := d.getEncoder(guest)

	roomLength := d.GetRoomLength(station)

	// global
	packet.Add("Terminal", d.formatNumeric(encoder, 3))
	packet.Add("KeyNumber", d.formatString(d.getRoom(station, guest), roomLength))
	packet.Add("KeyLevel", d.formatNumeric(d.getKeyLevel(station), 1))
	packet.Add("Password", d.formatString(d.getPassword(station, guest), 7))

	switch packetName {

	case template.PacketKeyRequest:

		packet.Add("EncoderStation", d.formatNumeric(encoder, 2))
		packet.Add("TXC", d.formatNumeric(d.getKeyType(guest), 3))
		packet.Add("KeyCount", d.formatNumeric(d.getKeyCount(guest), 2))
		packet.Add("CheckoutDate", d.formatDate(guest.Reservation.DepartureDate))
		packet.Add("CheckoutTime", d.formatTime(guest.Reservation.DepartureDate))
		packet.Add("KeyExpireDate", d.formatDate(guest.Reservation.DepartureDate))
		packet.Add("KeyExpireTime", d.formatTime(guest.Reservation.DepartureDate))
		packet.Add("PassNumberOption", d.formatNumeric(d.getPassNumberOption(station, guest), 1))
		packet.Add("PassNumber", d.getPassNumber(guest))
		packet.Add("TrackData", d.getTrackData(guest))

	case template.PacketKeyDelete:
		// nothing more to do

	case template.PacketBeaconResponse:
		// nothing more to do

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) GetKeyType(context interface{}) int {
	if guest, ok := context.(*record.Guest); ok {
		return d.getKeyType(guest)
	}
	return -1
}

func (d *Dispatcher) formatNumeric(value int, length int) []byte {
	data := cast.ToString(value)
	if len(data) < length {
		data = pad.Left(data, length, "0")
	}
	return []byte(data)
}

func (d *Dispatcher) formatString(data string, length int) []byte {
	if len(data) < length {
		data = pad.Right(data, length, " ")
	}
	return []byte(data)
}

func (d *Dispatcher) formatDate(t time.Time) []byte {
	data := t.Format("010206") // MMDDYY
	return []byte(data)
}

func (d *Dispatcher) formatTime(t time.Time) []byte {
	data := t.Format("1504") // HHMM
	return []byte(data)
}

func (d *Dispatcher) getEncoder(guest *record.Guest) int {
	if encoder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
		return cast.ToInt(encoder)
	}
	return 0
}

func (d *Dispatcher) getRoom(station uint64, guest *record.Guest) string {
	room := guest.Reservation.RoomNumber
	if !d.LeadingZeroes(station) {
		room = strings.TrimLeft(room, "0")
	}
	room = d.encode(station, room)
	return room
}

func (d *Dispatcher) getKeyCount(guest *record.Guest) int {
	if count, exist := guest.GetGeneric(defines.KeyCount); exist {
		return cast.ToInt(count)
	}
	return 1
}

func (d *Dispatcher) getKeyType(guest *record.Guest) int {

	if keyType, exist := guest.GetGeneric(defines.KeyType); exist {

		switch cast.ToString(keyType) {

		case "N": // new
			if guest.Reservation.SharedInd {
				return 3 // duplicate
			}
			return 1

		case "D": // duplicate
			return 3

		}
	}

	return 1
}

func (d *Dispatcher) getTrackData(guest *record.Guest) []byte {

	track := "*"
	data, _ := guest.GetGeneric(defines.Track2)
	t2 := cast.ToString(data)
	if len(t2) > 0 {
		track = track + "2#" + t2
	}
	return []byte(track)

}

func (d *Dispatcher) getKeyLevel(station uint64) int {

	data := d.GetConfig(station, "KeyLevel", "1")
	level := cast.ToInt(data)

	// 1 guest level
	// 2 connector level
	// 3 multi-connector level
	// 4 limited-use level

	if level > 0 && level < 5 {
		return level
	}

	return 1
}

func (d *Dispatcher) getPassNumberOption(station uint64, guest *record.Guest) int {

	state := d.GetConfig(station, "PassNumberAlwaysNull", "false") // bezieht sich auf das Feld PassNumberOption, 5 Euro in das "Opera macht das aber so" Schweindel.
	if cast.ToBool(state) {
		return 0
	}

	data, _ := guest.GetGeneric(defines.KeyOptions)
	ap := cast.ToString(data)
	for _, v := range ap {
		if v != '0' {
			return 1
		}
	}
	return 0
}

func (d *Dispatcher) getPassNumber(guest *record.Guest) []byte {

	reverse := func(s string) (result string) {
		for _, v := range s {
			result = string(v) + result
		}
		return
	}

	data, _ := guest.GetGeneric(defines.KeyOptions)
	ap := cast.ToString(data)
	ap = reverse(ap)
	ap = pad.Left(ap, 12, "0")

	return []byte(ap)
}

func (d *Dispatcher) getPassword(station uint64, guest *record.Guest) string {

	state := d.GetConfig(station, "PasswordFromUser", "false")
	if cast.ToBool(state) {
		data, _ := guest.GetGeneric(defines.UserID)
		pw := cast.ToString(data)
		if len(pw) > 0 && len(pw) < 8 {
			return pw
		}
	}

	data := d.GetConfig(station, "Password", "")
	pw := cast.ToString(data)
	if len(pw) < 8 {
		return pw
	}

	return "unknown"
}
