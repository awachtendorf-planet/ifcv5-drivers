package request

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/ahl/template"
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

		case template.PacketRoomStatus:
			p.handleRoomStatus(addr, packet)

		case template.PacketWakeupEvent:
			p.handleWakeupEvent(addr, packet)

		case template.PacketDataTransfer:
			p.handleDataTransfer(addr, packet)

		case template.PacketCallPacket, template.PacketCallPacketExtended:
			p.handleCallPacket(addr, packet)

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

func (p *Plugin) getField(packet *ifc.LogicalPacket, field string, trim bool) string {

	if packet == nil {
		return ""
	}

	data := packet.Data()
	value := cast.ToString(data[field])
	if trim {
		value = strings.Trim(value, " ")
	}
	return value
}

func (p *Plugin) constructDate(date string) (time.Time, error) {
	var layout string

	if len(date) == 0 {
		errors.New("empty date object")
	}

	switch len(date) {
	case 4: // hhmm
		layout = "1504"

	case 5: // hhmm(P/A)
		if date[len(date)-1] == 'A' || date[len(date)-1] == 'P' {
			layout = "1504PM"
			date += "M"
		}

	case 10: // ddmmyyhhmm
		layout = "0201061504"

	case 11: // ddmmyyhhmm(P/A)
		if date[len(date)-1] == 'A' || date[len(date)-1] == 'P' {
			layout = "0201061504PM"
			date += "M"
		}

	case 14: // ddmmyyyyhhmmss
		layout = "02012006150405"

	case 15: // ddmmyyyyhhmmss(P/A)
		if date[len(date)-1] == 'A' || date[len(date)-1] == 'P' {
			layout = "02012006150405PM"
			date += "M"
		}

	}

	constructed, err := time.Parse(layout, date)
	return constructed, err
}
