package mitel

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
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

// SendAliveRecord ...
func (d *Dispatcher) SendAliveRecord(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	data := d.GetConfig(station, "AliveRecord", "true")
	state := cast.ToBool(data)
	return state
}

// SendRestrictionRecord
func (d *Dispatcher) SendRestrictionRecord(addr string) bool {
	station, _ := d.GetStationAddr(addr)
	state := d.sendRestriction(station)
	return !state
}

func (d *Dispatcher) getExtensionWidth(station uint64) int {
	data := d.GetConfig(station, "ExtensionWidth", "5")
	data = strings.TrimLeft(data, " ")
	data = strings.TrimLeft(data, "0")
	if len(data) == 0 {
		return 4
	}
	return cast.ToInt(data)
}

func (d *Dispatcher) sendRestriction(station uint64) bool {
	data := d.GetConfig(station, "RestrictionRecord", "false")
	state := cast.ToBool(data)
	return state
}
