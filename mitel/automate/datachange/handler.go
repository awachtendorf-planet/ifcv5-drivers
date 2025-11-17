package datachange

import (
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/mitel/template"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	if !p.linkState(addr) {
		p.handleError(addr, job, "vendor disconnected")
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: nextRecord}

	switch state {

	case busy:

		if job == nil || !job.InProcess() { // order was canceled externally
			err = defines.ErrorJobFailed
			break
		}

		if job.Context == nil { // job context failed
			err = defines.ErrorJobContext
			break
		}

		// synchronization between datachange automate and linkcontrol automate
		if !p.blockCommunication(addr) {

			if job.Timestamp == 0 {
				job.Timestamp = time.Now().UnixNano()
			}

			elapsed := time.Since(time.Unix(0, job.Timestamp))
			if elapsed >= 10*time.Second {
				p.handleError(addr, job, "communication device is blocked")
				automate.NextAction(name, addr, shutdown, t)
				return
			}

			action.NextTimeout = 1 * time.Second
			log.Debug().Msgf("%s addr '%s' communication blocked, retry in %s (job runtime '%s')", name, addr, action.NextTimeout, elapsed)
			break
		}

		action.NextState = nextAction
		err = p.sendPacket(addr, &action, template.PacketEnq, "", nil)

	case nextAction:

		if job == nil || !job.InProcess() { // order was canceled externally
			err = defines.ErrorJobFailed
			break
		}

		if job.Context == nil { // job context failed
			err = defines.ErrorJobContext
			break
		}

		action.CurrentState = busy                        // on error fallback to send enquiry
		automate.SetErrorCounter(addr, uint16(job.Tries)) // error counter künstlich erzeugen, handling für ENQ -> ACK, Packet -> NAK

		// send packets defined at workflow
		if packet, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {
			err = p.sendPacket(addr, &action, packet, job.Initiator, job.Context)
		} else {
			// should not happend here
			p.handleSuccess(addr, job)
			automate.NextAction(name, addr, success, t)
			return
		}

		if err == nil {
			job.Tries++
		}

	case nextRecord:

		dispatcher := automate.Dispatcher()
		dispatcher.SetAlive(addr)

		job.Tries = 0
		job.Task++

		// exist more packets to send
		if packetName, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {

			// Achtung: Prüfung funktioniert nicht für den ersten Eintrag im Workflow, da der Automate im Busy State startet und danach direkt NextAction ausführt.
			// Sollte die Prüfung für den ersten Eintrag benötigt werden, muss der Automate im NextRecord State starten. Siehe dbsync.

			switch packetName {

			case template.PacketSetRestriction:

				if !p.driver.SendRestrictionRecord(addr) {
					automate.NextAction(name, addr, nextRecord, t)
					return
				}

			}

			automate.NextAction(name, addr, busy, t)
			return
		}

		// job done
		p.handleSuccess(addr, job)
		automate.NextAction(name, addr, success, t)
		return

	case success, shutdown:

		t.Tick()
		return

	default:

		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)

		if err == defines.ErrorJobContext || err == defines.ErrorJobFailed {
			p.handleError(addr, job, "job was aborted")
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		internalErrorCounter := automate.GetInternalErrorCounter(addr)
		if internalErrorCounter > 0 {
			p.handleError(addr, job, err.Error())
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
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
