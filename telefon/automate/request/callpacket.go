package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/weareplanet/ifcv5-drivers/telefon/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func (p *Plugin) handleCallPacket(addr string, packet *ifc.LogicalPacket) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	// extension
	extension, _ := p.getString(packet, template.Extension)
	if len(extension) == 0 {
		log.Error().Msgf("%s addr '%s' extension not found", name, addr)
		return errors.New("extension not found")
	}

	// dialed number
	dialedNumber, _ := p.getString(packet, template.DialedNumber)

	// pbx record
	pbx := record.PbxPosting{
		PhoneNumber: dialedNumber,
	}

	// helper function
	toInt := func(data string) int {
		data = strings.TrimLeft(data, "0")
		return cast.ToInt(data)
	}

	// construct duration
	duration, durationFormat := p.getString(packet, template.Duration)
	if d, err := p.constructDuration(duration, durationFormat); err == nil {
		pbx.Duration = toInt(d.Format("150405"))
	} else {
		log.Warn().
			Str("duration", duration).
			Str("durationformat", durationFormat).
			Msgf("%s addr '%s' construct duration failed, err=%s", name, addr, err)
	}

	// units
	pbx.Units = p.getNumeric(packet, template.Units)

	// calltype
	if callType, _ := p.getString(packet, template.CallType); len(callType) > 0 {
		pbx.CallType = callType
	}

	station, _ := dispatcher.GetStationAddr(addr)

	// posting record
	posting := record.SimplePosting{
		Station:     station,
		Context:     pbx,
		RoomName:    extension,
		PostingType: 2, // 2 = pbx, 1 = direct
	}

	if amount, amountFormat := p.getString(packet, template.Amount); len(amount) > 0 {

		var total float64

		if len(amountFormat) > 0 {
			dp := cast.ToInt32(amountFormat)
			if dp > 0 {
				amount = strings.Replace(amount, ",", "", -1)
				amount = strings.Replace(amount, ".", "", -1)
				data := cast.ToFloat64(amount)
				total = dispatcher.FormatVendorAmount(data, dp)
			}
		}

		if total == 0 {
			total, _ = dispatcher.GetAmountFromString(amount)
		}

		if total > 0 {
			posting.PostingType = 1
			posting.TotalAmount = total
		}
	}

	// construct call timestamp
	callTime, callTimeFormat := p.getString(packet, template.CallTime)
	callDate, callDateFormat := p.getString(packet, template.CallDate)
	if callTimestamp, err := p.constructDate(callDate, callDateFormat, callTime, callTimeFormat); err == nil {
		posting.Time = callTimestamp
	} else if len(callDate) > 0 || len(callTime) > 0 {
		log.Warn().
			Str("calltime", callTime).
			Str("calldate", callDate).
			Str("calltimeformat", callTimeFormat).
			Str("calldateformat", callDateFormat).
			Msgf("%s addr '%s' construct call time failed, err=%s", name, addr, err)
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

	return nil

}

func (p *Plugin) constructDuration(durationValue string, durationFormat string) (time.Time, error) {

	var ts time.Time
	if len(durationValue) == 0 {
		return ts, nil
	}

	durationFormat = strings.ToLower(durationFormat)

	if len(durationValue) < len(durationFormat) {
		durationValue = pad.Left(durationValue, len(durationFormat), "0")
	} else if len(durationValue) > len(durationFormat) {
		durationValue = durationValue[len(durationValue)-len(durationFormat):]
	}

	var hours, minutes, seconds = "0", "0", "0"

	if start, length := p.findPosition(durationFormat, 'h'); start >= 0 {
		hours = durationValue[start : start+length]
	}

	if start, length := p.findPosition(durationFormat, 'm'); start >= 0 {
		minutes = durationValue[start : start+length]
	}

	if start, length := p.findPosition(durationFormat, 's'); start >= 0 {
		seconds = durationValue[start : start+length]
	}

	parse := fmt.Sprintf("%sh%sm%ss", hours, minutes, seconds)
	duration, err := time.ParseDuration(parse)
	if err == nil {
		ts = ts.Add(duration)
	}

	return ts, err
}

func (p *Plugin) findPosition(str string, char byte) (int, int) {
	var length int
	index := strings.IndexByte(str, char)
	if index >= 0 {
		for i := index; i < len(str); i++ {
			if str[i] == char {
				length++
			} else if str[i] != char {
				break
			}
		}
	}
	return index, length
}
