package request

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleWakeup(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketWakeupAnswer
	// template.PacketWakeupRequest
	// template.PacketWakeupClear

	driver := p.driver
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()
	station, _ := dispatcher.GetStationAddr(addr)

	wakeupTime, err := driver.ParseTime(packet)
	if err != nil {
		log.Error().Msgf("%s addr '%s' wakeup time failed, err=%s", name, addr, err)
		return nil
	}

	roomNumber, exist := driver.GetField(packet, "RN")
	if !exist {
		log.Error().Msgf("%s addr '%s' wakeup room number failed, err=%s", name, addr, "missing value")
		return nil
	}

	switch packet.Name {

	case template.PacketWakeupRequest:
		wake := record.WakeupRequest{
			Station:  station,
			Time:     wakeupTime,
			RoomName: roomNumber,
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	case template.PacketWakeupClear:
		wake := record.WakeupClear{
			Station:  station,
			Time:     wakeupTime,
			RoomName: roomNumber,
		}
		dispatcher.CreatePmsJob(addr, packet, wake)

	case template.PacketWakeupAnswer:
		answerStatus, _ := driver.GetField(packet, "AS")
		success := (answerStatus == "OK" || answerStatus == "DE" || answerStatus == "SV")
		answer := driver.GetAnswerText(answerStatus)

		wake := record.WakeupPerformed{
			Station:  station,
			Success:  success,
			Message:  answer,
			Time:     wakeupTime,
			RoomName: roomNumber,
		}

		switch answerStatus {

		case "OK", "SV": // ok
			wake.ResultCode = 1

		case "BY": // busy
			wake.ResultCode = 4

		case "NR": // no response
			wake.ResultCode = 2

		case "RY": // retry
			wake.ResultCode = 3

		case "UR": // unprocessable
			wake.ResultCode = 6

		case "DE": // deleted
			wake.ResultCode = 5

		}

		dispatcher.CreatePmsJob(addr, packet, wake)
	}

	return nil
}
