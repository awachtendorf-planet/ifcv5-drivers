package callstar

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// IsMove ...
func (d *Dispatcher) IsMove(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		if data, exist := guest.GetGeneric(defines.OldRoomName); exist {
			oldRoom := cast.ToString(data)
			if oldRoom != guest.Reservation.RoomNumber {
				return true
			}
		}
	}
	return false
}

// IsSharer
func (d *Dispatcher) IsSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}
	return false
}

// EncodeString ...
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

// DecodeString ...
func (d *Dispatcher) DecodeString(station uint64, data []byte) string {
	encoding := d.GetEncodingByStation(station)
	if len(encoding) == 0 {
		return string(data)
	}
	if dec, err := d.Decode(data, encoding); err == nil && len(dec) > 0 {
		return string(dec)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return string(data)
}
