package inhova

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

	"github.com/weareplanet/ifcv5-drivers/inhova/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	guest, ok := context.(*record.Guest)
	if !ok {
		return packet, errors.Errorf("context '%T' not supported", context)
	}

	station, _ := d.GetStationAddr(addr)

	encoder := d.getEncoder(guest)
	packet.Add("Encoder", []byte(cast.ToString(encoder)))

	room := d.getRoom(station, guest)
	packet.Add("Room", []byte(room))

	operator := d.getOperator(guest)
	packet.Add("Operator", []byte(operator))

	switch packetName {

	case template.PacketKeyRequest:

		setEmpty := func(fields ...string) {
			for i := range fields {
				packet.Add(fields[i], []byte(""))
			}
		}

		setEmpty(
			"KeyCount",
			"ActivationDate", "ActivationTime",
			"ExpirationDate", "ExpirationTime",
			"Grants", "KeyPad", "CardOperation",
			"TesaHotelEncoder", // ?
			"Track1", "Track2",
			"Technology", // ?
			"Room2", "Room3", "Room4",
			"CardID", "CardType",
			"PhoneNumber", "Mail", "Mail2", "Mail3", "Mail4",
		)

		action, _ := d.getKeyType(guest)
		packet.Add("Action", []byte(action)) // I/G

		keyCount := d.getKeyCount(guest)
		if keyCount > 1 {
			packet.Add("KeyCount", []byte(cast.ToString(keyCount)))
		}

		if d.sendActivation(station) {
			activationDate := d.getDate(guest.Reservation.ArrivalDate)
			activationTime := d.getTime(guest.Reservation.ArrivalDate)
			if len(activationDate) == 10 && len(activationTime) == 5 {
				packet.Add("ActivationDate", []byte(activationDate))
				packet.Add("ActivationTime", []byte(activationTime))
			}
		}

		expiratenDate := d.getDate(guest.Reservation.DepartureDate)
		expiratenTime := d.getTime(guest.Reservation.DepartureDate)
		if len(expiratenDate) == 10 && len(expiratenTime) == 5 {
			packet.Add("ExpirationDate", []byte(expiratenDate))
			packet.Add("ExpirationTime", []byte(expiratenTime))
		}

		if d.sendTrack2(station) {
			if trackData := d.getTrack2Data(guest); len(trackData) > 0 {
				packet.Add("Track2", []byte(trackData))
			}
		}

		rooms := d.getAdditionalRooms(station, guest)
		for i := range rooms {
			if i > 2 || len(rooms[i]) == 0 {
				break
			}
			packet.Add(fmt.Sprintf("Room%d", i+2), []byte(rooms[i]))
		}

		grants := d.getGrants(station, guest)
		if len(grants) > 0 {
			packet.Add("Grants", []byte(grants))
		}

		packet.Add("ReturnCardID", []byte("1"))

	case template.PacketKeyDelete:
		break

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getEncoder(guest *record.Guest) int {
	if encoder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
		return cast.ToInt(encoder)
	}
	return 0
}

func (d *Dispatcher) getRoom(station uint64, guest *record.Guest) string {
	room := d.EncodeString(station, guest.Reservation.RoomNumber)
	return room
}

func (d *Dispatcher) getOperator(guest *record.Guest) string {
	userId, ok := guest.GetGeneric(defines.UserID)
	if !ok {
		return ""
	}
	operator := cast.ToString(userId)
	if len(operator) > 10 {
		operator = operator[:10]
	}
	return operator
}

func (d *Dispatcher) getKeyCount(guest *record.Guest) int {
	if count, exist := guest.GetGeneric(defines.KeyCount); exist {
		if keyCount := cast.ToInt(count); keyCount > 0 {
			return keyCount
		}
	}
	return 1
}

func (d *Dispatcher) getTrack2Data(guest *record.Guest) string {
	if track, exist := guest.GetGeneric(defines.Track2); exist {
		track2 := cast.ToString(track)
		return track2
	}
	return ""
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
		if index > 2 {
			break
		}
		room := d.EncodeString(station, rooms[i])
		if len(room) > 0 && room != mainRoom {
			roomList[index] = room
			index++
		}
	}

	return roomList
}

func (d *Dispatcher) getGrants(station uint64, guest *record.Guest) string {

	var grants []string

	keyOptions := d.getKeyOptions(station)

	if len(keyOptions) == 0 {
		if value, exist := guest.GetGeneric(defines.KeyOptions); exist {
			keyOptions = cast.ToString(value)
		}
	}

	if len(keyOptions) == 0 {
		return ""
	}

	length := len(keyOptions)

	for i := 0; i < length; i++ {

		if set := keyOptions[length-i-1] == '1'; set {
			index := cast.ToString(i + 1)
			if grant, err := d.GetMapping("accesspoint", station, index, true); err == nil && len(grant) > 0 {
				grants = append(grants, grant)
			}
		}

	}

	keyOptions = strings.Join(grants, ",")

	return keyOptions
}

func (d *Dispatcher) getDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("02/01/2006")
}

func (d *Dispatcher) getTime(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("15:04")
}

func (d *Dispatcher) getKeyType(guest *record.Guest) (string, error) {

	if guest == nil {
		return "", errors.Errorf("guest is nil")
	}

	if keyType, exist := guest.GetGeneric(defines.KeyType); exist {

		switch cast.ToString(keyType) {

		case "N": // new
			return "I", nil // guest checkin

		case "D": // duplicate
			return "G", nil // copy guest card

		default:
			return "", errors.Errorf("unsupported key type '%s'", keyType)

		}
	}

	return "", errors.Errorf("key type cannot be determined")
}

// GetKeyType maps the pms KeyType and the vendor KeyType
func (d *Dispatcher) GetKeyType(job *order.Job) (string, error) {

	if job == nil {
		return "", errors.Errorf("context is nil")
	}

	if job.Action == order.KeyDelete {
		return "O", nil
	}

	if guest, ok := job.Context.(*record.Guest); ok {
		return d.getKeyType(guest)
	}

	return "", errors.Errorf("key type cannot be determined")
}

// GetKeyCount ...
func (d *Dispatcher) GetKeyCount(job *order.Job) int {

	if job == nil {
		return 1
	}

	if guest, ok := job.Context.(*record.Guest); ok {
		return d.getKeyCount(guest)
	}

	return 1
}

// EncodeString ...
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
