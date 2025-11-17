package request

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/plabor/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

const (
	packetMsgItem = "msg item"
	packetMsgEnd  = "msg end"
)

func (p *Plugin) handleReqMessage(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	automate := p.automate

	dispatcher := automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	roomName := p.getField(packet, "RoomNumber")
	if p.driver.LeadingZeros(station) {
		roomName = strings.TrimLeft(roomName, "0")
	}

	reservationID := p.getField(packet, "AccountID")

	request := record.GuestMessageRequest{
		ReservationID: reservationID,
		RoomNumber:    roomName,
		Station:       station,
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	reply, pmsErr, sendErr := automate.PmsRequest(station, request, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}

	if pmsErr != nil {
		return pmsErr
	}

	p.clearRecords(addr, packet.Name)

	messages := []record.GuestMessage{}

	switch reply.(type) {

	case record.GuestMessage:
		msg := reply.(record.GuestMessage)
		messages = append(messages, msg)

	case []record.GuestMessage:
		messages = reply.([]record.GuestMessage)

	}

	if len(messages) == 0 { // pms error or no item exist
		item := record.GuestMessage{
			ReservationID: request.ReservationID,
			RoomNumber:    request.RoomNumber,
		}

		err := p.send(addr, action, template.PacketMessageEnd, trackingId, item)
		p.send(addr, action, template.PacketEot, trackingId, nil)
		return err

	}

	for i := range messages { // queue bill items -> send with nextaction
		item := messages[i]

		if len(item.ReservationID) == 0 {
			item.ReservationID = request.ReservationID
		}
		if len(item.RoomNumber) == 0 {
			item.RoomNumber = request.RoomNumber
		}

		var answer *record.Generic
		answer.Set("Record", item)

		if i < len(messages)-1 && p.driver.Protocol(station) == 0 {

			answer.Set(packetType, packetMsgItem)
			p.addNextRecord(addr, template.PacketMessageBlock, answer)
		} else {

			answer.Set(packetType, packetMsgEnd)
			p.addNextRecord(addr, template.PacketMessageEnd, answer)
		}
	}
	p.addNextRecord(addr, template.PacketEot, nil)

	return nil
}

func (p *Plugin) sendMessagePart(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketMessageBlock, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendMessageEnd(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketMessageEnd, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) handleMessagePart(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) (bool, error) {

	// bill item's/balance response to template.PacketGuestBillRequest

	answer := p.getNextRecord(addr, packet.Name)
	if answer == nil {
		return false, nil
	}

	var err error
	next := true

	switch answer.(type) {

	case *record.Generic:

		record := answer.(*record.Generic)
		recordType, _ := record.Get(packetType)
		item, _ := record.Get("Record")

		switch recordType {

		case packetMsgItem:
			err = p.sendMessagePart(addr, action, "", item)

		case packetMsgEnd:
			err = p.sendMessageEnd(addr, action, "", item)

		}

	}

	if err == nil {
		next = p.dropRecord(addr, packet.Name)
	}

	return next, err
}
