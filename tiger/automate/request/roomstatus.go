package request

import (
	"strings"

	// "github.com/weareplanet/ifcv5-drivers/tiger/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	// "github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "RoomNumber", true)
	extension = strings.TrimLeft(extension, "0")

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	status := p.getField(packet, "RoomStatusCode", true)
	status = dispatcherObj.GetPMSMapping(packet.Name, station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
	}

	dispatcherObj.CreatePmsJob(addr, packet, roomStatus)

}
