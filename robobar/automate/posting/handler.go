package posting

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/robobar/template"

	"github.com/pkg/errors"
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

		case template.PacketSale:
			err = p.handlePostCharge(addr, packet)

		default:
			err = errors.Errorf("missing handler for '%s'", packet.Name)
		}

		if err != nil {
			log.Error().Msgf("%s addr '%s' process packet '%s' failed, err=%s", name, addr, packet.Name, err.Error())

		}

		if automate.Dispatcher().IsShutdown(err) {
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
