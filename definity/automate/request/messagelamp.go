package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-drivers/definity/template"

	"github.com/pkg/errors"
)

func (p *Plugin) handleMessageLamp(addr string, _ *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	room := p.getRoom(packet)

	if len(room) == 0 {
		return errors.New("empty room")
	}

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: room,
	}

	switch packet.Name {

	case template.PacketMessageLampOn:

		roomStatus.MessageLightStatus = true

	case template.PacketMessageLampOff:

		// false make no sense because of xsd omit empty definition. good luck with this.
		roomStatus.MessageLightStatus = false

	default:
		return nil
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

	return nil
}
