package informatel

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
	d.initOverwrite()
	d.StartupDispatcher()
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
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

// SendVIP ...
func (d *Dispatcher) SendVIP(station uint64) bool {
	data := d.GetConfig(station, "SendVIP", "false")
	state := cast.ToBool(data)
	return state
}
