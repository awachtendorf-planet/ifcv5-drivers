package telefon

import (
	"strings"

	"github.com/weareplanet/ifcv5-main/ifc/defines"

	"github.com/spf13/cast"
)

// GetProtocolByStation ...
func (d *Dispatcher) GetProtocolByStation(station uint64) string {
	protocol := d.GetConfig(station, defines.Protocol, "")
	vendor := cast.ToString(protocol)
	vendor = strings.ToUpper(vendor)
	return vendor
}

// GetProtocol ...
func (d *Dispatcher) GetProtocol(addr string) string {
	station, _ := d.GetStationAddr(addr)
	protocol := d.GetConfig(station, defines.Protocol, "")
	vendor := cast.ToString(protocol)
	vendor = strings.ToUpper(vendor)
	return vendor
}

// GetParserSlot ...
func (d *Dispatcher) GetParserSlot(addr string) uint {
	vendor := d.GetProtocol(addr)
	slot := d.getSlot(vendor)
	return slot
}

// GetSlotFromProtocol ...
func (d *Dispatcher) GetProtocolFromSlot(slot uint) string {
	d.slotGuard.RLock()
	for k, v := range d.slots {
		if v == slot {
			d.slotGuard.RUnlock()
			return k
		}
	}
	d.slotGuard.RUnlock()
	return ""
}

func (d *Dispatcher) getSlot(vendor string) uint {

	vendor = strings.ToUpper(vendor)
	if len(vendor) == 0 {
		return 0
	}

	d.slotGuard.RLock()
	slot := d.slots[vendor]
	d.slotGuard.RUnlock()

	return slot
}

func (d *Dispatcher) newSlot(vendor string) uint {

	vendor = strings.ToUpper(vendor)
	if len(vendor) == 0 {
		return 0
	}

	d.slotGuard.Lock()
	d.slotCounter++
	slot := d.slotCounter
	d.slots[vendor] = slot
	d.slotGuard.Unlock()

	return slot
}
