package requestitems

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleGuestMessageRequest(addr string, packet *ifc.LogicalPacket) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()
	//driver := p.driver

	accountNumber, err := p.getAccountNumber(addr, packet)
	if err != nil {
		return unknownAccount, err
	}

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	station, _ := dispatcher.GetStationAddr(addr)

	message := record.GuestMessageRequest{
		Station:       station,
		MarkAsRead:    false,
		RoomNumber:    roomNumber,
		ReservationID: accountNumber,
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)
	packet.Tracking = correlationId

	reply, pmsErr, sendErr := automate.PmsRequest(station, message, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = errPmsUnavailable
	}

	if pmsErr != nil {
		return 0, pmsErr
	}

	accountNumberLength := p.driver.GetAccountNumberLength(station)

	messages := []record.GuestMessage{}

	switch msg := reply.(type) {

	case record.GuestMessage:

		if len(msg.RoomNumber) <= 6 && len(msg.ReservationID) <= accountNumberLength && len(msg.MessageID) <= 6 {
			messages = append(messages, msg)
		} else {
			log.Warn().Msgf("%s addr '%s' ignore guest message-id '%s' because of roomnumber or reservation-id or message-id is too long", name, addr, msg.MessageID)
		}

	case []record.GuestMessage:

		for i := range msg {
			if len(msg[i].RoomNumber) <= 6 && len(msg[i].ReservationID) <= accountNumberLength && len(msg[i].MessageID) <= 6 {
				messages = append(messages, msg[i])
			} else {
				log.Warn().Msgf("%s addr '%s' ignore guest message-id '%s' because of roomnumber or reservation-id or message-id is too long", name, addr, msg[i].MessageID)
			}
		}
	}

	if len(messages) == 0 {
		return emptyGuestMessage, errPmsEmptyResponse
	}

	identifier := p.getTransactionIdentifier(packet)
	sequence := int(0)

	add := func(name string, context interface{}, last bool) {
		if last {
			sequence = 9999
		}
		p.addNextRecord(addr, packet.Name, identifier, &sequence, name, context, correlationId, trackingId)
	}

	create := func(context record.GuestMessage, last bool) {
		add(template.PacketGuestMessageHeader, context, false)
		add(template.PacketGuestMessageCaller, context, false)
		lastTextMessage, textMessages := splitMessages(context)
		for _, textMessage := range textMessages {
			add(template.PacketGuestMessageText, textMessage, false)
		}
		add(template.PacketGuestMessageText, lastTextMessage, last)
	}

	records := len(messages)
	for i := 0; i < records; i++ {
		last := (i == records-1)
		create(messages[i], last)
	}

	return 0, nil
}

func split(input string, length int) (string, string) {
	message := input[:length]
	spaceIndex := strings.LastIndex(message, " ")
	if spaceIndex == -1 {
		spaceIndex = length
	}
	message = input[:spaceIndex]
	rest := input[spaceIndex:]
	return message, rest
}

func splitMessages(input record.GuestMessage) (record.GuestMessage, []record.GuestMessage) {
	var guestMessages []record.GuestMessage
	text := input.Text

	length := 64

	for len(text) > length {
		part := ""
		part, text = split(text, length)
		guestMessage := record.GuestMessage{
			Station:       input.Station,
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			DisplayName:   input.DisplayName,
			UserID:        input.UserID,
			Time:          input.Time,
			ReservationID: input.ReservationID,
			RoomNumber:    input.RoomNumber,
			MessageID:     input.MessageID,
			Text:          part,
		}
		if guestMessages == nil {
			guestMessages = []record.GuestMessage{}
		}
		guestMessages = append(guestMessages, guestMessage)

	}

	output := record.GuestMessage{
		Station:       input.Station,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		DisplayName:   input.DisplayName,
		UserID:        input.UserID,
		Time:          input.Time,
		ReservationID: input.ReservationID,
		RoomNumber:    input.RoomNumber,
		MessageID:     input.MessageID,
		Text:          text,
	}
	return output, guestMessages
}
