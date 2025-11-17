package dummy

/*
import (
	"reflect"
	"strings"

	"github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	tagField = "dummy"
)

// UnmarshalPacket construct a record object from a logical packet
func (d *Dispatcher) UnmarshalPacket(packet *ifc.LogicalPacket, p interface{}) error {

	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		return errors.New("must be a pointer")
	}
	if packet == nil {
		return errors.New("empty input buffer")
	}

	station, err := d.GetStationAddr(packet.Addr)

	reflectType := reflect.TypeOf(p).Elem()
	reflectValue := reflect.ValueOf(p).Elem()

	for i := 0; i < reflectType.NumField(); i++ {

		field := reflectType.Field(i)
		key, ok := field.Tag.Lookup(tagField)

		if !ok || len(key) == 0 || key == "-" {
			continue
		}

		args := strings.Split(key, ",")
		key = args[0]

		if value, exist := d.GetField(packet, key); exist {

			value = d.GetPMSMapping(packet.Name, station, key, value) // backward resolution vendor -> pms

			valueType := reflectValue.Field(i).Type()

			switch valueType.Kind() {

			case reflect.String:
				reflectValue.Field(i).SetString(value)

			case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
				value = strings.TrimLeft(value, "0")
				number := cast.ToInt64(value)
				reflectValue.Field(i).SetInt(number)

			case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
				value = strings.TrimLeft(value, "0")
				number := cast.ToUint64(value)
				reflectValue.Field(i).SetUint(number)

			case reflect.Float32, reflect.Float64:
				value = strings.TrimLeft(value, "0")
				number := cast.ToFloat64(value)
				reflectValue.Field(i).SetFloat(number)

			case reflect.Bool:
				if value == "Y" {
					reflectValue.Field(i).SetBool(true)
				} else if value == "N" {
					reflectValue.Field(i).SetBool(false)
				} else {
					state := cast.ToBool(value)
					reflectValue.Field(i).SetBool(state)
				}

			case reflect.Struct:
				err = errors.New("struct are not supported")

			default:
				err = errors.Errorf("unknown value type '%s'", valueType.Kind())
			}

		}
	}

	return err

}

// MarshalPacket construct a generic record object from a record
func (d *Dispatcher) MarshalPacket(p interface{}) (*record.Generic, error) {
	answer := record.NewGeneric()

	if reflect.TypeOf(p).Kind() == reflect.Ptr {
		return answer, errors.New("must not be a pointer")
	}

	var err error

	reflectType := reflect.TypeOf(p)
	reflectValue := reflect.ValueOf(p)

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectType.Field(i)
		key, ok := field.Tag.Lookup(tagField)

		if !ok || len(key) == 0 || key == "-" {
			continue
		}

		args := strings.Split(key, ",")
		key = args[0]

		valueType := reflectValue.Field(i).Type()

		switch valueType.Kind() {
		case reflect.String:
			value := reflectValue.Field(i).String()
			answer.Set(key, value)

		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			value := reflectValue.Field(i).Int()
			answer.Set(key, value)

		case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
			value := reflectValue.Field(i).Uint()
			answer.Set(key, value)

		case reflect.Float32, reflect.Float64:
			value := reflectValue.Field(i).Float()
			answer.Set(key, value)

		case reflect.Bool:
			value := reflectValue.Field(i).Bool()
			answer.Set(key, value)

		case reflect.Struct:
			err = errors.New("struct are not supported")

		default:
			err = errors.Errorf("unknown value type '%s'", valueType.Kind())
		}
	}

	return answer, err

}
*/
