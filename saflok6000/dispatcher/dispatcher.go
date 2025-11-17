package saflok6000

import (
	"github.com/spf13/cast"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	keyDelete map[uint64]bool // station <-> config
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

func (d *Dispatcher) encode(station uint64, data string) string {
	encoding := d.GetEncodingByStation(station)
	if len(encoding) == 0 {
		return data
	}
	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
		return string(enc)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return data
}

// Get Room Length by Protocol
func (d *Dispatcher) GetRoomLength(station uint64) int {
	data := d.GetConfig(station, "Protocol", "1")
	protocol := cast.ToInt(data)
	switch protocol {
	case 1:
		// Default Length
		return 5

	case 2:
		// Extended Length
		return 15

	default:
		// Fallback to Default
		return 5
	}
}

// Cut Leading Zeroes
func (d *Dispatcher) LeadingZeroes(station uint64) bool {
	data := d.GetConfig(station, "LeadingZeroes", "false")
	leadingzeroes := cast.ToBool(data)
	return leadingzeroes
}
