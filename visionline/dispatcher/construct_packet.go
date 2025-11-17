package visionline

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	command := d.getCommand(addr, packetName, context)
	if command == 0 {
		return nil, errors.New("command not supported")
	}

	recordName := string(command)
	records, exist := commands[command]
	if !exist {
		return nil, errors.Errorf("command '%s' not supported", recordName)
	}

	station, _ := d.GetStationAddr(addr)
	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	var buffer []byte
	data := make(map[string]string)

	for i := range records {
		key := records[i].Field
		record, errConstruct := d.constructRecord(station, recordName, key, context)

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

	// use range records for a sorted map
	for i := range records {
		key := records[i].Field
		if value, exist := data[key]; exist {
			buffer = append(buffer, key...)
			buffer = append(buffer, value...)
			buffer = append(buffer, 0x3b)
		}
	}

	packet.Add(payload, buffer)
	return packet, nil
}

func (d *Dispatcher) constructRecord(station uint64, recordName string, key string, context interface{}) (string, error) {

	switch recordName {

	case "C":

		switch key {

		case "EA":
			return "HEARTBEAT", nil

		case "AM":
			return "1", nil

		default:
			return "", nil

		}

	default:
		// get udf mapping
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

	switch recordName {

	case "F": // change card or move guest

		if d.isMove(guest) { // switch new room and old room

			switch key {

			case "NR": // new room
				key = "GR"

			case "GR": // guest room
				key = "OR"
			}

		}

	case "G": // checkout guest

		// suppress UI field (reservation number) if not sharer reservation
		// Visionline kann offensichtlich nicht damit umgehen, sondern erwartet beim Checkout des letzten Gastes, das das UI Feld nicht gefüllt ist.
		// Das hat protel die letzten 25 Jahre falsch gemacht, und 5 Euro in das "Opera macht das aber so" Schweindel.

		if key == "UI" && guest.Reservation.SharedInd == false {
			return "", nil
		}

	}

	switch key {

	case "AM":

		return "1", nil

	case "CR":

		rooms := d.getCommonRooms(station, guest)
		return rooms, nil

	case "GR":

		rooms := d.getGuestRooms(guest)
		return rooms, nil

	case "JR":

		keyType, _ := guest.GetGeneric(defines.KeyType)
		if keyType == "D" || (guest.Reservation.SharedInd && keyType == "N") {
			return "1", nil
		}
		return "", nil

	case "OT":

		keyType, _ := guest.GetGeneric(defines.KeyType)
		if keyType == "O" {
			return "1", nil
		}
		return "", nil

	case "PP", "PF":

		if d.sendTrack4(station, key) {
			key = defines.Track4
		}

	case "SR":

		if recordName == "B" { // key read
			return "?", nil
		}

		if value := d.getTrack2Origin(station); value == 2 {
			return "?", nil
		}

	case "T2":

		if value := d.getTrack2Origin(station); value != 1 {
			return "", nil
		}
	}

	// visionline key to guest struct member, eg. UI -> ReservationNumber
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

		return "", errors.Errorf("record '%s' field '%s' not mapped (guest record)", recordName, key)
	}

	return "", errors.Errorf("record '%s' field '%s' not found (guest record)", recordName, key)
}

func (d *Dispatcher) formatRecord(station uint64, recordName string, key string, v interface{}) (string, error) {

	var value string

	switch key {

	case "CI", "CO": // date

		switch v.(type) {

		case time.Time:
			t := v.(time.Time)
			if t.IsZero() && key == "CO" && recordName == "B" {
				t = time.Now()
			}
			if t.IsZero() {
				return "", nil
			}
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
	}

	return value, nil
}
