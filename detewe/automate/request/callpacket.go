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

func (p *Plugin) handleCallPacket(addr string, packet *ifc.LogicalPacket) {

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	extension := p.getField(packet, "Participant", true)

	if len(extension) == 0 {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	station, _ := dispatcherObj.GetStationAddr(addr)

	dialledNumber := p.getField(packet, "PhoneNumber", true)
	callType := p.getField(packet, "SortSign", false)

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
	if len(duration) == 9 { // hhHmmMssS

		h := duration[0:2]

		m := duration[3:5]

		s := duration[6:8]

		pbx.Duration = toInt(h + m + s)

	} else {
		log.Warn().Msgf("%s addr '%s' construct duration failed, unknown length: %d", name, addr, len(duration))
	}

	units := p.getField(packet, "Taxe", true)
	pbx.Units = toInt(units)

	posting := record.SimplePosting{
		Station:     station,
		Context:     pbx,
		RoomName:    extension,
		PostingType: 2, // 2 = pbx
	}

	// date
	// ddmmyyhhmm

	date := p.getField(packet, "Date", true)
	t := p.getField(packet, "Time", true)

	if timestamp, err := time.Parse("06010215:04", date+t); err == nil {

		posting.Time = timestamp
	} else {

		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}
	dispatcherObj.CreatePmsJob(addr, packet, posting)

	// Send the answer telegram

	NotificationNumber := p.getField(packet, "NotificationNumber", false)

	answer := ifc.NewLogicalPacket(template.PacketTG10Answer, addr, packet.Tracking)

	answer.Add("NotificationNumber", []byte(NotificationNumber))
	answer.Add("Result", []byte("1"))

	p.SendPacket(addr, answer, &dispatcher.StateAction{})

}
