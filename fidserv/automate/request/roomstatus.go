package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketRoomData

	driver := p.driver
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	station, _ := dispatcher.GetStationAddr(addr)

	roomStatus := record.RoomStatus{
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &roomStatus); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

	return nil
}
