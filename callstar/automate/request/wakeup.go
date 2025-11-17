package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleWakeupSet(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = ?
	// 1 = hhmm

	if len(roomNumber) == 0 || len(data) < 2 {
		return
	}

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	timestamp, err := p.getTime(data[1])
	if err != nil {
		log.Warn().Msgf("%s addr '%s' construct time failed, err=%s", name, addr, err)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	wake := record.WakeupRequest{
		Station:  station,
		RoomName: roomNumber,
		Time:     timestamp,
	}

	dispatcher.CreatePmsJob(addr, packet, wake)

}
