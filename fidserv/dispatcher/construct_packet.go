package fidserv

import (
	// "fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, recordName string, context interface{}, fields ...string) *ifc.LogicalPacket {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	var data string

	station, _ := d.GetStationAddr(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	var records []string

	if len(fields) == 0 { // get records from fias greeting, eg. GI -> G#GNGFGAGD
		records = d.GetLinkRecord(addr, recordName)
		if len(records) == 0 {
			log.Warn().Msgf("%T addr '%s' record name '%s' not requested", d, addr, recordName)
		}
	} else {
		records = fields // get records from provided fields
	}

	for _, key := range records {

		// pre check
		if d.filterField(packetName, key) {
			continue
		}

		if record, err := d.ConstructRecord(station, recordName, key, context); err == nil && len(record) > 0 {
			data = data + "|" + key + record
		} else {
			if d.AppendEmpty(station, recordName, key) {
				log.Warn().Msgf("%T addr '%s' record '%s' field '%s' not filled but delivered because of mandatory value", d, addr, recordName, key)
				data = data + "|" + key
			}
			if err != nil {
				log.Warn().Msgf("%T addr '%s' %s", d, addr, err)
			}
		}
	}
	if len(data) > 0 {
		data = data + "|"
	}
	packet.Add(payload, []byte(data))
	return packet
}

// ConstructRecord construct a fias record
func (d *Dispatcher) ConstructRecord(station uint64, recordName string, key string, context interface{}) (string, error) {

	switch recordName {

	case "LA", "LS", "DS", "DE", "LE", "NS", "NE", "WR", "WC":
		// pre-check to ignore key's at udf mapping
		break

	default:
		// get udf mapping, eg. GIA0 -> GN
		if udf, err := d.GetMapping("udf", station, recordName+key, false); err == nil {
			key = udf
		}
	}

	if context == nil || (recordName != "WR" && recordName != "WC" && recordName != "NS" && recordName != "NE" && recordName != "CK" && recordName != "XI" && recordName != "RE") {
		switch key {
		case "DA":
			return time.Now().Format("060102"), nil // YYMMDD

		case "TI":
			return time.Now().Format("150405"), nil // HHMMSS
		}
	}

	switch context.(type) {

	case *record.Guest:
		guest := context.(*record.Guest)
		value, err := d.getFromGuestRecord(station, recordName, key, guest)
		return value, err

	case *record.Generic:
		generic := context.(*record.Generic)
		value, err := d.getFromGenericRecord(station, recordName, key, generic)
		return value, err

	}

	return "", errors.Errorf("record '%s' field '%s' not found (context object '%T' unknown)", recordName, key, context)
}

// AppendEmpty returns true if a record/field is mandatory
func (d *Dispatcher) AppendEmpty(station uint64, recordName string, key string) bool {
	return d.MandatoryField(station, recordName, key)
}

func (d *Dispatcher) getFromGuestRecord(station uint64, recordName string, key string, guest *record.Guest) (string, error) {

	if guest == nil {
		return "", errors.Errorf("record '%s' field '%s' not found (guest record is nil)", recordName, key)
	}

	// fias key to guest struct member, eg. G# -> ReservationNumber
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
		// ESB ASW (Gx) has no RoomStatus field
		if recordName == "RE" && key == "RS" {
			return "", nil
		}

		// GC without OldRoomNumber, normal behaviour for guest data change without move
		if recordName == "GC" && key == "RO" {
			return "", nil
		}

		// virtual numbers not supported
		if (recordName == "GI" || recordName == "GO" || recordName == "GC") && (key == "EN" || key == "ES" || key == "EP") {
			return "", nil
		}
		if recordName == "GC" && (key == "EO" || key == "ET" || key == "EI") {
			return "", nil
		}

		return "", errors.Errorf("record '%s' field '%s' not mapped (guest record)", recordName, key)
	}

	return "", errors.Errorf("record '%s' field '%s' not found (guest record)", recordName, key)
}

func (d *Dispatcher) getFromGenericRecord(station uint64, recordName string, key string, generic *record.Generic) (string, error) {

	if generic == nil {
		return "", errors.Errorf("record '%s' field '%s' not found (generic record is nil)", recordName, key)
	}

	if value, exist := generic.Get(key); exist {
		return d.formatRecord(station, recordName, key, value)
	}

	if mapped, exist := d.GetMappedGenericField(key); exist {
		if value, exist := generic.Get(mapped); exist {
			return d.formatRecord(station, recordName, key, value)
		}
	}

	// suppress missing field warning

	// PA -> AS not OK
	if recordName == "PA" && (key == "G#" || key == "GN" || key == "RN") {
		return "", nil
	}
	// XT -> without MI/MT signals that no guest message exist
	if recordName == "XT" && (key == "MI" || key == "MT") {
		return "", nil
	}

	return "", errors.Errorf("record '%s' field '%s' not mapped (generic record)", recordName, key)
}

func (d *Dispatcher) formatRecord(station uint64, recordName string, key string, v interface{}) (string, error) {

	var value string

	switch key {

	case "DA", "GA", "GD": // date

		switch v.(type) {

		case time.Time:
			t := v.(time.Time)
			return t.Format("060102"), nil

		case string:
			s := v.(string)
			return d.formatDate(s)

		}

		return value, errors.Errorf("%s%s wrong data type '%T' expected 'time.Time' or 'string'", recordName, key, v)

	case "TI": // time HHMMSS

		switch v.(type) {

		case time.Time:
			t := v.(time.Time)
			return t.Format("150405"), nil

		case string:
			s := v.(string)
			return d.formatTime(s, "150405")
		}

		return value, errors.Errorf("%s%s wrong data type '%T' expected 'time.Time' or 'string'", recordName, key, v)

	case "DT": // departure time HH:MM

		switch v.(type) {

		case time.Time:
			t := v.(time.Time)
			return t.Format("15:04"), nil

		case string:
			s := v.(string)
			return d.formatTime(s, "15:04")
		}

		return value, errors.Errorf("%s%s wrong data type '%T' expected 'time.Time' or 'string'", recordName, key, v)

	case "CT": // clear text

		value = cast.ToString(v)

		if len(value) > 50 {
			value = value[:50]
		}
		value = strings.ToUpper(value)
		return value, nil

	default:
		// alias to module name, eg. TV = paytvright, GL = languagecode
		// get mapping from module, eg. paytvright 0 -> TU, languagecode en_US -> EA

		// cast.ToString casted boolean to "true/false"
		// correct it to "Y/N"
		switch v.(type) {
		case bool:
			b := v.(bool)
			if b {
				return "Y", nil
			}
			return "N", nil
		}

		value = cast.ToString(v)

		alias := d.ReadSetting("global", "alias", key, "")
		if len(alias) > 0 {
			if match, err := d.GetMapping(alias, station, value, false); err == nil {
				return match, nil
			}
		}
	}

	return value, nil
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
	return t.Format("060102"), nil // YYMMDD
}

func (d *Dispatcher) formatTime(value string, format string) (string, error) {
	var t = time.Time{}
	var err error

	if len(value) == 6 {
		t, err = time.Parse("150405", value)
	} else {
		t, err = d.parseTime(value)
	}

	if err != nil {
		return "", err
	}
	return t.Format(format), nil
}

func (d *Dispatcher) parseTime(value string) (time.Time, error) {
	// t, err := time.Parse("2006-01-02T15:04:05", value)
	// if err != nil {
	// 	t, err = time.Parse("2006-01-02T15:04:05Z", value)
	// }
	t, err := time.Parse(time.RFC3339, value)
	return t, err
}
