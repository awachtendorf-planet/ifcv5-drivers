package ht24

import (
	"github.com/spf13/cast"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	sendKeyDelete  map[uint64]bool // station <-> config
	sendTrack2Data map[uint64]bool // station <-> config
	sendUDF        map[uint64]bool // station <-> config
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[uint64]bool),
		make(map[uint64]bool),
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

func (d *Dispatcher) SendKeyDelete(station uint64) bool {
	data := d.GetConfig(station, "SendKeyDelete", "true")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) SendTrack2(station uint64) bool {
	data := d.GetConfig(station, "SendTrack2", "true")
	state := cast.ToBool(data)
	return state
}

func (d *Dispatcher) GetAccesPointOverlay(station uint64) string {
	data := d.GetConfig(station, "AccessPointOverlay", "")
	overlay := cast.ToString(data)
	return overlay
}

func (d *Dispatcher) GetProtocol(station uint64) int {
	data := d.GetConfig(station, "Protocol", "1")
	protocol := cast.ToInt(data)
	switch protocol {
	case 2, 14:
		protocol = 2
	case 3, 15:
		protocol = 3
	default:
		protocol = 1
	}
	return protocol
}

func (d *Dispatcher) SendUDF(station uint64) bool {
	if d.GetProtocol(station) == 1 || len(d.GetKeyID(station)) == 0 {
		return false
	}
	return true

}

func (d *Dispatcher) GetKeyID(station uint64) string {
	data := d.GetConfig(station, "KeyID", "")
	keyId := cast.ToString(data)
	return keyId
}
