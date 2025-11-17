package request

import (
	"time"

	vendor "github.com/weareplanet/ifcv5-drivers/xmlpos/record"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines/errors"
	defines "github.com/weareplanet/ifcv5-main/ifc/defines/errors"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-drivers/xmlpos/record"

	"github.com/pkg/errors"
)

func (p *Plugin) handlePostInquiry(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	req, ok := packet.Context.(*vendor.PostInquiry)

	if !ok {
		return errors.New("wrong context")
	}

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	hotelID := 1 // 1stellig

	now := time.Now()
	ts := p.getTime(req.Date, req.Time)

	requestType := p.getRequestType(req.RequestType)

	inquiry := record.PostingInquiry{
		Time:           ts,
		Station:        station,
		RequestType:    p.toString(requestType),
		SequenceNumber: p.toString(req.SequenceNumber),
		Inquiry:        req.InquiryInformation,
		MaximumMatches: req.MaximumReturnedMatches,
		WorkStation:    req.WorkstationId,
		PaymentMethod:  req.PaymentMethod,
		UserID:         req.WaiterId,
	}

	reply, err := p.PmsRequest(addr, inquiry, packet)

	var guests []*record.Guest

	switch reply.(type) {

	case *record.Guest:
		guest := reply.(*record.Guest)
		guests = append(guests, guest)

	case []*record.Guest:
		guests = reply.([]*record.Guest)

	}

	// handle guest list, send PostList

	if err == nil && len(guests) > 0 {

		ret := &vendor.PostList{
			Date:           now,
			Time:           now,
			HotelId:        p.toString(hotelID),
			SequenceNumber: req.SequenceNumber,
			PaymentMethod:  req.PaymentMethod,
			RevenueCenter:  req.RevenueCenter,
			WaiterId:       req.WaiterId,
			WorkstationId:  req.WorkstationId,
		}

		for i := range guests {

			guest := guests[i]

			item := vendor.PostListItem{
				RoomNumber:    guest.Reservation.RoomNumber,
				ReservationId: p.toInt(guest.Reservation.ReservationID),
				ProfileId:     p.toInt(guest.ProfileID),
				LastName:      guest.LastName,
				FirstName:     guest.FirstName,
				Title:         guest.Title,
				Vip:           guest.VIPStatus,
				PaymentMethod: req.PaymentMethod,
				CreditLimit:   guest.Reservation.CreditLimit,
				HotelId:       hotelID,
			}

			if guest.Rights.NoPostingInd {
				item.NoPost = "1"
			} else {
				item.NoPost = "0"
			}

			ret.PostListItem = append(ret.PostListItem, item)

		}

		err = p.send(addr, action, "", ret)
		return err
	}

	// handle error response, send negativ PostAnswer

	if err == nil && len(guests) == 0 {
		err = defines.ErrReservationNotFound
	}

	ret := &vendor.PostAnswer{
		Date:           now,
		Time:           now,
		HotelId:        hotelID,
		SequenceNumber: req.SequenceNumber,
		PaymentMethod:  req.PaymentMethod,
		RevenueCenter:  req.RevenueCenter,
		WaiterId:       req.WaiterId,
		WorkstationId:  req.WorkstationId,
	}

	if err != nil {

		switch err {

		case defines.ErrCreditLimitExceeded:
			ret.AnswerStatus = "CO"

		case defines.ErrNoPostSet:
			ret.AnswerStatus = "NP"

		case defines.ErrReservationNotFound:
			ret.AnswerStatus = "NG"

		default:
			ret.AnswerStatus = "UR"

		}

		ret.ResponseText = err.Error()
	}

	err = p.send(addr, action, "", ret)

	return err
}
