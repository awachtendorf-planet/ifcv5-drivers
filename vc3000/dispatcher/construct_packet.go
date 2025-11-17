package vc3000

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	appName = "protel"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, isSerialLayer bool, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	if isSerialLayer {
		packet, err := d.constructPacketSerial(addr, packetName, tracking, context)
		return packet, err
	} else {
		packet, err := d.constructPacketSocket(addr, packetName, tracking, context)
		return packet, err
	}

	return nil, nil
}

func (d *Dispatcher) constructRecord(station uint64, recordName string, key string, context interface{}) (string, error) {

	switch key {

	case "A": // access points, check if configured
		if ap, exist := d.getAccessPoints(station, context); exist {
			return ap, nil
		}

	case "T": // user type
		userType := d.getUserType(station, context)
		return userType, nil

	case "U": // room type
		userGroup := d.getUserGroup(station, context)
		return userGroup, nil

	case "R": // guest room, suppress if we have additional rooms
		is2800 := d.is2800FromStation(station)
		exist := d.existAdditionalRooms(station, context)
		if recordName == "H" || recordName == "B" || (is2800 && recordName == "A") { // ignore additional rooms, im H und B packet existiert kein L (roomlist/additional rooms) Feld
			if exist {
				log.Warn().Msgf("%T additional rooms are not supported in %s record", d, recordName)
			}
			break // use R field
		}
		if exist {
			return "", nil // supress R field because we use the L field
		}

	case "L": // additional rooms, included main room if additional rooms exist
		rooms := d.getAdditionalRooms(station, context)
		return rooms, nil
	}

	switch recordName {

	// case "A", "B", "H" ... :
	// 	// pre-check to ignore key's at udf mapping
	// 	break

	default:
		// get udf mapping, eg. IF -> N
		if udf, err := d.GetMapping("udf", station, recordName+key, true); err == nil {
			key = udf
		}
	}

	switch context.(type) {

	case *record.Guest:
		guest := context.(*record.Guest)
		value, err := d.getFromGuestRecord(station, recordName, key, guest)
		return value, err

	}

	return "", errors.Errorf("record '%s' field '%s' not found (context object '%T' unknown)", recordName, key, context)
}

func (d *Dispatcher) getFromGuestRecord(station uint64, recordName string, key string, guest *record.Guest) (string, error) {

	if guest == nil {
		return "", errors.Errorf("record '%s' field '%s' not found (guest record is nil)", recordName, key)
	}

	// vc3000 key to guest struct member, eg. P -> ReservationNumber
	if mapped, exist := d.GetMappedGuestField(key); exist {
		if value, valid := guest.Get(mapped); valid {
			return d.formatRecord(station, recordName, key, value)
		}
	} else if !exist {

		if value, exist := guest.GetGeneric(key); exist {
			return d.formatRecord(station, recordName, key, value)
		}

		if mapped, exist := d.GetMappedGenericField(key); exist {
			if value, exist := guest.GetGeneric(mapped); exist {
				return d.formatRecord(station, recordName, key, value)
			}
		}

		// suppress missing field warning
		// if recordName == "A" && key == "R" {
		// 	return "", nil
		// }

		return "", errors.Errorf("record '%s' field '%s' not mapped (guest record)", recordName, key)
	}

	return "", errors.Errorf("record '%s' field '%s' not found (guest record)", recordName, key)
}

func (d *Dispatcher) formatRecord(station uint64, recordName string, key string, v interface{}) (string, error) {

	var value string

	switch key {

	case "D", "O": // date

		switch v.(type) {

		case time.Time:
			t := v.(time.Time)
			return t.Format("200601021504"), nil // YYYYMMDDHHMM

		case string:
			s := v.(string)
			return d.formatDate(s)

		}

		return value, errors.Errorf("%s%s wrong data type '%T' expected 'time.Time' or 'string'", recordName, key, v)

	default:
		// cast.ToString casted boolean to "true/false" correct it to "1/0"
		switch v.(type) {
		case bool:
			b := v.(bool)
			if b {
				return "1", nil
			}
			return "0", nil
		}

		value = cast.ToString(v)
		/*
			alias := d.ReadSetting("global", "alias", key, "")
			if len(alias) > 0 {
				if match, err := d.GetMapping(alias, station, value, false); err == nil {
					return match, nil
				}
			}
		*/
	}

	return value, nil
}

func (d *Dispatcher) constructData(addr string, recordName string, station uint64, context interface{}, dst *[]byte) error {

	if len(recordName) == 0 {
		return errors.New("internal error: record name is empty")
	}

	if commands == nil {
		return errors.New("internal error: commands is nil")
	}

	// slot 0 = vision, 1 = 2800
	is2800 := d.is2800(addr)
	slot := cast.ToUint(is2800)

	var buffer []byte
	data := make(map[string]string)

	records, exist := commands[slot][recordName[0]]
	if !exist {
		if is2800 {
			return errors.Errorf("command '%s' not supported (Vingcard 2800)", recordName)
		}
		return errors.Errorf("command '%s' not supported (Vingcard Vision)", recordName)
	}

	for i := range records {
		key := records[i].Field
		record, errConstruct := d.constructRecord(station, recordName, key, context)

		if errConstruct == nil {

			// encode
			if (key == "F" || key == "N" || key == "R" || key == "L" || key == "T" || key == "U" || key == "?") && len(record) > 0 {
				record = d.encode(addr, record)
			}

			// auto shrink fields
			if (key == "F" || key == "N" || key == "?") && len(record) > 15 {
				record = record[:15]
			} else if is2800 && key == "A" && len(record) > 1 {
				record = record[:1]
			} else if key == "A" && len(record) > 53 {
				record = record[:53]
			}
		}

		err := d.validate(record, records[i].Validate)

		if err != nil {
			required := records[i].Required
			if required {
				if errConstruct != nil {
					log.Error().Msgf("%T addr '%s' %s, validator err=%s, value '%s'", d, addr, errConstruct, err.Error(), record)
				} else {
					log.Error().Msgf("%T addr '%s' record '%s' field '%s' required, validator err=%s, value '%s'", d, addr, recordName, key, err.Error(), record)
				}
			} else {
				if errConstruct != nil {
					log.Warn().Msgf("%T addr '%s' %s, validator err=%s, value '%s'", d, addr, errConstruct, err.Error(), record)
				} else {
					log.Warn().Msgf("%T addr '%s' record '%s' field '%s', validator err=%s, value '%s'", d, addr, recordName, key, err.Error(), record)
				}
			}
		}

		if err == nil && len(record) > 0 {
			data[key] = record
		}

	}

	// request "S" if not checkout
	// means return the card serial number if available
	if !is2800 && recordName != "B" {
		if _, exist := data["S"]; !exist {
			data["S"] = ""
		}
	}

	// note:
	// range data -> unsorted map
	// field order are undefined

	// use range records for a sorted map
	for i := range records {
		key := records[i].Field
		if value, exist := data[key]; exist {
			buffer = append(buffer, 0x1e)
			buffer = append(buffer, key...)
			buffer = append(buffer, value...)
		}
	}

	// create dst buffer if not exist (serial layer variable length)
	// at socket layer it should be already 511 bytes alocated
	if dst == nil || len(*dst) == 0 {
		*dst = make([]byte, len(buffer))
	}

	if len(buffer) > len(*dst) {
		copy(*dst, buffer[:len(*dst)-1])
	} else {
		copy(*dst, buffer)
	}

	return nil

}
