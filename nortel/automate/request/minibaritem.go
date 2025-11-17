package request

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleMinibar(addr string, packet *ifc.LogicalPacket) {

	name := p.GetName()

	extension := p.getExtension(packet)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	article := p.getValueAsNumeric(packet, "ItemCode")
	quantity := p.getValueAsNumeric(packet, "Quantity")
	unitPrice := p.getValueAsNumeric(packet, "UnitPrice")

	total := p.getValueAsNumeric(packet, "Total")
	subTotal := p.getValueAsNumeric(packet, "SubTotal")

	tax1 := p.getValueAsNumeric(packet, "Tax1")
	tax2 := p.getValueAsNumeric(packet, "Tax2")
	tax3 := p.getValueAsNumeric(packet, "Tax3")

	description := p.getValueAsString(packet, "Description")

	posDate := p.getValueAsString(packet, "Date")
	posTime := p.getValueAsString(packet, "Time")

	postingTime, err := time.Parse("200601021504", posDate+posTime)
	if err != nil {
		postingTime = time.Now()
	}

	taxes := []record.ChargeItem{}
	if tax1 != 0 {
		taxes = append(taxes, record.ChargeItem{
			ID:    1,
			Value: cast.ToFloat64(tax1),
		})
	}
	if tax2 != 0 {
		taxes = append(taxes, record.ChargeItem{
			ID:    2,
			Value: cast.ToFloat64(tax2),
		})
	}
	if tax3 != 0 {
		taxes = append(taxes, record.ChargeItem{
			ID:    3,
			Value: cast.ToFloat64(tax3),
		})
	}

	charge := record.ChargePosting{}

	if article != 0 {
		item := record.Article{
			Number:   article,
			Quantity: quantity,
		}
		charge.MinibarInfo.Items = append(charge.MinibarInfo.Items, item)
	}

	if unitPrice != 0 {
		charge.References = append(charge.References, record.ChargeItem{
			ID:    1,
			Name:  description,
			Value: cast.ToFloat64(unitPrice),
		})
	}

	if subTotal != 0 {
		charge.SubTotals = append(charge.SubTotals, record.ChargeItem{
			ID:    1,
			Value: cast.ToFloat64(subTotal),
		})
	}

	if len(taxes) > 0 {
		charge.Taxes = taxes
	}

	posting := record.SimplePosting{
		PostingType: 1,
		Time:        postingTime,
		Station:     station,
		Context:     charge,
		RoomName:    extension,
		TotalAmount: cast.ToFloat64(total),
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

}
