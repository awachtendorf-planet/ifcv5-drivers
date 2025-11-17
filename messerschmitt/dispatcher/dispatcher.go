package messerschmitt

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
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

// GetTImeZone return the configured time zone
func (d *Dispatcher) GetTimeZone(station uint64) string {
	timezone := d.GetConfig(station, "TimeZone", "")
	return timezone
}

// GetParserSlot return 1 for socke layer and 2 for serial layer
func (d *Dispatcher) GetParserSlot(addr string) uint {
	if d.IsSerialDevice(addr) {
		return 1
	}
	return 2
}
