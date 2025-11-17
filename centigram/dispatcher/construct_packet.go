package centigram

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

	"github.com/weareplanet/ifcv5-drivers/centigram/template"

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

	}

	// extension pre check
	extension := d.getExtension(context)
	extensionWidth := 6

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

	extension = pad.Right(extension, extensionWidth, " ")

	packet.Add("Extension", []byte(extension))

	//

	switch packetName {

	case template.PacketCheckIn, template.PacketCheckOut, template.PacketMessageLampOn:

	case template.PacketRoomChange:

		guest, _ := context.(*record.Guest)
		source := d.getOldExtension(guest)

		if len(source) == 0 {
			return nil, errors.Errorf("old extension not found in context '%T'", context)
		}
		if len(source) > extensionWidth {
			return nil, errors.Errorf("old extension '%s' too long (maximum width: %d)", source, extensionWidth)
		}

		ext := cast.ToInt16(source)
		if ext == 0 {
			return nil, errors.Errorf("old extension '%s' is not numerical", source)
		}

		source = pad.Right(source, extensionWidth, " ")
		packet.Add("Source", []byte(source))

	case template.PacketChangeName:

		guest, _ := context.(*record.Guest)

		name := d.getGuestName(guest)
		name = strings.Trim(name, ",")
		name = d.encodeString(addr, name)

		packet.Add("Name", d.formatString(name, 6))

	case template.PacketChangeFCOS:

		code := 0

		if guest, ok := context.(*record.Guest); ok {
			code = d.getGuestLanguage(station, guest)
			if code > 99 {
				code = 99
			}
		}

		packet.Add("Code", d.formatNumeric(code, 2))

	case template.PacketChangeVoiceMail:

		code := 0

		if guest, ok := context.(*record.Guest); ok {
			voiceMail := strings.Trim(guest.Reservation.Voicemail, " ")
			voiceMail = strings.TrimLeft(voiceMail, "0")
			if voiceMail == "Y" || voiceMail == "y" {
				code = 1
			} else {
				code = cast.ToInt(guest.Reservation.Voicemail)
			}
			if code > 99 {
				code = 99
			}
		}

		packet.Add("Code", d.formatNumeric(code, 2))

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

func (d *Dispatcher) getOldExtension(guest *record.Guest) string {
	oldroom := ""
	if guest != nil {
		data, _ := guest.GetGeneric(defines.OldRoomName)
		oldroom = cast.ToString(data)
		oldroom = strings.TrimLeft(oldroom, "0")
	}
	return oldroom
}

func (d *Dispatcher) getGuestLanguage(station uint64, guest *record.Guest) int {
	if guest == nil || len(guest.Language) == 0 {
		return 0
	}
	language := strings.ToUpper(guest.Language)
	if languageCode, err := d.GetMapping("languagecode", station, language, false); err == nil && len(languageCode) <= 2 {
		code := cast.ToInt(languageCode)
		return code
	}
	return 0
}

func (d *Dispatcher) formatString(data string, width int) []byte {
	return d.leftJustifiedString(data, width, " ")
}

func (d *Dispatcher) formatNumeric(value int, width int) []byte {
	data := cast.ToString(value)
	return d.rightJustifiedString(data, width, "0")
}

func (d *Dispatcher) leftJustifiedString(data string, width int, padding string) []byte {
	if len(data) == width {
		return []byte(data)
	}
	if len(data) > width {
		return []byte(data[0:width])
	}
	str := pad.Right(data, width, padding)
	return []byte(str)
}

func (d *Dispatcher) rightJustifiedString(data string, width int, padding string) []byte {
	if len(data) == width {
		return []byte(data)
	}
	if len(data) > width {
		return []byte(data[0:width])
	}
	str := pad.Left(data, width, padding)
	return []byte(str)
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
