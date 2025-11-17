package tiger

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"

	// "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	// "github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/tiger/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")
	// normalise device address
	addr = ifc.DeviceAddress(addr)

	station, _ := d.GetStationAddr(addr)

	protocol := d.Protocol(station)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	switch packetName {

	case template.PacketRoomBasedCheckIn, template.PacketAdditionalGuest:
		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		if protocol < 3 {
			record := "LN"
			if protocol == 2 && guest.Reservation.SharedInd {
				record = "LS"
			}
			packet.Add("Record", []byte(record))

		}

		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)

		reservationID := d.formatString(guest.Reservation.ReservationID, d.GetReservationIDLength(station))

		if protocol == 3 {

			vipStatus := guest.VIPStatus
			groupID := guest.Reservation.Group

			checkInDate := guest.Reservation.ArrivalDate.Format("020106")
			checkInTime := guest.Reservation.ArrivalDate.Format("1504")

			vipStatus = d.formatString(vipStatus, 2)
			groupID = d.formatString(groupID, 4)

			packet.Add("Date", []byte(checkInDate))
			packet.Add("Time", []byte(checkInTime))

			packet.Add("VipStatus", []byte(vipStatus))
			packet.Add("GroupID", []byte(groupID))
			packet.Add("PinNumber", []byte(strings.Repeat(" ", 5)))

		}
		languageCode := d.getGuestLanguage(guest, station)

		useDN := d.UseDisplayname(station)
		name := ""
		if useDN {
			name = d.getDisplayLine(station, guest)
		}

		actionCode := "3"

		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("AccountNumber", []byte(strings.Repeat(" ", 6)))
		packet.Add("ActionCode", []byte(actionCode))
		packet.Add("ReservationNumber", []byte(reservationID))
		packet.Add("LanguageDescription", []byte(strings.Repeat(" ", 14)))
		packet.Add("LanguageCode", []byte(languageCode))
		packet.Add("GuestName", []byte(name))

	case template.PacketRoomBasedCheckOut:
		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		if protocol < 3 {
			record := "LN"
			if protocol == 2 && guest.Reservation.SharedInd {
				record = "LS"
			}
			packet.Add("Record", []byte(record))

		}
		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, d.GetReservationIDLength(station))

		if protocol == 3 {

			vipStatus := guest.VIPStatus
			groupID := guest.Reservation.Group

			stampy := time.Now()
			checkInDate := stampy.Format("020106")
			checkInTime := stampy.Format("1504")

			vipStatus = d.formatString(vipStatus, 2)
			groupID = d.formatString(groupID, 4)

			packet.Add("Date", []byte(checkInDate))
			packet.Add("Time", []byte(checkInTime))

			packet.Add("VipStatus", []byte(vipStatus))
			packet.Add("GroupID", []byte(groupID))
			packet.Add("PinNumber", []byte(strings.Repeat(" ", 5)))

		}
		languageCode := d.getGuestLanguage(guest, station)

		useDN := d.UseDisplayname(station)
		name := ""
		if useDN {
			name = d.getDisplayLine(station, guest)
		}
		actionCode := "2"

		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("AccountNumber", []byte(strings.Repeat(" ", 6)))
		packet.Add("ActionCode", []byte(actionCode))
		packet.Add("ReservationNumber", []byte(reservationID))
		packet.Add("LanguageDescription", []byte(strings.Repeat(" ", 14)))
		packet.Add("LanguageCode", []byte(languageCode))
		packet.Add("GuestName", []byte(name))

	case template.PacketInformationUpdate:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		if protocol < 3 {
			record := "SR"
			if protocol == 2 && guest.Reservation.SharedInd {
				record = "SS"
			}
			packet.Add("Record", []byte(record))

		}

		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, d.GetReservationIDLength(station))

		vipStatus := guest.VIPStatus
		groupID := guest.Reservation.Group

		languageCode := d.getGuestLanguage(guest, station)
		lamp := cast.ToInt(guest.Reservation.MessageLightStatus)

		useDN := d.UseDisplayname(station)
		name := ""
		if useDN {
			name = d.getDisplayLine(station, guest)
		}

		checkInDate := time.Now().Format("020106")
		checkInTime := time.Now().Format("1504")

		var actionCode string
		cos := guest.Rights.ClassOfService

		if cos > 0 {
			actionCode = "0"
		} else {
			actionCode = "2"
			if protocol < 3 {
				actionCode = "1"
			}
		}

		groupID = d.formatString(groupID, 4)

		vipStatus = d.formatString(vipStatus, 2)

		packet.Add("Date", []byte(checkInDate))
		packet.Add("Time", []byte(checkInTime))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("GroupID", []byte(groupID))
		packet.Add("AccountNumber", []byte(strings.Repeat(" ", 6)))
		packet.Add("FirstWake", []byte("9999"))
		packet.Add("SecondWake", []byte("9999"))
		packet.Add("RoomStatusCode", []byte("99"))
		packet.Add("MessageWaitingAction", []byte(cast.ToString(lamp)))
		packet.Add("Lvl9Occupancy", []byte(actionCode))
		packet.Add("VipStatus", []byte(vipStatus))
		packet.Add("ReservationNumber", []byte(reservationID))
		packet.Add("LanguageDescription", []byte(strings.Repeat(" ", 14)))
		packet.Add("LanguageCode", []byte(languageCode))
		packet.Add("GuestName", []byte(name))

	case template.PacketVideoRights:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}
		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)

		checkInDate := time.Now().Format("020106")
		checkInTime := time.Now().Format("1504")

		videorights := guest.Rights.Video

		viewBill := cast.ToInt(videorights > 0)

		packet.Add("Date", []byte(checkInDate))
		packet.Add("Time", []byte(checkInTime))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("AccountNumber", []byte(strings.Repeat(" ", 6)))

		packet.Add("ViewBill", []byte(cast.ToString(viewBill)))
		packet.Add("VideoCheckOut", []byte("0"))
		packet.Add("CFlag", []byte("0"))
		packet.Add("DFlag", []byte("0"))
		packet.Add("EFlag", []byte("0"))

	case template.PacketRoomtransfer:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}
		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, d.GetReservationIDLength(station))

		oldRoom := d.getOldRoom(guest, station)
		checkInDate := time.Now().Format("020106")
		checkInTime := time.Now().Format("1504")

		packet.Add("Date", []byte(checkInDate))
		packet.Add("Time", []byte(checkInTime))
		packet.Add("RoomNumberOld", []byte(oldRoom))
		packet.Add("ReservationID", []byte(reservationID))
		packet.Add("RoomNumberNew", []byte(roomNumber))

	case template.PacketDND:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}
		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)

		dnd := cast.ToInt(guest.Reservation.DoNotDisturb)

		checkInDate := time.Now().Format("020106")
		checkInTime := time.Now().Format("1504")

		packet.Add("Date", []byte(checkInDate))
		packet.Add("Time", []byte(checkInTime))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("DND", []byte(cast.ToString(dnd)))

	case template.PacketWakeUpSet:

		generic, ok := context.(*record.Generic)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}
		roomnumber, _ := generic.Get(defines.WakeExtension)

		roomNumber := d.formatRoomnumber(cast.ToString(roomnumber), station)

		wakeTime := time.Now()
		calltime := time.Now()

		waketime, exist := generic.Get(defines.Timestamp)
		if exist {

			wakeTime = cast.ToTime(waketime)
		}

		hhmmWake := wakeTime.Format("1504")
		ddmmyyWake := wakeTime.Format("020106")

		hhmmCall := calltime.Format("1504")
		ddmmyyCall := calltime.Format("020106")

		action := "1"

		packet.Add("DateStamp", []byte(ddmmyyCall))
		packet.Add("TimeStamp", []byte(hhmmCall))
		packet.Add("DateCall", []byte(ddmmyyWake))
		packet.Add("TimeCall", []byte(hhmmWake))
		packet.Add("PinNumber", []byte(strings.Repeat(" ", 5)))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("GroupID", []byte(strings.Repeat(" ", 4)))
		packet.Add("LanguageCode", []byte("XX"))
		packet.Add("ActionCode", []byte(action))

	case template.PacketWakeUpClear:

		generic, ok := context.(*record.Generic)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}
		roomnumber, _ := generic.Get(defines.WakeExtension)

		roomNumber := d.formatRoomnumber(cast.ToString(roomnumber), station)

		wakeTime := time.Now()
		calltime := time.Now()

		waketime, exist := generic.Get(defines.Timestamp)
		if exist {
			wakeTime = cast.ToTime(waketime)
		}

		hhmmWake := wakeTime.Format("1504")
		ddmmyyWake := wakeTime.Format("020106")

		hhmmCall := calltime.Format("1504")
		ddmmyyCall := calltime.Format("020106")

		if protocol < 3 {
			hhmmWake = "0000"
		}

		action := "2"

		packet.Add("DateStamp", []byte(ddmmyyCall))
		packet.Add("TimeStamp", []byte(hhmmCall))
		packet.Add("DateCall", []byte(ddmmyyWake))
		packet.Add("TimeCall", []byte(hhmmWake))
		packet.Add("PinNumber", []byte(strings.Repeat(" ", 5)))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("GroupID", []byte(strings.Repeat(" ", 4)))
		packet.Add("LanguageCode", []byte("XX"))
		packet.Add("ActionCode", []byte(action))

	case template.PacketMessageWaitingGuest:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}
		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)

		mwl := "1"

		if !guest.Reservation.MessageLightStatus {
			mwl = d.MessageOffCharacter(station)
		}

		date := time.Now().Format("020106")
		time := time.Now().Format("1504")

		packet.Add("Date", []byte(date))
		packet.Add("Time", []byte(time))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("ActionCode", []byte(mwl))
		packet.Add("MessageWaitingCount", []byte("00"))

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) EncodeString(addr string, data string) string {
	encoding := d.GetEncoding(addr)
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

func (d *Dispatcher) DecodeString(addr string, data []byte) string {
	encoding := d.GetEncoding(addr)
	if len(encoding) == 0 {
		return string(data)
	}
	if dec, err := d.Decode(data, encoding); err == nil && len(dec) > 0 {
		return string(dec)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return string(data)
}

func (d *Dispatcher) getDisplayLine(station uint64, guest *record.Guest) string {
	displayName := ""

	nameLen := 23
	adNameLen := d.AdjustnameLength(station)
	if adNameLen {
		nameLen = 25
	}

	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {

		displayName = strings.Repeat(" ", nameLen)
	} else if len(guest.DisplayName) > 0 {

		displayName = guest.DisplayName
	} else if len(guest.FirstName) == 0 {

		displayName = guest.LastName
	} else {

		displayName = fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
	}
	if len(displayName) > nameLen {

		displayName = displayName[:nameLen]
	} else if len(displayName) < nameLen {

		numOfSpaces := nameLen - len(displayName)

		displayName = displayName + strings.Repeat(" ", numOfSpaces)
	}
	return displayName
}

func (d *Dispatcher) formatRoomnumber(roomnumber string, station uint64) string {

	roomLen := d.GetRoomNameLength(station)
	if len(roomnumber) < roomLen {

		numOfSpaces := roomLen - len(roomnumber)
		roomnumber = strings.Repeat("0", numOfSpaces) + roomnumber

	}

	return roomnumber
}

func (d *Dispatcher) formatString(in string, length int) string {

	if len(in) < length {
		numOfSpaces := length - len(in)
		in = strings.Repeat(" ", numOfSpaces) + in

	} else if len(in) > length {
		in = in[0:length]
	}
	return in
}

func (d *Dispatcher) getGuestLanguage(guest *record.Guest, station uint64) string {

	var language string
	if guest == nil || len(guest.Language) == 0 {
		language = ""
	} else {
		language = strings.ToUpper(guest.Language)
	}

	if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil {

		return d.formatString(languageCode, 2)
	}
	return d.formatString(language, 2)
}

func (d *Dispatcher) getOldRoom(context interface{}, station uint64) string {

	if guest, ok := context.(*record.Guest); ok {

		if data, exist := guest.GetGeneric(defines.OldRoomName); exist {

			oldRoom := cast.ToString(data)
			return d.formatRoomnumber(oldRoom, station)
		}
	}
	return ""
}

func (d *Dispatcher) GetRoomNameLength(station uint64) int {
	protocol := d.Protocol(station)
	roomLen := 6
	if d.ShortRoomname(station) || protocol < 3 {
		roomLen = 5
	}
	return roomLen
}

func (d *Dispatcher) GetReservationIDLength(station uint64) int {
	protocol := d.Protocol(station)
	reservationIDLength := 10
	if protocol == 2 {
		reservationIDLength = 8
	} else if protocol == 1 {
		reservationIDLength = 6
	}
	return reservationIDLength
}
