package datachange

import (
	"github.com/weareplanet/ifcv5-drivers/mitel/template"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
)

func (p *Plugin) setWorkflow() {

	p.registerWorkflow(order.Checkin,
		template.PacketCheckIn,
		template.PacketChangeName,
		template.PacketMessageLamp,
		template.PacketSetRestriction,
	)

	p.registerWorkflow(order.DataChange,
		template.PacketChangeName,
		template.PacketMessageLamp,
		template.PacketSetRestriction,
	)

	p.registerWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.registerWorkflow(order.WakeupRequest,
		template.PacketWakeupSet,
	)

	p.registerWorkflow(order.WakeupClear,
		template.PacketWakeupClear,
	)

	p.registerWorkflow(order.RoomStatus,
		template.PacketMessageLamp,
		template.PacketSetRestriction,
	)

}

func (p *Plugin) registerWorkflow(scope order.Action, templates ...string) {

	count := uint64(0)

	for i := range templates {
		p.workflow.Set(int(scope), count, templates[i])
		count++
	}

}
