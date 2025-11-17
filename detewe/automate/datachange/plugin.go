package datachange

import (
	"fmt"

	detewe "github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	dc "github.com/weareplanet/ifcv5-main/ifc/automate/generic"

	"github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/detewe/template"
)

// Plugin ...
type Plugin struct {
	*dc.Plugin
	driver *detewe.Dispatcher
}

// New return a new plugin
func New(parent *detewe.Dispatcher) *Plugin {

	p := &Plugin{
		dc.New(),
		parent,
	}

	p.Setup(dc.Config{

		Name: fmt.Sprintf("%T", p),

		WaitForReplyPacket: func(addr string, job *order.Job) bool { return true },

		InitHandler:     p.init,
		SendPacket:      p.send,
		ProcessWorkflow: p.processWorkflow,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	p.RegisterRule(template.PacketAnswerTG41, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{})
	p.RegisterRule(template.PacketAnswerTG67, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{})
	p.RegisterRule(template.PacketAnswerTG71, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{})
	p.RegisterRule(template.PacketAnswerTG60, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{})
	p.RegisterRule(template.PacketAnswerTG70, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{})
	p.RegisterRule(template.PacketAnswerTG80, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{})

}

func (p *Plugin) setWorkflow() {
	p.RegisterWorkflow(order.Checkin,
		template.PacketTG60,
		template.PacketTG41,
	)

	p.RegisterWorkflow(order.Checkout,
		template.PacketTG80,
		template.PacketTG67,
		template.PacketTG60,
		template.PacketTG41,
	)

	p.RegisterWorkflow(order.DataChange,
		template.PacketTG41,
	)

	p.RegisterWorkflow(order.RoomStatus,
		template.PacketTG60,
		template.PacketTG80,
	)

	p.RegisterWorkflow(order.WakeupRequest,
		template.PacketTG70,
		template.PacketTG71,
	)

	p.RegisterWorkflow(order.WakeupClear,
		template.PacketTG71,
	)

}

func (p *Plugin) processWorkflow(addr string, packetName string, job *order.Job) bool {

	station, _ := p.driver.GetStationAddr(addr)

	switch packetName {

	case template.PacketTG41:

		if !p.driver.GDXMode(station) && job.Action != order.DataChange {

			return false
		}

	case template.PacketTG67:

		if !p.driver.SendTG67(station) {

			return false
		}

	case template.PacketTG80:

		if !p.driver.SendTG80(station) && job.Action != order.RoomStatus {

			return false
		}

	case template.PacketTG70:

		if !p.driver.DirectWakeupMode(station) {

			return false
		}

	case template.PacketTG71:

		if p.driver.DirectWakeupMode(station) {

			return false
		}

	}
	return true
}

func (p *Plugin) handleReply(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	CMD := in.Data()["CMD"]
	result := in.Data()["Result"]

	expectedPacketName, _ := p.GetCurrentWorkflow(job)

	if in.Name != expectedPacketName {

		log.Warn().Msgf("%s addr '%s' unexpected answer '%s'", p.GetName(), addr, in.Name)
		p.CalculateNextTimeout(addr, action, job)
		return nil
	}

	station, _ := p.GetDispatcher().GetStationAddr(addr)

	expectedExtension := p.driver.GetExtension(station, job.Context)

	telegramExtension, ok := in.Data()["Participant"]
	if !ok || string(telegramExtension) != expectedExtension {

		log.Warn().Msgf("%s addr '%s' unexpected participant '%s' in '%s'", p.GetName(), addr, string(in.Data()["Participant"]), in.Name)
		p.CalculateNextTimeout(addr, action, job)
		return nil
	}

	if len(CMD) < 2 || len(result) == 0 {

		p.handleAnswerNegativ(addr, job, string(CMD)+" corrupted answer")
		p.ChangeState(addr, automatestate.Success)
		return nil
	}

	if (CMD[1] == 0x31 && CMD[0] != 0x37) || (CMD[1] == 0x31 && CMD[0] == 0x37 && job.Action == order.WakeupClear) {

		if result[0] != 0x30 {

			p.handleAnswerNegativ(addr, job, string(CMD)+string(result))
			p.ChangeState(addr, automatestate.Success)
			return nil
		}

	} else {

		if result[0] != 0x31 {

			p.handleAnswerNegativ(addr, job, string(CMD)+string(result))
			p.ChangeState(addr, automatestate.Success)
			return nil
		}
	}

	p.ChangeState(addr, automatestate.NextRecord)

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, job.Action, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}

func (p *Plugin) handleAnswerNegativ(addr string, job *order.Job, code string) {
	driver := p.driver
	reason := driver.GetAnswerText(code)
	p.HandleError(addr, job, reason)
}
