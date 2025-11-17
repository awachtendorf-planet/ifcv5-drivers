package plabor

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
	swapState      map[uint64]bool
	swapStateGuard sync.RWMutex
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[uint64]bool),
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

func (d *Dispatcher) DirectWakeupMode(station uint64) bool {
	data := d.GetConfig(station, "DirectWakeupMode", "false")
	return cast.ToBool(data)
}

func (d *Dispatcher) WakeUpDateWithYear(station uint64) bool {
	data := d.GetConfig(station, "WakeUpDateWithYear", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) LeadingZeros(station uint64) bool {
	data := d.GetConfig(station, "LeadingZeros", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) SyncCheckInsOnly(station uint64) bool {
	data := d.GetConfig(station, "SyncCheckInsOnly", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) Protocol(station uint64) int {

	// 0 = Plabor, 1 == Prodac
	data := d.GetConfig(station, "Protocol", "0")
	protocol := cast.ToInt(data)
	if protocol < 0 {
		protocol = 0
	} else if protocol > 1 {
		protocol = 1
	}
	return protocol
}

func (d *Dispatcher) RoomNumberAsReservationID(station uint64) bool {
	data := d.GetConfig(station, "RoomNumberAsReservationID", "false")
	return cast.ToBool(data)
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
func (d *Dispatcher) GetSwapState(station uint64) bool {
	d.swapStateGuard.RLock()
	state, exist := d.swapState[station]
	d.swapStateGuard.RUnlock()
	return exist && state
}

// SetSwapSate ...
func (d *Dispatcher) SetSwapState(station uint64, state bool) {
	d.swapStateGuard.Lock()
	if state {
		d.swapState[station] = true
	} else {
		delete(d.swapState, station)
	}
	d.swapStateGuard.Unlock()
}
