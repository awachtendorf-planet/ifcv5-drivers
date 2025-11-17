package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getExtension(packet)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	status := p.getStatus(packet)

	status = dispatcher.GetPMSMapping(packet.Name, station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}
