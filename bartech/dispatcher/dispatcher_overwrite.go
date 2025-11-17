package bartech

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
		"GI": {"RN", "GF", "GN", "GDN", "GS"},
		"GC": {"RN", "GF", "GN", "GDN", "GS"},
		"GO": {"RN", "GS"},
		"RE": {"RN", "MR"},
		"NE": {"DA"},
	}
)

const (
	ifcType = "MBS"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	d.Acknowledgement = func(addr string, packetName string) bool { return true }

	// register async response packets
	d.ConfigureESBHandler = d.configureESBHandler

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.preCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

}

func (d *Dispatcher) configureESBHandler() {

	// CI/CO/DC
	d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCGuestDataChangeNotifRQType{})

	// room status
	d.RegisterAsyncResponse(esb.IFCRoomEquipmentNotifRQType{})

	// end of day
	d.RegisterAsyncResponse(esb.IOSysAdministrationCompleteNotifRQ{})

}

func (d *Dispatcher) preCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {

	case order.Checkin, order.Checkout, order.RoomStatus:
		// OK

	case order.NightAuditEnd:
		return nil

	case order.DataChange:

		if d.IsMove(job.Context) {
			return errors.New("vendor does not support room move")
		}

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())

	}

	extensionWidth := d.getExtensionWidth(job.Station)
	extension := d.getExtension(job.Context)

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

func (d *Dispatcher) PreCheckDataSwap(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {

	case order.Checkin, order.Checkout:

		if d.IsSharer(job.Context) {
			return errors.New("sharer flag is set")
		}

	}

	return d.preCheck(job)
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
