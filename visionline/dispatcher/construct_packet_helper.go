package visionline

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-drivers/visionline/template"

	"github.com/spf13/cast"
)

func (d *Dispatcher) getCommand(addr string, packetName string, context interface{}) byte {

	switch packetName {

	case template.PacketCodeCard:
		return 'A'

	case template.PacketReadKey:
		return 'B'

	case template.PacketAlive:
		return 'C'

	case template.PacketCardUpdate:
		return 'F'

	case template.PacketCheckout:
		return 'G'

	}

	return 0
}

func (d *Dispatcher) isMove(guest *record.Guest) bool {
	if guest != nil {
		if oldroom, exist := guest.GetGeneric(defines.OldRoomName); exist && oldroom != guest.Reservation.RoomNumber {
			return true
		}
	}
	return false
}

func (d *Dispatcher) getTrack2Origin(station uint64) uint {

	// 1 = PMS provides Track2 data to be sent to Vindcard
	// 2 = Vingcard returns Track2 data in KeyAnswer

	state := d.GetConfig(station, "Track2Origin", "0")
	return cast.ToUint(state)
}

func (d *Dispatcher) sendTrack4(station uint64, key string) bool {

	// send content of Track4 data in PP or PF field

	state := d.GetConfig(station, "Track4"+key, "false")
	return cast.ToBool(state)
}

func (d *Dispatcher) getCommonRooms(station uint64, guest *record.Guest) string {
	rooms := ""
	if guest == nil {
		return rooms
	}
	if value, exist := guest.GetGeneric(defines.KeyOptions); exist {
		ap := cast.ToString(value)
		index := 0
		for i := range ap {
			index++
			if ap[i] == '1' {
				if room, err := d.GetMapping("accesspoint", station, cast.ToString(index), true); err == nil {
					room = strings.Replace(room, ";", ",", -1)
					if len(room) > 0 {
						rooms = rooms + room + ","
					}
				}
			}
		}
		rooms = strings.TrimRight(rooms, ",")
	}
	return rooms
}

func (d *Dispatcher) getGuestRooms(guest *record.Guest) string {
	rooms := ""
	if guest == nil {
		return rooms
	}
	rooms += guest.Reservation.RoomNumber

	if data, exist := guest.GetGeneric(defines.AdditionalRooms); exist {
		if additionRooms, ok := data.(string); ok {
			additionRooms = strings.Replace(additionRooms, ";", ",", -1)
			if len(additionRooms) > 0 {
				rooms += ","
				rooms += additionRooms
				rooms = strings.TrimRight(rooms, ",")
			}
		}
	}

	return rooms
}

func (d *Dispatcher) formatDate(value string) (string, error) {
	var t = time.Time{}
	var err error

	if len(value) == 6 {
		t, err = time.Parse("060102", value)
	} else {
		t, err = d.parseTime(value)
	}

	if err != nil {
		return "", err
	}
	return t.Format("200601021504"), nil // YYYYMMDDHHMM
}

func (d *Dispatcher) parseTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	return t, err
}
