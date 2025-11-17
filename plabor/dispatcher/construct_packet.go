package plabor

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	// "github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/plabor/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	station, _ := d.GetStationAddr(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq, template.PacketEot, template.PacketDBSyncFrameStart, template.PacketDBSyncFrameStopp:
		return packet, nil

	}

	switch packetName {

	case template.PacketCheckIn, template.PacketDataChangeUpdate:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, 7)

		name := d.getDisplayLine(addr, guest)
		checkInDate := guest.Reservation.ArrivalDate.Format("02.01.2006")
		checkOutDate := guest.Reservation.DepartureDate.Format("02.01.2006")

		videoRights := guest.Rights.Video
		tvRights := guest.Rights.TV

		remoteCheckoutFlag := "0"
		billViewFlag := cast.ToInt(videoRights > 0)
		tvProgrammFlag := cast.ToInt(tvRights > 0)
		standVideoFlag := cast.ToInt(tvRights == 1 || tvRights == 3)
		adultVideoFlag := cast.ToInt(tvRights == 3)
		seminarFlag := "0"
		switchOnTVFlag := "0"
		codeNrHighFlag := "0"
		codeNrLowFlag := "0"
		languageFlag := d.getGuestLanguage(guest, station)

		packet.Add("CIDate", []byte(checkInDate))
		packet.Add("RoomNumber", []byte(roomNumber))

		if d.RoomNumberAsReservationID(station) {
			packet.Add("AccountID", []byte(roomNumber))
		} else {
			packet.Add("AccountID", []byte(reservationID))
		}

		packet.Add("GuestName", []byte(name))
		packet.Add("CODate", []byte(checkOutDate))
		packet.Add("RemoteCheckoutFlag", []byte(remoteCheckoutFlag))
		packet.Add("BillViewFlag", []byte(cast.ToString(billViewFlag)))
		packet.Add("TVProgrammFlag", []byte(cast.ToString(tvProgrammFlag)))
		packet.Add("StandardVideoFlag", []byte(cast.ToString(standVideoFlag)))
		packet.Add("AdultVideoFlag", []byte(cast.ToString(adultVideoFlag)))
		packet.Add("SeminarFlag", []byte(seminarFlag))
		packet.Add("SwitchOnTVFlag", []byte(switchOnTVFlag))
		packet.Add("CodeNrHighFlag", []byte(codeNrHighFlag))
		packet.Add("CodeNrLowFlag", []byte(codeNrLowFlag))
		packet.Add("LanguageFlag", []byte(languageFlag))

	case template.PacketDataChangeMove:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, 7)

		roomnumberold := d.getOldRoom(context, station)

		packet.Add("RoomNumberOld", []byte(roomnumberold))

		if d.RoomNumberAsReservationID(station) {
			packet.Add("AccountID", []byte(roomnumberold))
		} else {
			packet.Add("AccountID", []byte(reservationID))
		}

		packet.Add("RoomNumberNew", []byte(roomNumber))

	case template.PacketCheckOut:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, 7)

		checkOutDate := guest.Reservation.DepartureDate.Format("02.01.2006")

		packet.Add("CODate", []byte(checkOutDate))
		packet.Add("RoomNumber", []byte(roomNumber))
		if d.RoomNumberAsReservationID(station) {
			packet.Add("AccountID", []byte(roomNumber))
		} else {
			packet.Add("AccountID", []byte(reservationID))
		}

	case template.PacketWakeupDirect, template.PacketWakeupIndirect:

		generic, ok := context.(*record.Generic)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}

		wakeTime := time.Now()
		roomnumber, _ := generic.Get(defines.WakeExtension)

		roomNumber := d.formatRoomnumber(cast.ToString(roomnumber), station)

		waketime, exist := generic.Get(defines.Timestamp)
		if exist {
			wakeTime = cast.ToTime(waketime)
		}

		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("Time", []byte(wakeTime.Format("15:04")))

		packet.Add("Date", []byte(wakeTime.Format("02.01.2006")))

	case template.PacketWakeupDirectk, template.PacketWakeupIndirectk:

		generic, ok := context.(*record.Generic)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}

		wakeTime := time.Now()
		roomnumber, _ := generic.Get(defines.WakeExtension)

		roomNumber := d.formatRoomnumber(cast.ToString(roomnumber), station)

		waketime, exist := generic.Get(defines.Timestamp)
		if exist {
			wakeTime = cast.ToTime(waketime)
		}

		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("Time", []byte(wakeTime.Format("15:04")))

		packet.Add("Date", []byte(wakeTime.Format("02.01")))

	case template.PacketBillPart:

		item, ok := context.(record.BillPreviewItem)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}

		roomNumber := d.formatRoomnumber(item.RoomNumber, station)

		booktime := item.Time.Format("02.01.2006")

		accountID := d.formatString(item.ReservationID, 7)

		value := d.formatString(cast.ToString(item.Amount), 9)

		orderText := d.formatString(d.EncodeString(addr, item.Description), 40)

		articleNumber := d.formatString(cast.ToString(item.FolioWindowNumber), 2)

		packet.Add("Date", []byte(booktime))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("AccountID", []byte(accountID))
		packet.Add("ArticleNumber", []byte(articleNumber))
		packet.Add("OrderText", []byte(orderText))
		packet.Add("Value", []byte(value))

	case template.PacketBalance:

		balance, ok := context.(record.BillPreviewBalance)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}

		roomNumber := d.formatRoomnumber(balance.RoomNumber, station)

		booktime := time.Now().Format("02.01.2006")

		accountID := d.formatString(balance.ReservationID, 7)

		value := d.formatString(cast.ToString(balance.Amount), 9)

		packet.Add("Date", []byte(booktime))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("AccountID", []byte(accountID))
		packet.Add("Value", []byte(value))

	case template.PacketMessageBlock, template.PacketMessageEnd:

		messagePart, ok := context.(*record.GuestMessage)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}

		roomNumber := d.formatRoomnumber(messagePart.RoomNumber, station)

		msgTime := messagePart.Time

		accountID := d.formatString(messagePart.ReservationID, 7)

		messageText := d.formatString(d.EncodeString(addr, messagePart.Text), 40)

		packet.Add("Date", []byte(msgTime.Format("02.01.2006")))
		packet.Add("RoomNumber", []byte(roomNumber))
		packet.Add("AccountID", []byte(accountID))
		packet.Add("Time", []byte(msgTime.Format("15:04:05")))
		packet.Add("Text", []byte(messageText))

	case template.PacketMessageSignal, template.PacketMessageDelete:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("context '%T' not supported", context)
		}

		roomNumber := d.formatRoomnumber(guest.Reservation.RoomNumber, station)
		reservationID := d.formatString(guest.Reservation.ReservationID, 7)

		packet.Add("RoomNumber", []byte(roomNumber))

		if d.RoomNumberAsReservationID(station) {

			packet.Add("AccountID", []byte(roomNumber))
		} else {

			packet.Add("AccountID", []byte(reservationID))
		}
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

func (d *Dispatcher) getDisplayLine(addr string, guest *record.Guest) string {
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

	return d.formatString(d.EncodeString(addr, displayName), 40)
}

func (d *Dispatcher) formatRoomnumber(roomnumber string, station uint64) string {
	if len(roomnumber) < 4 {

		numOfSpaces := 4 - len(roomnumber)
		if d.LeadingZeros(station) {

			roomnumber = strings.Repeat("0", numOfSpaces) + roomnumber
		} else {

			roomnumber = strings.Repeat(" ", numOfSpaces) + roomnumber
		}
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
	if guest == nil || len(guest.Language) == 0 {
		return ""
	}
	language := strings.ToUpper(guest.Language)
	if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil {
		return languageCode
	}
	return guest.Language
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
