package mitel

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/mitel/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	case template.PacketAlive, template.PacketSwapStart, template.PacketSwapEnd:
		return packet, nil

	}

	// extension pre check
	extension := d.getExtension(context)
	extensionWidth := d.getExtensionWidth(station)

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

	extension = pad.Left(extension, extensionWidth, " ")

	packet.Add("Extension", []byte(extension))

	//

	switch packetName {

	case template.PacketCheckIn, template.PacketCheckOut:

	case template.PacketChangeName:

		status := 2 // replace
		guest, ok := context.(*record.Guest)
		if ok && guest.Reservation.SharedInd {
			status = 1 // add
		}

		name := d.getGuestName(guest)
		name = strings.Trim(name, ",")
		name = d.encodeString(addr, name)

		packet.Add("Status", d.leftJustifiedNumeric(status, 1)) // 1 = add, 2 = replace, 3 = delete
		packet.Add("Name", d.leftJustifiedString(name, 21))

	case template.PacketMessageLamp:

		status := 0
		if guest, ok := context.(*record.Guest); ok && guest.Reservation.MessageLightStatus {
			status = 1
		}

		packet.Add("Status", d.leftJustifiedNumeric(status, 1)) // 1 = on, 0 = off

	case template.PacketSetRestriction:

		status := 0
		if guest, ok := context.(*record.Guest); ok {
			status = d.getClassOfService(station, guest)
		}

		packet.Add("Status", d.leftJustifiedNumeric(status, 1)) // class of service

	case template.PacketWakeupSet:

		status := d.getWakeupTime(context)
		packet.Add("Status", d.leftJustifiedString(status, 4)) // hhmm

	case template.PacketWakeupClear:

		packet.Add("Status", d.leftJustifiedString("", 4)) // clear time

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getExtension(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		return strings.TrimLeft(guest.Reservation.RoomNumber, "0")
	}
	if generic, ok := context.(*record.Generic); ok {
		if extension, exist := generic.Get(defines.WakeExtension); exist {
			extension := cast.ToString(extension)
			return strings.TrimLeft(extension, "0")
		}
	}
	return ""
}

func (d *Dispatcher) getGuestName(guest *record.Guest) string {
	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {
		return strings.Repeat("0", 21)
	}
	if len(guest.DisplayName) > 0 {
		return guest.DisplayName
	}
	return fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
}

func (d *Dispatcher) getClassOfService(station uint64, guest *record.Guest) int {
	if guest == nil {
		return 0
	}
	cos := cast.ToString(guest.Rights.ClassOfService)
	if data, err := d.GetMapping("classofservice", station, cos, false); err == nil {
		return cast.ToInt(data)
	}
	return guest.Rights.ClassOfService
}

func (d *Dispatcher) getWakeupTime(context interface{}) string {
	if generic, ok := context.(*record.Generic); ok {
		if wakeupTime, exist := generic.Get(defines.Timestamp); exist {
			if t, ok := wakeupTime.(time.Time); ok && !t.IsZero() {
				return t.Format("1504") // hhmm
			}
			if s, ok := wakeupTime.(string); ok && len(s) == 4 { // preformated value
				return s
			}
		}
	}
	return ""
}

func (d *Dispatcher) leftJustifiedString(data string, width int) []byte {
	if len(data) == width {
		return []byte(data)
	}
	if len(data) > width {
		return []byte(data[0:width])
	}
	str := pad.Right(data, width, " ")
	return []byte(str)
}

func (d *Dispatcher) leftJustifiedNumeric(value int, width int) []byte {
	data := cast.ToString(value)
	return d.leftJustifiedString(data, width)
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
