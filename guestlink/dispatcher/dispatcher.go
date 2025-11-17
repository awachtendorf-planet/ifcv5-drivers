package guestlink

import (
	"sync"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	transaction      map[string]interface{}
	owner            map[string]Transaction
	transactionGuard sync.RWMutex
	router           map[string]*ifc.Route
	routerGuard      sync.RWMutex
	tan              map[uint64]uint16
	tanGuard         sync.RWMutex
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[string]interface{}),
		make(map[string]Transaction),
		sync.RWMutex{},
		make(map[string]*ifc.Route),
		sync.RWMutex{},
		make(map[uint64]uint16),
		sync.RWMutex{},
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

// GetProtocolType ...
func (d *Dispatcher) GetProtocolType(station uint64) uint {
	protocol := d.GetConfig(station, defines.Protocol, cast.ToString(DrvGuestlink))
	return cast.ToUint(protocol)
}

// GetAccountNumberLength Vorbereitung für Maginet Enhanced
func (d *Dispatcher) GetAccountNumberLength(station uint64) int {
	drv := d.GetProtocolType(station)
	if drv == DrvMaginetEnhanced {
		return 8
	}
	return 6
}

// GetAccountNumberAlignment returns left or right, default left
func (d *Dispatcher) GetAccountNumberAlignment(station uint64) int {
	drv := d.GetProtocolType(station)
	if drv == DrvMovielink || drv == DrvQuadriga {
		return JustifiedRight
	}
	return JustifiedLeft
}

// SendHeloPacket sends the HELO packet if TripleGuest or Eclipse
func (d *Dispatcher) SendHeloPacket(station uint64) bool {
	drv := d.GetProtocolType(station)
	if drv == DrvTripleGuest || drv == DrvEclipse {
		return true
	}
	return false
}

// GetParserSlot Vorbereitung Maginet Enhanced/OnCommand (Guest Message)
func (d *Dispatcher) GetParserSlot(addr string) uint {
	return 1
}
