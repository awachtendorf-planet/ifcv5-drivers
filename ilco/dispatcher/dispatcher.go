package ilco

import (
	"github.com/spf13/cast"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	sendTrack2Data map[uint64]bool // station <-> config
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[uint64]bool),
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

func (d *Dispatcher) SendTrack2(station uint64) bool {
	data := d.GetConfig(station, "SendTrack2", "true")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) GetKeyOptionsDefault(station uint64) string {
	data := d.GetConfig(station, "KeyOptionsDefault", "")
	overlay := cast.ToString(data)
	return overlay
}

func (d *Dispatcher) IsGatewayMode(station uint64) bool {
	data := d.GetConfig(station, "Protocol", "0")
	gotGateway := cast.ToBool(data)
	return gotGateway
}

func (d *Dispatcher) GetAuthNumber(station uint64) string {
	data := d.GetConfig(station, "AuthNumber", "")
	authnumber := cast.ToString(data)
	return authnumber
}
