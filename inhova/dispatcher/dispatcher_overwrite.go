package inhova

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"KR": {"KC", "KT", "K#", "KO", "SI", "RN", "GS", "G#", "GA", "GD", "ID", "T2"},
		"KD": {"KC", "GS", "RN"},
	}
)

const (
	ifcType = "KES"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	d.Acknowledgement = func(addr string, packetName string) bool { return true }

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.PreCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

	// config changed trigger
	d.ConfigChanged = d.configChanged

}

func (d *Dispatcher) PreCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {
	case order.KeyDelete, order.KeyRequest: // ok
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}

	guest, ok := job.Context.(*record.Guest)
	if !ok {
		return errors.Errorf("context '%T' not supported", job.Context)
	}

	if _, err := d.GetKeyType(job); err != nil {
		return err
	}

	room := d.getRoom(job.Station, guest)

	if len(room) > 18 {
		return errors.Errorf("room name '%s' too long (max: %d)", room, 18)
	} else if len(room) == 0 {
		return errors.New("empty room name")
	}

	if job.Action == order.KeyDelete && guest.Reservation.SharedInd {
		return errors.New("key delete with sharer flag is not allowed")
	}

	return nil

}

func (d *Dispatcher) configChanged(station uint64) {

	// map are safe without mutex, always called from one thread and never multiple times

	oldStateActivitionTime := d.sendActivationTime[station]
	oldStateTrack2Data := d.sendTrack2Data[station]

	newStateActivitionTime := d.sendActivation(station)
	newStateTrack2Data := d.sendTrack2(station)

	d.sendActivationTime[station] = newStateActivitionTime
	d.sendTrack2Data[station] = newStateTrack2Data

	if d.IsReady() && d.IsStationActive(station) && (oldStateActivitionTime != newStateActivitionTime || oldStateTrack2Data != newStateTrack2Data) {
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

	subscribe.MessageName = make(map[string][]string)

	for k, v := range interestedIn {

		if k == "KR" {

			var values []string

			for _, element := range v {

				if element == "GA" && !d.sendActivation(station) {
					continue
				}

				if element == "T2" && !d.sendTrack2(station) {
					continue
				}

				values = append(values, element)

			}

			subscribe.MessageName[k] = values
			continue

		}

		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
