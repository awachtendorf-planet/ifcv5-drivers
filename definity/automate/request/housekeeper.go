package request

import (
	"fmt"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-drivers/definity/template"

	"github.com/pkg/errors"
)

func (p *Plugin) handleHouseKeeperStatus(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	room := p.getRoom(packet)

	if len(room) == 0 {
		return errors.New("empty room")
	}

	dig := p.getString(packet, "DIG")

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: room,
		UserID:     dig,
	}

	proc := p.getProc(packet)

	switch packet.Name {

	case template.PacketHouseKeeperRoom, template.PacketHouseKeeperStation:

		switch proc {

		case 0x1, // housekeeper in room
			0x2, // room clean - vacant
			0x3, // room clean - occupied
			0x4: // room not clean - vacant

		case 0x5, // room not clean - occupied
			0x6: // room clean - need inspection

			if packet.Name == template.PacketHouseKeeperRoom {
				break
			}

			fallthrough

		default:

			var rejected string

			if packet.Name == template.PacketHouseKeeperRoom {
				rejected = template.PacketHouseKeeperRoomRejected
			} else {
				rejected = template.PacketHouseKeeperStationRejected
			}

			p.send(addr, action, rejected, "", roomStatus)

			return errors.Errorf("unknown process code: %d", proc)

		}

	default:
		return nil

	}

	status := fmt.Sprintf("%d", proc)

	status = dispatcher.GetPMSMapping("Room Status", station, "RS", status)

	roomStatus.RoomStatus = status

	var accepted string

	if packet.Name == template.PacketHouseKeeperRoom {
		accepted = template.PacketHouseKeeperRoomAccepted
	} else {
		accepted = template.PacketHouseKeeperStationAccepted
	}

	p.send(addr, action, accepted, "", roomStatus)

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

	return nil
}
