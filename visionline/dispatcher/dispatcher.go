package visionline

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
	d.setDefaults()
	d.initOverwrite()
	d.StartupDispatcher()
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
}

// GetParserSlot return 1 for socke layer and 2 for serial layer
func (d *Dispatcher) GetParserSlot(addr string) uint {
	if d.IsSerialDevice(addr) {
		return 2
	}
	return 1
}

// LogOutgoingPacket print a logical packet
func (d *Dispatcher) LogOutgoingPacket(name string, packet *ifc.LogicalPacket) {
	payload, _ := d.GetPayload(packet)

	sl := log.With().Logger()
	if len(packet.Tracking) > 0 {
		sl = sl.With().Str("ref", packet.Tracking).Logger()
	}

	log.Info().Msgf("%s addr '%s' send packet '%s' %s", name, packet.Addr, packet.Name, payload)

}
