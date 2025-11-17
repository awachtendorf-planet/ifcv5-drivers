package requestitems

import (
	"time"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
)

func (p *Plugin) handleDisplayRequest(addr string, packet *ifc.LogicalPacket) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	accountNumber, err := p.getAccountNumber(addr, packet)
	roomNumber, err := p.getRoomNumber(addr, packet)

	if len(accountNumber) == 0 {
		return unknownAccount, err
	}
	if len(roomNumber) == 0 {
		return unknownRoom, err
	}

	station, _ := dispatcher.GetStationAddr(addr)

	// inquiry request

	inquiry := record.PostingInquiry{
		Time:           time.Now(),
		Station:        station,
		Inquiry:        accountNumber,
		MaximumMatches: 1,
		RequestType:    "8",
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)
	packet.Tracking = correlationId

	reply, pmsErr, sendErr := automate.PmsRequest(station, inquiry, pmsTimeOut, correlationId, trackingId)

	if sendErr != nil {
		pmsErr = errPmsUnavailable
	}

	if pmsErr != nil {
		return 0, pmsErr
	}

	accountNumberLength := p.driver.GetAccountNumberLength(station)

	var guests []*record.Guest

	switch reply.(type) {

	case *record.Guest:

		guest := reply.(*record.Guest)
		if len(guest.Reservation.RoomNumber) <= 6 && len(guest.Reservation.ReservationID) <= accountNumberLength {
			guests = append(guests, guest)
		} else {
			log.Warn().Msgf("%s addr '%s' ignore reservation-id: %s because of roomnumber or reservation-id is too long", name, addr, guest.Reservation.ReservationID)
		}

	case []*record.Guest:

		guest := reply.([]*record.Guest)
		for i := range guest {
			if len(guest[i].Reservation.RoomNumber) <= 6 && len(guest[i].Reservation.ReservationID) <= accountNumberLength {
				guests = append(guests, guest[i])
			} else {
				log.Warn().Msgf("%s addr '%s' ignore reservation-id: %s because of roomnumber or reservation-id is too long", name, addr, guest[i].Reservation.ReservationID)
			}
		}

	}

	if len(guests) == 0 {
		return unknownAccount, errPmsEmptyResponse
	}

	guest := guests[0]

	// folio request

	// todo: ?
	// guest.Reservation.RoomNumber == roomNumber
	// guest.Reservation.ReservationID == accountNumber

	preview := record.BillPreviewRequest{
		Station:       station,
		RoomNumber:    guest.Reservation.RoomNumber,
		ReservationID: guest.Reservation.ReservationID,
	}

	correlationId = dispatcher.NewCorrelationID(station)

	reply, pmsErr, sendErr = automate.PmsRequest(station, preview, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = errPmsUnavailable
	}

	if pmsErr != nil {
		return 0, pmsErr
	}

	invoice := record.BillPreviewInvoice{}

	switch reply.(type) {

	case record.BillPreviewInvoice:
		invoice = reply.(record.BillPreviewInvoice)
	}

	switch invoice.ResponseCode {

	case 1: // ok
		break

	case 2: // reservation or room not found
		return unknownAccount, errors.New("reservation or room not found")

	case 3: // feature not enabled
		return lockedFolio, errors.New("feature not enabled")

	}

	identifier := p.getTransactionIdentifier(packet)
	sequence := int(0)

	add := func(name string, context interface{}, last bool) {
		if last {
			sequence = 9999
		}
		p.addNextRecord(addr, packet.Name, identifier, &sequence, name, context, correlationId, trackingId)
	}

	// create NAME record
	add(template.PacketNameReply, guest, false)

	// create ITEM records
	for i := range invoice.Items {
		add(template.PacketItemReply, invoice.Items[i], false)
	}

	// create BAL record
	add(template.PacketBalanceReply, invoice.Balance, true)

	return 0, nil
}
