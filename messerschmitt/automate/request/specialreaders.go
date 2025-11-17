package request

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
)

func (p *Plugin) handleSpecialReaders(addr string, packet *ifc.LogicalPacket) error {

	room := p.getRoom(addr, packet)
	if len(room) == 0 {
		return errors.New("room number is empty")
	}

	reader := p.getKeyAsInt("Reader", packet)
	if reader == 0 {
		return errors.New("room number of special reader is empty or not numeric")
	}

	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	resno := p.getKeyAsInt("GuestIndex", packet)

	item := record.Article{
		Number:   reader,
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

	dispatcher.CreatePmsJob(addr, packet, posting)

	return nil
}
