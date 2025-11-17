package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) {

	dispatcher := p.GetDispatcher()
	name := p.GetName()

	extension := p.getExtension(packet)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	status := p.getValueAsString(packet, "State")

	switch status {
	case "RE", "PR", "CL", "PA", "FA", "SK", "NS":
	default:
		log.Warn().Msgf("%s addr '%s' unknown room status '%s'", name, addr, status)
		return
	}

	status = dispatcher.GetPMSMapping(packet.Name, station, "RS", status)

	maid := p.getValueAsString(packet, "Maid")

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
		UserID:     maid,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}
