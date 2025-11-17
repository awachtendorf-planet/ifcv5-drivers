package dbsync

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, _ *order.Job) {

	automate := p.automate
	name := automate.Name()

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case busy:

		dispatcher := p.automate.Dispatcher()

		for {

			job := dispatcher.GetSyncRecord(addr)
			if job == nil {
				break
			}

			switch representer := job.Context.(type) {

			case *record.Guest:

				representer.SetGeneric(defines.SyncInd, true)
				dispatcher.CreateDriverJob(job.Station, job.Action, representer, "")
			}

			if !dispatcher.GetNextSyncRecord(addr) {
				break
			}
		}

		automate.ChangeState(name, addr, success)

	case success:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}
