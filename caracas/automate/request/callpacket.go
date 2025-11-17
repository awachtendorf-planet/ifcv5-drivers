package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	// "github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

func (p *Plugin) handleCallPacket(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketPostingRequest
	// PI = search inquiry

	dispatcherObj := p.GetDispatcher()
	name := p.GetName()

	station, _ := dispatcherObj.GetStationAddr(addr)

	extension := p.getField(packet, "Extension", true)

	if len(extension) == 0 {
		return fmt.Errorf("%s addr '%s' extension not found", name, addr)
	}

	// dateRaw := p.getField(packet, "CallDate", false)
	callStartTimeRaw := p.getField(packet, "CallStartTime", false)
	callEndTimeRaw := p.getField(packet, "CallEndTime", false)

	userid := p.getField(packet, "SystemID", true)

	// caller := p.getField(packet, "CallExtension", true) // Might be the same as Extension
	// exchangeLineNumber := p.getField(packet, "ExchangeLineNumber", true) // Like, what?

	dialledNumber := p.getField(packet, "TargetNumber", true)
	// nummernKennung := p.getField(packet, "TargetNumberID", true)

	pbx := record.PbxPosting{
		PhoneNumber: dialledNumber,
	}

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	fmtDuration := func(d time.Duration) string {

		d = d.Round(time.Second)
		h := d / time.Hour
		d = d - (h * time.Hour)
		m := d / time.Minute
		d = d - (m * time.Minute)
		s := d / time.Second

		return fmt.Sprintf("%02d%02d%02d", h, m, s) // HHMMSS
	}

	units := p.getField(packet, "Units", false)

	start, _ := time.Parse("150405", callStartTimeRaw)
	end, _ := time.Parse("150405", callEndTimeRaw)
	duration := end.Sub(start)
	durationString := fmtDuration(duration)

	pbx.Units = toInt(units)
	pbx.Duration = toInt(durationString)

	posting := record.SimplePosting{
		Station:  station,
		Context:  pbx,
		RoomName: extension,
		UserID:   userid,
		Time:     time.Now(),

		PostingType: 2, // 2 = pbx
	}

	dispatcherObj.CreatePmsJob(addr, packet, posting)

	return nil

}
