package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleVoiceCount(addr string, packet *ifc.LogicalPacket) {

	dispatcher := p.GetDispatcher()
	name := p.GetName()

	extension := p.getExtension(packet)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	unread := p.getValueAsNumeric(packet, "Unread")
	//read := p.getValueAsNumeric(packet, "Read")

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		Voicemail:  cast.ToString(unread),
	}

	dispatcher.CreatePmsJob(addr, packet, roomStatus)

}
