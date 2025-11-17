package request

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-drivers/telefon/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func (p *Plugin) handlePosting(addr string, packet *ifc.LogicalPacket) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	// extension
	extension, _ := p.getString(packet, template.Extension)
	if len(extension) == 0 {
		log.Error().Msgf("%s addr '%s' extension not found", name, addr)
		return errors.New("extension not found")
	}

	// total amount and article
	var total float64
	articleNumber := p.getNumeric(packet, template.Article)

	if amount, amountFormat := p.getString(packet, template.Amount); len(amount) > 0 {

		if len(amountFormat) > 0 {
			dp := cast.ToInt32(amountFormat)
			if dp > 0 {
				amount = strings.Replace(amount, ",", "", -1)
				amount = strings.Replace(amount, ".", "", -1)
				data := cast.ToFloat64(amount)
				total = dispatcher.FormatVendorAmount(data, dp)
			}
		}

		if total == 0 {
			total, _ = dispatcher.GetAmountFromString(amount)
		}

	}

	if total == 0 && articleNumber == 0 {
		log.Error().Msgf("%s addr '%s' article or amount to be posted not found", name, addr)
		return errors.New("article or amount to be posted not found")
	}

	station, _ := dispatcher.GetStationAddr(addr)
	user, _ := p.getString(packet, template.User)
	outlet := p.getNumeric(packet, template.Outlet)

	// quantity
	quantity := p.getNumeric(packet, template.Quantity)
	if quantity == 0 {
		quantity = 1
	}

	// posting
	postingType := 1 // direct charge
	if total == 0 {
		postingType = 3 // minibar article
	}

	charge := record.ChargePosting{}

	chargeItem := record.ChargeItem{
		ID:    1,
		Value: total,
		//Name:  ,
	}
	charge.References = append(charge.References, chargeItem)

	article := record.Article{
		Quantity: quantity,
		Number:   articleNumber,
	}
	charge.MinibarInfo.Items = append(charge.MinibarInfo.Items, article)

	posting := record.SimplePosting{
		//Time:        time.Now(),
		Station:      station,
		Context:      charge,
		PostingType:  postingType,
		RoomName:     extension,
		TotalAmount:  total,
		UserID:       user,
		OutletNumber: outlet,
		//SequenceNumber:,

	}

	dispatcher.CreatePmsJob(addr, packet, posting)

	return nil
}
