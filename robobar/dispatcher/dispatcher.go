package robobar

import (
	"strings"

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
	d.setDefaults()
	d.initOverwrite()
	d.StartupDispatcher()
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
}

// GetSequenceNumber load the last sale number
func (d *Dispatcher) GetSequenceNumber(station uint64) int {
	data := d.ReadValue(station, defines.Transaction, "-1")
	data = strings.TrimLeft(data, " ")
	data = strings.TrimLeft(data, "0")
	if len(data) == 0 {
		return -1
	}
	return cast.ToInt(data)
}

// SetSequenceNumber save the last sale number
func (d *Dispatcher) SetSequenceNumber(station uint64, sequence int) {
	d.StoreValue(station, defines.Transaction, cast.ToString(sequence))
}

func (d *Dispatcher) GetRoomLength(station uint64) int {
	data := d.GetConfig(station, "RoomLength", "4")
	data = strings.TrimLeft(data, " ")
	data = strings.TrimLeft(data, "0")
	if len(data) == 0 {
		return 4
	}
	return cast.ToInt(data)
}
