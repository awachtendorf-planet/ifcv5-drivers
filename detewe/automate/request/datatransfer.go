package request

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/detewe/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleDataTransfer(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "Participant", true)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	data := p.getField(packet, "Data", false)

	article := data[:2]
	articlecount := string(data[2])
	quantity := cast.ToInt(articlecount)
	article = strings.TrimLeft(article, " ")
	article = strings.TrimLeft(article, "0")

	item := record.Article{
		Number:   cast.ToInt(article),
		Quantity: quantity,
	}

	date := p.getField(packet, "Date", true)
	t := p.getField(packet, "Time", true)

	timestamp, err := time.Parse("06010215:04", date+t)
	if err != nil {
		timestamp = time.Now()
		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}

	if quantity != 0 && item.Number > 0 {

		articles := record.ArticlePosting{}
		articles.Items = append(articles.Items, item)

		posting := record.SimplePosting{
			Time:        timestamp,
			Station:     station,
			Context:     articles,
			PostingType: 3,
			RoomName:    extension,
		}

		dispatcherObj.CreatePmsJob(addr, packet, posting)

		// Send the answer telegram

		NotificationNumber := p.getField(packet, "NotificationNumber", false)

		answer := ifc.NewLogicalPacket(template.PacketTG40Answer, addr, packet.Tracking)

		answer.Add("NotificationNumber", []byte(NotificationNumber))
		answer.Add("Result", []byte("1"))

		p.SendPacket(addr, answer, &dispatcher.StateAction{})

	}
}
