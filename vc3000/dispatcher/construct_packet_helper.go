package vc3000

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/vc3000/template"

	"github.com/spf13/cast"
)

func (d *Dispatcher) Is2800(addr string) bool {
	return d.is2800(addr)
}

func (d *Dispatcher) UseRmtCommandForKeyDelete(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	rmt := d.GetConfig(station, "UseRmtForKeyDelete", "false")
	state := cast.ToBool(rmt)
	return state
}

func (d *Dispatcher) formatUint32(value uint32) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err == nil {
		return buf.Bytes()
	}
	return []byte{0x0, 0x0, 0x0, 0x0}
}

func (d *Dispatcher) formatString(data string, size int16) []byte {
	x := make([]byte, size, size)
	if len(data) > 19 {
		copy(x[0:], data[0:19])
	} else {
		copy(x[0:], data)
	}
	return x
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

func (d *Dispatcher) getEncoder(context interface{}) uint16 {
	if g, ok := context.(*record.Guest); ok {
		if encoder, exist := g.GetGeneric(defines.EncoderNumber); exist {
			return cast.ToUint16(encoder)
		}
	}
	return 0
}

func (d *Dispatcher) getDestination(context interface{}) string {
	encoder := d.getEncoder(context)
	destination := cast.ToString(encoder)
	for len(destination) < 2 {
		destination = "0" + destination
	}
	return destination
}

func (d *Dispatcher) getAccessPoints(station uint64, context interface{}) (string, bool) {
	ap, err := d.GetConfigIfExist(station, defines.AccessPoints)
	return ap, err == nil
}

func (d *Dispatcher) getUserType(station uint64, context interface{}) string {
	key := d.getUserTypeKey(station, context)
	if match, err := d.GetMapping("usertype", station, key, false); err == nil {
		return match
	}
	return "DEFAULT"
}

func (d *Dispatcher) getUserGroup(station uint64, context interface{}) string {
	key := d.getUserGroupKey(station, context)
	if match, err := d.GetMapping("usergroup", station, key, false); err == nil {
		return match
	}
	return "DEFAULT"
}

func (d *Dispatcher) getAdditionalRooms(station uint64, context interface{}) string {
	if g, ok := context.(*record.Guest); ok {

		data, exist := g.GetGeneric(defines.AdditionalRooms)
		if !exist {
			return ""
		}

		additionalRooms, ok := data.(string)
		if !ok || len(additionalRooms) == 0 {
			return ""
		}

		var roomList []string
		if len(g.Reservation.RoomNumber) <= 7 {
			roomList = append(roomList, g.Reservation.RoomNumber)
		}

		rooms := strings.Split(additionalRooms, ";")
		for i := range rooms {
			room := rooms[i]
			if len(room) > 0 && len(room) <= 7 {
				roomList = append(roomList, room)
			}
			if len(roomList) >= 3 {
				break
			}
		}

		if len(roomList) > 0 {
			return strings.Join(roomList, ",")
		}
		return strings.Replace(additionalRooms, ";", ",", -1)

	}
	return ""
}

func (d *Dispatcher) existAdditionalRooms(station uint64, context interface{}) bool {
	if g, ok := context.(*record.Guest); ok {
		if data, exist := g.GetGeneric(defines.AdditionalRooms); exist {
			additionalRooms, ok := data.(string)
			if ok && len(additionalRooms) > 0 {
				return true
			}
		}
	}
	return false
}

func (d *Dispatcher) getUserTypeKey(station uint64, context interface{}) string {

	// value on premise = roomname
	// value radisson = x.te stelle aus accesspoints

	radisson := d.isRadisson(station)

	if g, ok := context.(*record.Guest); ok {
		if radisson {
			if value, exist := g.GetGeneric(defines.KeyOptions); exist {
				position := 1
				ap := cast.ToString(value)
				if position > 0 && len(ap) >= position {
					return ap[position-1 : position]
				}
			}
		} else {
			return g.Reservation.RoomNumber
		}
	}
	return ""
}

func (d *Dispatcher) getUserGroupKey(station uint64, context interface{}) string {

	// value on premise = ?
	// value radisson = x.te stelle aus accesspoints

	radisson := d.isRadisson(station)

	if g, ok := context.(*record.Guest); ok {
		if radisson {
			if value, exist := g.GetGeneric(defines.KeyOptions); exist {
				position := 2
				ap := cast.ToString(value)
				if position > 0 && len(ap) >= position {
					return ap[position-1 : position]
				}
			}
		} else {
			// ?
			//return g.Reservation.RoomNumber
		}
	}
	return ""
}

func (d *Dispatcher) isRadisson(station uint64) bool {
	isRadisson := d.GetConfig(station, "Radisson", "false")
	state := cast.ToBool(isRadisson)
	return state
}

func (d *Dispatcher) is2800(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	state := d.is2800FromStation(station)
	return state
}

func (d *Dispatcher) is2800FromStation(station uint64) bool {
	is2800 := d.GetConfig(station, "2800", "false")
	state := cast.ToBool(is2800)
	return state
}

func (d *Dispatcher) getCodeCardCommand(addr string, packetName string, context interface{}) byte {

	if g, ok := context.(*record.Guest); ok {

		is2800 := d.is2800(addr)
		keyType, _ := g.GetGeneric(defines.KeyType)

		switch packetName {

		case template.PacketCodeCard:

			switch cast.ToString(keyType) {

			case "N": // new
				if is2800 {
					return 'A'
				}
				return 'I' // check out old guest, check in new guest

			case "D": // duplicate
				if is2800 {
					return 0
				}
				return 'H' // add guest
			}

		case template.PacketCodeCardModify:

			if is2800 {
				return 0
			}

			switch cast.ToString(keyType) {

			case "N": // new
				return 'D' // replace guest card

			case "D": // duplicate
				return 'F' // replace guest user id
			}

		case template.PacketCheckout:
			return 'B'

		case template.PacketReadKey:
			return 'E'

		}

	}

	return 0
}
