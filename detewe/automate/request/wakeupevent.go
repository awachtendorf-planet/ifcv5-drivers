package request

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleWakeupEvent(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	// extension
	extension := p.getField(packet, "Participant", true)
	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	controlcode := p.getField(packet, "ControlCode", false)

	// result code
	result := p.getField(packet, "Result", false)

	year := p.getField(packet, "Year", false)
	month := p.getField(packet, "Month", false)
	day := p.getField(packet, "Day", false)
	hour := p.getField(packet, "Hour", false)
	minute := p.getField(packet, "Minute", false)

	// wakeup time
	timestamp, err := time.Parse("200601021504", year+month+day+hour+minute)
	if err != nil {
		timestamp = time.Now()
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	switch controlcode {

	case "0": // delete
		wake := record.WakeupClear{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	case "2": // add
		wake := record.WakeupRequest{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	case "W", "A": // normal answer
		wake := record.WakeupPerformed{
			Station:    station,
			RoomName:   extension,
			Time:       timestamp,
			Success:    true,
			ResultCode: 1,
		}

		switch result {
		case "72":
			wake.ResultCode = 4 // busy
			// wake.Message = "busy line"
		case "73":
			wake.ResultCode = 2 // no response
			// wake.Message = "no answer"
		case "74":
			wake.ResultCode = 6 // unprocessable
			wake.Message = "invalid number"
		}
		dispatcherObj.CreatePmsJob(addr, packet, wake)

	default:
		log.Warn().Msgf("%s addr '%s' unknown reply cause '%s'", name, addr, result)
	}

}
