package callstar

import (
	"sync"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
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

// SendSwapLabel ...
func (d *Dispatcher) SendSwapLabel(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	data := d.GetConfig(station, "SwapLabel", "true")
	state := cast.ToBool(data)
	return state
}

// TicketAsOutlet ...
func (d *Dispatcher) TicketAsOutlet(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	data := d.GetConfig(station, "MinibarMode", "0")
	state := cast.ToInt(data)
	return state == 1
}
