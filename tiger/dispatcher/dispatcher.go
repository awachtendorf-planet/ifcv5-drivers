package tiger

import (
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

// IsMove ...
func (d *Dispatcher) IsMove(context interface{}) bool {
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

func (d *Dispatcher) AdjustnameLength(station uint64) bool {
	data := d.GetConfig(station, "AdjustnameLength", "false")
	return cast.ToBool(data)
}

func (d *Dispatcher) UseDisplayname(station uint64) bool {
	data := d.GetConfig(station, "UseNameDisplay", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) MessageOffCharacter(station uint64) string {
	data := d.GetConfig(station, "MessageOffCharacter", "0")
	if data != "0" && data != "2" {
		data = "0"
	}
	return data
}

func (d *Dispatcher) ShortRoomname(station uint64) bool {
	data := d.GetConfig(station, "ShortRoomname", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) SupportDND(station uint64) bool {
	data := d.GetConfig(station, "SupportDND", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) GetParserSlot(addr string) uint {
	station, _ := d.GetStationAddr(addr)
	return cast.ToUint(d.Protocol(station))
}

func (d *Dispatcher) Protocol(station uint64) int {
	data := d.GetConfig(station, "Protocol", "3")
	val := cast.ToInt(data)
	if val < 0 || val > 3 {
		val = 3
	}
	return val
}
