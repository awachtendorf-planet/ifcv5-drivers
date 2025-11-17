package request

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/bartech/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	errInternalPacket = errors.New("object packet is nil")
	errInternalJob    = errors.New("job failed")
	errPMSUnavailable = errors.New("PMS is not available")
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// if respectLinkState {
	// 	if !p.linkState(addr) {
	// 		automate.NextAction(name, addr, shutdown, t)
	// 		return
	// 	}
	// }

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case busy:
		if job == nil || !job.InProcess() {
			err = errInternalJob
			break
		}

		if job.Context == nil {
			err = errInternalPacket
			break
		}

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = errInternalPacket
			break
		}

		switch packet.Name {

		case template.PacketSyncRequest:
			p.handleSyncRequest(addr, packet)

		case template.PacketRoomStatus:
			p.handleRoomStatus(addr, packet)

		case template.PacketInvoice:
			p.handleInvoice(addr, packet)

		default:
			log.Warn().Msgf("%s addr '%s' unknown packet '%s'", name, addr, packet.Name)
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		automate.NextAction(name, addr, success, t)
		return

	case success, shutdown:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		if err == errInternalPacket || err == errInternalJob {
			automate.NextAction(name, addr, shutdown, t)
			return
		}
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) getString(packet *ifc.LogicalPacket, field string) string {

	if packet == nil {
		return ""
	}

	raw := packet.Data()
	data := raw[field]
	str := cast.ToString(data)
	str = strings.Trim(str, " ")
	return str
}

func (p *Plugin) getNumeric(packet *ifc.LogicalPacket, field string) int {

	if packet == nil || len(field) == 0 {
		return 0
	}

	raw := packet.Data()
	data := raw[field]

	value := cast.ToString(data)
	value = strings.Trim(value, " ")
	value = strings.TrimLeft(value, "+")
	negative := strings.HasPrefix(value, "-")
	if negative {
		value = strings.TrimLeft(value, "-")
	}
	value = strings.TrimLeft(value, "0")

	ret := cast.ToInt(value)
	if negative {
		ret = ret * -1
	}

	return ret
}

func (p *Plugin) getDate(packet *ifc.LogicalPacket) (time.Time, error) {

	if packet == nil {
		return time.Time{}, errInternalPacket
	}

	dateValue := p.getString(packet, "Date")
	timeValue := p.getString(packet, "Time")

	var dateFormat string
	var timeFormat string

	if len(dateValue) == 6 {
		dateFormat = "020106" //ddmmyy
	} else if len(dateValue) == 8 {
		dateFormat = "01022006" // mmddccyy
	}

	if len(timeValue) == 4 {
		timeFormat = "1504" // hhmm
	} else if len(timeValue) == 6 {
		timeFormat = "150405" // hhmmss
	}

	timestamp, err := p.constructDate(dateValue, dateFormat, timeValue, timeFormat)
	return timestamp, err

}

func (p *Plugin) constructDate(dateValue string, dateFormat string, timeValue string, timeFormat string) (time.Time, error) {

	var timestamp time.Time

	t, _ := time.Parse(timeFormat, timeValue)
	haveTime := !t.IsZero() && len(timeValue) > 0 && len(timeFormat) > 0

	d, err := time.Parse(dateFormat, dateValue)
	if err == nil {
		timestamp = d
		if haveTime {
			timestamp = t.AddDate(d.Year(), int(d.Month()-1), d.Day()-1)
		}
	} else if haveTime {
		timestamp = t
	}

	return timestamp, err
}
