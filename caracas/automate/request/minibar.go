package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/spf13/cast"
)

func (p *Plugin) handleMinibarCharge(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	dispatcherObj := p.GetDispatcher()

	extension := p.getField(packet, "Extension", true)
	name := p.GetName()

	if len(extension) == 0 {
		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	articleNumber := p.getField(packet, "ArticleNumber", true)
	articleAmount := p.getField(packet, "ArticleAmount", true)
	// articleCosts := p.getField(packet, "ArticleCosts", true)

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	// toFloat64 := func(data string) float64 {
	// 	data = strings.TrimLeft(data, "0")
	// 	return cast.ToFloat64(data)
	// }

	article := toInt(articleNumber)
	quantity := toInt(articleAmount)
	if len(articleAmount) == 0 {
		quantity = 1
	}
	// costs := toFloat64(articleCosts)

	item := record.Article{
		Number:   article,
		Quantity: quantity,
	}

	articles := record.ArticlePosting{}
	articles.Items = append(articles.Items, item)

	posting := record.SimplePosting{
		Time:        time.Now(),
		Station:     station,
		Context:     articles,
		PostingType: 3,
		RoomName:    extension,
	}

	// if costs > 0.0 {
	// 	totalAmount := costs * float64(quantity)

	// 	posting.PostingType = 1
	// 	posting.TotalAmount = totalAmount
	// }

	dispatcherObj.CreatePmsJob(addr, packet, posting)

	return nil

}
