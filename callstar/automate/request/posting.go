package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleCallPacket(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = charge in pence
	// 1 = hhmm
	// 2 = call duration in tenths of a minute
	// 3 = dialed number
	// 4 = units

	// optional:
	// 5 = message count
	// 6 = call type

	if len(roomNumber) == 0 || len(data) < 5 {
		return
	}

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	station, _ := dispatcher.GetStationAddr(addr)

	pbx := record.PbxPosting{
		PhoneNumber: p.getString(data[3]),
		Units:       p.getNumeric(data[4]),
	}

	duration := p.getNumeric(data[2]) * 6

	var ts time.Time
	if d, err := time.ParseDuration(fmt.Sprintf("%ds", duration)); err == nil {
		ts = ts.Add(d)
		pbx.Duration = p.getNumeric(ts.Format("150405"))
	}

	if len(data) > 6 {
		pbx.CallType = p.getString(data[6])
	}

	posting := record.SimplePosting{
		PostingType: 2, // pbx
		Station:     station,
		Context:     pbx,
		RoomName:    roomNumber,
	}

	if timestamp, err := p.getTime(data[1]); err == nil {
		posting.Time = timestamp
	} else {
		log.Warn().Msgf("%s addr '%s' construct time failed, err=%s", name, addr, err)
	}

	totalAmount := p.getFloat64(data[0])

	if totalAmount > 0 {
		posting.PostingType = 1 // direct
		posting.TotalAmount = totalAmount
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

}

func (p *Plugin) handlePostCharge(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = charge in pence
	// 1 = hhmm
	// 2 = operator
	// 3 = reason

	if len(roomNumber) == 0 || len(data) < 2 {
		return
	}

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	totalAmount := p.getFloat64(data[0])
	if totalAmount == 0 {
		log.Warn().Msgf("%s addr '%s' zero total amount", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	var operator string
	var description string

	if len(data) > 2 {
		operator = p.getString(data[2])
	}
	if len(data) > 3 {
		driver := p.driver
		station, _ := dispatcher.GetStationAddr(addr)
		description = p.getString(data[3])
		description = driver.DecodeString(station, []byte(description))
	}

	charge := record.ChargePosting{}

	chargeItem := record.ChargeItem{
		Name:  description,
		Value: totalAmount,
	}
	charge.References = append(charge.References, chargeItem)

	posting := record.SimplePosting{
		PostingType: 1, // direct
		Station:     station,
		Context:     charge,
		RoomName:    roomNumber,
		UserID:      operator,
		TotalAmount: totalAmount,
	}

	if timestamp, err := p.getTime(data[1]); err == nil {
		posting.Time = timestamp
	} else {
		log.Warn().Msgf("%s addr '%s' construct time failed, err=%s", name, addr, err)
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

}

func (p *Plugin) handleVoiceCharge(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = reason digit code
	// 1 = charge in pence

	if len(roomNumber) == 0 || len(data) < 2 {
		return
	}

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	totalAmount := p.getFloat64(data[1])
	if totalAmount == 0 {
		log.Warn().Msgf("%s addr '%s' zero total amount", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	description := p.getString(data[0])

	charge := record.ChargePosting{}

	chargeItem := record.ChargeItem{
		Name:  description,
		Value: totalAmount,
	}
	charge.References = append(charge.References, chargeItem)

	posting := record.SimplePosting{
		PostingType: 1, // direct
		Station:     station,
		Context:     charge,
		RoomName:    roomNumber,
		TotalAmount: totalAmount,
	}

	dispatcher.CreatePmsJob(addr, packet, posting)
}

func (p *Plugin) handleMinibar(addr string, roomNumber string, data []string, packet *ifc.LogicalPacket) {

	// 0 = ticket
	// 1 = waiter code
	// 2 = xxyy = quantity of yy, multiple times, blank separeted
	// 3 = total charge in pence

	if len(roomNumber) == 0 || len(data) < 4 {
		return
	}

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	totalAmount := p.getFloat64(data[3])
	if totalAmount == 0 {
		log.Warn().Msgf("%s addr '%s' zero total amount", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	ticket := p.getString(data[0])
	operator := p.getString(data[1])

	charge := record.ChargePosting{}

	articles := strings.Split(p.getString(data[2]), " ")

	for i := range articles {

		article := articles[i]
		if len(article) != 4 { // xxyy
			continue
		}

		item := record.Article{
			Quantity: p.getNumeric(article[:2]),
			Number:   p.getNumeric(article[2:]),
		}

		charge.MinibarInfo.Items = append(charge.MinibarInfo.Items, item)
	}

	posting := record.SimplePosting{
		PostingType:    3, // minibar
		Station:        station,
		Context:        charge,
		TotalAmount:    totalAmount,
		RoomName:       roomNumber,
		UserID:         operator,
		SequenceNumber: ticket,
	}

	if p.driver.TicketAsOutlet(addr) {
		posting.OutletNumber = dispatcher.GetPMSOutlet(station, ticket)
	}

	if len(charge.MinibarInfo.Items) == 0 {
		posting.PostingType = 1 // direct
	}

	dispatcher.CreatePmsJob(addr, packet, posting)
}
