package callstar

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/weareplanet/ifcv5-drivers/callstar/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	// send swap packet with Z+/Z- or as normal live packet
	switch packetName {

	case template.PacketCheckInSwap, template.PacketCheckOutSwap, template.PacketWakeupRequestSwap, template.PacketWakeupClearSwap:

		if d.SendSwapLabel(addr) {
			break
		}

		switch packetName {

		case template.PacketCheckInSwap:
			packetName = template.PacketCheckIn

		case template.PacketCheckOutSwap:
			packetName = template.PacketCheckOut

		case template.PacketWakeupRequestSwap:
			packetName = template.PacketWakeupRequest

		case template.PacketWakeupClearSwap:
			packetName = template.PacketWakeupClear

		case template.PacketRoomStatusSwap:
			packetName = template.PacketRoomStatus

		}

	}

	// data change or room move
	if packetName == template.PacketDataChange && d.IsMove(context) {
		packetName = template.PacketRoomMove
	}

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	station, _ := d.GetStationAddr(addr)

	var data string

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	case template.PacketCheckIn, template.PacketCheckInSwap:

		data = d.append(data, "I")

		name := d.constructName(station, context)
		if len(name) > 0 {
			data = d.append(data, "N"+name)
		}

	case template.PacketDataChange:

		name := d.constructName(station, context)
		if len(name) > 0 {
			data = d.append(data, "N"+name)
		}

	case template.PacketRoomMove:

		oldRoom := d.getOldRoom(context)
		oldRoom = d.normaliseString(oldRoom)
		oldRoom = d.EncodeString(station, oldRoom)

		packet.Add("OldRoom", []byte(oldRoom))

		name := d.constructName(station, context)
		data = d.append(data, name)

	case template.PacketRoomStatus, template.PacketRoomStatusSwap:

		if d.getMessageLightStatus(context) {
			data = d.append(data, "L+")
		} else {
			data = d.append(data, "L-")
		}

		// if d.getVoiceMailStatus(context) {
		// 	data = d.append(data, "V+")
		// } else {
		// 	data = d.append(data, "V-")
		// }

		if d.getClassOfService(context) {
			data = d.append(data, "U")
		} else {
			data = d.append(data, "B")
		}

	case template.PacketCheckOut, template.PacketCheckOutSwap:

		data = d.append(data, "O")

	case template.PacketWakeupRequest, template.PacketWakeupRequestSwap:

		t := d.getWakeupTime(context)
		data = d.append(data, "A"+t)

	case template.PacketWakeupClear, template.PacketWakeupClearSwap:

		data = d.append(data, "A")

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	room := d.getRoom(context)
	room = d.normaliseString(room)
	room = d.EncodeString(station, room)

	packet.Add("Room", []byte(room))
	packet.Add("Data", []byte(data))

	return packet, nil
}

func (d *Dispatcher) append(data string, value string) string {
	data = data + value + "\\"
	return data
}

func (d *Dispatcher) shrink(data string, width int) string {
	if len(data) > width {
		return data[0:width]
	}
	return data
}

func (d *Dispatcher) getRoom(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Reservation.RoomNumber
	}
	if generic, ok := context.(*record.Generic); ok {
		if extension, exist := generic.Get(defines.WakeExtension); exist {
			return cast.ToString(extension)
		}
	}
	return ""
}

func (d *Dispatcher) getOldRoom(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		if data, exist := guest.GetGeneric(defines.OldRoomName); exist {
			oldRoom := cast.ToString(data)
			return oldRoom
		}
	}
	return ""
}

func (d *Dispatcher) getMessageLightStatus(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Reservation.MessageLightStatus
	}
	return false
}

// func (d *Dispatcher) getVoiceMailStatus(context interface{}) bool {
// 	if guest, ok := context.(*record.Guest); ok {
// 		status := strings.ToUpper(guest.Reservation.Voicemail)
// 		if len(status) == 0 || status == "N" || status == "0" {
// 			return false
// 		}
// 		return true
// 	}
// 	return false
// }

func (d *Dispatcher) getClassOfService(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Rights.ClassOfService > 0
	}
	return false
}

func (d *Dispatcher) getGuestLanguage(station uint64, guest *record.Guest) string {
	if guest == nil || len(guest.Language) == 0 {
		return ""
	}
	language := strings.ToUpper(guest.Language)
	if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil {
		return languageCode
	}
	return guest.Language
}

func (d *Dispatcher) getVIPCode(station uint64, guest *record.Guest) string {
	if guest == nil {
		return ""
	}
	status := strings.ToUpper(guest.VIPStatus)
	if code, err := d.GetMapping("vipstatus", station, status, false); err == nil {
		return code
	}
	return guest.VIPStatus
}

func (d *Dispatcher) getWakeupTime(context interface{}) string {
	if generic, ok := context.(*record.Generic); ok {
		if wakeupTime, exist := generic.Get(defines.Timestamp); exist {
			if t, ok := wakeupTime.(time.Time); ok && !t.IsZero() {
				return t.Format("1504") // hhmm
			}
		}
	}
	return ""
}

func (d *Dispatcher) normaliseString(data string) string {
	data = strings.Replace(data, "_", "-", -1)
	data = strings.Replace(data, "\\", "", -1)
	data = strings.Trim(data, " ")
	return data
}

func (d *Dispatcher) constructName(station uint64, context interface{}) string {

	guest, ok := context.(*record.Guest)
	if !ok {
		return ""
	}

	guestName := guest.LastName + "  " + guest.FirstName
	salutation := guest.Salutation

	guestName = d.normaliseString(guestName)
	salutation = d.normaliseString(salutation)

	guestName = d.EncodeString(station, guestName)
	salutation = d.EncodeString(station, salutation)

	language := d.getGuestLanguage(station, guest)
	status := d.getVIPCode(station, guest)
	group := guest.Reservation.Group
	id := guest.Reservation.ReservationID

	guestName = d.shrink(guestName, 24)
	salutation = d.shrink(salutation, 8)
	language = d.shrink(language, 3)
	status = d.shrink(status, 8)

	if len(group) > 6 {
		group = ""
	}

	name := fmt.Sprintf("%s_%s_%s_%s_%s_%s", guestName, salutation, language, status, group, id)

	return name
}
