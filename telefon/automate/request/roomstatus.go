package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/telefon/template"

	"github.com/pkg/errors"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	// extension
	extension, _ := p.getString(packet, template.Extension)
	if len(extension) == 0 {
		log.Error().Msgf("%s addr '%s' extension not found", name, addr)
		return errors.New("extension not found")
	}

	station, _ := dispatcher.GetStationAddr(addr)

	// maid/roomstatus
	maid, _ := p.getString(packet, template.User)
	status, _ := p.getString(packet, template.RoomStatus)

	status = dispatcher.GetPMSMapping(packet.Name, station, "RS", status)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
		UserID:     maid,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

	return nil
}
