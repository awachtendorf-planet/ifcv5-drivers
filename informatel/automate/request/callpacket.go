package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

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

	// roomstatus := p.getField(packet, "RoomStatus", true)

	datetimeraw := p.getField(packet, "DateTime", false)

	dialledNumber := p.getField(packet, "DialledDigits", true)

	pbx := record.PbxPosting{
		PhoneNumber: dialledNumber,
	}

	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	units := p.getField(packet, "Units", false)
	duration := p.getField(packet, "Duration", false)

	pbx.Units = toInt(units)
	pbx.Duration = toInt(duration)

	posting := record.SimplePosting{
		Station:  station,
		Context:  pbx,
		RoomName: extension,

		PostingType: 2, // 2 = pbx
	}

	// date
	// YYMMDDHHMM

	dateformatString := "0601021504"

	if timestamp, err := time.Parse(dateformatString, datetimeraw); err == nil {

		posting.Time = timestamp
	} else {

		log.Warn().Msgf("%s addr '%s' construct date failed, err=%s", name, addr, err)
	}
	dispatcherObj.CreatePmsJob(addr, packet, posting)

	return nil

}
