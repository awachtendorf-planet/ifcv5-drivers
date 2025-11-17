package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleRoomStatus(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = 2 digit status
	// 1 = hhmm

	if len(roomNumber) == 0 || len(data) == 0 || len(data[0]) == 0 {
		return
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()

	station, _ := dispatcher.GetStationAddr(addr)

	status := p.getString(data[0])

	status = dispatcher.GetPMSMapping("Room Status", station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: roomNumber,
		RoomStatus: status,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}

func (p *Plugin) handleMessageLamp(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = + or -

	if len(roomNumber) == 0 || len(data) == 0 || len(data[0]) == 0 {
		return
	}

	if data[0] != "+" && data[0] != "-" {
		return
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()

	station, _ := dispatcher.GetStationAddr(addr)

	status := p.getString(data[0]) == "+"

	roomStatus := record.RoomStatus{
		Station:            station,
		RoomNumber:         roomNumber,
		MessageLightStatus: status,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}
