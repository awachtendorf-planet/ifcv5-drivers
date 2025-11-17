package bartech

import (
	"strings"
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
// func (d *Dispatcher) SwapRequest(addr string) bool {
// 	station, _ := d.GetStationAddr(addr)
// 	data := d.GetConfig(station, "SwapRequest", "true")
// 	state := cast.ToBool(data)
// 	return state
// }

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

// IsSharer
func (d *Dispatcher) IsSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}
	return false
}

// UnlockBar ...
func (d *Dispatcher) UnlockBar(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		return guest.Rights.Minibar == 2
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

// GetParserSlot ...
func (d *Dispatcher) GetParserSlot(addr string) uint {
	return 1
}

func (d *Dispatcher) getExtensionWidth(station uint64) int {
	data := d.GetConfig(station, "ExtensionWidth", "6")
	data = strings.TrimLeft(data, " ")
	data = strings.TrimLeft(data, "0")
	if len(data) == 0 || data != "4" {
		return 6
	}
	return 4
}

func (d *Dispatcher) simpleCheckin(station uint64) bool {
	data := d.GetConfig(station, "SimpleCheckin", "false")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) sendHappyHour(station uint64) bool {
	data := d.GetConfig(station, "SendHappyHour", "false")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) sendGuestName(station uint64) bool {
	data := d.GetConfig(station, "SendGuestName", "false")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) sendCheckoutDate(station uint64) bool {
	data := d.GetConfig(station, "SendCheckoutDate", "false")
	state := cast.ToBool(data)
	return state
}
