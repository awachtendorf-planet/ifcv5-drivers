package detewe

import (
	"github.com/spf13/cast"
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

// GetParserSlot return 1 for socke layer and 2 for serial layer
func (d *Dispatcher) GetParserSlot(addr string) uint {
	if d.IsSerialDevice(addr) {
		return 1
	}
	return 2
}

func (d *Dispatcher) DisplayType(station uint64) string {
	data := d.GetConfig(station, "DisplayType", "0")
	displayType := cast.ToString(data)
	return displayType
}

func (d *Dispatcher) GDXMode(station uint64) bool {
	data := d.GetConfig(station, "GDXMode", "true")
	mode := cast.ToBool(data)
	return mode
}

func (d *Dispatcher) SendTG67(station uint64) bool {
	data := d.GetConfig(station, "SendTG67", "false")
	sendTG67 := cast.ToBool(data)
	return sendTG67
}

func (d *Dispatcher) SendTG80(station uint64) bool {
	data := d.GetConfig(station, "SendTG80", "false")
	sendTG80 := cast.ToBool(data)
	return sendTG80
}

func (d *Dispatcher) DirectWakeupMode(station uint64) bool {
	data := d.GetConfig(station, "DirectWakeupMode", "false")
	wakeupMode := cast.ToBool(data)
	return wakeupMode
}
