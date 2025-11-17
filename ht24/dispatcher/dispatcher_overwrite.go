package ht24

import (
	"github.com/pkg/errors"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

var (
	interestedIn = map[string][]string{
		"KR": {"KC", "KT", "K#", "KO", "SI", "RN", "GS", "G#", "GA", "GD", "ID", "T2", "UDF"},
		"KD": {"KC", "GS", "RN"},
	}
)

const (
	ifcType = "KES"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	d.Acknowledgement = func(addr string, packetName string) bool { return true }

	// register async response packets
	d.ConfigureESBHandler = d.configureESBHandler

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.PreCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

	d.ConfigChanged = d.configChanged

}

func (d *Dispatcher) configureESBHandler() {

	// CI/CO/DC
	// d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})
	// d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})

	// room status
	// d.RegisterAsyncResponse(esb.IFCRoomEquipmentNotifRQType{})

}

func (d *Dispatcher) PreCheck(job *order.Job) error {
	if job == nil {
		return nil
	}
	switch job.Action {
	case order.KeyDelete, order.KeyRequest:
		// OK
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}
	guest, ok := job.Context.(*record.Guest)
	if !ok {
		return errors.Errorf("context '%T' not supported", job.Context)
	}

	mainRoom := d.getRoom(job.Station, guest)
	if len(mainRoom) > RoomNumberLength {
		return errors.Errorf("room name '%s' too long (max: %d)", mainRoom, RoomNumberLength)
	} else if len(mainRoom) == 0 {
		return errors.New("empty room name")
	}

	encoder := d.getEncoder(guest)

	if job.Action == order.KeyRequest {
		if encoder < 0 {
			return errors.New("encoder number must be 0 or greater")
		}

		if encoder > 24 {
			return errors.New("encoder number must be 24 less")
		}
		keyType := d.GetKeyType(guest)
		if keyType == "N" && guest.Reservation.SharedInd {
			return errors.New("key request with key type new and sharer flag is not allowed")
		}
	}

	if job.Action == order.KeyDelete && guest.Reservation.SharedInd {
		if d.SendKeyDelete(job.Station) {
			return errors.New("key delete with sharer flag is not allowed")
		}
	}

	return nil
}

func (d *Dispatcher) configChanged(station uint64) {
	oldStateKD := d.sendKeyDelete[station]
	oldStateT2 := d.sendTrack2Data[station]
	oldStateUDF := d.sendUDF[station]

	d.sendKeyDelete[station] = d.SendKeyDelete(station)
	d.sendTrack2Data[station] = d.SendTrack2(station)
	d.sendUDF[station] = d.SendUDF(station)

	if d.IsReady() && d.IsStationActive(station) && (oldStateKD != d.sendKeyDelete[station] || oldStateT2 != d.sendTrack2Data[station] || oldStateUDF != d.sendUDF[station]) {
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

		if k == "KD" && !d.SendKeyDelete(station) {

			continue

		}
		if k == "KR" {

			if !d.SendTrack2(station) {

				for index, element := range v {

					if element == "T2" {

						v = append(v[:index], v[index+1:]...)

						break

					}
				}
			}

			if !d.SendUDF(station) {

				for index, element := range v {

					if element == "UDF" {

						v = append(v[:index], v[index+1:]...)

						break
					}
				}
			}
		}
		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
