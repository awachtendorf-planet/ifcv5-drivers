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

	amount := p.getField(packet, "ArticleAmount", true)
	articleNumber := p.getField(packet, "ArticleNumber", true)

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	article := toInt(articleNumber)
	quantity := toInt(amount)

	item := record.Article{
		Number:   article,
		Quantity: quantity,
	}

	if quantity != 0 && item.Number > 0 {

		articles := record.ArticlePosting{}
		articles.Items = append(articles.Items, item)

		posting := record.SimplePosting{
			Time:        time.Now(),
			Station:     station,
			Context:     articles,
			PostingType: 3,
			RoomName:    extension,
		}

		dispatcherObj.CreatePmsJob(addr, packet, posting)
	}

	return nil

}
