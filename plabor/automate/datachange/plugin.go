package datachange

import (
	"fmt"

	//"github.com/weareplanet/ifcv5-main/ifc/defines"
	plabor "github.com/weareplanet/ifcv5-drivers/plabor/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	dc "github.com/weareplanet/ifcv5-main/ifc/automate/generic"

	"github.com/weareplanet/ifcv5-drivers/plabor/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/plabor/template"

	// "github.com/spf13/cast"
	"github.com/pkg/errors"
)

// Plugin ...
type Plugin struct {
	*dc.Plugin
	driver *plabor.Dispatcher
}

// New return a new plugin
func New(parent *plabor.Dispatcher) *Plugin {

	p := &Plugin{
		dc.New(),
		parent,
	}

	p.Setup(dc.Config{

		Name: fmt.Sprintf("%T", p),

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,
		TemplateEnq: template.PacketEnq,
		TemplateEOT: template.PacketEot,

		InitHandler:     p.init,
		SendPacket:      p.send,
		ProcessWorkflow: p.processWorkflow,
		PreCheck:        p.preCheck,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.Checkin,
		template.PacketCheckIn,
	)

	p.RegisterWorkflow(order.DataChange,
		template.PacketDataChangeUpdate,
		template.PacketDataChangeMove,
	)

	p.RegisterWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.RegisterWorkflow(order.WakeupRequest,
		template.PacketWakeupIndirect,
		template.PacketWakeupIndirectk,
		template.PacketWakeupDirect,
		template.PacketWakeupDirectk,
	)

	p.RegisterWorkflow(order.GuestMessageOnline,
		template.PacketMessageSignal,
	)

	p.RegisterWorkflow(order.GuestMessageDelete,
		template.PacketMessageDelete,
	)

}

func (p *Plugin) preCheck(addr string, job *order.Job) (bool, error) {

	// database swap running
	if p.driver.GetSwapState(job.Station) {
		// p.HandleError(addr, job, "database swap running")
		return false, errors.New("database swap running")
	}
	return true, nil
}

func (p *Plugin) processWorkflow(addr string, packetName string, job *order.Job) bool {

	packet, _ := p.GetCurrentWorkflow(job)

	switch packet {

	case template.PacketWakeupIndirect:
		station, _ := p.driver.GetStationAddr(addr)
		if !p.driver.WakeUpDateWithYear(station) {
			return false
		}
		if p.driver.DirectWakeupMode(station) {
			return false
		}

	case template.PacketWakeupIndirectk:
		station, _ := p.driver.GetStationAddr(addr)
		if p.driver.WakeUpDateWithYear(station) {
			return false
		}
		if p.driver.DirectWakeupMode(station) {
			return false
		}

	case template.PacketWakeupDirect:
		station, _ := p.driver.GetStationAddr(addr)
		if !p.driver.WakeUpDateWithYear(station) {
			return false
		}
		if !p.driver.DirectWakeupMode(station) {
			return false
		}

	case template.PacketWakeupDirectk:
		station, _ := p.driver.GetStationAddr(addr)
		if p.driver.WakeUpDateWithYear(station) {
			return false
		}
		if !p.driver.DirectWakeupMode(station) {
			return false
		}

	case template.PacketDataChangeUpdate:

		return !p.driver.IsMove(job.Context)

	case template.PacketDataChangeMove:

		return p.driver.IsMove(job.Context)

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
