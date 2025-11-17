package caracas

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
		"GI": {"RN", "G#", "GF", "GN", "GDN", "GV", "GG" /*, "CL"*/, "GL", "GS"},
		"GC": {"RN", "G#", "GF", "GN", "GDN", "GV", "GG" /*, "CL"*/},
		"GO": {"RN", "G#", "GS"},
		"RE": {"RN", "RS", "ML", "DN", "CS"},
		"WR": {"RN", "DA", "TI"},
		"WC": {"RN", "DA", "TI"},
	}
)

const (
	ifcType = "PBX"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	// d.Acknowledgement = func(addr string, packetName string) bool { return true }

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
	case order.Checkin, order.Checkout, order.DataChange, order.RoomStatus, order.WakeupRequest, order.WakeupClear, order.Internal:
		// OK
	default:
		return errors.Errorf("vendor does not supported this command (%s)", job.Action.String())
	}

	extension := d.getRoom(job.Station, job.Context)

	if len(extension) > 6 {
		return errors.Errorf("vendor supports max roomnumber length of 6 digits (%s)", extension)
	}

	resID := d.getReservationID(job.Station, job.Context)

	if len(resID) > 10 {
		return errors.Errorf("vendor supports max reservationID length of 10 digits (%s)", resID)
	}

	return nil
}

func (d *Dispatcher) configChanged(station uint64) {

	oldBindMode := d.bindMode
	oldNewVersionMode := d.newVersionMode

	d.bindMode = d.GetBindMode(station)
	d.newVersionMode = d.GetNewVersionMode(station)

	if d.IsReady() && d.IsStationActive(station) && (oldBindMode != d.bindMode || oldNewVersionMode != d.newVersionMode) {
		d.loginStation(station)
	}

}

// func (d *Dispatcher) encodeString(addr string, data string) string {
// 	encoding := d.GetEncoding(addr)
// 	if len(encoding) == 0 {
// 		return data
// 	}
// 	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
// 		return string(enc)
// 	} else if err != nil && len(encoding) > 0 {
// 		log.Warn().Msgf("%s", err)
// 	}
// 	return data
// }

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

		if k == "GC" && d.GetNewVersionMode(station) {

			v = append(v, "RO")

		}

		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
