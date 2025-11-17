package centigram

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-drivers/centigram/template"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "GL", "ML", "VM", "GS"},
		"GC": {"RN", "GF", "GN", "GDN", "GL", "ML", "VM", "GS", "RO"},
		"GO": {"RN", "GS"},
		"RE": {"RN", "VM"},
	}
)

const (
	ifcType = "VMS"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	d.Acknowledgement = func(addr string, packetName string) bool { return true }

	// register addition Refusion/MessageWaitingStatus packet for ack/nak router
	d.RegisterAcknowledgementPackets = func() []string {
		packets := []string{
			template.PacketRefusion,
			template.PacketMessageWaitingStatusOnSwap,
		}
		return packets
	}

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

}

func (d *Dispatcher) preCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	extensionWidth := 6

	switch job.Action {

	case order.Checkin, order.Checkout, order.RoomStatus:
		// OK

	case order.DataChange:

		if d.IsMove(job.Context) {

			// check old room/extension

			guest, _ := job.Context.(*record.Guest)
			source := d.getOldExtension(guest)

			if len(source) == 0 {
				return errors.Errorf("old extension not found in context '%T'", job.Context)
			}
			if len(source) > extensionWidth {
				return errors.Errorf("old extension '%s' too long (maximum width: %d)", source, extensionWidth)
			}

			ext := cast.ToInt16(source)
			if ext == 0 {
				return errors.Errorf("old extension '%s' is not numerical", source)
			}

		}

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())

	}

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

		if guest, ok := job.Context.(*record.Guest); ok {
			if guest.Reservation.SharedInd {
				return errors.New("sharer flag is set")
			}
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
