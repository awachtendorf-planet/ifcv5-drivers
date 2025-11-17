package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getExtension(packet)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	unread := p.getNumeric(packet, "Unread")
	voiceMail := cast.ToString(unread)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		Voicemail:  voiceMail,
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}
