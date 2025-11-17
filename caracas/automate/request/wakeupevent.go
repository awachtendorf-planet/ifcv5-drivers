package request

import (
	"fmt"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleWakeupOrder(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	station, _ := dispatcherObj.GetStationAddr(addr)

	// extension
	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {

		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	wakeuptime := p.getField(packet, "Time", false)
	order := p.getField(packet, "Status", false)

	// wakeup time
	timestamp, err := time.Parse("1504", wakeuptime)
	if err != nil {

		timestamp = time.Now()
	}

	switch order {

	case "0": // delete
		wake := record.WakeupClear{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	case "1": // set
		wake := record.WakeupRequest{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	default:

		return fmt.Errorf("%s addr '%s' unknown order '%s'", name, addr, order)
	}

	return nil

}

func (p *Plugin) handleWakeupResult(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	station, _ := dispatcherObj.GetStationAddr(addr)

	// extension
	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {

		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	wakeupdatetime := p.getField(packet, "DateTime", false)
	status := p.getField(packet, "Status", false)

	// wakeup time
	timestamp, err := time.Parse("0601021504", wakeupdatetime)
	if err != nil {

		timestamp = time.Now()
	}

	resultCode := 2
	if status == "1" {
		resultCode = 1
	}

	wakeupresult := record.WakeupPerformed{
		Station:    station,
		RoomName:   extension,
		Time:       timestamp,
		ResultCode: resultCode,
	}
	dispatcherObj.CreatePmsJob(addr, packet, wakeupresult)

	return nil

}
