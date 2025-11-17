package informatel

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/informatel/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, jobAction order.Action, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketAck, template.PacketNak:
		return packet, nil

	}
	extension := d.GetRoom(station, context)

	extensionLength := 10

	if packetName == template.PacketRoomMove {
		extensionLength = 8
	}
	packet.Add("Extension", []byte(d.formatString(station, extension, extensionLength, false)))

	if packetName != template.PacketWakeupOutgoingCall {

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Guest'", context)
		}
		switch packetName {

		case template.PacketCheckIn, template.PacketNameChange:

			name := d.getDisplayName(station, guest)
			vip := d.getVIPCode(station, guest)
			languageCode := d.getLanguageCode(station, guest)

			packet.Add("Name", []byte(d.formatString(station, name, 20, false)))
			packet.Add("VIP", []byte(d.formatString(station, vip, 3, true)))
			packet.Add("LanguageCode", []byte(d.formatString(station, languageCode, 1, false)))

		case template.PacketCheckOut:
			//

		case template.PacketRoomMove:

			oldRoomnumber := d.getOldRoom(station, guest)

			packet.Add("OldRoomnumber", []byte(d.formatString(station, oldRoomnumber, 10, false)))

		case template.PacketMessageLightStatus:

			messageLightStatus := d.getMessageLightStatus(guest)

			packet.Add("MessageStatus", []byte(d.formatString(station, messageLightStatus, 3, false)))

		case template.PacketDND:

			dnd := d.getDND(guest)

			packet.Add("DND", []byte(d.formatString(station, dnd, 3, false)))

		case template.PacketClassOfService:

			cos := d.getClassOfService(guest)

			packet.Add("ClassOfService", []byte(d.formatString(station, cos, 3, false)))

		default:
			return packet, errors.Errorf("packet '%s' handler not defined", packetName)

		}
	} else {

		// WakeupOutgoing

		wakeupAction := "PSV" // PSV == SET, ASV == CLEAR

		if jobAction == order.WakeupClear {
			wakeupAction = "ASV"
		}

		wakeupTime := d.getWakeupTime(context)

		packet.Add("WakeupStatus", []byte(d.formatString(station, wakeupAction, 3, false)))
		packet.Add("WakeupTime", []byte(d.formatString(station, wakeupTime, 4, false)))

	}

	return packet, nil
}

func (d *Dispatcher) GetRoom(station uint64, context interface{}) string {
	var extension string

	switch representer := context.(type) {

	case *record.Guest:
		extension = d.EncodeString(station, representer.Reservation.RoomNumber)

	case *record.Generic:
		roomNo, _ := representer.Get(defines.WakeExtension)
		roomNumber := cast.ToString(roomNo)
		extension = d.EncodeString(station, roomNumber)
	}

	return extension
}

func (d *Dispatcher) getOldRoom(station uint64, guest *record.Guest) string {

	if data, exist := guest.GetGeneric(defines.OldRoomName); exist {

		oldRoom := cast.ToString(data)
		return d.EncodeString(station, oldRoom)
	}
	return ""
}

func (d *Dispatcher) getDisplayName(station uint64, guest *record.Guest) string {
	displayName := ""
	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {

		displayName = strings.Repeat(" ", 40)
	} else if len(guest.DisplayName) > 0 {

		displayName = guest.DisplayName
	} else if len(guest.FirstName) == 0 {

		displayName = guest.LastName
	} else {

		displayName = fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
	}

	return d.EncodeString(station, displayName)
}

func (d *Dispatcher) getLanguageCode(station uint64, guest *record.Guest) string {
	key := guest.Language
	if match, err := d.GetMapping("languagecode", station, key, false); err == nil {
		return d.EncodeString(station, match)
	}
	return key
}

func (d *Dispatcher) getVIPCode(station uint64, guest *record.Guest) string {
	if d.SendVIP(station) {
		return guest.VIPStatus

	}
	return " "
}

func (d *Dispatcher) getMessageLightStatus(guest *record.Guest) string {
	messageLightStatus := guest.Reservation.MessageLightStatus

	mlt := cast.ToBool(messageLightStatus)
	if mlt {
		return "ON"
	}

	return "OFF"
}

func (d *Dispatcher) getDND(guest *record.Guest) string {
	dnd := guest.Reservation.DoNotDisturb

	if dnd {
		return "ON"
	}

	return "OFF"
}

func (d *Dispatcher) getClassOfService(guest *record.Guest) string {
	cos := guest.Rights.ClassOfService

	switch cos {
	case 0:
		return "HSR"
	case 1:
		return "REG"
	case 2:
		return "NAT"
	case 3:
		return "INT"
	}

	return "HSR"
}

func (d *Dispatcher) getWakeupTime(context interface{}) string {

	timestamp := time.Now()

	generic := context.(*record.Generic)
	if wakeupTime, exist := generic.Get(defines.Timestamp); exist {
		if t, ok := wakeupTime.(time.Time); ok && !t.IsZero() {
			timestamp = t
		}
	}
	return timestamp.Format("1504")
}

func (d *Dispatcher) formatString(station uint64, data string, width int, encode bool) []byte {
	if encode {
		data = d.EncodeString(station, data)
	}
	if len(data) > width {
		data = data[0:width]
	}
	if len(data) == width {
		return []byte(data)
	}
	data = pad.Right(data, width, " ")
	return []byte(data)
}

func (d *Dispatcher) EncodeString(station uint64, data string) string {
	encoding := d.GetEncodingByStation(station)
	if len(encoding) == 0 {
		return data
	}
	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
		return string(enc)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return data
}
