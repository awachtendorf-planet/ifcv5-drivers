package requestsimple

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket, retry bool) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	if retry {
		return 0, nil
	}

	station, _ := dispatcher.GetStationAddr(addr)

	status := p.getField(packet, "StatusCode", true)
	status = strings.TrimPrefix(status, "0")

	status = dispatcher.GetPMSMapping(packet.Name, station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: roomNumber,
		RoomStatus: status,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

	return 0, nil

}
