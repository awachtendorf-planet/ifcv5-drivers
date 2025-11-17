package datachange

import (
	"fmt"
	"time"

	//"github.com/weareplanet/ifcv5-main/ifc/defines"
	informatel "github.com/weareplanet/ifcv5-drivers/informatel/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	// "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"

	dc "github.com/weareplanet/ifcv5-main/ifc/automate/generic"

	"github.com/weareplanet/ifcv5-drivers/informatel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/informatel/template"
	// "github.com/spf13/cast"
	// "github.com/pkg/errors"
)

// Plugin ...
type Plugin struct {
	*dc.Plugin
	driver *informatel.Dispatcher
}

// New return a new plugin
func New(parent *informatel.Dispatcher) *Plugin {

	p := &Plugin{
		dc.New(),
		parent,
	}

	p.Setup(dc.Config{

		Name: fmt.Sprintf("%T", p),

		OrderType:          order.ASW,
		AnswerTimeout:      8 * time.Second,
		GetAnswerTimeout:   p.getAnswerTimeout,
		WaitForReplyPacket: p.waitForReplyPacket,

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,
		//TemplateEnq: template.PacketEnq,

		InitHandler:     p.init,
		SendPacket:      p.send,
		ProcessWorkflow: p.processWorkflow,
		PreCheck:        p.preCheck,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	// p.RegisterRule(template.PacketAck, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{NextState: automatestate.NextRecord})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.Checkin,
		template.PacketCheckIn,
	)

	p.RegisterWorkflow(order.DataChange,
		template.PacketRoomMove,
		template.PacketNameChange,
	)

	p.RegisterWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.RegisterWorkflow(order.WakeupRequest,
		template.PacketWakeupOutgoingCall,
	)

	p.RegisterWorkflow(order.WakeupClear,
		template.PacketWakeupOutgoingCall,
	)

	p.RegisterWorkflow(order.RoomStatus,
		template.PacketClassOfService,
		template.PacketMessageLightStatus,
		template.PacketDND,
	)

}

func (p *Plugin) waitForReplyPacket(addr string, job *order.Job) bool {

	// station, _ := p.driver.GetStationAddr(addr)
	// state := p.driver.GetConfig(station, "WaitForReplyPacket", "false")

	// return cast.ToBool(state)

	return false
}

func (p *Plugin) getAnswerTimeout(addr string, job *order.Job) time.Duration {

	// if timeout, ok := p.GetSetting(addr, defines.KeyAnswerTimeout, uint(0)).(uint64); ok && timeout > 0 {
	// 	if duration := cast.ToDuration(time.Duration(timeout) * time.Second); duration > 0 {
	// 		return duration
	// 	}
	// }

	return 15 * time.Second
}

func (p *Plugin) preCheck(addr string, job *order.Job) (bool, error) {
	// return errors.New("job pre check failed")
	return true, nil
}

func (p *Plugin) processWorkflow(addr string, packetName string, job *order.Job) bool {

	if packetName == template.PacketRoomMove {

		return p.driver.IsMove(job.Context)
	}

	return true
}

func (p *Plugin) handleReply(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	// success -> auto pms reply
	p.ChangeState(addr, action.NextState)

	// error -> auto pms reply
	// return errors.New("error reason")

	// error -> create pms reply
	// p.HandleError(addr, job, "error reason")
	// p.ChangeState(addr, automatestate.Success)

	// ignore packet, re-calculate answer timeout
	// p.CalculateNextTimeout(addr, action, job)

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
