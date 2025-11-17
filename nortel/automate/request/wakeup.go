package request

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/nortel/template"
)

func (p *Plugin) handleWakeup(addr string, packet *ifc.LogicalPacket) {

	dispatcher := p.GetDispatcher()
	name := p.GetName()

	extension := p.getExtension(packet)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	switch packet.Name {

	case template.PacketWakeupSet:

		timestamp := p.getValueAsString(packet, "Time")

		if len(timestamp) == 3 { // expected hhmm
			timestamp = "0" + timestamp
		}

		if constructed, err := time.Parse("1504", timestamp); err == nil {

			wake := record.WakeupRequest{
				Station:  station,
				RoomName: extension,
				Time:     constructed,
			}
			dispatcher.CreatePmsJob(addr, packet, wake)

		} else {
			log.Warn().Msgf("%s addr '%s' wakeup set failed, parse timestamp err=", name, addr, err)
		}

	case template.PacketWakeupClear:

		wake := record.WakeupClear{
			Station:  station,
			RoomName: extension,
		}

		dispatcher.CreatePmsJob(addr, packet, wake)

	case template.PacketWakeupAnswer:

		state := p.getValueAsString(packet, "State")

		wake := record.WakeupPerformed{
			Station:  station,
			RoomName: extension,
		}

		switch state {

		case "AN":
			wake.ResultCode = 1 // ok

		case "BL":
			wake.ResultCode = 6 // unprocessable

		case "RE":
			wake.ResultCode = 2 // no response

		default:
			wake.ResultCode = 6 // unprocessable

		}

		if wake.ResultCode == 1 {
			wake.Success = true
		}

		dispatcher.CreatePmsJob(addr, packet, wake)

	}

}
