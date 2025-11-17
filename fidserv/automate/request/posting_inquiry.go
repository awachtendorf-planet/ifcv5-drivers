package request

import (
	"sort"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handlePostingInquiry(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketPostingRequest
	// PI = search inquiry

	driver := p.driver
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()
	station, _ := dispatcher.GetStationAddr(addr)

	postingTime, _ := driver.ParseTime(packet)

	inquiry := record.PostingInquiry{
		Time:    postingTime,
		Station: station,
		//RequestType: "15",
	}

	if err := driver.UnmarshalPacket(packet, &inquiry); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	correlationId := dispatcher.NewCorrelationID(station)
	trackingId := dispatcher.GetTrackingId(packet)

	reply, pmsErr, sendErr := automate.PmsRequest(station, inquiry, pmsTimeOut, correlationId, trackingId)
	if sendErr != nil {
		pmsErr = pmsUnavailable
	}
	if dispatcher.IsShutdown(pmsErr) {
		return pmsErr
	}

	answer, err := driver.MarshalPacket(inquiry)
	if err != nil {
		log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, inquiry, err)
	}

	if pmsErr != nil {
		answer.Set("AS", "UR")
		// todo: aus internen Error (no handler found, etc) für den Bediener verständliche Meldung erzeugen
		answer.Set("CT", pmsErr.Error())
		err := p.sendPostingAnswer(addr, action, correlationId, answer)
		return err
	}

	var guests []*record.Guest

	switch reply.(type) {
	case *record.Guest:
		guest := reply.(*record.Guest)
		guests = append(guests, guest)

	case []*record.Guest:
		guests = reply.([]*record.Guest)
	}

	if len(guests) == 0 {
		answer.Set("AS", "NG")
		answer.Set("CT", "no guest found")
		err := p.sendPostingAnswer(addr, action, correlationId, answer)
		return err
	}

	// hier wird es etwas tricky
	// da mehrere Gäste aufgelistet werden können, kann der standard packet constructor nicht verwendet werden
	// und wir müssen das daten packet händisch erstellen

	maxGuests := int(16)
	if max, exist := driver.GetField(packet, "MX"); exist {
		maxGuests = cast.ToInt(max)
	}

	deviceAddr := ifc.DeviceAddress(addr)
	records := driver.GetLinkRecord(deviceAddr, "PL")

	remaining := make(map[string]bool)

	var packetData string

	for i := range guests {

		if i >= maxGuests {
			break
		}

		var data string
		for _, key := range records {

			switch key {
			// this fields are only needed once and do not belong to counting guests
			case "C#", "P#", "PM", "ID", "SO", "DA", "TI", "WS":
				remaining[key] = true
				continue
			}

			if record, err := driver.ConstructRecord(station, "PL", key, guests[i]); err == nil {
				data = data + "|" + key + record
			} else {
				if driver.AppendEmpty(station, "PL", key) {
					data = data + "|" + key
				}
				log.Warn().Msgf("%s", err)
			}

		}

		packetData = packetData + data

	}

	// remaining is a unsorted map, create a sorted slice
	fields := make([]string, len(remaining))
	i := 0
	for key := range remaining {
		fields[i] = key
		i++
	}
	sort.Strings(fields)

	var data string
	for i := range fields {

		key := fields[i]
		if record, err := driver.ConstructRecord(station, "PL", key, answer); err == nil && len(record) > 0 {
			data = data + "|" + key + record
		} else {
			if driver.AppendEmpty(station, "PL", key) {
				data = data + "|" + key
			}
			if err != nil {
				log.Warn().Msgf("%s", err)
			}
		}

	}
	packetData = packetData + data + "|"

	err = p.sendPostingList(addr, action, packetData, correlationId)
	return err
}

func (p *Plugin) sendPostingList(addr string, action *dispatcher.StateAction, data string, tracking string) error {
	packet := ifc.NewLogicalPacket(template.PacketPostingList, addr, tracking)
	packet.Add("Data", []byte(data))
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
