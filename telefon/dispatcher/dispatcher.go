package telefon

import (
	"sync"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"

	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparserdynamisch"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	parser        *analyser.GenericProtocol
	lrcHandler    map[string]calculateFn
	lrcGuard      sync.RWMutex
	slots         map[string]uint
	slotCounter   uint
	slotGuard     sync.RWMutex
	protocol      map[uint]Protocol
	protocolGuard sync.RWMutex
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network, parser *analyser.GenericProtocol) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		parser,
		make(map[string]calculateFn),
		sync.RWMutex{},
		make(map[string]uint),
		0,
		sync.RWMutex{},
		make(map[uint]Protocol),
		sync.RWMutex{},
	}
}

// Startup start the driver and generic dispatcher
func (d *Dispatcher) Startup() {
	log.Info().Msgf("startup %T", d)
	d.setDefaults()
	d.initOverwrite()
	d.initLRCHandler()
	d.StartupDispatcher()
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
}

func (d *Dispatcher) getStationSetting(addr string, key string, fallback string) string {
	var value string
	if station, err := d.GetStationAddr(addr); err == nil {
		value = d.GetConfig(station, key, fallback)
	}
	return value
}
