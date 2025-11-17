package detewe

import (
	"fmt"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-main/utils"
	// "github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/detewe/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, jobAction order.Action, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	case template.PacketTCPLogin:

		packet.Add("Config", []byte("+10+20+40+41+60+67+70+71+72+80"))

		return packet, nil

	case template.PacketTG00:

		return packet, nil
	}

	station, _ := d.GetStationAddr(addr)

	extension := d.GetExtension(station, context)

	packet.Add("Participant", []byte(extension))

	switch packetName {

	case template.PacketTG41:

		displayType := d.DisplayType(station)

		var displayName string

		if jobAction == order.Checkout {

			displayName = fmt.Sprintf("% 16s", " ")
		} else {

			displayName = d.getDisplayName(station, context)
		}

		packet.Add("Text", []byte(displayName))
		packet.Add("DisplayType", []byte(displayType))

	case template.PacketTG60:

		right := "7" // 0x37 Checkin 0x30 Checkout

		if jobAction == order.Checkout {

			right = "0"
		}

		packet.Add("Right", []byte(right))

	case template.PacketTG67:
		//

	case template.PacketTG80:

		action := "0"

		if jobAction == order.RoomStatus {

			action = d.getMessageLightStatus(context) // 0x30 ausschalten 0x31 einschalten
		}

		packet.Add("Action", []byte(action))

	case template.PacketTG70:
		//

	case template.PacketTG71:

		wakeupAction := "2" // WakeupSet 0x32; WakeupClear 0x30

		if jobAction == order.WakeupClear {
			wakeupAction = "0"
		}

		packet.Add("ControlCode", []byte(wakeupAction))

		wakeupTime, _ := d.getWakeupTime(context)
		packet.Add("Year", []byte(wakeupTime.Format("2006")))
		packet.Add("Month", []byte(wakeupTime.Format("01")))
		packet.Add("Day", []byte(wakeupTime.Format("02")))
		packet.Add("Hour", []byte(wakeupTime.Format("15")))
		packet.Add("Minute", []byte(wakeupTime.Format("04")))

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

func (d *Dispatcher) GetExtension(station uint64, context interface{}) string {
	var extension string

	switch representer := context.(type) {

	case *record.Guest:
		extension = d.EncodeString(station, representer.Reservation.RoomNumber)

	case *record.Generic:
		roomNo, _ := representer.Get(defines.WakeExtension)
		roomNumber := cast.ToString(roomNo)
		extension = d.EncodeString(station, roomNumber)
	}

	if len(extension) < 4 {

		extension = fmt.Sprintf("% 4s", extension)
	}

	return extension
}

func (d *Dispatcher) getDisplayName(station uint64, context interface{}) string {

	guest := context.(*record.Guest)

	var displayName string

	if guest == nil || (len(guest.DisplayName) == 0 && len(guest.LastName) == 0 && len(guest.FirstName) == 0) {

		displayName = fmt.Sprintf("% 16s", "not set")
	} else if len(guest.DisplayName) > 0 {

		displayName = guest.DisplayName
	} else if len(guest.FirstName) == 0 {

		displayName = guest.LastName
	} else {

		displayName = fmt.Sprintf("%s, %s", guest.LastName, guest.FirstName)
	}

	displayName = d.EncodeString(station, displayName)

	if len(displayName) > 16 {

		displayName = displayName[:16]

	} else if len(displayName) < 16 {

		displayName = fmt.Sprintf("% 16s", displayName)
	}
	return displayName
}

func (d *Dispatcher) getWakeupTime(context interface{}) (time.Time, bool) {
	generic := context.(*record.Generic)
	if wakeupTime, exist := generic.Get(defines.Timestamp); exist {
		if t, ok := wakeupTime.(time.Time); ok && !t.IsZero() {
			return t, true
		}
	}
	return time.Now(), false
}

func (d *Dispatcher) getMessageLightStatus(context interface{}) string {
	guest := context.(*record.Guest)
	messageLightStatus := guest.Reservation.MessageLightStatus

	mlt := cast.ToBool(messageLightStatus)
	if mlt {
		return "1"
	}

	return "0"
}
