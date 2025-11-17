package caracas

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/utils"

	// "github.com/weareplanet/ifcv5-main/utils/pad"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/caracas/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, jobAction order.Action, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)
	station, _ := d.GetStationAddr(addr)

	if d.GetNewVersionMode(station) {
		switch packetName {
		case template.PacketCheckIn, template.PacketCheckOut, template.PacketUpdateCheckInName:
			packetName = template.PacketDataChangeAdvanced

		case template.PacketMessageLight:
			packetName = template.PacketMessageLightAlternative

		}
	}

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	//station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	switch packetName {
	case template.PacketAck, template.PacketDatetime:
		// Do Nothing

	default:
		extension := d.getExtension(station, context)
		packet.Add("Extension", []byte(extension))

	}

	switch packetName {

	case template.PacketCheckIn, template.PacketCheckOut, template.PacketDataChangeAdvanced:

		languageCode := d.getLanguageCode(station, context)
		guestName := d.getDisplayName(station, context)

		packet.Add("Language", []byte(languageCode))
		packet.Add("GuestName", []byte(guestName))

	case template.PacketMessageLight, template.PacketMessageLightAlternative:
		messageLight := d.getMessageLight(station, context)

		packet.Add("MessageLightStatus", []byte(messageLight))

	}

	switch packetName {

	case template.PacketBind:
		packet.Add("Status", []byte{0x31})

	case template.PacketDatetime:
		var now time.Time

		// get time zone from config

		if timezone := d.GetTimeZone(station); len(timezone) > 0 {

			if loc, err := time.LoadLocation(timezone); err == nil {
				now = time.Now().In(loc)
			} else {
				log.Warn().Msgf("%T addr '%s' get time zone failed, err=%s", d, addr, err)
			}

		}
		if now.IsZero() { // fallback

			now = time.Now()

		}

		datetime := now.Format("020106150405")
		dayOfWeek := now.Weekday()
		dowInt := int(dayOfWeek)
		if dowInt == 0 {
			dowInt = 7
		}

		packet.Add("DateTime", []byte(datetime))
		packet.Add("DayOfWeek", []byte(cast.ToString(dowInt)))

	case template.PacketDataChangeAdvanced:
		sharer := d.getSharer(station, context)
		order := d.getASWOrder(sharer, jobAction)
		syncInd := d.getSyncInd(station, context)
		reservationID := d.getReservationID(station, context)
		group := d.getGroupNumber(station, context)
		oldRoom := d.getOldRoomNumber(station, context)
		vip := d.getVIP(station, context)
		creditLimit := d.getCreditLimit(station, context)

		packet.Add("Order", []byte(order))
		packet.Add("SyncInd", []byte(syncInd))
		packet.Add("ReservationID", []byte(d.formatString(reservationID, 10)))
		packet.Add("Group", []byte(group))
		packet.Add("OldRoomNumber", []byte(oldRoom))
		packet.Add("VIP", []byte(vip))
		packet.Add("GroupName", []byte(d.formatString("", 40)))
		packet.Add("CreditLimit", []byte(creditLimit))

	case template.PacketMessageLightAlternative:
		messageMode := d.getMessageMode(station, context)

		packet.Add("Mode", []byte(messageMode))

	case template.PacketWakeupOrder:
		datetime := d.getWakeupDatetime(station, context)
		order := d.getWakeupOrder(jobAction)

		packet.Add("DateTime", []byte(datetime))
		packet.Add("Order", []byte(order))
		packet.Add("Mode", []byte{0x30})

	case template.PacketUpdateCheckInName:

		guestName := d.getDisplayName(station, context)

		packet.Add("GuestName", []byte(guestName))

	case template.PacketClassOfService:

		cos := d.getClassOfService(station, context)

		packet.Add("Status", []byte(cos))

	case template.PacketDND:

		dnd := d.getDND(station, context)

		packet.Add("Status", []byte(dnd))

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getASWOrder(sharer bool, jobAction order.Action) string {
	orderInt := 0
	switch jobAction {
	case order.Checkin:
		orderInt = 1
	case order.DataChange:
		orderInt = 3
	case order.Checkout:
		orderInt = 6
	}

	if orderInt != 0 && sharer && orderInt < 5 {
		orderInt++
	} else if orderInt == 6 && sharer {
		orderInt--
	}
	return cast.ToString(orderInt)
}

func (d *Dispatcher) getWakeupOrder(jobAction order.Action) string {
	orderInt := 0
	switch jobAction {
	case order.WakeupRequest:
		orderInt = 1
	}

	return cast.ToString(orderInt)
}

func (d *Dispatcher) getWakeupDatetime(station uint64, context interface{}) string {

	wakeTime := time.Now()
	switch representer := context.(type) {
	case *record.Generic:
		waketime, exist := representer.Get(defines.Timestamp)
		if exist {
			wakeTime = cast.ToTime(waketime)
		}
	}

	return wakeTime.Format("0201061504")
}

func (d *Dispatcher) getSyncInd(station uint64, context interface{}) string {

	switch representer := context.(type) {

	case *record.Guest:
		syncBool, _ := representer.GetGeneric(defines.SyncInd)
		syncInt := cast.ToInt(syncBool)
		return cast.ToString(syncInt)
	}

	return "0"
}

func (d *Dispatcher) getDisplayName(station uint64, context interface{}) string {
	displayName := ""
	maxLen := 40

	switch representer := context.(type) {

	case *record.Guest:

		if representer == nil || (len(representer.DisplayName) == 0 && len(representer.LastName) == 0 && len(representer.FirstName) == 0) {

			displayName = strings.Repeat(" ", maxLen)
		} else if len(representer.DisplayName) > 0 {

			displayName = representer.DisplayName
		} else if len(representer.FirstName) == 0 {

			displayName = representer.LastName
		} else {

			displayName = fmt.Sprintf("%s, %s", representer.LastName, representer.FirstName)
		}
	}

	return d.formatString(d.EncodeString(station, displayName), maxLen)
}

func (d *Dispatcher) formatString(in string, length int) string {

	if len(in) < length {
		numOfSpaces := length - len(in)
		in = in + strings.Repeat(" ", numOfSpaces)

	} else if len(in) > length {
		in = in[0:length]
	}
	return in
}

func (d *Dispatcher) getSharer(station uint64, context interface{}) bool {

	switch representer := context.(type) {

	case *record.Guest:

		return representer.Reservation.SharedInd
	}

	log.Warn().Msgf("Station %d no guest record found, sharer = false", station)

	return false
}

func (d *Dispatcher) getMessageLight(station uint64, context interface{}) string {

	switch representer := context.(type) {

	case *record.Guest:

		messageLightStatus := representer.Reservation.MessageLightStatus

		voiceMail := representer.Reservation.Voicemail
		vmail := cast.ToString(voiceMail)
		mlt := cast.ToBool(messageLightStatus)

		if mlt || vmail == "Y" {
			return "1"
		}

	}

	return "0"
}

func (d *Dispatcher) getClassOfService(station uint64, context interface{}) string {
	classOfService := "5"
	switch representer := context.(type) {

	case *record.Guest:

		cos := representer.Rights.ClassOfService

		switch cos {
		case 0:
			classOfService = "5"
		case 1:
			classOfService = "4"
		case 2:
			classOfService = "3"
		case 3:
			classOfService = "2"
		}

	}

	return classOfService
}

func (d *Dispatcher) getDND(station uint64, context interface{}) string {
	doNotDisturb := "0"
	switch representer := context.(type) {

	case *record.Guest:

		dnd := representer.Reservation.DoNotDisturb
		if dnd {
			doNotDisturb = "1"
		}

	}

	return doNotDisturb
}

func (d *Dispatcher) getMessageMode(station uint64, context interface{}) string {
	mode := 0
	switch representer := context.(type) {

	case *record.Guest:

		messageLightStatus := representer.Reservation.MessageLightStatus

		mlt := cast.ToBool(messageLightStatus)
		if mlt {
			mode++
		}

		voiceMail := representer.Reservation.Voicemail

		vmail := cast.ToString(voiceMail)
		if vmail == "Y" {
			mode++
		}

	}

	return cast.ToString(mode)
}

func (d *Dispatcher) getLanguageCode(station uint64, context interface{}) string {

	switch representer := context.(type) {

	case *record.Guest:

		language := strings.ToUpper(representer.Language)
		if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil {
			if len(languageCode) == 2 {
				return languageCode
			} else {
				log.Error().Msgf("station '%d' language code exeeds 2 digits '%s'", station, languageCode)
			}
		}
	}

	return "01"
}

func (d *Dispatcher) getGroupNumber(station uint64, context interface{}) string {
	var group string
	switch representer := context.(type) {

	case *record.Guest:
		group = d.EncodeString(station, representer.Reservation.Group)
	}

	return d.formatString(group, 10)
}

func (d *Dispatcher) getOldRoomNumber(station uint64, context interface{}) string {
	var oldRoom string
	switch representer := context.(type) {

	case *record.Guest:
		olRoom, _ := representer.GetGeneric(defines.OldRoomName)
		oldRoom = d.EncodeString(station, cast.ToString(olRoom))
	}

	return d.formatString(oldRoom, 6)
}

func (d *Dispatcher) getVIP(station uint64, context interface{}) string {
	var vip string
	switch representer := context.(type) {

	case *record.Guest:
		vip = d.EncodeString(station, representer.VIPStatus)
	}

	return d.formatString(vip, 3)
}

func (d *Dispatcher) getCreditLimit(station uint64, context interface{}) string {
	var creditLimit string
	switch representer := context.(type) {

	case *record.Guest:
		creditLimit = cast.ToString(representer.Reservation.CreditLimit)
	}

	return d.formatString(creditLimit, 10)
}

func (d *Dispatcher) getReservationID(station uint64, context interface{}) string {
	var reservationID string
	switch representer := context.(type) {

	case *record.Guest:
		reservationID = d.EncodeString(station, representer.Reservation.ReservationID)
	}

	return reservationID
}

func (d *Dispatcher) getExtension(station uint64, context interface{}) string {
	var extension string

	extension = d.getRoom(station, context)

	return d.formatString(extension, 6)
}

func (d *Dispatcher) getRoom(station uint64, context interface{}) string {
	var room string

	switch representer := context.(type) {

	case *record.Guest:
		room = d.EncodeString(station, representer.Reservation.RoomNumber)

	case *record.Generic:
		roomNo, _ := representer.Get(defines.WakeExtension)
		roomNumber := cast.ToString(roomNo)
		room = d.EncodeString(station, roomNumber)
	}
	return room
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
