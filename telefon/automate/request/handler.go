package request

import (
	"fmt"
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/telefon/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	internalPacketError = errors.New("object packet is nil")
	internalJobError    = errors.New("job failed")
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	if respectLinkState {
		if !p.linkState(addr) {
			automate.NextAction(name, addr, shutdown, t)
			return
		}
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case busy:
		if job == nil || !job.InProcess() {
			err = internalJobError
			break
		}

		if job.Context == nil {
			err = internalPacketError
			break
		}

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = internalPacketError
			break
		}

		switch packet.Name {

		case template.CallPacket:
			err = p.handleCallPacket(addr, packet)

		case template.RoomStatus:
			err = p.handleRoomStatus(addr, packet)

		case template.Posting:
			err = p.handlePosting(addr, packet)

		default:
			log.Warn().Msgf("%s addr '%s' unknown packet '%s'", name, addr, packet.Name)
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		if err != nil {
			dispatcher := automate.Dispatcher()
			dispatcher.SetAlive(addr)
			automate.NextAction(name, addr, shutdown, t)
		} else {
			automate.NextAction(name, addr, success, t)
		}
		return

	case success, shutdown:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		if err == internalPacketError || err == internalJobError {
			automate.NextAction(name, addr, shutdown, t)
			return
		}
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) getString(packet *ifc.LogicalPacket, name string) (string, string) {

	value := ""
	format := ""

	if packet == nil {
		return value, format
	}

	data := packet.Data()

	if raw, exist := data[name]; exist {
		value = cast.ToString(raw)
		value = strings.Trim(value, " ")
	}

	if raw, exist := data[fmt.Sprintf("i:format:%s", name)]; exist {
		format = cast.ToString(raw)
	}

	return value, format
}

func (p *Plugin) getNumeric(packet *ifc.LogicalPacket, name string) int {

	value := 0

	if packet == nil {
		return value
	}

	data := packet.Data()

	if raw, exist := data[name]; exist {
		str := cast.ToString(raw)
		str = strings.Trim(str, " ")
		str = strings.TrimLeft(str, "0")
		value = cast.ToInt(str)
	}

	return value
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
