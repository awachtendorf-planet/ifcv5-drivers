package request

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleWakeupEvent(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	// extension
	extension := p.getField(packet, "Extension", true)
	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	// reply code
	code := p.getField(packet, "Code", false)
	if len(code) < 4 /*|| code[0] != 'W' */ {
		log.Warn().Msgf("%s addr '%s' unknown reply code '%s'", name, addr, code)
		return
	}

	// wakeup time
	timestamp := time.Now()
	wakeupTime := p.getField(packet, "WakeupTime", true)
	if constructed, err := p.constructDate(wakeupTime); err == nil {
		timestamp = constructed
	} else {
		log.Warn().Msgf("%s addr '%s' construct wakeup time failed, err=%s", name, addr, err)
	}

	station, _ := dispatcher.GetStationAddr(addr)

	switch code[2] {

	case 'P': // add
		wake := record.WakeupRequest{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	case 'M': // modify
		wake := record.WakeupRequest{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	case 'C': // delete
		wake := record.WakeupClear{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	case 'A': // normal answer
		wake := record.WakeupPerformed{
			Station:    station,
			RoomName:   extension,
			Time:       timestamp,
			Success:    true,
			ResultCode: 1,
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	case 'B', 'N', 'O', 'R': // failed
		wake := record.WakeupPerformed{
			Station:  station,
			RoomName: extension,
			Time:     timestamp,
			Success:  false,
		}

		switch code[2] {
		case 'B':
			wake.ResultCode = 4 // busy
			// wake.Message = "busy line"
		case 'N':
			wake.ResultCode = 2 // no response
			// wake.Message = "no answer"
		case 'O':
			wake.ResultCode = 6 // unprocessable
			// wake.Message = "out of order line"
		case 'R':
			wake.ResultCode = 6 // unprocessable
			// wake.Message = "reject, max. number of wake-up already programmed"
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	default:
		log.Warn().Msgf("%s addr '%s' unknown reply cause '%c'", name, addr, code[2])
	}

}
