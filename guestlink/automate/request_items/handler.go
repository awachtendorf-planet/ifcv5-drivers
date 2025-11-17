package requestitems

import (
	"strings"

	guestlink "github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	errInternalPacket     = errors.New("object packet is nil")
	errInternalJob        = errors.New("job failed")
	errNotLastSequence    = errors.New("was not the last sequence")
	errPmsUnavailable     = errors.New("PMS is not available")
	errPmsEmptyResponse   = errors.New("PMS send empty response")
	errPmsRejectedRequest = errors.New("PMS has rejected the request")
)

const (
	unknownError      = 0
	unknownCommand    = 1
	unknownRoom       = 2
	roomUnoccupied    = 3
	unknownAccount    = 4
	lockedFolio       = 10
	emptyGuestMessage = 12
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
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: nextAction}

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

		var errorCode int

		switch packet.Name {

		case template.PacketGuestMessageRequest: // MSGR -> MHDR + MCLR + MTXT
			errorCode, err = p.handleGuestMessageRequest(addr, packet)

		case template.PacketLookupRequest: // LOOK -> NAME
			errorCode, err = p.handleLookupRequest(addr, packet)

		case template.PacketStatusRequest: // STAT -> INFO
			errorCode, err = p.handleLookupRequest(addr, packet)

		case template.PacketDisplayRequest: // DISP -> NAME + ITEM + BAL
			errorCode, err = p.handleDisplayRequest(addr, packet)

		default:
			errorCode = unknownCommand
			err = errors.Errorf("missing handler for '%s'", packet.Name)
		}

		if err != nil {

			action.NextState = success

			if err != errPmsEmptyResponse && err != errPmsRejectedRequest { // suppress warning if not necessary
				action.NextState = shutdown
				log.Warn().Msgf("%s addr '%s' process packet '%s' failed, err=%s'", name, addr, packet.Name, err.Error())
			}

			err = p.sendCommandRefusion(addr, packet, &action, errorCode, packet.Tracking)

			if automate.Dispatcher().IsShutdown(err) {
				automate.NextAction(name, addr, shutdown, t)
			} else if err == errNotLastSequence {
				automate.NextAction(name, addr, success, t)
			}

			break
		}

		if err == nil {
			automate.NextAction(name, addr, nextAction, t)
			return
		}

	case nextAction:

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = errInternalPacket
			break
		}

		action.NextState = nextRecord
		if record, exist := p.getNextRecord(addr, packet.Name); exist {
			err = p.sendPacket(addr, &action, record)
		} else {
			automate.NextAction(name, addr, success, t)
			return
		}

	case nextRecord:

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = errInternalPacket
			break
		}

		pending := p.dropRecord(addr, packet.Name)
		if pending {
			automate.NextAction(name, addr, nextAction, t)
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

func (p *Plugin) decode(addr string, data string) string {
	dispatcher := p.automate.Dispatcher()
	encoding := dispatcher.GetEncoding(addr)
	if len(encoding) == 0 {
		return data
	}
	if dec, err := dispatcher.Decode([]byte(data), encoding); err == nil && len(dec) > 0 {
		return string(dec)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return data
}

func (p *Plugin) getRoomNumber(addr string, packet *ifc.LogicalPacket) (string, error) {
	roomNumber := p.getField(packet, "RoomNumber", true)
	roomNumber = p.decode(addr, roomNumber)
	if len(roomNumber) == 0 {
		return roomNumber, errors.Errorf("unknown room number")
	}
	return roomNumber, nil
}

func (p *Plugin) getAccountNumber(_ string, packet *ifc.LogicalPacket) (string, error) {
	accountNumber := p.getField(packet, "AccountNumber", true)
	accountNumber = strings.TrimLeft(accountNumber, "0")
	if len(accountNumber) == 0 {
		return accountNumber, errors.Errorf("unknown account number")
	}
	return accountNumber, nil
}

func (p *Plugin) getTransactionIdentifier(packet *ifc.LogicalPacket) string {
	driver := p.driver
	transaction := driver.GetTransaction(packet)
	return transaction.Identifier
}

func (p *Plugin) sendCommandRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, errorCode int, tracking string) error {
	driver := p.driver
	transaction := driver.GetTransaction(in)

	if transaction.IsLastSequence() {
		packet, err := driver.ConstructPacket(addr, template.PacketError, tracking, errorCode, transaction)
		if err == nil {
			err = p.automate.SendPacket(addr, packet, action)
		}
		return err
	}

	return errNotLastSequence
}

func (p *Plugin) sendPacket(addr string, action *dispatcher.StateAction, record nextaction) error {

	transaction := guestlink.Transaction{
		Identifier: record.identifier,
		Sequence:   record.sequence,
	}

	packet, err := p.driver.ConstructPacket(addr, record.TemplateName, record.CorrelationId, record.Context, transaction)
	if err == nil {
		if err = p.automate.SendPacket(addr, packet, action); err == nil {
			return nil
		}
	}

	return err
}
