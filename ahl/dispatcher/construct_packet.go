package ahl

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

	"github.com/weareplanet/ifcv5-drivers/ahl/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	// station, _ := d.GetStationAddr(addr)
	// protocolType := d.GetProtocolType(station)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketLinkStart, template.PacketLinkAlive:
		return packet, nil

	case template.PacketWakeupClear:
		if generic, ok := context.(*record.Generic); ok {
			//generic.Delete(defines.Timestamp)
			generic.Set(defines.Timestamp, "00000") // suppression of all wake-up calls
		}

	}

	station, _ := d.GetStationAddr(addr)
	protocolType := d.GetProtocolType(station)

	// check extension width
	extensionWidth := d.GetExtensionWidth(protocolType)

	extension := d.GetExtension(context)
	if len(extension) == 0 {
		return nil, errors.Errorf("extension not found in context '%T'", context)
	}
	if len(extension) > extensionWidth {
		return nil, errors.Errorf("extension '%s' too long (maximum width: %d)", extension, extensionWidth)
	}

	// construct packet
	packet.Add("LRC", []byte("??"))
	packet.Add("Extension", d.rightJustified(extension, extensionWidth))

	switch packetName {

	case template.PacketCheckOut:
		// nothing more to do, only extension is needed

	case template.PacketCheckIn, template.PacketDataChange, template.PacketWakeupSet, template.PacketWakeupClear:

		if protocolType == AHL4400_5 || protocolType == AHL4400_8 {

			guest, _ := context.(*record.Guest)

			// encode guestname
			guestName := d.encode(addr, d.getGuestName(guest))
			packet.Add("GuestName", d.leftJustified(guestName, 20))
			packet.Add("Occupation", d.rightJustified(d.getOccupation(guest), 1))
			packet.Add("GuestLanguage", d.rightJustified(d.getGuestLanguage(station, guest), 1))
			packet.Add("VIPState", d.rightJustified(d.getVIPState(guest), 1))
			packet.Add("MessageWaiting", d.rightJustified(d.getMessageState(guest), 1))
			packet.Add("DOD", d.rightJustified(d.getDOD(station, guest), 2))
			packet.Add("DND", d.rightJustified(d.getDND(guest), 1))
			packet.Add("GroupName", d.rightJustified(d.getGroupName(guest), 3))
			packet.Add("Password", d.rightJustified("", 4))
			packet.Add("DepositAmount", d.rightJustified("", 9))
			packet.Add("WakeupTime", d.rightJustified(d.getWakeupTime(context), 5))

		} else {
			return nil, errors.Errorf("no handler for template '%s' protocol type: %d defined", packetName, protocolType)
		}

	case template.PacketRoomStatus:

		if protocolType == AHL4400_5 || protocolType == AHL4400_8 {

			guest, _ := context.(*record.Guest)

			packet.Add("Occupation", d.rightJustified(d.getOccupation(guest), 1))
			packet.Add("VIPState", d.rightJustified(d.getVIPState(guest), 1))
			packet.Add("MessageWaiting", d.rightJustified(d.getMessageState(guest), 1))
			packet.Add("DOD", d.rightJustified(d.getDOD(station, guest), 2))
			packet.Add("DND", d.rightJustified(d.getDND(guest), 1))
			// do not change values in pbx
			packet.Add("GuestName", d.rightJustified(strings.Repeat("0", 20), 20))
			packet.Add("GuestLanguage", d.rightJustified("0", 1))
			packet.Add("GroupName", d.rightJustified("", 3))
			packet.Add("WakeupTime", d.rightJustified("", 5))
			packet.Add("Password", d.rightJustified("", 4))
			packet.Add("DepositAmount", d.rightJustified("", 9))

		} else {
			return nil, errors.Errorf("no handler for template '%s' protocol type: %d defined", packetName, protocolType)
		}

	default:
		return nil, errors.Errorf("no handler for template '%s' defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) GetExtension(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Reservation.RoomNumber
	}
	if generic, ok := context.(*record.Generic); ok {
		if extension, exist := generic.Get(defines.WakeExtension); exist {
			return cast.ToString(extension)
		}
	}
	return ""
}

func (d *Dispatcher) rightJustified(data string, width int) []byte {
	if len(data) == width {
		return []byte(data)
	}
	if len(data) > width {
		return []byte(data[0:width])
	}
	str := pad.Left(data, width, " ")
	return []byte(str)
}

func (d *Dispatcher) leftJustified(data string, width int) []byte {
	if len(data) == width {
		return []byte(data)
	}
	if len(data) > width {
		return []byte(data[0:width])
	}
	str := pad.Right(data, width, " ")
	return []byte(str)
}

func (d *Dispatcher) encode(addr string, data string) string {
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

func (d *Dispatcher) getOccupation(guest *record.Guest) string {
	if guest == nil {
		return "0"
	}
	if guest.Reservation.SharedInd {
		return "*"
	}
	return ""
}

func (d *Dispatcher) getGuestName(guest *record.Guest) string {
	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {
		return strings.Repeat("0", 20)
	}
	if len(guest.DisplayName) > 0 {
		return guest.DisplayName
	}
	return fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
}

func (d *Dispatcher) getGuestLanguage(station uint64, guest *record.Guest) string {
	if guest == nil || len(guest.Language) == 0 {
		return "0"
	}
	language := strings.ToUpper(guest.Language)
	if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil && len(languageCode) == 1 {
		return languageCode
	}
	return ""
}

func (d *Dispatcher) getVIPState(guest *record.Guest) string {
	if guest == nil {
		return "0"
	}
	if len(guest.VIPStatus) > 0 {
		return "V"
	}
	return ""
}

func (d *Dispatcher) getGroupName(guest *record.Guest) string {
	return "" // unused in 4400
}

func (d *Dispatcher) getDOD(station uint64, guest *record.Guest) string {
	if guest == nil {
		return "  "
	}
	cos := cast.ToString(guest.Rights.ClassOfService)
	if dod, err := d.GetMapping("classofservice", station, cos, false); err == nil {
		return dod
	}
	return cos
}

func (d *Dispatcher) getDND(guest *record.Guest) string {
	if guest == nil {
		return " "
	}
	if guest.Reservation.DoNotDisturb {
		return "2"
	}
	return "0"
}

func (d *Dispatcher) getMessageState(guest *record.Guest) string {
	if guest == nil {
		return " "
	}
	if guest.Reservation.MessageLightStatus {
		return "1"
	}
	return ""
}

func (d *Dispatcher) getWakeupTime(context interface{}) string {
	if generic, ok := context.(*record.Generic); ok {
		if wakeupTime, exist := generic.Get(defines.Timestamp); exist {
			if t, ok := wakeupTime.(time.Time); ok && !t.IsZero() {
				return fmt.Sprintf("%s ", t.Format("1504")) // hhmm(blank) -> military time
			}
			if s, ok := wakeupTime.(string); ok && len(s) == 5 { // preformated value
				return s
			}
		}
	}
	return ""
}
