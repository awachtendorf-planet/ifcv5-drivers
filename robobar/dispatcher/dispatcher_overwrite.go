package robobar

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GS"},
		"GO": {"RN", "GS"},
		"RE": {"RN", "MR"},
	}
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

	// db swap pre/post processor to create RE packets (lock/unlock bar)
	// d.DatabaseSyncPreProcess = d.databaseSyncPreProcess
	// d.DatabaseSyncPostProcess = d.databaseSyncPostProcess
}

func (d *Dispatcher) configureESBHandler() {

	// CI/CO/DC
	d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})

	// room status
	d.RegisterAsyncResponse(esb.IFCRoomEquipmentNotifRQType{})

}

func (d *Dispatcher) PreCheck(job *order.Job) error {
	if job == nil {
		return nil
	}

	switch job.Action {
	case order.Checkin, order.Checkout, order.RoomStatus:
		// OK
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}

	if d.IsSharer(job.Context) {
		return errors.New("sharer are not supported")
	}

	pmsRoom := d.GetRoom(job.Context)
	room := d.encodeString(d.GetVendorAddr(job.Station, 1), pmsRoom)

	if len(pmsRoom) == 0 {
		return errors.New("empty room")
	}

	roomLength := d.GetRoomLength(job.Station)
	if roomLength > 0 && len(room) > roomLength {
		return errors.Errorf("room '%s' too long (maximum width: %d)", pmsRoom, roomLength)
	}

	for i := range room {
		if room[i] < 0x20 || room[i] > 0x7e {
			return errors.Errorf("room '%s' contains an invalid character '%c' ", pmsRoom, room[i])
		}
	}

	if job.Action == order.RoomStatus {
		mb := d.GetMinibarRight(job.Context)
		if mb != 0 && mb != 2 {
			return errors.Errorf("unknown minibar right: %d (expected 0 or 2)", mb)
		}
	}

	return nil
}

func (d *Dispatcher) encodeString(addr string, data string) string {
	encoding := d.GetEncoding(addr)
	if len(encoding) == 0 {
		return data
	}
	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
		return string(enc)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return data
}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	addr := d.GetVendorAddr(station, 0)

	ifcType := d.GetConfig(station, defines.IFCType, "MBS")

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
