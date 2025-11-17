package callpacket

import (
	"fmt"

	"math"
	"regexp"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

var (
	pattern, _       = regexp.Compile("( [0-9]{5} [0-9]{6} )")
	logStoredRecords = true
)

type PbxRecord struct {
	Station       uint64
	LineNumber    int
	RecordNumber  string
	RecordType    string
	DialledNumber string
	Extension     string
	Amount        float64
	Units         int
	CallTime      time.Time
	Duration      time.Duration
}

func (p *Plugin) handleCallPacket(addr string, packet *ifc.LogicalPacket) {

	name := p.GetName()
	dispatcher := p.GetDispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	// N, S, E, X, L
	recordType := p.getValueAsString(packet, "RecordType")
	if recordType != "N" && recordType != "S" && recordType != "E" && recordType != "X" && recordType != "L" {
		log.Debug().Msgf("%s addr '%s' unknown record type '%s'", name, addr, recordType)
		return
	}

	if recordType == "L" && !p.driver.PostInternalCall(station) {
		return
	}

	extension := p.getExtension(packet)
	extension = p.normalize(extension)

	if len(extension) == 0 && (recordType == "N" || recordType == "L") {
		log.Warn().Msgf("%s addr '%s' extension not found", name, addr)
		return
	}

	// T000007 11/12 21:20 00:05:26 900492319159350
	pl1 := p.getValueAsString(packet, "Payload1")

	lineNumber := p.getLineNumber(pl1)
	if lineNumber <= 0 {
		log.Warn().Msgf("%s addr '%s' line number not found, or not three digits numeric", name, addr)
		return
	}

	// make sure we have mm/dd in front
	// 11/12 21:20 00:05:26 900492319159350
	index := strings.Index(pl1, "/")
	if index <= 2 {
		log.Warn().Msgf("%s addr '%s' can not extract call information, date offset not found", name, addr)
		return
	}
	pl1 = pl1[index-2:]

	data := strings.Split(pl1, " ")
	if len(data) < 3 {
		log.Warn().Msgf("%s addr '%s' can not extract call information, missing parameter", name, addr)
		return
	}

	// trim whitespace
	for i := range data {
		data[i] = strings.Trim(data[i], " ")
	}

	var err error
	var timestamp time.Time
	var duration time.Duration
	var units int
	var amount float64
	var dialledNumber string

	callTimestamp := data[0] + data[1]
	callDuration := data[2]

	callTimestamp = p.normalize(callTimestamp)
	callDuration = p.normalize(callDuration)

	// extract call time
	switch len(callTimestamp) {

	case 10: // mm/ddhh:mm
		timestamp, err = time.Parse("01/0215:04", callTimestamp)

	case 13: // mm/ddhh:mm:ss
		timestamp, err = time.Parse("01/0215:04:05", callTimestamp)

	default:
		log.Warn().Msgf("%s addr '%s' can not extract time information, unknown date '%s' and time '%s' format", name, addr, data[0], data[1])
		return
	}

	if err != nil {
		log.Warn().Msgf("%s addr '%s' can not extract time information, parse date '%s' and time '%s' failed, err=%s", name, addr, data[0], data[1], err)
		return
	}

	// extract call duration
	switch len(callDuration) {

	case 5: // hh:mm
		duration, err = time.ParseDuration(fmt.Sprintf("%dh%dm", callDuration[:2], callDuration[3]))

	case 8, 10: // hh:mm:ss, hh:mm:ss.x
		duration, err = time.ParseDuration(fmt.Sprintf("%sh%sm%ss", callDuration[:2], callDuration[3:5], callDuration[6:8]))

	default:
		log.Warn().Msgf("%s addr '%s' can not extract call durationn, unknown '%s' format", name, addr, data[2])
		return
	}

	if err != nil {
		log.Warn().Msgf("%s addr '%s' can not extract call duration, parse duration '%s' failed, err=%s", name, addr, data[2], err)
		return
	}

	// extract dialled number
	for i := 3; i < len(data); i++ {
		if len(data[i]) > 1 {
			dialledNumber = data[i]
			break
		}
	}

	dialledNumber = p.extractDialledNumber(station, dialledNumber)

	// extract units, amount from second line
	if pattern != nil {
		if pl2, exist := p.getValue(packet, "Payload2"); exist && pattern.Match(pl2) { // match pattern " xxxxx xxxxxx "
			if match := pattern.FindSubmatch(pl2); len(match) > 0 {
				data := strings.Trim(string(match[0]), " ")
				units = p.toInt(data[:5])
				amount = p.toFloat64(data[6:])
			}
		}
	}

	recordNumber := p.getValueAsString(packet, "RecordNumber")

	pbx := PbxRecord{
		Station:       station,
		LineNumber:    lineNumber,
		RecordType:    recordType,
		RecordNumber:  recordNumber,
		Extension:     extension,
		DialledNumber: dialledNumber,
		CallTime:      timestamp,
		Duration:      duration,
		Units:         units,
		Amount:        amount,
	}

	if recordType == "N" || recordType == "L" {
		p.createPosting(addr, packet, pbx)
		return
	}

	postAllRecords := p.driver.PostAllRecords(station)
	postLastRecordOnly := p.driver.PostLastRecordOnly(station)
	mergeRecords := p.driver.MergeRecords(station)

	idx := p.getIndex(station, lineNumber)

	switch recordType {

	case "S":

		p.storeRecord(idx, pbx)

		if postAllRecords {
			p.createPosting(addr, packet, pbx)
		}

	case "X":

		stored, exist := p.getRecord(idx)
		if exist {
			stored.Duration = stored.Duration + pbx.Duration
			stored.Units = stored.Units + pbx.Units
			stored.Amount = stored.Amount + pbx.Amount

			p.storeRecord(idx, stored)
		}

		if postAllRecords {
			if len(pbx.DialledNumber) == 0 && exist {
				pbx.DialledNumber = stored.DialledNumber
			}
			if len(pbx.RecordNumber) == 0 && exist {
				pbx.RecordNumber = stored.RecordNumber
			}

			p.createPosting(addr, packet, pbx)
		}

	case "E":

		stored, exist := p.getRecord(idx)
		if exist {
			stored.Duration = stored.Duration + pbx.Duration
			stored.Units = stored.Units + pbx.Units
			stored.Amount = stored.Amount + pbx.Amount
		}

		if postAllRecords || postLastRecordOnly {

			if len(pbx.DialledNumber) == 0 && exist {
				pbx.DialledNumber = stored.DialledNumber
			}
			if len(pbx.RecordNumber) == 0 && exist {
				pbx.RecordNumber = stored.RecordNumber
			}

			p.createPosting(addr, packet, pbx)

		} else if mergeRecords && exist {

			if len(pbx.RecordNumber) != 0 {
				stored.RecordNumber = pbx.RecordNumber
			}

			stored.Extension = pbx.Extension

			p.createPosting(addr, packet, stored)

		}

		p.removeRecord(idx)

	}

}

func (p *Plugin) createPosting(addr string, packet *ifc.LogicalPacket, pbx PbxRecord) {

	if len(pbx.Extension) == 0 {
		name := p.GetName()
		log.Error().
			Interface("record", pbx).
			Msgf("%s addr '%s' ignore posting, extension not found", name, addr)
		return
	}

	dispatcher := p.GetDispatcher()

	ctx := record.PbxPosting{
		PhoneNumber: pbx.DialledNumber,
		Units:       pbx.Units,
	}

	if totalSeconds := int(math.Floor(pbx.Duration.Seconds())); totalSeconds > 0 {
		hours := int(totalSeconds / 3600)
		totalSeconds = totalSeconds % 3600
		minutes := int(math.Floor(float64(totalSeconds) / 60))
		seconds := int(totalSeconds % 60)
		ctx.Duration = p.toInt(fmt.Sprintf("%02d%02d%02d", hours, minutes, seconds))
	}

	posting := record.SimplePosting{
		Context:     ctx,
		Station:     pbx.Station,
		Time:        pbx.CallTime,
		PostingType: 2, // 2 = pbx
		RoomName:    pbx.Extension,
		CheckNumber: pbx.RecordNumber,
	}

	if pbx.Amount > 0 {
		posting.TotalAmount = pbx.Amount
		posting.PostingType = 1 // 1 = direct
	}

	dispatcher.CreatePmsJob(addr, packet, posting)
}

func (p *Plugin) getIndex(station uint64, lineNumber int) string {
	return fmt.Sprintf("%d-%d", station, lineNumber)
}

func (p *Plugin) getLineNumber(data string) int {

	// Trrrmm or Arrrmmm
	// rrr = route number
	// mmm = member number -> line number

	data = p.normalize(data)

	if len(data) < 7 {
		return 0
	}

	v := data[4 : 4+3]

	v = strings.TrimLeft(v, "0")

	return cast.ToInt(v)

}

func (p *Plugin) extractDialledNumber(station uint64, str string) string {

	pattern, err := p.driver.GetDialledNumberRegex(station)
	if err != nil || pattern == nil {
		return str
	}

	data := []byte(str)

	if !pattern.Match(data) {
		return str
	}

	if match := pattern.FindSubmatch(data); len(match) > 1 {
		return string(match[1])
	}

	return str

}

func (p *Plugin) getRecord(index string) (PbxRecord, bool) {
	p.recordsGuard.RLock()
	pbx, exist := p.records[index]
	p.recordsGuard.RUnlock()
	return pbx, exist
}

func (p *Plugin) storeRecord(index string, pbx PbxRecord) {

	p.recordsGuard.Lock()
	old, exist := p.records[index]
	p.records[index] = pbx
	p.recordsGuard.Unlock()

	if !logStoredRecords {
		return
	}

	name := p.GetName()
	if exist {
		log.Debug().
			Interface("record", pbx).
			Interface("previous", old).
			Msgf("%s pbx record index '%s' updated", name, index)
	} else {
		log.Debug().
			Interface("record", pbx).
			Msgf("%s pbx record index '%s' stored", name, index)
	}

}

func (p *Plugin) removeRecord(index string) {

	p.recordsGuard.Lock()
	_, exist := p.records[index]
	delete(p.records, index)
	p.recordsGuard.Unlock()

	if !logStoredRecords {
		return
	}

	if exist {
		name := p.GetName()
		log.Debug().Msgf("%s pbx record index '%s' removed", name, index)
	}
}

func (p *Plugin) normalize(data string) string {
	data = strings.Replace(data, string(byte(0x0)), "", -1)
	return data
}

func (p *Plugin) toInt(data string) int {
	data = strings.TrimLeft(data, "0")
	return cast.ToInt(data)
}

func (p *Plugin) toFloat64(data string) float64 {
	data = strings.TrimLeft(data, "0")
	return cast.ToFloat64(data)
}
