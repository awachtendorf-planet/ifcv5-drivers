package request

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleVoiceMessage(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = [folio number][$]

	if len(roomNumber) == 0 || len(data) == 0 || len(data[0]) == 0 {
		return
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()

	station, _ := dispatcher.GetStationAddr(addr)

	if strings.HasSuffix(data[0], "+") || strings.HasSuffix(data[0], "-") {

		status := "N"

		if strings.HasSuffix(data[0], "+") {
			status = "Y"
		}

		roomStatus := record.RoomStatus{
			Station:    station,
			RoomNumber: roomNumber,
			Voicemail:  status,
		}

		dispatcher.CreatePmsJob(addr, packet, roomStatus)

	} else {

		// todo: send error packet, folio number does not exist (no pms function available for this request)
		// QINVALID Folio Number data[0]

	}

}
