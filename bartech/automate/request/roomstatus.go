package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getNumeric(packet, "Extension")

	if extension == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	roomNumber := cast.ToString(extension)

	station, _ := dispatcher.GetStationAddr(addr)

	status := p.getString(packet, "Status")

	status = dispatcher.GetPMSMapping(packet.Name, station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: roomNumber,
		RoomStatus: status,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}
