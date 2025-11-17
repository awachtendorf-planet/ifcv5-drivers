package callstar

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"GI": {"RN", "G#", "GF", "GN", "SL", "GL", "GV", "GS"},
		"GC": {"RN", "G#", "GF", "GN", "SL", "GL", "GV", "GS", "RO"},
		"GO": {"RN", "G#", "GS"},
		"RE": {"RN", "G#", "CS", "ML", "GS"},
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
	d.PreCheckDriverJob = d.preCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

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

func (d *Dispatcher) preCheck(job *order.Job) error {
	if job == nil {
		return nil
	}

	switch job.Action {
	case order.Checkin, order.Checkout, order.DataChange, order.RoomStatus, order.WakeupRequest, order.WakeupClear:
		// OK
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
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

// GetSwapState ...
func (d *Dispatcher) GetSwapState(addr string) bool {
	d.swapStateGuard.RLock()
	state, exist := d.swapState[addr]
	d.swapStateGuard.RUnlock()
	return exist && state
}

// SetSwapSate ...
func (d *Dispatcher) SetSwapState(addr string, state bool) {
	d.swapStateGuard.Lock()
	if state {
		d.swapState[addr] = true
	} else {
		delete(d.swapState, addr)
	}
	d.swapStateGuard.Unlock()
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
