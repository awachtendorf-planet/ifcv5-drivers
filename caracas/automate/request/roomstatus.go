package request

import (
	"fmt"

	// "github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/spf13/cast"
)

func (p *Plugin) handleRoomStatus(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {

		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	status := p.getField(packet, "Roomstatus", true)
	userID := p.getField(packet, "UserID", true)

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		RoomStatus: status,
		UserID:     userID,
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

	voicemailStatus := p.getField(packet, "Status", true)
	voicemailAmount := cast.ToInt(p.getField(packet, "AmountNewMessages", true))
	voicemailBool := "N"

	if voicemailStatus == "1" || voicemailAmount > 0 {

		voicemailBool = "Y"

	}

	roomStatus := record.RoomStatus{
		Station:    station,
		RoomNumber: extension,
		Voicemail:  voicemailBool,
	}

	dispatcherObj.CreatePmsJob(addr, packet, roomStatus)

	return nil

}

// func (p *Plugin) handleDND(addr string, packet *ifc.LogicalPacket) error {

// 	dispatcherObj := p.GetDispatcher()
// 	name := p.GetName()

// 	extension := p.getField(packet, "Extension", true)

// 	if len(extension) == 0 {

// 		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
// 	}

// 	station, _ := dispatcherObj.GetStationAddr(addr)

// 	dnd := cast.ToBool(p.getField(packet, "DND", false))

// 	roomStatus := record.RoomStatus{
// 		Station:      station,
// 		RoomNumber:   extension,
// 		DoNotDisturb: dnd,
// 	}

// 	dispatcherObj.CreatePmsJob(addr, packet, roomStatus)

// 	return nil

// }
