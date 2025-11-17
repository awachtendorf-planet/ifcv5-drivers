package datachange

import (
	"fmt"
	"strings"
	"time"

	definity "github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	dc "github.com/weareplanet/ifcv5-main/ifc/automate/generic"

	"github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/definity/template"
)

// Plugin ...
type Plugin struct {
	*dc.Plugin
	driver *definity.Dispatcher
}

// New return a new plugin
func New(parent *definity.Dispatcher) *Plugin {

	p := &Plugin{
		dc.New(),
		parent,
	}

	p.Setup(dc.Config{

		Name: fmt.Sprintf("%T", p),

		GetAnswerTimeout:   p.getAnswerTimeout,
		WaitForReplyPacket: p.waitForReplyPacket,

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,

		InitHandler: p.init,
		SendPacket:  p.send,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	p.RegisterRule(template.PacketCheckIn, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{NextState: automatestate.NextRecord})
	p.RegisterRule(template.PacketCheckOut, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{NextState: automatestate.NextRecord})
	p.RegisterRule(template.PacketDataChange, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{NextState: automatestate.NextRecord})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.Checkin,
		template.PacketCheckIn,
		template.PacketSwitchMessageLamp,
		template.PacketSetRestriction,
	)

	p.RegisterWorkflow(order.DataChange,
		template.PacketDataChange,
		template.PacketSwitchMessageLamp,
		template.PacketSetRestriction,
	)

	p.RegisterWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.RegisterWorkflow(order.RoomStatus,
		template.PacketSwitchMessageLamp,
		template.PacketSetRestriction,
	)

}

func (p *Plugin) waitForReplyPacket(_ string, job *order.Job) bool {

	if packet, exist := p.GetCurrentWorkflow(job); exist {

		if packet == template.PacketSwitchMessageLamp || packet == template.PacketSetRestriction {
			return false
		}

	}

	return true
}

func (p *Plugin) getAnswerTimeout(_ string, _ *order.Job) time.Duration {
	return definity.AnswerTimeout
}

func (p *Plugin) handleReply(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	if job == nil { // should not happend
		p.ChangeState(addr, action.NextState)
		return nil
	}

	name := p.GetName()

	// did the incoming packet match the workflow

	if packet, exist := p.GetCurrentWorkflow(job); !exist || packet != in.Name {

		if exist {
			log.Warn().Msgf("%s addr '%s' ignore packet, vendor reply with '%s', expected '%s'", name, addr, in.Name, packet)
		}

		p.CalculateNextTimeout(addr, action, job)
		return nil

	}

	driver := p.driver

	expectedRoom := driver.GetRoom(job.Context)

	room := p.getRoom(in)

	if room != expectedRoom {

		log.Warn().Msgf("%s addr '%s' ignore packet, vendor has answered for another room, expected '%s', got '%s'", name, addr, expectedRoom, room)

		p.CalculateNextTimeout(addr, action, job)
		return nil
	}

	proc := p.getProc(in)

	if proc < 0x1 {
		p.HandleError(addr, job, "vendor reply with unknown process code")
		p.ChangeState(addr, automatestate.Success)
		return nil
	}

	switch in.Name {

	case template.PacketCheckIn:

		if proc == 0x2 {
			p.HandleError(addr, job, "room already occupied")
			p.ChangeState(addr, automatestate.Success)
			return nil
		}

	case template.PacketDataChange:

		if proc == 0x3 {
			p.HandleError(addr, job, "room vacant")
			p.ChangeState(addr, automatestate.Success)
			return nil
		}

	}

	p.ChangeState(addr, action.NextState)

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, _ interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, job)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
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
