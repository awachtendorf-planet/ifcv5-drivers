package mitel

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "ML", "CS", "GS"},
		"GC": {"RN", "GF", "GN", "GDN", "ML", "CS", "GS"},
		"GO": {"RN", "GS"},
		"RE": {"RN", "ML", "CS"},
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

	// database swap
	d.DatabaseSyncStart = d.databaseSyncStart
	d.DatabaseSyncEnd = d.databaseSyncEnd

	// possibly SendRestrictionRecord has changed, subscribe again
	d.ConfigChanged = d.configChanged

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

	case order.Checkin, order.Checkout, order.DataChange, order.WakeupRequest, order.WakeupClear, order.RoomStatus:
		// OK

	case order.NightAuditStart, order.NightAuditEnd:
		// fake job's for database swap start/end
		return nil

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())

	}

	extension := d.getExtension(job.Context)
	extensionWidth := d.getExtensionWidth(job.Station)

	if len(extension) == 0 {
		return errors.Errorf("extension not found in context '%T'", job.Context)
	}
	if len(extension) > extensionWidth {
		return errors.Errorf("extension '%s' too long (maximum width: %d)", extension, extensionWidth)
	}

	ext := cast.ToInt16(extension)
	if ext == 0 {
		return errors.Errorf("extension '%s' is not numerical", extension)
	}

	return nil
}

func (d *Dispatcher) configChanged(station uint64) {
	if !d.IsReady() {
		return
	}
	if d.IsStationActive(station) {
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

	sendRestriction := d.sendRestriction(station)

	for k, v := range interestedIn {
		for i := range v {
			if v[i] == "CS" && !sendRestriction {
				continue
			}
			subscribe.MessageName[k] = append(subscribe.MessageName[k], v[i])
		}
	}

	d.CreatePmsJob(addr, nil, subscribe)
}

func (d *Dispatcher) databaseSyncStart(station uint64) *order.Job {

	// create fake job for swap start

	job := &order.Job{
		Station: station,
		Action:  order.NightAuditStart,
		Context: nil,
	}
	return job
}

func (d *Dispatcher) databaseSyncEnd(station uint64) *order.Job {

	// create fake job for swap end

	job := &order.Job{
		Station: station,
		Action:  order.NightAuditEnd,
		Context: nil,
	}
	return job
}
