package ahl

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/ahl/template"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "GF", "GN", "GDN", "GV", "DN", "ML", "GL", "GS", "CS"},
		"GC": {"RN", "GF", "GN", "GDN", "GV", "DN", "ML", "GL", "GS", "CS"},
		//"GO": {"RN", "GF", "GN", "GDN", "GV", "DN", "ML", "GL", "GS", "CS"},
		"GO": {"RN"},
		"RE": {"RN", "DN", "ML", "CS"},
		"WR": {"RN", "DA"},
		"WC": {"RN", "DA"},
	}
)

func (d *Dispatcher) initOverwrite() {

	// register async response packets
	d.ConfigureESBHandler = d.configureESBHandler

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.precheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

	// always true beside LinkAlive packet because fucking ahl protocol logic
	d.Acknowledgement = d.acknowledgement

	// auto create CO/CI on DC if moves is not allowed
	d.MovesAllowed = func(station uint64, deviceNumber int) bool { return true }

}

func (d *Dispatcher) configureESBHandler() {

	// CI/CO/DC
	d.RegisterAsyncResponse(esb.IFCGuestDataChangeNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckOutNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCCheckInNotifRQType{})

	// wake
	d.RegisterAsyncResponse(esb.IFCWakeupSetNotifRQType{})
	d.RegisterAsyncResponse(esb.IFCWakeupDeleteNotifRQ{})

	// room status
	d.RegisterAsyncResponse(esb.IFCRoomEquipmentNotifRQType{})

}

func (d *Dispatcher) acknowledgement(addr string, packetName string) bool {
	if packetName == template.PacketLinkAlive && !d.IsSerialDevice(addr) {
		return false
	}
	return true
}

func (d *Dispatcher) precheck(job *order.Job) error {
	if job == nil {
		return nil
	}

	switch job.Action {
	case order.Checkin, order.Checkout, order.DataChange:
		// OK
	case order.WakeupRequest, order.WakeupClear:
		// OK
	case order.RoomStatus:
		// OK
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}

	protocolType := d.GetProtocolType(job.Station)
	extensionWidth := d.GetExtensionWidth(protocolType)
	extension := d.GetExtension(job.Context)

	if len(extension) == 0 {
		return errors.New("extension empty")
	}
	if len(extension) > extensionWidth {
		return errors.Errorf("extension '%s' too long (maximum width: %d)", extension, extensionWidth)
	}
	if strings.HasPrefix(extension, " ") {
		return errors.Errorf("extension '%s' has leading spaces", extension)
	}

	return nil
}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	addr := d.GetVendorAddr(station, 0)

	ifcType := d.GetConfig(station, defines.IFCType, "PBX")

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
