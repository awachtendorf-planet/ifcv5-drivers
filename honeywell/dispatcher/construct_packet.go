package honeywell

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/honeywell/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	switch packetName {

	case template.PacketCheckIn, template.PacketCheckOut:

		room := d.GetRoom(context)

		if err := d.CheckRoom(room); err != nil {
			return packet, err
		}

		var data string

		station, _ := d.GetStationAddr(addr)
		protocol := d.GetProtocolType(station)

		switch protocol {

		case HONEYWELL_PROTOCOL:
			data = pad.Right(room, 4, " ")

		case ALERTON_PROTOCOL_1, ALERTON_PROTOCOL_2:
			data = pad.Left(room, 4, "0")

		default:
			return packet, errors.Errorf("unknown protocol type: %d ", protocol)

		}

		packet.Add("Room", []byte(data))

	default:
		return packet, errors.Errorf("packet '%s' handler not defined", packetName)

	}

	return packet, nil
}

func (d *Dispatcher) GetRoom(context interface{}) string {

	if guest, ok := context.(*record.Guest); ok {
		// return strings.TrimLeft(guest.Reservation.RoomNumber, "0") // wurde mit commit 82be86305bae8756cbdbc36545d51bf45baa9391 entfernt
		return guest.Reservation.RoomNumber
	}

	return ""

}

func (d *Dispatcher) CheckRoom(room string) error {

	if len(room) == 0 {
		return errors.New("empty room")
	}

	if len(room) > 4 {
		return errors.Errorf("invalid room length '%s' (max length 4)", room)
	}

	roomNumber := cast.ToInt(strings.TrimLeft(room, "0"))

	if roomNumber < 1 || roomNumber > 9999 {
		return errors.Errorf("invalid room '%s' (must be 1...9999)", room)
	}

	return nil

}

func (d *Dispatcher) IsSharer(context interface{}) bool {

	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}

	return false

}
