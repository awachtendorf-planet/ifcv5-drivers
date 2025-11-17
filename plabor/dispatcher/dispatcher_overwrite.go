package plabor

import (
	"github.com/pkg/errors"
	"github.com/weareplanet/ifcv5-drivers/plabor/template"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "DN", "ML", "GL", "TV", "VR"},
		"GC": {"RN", "GF", "GN", "GDN", "DN", "ML", "GL", "TV", "VR"},
		"GO": {"RN"},
		"XL": {"RN", "MT", "MI"},
		"XD": {"RN", "MT", "MI"},
	}
)

const (
	ifcType = "PTV"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	d.Acknowledgement = d.isAcknowledgement

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

func (d *Dispatcher) isAcknowledgement(addr string, packetName string) bool {

	return packetName != template.PacketEot

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
	case order.Checkin, order.Checkout, order.GuestMessageDelete, order.GuestMessageOnline:

		guest, ok := job.Context.(*record.Guest)
		if !ok {
			return errors.Errorf("context '%T' not supported", job.Context)
		}
		room := guest.Reservation.RoomNumber
		if len(room) > 4 {
			return errors.Errorf("room name '%s' too long (max: 4)", room)
		}
		resID := guest.Reservation.ReservationID
		if len(guest.Reservation.ReservationID) > 7 {
			return errors.Errorf("reservation id '%s' too long (max: 7)", resID)
		}
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
