package bartech

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/bartech/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	station, _ := d.GetStationAddr(addr)

	guestNameLength := 30

	switch packetName {

	case template.PacketCheckIn:
		if d.simpleCheckin(station) {
			packetName = template.PacketCheckInSimple
		} else if !d.sendHappyHour(station) {
			packetName = template.PacketCheckInWithoutHappyHour
			guestNameLength = 32
		}

	}

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	extensionWidth := d.getExtensionWidth(station)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	case template.PacketEndOfDay:
		opt := pad.Left("0", extensionWidth, "0")
		packet.Add("Opt", []byte(opt))
		return packet, nil

	}

	// extension pre check
	extension := d.getExtension(context)

	if len(extension) == 0 {
		return nil, errors.Errorf("extension not found in context '%T'", context)
	}
	if len(extension) > extensionWidth {
		return nil, errors.Errorf("extension '%s' too long (maximum width: %d)", extension, extensionWidth)
	}

	ext := cast.ToInt16(extension)
	if ext == 0 {
		return nil, errors.Errorf("extension '%s' is not numerical", extension)
	}

	extension = pad.Left(extension, extensionWidth, "0")

	packet.Add("Extension", []byte(extension))

	//

	switch packetName {

	case template.PacketCheckOut, template.PacketLockBar, template.PacketUnlockBar:
		// nothing more to do

	case template.PacketCheckIn, template.PacketCheckInSimple, template.PacketCheckInWithoutHappyHour:

		// Cmd 32 = unlocked, 35 = locked
		if d.UnlockBar(context) {
			packet.Add("Cmd", []byte("32"))
		} else {
			packet.Add("Cmd", []byte("35"))
		}

		if packetName == template.PacketCheckInSimple {
			break
		}

		guest, _ := context.(*record.Guest)

		name := ""
		if d.sendGuestName(station) {
			name = d.getGuestName(guest)
		}
		name = strings.Trim(name, ",")
		name = d.encodeString(addr, name)

		if len(name) > guestNameLength {
			name = name[0:guestNameLength]
			name = strings.Trim(name, ",")
		}
		if len(name) < guestNameLength {
			name = pad.Right(name, guestNameLength, " ")
		}

		packet.Add("GuestName", []byte(name)) // 30 stellig oder 32 stellig wenn Happy Hour disabled ist, left justified

		dt := ""
		if d.sendCheckoutDate(station) {
			dt = d.getCheckoutDate(guest)
		}
		if len(dt) != 8 {
			dt = "        " // fallback damit der parser nicht meckert
		}
		packet.Add("CODate", []byte(dt)) // mmddccyy

		packet.Add("HappyHour", []byte("00"))    // 00 = no, 01 = free, 02 = first comsumption free, 03 = first day free
		packet.Add("Password", []byte("000000")) // 6 stellig, right justified, 0 filled optional
		packet.Add("Opt1", []byte("Y"))          // safe used
		packet.Add("Opt2", []byte("N"))          // extra command
		packet.Add("Opt3", []byte("N"))          // extra command
		packet.Add("Opt4", []byte("N"))          // extra command

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getExtension(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		return strings.TrimLeft(guest.Reservation.RoomNumber, "0")
	}
	return ""
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

func (d *Dispatcher) getCheckoutDate(guest *record.Guest) string {
	dt := time.Now()
	if guest != nil && !guest.Reservation.ArrivalDate.IsZero() {
		dt = guest.Reservation.DepartureDate
	}
	return dt.Format("01022006") // mmddccyy
}

func (d *Dispatcher) encodeString(addr string, data string) string {
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
