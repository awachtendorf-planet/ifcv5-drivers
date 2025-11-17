package honeywell

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GS"},
		"GO": {"RN", "GS"},
	}
)

const (
	ifcType = "BMS"
)

func (d *Dispatcher) initOverwrite() {

	// no ack/nak
	d.Acknowledgement = d.acknowledgement

	// register async response packets
	d.ConfigureESBHandler = d.configureESBHandler

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.PreCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

}

func (d *Dispatcher) configureESBHandler() {

	d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})

}

func (d *Dispatcher) acknowledgement(addr string, _ string) bool {

	station, _ := d.GetStationAddr(addr)
	protocol := d.GetProtocolType(station)

	switch protocol {

	case HONEYWELL_PROTOCOL: // no low level ACK/NAK
		return false

	case ALERTON_PROTOCOL_1, ALERTON_PROTOCOL_2: // use low level ACK/NAK
		return true

	}

	return false
}

func (d *Dispatcher) PreCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {

	case order.Checkin, order.Checkout: // OK

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())

	}

	room := d.GetRoom(job.Context)

	err := d.CheckRoom(room)

	return err
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
		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
