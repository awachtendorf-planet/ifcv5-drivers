package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleCallPacket(addr string, packet *ifc.LogicalPacket) {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcher.GetStationAddr(addr)

	dialledNumber := p.getField(packet, "CallingNumber", true)
	callType := p.getField(packet, "CallType", true)

	pbx := record.PbxPosting{
		PhoneNumber: dialledNumber,
		CallType:    callType,
	}

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	// duration
	duration := p.getField(packet, "Duration", false)
	if len(duration) == 5 { // mmmss
		duration = strings.Replace(duration, " ", "0", -1)
		if d, err := time.ParseDuration(fmt.Sprintf("%dm", toInt(duration[:3]))); err == nil {
			d = d.Round(time.Minute)
			h := d / time.Hour
			d -= h * time.Hour
			m := d / time.Minute
			s := toInt(duration[3:])
			pbx.Duration = toInt(fmt.Sprintf("%02d%02d%02d", h, m, s))
		} else {
			log.Warn().Msgf("%s addr '%s' construct duration failed, err=%s", name, addr, err)
		}
	} else if len(duration) == 6 { // hhmmss
		duration = strings.Replace(duration, " ", "0", -1)
		pbx.Duration = toInt(duration)
	} else {
		log.Warn().Msgf("%s addr '%s' construct duration failed, unknown length: %d", name, addr, len(duration))
	}

	// cost amount or units
	var totalAmount float64
	cost := p.getField(packet, "Cost", false)
	cost = strings.TrimLeft(cost, " ")
	cost = strings.Replace(cost, ",", ".", -1)

	if index := strings.Index(cost, "."); index > 0 {
		totalAmount = cast.ToFloat64(cost)
	} else {
		pbx.Units = toInt(cost)
	}

	posting := record.SimplePosting{
		Station:     station,
		Context:     pbx,
		RoomName:    extension,
		PostingType: 2, // 2 = pbx
	}

	// date
	date := p.getField(packet, "Date", true)
	if timestamp, err := p.constructDate(date); err == nil {
		posting.Time = timestamp
	} else {
		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}

	if totalAmount > 0 {
		posting.TotalAmount = totalAmount
		posting.PostingType = 1 // 1 = direct
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

}
