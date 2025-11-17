package request

import (
	"fmt"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleWakeupEvent(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	station, _ := dispatcherObj.GetStationAddr(addr)

	// extension
	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {

		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	// result code
	wakeupstatus := p.getField(packet, "WakeupStatus", false)

	wakeuptime := p.getField(packet, "WakeupTime", false)

	// wakeup time
	timestamp, err := time.Parse("1504", wakeuptime)
	if err != nil {

		timestamp = time.Now()
	}

	switch wakeupstatus {

	case "AC": // delete
		wake := record.WakeupClear{
			Station:  station,
			RoomName: extension,
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	case "DC": // add
		wake := record.WakeupRequest{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	case "EF", "OC", "SR", "FA": // normal answer
		wake := record.WakeupPerformed{
			Station:    station,
			RoomName:   extension,
			Success:    true,
			ResultCode: 1,
		}

		switch wakeupstatus {
		case "EF":
			wake.ResultCode = 1 // success
			// wake.Message = "OK"
		case "SR":
			wake.ResultCode = 2 // no response
			// wake.Message = "NR"
		case "OC":
			wake.ResultCode = 4 // busy
			// wake.Message = "BY"
		case "FA":
			wake.ResultCode = 6 // unprocessable
			// wake.Message = "UR"
		default:
			wake.ResultCode = 3 // retry
			// wake.Message = "RY"
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	default:

		fmt.Errorf("%s addr '%s' unknown reply cause '%s'", name, addr, wakeupstatus)
	}

	return nil

}
