package saflok6000

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	interestedIn = map[string][]string{
		"KR": {"KC", "KT", "K#", "KO", "RN", "GS", "GD", "ID", "T2"},
		"KD": {"KC", "RN", "GS", "ID"},
	}
)

const (
	ifcType = "KES"
)

func (d *Dispatcher) initOverwrite() {

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.PreCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

	d.ConfigChanged = d.configChanged
}

func (d *Dispatcher) PreCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {

	case order.KeyRequest, order.KeyDelete:
		// OK

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())

	}

	if guest, ok := job.Context.(*record.Guest); ok {

		roomLength := d.GetRoomLength(job.Station)

		room := d.encode(job.Station, guest.Reservation.RoomNumber)
		if len(room) > roomLength {
			return errors.Errorf("room number '%s' too long (max: %d)", guest.Reservation.RoomNumber, roomLength)
		}

		if data, exist := guest.GetGeneric(defines.KeyOptions); exist {
			ap := cast.ToString(data)
			if len(ap) > 12 {
				return errors.Errorf("keyoptions '%s' too long (max: 12)", ap)
			}
		}

		if job.Action == order.KeyDelete && guest.Reservation.SharedInd {
			data := d.GetConfig(job.Station, "KeyDelete", "true")
			if cast.ToBool(data) {
				return errors.New("key delete with sharer flag is not allowed")
			}
		}

		data, _ := guest.GetGeneric(defines.EncoderNumber)
		encoder := cast.ToInt(data)

		min := 1
		if job.Action == order.KeyDelete {
			min = 0
		}
		if encoder < min || encoder > 99 {
			return errors.Errorf("encoder number must be in the range from %d to 99 (encoder: %d - %s)", min, encoder, job.Action.String())
		}

	}

	return nil
}

func (d *Dispatcher) configChanged(station uint64) {

	oldState := d.keyDelete[station]

	data := d.GetConfig(station, "KeyDelete", "true")
	newState := cast.ToBool(data)

	d.keyDelete[station] = newState

	if d.IsReady() && d.IsStationActive(station) && oldState != newState {
		d.loginStation(station)
	}

}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	addr := d.GetVendorAddr(station, 0)

	ifcType := d.GetConfig(station, defines.IFCType, ifcType)

	subscribe := record.Subscribe{
		Station: station,
		IFCType: ifcType,
	}

	data := d.GetConfig(station, "KeyDelete", "true")
	keyDelete := cast.ToBool(data)

	subscribe.MessageName = make(map[string][]string)

	for k, v := range interestedIn {
		if k == "KD" && keyDelete == false {
			continue
		}
		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
