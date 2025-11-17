package dbsync

import (
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/robobar/template"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, Job *order.Job) {

	automate := p.automate
	name := automate.Name()

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case resyncStart:
		automate.ChangeState(name, addr, resyncRecord)

	case resyncRecord:
		// send record
		action.NextState = resyncNextRecord
		err = p.sendRecord(addr, &action)

	case resyncNextRecord:
		// get next record
		// next state -> resyncRecord or resyncEnd if no more records
		action.NextState = resyncRecord
		err = p.getNextRecord(addr, &action)

	case resyncEnd:
		automate.ChangeState(name, addr, success)

	case busy:
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

func (p *Plugin) sendRecord(addr string, action *dispatcher.StateAction) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := p.driver
	name := automate.Name()

	storeJob := func(addr string, job *order.Job) {
		p.guard.Lock()
		p.previousJob[addr] = job
		p.guard.Unlock()
	}

	job := dispatcher.GetSyncRecord(addr)
	if job != nil {

		if err := driver.PreCheck(job); err == nil {

			switch job.Action {

			case order.Checkin:
				err = p.sendPacket(addr, action, template.PacketUpdateCheckIn, job.Initiator, job.Context)

			case order.Checkout:
				err = p.sendPacket(addr, action, template.PacketUpdateCheckOut, job.Initiator, job.Context)

			case order.RoomStatus:
				mb := driver.GetMinibarRight(job.Context)

				switch mb {

				case 0:
					err = p.sendPacket(addr, action, template.PacketUpdateLockBar, job.Initiator, job.Context)

				case 2:
					err = p.sendPacket(addr, action, template.PacketUpdateUnlockBar, job.Initiator, job.Context)

					// default: save with because of PreCheck
				}

			default:
				log.Debug().Msgf("%s addr '%s' drop record '%s'", name, addr, job.Action.String())
				return nil

			}

			storeJob(addr, job)
			return err

		} else {

			log.Warn().Msgf("%s addr '%s' ignore record '%s', err=%s", name, addr, job.Action.String(), err)

		}

		if action != nil {
			action.NextState = resyncNextRecord
			automate.ChangeState(name, addr, action.NextState)
		}
		return nil
	}

	if action != nil {
		action.NextState = resyncEnd
		automate.ChangeState(name, addr, action.NextState)
	}
	return nil
}

func (p *Plugin) getNextRecord(addr string, action *dispatcher.StateAction) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	if dispatcher.GetNextSyncRecord(addr) {
		if action != nil {
			automate.ChangeState(name, addr, action.NextState)
		}
		return nil
	}

	if action != nil {
		action.NextState = resyncEnd
		automate.ChangeState(name, addr, action.NextState)
	}
	return nil
}

func (p *Plugin) sendPacket(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}) error {
	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		p.automate.AddInternalErrorCounter(addr)
		return err
	}
	err = p.automate.SendPacket(addr, packet, action)
	return err
}
