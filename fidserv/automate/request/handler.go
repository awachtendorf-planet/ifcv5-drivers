package request

import (
	"github.com/pkg/errors"
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"
)

var (
	internalPacketError = errors.New("object packet is nil")
	internalJobError    = errors.New("job failed")
	pmsUnavailable      = errors.New("PMS is not available")
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
			err = internalJobError
			break
		}

		if job.Context == nil {
			err = internalPacketError
			break
		}

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = internalPacketError
			break
		}

		pendingRecords := false

		switch packet.Name {

		case template.PacketPostingSimple:
			err = p.handleSimplePosting(addr, packet, &action)

		case template.PacketPostingRequest:
			driver := p.driver
			if _, exist := driver.GetField(packet, "PI"); exist {
				err = p.handlePostingInquiry(addr, packet, &action)
			} else {
				err = p.handlePostingRequest(addr, packet, &action)
			}

		case template.PacketWakeupAnswer, template.PacketWakeupRequest, template.PacketWakeupClear:
			err = p.handleWakeup(addr, packet, &action)

		case template.PacketRoomData:
			err = p.handleRoomStatus(addr, packet, &action)

		case template.PacketRemoteCheckOut:
			err = p.handleRemoteCheckout(addr, packet, &action)

		case template.PacketGuestCheckDetails:
			err = p.handleGuestCheckDetails(addr, packet, &action)

		case template.PacketGuestMessageDelete:
			err = p.handleGuestMessageDelete(addr, packet, &action)

		case template.PacketGuestMessageRequest:
			pendingRecords, err = p.handleGuestMessageRequest(addr, packet, &action)

		case template.PacketGuestBillRequest:
			pendingRecords, err = p.handleBillPreview(addr, packet, &action)

		default:
			log.Warn().Msgf("%s addr '%s' unknown packet '%s'", name, addr, packet.Name)
		}

		if automate.Dispatcher().IsShutdown(err) {
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		if pendingRecords {
			automate.NextAction(name, addr, nextAction, t)
		} else {
			automate.NextAction(name, addr, success, t)
		}
		return

	case nextAction:

		if job == nil || !job.InProcess() {
			err = internalJobError
			break
		}

		if job.Context == nil {
			err = internalPacketError
			break
		}

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = internalPacketError
			break
		}

		pendingRecords := false

		switch packet.Name {

		case template.PacketGuestMessageRequest:
			pendingRecords, err = p.handleGuestMessageItem(addr, packet, &action)

		case template.PacketGuestBillRequest:
			pendingRecords, err = p.handleBillInvoice(addr, packet, &action)

		default:
			automate.NextAction(name, addr, success, t)

		}

		if pendingRecords {
			automate.NextAction(name, addr, nextAction, t)
		} else {
			automate.NextAction(name, addr, success, t)
		}

		return

	case success:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		if err == internalPacketError || err == internalJobError {
			automate.NextAction(name, addr, shutdown, t)
			return
		}
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}
