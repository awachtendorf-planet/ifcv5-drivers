package miwa

import (
	// "strings"
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/weareplanet/ifcv5-drivers/miwa/template"

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

	guest, ok := context.(*record.Guest)
	if !ok {
		return packet, errors.Errorf("context '%T' not supported", context)
	}

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil
	}

	encoderID := d.getEncoderID(guest)
	packet.Add("ID", []byte(encoderID))

	switch packetName {

	case template.PacketKeyCreate:

		keyType := d.GetCardType(guest)
		if keyType == "O" {
			return packet, errors.Errorf("keytype %s is not supported", keyType)
		}

		arrival := d.GetArrivalDate(guest)
		departure := d.GetDapartureDate(guest)
		room := d.getRoom(station, guest)
		count := d.GetCount(guest)
		rooms := d.getAdditionalRooms(station, guest)
		if len(rooms) < 5 {
			return packet, errors.Errorf("additional rooms array length is %d, need atleast a length of 3", len(rooms))
		}

		staffID := d.getOperator(guest)

		packet.Add("IssueType", []byte(keyType))
		packet.Add("CardType", []byte("GA"))
		packet.Add("CheckIn", []byte(arrival))    // YYMMDDHHMM
		packet.Add("CheckOut", []byte(departure)) // YYMMDDHHMM
		packet.Add("MainRoom", []byte(room))      // 8 bytes
		packet.Add("ExtraRoom1", []byte(rooms[0]))
		packet.Add("ExtraRoom2", []byte(rooms[1]))
		packet.Add("ExtraRoom3", []byte(rooms[2]))
		packet.Add("ExtraRoom4", []byte(rooms[3]))
		packet.Add("ExtraRoom5", []byte(rooms[4]))
		packet.Add("Reserve1", []byte((strings.Repeat(" ", 32)))) // filled with blanks
		packet.Add("SpecialFlag", []byte((strings.Repeat("0", 40))))
		packet.Add("IssueNumber", []byte(count))                 // nr of cards to be issued
		packet.Add("POSInfo", []byte((strings.Repeat(" ", 37)))) // filled with blanks
		packet.Add("StaffCode", []byte(staffID))
		packet.Add("Reserve2", []byte(strings.Repeat(" ", 4))) //

	case template.PacketKeyRead:

	case template.PacketStatusRequest:

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getAdditionalRooms(station uint64, guest *record.Guest) []string {

	roomList := []string{"00000000", "00000000", "00000000", "00000000", "00000000"}

	data, exist := guest.GetGeneric(defines.AdditionalRooms)
	if !exist {
		return roomList
	}

	additionalRooms, ok := data.(string)
	if !ok || len(additionalRooms) == 0 {
		return roomList
	}

	rooms := strings.Split(additionalRooms, ";")
	mainRoom := d.getRoom(station, guest)
	index := 0
	for i := range rooms {
		if index < 5 {
			room := d.EncodeString(station, rooms[i])
			if len(room) > 0 && len(room) <= 8 && room != mainRoom {
				if len(room) <= 8 {
					room = fmt.Sprintf("%0*s", 8, room)
				}
				roomList[index] = room
				index++
			}
		}
	}
	return roomList
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

func (d *Dispatcher) getRoom(station uint64, guest *record.Guest) string {
	room := d.EncodeString(station, guest.Reservation.RoomNumber)
	if len(room) < 8 {
		room = fmt.Sprintf("%0*s", 8, room)
	}

	return room
}

// GetKeyType maps the pms KeyType and the vendor KeyType
func (d *Dispatcher) GetCardType(guest *record.Guest) string {
	if keyType, exist := guest.GetGeneric(defines.KeyType); exist {
		// N (new), D (duplicate), O (oneshot)

		keyT := cast.ToString(keyType)

		switch keyT {
		case "N":
			return "1"
		case "D":
			return "3"
		case "O":
			return "O"

		default:
			return "1"
		}
	}
	return "1"
}

func (d *Dispatcher) getEncoderID(guest *record.Guest) string {
	encoder := "00"
	if coder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
		encoder = cast.ToString(coder)
		if len(encoder) == 1 {
			encoder = "0" + encoder
		}
		return encoder
	}

	return encoder
}

func (d *Dispatcher) getOperator(guest *record.Guest) string {
	userId, ok := guest.GetGeneric(defines.UserID)
	if !ok || len(cast.ToString(userId)) == 0 {
		return "000000"
	}

	operator := cast.ToString(userId)
	if len(operator) == 0 {
		return "000000"
	}

	if len(operator) > 6 {
		operator = operator[:6]
	}

	return operator
}

func (d *Dispatcher) GetArrivalDate(guest *record.Guest) string {
	arrivalDate := guest.Reservation.ArrivalDate
	formatted := arrivalDate.Format("0601021504")

	return formatted
}

func (d *Dispatcher) GetDapartureDate(guest *record.Guest) string {
	departureDate := guest.Reservation.DepartureDate
	formatted := departureDate.Format("0601021504")

	return formatted
}

func (d *Dispatcher) GetRoom(guest *record.Guest) string {
	room := guest.Reservation.RoomNumber

	if len(room) < 8 {
		room = fmt.Sprintf("%0*s", 8, room)
	}

	return room
}

func (d *Dispatcher) GetCount(guest *record.Guest) string {
	if count, exist := guest.GetGeneric(defines.KeyCount); exist {
		counts := cast.ToString(count)
		counts = fmt.Sprintf("%0*s", 2, counts)
		return counts
	}

	return "01"
}
