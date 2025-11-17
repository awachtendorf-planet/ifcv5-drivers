package request

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/weareplanet/ifcv5-drivers/simphony/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleNonFrontOfficePostingCharge(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketPostingRequest
	// PI = search inquiry

	dispatcherObj := p.GetDispatcher()

	station, _ := dispatcherObj.GetStationAddr(addr)
	retransmit := setLength((p.getField(packet, "MessageRetransmitFlag", true)), 1)
	sourceID := p.getField(packet, "SourceID", true)
	workStation := p.getWorkStation(sourceID, station)

	//requestInquiry := p.getField(packet, "GuestID", true)
	sequenceNumber := p.getField(packet, "SequenceNumber", true)
	//guestName := p.getField(packet, "GuestName", true)
	//selectionNumber := cast.ToInt(p.getField(packet, "SelectionNumber", true))
	paymentType := p.getField(packet, "PaymentType", true)
	tenderamount := cast.ToFloat64(p.getField(packet, "TenderAmount", true))
	numSalesItemizer := cast.ToInt(p.getField(packet, "NumSalesItemizer", true))

	data := p.getField(packet, "Data", true)

	// Itemizer = SubTotals
	splittedData := strings.Split(data, string(rune(0x1c)))
	subtotals := p.getTotals(splittedData, numSalesItemizer)
	splittedData = splittedData[numSalesItemizer:]

	// Discounts

	numDiscounts := cast.ToInt(splittedData[0])
	splittedData = splittedData[1:]
	discounts := p.getTotals(splittedData, numDiscounts)
	splittedData = splittedData[numDiscounts:]

	// ServiceCharges
	numServiceCharges := cast.ToInt(splittedData[0])
	splittedData = splittedData[1:]
	serviceCharges := p.getTotals(splittedData, numServiceCharges)
	splittedData = splittedData[numServiceCharges:]

	// Taxes
	numTaxes := cast.ToInt(splittedData[0])
	splittedData = splittedData[1:]
	taxes := p.getTotals(splittedData, numTaxes)
	splittedData = splittedData[numTaxes:]

	charge := record.ChargePosting{
		SubTotals:      subtotals,
		Discounts:      discounts,
		ServiceCharges: serviceCharges,
		Taxes:          taxes,
	}

	guestCheckNumber := splittedData[0]
	checkEmployeeNumber := splittedData[1]
	servingPeriod := splittedData[2]
	revenueCenterNumber := cast.ToInt(splittedData[3])
	numberOfGuests := cast.ToInt(splittedData[4])
	dateYYYYMMDD := splittedData[5]
	timeHH_MM_SS := splittedData[6]
	//previousPayment := splittedData[7]

	date, err := time.Parse("2006010215:04:05", dateYYYYMMDD+timeHH_MM_SS)
	if err != nil {
		date = time.Now()
	}
	posting := record.PostingRequest{
		PostingType:    1,
		Station:        station,
		Context:        charge,
		TotalAmount:    tenderamount,
		Time:           date,
		CheckNumber:    guestCheckNumber,
		SequenceNumber: sequenceNumber,
		ServingTime:    servingPeriod,
		Covers:         numberOfGuests,
		OutletNumber:   revenueCenterNumber,
		PaymentMethod:  paymentType,
		UserID:         checkEmployeeNumber,
		WorkStation:    workStation,
	}

	//sequenceNumber = p.incrementSequenceNumber(sequenceNumber)

	response := ifc.NewLogicalPacket(template.PacketChargePostingAck, addr, packet.Tracking)

	reply, err := p.PmsRequest(addr, posting, packet)
	if err != nil {
		errResp := p.buildErrorResponse(err.Error(), retransmit, addr, packet.Tracking, sourceID, sequenceNumber)
		err := p.SendPacket(addr, errResp, action)
		return err
	}

	var postingResponse record.PostingResponse

	{
		var ok bool

		if postingResponse, ok = reply.(record.PostingResponse); !ok {
			errResp := p.buildErrorResponse("posting response not valid", retransmit, addr, packet.Tracking, sourceID, sequenceNumber)
			err := p.SendPacket(addr, errResp, action)
			return err
		}
	}

	successMessage := setLength(paymentType, 30)

	switch postingResponse.ResponseCode {
	case 1: // posted successfully
		successMessage = setLength(fmt.Sprintf("%s", postingResponse.ResponseText), 30)
	case 2:
		successMessage = setLength("credit lim. excd", 30) // Denied due to exceeding credit limit
	case 3:
		successMessage = setLength("nopost active", 30) // Denied due to NoPost activated for this guest
	case 4:
		successMessage = setLength("res. not found", 30) // Reservation not found

	default: // Denied due to a not defined reason
		successMessage = setLength("unknown reason", 30)
	}

	response.Add("SourceID", []byte(setLength(sourceID, 25)))
	response.Add("Message", []byte(successMessage))
	response.Add("SequenceNumber", []byte(sequenceNumber))
	response.Add("MessageRetransmitFlag", []byte(retransmit))
	response.Add("MessageStatus", []byte("N"))
	response.Add("Status", []byte("P"))
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
