package requestsimple

import (
	"fmt"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleWakeupSet(addr string, packet *ifc.LogicalPacket, retry bool) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	wakeTime, err := p.getTime(packet)
	if err != nil {
		return unknownError, err
	}

	if retry {
		if err, exist := p.lastError[addr]; exist {
			return 0, err
		}
		return 0, nil
	}

	station, _ := dispatcher.GetStationAddr(addr)

	wakeup := record.WakeupRequest{
		Station:  station,
		RoomName: roomNumber,
		Time:     wakeTime,
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	_, pmsErr, sendErr := automate.PmsRequest(station, wakeup, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = sendErr
	}

	if pmsErr != nil {
		p.lastError[addr] = pmsErr
	} else {
		delete(p.lastError, addr)
	}

	return 0, pmsErr

}

func (p *Plugin) handleWakeupClear(addr string, packet *ifc.LogicalPacket, retry bool) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	wakeTime, err := p.getTime(packet)
	if err != nil {
		return unknownError, err
	}

	if retry {
		if err, exist := p.lastError[addr]; exist {
			return 0, err
		}
		return 0, nil
	}

	station, _ := dispatcher.GetStationAddr(addr)

	wakeup := record.WakeupClear{
		Station:  station,
		RoomName: roomNumber,
		Time:     wakeTime,
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	_, pmsErr, sendErr := automate.PmsRequest(station, wakeup, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = sendErr
	}

	if pmsErr != nil {
		p.lastError[addr] = pmsErr
	} else {
		delete(p.lastError, addr)
	}

	return 0, pmsErr

}

func (p *Plugin) handleWakeupResult(addr string, packet *ifc.LogicalPacket, retry bool) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	wakeTime, err := p.getTime(packet)
	if err != nil {
		return unknownError, err
	}

	if retry {
		return 0, nil
	}

	station, _ := dispatcher.GetStationAddr(addr)

	wakeup := record.WakeupPerformed{
		Station:  station,
		RoomName: roomNumber,
		Time:     wakeTime,
	}

	statusCode := p.getField(packet, "StatusCode", true)

	switch statusCode {

	case "D": // success
		wakeup.Success = true
		wakeup.ResultCode = 1
		wakeup.Message = "wakeup call was ssuccessfully received by guest"

	case "T": // timeout
		wakeup.ResultCode = 2
		wakeup.Message = "wakeup call was not successfully received by guest, timeout error"

	case "H": // hardware error
		wakeup.ResultCode = 6
		wakeup.Message = "wakeup call was not successfully received by guest, hardware error"

	case "U": // tripleguest
		wakeup.ResultCode = 2
		wakeup.Message = "call undeliverable"

	default:
		wakeup.ResultCode = 2
		wakeup.Message = fmt.Sprintf("vendor response with unknown status code '%s'", statusCode)

	}

	dispatcher.CreatePmsJob(addr, packet, wakeup)

	return 0, nil

}
