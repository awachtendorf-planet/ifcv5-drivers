package informatel

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	// "github.com/weareplanet/ifcv5-main/log"
	"github.com/pkg/errors"
	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "GV", "GL"},
		"GC": {"RN", "GF", "GN", "GDN", "GV", "GL"},
		"GO": {"RN"},
		"RE": {"RN", "DN", "ML", "CS"},
		"WR": {"RN", "DA"},
		"WC": {"RN", "DA"},
	}
)

const (
	ifcType = "PBX"
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

	// db swap pre/post processor
	// d.DatabaseSyncPreProcess = d.databaseSyncPreProcess
	// d.DatabaseSyncPostProcess = d.databaseSyncPostProcess
}

func (d *Dispatcher) configureESBHandler() {

	// CI/CO/DC
	d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCGuestDataChangeNotifRQType{})

	// wake
	d.RegisterAsyncResponse(esb.IFCWakeupSetNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCWakeupDeleteNotifRQ{})

	// room status
	d.RegisterAsyncResponse(esb.IFCRoomEquipmentNotifRQType{})

}

func (d *Dispatcher) PreCheck(job *order.Job) error {
	if job == nil {
		return nil
	}

	switch job.Action {

	case order.Checkin, order.Checkout, order.RoomStatus, order.DataChange, order.WakeupRequest, order.WakeupClear:
		// OK

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}

	extension := d.GetRoom(job.Station, job.Context)

	if len(extension) > 10 {

		return errors.Errorf("extension '%s' is too long, max 10 characters allowed", extension)
	}

	return nil
}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	// TODO: remove in production
	if d.IsDebugMode(station) {
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
