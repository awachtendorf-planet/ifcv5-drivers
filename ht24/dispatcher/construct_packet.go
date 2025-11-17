package ht24

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/weareplanet/ifcv5-drivers/ht24/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	station, _ := d.GetStationAddr(addr)

	protocol := d.GetProtocol(station)

	if packetName == template.PacketCodeCard {
		switch protocol {
		case 2:
			packetName = template.PacketCodeCard14
		case 3:
			packetName = template.PacketCodeCard15
		default:
			// OK incl. 1,13
		}
	}

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil
	}

	guest, ok := context.(*record.Guest)
	if !ok {
		return packet, errors.Errorf("context '%T' not supported", context)
	}

	encoder := d.getEncoder(guest)

	roomNumber := d.getRoom(station, guest)
	if len(roomNumber) > RoomNumberLength {
		return packet, errors.Errorf("room number '%s' exceeds the limit of %d characters", roomNumber, RoomNumberLength)
	} else if len(roomNumber) == 0 {
		return packet, errors.New("room number is empty")
	}

	packet.Add("EncoderNumber", []byte(cast.ToString(encoder)))
	packet.Add("Room", []byte(roomNumber))

	switch packetName {

	case template.PacketCheckout:
		//OK
	case template.PacketCodeCard, template.PacketCodeCard14, template.PacketCodeCard15:

		keyType := d.GetKeyType(guest)

		cardCount := d.getKeyCount(guest)

		operator := d.getOperator(guest)

		initialDate := d.getFormattedDate(guest.Reservation.ArrivalDate)

		expireDate := ""
		if len(initialDate) > 0 {
			expireDate = d.getFormattedDate(guest.Reservation.DepartureDate)
		}

		rooms := d.getAdditionalRooms(station, guest)
		if len(rooms) < 3 {

			return packet, errors.Errorf("additional rooms array length is %d, need atleast a length of 3", len(rooms))
		}

		allowedAuthorisations, deniedAuthorisations := d.getAccessPoints(station, guest)

		track2Data := ""
		if d.SendTrack2(station) {

			track2Data = d.getTrack2Data(guest)
		}

		keyId := ""
		if keyTag := d.GetKeyID(station); keyTag != "" {

			keyId = d.getKeyID(keyTag, guest)
		}

		packet.Add("CardType", []byte(cast.ToString(keyType)))
		packet.Add("CardCount", []byte(cast.ToString(cardCount)))
		packet.Add("Room2", []byte(rooms[0]))
		packet.Add("Room3", []byte(rooms[1]))
		packet.Add("Room4", []byte(rooms[2]))
		packet.Add("AA", []byte(allowedAuthorisations))
		packet.Add("AD", []byte(deniedAuthorisations))
		packet.Add("InitalDate", []byte(initialDate))
		packet.Add("ExpireDate", []byte(expireDate))
		packet.Add("Operator", []byte(operator))
		packet.Add("Track1", []byte(""))
		packet.Add("Track2", []byte(track2Data))
		packet.Add("KeyID", []byte(keyId)) // Just relevant for PacketCodeCard14 + *15

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
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

func (d *Dispatcher) getEncoder(guest *record.Guest) int {
	if encoder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
		return cast.ToInt(encoder)
	}
	return 0
}

func (d *Dispatcher) getOperator(guest *record.Guest) string {
	userId, ok := guest.GetGeneric(defines.UserID)
	if !ok {
		return ""
	}
	operator := cast.ToString(userId)
	if len(operator) > 20 {
		operator = operator[:20]
	}
	return operator
}

func (d *Dispatcher) getFormattedDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("15020106")
}

func (d *Dispatcher) getRoom(station uint64, guest *record.Guest) string {
	room := d.EncodeString(station, guest.Reservation.RoomNumber)
	return room
}

func (d *Dispatcher) getAdditionalRooms(station uint64, guest *record.Guest) []string {

	roomList := []string{"", "", ""}

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
		if index < 3 {
			room := d.EncodeString(station, rooms[i])
			if len(room) > 0 && len(room) <= RoomNumberLength && room != mainRoom {
				roomList[index] = room
				index++
			}
		}
	}
	return roomList
}

func (d *Dispatcher) getKeyCount(guest *record.Guest) int {
	if count, exist := guest.GetGeneric(defines.KeyCount); exist {
		countInt := cast.ToInt(count)
		if countInt > 0 {
			if countInt < 9 {
				return countInt
			}
			return 9
		}
	}
	return 1
}

func (d *Dispatcher) getTrack2Data(guest *record.Guest) string {
	if track, exist := guest.GetGeneric(defines.Track2); exist {
		track2 := cast.ToString(track)
		if len(track2) > 17 {
			return "" //ERROR
		}
		return track2
	}
	return ""
}

// getAccessPoints returns allowedAccessPoints, deniedAccessPoints
func (d *Dispatcher) getAccessPoints(station uint64, guest *record.Guest) (string, string) {
	keyOptions := []byte{}

	if keyOp, exist := guest.GetGeneric(defines.KeyOptions); exist {
		keyOpts := cast.ToString(keyOp)
		keyOptions = []byte(keyOpts)
	}

	keyOptionsOverlay := d.GetAccesPointOverlay(station)

	if len(keyOptions) > 0 {

		for index := range keyOptionsOverlay {

			if index < len(keyOptions) {

				letter := keyOptionsOverlay[index]

				if letter == '1' || letter == '0' {

					keyOptions[index] = letter

				}

			} else {
				break
			}
		}
	} else {
		keyOptions = []byte(keyOptionsOverlay)
	}

	allowedAccessPoints := ""
	deniedAccessPoints := ""

	for index, letter := range keyOptions {

		if letter == '1' && index < 8 {

			allowedAccessPoints += cast.ToString(index + 1)

		} else if letter == '0' && index < 8 {

			deniedAccessPoints += cast.ToString(index + 1)

		}
	}

	return allowedAccessPoints, deniedAccessPoints
}

// GetKeyType maps the pms KeyType and the vendor KeyType
func (d *Dispatcher) GetKeyType(guest *record.Guest) string {

	if keyType, exist := guest.GetGeneric(defines.KeyType); exist {

		switch cast.ToString(keyType) {

		case "N": // new
			return "N"

		case "D": // duplicate
			return "C"

		case "O": // oneshot
			return "A"
		}
	}

	return "N"
}

func (d *Dispatcher) getKeyID(udfx string, guest *record.Guest) string {
	if keyType, exist := guest.GetGeneric(udfx); exist {
		return cast.ToString(keyType)
	}

	return ""
}
