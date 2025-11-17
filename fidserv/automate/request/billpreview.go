package request

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

const (
	packetType    = "type"
	packetItem    = "item"
	packetBalance = "balance"
)

func (p *Plugin) handleBillPreview(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) (bool, error) {

	// template.PacketGuestBillRequest

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	dispatcher := automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	preview := record.BillPreviewRequest{
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &preview); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	reply, pmsErr, sendErr := automate.PmsRequest(station, preview, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}

	if dispatcher.IsShutdown(pmsErr) {
		return false, pmsErr
	}

	p.clearRecords(addr, packet.Name)

	invoice := record.BillPreviewInvoice{}

	if pmsErr == nil {

		switch reply.(type) {

		case record.BillPreviewInvoice:
			invoice = reply.(record.BillPreviewInvoice)
		}
	}

	if len(invoice.Items) == 0 { // pms error or no item exist
		item := record.BillPreviewBalance{
			ReservationID: preview.ReservationID,
			RoomNumber:    preview.RoomNumber,
			Amount:        0,
		}
		answer, err := driver.MarshalPacket(item)
		if err != nil {
			log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, item, err)
		}

		switch invoice.ResponseCode {

		case 1:
			answer.Set("AS", "OK")

		case 2:
			answer.Set("AS", "NG")

		default:
			answer.Set("AS", "UR")
		}

		err = p.sendBillBalance(addr, action, correlationId, answer)
		return false, err

	}

	for i := range invoice.Items { // queue bill items -> send with nextaction
		item := invoice.Items[i]

		if len(item.ReservationID) == 0 {
			item.ReservationID = invoice.Balance.ReservationID
		}
		if len(item.RoomNumber) == 0 {
			item.RoomNumber = invoice.Balance.RoomNumber
		}

		answer, err := driver.MarshalPacket(item)
		if err != nil {
			log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, item, err)
		}

		if !item.Time.IsZero() {
			if _, exist := answer.Get("DA"); !exist {
				answer.Set("DA", item.Time)
			}
			if _, exist := answer.Get("TI"); !exist {
				answer.Set("TI", item.Time)
			}
		}

		answer.Set(packetType, packetItem)
		p.addNextRecord(addr, packet.Name, answer)
	}

	// queue bill balance -> send with nextaction

	item := record.BillPreviewBalance{
		ReservationID: invoice.Balance.ReservationID,
		RoomNumber:    invoice.Balance.RoomNumber,
		Amount:        invoice.Balance.Amount,
	}

	answer, err := driver.MarshalPacket(item)
	if err != nil {
		log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, item, err)
	}
	answer.Set("AS", "OK")

	answer.Set(packetType, packetBalance)
	p.addNextRecord(addr, packet.Name, answer)

	return true, nil
}

func (p *Plugin) sendBillItem(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketGuestBillItem, tracking, "XI", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendBillBalance(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketGuestBillBalance, tracking, "XB", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) handleBillInvoice(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) (bool, error) {

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

		switch recordType {

		case packetItem:
			err = p.sendBillItem(addr, action, "", record)

		case packetBalance:
			err = p.sendBillBalance(addr, action, "", record)

		}

	}

	if err == nil {
		next = p.dropRecord(addr, packet.Name)
	}

	return next, err
}
