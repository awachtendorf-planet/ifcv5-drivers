package request

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleSyncRequest(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getNumeric(packet, "Extension")

	if extension == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	roomNumber := cast.ToString(extension)
	sequence := p.getString(packet, "Ticket")

	// date := p.getString(packet, "Date")
	// time := p.getString(packet, "Time")
	// fmt.Println(date, time)

	station, _ := dispatcher.GetStationAddr(addr)

	inquiry := record.PostingInquiry{
		Time:           time.Now(),
		Station:        station,
		Inquiry:        roomNumber,
		SequenceNumber: sequence,
		MaximumMatches: 1,
		RequestType:    "2",
	}
	reply, pmsErr, sendErr := automate.PmsRequest(station, inquiry, pmsTimeOut, "", "")

	if sendErr != nil {
		pmsErr = errPMSUnavailable
	}

	if pmsErr != nil {
		log.Error().Msgf("%s addr '%s' request failed, err=%s", name, addr, pmsErr)
		return
	}

	var guests []*record.Guest

	switch guest := reply.(type) {

	case *record.Guest:

		if guest.Reservation.RoomNumber == roomNumber {
			guests = append(guests, guest)
		} else {
			log.Warn().Msgf("%s addr '%s' ignore reservation-id: %s because of roomnumber '%s' does not match inquiry '%s'", name, addr, guest.Reservation.ReservationID, guest.Reservation.RoomNumber, roomNumber)
		}

	case []*record.Guest:

		for i := range guest {
			if guest[i].Reservation.RoomNumber == roomNumber {
				guests = append(guests, guest[i])
			} else {
				log.Warn().Msgf("%s addr '%s' ignore reservation-id: %s because of roomnumber '%s' does not match inquiry '%s'", name, addr, guest[i].Reservation.ReservationID, guest[i].Reservation.RoomNumber, roomNumber)
			}
		}
	}

	// create room checkout if no guest found
	if len(guests) == 0 {
		guest := &record.Guest{Reservation: record.Reservation{RoomNumber: roomNumber}}
		dispatcher.CreateDriverJob(station, order.Checkout, guest, "")
		return
	}

	// create guest checkins
	for i := range guests {
		guest := guests[i]
		dispatcher.CreateDriverJob(station, order.Checkin, guest, "")
	}

}
