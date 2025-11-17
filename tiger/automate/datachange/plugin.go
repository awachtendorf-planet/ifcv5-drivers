package datachange

import (
	"fmt"

	//"github.com/weareplanet/ifcv5-main/ifc/defines"
	tiger "github.com/weareplanet/ifcv5-drivers/tiger/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	dc "github.com/weareplanet/ifcv5-main/ifc/automate/generic"

	"github.com/weareplanet/ifcv5-drivers/tiger/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/tiger/template"
	// "github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*dc.Plugin
	driver *tiger.Dispatcher
}

// New return a new plugin
func New(parent *tiger.Dispatcher) *Plugin {

	p := &Plugin{
		dc.New(),
		parent,
	}

	p.Setup(dc.Config{

		Name: fmt.Sprintf("%T", p),

		HandleNoAnswer: p.handleNoAnswer,

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,
		TemplateEnq: template.PacketEnq,

		InitHandler:     p.init,
		SendPacket:      p.send,
		ProcessWorkflow: p.processWorkflow,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.Checkin,
		template.PacketRoomBasedCheckIn,
		template.PacketAdditionalGuest,
		template.PacketVideoRights,
	)

	p.RegisterWorkflow(order.Checkout,
		template.PacketRoomBasedCheckOut,
	)

	p.RegisterWorkflow(order.DataChange,
		template.PacketRoomtransfer,
		template.PacketInformationUpdate,
	)

	p.RegisterWorkflow(order.RoomStatus,
		template.PacketDND,
	)

	p.RegisterWorkflow(order.WakeupRequest,
		template.PacketWakeUpSet,
	)

	p.RegisterWorkflow(order.WakeupClear,
		template.PacketWakeUpClear,
	)

	p.RegisterWorkflow(order.GuestMessageOnline,
		template.PacketMessageWaitingGuest,
	)

}

func (p *Plugin) handleNoAnswer(_ string, _ *order.Job) bool {
	// no answer mean no error reply, continue with workflow
	return true
}

func (p *Plugin) processWorkflow(addr string, packetName string, job *order.Job) bool {

	packet, _ := p.GetCurrentWorkflow(job)

	switch packet {

	case template.PacketRoomBasedCheckIn:

		guest, ok := job.Context.(*record.Guest)
		station, _ := p.driver.GetStationAddr(addr)

		return ok && !guest.Reservation.SharedInd || p.driver.Protocol(station) < 3

	case template.PacketVideoRights:

		station, _ := p.driver.GetStationAddr(addr)
		return p.driver.Protocol(station) == 3

	case template.PacketAdditionalGuest:

		guest, ok := job.Context.(*record.Guest)

		station, _ := p.driver.GetStationAddr(addr)
		return ok && guest.Reservation.SharedInd && p.driver.Protocol(station) == 3

	case template.PacketInformationUpdate:

		return !p.driver.IsMove(job.Context)

	case template.PacketRoomtransfer:

		station, _ := p.driver.GetStationAddr(addr)
		return p.driver.IsMove(job.Context) && p.driver.Protocol(station) > 2

	}

	return true
}

func (p *Plugin) handleReply(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	// success -> auto pms reply
	p.ChangeState(addr, action.NextState)

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)

	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
