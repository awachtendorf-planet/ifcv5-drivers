package definity

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/definity/template"

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

	station, _ := d.GetStationAddr(addr)

	// message counter
	// ascii, transparent mode: 1 nipple counter (modulo 10) + 2, 4-bits), 2 nipple 0xf, range 0x2 to 0xb, restart 0x2

	counter := d.getMessageCounter(station) // 0-9

	counter = (counter % 10) + 2
	counter = (counter << 4) + 0xf

	packet.Add("MSGCT", []byte{counter})

	switch packetName {

	case template.PacketHeartbeat, template.PacketLinkEndConfirmed, template.PacketSyncStart, template.PacketSyncEnd:
		return packet, nil

	case template.PacketHouseKeeperRoomAccepted, template.PacketHouseKeeperStationAccepted, template.PacketHouseKeeperRoomRejected, template.PacketHouseKeeperStationRejected:

		roomStatus, ok := context.(record.RoomStatus)
		if !ok {
			return packet, errors.Errorf("construct packet '%s' failed, unexpected context type '%T'", packetName, context)
		}

		roomLength := d.getRoomLength(station)
		room := pad.Left(roomStatus.RoomNumber, roomLength, "0")

		packet.Add("RSN", []byte(room))

		dig := pad.Left(roomStatus.UserID, 6, "0")
		packet.Add("DIG", []byte(dig))

		return packet, nil
	}

	job, ok := context.(*order.Job)
	if !ok {
		return packet, errors.Errorf("construct packet '%s' failed, unexpected context type '%T'", packetName, context)
	}

	guest, _ := job.Context.(*record.Guest)

	// set room number
	room := d.GetRoom(job.Context)
	roomLength := d.getRoomLength(station)
	room = pad.Left(room, roomLength, "0")

	packet.Add("RSN", []byte(room))

	// set coverage path
	if job.Action != order.Checkout &&
		packetName != template.PacketSwitchMessageLamp &&
		packetName != template.PacketSetRestriction {

		cp := d.getCoveragePath(station)
		cp = pad.Left(cp, 4, "0")
		packet.Add("COVERAGE_PATH", []byte(cp))

	}

	switch packetName {

	case template.PacketCheckIn:

		packet.Add("REQ_DID", []byte(" ")) // y to request DID
		fallthrough

	case template.PacketDataChange:

		packet.Add("VM_PASSWORD", []byte(DEFAULT_PASSWORD))

		language := d.getGuestLanguage(station, guest)
		packet.Add("VM_LANGUAGE", []byte(language))

		var name string
		nameLength := d.getNameLength(station)
		name = d.getName(station, guest)
		if len(name) > nameLength {
			name = name[:nameLength]
		}
		name = pad.Right(name, nameLength, " ")
		packet.Add("DISPLAY_NAME", []byte(name))

	case template.PacketCheckOut:
		// nothing more to do

	case template.PacketRoomDataImageSwap:

		// set defaults for vacant
		packet.Add("OCCUPIED", []byte("0"))
		packet.Add("MESSAGE_WAITING", []byte("0"))
		packet.Add("RESTRICT_LEVEL", []byte(DEFAULT_RESTRICT_LEVEL))
		packet.Add("VM_PASSWORD", []byte(DEFAULT_PASSWORD))
		packet.Add("VM_LANGUAGE", []byte(DEFAULT_LANGUAGE))

		var name string
		nameLength := d.getNameLength(station)

		if job.Action != order.Checkout {

			packet.Add("OCCUPIED", []byte("1"))

			if d.getMessageState(guest) {
				packet.Add("MESSAGE_WAITING", []byte("1"))
			}

			if d.getClassOfService(guest) > 0 {

				if d.getDND(guest) {
					packet.Add("RESTRICT_LEVEL", []byte("05")) // denies all calls to the room
				} else {
					packet.Add("RESTRICT_LEVEL", []byte("00")) // no restriction
				}

			}

			language := d.getGuestLanguage(station, guest)
			packet.Add("VM_LANGUAGE", []byte(language))

			name = d.getName(station, guest)
			if len(name) > nameLength {
				name = name[:nameLength]
			}
		}

		name = pad.Right(name, nameLength, " ")
		packet.Add("DISPLAY_NAME", []byte(name))

	case template.PacketSwitchMessageLamp:

		if d.getMessageState(guest) {
			packet.Add("PROC", []byte{0x1}) // turn on
		} else {
			packet.Add("PROC", []byte{0x2}) // turn off
		}

	case template.PacketSetRestriction:

		if d.getClassOfService(guest) > 0 {

			if d.getDND(guest) {
				packet.Add("RESTRICT_LEVEL", []byte("05")) // denies all calls to the room
			} else {
				packet.Add("RESTRICT_LEVEL", []byte("00")) // no restriction
			}

		} else {
			packet.Add("RESTRICT_LEVEL", []byte(DEFAULT_RESTRICT_LEVEL))
		}

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) GetRoom(context interface{}) string {

	if guest, ok := context.(*record.Guest); ok {
		return strings.TrimLeft(guest.Reservation.RoomNumber, "0")
	}

	if generic, ok := context.(*record.Generic); ok {
		if extension, exist := generic.Get(defines.WakeExtension); exist {
			return cast.ToString(extension)
		}
	}

	return ""
}

func (d *Dispatcher) getName(station uint64, guest *record.Guest) string {

	if guest == nil {
		return ""
	}

	var name string

	name = d.getGuestName(guest)
	name = d.encodeString(station, name)
	name = d.forceASCII(name)

	return name
}

func (d *Dispatcher) getGuestName(guest *record.Guest) string {

	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {
		return ""
	}

	if len(guest.DisplayName) > 0 {
		return guest.DisplayName
	}

	if len(guest.FirstName) == 0 {
		return guest.LastName
	}

	return fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
}

func (d *Dispatcher) getGuestLanguage(station uint64, guest *record.Guest) string {

	if guest == nil || len(guest.Language) == 0 {
		return DEFAULT_LANGUAGE
	}

	language := strings.ToUpper(guest.Language)
	if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil && len(languageCode) == 2 {
		return languageCode
	}

	return DEFAULT_LANGUAGE
}

func (d *Dispatcher) getClassOfService(guest *record.Guest) int {

	if guest == nil {
		return 0 // off
	}

	return guest.Rights.ClassOfService
}

func (d *Dispatcher) getDND(guest *record.Guest) bool {

	if guest == nil {
		return false
	}

	return guest.Reservation.DoNotDisturb
}

func (d *Dispatcher) getMessageState(guest *record.Guest) bool {

	if guest == nil {
		return false
	}

	return guest.Reservation.MessageLightStatus
}

func (d *Dispatcher) forceASCII(s string) string {

	rs := make([]rune, 0, len(s))

	for _, r := range s {

		if r <= 126 {
			rs = append(rs, r)
		} else {
			rs = append(rs, '?')
		}

	}

	return string(rs)
}
