package nortel

import (
	"strings"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	esb "github.com/weareplanet/ifcv5-main/ifc/record/esb/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	interestedIn = map[string][]string{

		"GI": {"RN", "GF", "GN", "GDN", "GV", "GL", "GS"}, // "DN", "ML", "CS",
		"GC": {"RN", "GF", "GN", "GDN", "GV", "GL", "GS"}, // "DN", "ML", "CS",
		"GO": {"RN", "GS"},

		"RE": {"RN", "DN", "ML", "CS"},

		"WR": {"RN", "DA"},
		"WC": {"RN", "DA"},
	}
)

const (
	ifcType = "PBX"
)

type deviceType struct {
	Management         bool
	BackgroundTerminal bool
}

func (d *Dispatcher) initOverwrite() {

	// ack on management device if not background terminal
	d.Acknowledgement = d.acknowledgement

	// register async response packets
	d.ConfigureESBHandler = d.configureESBHandler

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.preCheck

	// find addr of management device if configured
	d.GetDriverAddr = d.getDriverAddr

	// subscribe if station is ready
	d.LoginStation = d.loginStation

	// cache handling
	d.HandleLinkState = d.handleLinkState
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

func (d *Dispatcher) configChanged(station uint64) {

	if !d.IsReady() {
		return
	}

	// clear device type cache, maybe "ManagementDevice" number has been changed
	if addrs, exist := d.GetStationLinkUpAddrs(station); exist {
		for i := range addrs {
			d.deviceCache.Delete(addrs[i])
		}
	}

	// clear regexp cache, maybe "DialledNumberRegexp" has been changed
	d.regexpCache.Delete(station)

	// is "ManagementDevice" 0 or >0 changed
	devNumber := d.getManagementDeviceNumber(station)
	currentState := devNumber > 0

	d.managementGuard.RLock()
	oldState, exist := d.management[station]
	d.managementGuard.RUnlock()

	login := false

	if !exist || oldState != currentState {

		d.managementGuard.Lock()
		delete(d.management, station)
		d.managementGuard.Unlock()

		login = true
	}

	// is "HandleWakeup" changed
	wake := d.handleWakeup(station)
	if oldWake, exist := d.wake[station]; !exist || oldWake != wake {
		d.wake[station] = wake
		login = true
	}

	// is "SendClassOfService" changed
	cos := d.SendClassOfService(station)
	if oldCos, exist := d.cos[station]; !exist || oldCos != cos {
		d.cos[station] = cos
		login = true
	}

	if login {
		d.loginStation(station)
	}

}

func (d *Dispatcher) handleLinkState(addr string, state bool) {
	if state {
		d.setDeviceType(addr)
	} else {
		d.deviceCache.Delete(addr)
	}
}

// acknowledgement return true if management device and not background terminal
func (d *Dispatcher) acknowledgement(addr string, _ string) bool {
	dt := d.getDeviceType(addr)
	if dt.Management && !dt.BackgroundTerminal {
		return true
	}
	return false
}

func (d *Dispatcher) getDeviceType(addr string) deviceType {
	if dt, exist := d.deviceCache.Get(addr); exist {
		if dt, ok := dt.(deviceType); ok {
			return dt
		}
	}
	dt := d.setDeviceType(addr)
	return dt
}

func (d *Dispatcher) setDeviceType(addr string) deviceType {
	dt := deviceType{}
	dt.Management = d.isManagementDevice(addr)
	if dt.Management {
		dt.BackgroundTerminal = d.isBackgroundTerminal(addr)
	}
	d.deviceCache.Set(addr, dt)
	return dt
}

func (d *Dispatcher) haveManagementDevice(station uint64) bool {

	d.managementGuard.RLock()
	state, exist := d.management[station]
	d.managementGuard.RUnlock()

	if exist && state {
		return true
	}

	devNumber := d.getManagementDeviceNumber(station)
	state = devNumber > 0

	d.managementGuard.Lock()
	d.management[station] = state
	d.managementGuard.Unlock()

	return state
}

// existManagementDevice return true if a management device is configured
func (d *Dispatcher) existManagementDevice(station uint64) bool {
	if devNumber := d.getManagementDeviceNumber(station); devNumber > 0 {
		if _, err := d.GetDeviceConfigByOwner(station, uint16(devNumber)); err == nil {
			return true
		}
	}
	return false
}

// getManagementDeviceAddr return the physical addr of the management device
// pcon-local:1234567890:1:127.0.0.1#43550
func (d *Dispatcher) getManagementDeviceAddr(station uint64) (string, bool) {
	devNumber := d.getManagementDeviceNumber(station)
	if devNumber == 0 {
		return "", false
	}
	if addrs, exist := d.GetStationLinkUpAddrs(station); exist {
		for i := range addrs {
			addr := addrs[i]
			if dev, _ := d.GetDeviceNumber(addr); dev == uint16(devNumber) {
				return addr, true
			}
		}
	}
	return "", false
}

// isManagementDevice return true if the physical addr matched the management device number
func (d *Dispatcher) isManagementDevice(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	if devNumber := d.getManagementDeviceNumber(station); devNumber > 0 {
		currentDevNumber, _ := d.GetDeviceNumber(addr)
		if uint16(devNumber) == currentDevNumber {
			return true
		}
	}
	return false
}

// isBackgroundTerminal return true if the physical addr matched the management device
// and the managament device are configured for background terminal mode
func (d *Dispatcher) isBackgroundTerminal(addr string) bool {
	// normalize addr pcon-local:1234567890:1:COM1 -> pcon-local:1234567890:1
	addr, _ = d.GetCustomerAddr(addr)
	if config, err := d.GetDeviceConfigByAddress(addr); err == nil {
		return config.IsBackgroundTerminal()
	}
	return false
}

// getDriverAddr find a physical addr from the management device
// pcon-local:1234567890:1:127.0.0.1#42832
func (d *Dispatcher) getDriverAddr(station uint64, _ order.QueueType, _ *order.Job) ([]string, bool) {
	var addrs []string
	if addr, exist := d.getManagementDeviceAddr(station); exist {
		addrs = append(addrs, addr)
	}
	return addrs, len(addrs) > 0
}

func (d *Dispatcher) preCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	if !d.existManagementDevice(job.Station) {
		return errors.New("no management device found")
	}

	switch job.Action {

	case order.Checkin, order.Checkout, order.RoomStatus:

	case order.WakeupRequest, order.WakeupClear:

		if !d.handleWakeup(job.Station) {
			return errors.New("wakeup handling are set to inactive")
		}

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

	extension = strings.TrimLeft(extension, "0")
	if ext := cast.ToInt16(extension); ext == 0 {
		return errors.Errorf("extension '%s' is not numerical", extension)
	}

	if addr, exist := d.getManagementDeviceAddr(job.Station); !exist || d.GetLinkState(addr) != automatestate.LinkUp {
		return errors.Errorf("no link-up device for InterfaceID %d found", job.Station)
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

	if d.haveManagementDevice(station) {

		handleWake := d.handleWakeup(station)
		handleCOS := d.SendClassOfService(station)

		for k, v := range interestedIn {

			if !handleWake && (k == "WR" || k == "WC") {
				continue
			}

			if !handleCOS && k == "RE" {

				for index, element := range v {

					if element == "CS" {
						v = append(v[:index], v[index+1:]...)
						break
					}
				}

			}

			subscribe.MessageName[k] = v
		}
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
