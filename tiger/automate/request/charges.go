package request

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/tiger/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"

	// "github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleCharges(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "RoomNumber", true)
	extension = strings.TrimLeft(extension, "0")

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	productID := cast.ToInt(p.getField(packet, "PPID", false))
	description := p.getField(packet, "Description", false)

	quantity := cast.ToInt(p.getField(packet, "Quantity", false))

	cost := cast.ToFloat64(p.getField(packet, "Cost", true))

	duration := p.getField(packet, "Duration", true)
	dateformatString := "150403"
	if p.driver.Protocol(station) == 0 {
		dateformatString = "1504"
	} else if p.driver.Protocol(station) == 3 {

		date := p.getField(packet, "Date", true)
		t := p.getField(packet, "Time", true)
		dateformatString = "0201061504"

		if len(date) == 4 {
			dateformatString = "02011504"
		}
		duration = date + t

	}

	timestamp, err := time.Parse(dateformatString, duration)
	if err != nil {
		timestamp = time.Now()
		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}

	posting := record.SimplePosting{
		Time:     timestamp,
		Station:  station,
		RoomName: extension,

		TotalAmount: cost,
	}

	switch packet.Name {
	case template.PacketMiniBarBilling:
		posting.PostingType = 3

		item := record.Article{
			Number:   productID,
			Quantity: quantity,
		}

		articles := record.ArticlePosting{}
		articles.Items = append(articles.Items, item)
		posting.Context = articles

	case template.PacketOtherChargePosting:
		posting.PostingType = 1

		charge := record.ChargeItem{
			ID:    productID,
			Name:  description,
			Value: cost,
		}

		charges := record.ChargePosting{}
		charges.SubTotals = append(charges.SubTotals, charge)
		posting.Context = charges
	}

	dispatcherObj.CreatePmsJob(addr, packet, posting)

}
