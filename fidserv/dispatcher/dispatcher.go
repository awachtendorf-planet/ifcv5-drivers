package fidserv

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"

	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	records *utils.S3Map
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		utils.NewS3Map(),
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

// LogOutgoingPacket print a logical packet
func (d *Dispatcher) LogOutgoingPacket(name string, packet *ifc.LogicalPacket) {
	payload, _ := d.GetPayload(packet)

	sl := log.With().Logger()
	if len(packet.Tracking) > 0 {
		sl = sl.With().Str("ref", packet.Tracking).Logger()
	}

	sl.Info().Msgf("%s addr '%s' send packet '%s' %s", name, packet.Addr, packet.Name, payload)

}

func (d *Dispatcher) getGuestMessageHandling(station uint64) int {
	data := d.GetConfig(station, "GuestMessageHandling", "0")
	value := cast.ToInt(data)
	return value
}
