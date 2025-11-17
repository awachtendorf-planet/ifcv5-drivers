package request

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleCallPacket(addr string, packet *ifc.LogicalPacket) {

	station, _ := p.driver.GetStationAddr(addr)
	if p.driver.Protocol(station) < 3 {
		p.protocolOpera(addr, packet)
		return
	}
	p.protocolTiger(addr, packet)

}

func (p *Plugin) protocolOpera(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "RoomNumber", true)
	extension = strings.TrimLeft(extension, "0")

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	dialledNumber := p.getField(packet, "DialledNumber", true)

	pbx := record.PbxPosting{
		PhoneNumber: dialledNumber,
	}

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	// duration
	duration := p.getField(packet, "Duration", false)

	pbx.Duration = toInt(duration)

	cost := cast.ToFloat64(p.getField(packet, "Cost", true))

	posting := record.SimplePosting{
		Station:     station,
		Context:     pbx,
		RoomName:    extension,
		TotalAmount: cost,

		PostingType: 2, // 2 = pbx
	}

	// date
	// ddmmyyhhmm

	t := p.getField(packet, "Time", true)
	dateformatString := "1504"

	if timestamp, err := time.Parse(dateformatString, t); err == nil {

		posting.Time = timestamp
	} else {

		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}
	dispatcherObj.CreatePmsJob(addr, packet, posting)

}

func (p *Plugin) protocolTiger(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "RoomNumber", true)
	extension = strings.TrimLeft(extension, "0")

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	dialledNumber := p.getField(packet, "DialledNumber", true)
	conditionCode := p.getField(packet, "ConditionCode", true)

	pbx := record.PbxPosting{
		PhoneNumber: dialledNumber,
		CallType:    string(conditionCode[0]),
	}

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	// duration
	duration := p.getField(packet, "Duration", false)
	if len(duration) == 6 { // hhmmss

		pbx.Duration = toInt(duration)

	} else {
		log.Warn().Msgf("%s addr '%s' construct duration failed, unknown length: %d", name, addr, len(duration))
	}

	sequence := p.getField(packet, "Sequence", true)
	cost := cast.ToFloat64(p.getField(packet, "Cost", true))

	posting := record.SimplePosting{
		Station:        station,
		Context:        pbx,
		RoomName:       extension,
		SequenceNumber: sequence,
		TotalAmount:    cost,

		PostingType: 2, // 2 = pbx
	}

	// date
	// ddmmyyhhmm

	date := p.getField(packet, "Date", true)
	t := p.getField(packet, "Time", true)
	dateformatString := "0201061504"

	if len(date) == 4 {
		dateformatString = "02011504"
	}

	if timestamp, err := time.Parse(dateformatString, date+t); err == nil {

		posting.Time = timestamp
	} else {

		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}
	dispatcherObj.CreatePmsJob(addr, packet, posting)

}
