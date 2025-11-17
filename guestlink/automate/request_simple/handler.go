package requestsimple

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	errInternalPacket  = errors.New("object packet is nil")
	errInternalJob     = errors.New("job failed")
	errNotLastSequence = errors.New("was not the last sequence")
)

const (
	unknownError   = 0
	unknownCommand = 1
	unknownRoom    = 2
	unknownAccount = 4
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

		job.Tries++

		var errorCode int
		var retry = job.Tries > 1

		switch packet.Name {

		case template.PacketRoomStatus:
			errorCode, err = p.handleRoomStatus(addr, packet, retry)

		case template.PacketGuestMessageRead:
			errorCode, err = p.handleGuestMessageRead(addr, packet, retry)

		case template.PacketWakeupSet:
			errorCode, err = p.handleWakeupSet(addr, packet, retry)

		case template.PacketWakeupClear:
			errorCode, err = p.handleWakeupClear(addr, packet, retry)

		case template.PacketWakeupResult:
			errorCode, err = p.handleWakeupResult(addr, packet, retry)

		case template.PacketPostCharge:
			errorCode, err = p.handlePostCharge(addr, packet, retry)

		case template.PacketInit:
			errorCode, err = p.handleInitRequest(addr, packet, retry)

		case template.PacketCheckoutRequest, template.PacketUnknownCommand:
			errorCode = unknownCommand
			err = errors.New("not supported")

		default:
			errorCode = unknownCommand
			err = errors.Errorf("missing handler for '%s'", packet.Name)
		}

		if err != nil {
			log.Error().Msgf("%s addr '%s' process packet '%s' failed, err=%s", name, addr, packet.Name, err.Error())
			err = p.sendCommandRefusion(addr, packet, &action, errorCode)
		} else {
			err = p.sendCommandAcknowledge(addr, packet, &action)
		}

		if automate.Dispatcher().IsShutdown(err) {
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		// we did not send command acknowledge or refusion, because it was not the last sequence number
		// change automate state to success
		if err == errNotLastSequence {
			automate.NextAction(name, addr, success, t)
			return
		}

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

func (p *Plugin) getTime(packet *ifc.LogicalPacket) (time.Time, error) {
	var layout string

	dt := p.getField(packet, "Date", false) // MMDDYY
	ti := p.getField(packet, "Time", false) // HHMM or HHMMSS

	date := dt + ti

	if len(date) == 0 {
		errors.New("empty date object")
	}

	switch len(date) {
	case 10: // MMDDYYHHMM
		layout = "0102061504"

	case 12: // MMDDYYHHMMSS
		layout = "010206150405"

	}

	constructed, err := time.Parse(layout, date)
	return constructed, err
}

func (p *Plugin) sendCommandAcknowledge(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction) error {
	driver := p.driver
	transaction := driver.GetTransaction(in)
	if transaction.IsLastSequence() {
		packet, err := driver.ConstructPacket(addr, template.PacketVerify, "", nil, transaction)
		if err == nil {
			err = p.automate.SendPacket(addr, packet, action)
		}
		return err
	}
	return errNotLastSequence
}

func (p *Plugin) sendCommandRefusion(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, errorCode int) error {
	driver := p.driver
	transaction := driver.GetTransaction(in)
	if transaction.IsLastSequence() {
		packet, err := driver.ConstructPacket(addr, template.PacketError, "", errorCode, transaction)
		if err == nil {
			err = p.automate.SendPacket(addr, packet, action)
		}
		return err
	}
	return errNotLastSequence
}
