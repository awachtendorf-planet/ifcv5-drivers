package request

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleGuestMessageRequest(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) (bool, error) {

	// template.PacketGuestMessageRequest

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	dispatcher := automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	message := record.GuestMessageRequest{
		Station:    station,
		MarkAsRead: false,
	}

	if err := driver.UnmarshalPacket(packet, &message); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	reply, pmsErr, sendErr := automate.PmsRequest(station, message, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}

	if dispatcher.IsShutdown(pmsErr) {
		return false, pmsErr
	}

	p.clearRecords(addr, packet.Name)

	messages := []record.GuestMessage{}

	if pmsErr == nil {

		switch reply.(type) {

		case record.GuestMessage:
			msg := reply.(record.GuestMessage)
			messages = append(messages, msg)

		case []record.GuestMessage:
			messages = reply.([]record.GuestMessage)

		}
	}

	if len(messages) == 0 { // pms error or no message exist
		item := record.GuestMessage{
			ReservationID: message.ReservationID,
			RoomNumber:    message.RoomNumber,
		}
		answer, err := driver.MarshalPacket(item)
		if err != nil {
			log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, item, err)
		}

		// without MI/MT signals that no message exist
		answer.Delete("MI")
		answer.Delete("MT")

		// p.addNextRecord(addr, packet.Name, answer) // single item -> send direct without nextaction
		err = p.sendGuestMessageText(addr, action, correlationId, answer)
		return false, err

	}

	for i := range messages { // queue message items -> send with nextaction
		item := messages[i]
		answer, err := driver.MarshalPacket(item)
		if err != nil {
			log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, item, err)
		}

		p.addNextRecord(addr, packet.Name, answer)
	}

	return true, nil
}

func (p *Plugin) handleGuestMessageItem(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) (bool, error) {

	// message item's response to template.PacketGuestMessageRequest

	answer := p.getNextRecord(addr, packet.Name)
	if answer == nil {
		return false, nil
	}

	next := true

	err := p.sendGuestMessageText(addr, action, "", answer)
	if err == nil {
		next = p.dropRecord(addr, packet.Name)
	}

	return next, err
}

func (p *Plugin) handleGuestMessageDelete(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketGuestMessageDelete:

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	dispatcher := automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	message := record.GuestMessageDelete{
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &message); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	dispatcher.CreatePmsJob(addr, packet, message)
	return nil
}

func (p *Plugin) sendGuestMessageText(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketGuestMessageText, tracking, "XT", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
