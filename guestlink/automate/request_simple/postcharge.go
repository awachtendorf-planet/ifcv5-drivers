package requestsimple

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/spf13/cast"
)

func (p *Plugin) handlePostCharge(addr string, packet *ifc.LogicalPacket, retry bool) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	roomNumber, err := p.getRoomNumber(addr, packet)
	if err != nil {
		return unknownRoom, err
	}

	if retry {
		if err, exist := p.lastError[addr]; exist {
			return 0, err
		}
		return 0, nil
	}

	revenueCode := p.getField(packet, "RevenueCode", true)
	description := p.getField(packet, "Description", true)
	purchaseID := p.getField(packet, "PurchaseID", true) // POST enhanced packet

	description = p.decode(addr, description)
	purchaseID = p.decode(addr, purchaseID)

	value := p.getField(packet, "Amount", true)
	value = strings.Replace(value, ",", ".", -1)
	amount := cast.ToFloat64(value)

	item := record.ChargeItem{
		ID:    1,
		Name:  description,
		Value: amount,
	}

	charge := record.ChargePosting{}
	charge.References = append(charge.References, item)

	station, _ := dispatcher.GetStationAddr(addr)

	posting := record.SimplePosting{
		PostingType:  1,
		Station:      station,
		Context:      charge,
		RoomName:     roomNumber,
		OutletNumber: cast.ToInt(revenueCode),
		TotalAmount:  amount,
		CheckNumber:  purchaseID,
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	_, pmsErr, sendErr := automate.PmsRequest(station, posting, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = sendErr
	}

	if pmsErr != nil {
		p.lastError[addr] = pmsErr
	} else {
		delete(p.lastError, addr)
	}

	return 0, pmsErr

}
