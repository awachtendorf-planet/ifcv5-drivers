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

func (p *Plugin) handlePostRequest(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	req, ok := packet.Context.(*vendor.PostRequest)

	if !ok {
		return errors.New("wrong context")
	}

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	hotelID := 1 // 1stellig

	now := time.Now()
	ts := p.getTime(req.Date, req.Time)

	charge := p.getCharge(req)

	requestType := p.getRequestType(req.RequestType)

	posting := record.PostingRequest{
		PostingType:          1,
		Time:                 ts,
		Station:              station,
		Context:              charge,
		TotalAmount:          float64(req.TotalAmount),
		ReservationNumber:    uint64(req.ReservationId),
		SequenceNumber:       p.toString(req.SequenceNumber),
		CheckNumber:          p.toString(req.CheckNumber),
		ServingTime:          p.toString(req.ServingTime),
		RequestType:          p.toString(requestType),
		PaymentMethod:        req.PaymentMethod,
		UserID:               req.WaiterId,
		WorkStation:          req.WorkstationId,
		Covers:               req.Covers,
		Inquiry:              req.InquiryInformation,
		RoomName:             req.RoomNumber,
		GuestName:            req.LastName,
		MatchFromPostingList: req.MatchfromPostList,
	}

	if req.CreditlimitOverride == "Y" {
		posting.CreditlimitOverride = true
	}

	reply, err := p.PmsRequest(addr, posting, packet)

	ret := &vendor.PostAnswer{
		Date:           now,
		Time:           now,
		HotelId:        hotelID,
		SequenceNumber: req.SequenceNumber,
		PaymentMethod:  req.PaymentMethod,
		RevenueCenter:  req.RevenueCenter,
		WaiterId:       req.WaiterId,
		WorkstationId:  req.WorkstationId,
		CheckNumber:    req.CheckNumber,
		LastName:       req.LastName,
		ReservationId:  req.ReservationId,
		RoomNumber:     req.RoomNumber,
	}

	if err == nil {

		if response, ok := reply.(record.PostingResponse); ok {

			switch response.ResponseCode {

			case 1: // posted successfully
				ret.AnswerStatus = "OK"

			case 2: // denied credit limit
				ret.AnswerStatus = "CO"

			case 3: // denied nopost
				ret.AnswerStatus = "NP"

			case 4: // reservation not found
				ret.AnswerStatus = "NG"

			case 5: // denied not defined
				ret.AnswerStatus = "UR"

			}

			ret.ResponseText = response.ResponseText
			ret.CheckNumber = p.toInt(response.CheckNumber)
			ret.ReservationId = int(response.ReservationNumber)
			ret.LastName = response.GuestName
			ret.RoomNumber = response.RoomName

		} else {

			ret.AnswerStatus = "UR"
			ret.ResponseText = "unknown answer from PMS"

		}

	} else {

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

func (p *Plugin) getCharge(req *vendor.PostRequest) record.ChargePosting {

	charge := record.ChargePosting{}

	add := func(v *[]record.ChargeItem, items ...int) {

		for i := range items {

			if value := items[i]; value != 0 {
				item := record.ChargeItem{ID: i + 1, Value: float64(value)}
				*v = append(*v, item)
			}

		}

	}

	add(&charge.SubTotals,
		req.Subtotal1, req.Subtotal2, req.Subtotal3, req.Subtotal4, req.Subtotal5, req.Subtotal6, req.Subtotal7, req.Subtotal8, req.Subtotal9,
		req.Subtotal10, req.Subtotal11, req.Subtotal12, req.Subtotal13, req.Subtotal14, req.Subtotal15, req.Subtotal16,
	)

	add(&charge.Discounts,
		req.Discount1, req.Discount2, req.Discount3, req.Discount4, req.Discount5, req.Discount6, req.Discount7, req.Discount8, req.Discount9,
		req.Discount10, req.Discount11, req.Discount12, req.Discount13, req.Discount14, req.Discount15, req.Discount16,
	)

	add(&charge.Taxes,
		req.Tax1, req.Tax2, req.Tax3, req.Tax4, req.Tax5, req.Tax6, req.Tax7, req.Tax8, req.Tax9,
		req.Tax10, req.Tax11, req.Tax12, req.Tax13, req.Tax14, req.Tax15, req.Tax16,
	)

	add(&charge.ServiceCharges,
		req.ServiceCharge1, req.ServiceCharge2, req.ServiceCharge3, req.ServiceCharge4, req.ServiceCharge5, req.ServiceCharge6, req.ServiceCharge7, req.ServiceCharge8, req.ServiceCharge9,
		req.ServiceCharge10, req.ServiceCharge11, req.ServiceCharge12, req.ServiceCharge13, req.ServiceCharge14, req.ServiceCharge15, req.ServiceCharge16,
	)

	add(&charge.Tips,
		req.Tip,
	)

	return charge
}
