package ahl

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

func (d *Dispatcher) GetProtocolType(station uint64) uint {
	protocol := d.GetConfig(station, defines.Protocol, cast.ToString(AHL4400_8))
	return cast.ToUint(protocol)
}

func (d *Dispatcher) GetDataTransferType(station uint64) uint {
	transfer := d.GetConfig(station, "DataTransfer", "0")
	return cast.ToUint(transfer)
}

func (d *Dispatcher) GetParserSlot(addr string) uint {
	station, _ := d.GetStationAddr(addr)
	return d.GetProtocolType(station)
}

func (d *Dispatcher) GetExtensionWidth(protocolType uint) int {

	switch protocolType {

	case AHL4400_5:
		return 5

	case AHL4400_8:
		return 8

	}

	return 0
}
