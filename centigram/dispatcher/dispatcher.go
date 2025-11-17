package centigram

import (
	"sync"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	swapState      map[string]bool
	swapStateGuard sync.RWMutex
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[string]bool),
		sync.RWMutex{},
	}
}

// Startup start the driver and generic dispatcher
func (d *Dispatcher) Startup() {
	log.Info().Msgf("startup %T", d)
	d.setDefaults()
	d.initOverwrite()
	d.StartupDispatcher()
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
}

// SwapRequest ...
func (d *Dispatcher) SwapRequest(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	data := d.GetConfig(station, "SwapRequest", "true")
	state := cast.ToBool(data)
	return state
}

// SendGuestName ...
func (d *Dispatcher) SendGuestName(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	data := d.GetConfig(station, "SendGuestName", "false")
	state := cast.ToBool(data)
	return state
}

// SendGuestLanguage ...
func (d *Dispatcher) SendGuestLanguage(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	data := d.GetConfig(station, "SendGuestLanguage", "false")
	state := cast.ToBool(data)
	return state
}

// IsMove ...
func (d *Dispatcher) IsMove(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		if data, exist := guest.GetGeneric(defines.OldRoomName); exist {
			oldRoom := cast.ToString(data)
			if oldRoom != guest.Reservation.RoomNumber {
				return true
			}
		}
	}
	return false
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

// GetParserSlot return 1 for HIS and 2 for Encore Protocol
func (d *Dispatcher) GetParserSlot(addr string) uint {
	station, _ := d.GetStationAddr(addr)
	return d.getProtocolType(station)
}

func (d *Dispatcher) getProtocolType(station uint64) uint {
	protocol := d.GetConfig(station, defines.Protocol, cast.ToString(HIS_PROTOCOL))
	return cast.ToUint(protocol)
}
