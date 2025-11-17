package dbsync

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/definity/template"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, _ *order.Job) {

	automate := p.automate
	name := automate.Name()

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case resyncStart:

		// send sync start

		action.NextState = resyncRecord
		err = p.sendPacket(addr, &action, template.PacketSyncStart, "", nil)

	case resyncRecord:

		// send record

		// action.NextState = resyncNextRecord // if no answer is expected

		action.NextState = answer
		action.NextTimeout = answerTimeout

		err = p.sendRecord(addr, &action)

	case resyncNextRecord:

		// get next record, next state -> resyncRecord or resyncEnd if no more records

		action.NextState = resyncRecord
		err = p.getNextRecord(addr, &action)

	case resyncEnd:

		// send sync end

		action.NextState = success
		err = p.sendPacket(addr, &action, template.PacketSyncEnd, "", nil)

	case busy:

		automate.ChangeState(name, addr, success)

	case answer:

		log.Warn().Msgf("%s addr '%s' vendor did not answer, send next sync record", name, addr)
		automate.ChangeState(name, addr, resyncNextRecord)
		t.Tick()
		return

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

	job := dispatcher.GetSyncRecord(addr)

	if job != nil {

		if err := driver.PreCheck(job); err == nil {

			switch job.Action {

			case order.Checkin, order.Checkout:

				err = p.sendPacket(addr, action, template.PacketRoomDataImageSwap, "", job)
				return err

			default:

				log.Debug().Msgf("%s addr '%s' drop record '%s'", name, addr, job.Action.String())

			}

		} else {

			log.Warn().Msgf("%s addr '%s' ignore record '%s', err=%s", name, addr, job.Action.String(), err)

		}

		if action != nil {
			action.NextState = resyncNextRecord
			action.NextTimeout = nextActionDelay
			automate.ChangeState(name, addr, action.NextState)
		}

		return nil
	}

	// no more sync records

	if action != nil {
		action.NextState = resyncEnd
		action.NextTimeout = nextActionDelay
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

func (p *Plugin) handleReply(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, _ *order.Job) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := p.driver
	name := p.automate.Name()

	job := dispatcher.GetSyncRecord(addr)

	if job == nil { // should not happend
		log.Error().Msgf("%s addr '%s' can not find sync record", name, addr)
		automate.ChangeState(name, addr, action.NextState)
		return nil
	}

	expectedRoom := driver.GetRoom(job.Context)

	room := p.getRoom(packet)

	if room != expectedRoom {
		log.Error().Msgf("%s addr '%s' vendor has answered for another room, expected '%s', got '%s'", name, addr, expectedRoom, room)
	}

	automate.ChangeState(name, addr, action.NextState)

	return nil
}

// todo: create dispatcher function for transparent mode

func (p *Plugin) getString(packet *ifc.LogicalPacket, field string) string {

	if packet == nil || len(field) == 0 {
		return ""
	}

	value := string(packet.Data()[field])
	value = strings.TrimLeft(value, "0")

	return value
}

func (p *Plugin) getRoom(packet *ifc.LogicalPacket) string {
	return p.getString(packet, "RSN")
}

func (p *Plugin) getProc(packet *ifc.LogicalPacket) byte {

	if packet == nil {
		return 0
	}

	proc := string(packet.Data()["PROC"])

	if len(proc) > 0 {
		return proc[0]
	}

	return 0

}
