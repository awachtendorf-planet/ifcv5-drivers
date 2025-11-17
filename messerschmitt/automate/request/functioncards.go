package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func (p *Plugin) handleFunctionCards(addr string, packet *ifc.LogicalPacket) error {

	room := p.getRoom(addr, packet)
	if len(room) == 0 {
		return errors.New("room number is empty")
	}

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	function := p.getKeyAsInt("CardType", packet)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: room,
	}

	// serial 61-65
	// socket 00-31

	switch function {

	case 68, 06: // DND On
		roomStatus.DoNotDisturb = true

	case 69, 07: // DND Off

	default: // get room status from mapping
		status := dispatcher.GetPMSMapping("Room Status", station, "RS", cast.ToString(function))
		roomStatus.RoomStatus = status

	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

	return nil
}
