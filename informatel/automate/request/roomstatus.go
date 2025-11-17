package request

import (
	"fmt"

	// "github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {

		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	status := p.getField(packet, "RoomStatus", true)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
	}

	dispatcherObj.CreatePmsJob(addr, packet, roomStatus)

	return nil

}

func (p *Plugin) handleVoicemail(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {

		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	voicemailStr := p.getField(packet, "VoiceMail", true)
	voicemailBool := "N"

	if voicemailStr == "P" {

		voicemailBool = "Y"

	} else if voicemailStr != "A" {

		return fmt.Errorf("%s addr '%s' non-permissible voicemail state '%s'", name, addr, voicemailStr)
	}

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		Voicemail:  voicemailBool,
	}

	dispatcherObj.CreatePmsJob(addr, packet, roomStatus)

	return nil

}
