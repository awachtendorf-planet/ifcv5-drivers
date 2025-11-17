package definity

import (
	"sync"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	messageCounter      map[uint64]byte
	messageCounterMutex sync.RWMutex
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		make(map[uint64]byte),
		sync.RWMutex{},
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

// GetParserSlot return which template slot are to use
func (d *Dispatcher) GetParserSlot(addr string) uint {
	station, _ := d.GetStationAddr(addr)
	return d.getProtocolType(station)
}

func (d *Dispatcher) getProtocolType(station uint64) uint {
	protocol := d.GetConfig(station, defines.Protocol, cast.ToString(ASCII_MODE))
	return cast.ToUint(protocol)
}

func (d *Dispatcher) getProtocolFormat(station uint64) uint {
	format := d.GetConfig(station, "RecordFormat", cast.ToString(STANDARD_FORMAT))
	return cast.ToUint(format)
}

func (d *Dispatcher) getMessageCounter(station uint64) byte { // 0-9

	d.messageCounterMutex.Lock()
	defer d.messageCounterMutex.Unlock()

	counter, exist := d.messageCounter[station]
	if !exist || counter >= 9 {
		d.messageCounter[station] = 0
		return 0
	}

	counter++
	d.messageCounter[station] = counter

	return counter
}

func (d *Dispatcher) getRoomLength(station uint64) int {
	format := d.getProtocolFormat(station)
	switch format {
	case EXTENDED_FORMAT:
		return 7
	}
	return 5
}

func (d *Dispatcher) getNameLength(station uint64) int {
	format := d.getProtocolFormat(station)
	switch format {
	case EXTENDED_FORMAT:
		return 30
	}
	return 15
}

func (d *Dispatcher) getCoveragePath(station uint64) string {
	cp := d.GetConfig(station, "CoveragePath", "0000")
	if len(cp) > 4 {
		cp = cp[:4]
	}
	return cp
}

func (d *Dispatcher) isMove(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		if data, exist := guest.GetGeneric(defines.OldRoomName); exist {
			oldRoom := cast.ToString(data)
			if oldRoom != guest.Reservation.RoomNumber {
				return true
			}
		}
	}
	return false
}

func (d *Dispatcher) isSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}
	return false
}

func (d *Dispatcher) encodeString(station uint64, data string) string {

	encoding := d.GetEncodingByStation(station)

	if len(encoding) == 0 {
		encoding = DEFAULT_ENCODING
	}

	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
		return string(enc)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}

	return data
}
