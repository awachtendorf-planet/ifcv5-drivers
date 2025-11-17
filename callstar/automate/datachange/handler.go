package datachange

import (
	"fmt"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/callstar/template"
)

const (
	maxRefusion = 3
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// link state
	if !p.linkState(addr) {
		p.handleError(addr, job, "vendor disconnected")
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	// database swap running
	if p.driver.GetSwapState(addr) {
		p.handleError(addr, job, "database swap running")
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: nextRecord}

	switch state {

	case busy:

		if job != nil && job.Tries >= maxRefusion { // timeout: check record error counter
			p.handleError(addr, job, "vendor did not answer")
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		action.NextState = nextAction
		err = p.sendPacket(addr, &action, template.PacketEnq, "", nil)

	case nextAction:

		if job == nil { // should not happend, safty check
			automate.NextAction(name, addr, success, t)
			return
		}

		action.CurrentState = busy // on error fallback to send enquiry

		// send packets defined at workflow
		if packet, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {
			err = p.sendPacket(addr, &action, packet, job.Initiator, job.Context)
		} else {
			// should not happend, safty check
			automate.NextAction(name, addr, nextRecord, t)
			return
		}

		if err == nil {
			job.Tries++
		}

	case nextRecord:

		if job != nil {

			job.Tries = 0
			job.Task++

			// exist more packets to send
			if _, exist := p.workflow.Get(int(job.Action), uint64(job.Task)); exist {
				automate.NextAction(name, addr, busy, t) // send enquire
				return
			}
		}

		// done
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

func (p *Plugin) handlePacketRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	automate := p.automate
	name := automate.Name()

	if job != nil && job.Tries >= maxRefusion {
		p.handleError(addr, job, fmt.Sprintf("vendor refused the record %d times", maxRefusion))
		if action != nil {
			action.NextTimeout = nextActionDelay
		}
		automate.ChangeState(name, addr, shutdown)
		return nil
	}

	return p.automate.HandlePacketRefusion(addr, in, action, job)
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
