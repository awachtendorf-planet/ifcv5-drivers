package inhova

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	sendActivationTime map[uint64]bool // station <-> config
	sendTrack2Data     map[uint64]bool // station <-> config
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[uint64]bool),
		make(map[uint64]bool),
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

func (d *Dispatcher) sendActivation(station uint64) bool {
	data := d.GetConfig(station, "SendActivationTime", "false")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) sendTrack2(station uint64) bool {
	data := d.GetConfig(station, "SendTrack2", "true")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) getKeyOptions(station uint64) string {
	data := d.GetConfig(station, "AccessPoints", "")
	options := cast.ToString(data)
	return options
}
