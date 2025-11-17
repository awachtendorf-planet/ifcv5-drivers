package request

import (
	"regexp"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleDataTransfer(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	transfer := p.driver.GetDataTransferType(station)

	data := p.getField(packet, "Data", true)

	switch transfer {

	case 0: // data as article

		article := data
		quantity := 1

		// article zb. 4711
		// article*count zb. 4711*2

		if offset := strings.Index(data, "*"); offset > 0 && offset+1 < len(data) {
			article = data[:offset]
			count := data[offset+1:]
			count = strings.TrimRight(count, "#T")
			count = strings.TrimLeft(count, "0")
			quantity = cast.ToInt(count)
		}

		article = strings.TrimLeft(article, "0")

		item := record.Article{
			Number:   cast.ToInt(article),
			Quantity: quantity,
		}

		if quantity == 0 || item.Number == 0 {
			log.Warn().Msgf("%s addr '%s' cannot interpret data as article, expected 'article' or 'article*quantity' (%s)", name, addr, data)
			return
		}

		articles := record.ArticlePosting{}
		articles.Items = append(articles.Items, item)

		// outlet := p.getOutlet(station, packet)

		posting := record.SimplePosting{
			Time:        time.Now(),
			Station:     station,
			Context:     articles,
			PostingType: record.PostingByArticle,
			RoomName:    extension,
			// OutletNumber: outlet,
		}

		dispatcher.CreatePmsJob(addr, packet, posting)

	case 1: // data as total amount

		amount := cast.ToFloat64(data)

		if amount == 0 {
			log.Warn().Msgf("%s addr '%s' cannot interpret data as total amount (%s)", name, addr, data)
			return
		}

		// posting context is required
		charge := record.ChargePosting{}

		outlet := p.getOutlet(station, packet)

		posting := record.SimplePosting{
			Time:         time.Now(),
			Station:      station,
			Context:      charge,
			PostingType:  record.PostingByCharge,
			RoomName:     extension,
			TotalAmount:  amount,
			OutletNumber: outlet,
		}

		dispatcher.CreatePmsJob(addr, packet, posting)

	default:

		log.Warn().Msgf("%s addr '%s' unknown data transfer type: %d", name, addr, transfer)

	}

}

func (p *Plugin) getOutlet(station uint64, packet *ifc.LogicalPacket) int {

	if packet == nil {
		return 0
	}

	code := packet.Data()["Code"]

	if len(code) == 0 {
		return 0
	}

	dispatcher := p.automate.Dispatcher()

	reg := dispatcher.GetConfig(station, "OutletRegexp", ".{4}(\\d)") // default 5th position from dialled code
	if len(reg) == 0 {
		return 0
	}

	var outlet int

	if pattern, err := regexp.Compile(reg); err == nil && pattern.Match(code) {
		if match := pattern.FindSubmatch(code); len(match) > 1 {
			data := strings.Trim(string(match[1]), " ")
			data = strings.TrimLeft(data, "0")
			outlet = cast.ToInt(data)
		}
	}
	return outlet

}
