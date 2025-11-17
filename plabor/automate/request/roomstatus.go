package request

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.driver
	name := p.automate.Name

	station, _ := dispatcherObj.GetStationAddr(addr)

	extension := p.getField(packet, "RoomNumber")
	if p.driver.LeadingZeros(station) {
		extension = strings.TrimLeft(extension, "0")
	}

	userID := p.getField(packet, "Operator")

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	status := p.getField(packet, "Status")
	status = dispatcherObj.GetPMSMapping(packet.Name, station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
		UserID:     userID,
	}

	dispatcherObj.CreatePmsJob(addr, packet, roomStatus)

}
