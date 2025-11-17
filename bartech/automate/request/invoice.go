package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleInvoice(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getNumeric(packet, "Extension")

	if extension == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	roomNumber := cast.ToString(extension)

	station, _ := dispatcher.GetStationAddr(addr)

	// article number = TaxID, opera bullshit
	articleNumber := p.getNumeric(packet, "TaxID")
	quantity := p.getNumeric(packet, "Quantity")

	//taxID := p.getNumeric(packet, "TaxID")
	//taxAmount := float64(p.getNumeric(packet, "TaxAmount"))
	chargeAmount := float64(p.getNumeric(packet, "Charge"))
	ticket := p.getString(packet, "Ticket")
	description := p.getString(packet, "Text")

	description = p.decodeString(addr, description)

	charge := record.ChargePosting{}

	chargeItem := record.ChargeItem{
		ID:    articleNumber,
		Name:  description,
		Value: chargeAmount,
	}
	charge.References = append(charge.References, chargeItem)

	// tax := record.ChargeItem{
	// 	ID:    taxID,
	// 	Name:  description,
	// 	Value: taxAmount,
	// }

	// charge.Taxes = append(charge.Taxes, tax)

	article := record.Article{
		Quantity: quantity,
		Number:   articleNumber,
	}
	charge.MinibarInfo.Items = append(charge.MinibarInfo.Items, article)

	posting := record.SimplePosting{
		Station:        station,
		Context:        charge,
		PostingType:    1, // direct charge
		RoomName:       roomNumber,
		SequenceNumber: ticket,
		TotalAmount:    chargeAmount,
		OutletNumber:   articleNumber, // opera bullshit
	}

	if timestamp, err := p.getDate(packet); err == nil {
		posting.Time = timestamp
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

}

func (p *Plugin) decodeString(addr string, data string) string {
	dispatcher := p.automate.Dispatcher()

	encoding := dispatcher.GetEncoding(addr)
	if len(encoding) == 0 {
		return data
	}
	if dec, err := dispatcher.Decode([]byte(data), encoding); err == nil && len(dec) > 0 {
		return string(dec)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return data
}
