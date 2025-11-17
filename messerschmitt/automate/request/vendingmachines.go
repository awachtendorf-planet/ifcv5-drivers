package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
)

func (p *Plugin) handleVendingMachines(addr string, packet *ifc.LogicalPacket) error {

	room := p.getRoom(addr, packet)
	if len(room) == 0 {
		return errors.New("room number is empty")
	}

	article := p.getKeyAsInt("Article", packet)
	if article == 0 {
		return errors.New("article number is zero or not numeric")
	}

	reader := p.getKeyAsInt("Reader", packet)

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	resno := p.getKeyAsInt("GuestIndex", packet)

	item := record.Article{
		Number:   article,
		Quantity: 1,
	}

	articles := record.ArticlePosting{}
	articles.Items = append(articles.Items, item)

	ts := p.getTime(addr, packet)

	posting := record.SimplePosting{
		Station:           station,
		Time:              ts,
		OutletNumber:      reader,
		RoomName:          room,
		ReservationNumber: uint64(resno),
		Context:           articles,
		PostingType:       record.PostingByArticle,
	}

	if price := p.getKeyAsInt("Price", packet); price != 0 {
		posting.PostingType = record.PostingByCharge
		posting.TotalAmount = float64(price)
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

	return nil
}
