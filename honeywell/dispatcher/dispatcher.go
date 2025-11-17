package honeywell

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
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

// GetParserSlot return the template slot
func (d *Dispatcher) GetParserSlot(addr string) uint {

	station, _ := d.GetStationAddr(addr)
	protocol := d.GetProtocolType(station)

	switch protocol {

	case HONEYWELL_PROTOCOL:
		return HONEYWELL_PROTOCOL

	case ALERTON_PROTOCOL_1, ALERTON_PROTOCOL_2:
		return ALERTON_PROTOCOL_1

	}

	return 0
}

// GetProtocolType return the protocol type
func (d *Dispatcher) GetProtocolType(station uint64) uint {
	protocol := d.GetConfig(station, defines.Protocol, cast.ToString(HONEYWELL_PROTOCOL))
	return cast.ToUint(protocol)
}
