package caracas

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	bindMode       bool
	newVersionMode bool
	sendTime       int64
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		true,
		true,
		int64(0),
	}
}

// Startup start the driver and generic dispatcher
func (d *Dispatcher) Startup() {
	log.Info().Msgf("startup %T", d)
	d.initOverwrite()
	d.StartupDispatcher()
}

func (d *Dispatcher) GetBindMode(station uint64) bool {
	data := d.GetConfig(station, "BindMode", "true")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) GetNewVersionMode(station uint64) bool {
	data := d.GetConfig(station, "NewVersionMode", "true")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) GetSendTimeMS(station uint64) int64 {
	data := d.GetConfig(station, "SendTime", "0")
	number := cast.ToInt64(data)
	return number
}

// GetTimeZone return the configured time zone
func (d *Dispatcher) GetTimeZone(station uint64) string {
	timezone := d.GetConfig(station, "TimeZone", "")
	return timezone
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
}
