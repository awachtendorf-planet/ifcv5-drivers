package vc3000

// slot handling im socket layer

// das module sollte noch optimiert werden
// performance technisch ist das eher ranzig
// ein mutex über alle instanzen könnte zum problem werden

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type SlotEvent struct {
	Station uint64
}

func NewSlotEvent(station uint64) *SlotEvent {
	return &SlotEvent{Station: station}
}

type Slot struct {
	Addr      string
	Used      bool
	Timestamp int64
}

func (d *Dispatcher) AddSlot(addr string) {

	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return
	}

	d.slotGuard.Lock()
	defer func() {
		d.slotGuard.Unlock()
		d.slotAvailable(addr)
	}()

	slot, exist := d.slots[devAddr]
	if exist {

		for i := range slot {
			if slot[i].Addr == addr {
				if slot[i].Used == true {
					slot[i].Used = false
					slot[i].Timestamp = time.Now().UnixNano()
					d.slots[devAddr] = slot
				}
				return
			}
		}

		for i := range slot {
			if len(slot[i].Addr) == 0 {
				slot[i].Addr = addr
				slot[i].Used = false
				slot[i].Timestamp = time.Now().UnixNano()
				d.slots[devAddr] = slot
				return
			}
		}
	}
	d.slots[devAddr] = append(d.slots[devAddr], Slot{Addr: addr, Used: false, Timestamp: time.Now().UnixNano()})
}

func (d *Dispatcher) RemoveSlot(addr string) {

	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return
	}

	d.slotGuard.Lock()
	defer d.slotGuard.Unlock()

	slot, exist := d.slots[devAddr]
	if !exist {
		return
	}

	for i := range slot {
		if slot[i].Addr == addr {
			slot[i].Addr = ""
			slot[i].Used = false
			d.slots[devAddr] = slot
			return
		}
	}
}

func (d *Dispatcher) RemoveSlotIfUnsed(addr string) bool {

	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return false
	}

	d.slotGuard.Lock()
	defer d.slotGuard.Unlock()

	slot, exist := d.slots[devAddr]
	if !exist {
		return false
	}

	for i := range slot {
		if slot[i].Addr == addr {
			if !slot[i].Used {
				slot[i].Addr = ""
				d.slots[devAddr] = slot
				return true
			}
			break
		}
	}
	return false
}

func (d *Dispatcher) GetSlot(addr string) (string, error) {

	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return "", err
	}

	d.slotGuard.Lock()
	defer d.slotGuard.Unlock()

	slot, exist := d.slots[devAddr]
	if !exist {
		return "", errors.Errorf("slot '%s' does not exist", devAddr)
	}

	for i := range slot {
		if len(slot[i].Addr) > 0 && slot[i].Used == false {
			slotAddr := slot[i].Addr
			slot[i].Used = true
			d.slots[devAddr] = slot
			return slotAddr, nil
		}
	}
	return "", errors.Errorf("no free slot for '%s' found", devAddr)
}

func (d *Dispatcher) GetSlotInfo(addr string) (Slot, bool) {
	slot := Slot{}

	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return slot, false
	}

	d.slotGuard.RLock()
	defer d.slotGuard.RUnlock()

	slots, exist := d.slots[devAddr]
	if !exist {
		return slot, false
	}

	for i := range slots {
		if slots[i].Addr == addr {
			slot.Addr = slots[i].Addr
			slot.Used = slots[i].Used
			slot.Timestamp = slots[i].Timestamp
			return slot, true
		}
	}

	return slot, false
}

func (d *Dispatcher) GetSlotsCount(addr string) int {
	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return -1
	}

	d.slotGuard.RLock()
	slots, exist := d.slots[devAddr]

	if !exist {
		d.slotGuard.RUnlock()
		return -1
	}

	count := 0
	for i := range slots {
		if len(slots[i].Addr) > 0 {
			count++
		}
	}

	d.slotGuard.RUnlock()
	return count
}

func (d *Dispatcher) FreeSlot(addr string) {

	devAddr, err := d.GetDeviceAddr(addr)
	if err != nil {
		return
	}

	d.slotGuard.Lock()

	slot, exist := d.slots[devAddr]
	if !exist {
		d.slotGuard.Unlock()
		return
	}

	for i := range slot {
		if slot[i].Addr == addr {
			if slot[i].Used == true {
				slot[i].Used = false
				slot[i].Timestamp = time.Now().UnixNano()
				d.slots[devAddr] = slot
			}
			d.slotGuard.Unlock()
			d.slotAvailable(addr)
			return
		}
	}

	d.slotGuard.Unlock()
}

func (d *Dispatcher) slotAvailable(addr string) {
	station, err := d.GetStationAddr(addr)
	if err != nil {
		return
	}

	broker := d.Broker()
	if broker != nil {
		name := d.GetSlotEventName(station)
		broker.Broadcast(NewSlotEvent(station), name)
	}
}

func (d *Dispatcher) GetSlotEventName(station uint64) string {
	return fmt.Sprintf("FreeSlot:%d", station)
}

// func (d *Dispatcher) PrintSlots() {
// 	for k, v := range d.slots {
// 		fmt.Printf("slot: %s \n", k)
// 		slot := v
// 		for s := range slot {
// 			x := slot[s]
// 			fmt.Printf(" -> addr: %s used: %t\n", x.Addr, x.Used)

// 		}
// 	}
// }
