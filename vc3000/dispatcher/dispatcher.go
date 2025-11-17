package vc3000

import (
	"sync"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	slots     map[string][]Slot
	slotGuard sync.RWMutex
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[string][]Slot),
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

func (d *Dispatcher) IsRadisson(station uint64) bool {
	return d.isRadisson(station)
}

// ProcessIncomingBytes debug function for local testing
func (d *Dispatcher) ProcessIncomingBytes(addr string, data *[]byte, size int) {

	station, err := d.GetStationAddr(addr)
	if err != nil || !d.IsDebugMode(station) {
		return
	}

	if data == nil || size == 0 {
		return
	}

	raw := *data
	for i := 0; i < size; i++ {
		if raw[i] == 0x10 || raw[i] == 'x' {
			raw[i] = 0x0
		}
	}
}

// GetParserSlot return 1 for socke layer and 2 for serial layer
func (d *Dispatcher) GetParserSlot(addr string) uint {
	if d.IsSerialDevice(addr) {
		return 2
	}
	return 1
}
