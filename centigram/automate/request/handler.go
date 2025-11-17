package request

import (
	"strings"

	"github.com/weareplanet/ifcv5-drivers/centigram/template"
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

		case template.PacketMessageWaitingStatus:
			p.handleRoomStatus(addr, packet)

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

func (p *Plugin) getExtension(packet *ifc.LogicalPacket) string {

	if packet == nil {
		return ""
	}

	data := packet.Data()
	value := data["Extension"]
	extension := cast.ToString(value)
	extension = strings.Trim(extension, " ")
	return extension
}

func (p *Plugin) getNumeric(packet *ifc.LogicalPacket, field string) int {

	if packet == nil || len(field) == 0 {
		return 0
	}

	raw := packet.Data()
	data := raw[field]

	value := cast.ToString(data)
	value = strings.Trim(value, " ")
	value = strings.TrimLeft(value, "0")

	ret := cast.ToInt(value)
	return ret
}
