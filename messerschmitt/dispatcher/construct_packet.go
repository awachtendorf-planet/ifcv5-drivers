package messerschmitt

import (
	"fmt"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/messerschmitt/template"

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

	case template.PacketAck, template.PacketNak, template.PacketEnq, template.PacketEOT, template.PacketAck0, template.PacketAck1:
		return packet, nil

	}

	station, _ := d.GetStationAddr(addr)

	timeLayout := "150405" // hhmmss
	var dateLayout string

	isSerial := d.IsSerialDevice(addr)

	if isSerial {
		dateLayout = "020106" // ddmmyy
	} else {
		dateLayout = "02012006" // ddmmyyyy
	}

	if packetName == template.PacketSynchronisation {

		var now time.Time

		// get time zone from config

		if timezone := d.GetTimeZone(station); len(timezone) > 0 {

			if loc, err := time.LoadLocation(timezone); err == nil {
				now = time.Now().In(loc)
			} else {
				log.Warn().Msgf("%T addr '%s' get time zone failed, err=%s", d, addr, err)
			}

		}

		if now.IsZero() { // fallback mirror vendor time from incoming packet

			if packet, ok := context.(*ifc.LogicalPacket); ok {

				vendorDate := packet.Data()["Date"]
				vendorTime := packet.Data()["Time"]

				if t, err := time.Parse(dateLayout+timeLayout, fmt.Sprintf("%s%s", vendorDate, vendorTime)); err == nil {
					now = t
				}

			}

			if now.IsZero() { // fallback, use current time
				now = time.Now()
			}

		}

		date := now.Format(dateLayout)
		time := now.Format(timeLayout)

		packet.Add("Date", []byte(date))
		packet.Add("Time", []byte(time))

		return packet, nil
	}

	guest, ok := context.(*record.Guest)
	if !ok {
		return packet, errors.Errorf("context '%T' not supported", context)
	}

	room := d.formatRoom(station, guest, isSerial)
	packet.Add("Room", []byte(room))

	encoder := d.formatEncoder(station, guest, isSerial)
	packet.Add("Encoder", []byte(encoder))

	switch packetName {

	case template.PacketKeyRequest:

		keyType, _ := d.getKeyType(guest)
		packet.Add("Assignment", []byte(keyType))

		keyCount := d.getKeyCount(guest)
		packet.Add("Cards", []byte(cast.ToString(keyCount)))

		ap := d.formatAccesspoints(station, guest, isSerial)
		packet.Add("AccessPoints", []byte(ap))

		resno := pad.Left(guest.Reservation.ReservationID, 16, "0")
		departure := guest.Reservation.DepartureDate

		if isSerial {

			departureDate := departure.Format("020106")
			packet.Add("DepartureDate", []byte(departureDate)) // ddmmyy

			departureTime := departure.Format("1504")              // hhmm
			packet.Add("DepartureTime", []byte(departureTime[:3])) // hhm, array should be safe, because of it's from a time struct

			packet.Add("GuestIndex2", []byte(resno[8:]))
			packet.Add("GuestIndex1", []byte(resno[:8]))

			// default
			packet.Add("Suite", []byte("000000"))
			packet.Add("Data", []byte(pad.Left("", 19, "0")))

		} else {

			departureDate := departure.Format("02012006")
			packet.Add("DepartureDate", []byte(departureDate)) // ddmmyyyy

			// to next full hour
			departure = departure.Add(1 * time.Hour)
			departureTime := departure.Format("15")            // hh
			packet.Add("DepartureTime", []byte(departureTime)) // hh

			packet.Add("GuestIndex", []byte(resno))

		}

	case template.PacketKeyDelete:
		// nothing more to do

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getEncoder(guest *record.Guest) string {

	if guest == nil {
		return "0"
	}

	if encoder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
		return cast.ToString(encoder)
	}

	return "0"
}

func (d *Dispatcher) getKeyCount(guest *record.Guest) int {

	if guest == nil {
		return 0
	}

	if count, exist := guest.GetGeneric(defines.KeyCount); exist {
		if keyCount := cast.ToInt(count); keyCount > 0 {
			return keyCount
		}
	}

	return 1
}

func (d *Dispatcher) getKeyType(guest *record.Guest) (string, error) {

	if guest == nil {
		return "", errors.Errorf("guest is nil")
	}

	if keyType, exist := guest.GetGeneric(defines.KeyType); exist {

		switch cast.ToString(keyType) {

		case "N": // new
			return "0", nil // new assignment

		case "D": // duplicate
			return "1", nil // old assignment, duplicate

		default:
			return "", errors.Errorf("unsupported key type '%s'", keyType)

		}
	}

	return "", errors.Errorf("key type cannot be determined")
}

func (d *Dispatcher) formatAccesspoints(_ uint64, guest *record.Guest, isSerial bool) string {

	var ap []byte

	if data, exist := guest.GetGeneric(defines.KeyOptions); exist {
		ap = []byte(cast.ToString(data))
	}

	for i := range ap {
		if ap[i] != '0' && ap[i] != '1' {
			ap[i] = '0'
		}
	}

	var accessPoints string

	if isSerial {
		accessPoints = pad.Right(string(ap), 27, "0")
		if len(accessPoints) > 27 {
			accessPoints = accessPoints[:27]
		}
	} else {
		accessPoints = pad.Right(string(ap), 40, "0")
		if len(accessPoints) > 40 {
			accessPoints = accessPoints[:40]
		}
	}

	return accessPoints
}

func (d *Dispatcher) formatRoom(station uint64, guest *record.Guest, isSerial bool) string {

	room := d.GetRoom(station, guest, isSerial)

	if isSerial {
		room = pad.Left(room, 4, "0")
	} else {
		room = pad.Left(room, 10, " ")
	}

	return room
}

func (d *Dispatcher) formatEncoder(_ uint64, guest *record.Guest, isSerial bool) string {

	encoder := d.getEncoder(guest)

	if isSerial {

		number := cast.ToInt(encoder)
		if number < 10 {
			return encoder
		}
		number = 65 + number - 10

		return fmt.Sprintf("%c", number)
	}

	encoder = pad.Left(encoder, 2, "0")

	return encoder
}

// FormatEncoder ...
func (d *Dispatcher) FormatEncoder(guest *record.Guest, isSerial bool) string {
	return d.formatEncoder(0, guest, isSerial)
}

// GetEncoder ...
func (d *Dispatcher) GetEncoder(guest *record.Guest) int {
	encoder := cast.ToInt(d.getEncoder(guest))
	return encoder
}

// GetRoom ...
func (d *Dispatcher) GetRoom(station uint64, guest *record.Guest, isSerial bool) string {

	if guest == nil {
		return ""
	}

	if isSerial {
		return guest.Reservation.RoomNumber
	}

	room := d.EncodeString(station, guest.Reservation.RoomNumber)

	return room
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

func (d *Dispatcher) EncodeString(station uint64, data string) string {

	encoding := d.GetEncodingByStation(station)

	if len(encoding) == 0 {
		encoding = DEFAULT_ENCODING
	}

	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
		return string(enc)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}

	return data
}
