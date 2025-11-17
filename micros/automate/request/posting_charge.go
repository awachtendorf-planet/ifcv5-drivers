package request

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/weareplanet/ifcv5-drivers/micros/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handlePostingCharge(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketPostingRequest
	// PI = search inquiry

	dispatcherObj := p.GetDispatcher()

	station, _ := dispatcherObj.GetStationAddr(addr)

	workstation := p.getField(packet, "WorkStation", false)

	accountID := p.getField(packet, "AccountID", false)

	accountID = strings.TrimRight(accountID, " ")

	checknumber := p.getField(packet, "GuestCheckNumber", true)
	selectionInfo := p.getField(packet, "SelectionInfo", true)
	selectionNumber := cast.ToInt(p.getField(packet, "SelectionNumber", true))
	numberOfCovers := cast.ToInt(p.getField(packet, "NumberOfCovers", true))
	servingPeriod := p.getField(packet, "ServingPeriodNumber", true)
	discount := p.getField(packet, "DiscountTotal", true)
	outletNumber := cast.ToInt(p.getField(packet, "RevenueCenterNumber", true))
	paymentType := p.getField(packet, "CurrentPaymentNumber", true)
	paymentTypeMapped := p.driver.GetPaymentMethodMapping(station, paymentType)
	totalAmount := cast.ToFloat64(p.getField(packet, "CurrentPaymentAmount", true))

	waiter := p.getField(packet, "TransEmployeeNumber", true)
	cashier := p.getField(packet, "ChkEmployeeNumber", true)

	tip := cast.ToFloat64(p.getField(packet, "ServiceChargeTotalEntered", true))
	entertainment := cast.ToFloat64(p.getField(packet, "ServiceChargeTotalAutomatic", true))

	charge := record.ChargePosting{}

	charge.Taxes = p.getTotals(packet, "Tax", 8)

	charge.SubTotals = p.getTotals(packet, "Sales", 16)

	charge.ServiceCharges = []record.ChargeItem{
		record.ChargeItem{ID: 1, Value: entertainment},
	}

	charge.Tips = []record.ChargeItem{
		record.ChargeItem{ID: 1, Value: tip},
	}

	discounts := []record.ChargeItem{}

	discounts = append(discounts, record.ChargeItem{
		ID:    1,
		Value: cast.ToFloat64(discount),
	})

	charge.Discounts = discounts

	var userID string

	if len(waiter) > 0 {

		userID = waiter
	} else if len(cashier) > 0 {

		userID = cashier
	}

	posting := record.PostingRequest{
		PostingType:          1,
		Station:              station,
		Context:              charge,
		TotalAmount:          totalAmount,
		CheckNumber:          checknumber,
		ServingTime:          servingPeriod,
		Covers:               numberOfCovers,
		OutletNumber:         outletNumber,
		MatchFromPostingList: selectionNumber,
		PaymentMethod:        paymentTypeMapped,
		Inquiry:              accountID,
		UserID:               userID,
		GuestName:            selectionInfo, // split guestname and roomnumber
		WorkStation:          workstation,
	}

	response := ifc.NewLogicalPacket(template.PacketOutletChargeResponse, addr, packet.Tracking)

	response.Add("WorkStation", []byte(workstation))

	reply, err := p.PmsRequest(addr, posting, packet)

	if err != nil {
		response.Add("AcceptanceDenial", []byte(setLength("/"+err.Error(), 16)))
		response.Add("AdditionalMessages", []byte{})
		chksum := p.driver.CalcChecksum(response.Addr, response)
		response.Add("Checksum", chksum)
		err := p.SendPacket(addr, response, action)
		return err
	}

	var postingResponse record.PostingResponse
	{
		var ok bool

		if postingResponse, ok = reply.(record.PostingResponse); !ok {
			response.Add("AcceptanceDenial", []byte(setLength("/unknown answer", 16)))
			response.Add("AdditionalMessages", []byte{})
			chksum := p.driver.CalcChecksum(response.Addr, response)
			response.Add("Checksum", chksum)
			err := p.SendPacket(addr, response, action)
			return err
		}
	}

	successMessage := setLength(paymentType, 16)

	switch postingResponse.ResponseCode {
	case 1: // posted successfully
	case 2:
		successMessage = setLength("/credit lim. excd", 16) // Denied due to exceeding credit limit
	case 3:
		successMessage = setLength("/nopost active", 16) // Denied due to NoPost activated for this guest
	case 4:
		successMessage = setLength("/res. not found", 16) // Reservation not found

	default: // Denied due to a not defined reason
		successMessage = setLength("/unknown reason", 16)
	}

	response.Add("AcceptanceDenial", []byte(successMessage))
	response.Add("AdditionalMessages", []byte{})

	chksum := p.driver.CalcChecksum(response.Addr, response)
	response.Add("Checksum", chksum)

	if len(chksum) > 0 {
		if p.driver.CheckChecksum(packet) {
			err = p.SendPacket(addr, response, action)
			return err

		}
	} else {

		err = p.SendPacket(addr, response, action)
		return err
	}

	return errors.New("Checksum Mismatch")

}

func (p *Plugin) getTotals(packet *ifc.LogicalPacket, totalName string, limit int) []record.ChargeItem {

	var totals []record.ChargeItem

	for i := 1; i <= limit; i++ {

		refrence := fmt.Sprintf("%s%dTotal", totalName, i)
		subTotal := cast.ToFloat64(p.getField(packet, refrence, true))

		if subTotal != 0 {

			totals = append(totals, record.ChargeItem{
				ID:    i,
				Value: subTotal,
			})
		}
	}

	return totals
}
