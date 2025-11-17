package request

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/plabor/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

const (
	packetType    = "type"
	packetItem    = "item"
	packetBalance = "balance"
)

func (p *Plugin) handleReqBill(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketGuestBillRequest
	automate := p.automate

	dispatcher := automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	roomName := p.getField(packet, "RoomNumber")
	if p.driver.LeadingZeros(station) {
		roomName = strings.TrimLeft(roomName, "0")
	}

	reservationID := p.getField(packet, "AccountID")

	preview := record.BillPreviewRequest{
		ReservationID: reservationID,
		RoomNumber:    roomName,
		Station:       station,
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	reply, pmsErr, sendErr := automate.PmsRequest(station, preview, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}

	if pmsErr != nil {
		return pmsErr
	}

	p.clearRecords(addr, packet.Name)

	invoice := record.BillPreviewInvoice{}

	switch reply.(type) {

	case record.BillPreviewInvoice:
		invoice = reply.(record.BillPreviewInvoice)
	}

	if len(invoice.Items) == 0 { // pms error or no item exist

		err := p.send(addr, action, template.PacketBalance, trackingId, invoice.Balance)
		p.send(addr, action, template.PacketEot, trackingId, nil)
		return err

	}

	for i := range invoice.Items { // queue bill items -> send with nextaction
		item := invoice.Items[i]

		if len(item.ReservationID) == 0 {
			item.ReservationID = invoice.Balance.ReservationID
		}
		if len(item.RoomNumber) == 0 {
			item.RoomNumber = invoice.Balance.RoomNumber
		}

		p.addNextRecord(addr, template.PacketBillPart, item)
	}

	// queue bill balance -> send with nextaction

	p.addNextRecord(addr, template.PacketBalance, invoice.Balance)

	p.addNextRecord(addr, template.PacketEot, nil)

	return nil
}

func (p *Plugin) sendBillItem(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketBillPart, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendBillBalance(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, template.PacketBalance, tracking, context)
	if err != nil {
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
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

	case record.BillPreviewBalance:

		err = p.sendBillItem(addr, action, "", answer)

	case record.BillPreviewItem:
		err = p.sendBillBalance(addr, action, "", answer)

	}

	if err == nil {
		next = p.dropRecord(addr, packet.Name)
	}

	return next, err
}
