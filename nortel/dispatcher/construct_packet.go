package nortel

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
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/nortel/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}, action order.Action) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	station, _ := d.GetStationAddr(addr)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	width := d.getExtensionWidth(station)
	ext := d.getExtension(context)
	ext = pad.Left(ext, width, " ")
	packet.Add("Ext", []byte(ext))

	switch packetName {

	case template.PacketCheckInExtension, template.PacketCheckOutExtension:
		// ok, nothing more to fill

	case template.PacketDisplayName:

		var displayName string
		if action != order.Checkout {
			displayName = d.getGuestName(context)
		} else if d.sendRoomVacant(station) {
			displayName = "Room Vacant"
		}
		displayName = d.formatDisplayName(station, displayName)

		packet.Add("Name", []byte(displayName))

	case template.PacketRoomStatus:

		var state string
		if guest, ok := context.(*record.Guest); ok {
			if value, exist := guest.GetGeneric(defines.RoomStatus); exist {
				if state, ok = value.(string); ok {
					if mapped, err := d.GetMapping("roomstatus", station, state, false); err == nil && len(mapped) > 0 {
						state = mapped
					}
				}
			}
		}

		state = strings.Trim(state, " ")
		if len(state) != 2 {
			state = "RE" // cleaning request
		}

		packet.Add("State", []byte(state))

	case template.PacketLanguage:

		var language string
		if guest, ok := context.(*record.Guest); ok {
			language = guest.Language
			if mapped, err := d.GetMapping("languagecode", station, language, false); err == nil && len(mapped) > 0 {
				language = mapped
			}
		}

		language = strings.Trim(language, " ")

		if len(language) > 2 {
			language = language[:2]
		}

		if len(language) == 0 {
			language = "0" // default language
		}

		packet.Add("Lang", []byte(language))

	case template.PacketMessageLamp:

		state := "OF"
		if action != order.Checkout {
			if guest, ok := context.(*record.Guest); ok {
				if guest.Reservation.MessageLightStatus {
					state = "ON"
				}
			}
		}

		packet.Add("State", []byte(state))

	case template.PacketVipState:

		state := "OF"
		if action != order.Checkout {
			if guest, ok := context.(*record.Guest); ok {
				if guest.VIPStatus == "1" || guest.VIPStatus == "Y" {
					state = "ON"
				}
			}
		}

		packet.Add("State", []byte(state))

	case template.PacketDoNotDisturb:

		state := "OF"
		if action != order.Checkout {
			if guest, ok := context.(*record.Guest); ok {
				if guest.Reservation.DoNotDisturb {
					state = "ON"
				}
			}
		}

		packet.Add("State", []byte(state))

	case template.PacketClassOfService:

		// default block all

		level := "CO"
		state := "ON"

		if action != order.Checkout {

			if guest, ok := context.(*record.Guest); ok {

				switch guest.Rights.ClassOfService {

				case 0: // block all
					break

				case 1: // local

					level = "E1"
					state = "ON"

				case 2: // national

					level = "E2"
					state = "ON"

				default: // no restriction

					level = "CO"
					state = "OF"

				}

			}

		}

		packet.Add("Level", []byte(level))
		packet.Add("State", []byte(state))

		/*
			case template.PacketSetCCRS, template.PacketSetECC1, template.PacketSetECC2:

				state := "OF"

				// TODO ggf. noch abhängig von der job.Action

				// gesperrt	0 	CO ON
				// local 		1	E1 ON
				// national 	2	E2 ON
				// all 		3	CO OF



				if action != order.Checkout {

					if guest, ok := context.(*record.Guest); ok {

						if packetName == template.PacketSetCCRS {
							if guest.Rights.ClassOfService >= 1 {
								state = "ON"
							}
						}

						if packetName == template.PacketSetECC1 {
							if guest.Rights.ClassOfService >= 2 {
								state = "ON"
							}
						}

						if packetName == template.PacketSetECC2 {
							if guest.Rights.ClassOfService >= 3 {
								state = "ON"
							}
						}

					}
				}

				packet.Add("State", []byte(state))
		*/

	case template.PacketWakeupSet:

		timestamp := d.getWakeupTime(context)
		if len(timestamp) != 4 {
			return packet, errors.Errorf("packet '%s' failed, wake time invalid", packetName)
		}
		packet.Add("Time", []byte(timestamp))

	case template.PacketWakeupClear:
		// ok, nothing more to fill

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) getExtension(context interface{}) string {
	if guest, ok := context.(*record.Guest); ok {
		return strings.Trim(guest.Reservation.RoomNumber, " ")
	}
	if generic, ok := context.(*record.Generic); ok {
		if extension, exist := generic.Get(defines.WakeExtension); exist {
			extension := cast.ToString(extension)
			return strings.Trim(extension, " ")
		}
	}
	return ""
}

func (d *Dispatcher) getGuestName(context interface{}) string {
	guest, ok := context.(*record.Guest)
	if !ok {
		return ""
	}
	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {
		return ""
	}
	if len(guest.DisplayName) > 0 {
		return guest.DisplayName
	}
	return fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
}

func (d *Dispatcher) formatDisplayName(station uint64, name string) string {

	// remove illegal characters
	if len(name) > 0 {
		name = strings.Trim(name, " ")
		re := strings.NewReplacer("*", "", ":", "", "\"", "", "\r", "")
		name = re.Replace(name)
	}

	// encode name
	name = d.encode(station, name)

	// replace non-ascii characters
	if len(name) > 0 {
		check := []byte(name)
		found := false
		for i := range check {
			if check[i] >= 0x7f {
				check[i] = '?'
				found = true
			}
		}
		if found {
			name = string(check)
		}
	}

	// 27 maximum
	// 23 recommended
	// -2 for start and end quote character

	length := d.getGuestNameLength(station)
	if length < 3 {
		length = 3
	} else if length > 27 {
		length = 27
	}

	length = length - 2 // start and end quote

	if len(name) > length {
		name = name[:length]
	}

	if d.sendGuestNameLength(station) {

		// 5 minimum? -> siehe Feature Description Page 45
		// -2 for start and end quote character

		if len(name) < 3 {
			name = pad.Right(name, 3, " ")
		}
		return fmt.Sprintf("\"%s\" %d", name, len(name)+2)
	}

	return fmt.Sprintf("\"%s\"", name)
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

func (d *Dispatcher) encode(station uint64, data string) string {
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
