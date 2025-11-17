package ilco

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/weareplanet/ifcv5-drivers/ilco/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	station, _ := d.GetStationAddr(addr)

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq, template.PacketSof:
		return packet, nil
	}

	guest, ok := context.(*record.Guest)
	if !ok {
		return packet, errors.Errorf("context '%T' not supported", context)
	}

	switch packetName {

	case template.PacketCodeCard:

		withGateway := d.IsGatewayMode(station)

		displayLine := d.getDisplayLine(guest)

		keyType := d.GetKeyType(guest)

		cardCount := d.getKeyCount(guest)

		authNumber := d.GetAuthNumber(station)

		reservationNumber := d.getReservationID(guest)

		departureDate := d.getFormattedDate(guest.Reservation.DepartureDate)

		rooms := d.getRooms(station, guest)
		sortedRoomList := sortRooms(rooms)

		allowedAuthorisations := d.getAccessPoints(station, guest)

		msgLength := 8
		var payload []byte

		payload = d.constructRecord(payload, KeyDisplayLine, displayLine)

		for _, room := range sortedRoomList {
			payload = d.constructRecord(payload, KeyRoomNumber, room)
		}

		payload = d.constructRecord(payload, KeyCheckoutDate, departureDate)
		payload = d.constructRecord(payload, KeyCardType, keyType)
		payload = d.constructRecord(payload, KeyNumberKeys, cardCount)
		payload = d.constructRecord(payload, KeyAreas, allowedAuthorisations)
		if len(authNumber) > 0 {
			payload = d.constructRecord(payload, KeyAuthNo, authNumber)
		}
		payload = d.constructRecord(payload, KeyFolioNo, reservationNumber)

		track2Data := ""
		if d.SendTrack2(station) {

			track2Data = d.getTrack2Data(guest)
			payload = d.constructRecord(payload, KeyTrack2, track2Data)
		}

		if withGateway {

			encoder := d.getEncoder(guest)

			// Durch PreCheck 0 <= encoder <= 24
			encoderString := fmt.Sprintf("%02d", encoder)

			packet.Add("Addr", []byte(encoderString))

			msgLength += 2

		} else {

			packet.Add("Addr", []byte{})
		}

		packet.Add("Payload", payload)
		msgLength += len(payload)

		formattedLength := formatLength(msgLength)

		packet.Add("Len", []byte(formattedLength))

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func formatLength(length int) string {
	return fmt.Sprintf("%04d", length)
}

func sortRooms(rooms []string) []string {
	roomsAsInts := []int{}
	sortedRooms := []string{}

	for _, room := range rooms {
		val, err := strconv.Atoi(room)
		if err != nil {
			continue
		}
		roomsAsInts = append(roomsAsInts, val)
	}

	sort.Ints(roomsAsInts)

	for _, room := range roomsAsInts {
		val := strconv.Itoa(room)

		sortedRooms = append(sortedRooms, val)
	}

	return sortedRooms
}

func (d *Dispatcher) getEncoder(guest *record.Guest) int {
	if encoder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
		return cast.ToInt(encoder)
	}
	return 0
}

func (d *Dispatcher) constructRecord(payload []byte, key byte, value string) []byte {
	if len(value) > 0 {
		msg := []byte{key}
		msg = append(msg, []byte(value)...)
		msg = append(msg, FieldSeparator)

		payload = append(payload, msg...)
	}
	return payload
}

func (d *Dispatcher) getDisplayLine(guest *record.Guest) string {
	displayName := ""
	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {
		displayName = "2,"
	} else if len(guest.DisplayName) > 0 {
		displayName = "2," + guest.DisplayName
	} else if len(guest.FirstName) == 0 {
		displayName = "2," + guest.LastName
	} else {
		displayName = fmt.Sprintf("2,%s, %s", guest.LastName, guest.FirstName)
	}
	if len(displayName) > 20 {
		displayName = displayName[:20]
	}
	return displayName
}

func (d *Dispatcher) getReservationID(guest *record.Guest) string {
	reservationId := guest.Reservation.ReservationID
	folioWidth := 19
	reservationId = fmt.Sprintf("%0*s", folioWidth, reservationId)
	return reservationId
}

func (d *Dispatcher) getFormattedDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006/01/02 15:04")
}

func (d *Dispatcher) getRooms(_ uint64, guest *record.Guest) []string {

	var roomList []string

	mainRoom := guest.Reservation.RoomNumber

	roomList = append(roomList, mainRoom)

	data, exist := guest.GetGeneric(defines.AdditionalRooms)
	if !exist {
		return roomList
	}

	additionalRooms, ok := data.(string)
	if !ok || len(additionalRooms) == 0 {
		return roomList
	}

	rooms := strings.Split(additionalRooms, ";")
	for _, room := range rooms {

		if len(room) > 0 && len(room) <= RoomNumberLength && room != mainRoom {
			roomList = append(roomList, room)
		}

	}
	return roomList
}

func (d *Dispatcher) getKeyCount(guest *record.Guest) string {
	if count, exist := guest.GetGeneric(defines.KeyCount); exist {
		countInt := cast.ToInt(count)
		if countInt > 0 {
			if countInt < 256 {
				return cast.ToString(countInt)
			}
			return "256"
		}
	}
	return "1"
}

func (d *Dispatcher) getTrack2Data(guest *record.Guest) string {
	if track, exist := guest.GetGeneric(defines.Track2); exist {
		track2 := cast.ToString(track)
		return "202" + track2
	}
	return "202"
}

// getAccessPoints returns allowedAccessPoints
func (d *Dispatcher) getAccessPoints(station uint64, guest *record.Guest) string {
	keyOptions := []byte{}

	keyDefault := []byte(d.GetKeyOptionsDefault(station))

	if len(keyDefault) > 0 {

		keyOptions = keyDefault
	} else {

		if keyOp, exist := guest.GetGeneric(defines.KeyOptions); exist {
			keyOpts := cast.ToString(keyOp)
			keyOptions = []byte(keyOpts)
		}
	}

	allowedAccessPoints := ""

	for index, letter := range keyOptions {

		if letter == '1' && index < 8 {
			if len(allowedAccessPoints) > 0 {

				allowedAccessPoints += ","
			}

			allowedAccessPoints += cast.ToString(index + 1)

		}
	}

	return allowedAccessPoints
}

// GetKeyType maps the pms KeyType and the vendor KeyType
func (d *Dispatcher) GetKeyType(guest *record.Guest) string {

	if keyType, exist := guest.GetGeneric(defines.KeyType); exist {

		switch cast.ToString(keyType) {

		case "N": // new
			return "0"

		case "D": // duplicate
			return "1"

		}
	}

	return "0"
}
