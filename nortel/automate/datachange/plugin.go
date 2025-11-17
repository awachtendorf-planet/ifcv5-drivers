package datachange

import (
	"fmt"
	"time"

	nortel "github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	dc "github.com/weareplanet/ifcv5-main/ifc/automate/generic"

	"github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/nortel/template"

	"github.com/pkg/errors"
)

// Plugin ...
type Plugin struct {
	*dc.Plugin
	driver *nortel.Dispatcher
}

// New return a new plugin
func New(parent *nortel.Dispatcher) *Plugin {

	p := &Plugin{
		dc.New(),
		parent,
	}

	p.Setup(dc.Config{

		Name: fmt.Sprintf("%T", p),

		GetAnswerTimeout:         p.getAnswerTimeout,
		GetAcknowledgmentTimeout: p.getAcknowledgmentTimeout,
		WaitForReplyPacket:       p.waitForReplyPacket,
		HandleNoAnswer:           p.handleNoAnswer,

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,

		InitHandler:     p.init,
		SendPacket:      p.send,
		ProcessWorkflow: p.processWorkflow,
		PreCheck:        p.preCheck,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	p.RegisterRule(template.PacketError, automatestate.WaitForAnswer, p.handleErrorReply, dispatcher.StateAction{})
	p.RegisterRule(template.PacketError, automatestate.CommandSent, p.handleErrorReply, dispatcher.StateAction{})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.Checkin,
		template.PacketCheckInExtension,
		template.PacketDisplayName,
		template.PacketLanguage,
		template.PacketVipState,
	)

	p.RegisterWorkflow(order.DataChange,
		template.PacketDisplayName,
		template.PacketLanguage,
		template.PacketVipState,
	)

	p.RegisterWorkflow(order.Checkout,
		template.PacketDisplayName,
		template.PacketMessageLamp,
		template.PacketCheckOutExtension,
	)

	p.RegisterWorkflow(order.RoomStatus,
		template.PacketMessageLamp,
		template.PacketDoNotDisturb,
		// template.PacketSetCCRS,
		// template.PacketSetECC1,
		// template.PacketSetECC2,
		template.PacketClassOfService,
	)

	p.RegisterWorkflow(order.WakeupRequest,
		template.PacketWakeupSet,
	)

	p.RegisterWorkflow(order.WakeupClear,
		template.PacketWakeupClear,
	)

}

func (p *Plugin) getAcknowledgmentTimeout(addr string) time.Duration {
	driver := p.driver
	if driver.IsBackgroundTerminalMode(addr) {
		return driver.GetBackgroundTerminalTimeout(addr)
	}
	return driver.GetPacketTimeout(addr)
}

func (p *Plugin) getAnswerTimeout(addr string, _ *order.Job) time.Duration {
	// wait a second for possible error reply
	// todo: background terminal handling?
	driver := p.driver
	return driver.GetAnswerTimeout(addr)
}

func (p *Plugin) waitForReplyPacket(_ string, _ *order.Job) bool {
	// wait for a second for possible error reply
	// todo: background terminal handling?
	return true
}

func (p *Plugin) preCheck(_ string, job *order.Job) (bool, error) {
	if p.driver.IsSharer(job.Context) {
		return true, errors.New("ignore packet, because of sharer flag is set")
	}
	return true, nil
}

func (p *Plugin) processWorkflow(addr string, _ string, job *order.Job) bool {

	packet, _ := p.GetCurrentWorkflow(job)

	switch packet {

	case template.PacketDisplayName:
		if job.Action == order.Checkin || job.Action == order.DataChange {
			station, _ := p.driver.GetStationAddr(addr)
			if !p.driver.SendGuestName(station) {
				return false
			}
		}

	case template.PacketLanguage:
		station, _ := p.driver.GetStationAddr(addr)
		if !p.driver.SendGuestLanguage(station) {
			return false
		}

	case template.PacketVipState:
		station, _ := p.driver.GetStationAddr(addr)
		if !p.driver.SendGuestVIPState(station) {
			return false
		}

	case template.PacketClassOfService:
		station, _ := p.driver.GetStationAddr(addr)
		if !p.driver.SendClassOfService(station) {
			return false
		}

		/*
			case template.PacketSetCCRS, template.PacketSetECC1, template.PacketSetECC2:
				station, _ := p.driver.GetStationAddr(addr)
				if !p.driver.SendClassOfService(station) {
					return false
				}
		*/

	}

	return true
}

func (p *Plugin) handleNoAnswer(_ string, _ *order.Job) bool {
	// no answer mean no error reply, continue with workflow
	return true
}

func (p *Plugin) handleErrorReply(addr string, in *ifc.LogicalPacket, _ *dispatcher.StateAction, job *order.Job) error {

	msg := string(in.Data()["Msg"])

	switch msg {

	case "MNEMONIC":
		msg = "bad keyword"
		// continue with workflow ?

	case "NAME BIG":
		msg = "display name exceeds maximum length"

	case "DUPLICATE":
		msg = "mote then one set of quotation marks found"

	case "INPUT ERROR":
		packet, _ := p.GetCurrentWorkflow(job)
		if packet == template.PacketDisplayName {
			msg = "field value exceeds the maximum allowed name length"
		} else {
			msg = "error in data field"
		}

	case "NO DATA FOUND":
		msg = "specified extension is not a room phone"

	case "NO SET CPND DATA":
		msg = "specified extension is not set for CPND"

	case "NO CUST CPND DATA":
		msg = "the customer level CPND data block is not configured"

	case "NO CPND MEMORY":
		msg = "no CPND memory can be allocated for the given extension"

	default:
		msg = "unspecified error"
	}

	p.HandleError(addr, job, msg)
	p.ChangeState(addr, automatestate.Shutdown)

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context, job.Action)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
