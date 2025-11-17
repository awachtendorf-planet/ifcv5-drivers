package requestitems

import (
	"time"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	defines "github.com/weareplanet/ifcv5-main/ifc/defines/errors"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleLookupRequest(addr string, packet *ifc.LogicalPacket) (int, error) {

	// PacketLookupRequest
	// PacketStatusRequest

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	station, _ := dispatcher.GetStationAddr(addr)

	inquiry := record.PostingInquiry{
		Time:           time.Now(),
		Station:        station,
		Inquiry:        roomNumber,
		MaximumMatches: 6,
		RequestType:    "2",
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)
	packet.Tracking = correlationId

	reply, pmsErr, sendErr := automate.PmsRequest(station, inquiry, pmsTimeOut, correlationId, trackingId)

	if sendErr != nil {
		pmsErr = errPmsUnavailable // ERR 00
	}

	if pmsErr != nil {

		switch pmsErr {

		case defines.ErrCreditLimitExceeded, defines.ErrNoPostSet, defines.ErrReservationNotFound:
			return roomUnoccupied, errPmsRejectedRequest // ERR 03
		}

		return 0, pmsErr // ERR 00
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
		return roomUnoccupied, errPmsEmptyResponse
	}

	identifier := p.getTransactionIdentifier(packet)
	sequence := int(0)

	add := func(name string, context interface{}, last bool) {
		if last {
			sequence = 9999
		}
		p.addNextRecord(addr, packet.Name, identifier, &sequence, name, context, correlationId, trackingId)
	}

	// maximal 6 records
	records := len(guests)
	if records > 6 {
		records = 6
	}

	for i := 0; i < records; i++ {
		last := (i == records-1)
		switch packet.Name {

		case template.PacketLookupRequest:
			add(template.PacketNameReply, guests[i], last)

		case template.PacketStatusRequest:
			add(template.PacketInfoReply, guests[i], last)

		}
	}

	return 0, nil
}
