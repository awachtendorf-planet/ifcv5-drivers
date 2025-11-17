package definity

import (
	"strings"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "DN", "ML", "GL", "GS", "CS"},
		"GC": {"RN", "GF", "GN", "GDN", "DN", "ML", "GL", "GS", "CS"},
		"GO": {"RN", "GS"},
		"RE": {"RN", "DN", "ML", "CS"},
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

}

func (d *Dispatcher) configureESBHandler() {

	// CI/CO/DC
	d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCGuestDataChangeNotifRQType{})

	// room status
	d.RegisterAsyncResponse(esb.IFCRoomEquipmentNotifRQType{})

}

func (d *Dispatcher) PreCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {

	case order.Checkin, order.Checkout, order.DataChange, order.RoomStatus:
		// OK

	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())

	}

	extensionWidth := d.getRoomLength(job.Station)
	extension := d.GetRoom(job.Context)

	if len(extension) == 0 {
		return errors.New("extension empty")
	}

	if len(extension) > extensionWidth {
		return errors.Errorf("extension '%s' too long (maximum width: %d)", extension, extensionWidth)
	}

	if strings.HasPrefix(extension, " ") {
		return errors.Errorf("extension '%s' has leading spaces", extension)
	}

	ext := cast.ToInt16(extension)
	if ext == 0 {
		return errors.Errorf("extension '%s' is not numerical", extension)
	}

	if d.isSharer(job.Context) {
		return errors.New("sharer flag is set")
	}

	if d.isMove(job.Context) {
		return errors.New("room move is not allowed")
	}

	return nil
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
