package request

import (
	"strings"
	"time"

	"github.com/spf13/cast"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleChargeMinibar(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.driver

	station, _ := dispatcherObj.GetStationAddr(addr)

	roomName := p.getField(packet, "RoomNumber")
	if p.driver.LeadingZeros(station) {
		roomName = strings.TrimLeft(roomName, "0")
	}

	total := cast.ToFloat64(p.getField(packet, "Value"))
	articleNo := cast.ToInt(p.getField(packet, "ArticleNumber"))
	articleName := p.getField(packet, "TransText")
	amount := cast.ToInt(p.getField(packet, "Amount"))

	day := p.getField(packet, "Day")
	month := p.getField(packet, "Month")
	year := p.getField(packet, "Year")
	hour := p.getField(packet, "Hour")
	minute := p.getField(packet, "Minute")

	date, err := time.Parse("200601021504", year+month+day+hour+minute)
	if err != nil {
		date = time.Now()
	}

	posting := record.SimplePosting{
		PostingType: 1,
		Station:     station,
		RoomName:    roomName,
		Time:        date,
	}

	if amount > 0 {
		charge := record.ArticlePosting{}
		var article record.Article
		article.Number = articleNo
		article.Quantity = amount
		charge.Items = append(charge.Items, article)

		posting.Context = charge
	}

	if total > 0 {
		charge := record.ChargePosting{}
		var subTotal record.ChargeItem
		subTotal.Value = total
		subTotal.Name = articleName
		charge.SubTotals = append(charge.SubTotals, subTotal)

		posting.Context = charge
		posting.TotalAmount = total
	}

	correlationId := dispatcherObj.NewCorrelationID(station)
	trackingId := dispatcherObj.GetTrackingId(packet)

	_, pmsErr, sendErr := p.automate.PmsRequest(station, posting, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}
	return pmsErr

}

func (p *Plugin) handleChargePayTV(addr string, packet *ifc.LogicalPacket) error {

	dispatcherObj := p.driver

	station, _ := dispatcherObj.GetStationAddr(addr)

	roomName := p.getField(packet, "RoomNumber")
	if p.driver.LeadingZeros(station) {
		roomName = strings.TrimLeft(roomName, "0")
	}

	val := cast.ToFloat64(p.getField(packet, "Value"))
	channel := p.getField(packet, "Channel")

	day := p.getField(packet, "Day")
	month := p.getField(packet, "Month")
	year := p.getField(packet, "Year")
	hour := p.getField(packet, "Hour")
	minute := p.getField(packet, "Minute")

	duration, err := time.Parse("200601021504", year+month+day+hour+minute)
	if err != nil {
		duration = time.Now()
	}

	charge := record.ChargePosting{}
	var subTotal record.ChargeItem
	subTotal.Value = val
	subTotal.Name = channel
	charge.SubTotals = append(charge.SubTotals, subTotal)

	total := val

	posting := record.SimplePosting{
		PostingType: 1,
		Station:     station,
		Context:     charge,
		TotalAmount: total,
		RoomName:    roomName,
		Time:        duration,
	}

	correlationId := dispatcherObj.NewCorrelationID(station)
	trackingId := dispatcherObj.GetTrackingId(packet)

	_, pmsErr, sendErr := p.automate.PmsRequest(station, posting, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}
	return pmsErr

}
