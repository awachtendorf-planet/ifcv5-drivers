package request

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

// PostingType represents the posting context
type PostingType int

// posting context
const (
	Unknown PostingType = iota
	Direct
	Pbx
	Minibar
)

func (p *Plugin) handleSimplePosting(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketPostingSimple

	driver := p.driver
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()
	station, _ := dispatcher.GetStationAddr(addr)

	postingTime, _ := driver.ParseTime(packet)

	posting := record.SimplePosting{
		Time:    postingTime,
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &posting); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	context, postingType := p.getContext(addr, packet)
	posting.Context = context
	posting.PostingType = int(postingType)

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	_, pmsErr, sendErr := automate.PmsRequest(station, posting, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}
	if dispatcher.IsShutdown(pmsErr) {
		return pmsErr
	}

	answer, err := driver.MarshalPacket(posting)

	if err != nil {
		log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, posting, err)
	}

	if pmsErr == nil {
		answer.Set("AS", "OK")
		answer.Set("CT", "ACCEPTED")
	} else {
		answer.Set("AS", "UR")
		answer.Set("CT", pmsErr.Error())
	}

	err = p.sendPostingAnswer(addr, action, correlationId, answer)
	return err
}

func (p *Plugin) handlePostingRequest(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketPostingRequest

	driver := p.driver
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()
	station, _ := dispatcher.GetStationAddr(addr)

	postingTime, _ := driver.ParseTime(packet)

	posting := record.PostingRequest{
		Time:    postingTime,
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &posting); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	context, postingType := p.getContext(addr, packet)
	posting.Context = context
	posting.PostingType = int(postingType)

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	_, pmsErr, sendErr := automate.PmsRequest(station, posting, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}

	if dispatcher.IsShutdown(pmsErr) {
		return pmsErr
	}

	answer, err := driver.MarshalPacket(posting)

	if err != nil {
		log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, posting, err)
	}

	if pmsErr == nil {
		answer.Set("AS", "OK")
		answer.Set("CT", "ACCEPTED")
	} else {
		answer.Set("AS", "UR")
		answer.Set("CT", pmsErr.Error())
	}

	err = p.sendPostingAnswer(addr, action, correlationId, answer)
	return err
}

func (p *Plugin) getContext(addr string, packet *ifc.LogicalPacket) (interface{}, PostingType) {

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	postingType := p.getPostingType(packet)

	switch postingType {

	case Direct:
		direct := record.ChargePosting{}
		if err := driver.UnmarshalPacket(packet, &direct); err != nil {
			log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
		}

		getItem := func(key string, value string) {

			if len(key) < 2 || len(value) == 0 {
				return
			}

			id := cast.ToInt(string(key[1]))
			if id > 0 && id < 10 {
				switch key[0] {

				case 'S': // subtotal
					direct.SubTotals = append(direct.SubTotals, record.ChargeItem{Value: cast.ToFloat64(value), ID: id})

				case 'T': // tax
					direct.Taxes = append(direct.Taxes, record.ChargeItem{Value: cast.ToFloat64(value), ID: id})

				case 'D': // discount
					direct.Discounts = append(direct.Discounts, record.ChargeItem{Value: cast.ToFloat64(value), ID: id})

				}
				return
			}
			if key == "SC" {
				direct.ServiceCharges = append(direct.ServiceCharges, record.ChargeItem{Value: cast.ToFloat64(value), ID: 1})
			} else if key == "TP" {
				direct.Tips = append(direct.Tips, record.ChargeItem{Value: cast.ToFloat64(value), ID: 1})
			}
		}

		driver.EachField(packet, getItem)

		// did we have additional infos
		pbx := record.PbxPosting{}
		if err := driver.UnmarshalPacket(packet, &pbx); err == nil {
			direct.TelephoneInfo = pbx
		}

		return direct, postingType

	case Minibar:
		articles := record.ArticlePosting{}
		if err := driver.UnmarshalPacket(packet, &articles); err != nil {
			log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
		}

		var quantity int
		var article int

		getArticle := func(key string, value string) {
			if key == "MA" {
				value = strings.TrimLeft(value, "0")
				article = cast.ToInt(value)
			}
			if key == "M#" {
				value = strings.TrimLeft(value, "0")
				quantity = cast.ToInt(value)
			}
			if article != 0 && quantity > 0 {
				item := record.Article{
					Quantity: quantity,
					Number:   article,
				}
				articles.Items = append(articles.Items, item)
				article = 0
				quantity = 0
			}
		}

		driver.EachField(packet, getArticle)

		return articles, postingType

	case Pbx:
		pbx := record.PbxPosting{}
		if err := driver.UnmarshalPacket(packet, &pbx); err != nil {
			log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
		}

		return pbx, postingType

	case Unknown:
		log.Error().Msgf("%s addr '%s' unknown posting type for packet '%s'", name, addr, packet.Name)

	}

	return nil, Unknown
}

func (p *Plugin) getPostingType(packet *ifc.LogicalPacket) PostingType {

	driver := p.driver

	if postingType, exist := driver.GetField(packet, "PT"); exist {
		switch postingType {
		case "C":
			return Direct
		case "M":
			return Minibar
		case "T":
			return Pbx
		}
	}

	if driver.ExistField(packet, "MA") && driver.ExistField(packet, "M#") {
		return Minibar
	}

	if driver.ExistField(packet, "MP") || driver.ExistField(packet, "DU") {
		return Pbx
	}

	if driver.ExistField(packet, "TA") {
		return Direct
	}

	return Unknown
}

func (p *Plugin) sendPostingAnswer(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketPostingAnswer, tracking, "PA", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
