package dummy

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	// "github.com/weareplanet/ifcv5-main/log"
	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"
	// "github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "GV", "DN", "ML", "GL", "GS", "CS"},
		"GC": {"RN", "GF", "GN", "GDN", "GV", "DN", "ML", "GL", "GS", "CS"},
		//"GO": {"RN", "GF", "GN", "GDN", "GV", "DN", "ML", "GL", "GS", "CS"},
		"GO": {"RN", "GS"},
		"RE": {"RN", "MR"},
		"KR": {"KC", "KT", "K#", "KO", "SI", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS", "T1", "T2", "T3"},
		"KD": {"KC", "SI", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS"},
		"KZ": {"KC", "ID", "WS"},
	}
)

const (
	ifcType = "KES"
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
	// switch job.Action {
	// case order.Checkin, order.Checkout, order.RoomStatus:
	// 	// OK
	// default:
	// 	return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	// }
	return nil
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
