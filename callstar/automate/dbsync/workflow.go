package dbsync

import (
	"github.com/weareplanet/ifcv5-drivers/callstar/template"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
)

func (p *Plugin) setWorkflow() {

	p.registerWorkflow(order.Checkin,
		template.PacketCheckInSwap,
	)

	p.registerWorkflow(order.Checkout,
		template.PacketCheckOutSwap,
	)

	p.registerWorkflow(order.WakeupRequest,
		template.PacketWakeupRequestSwap,
	)

	p.registerWorkflow(order.RoomStatus,
		template.PacketRoomStatusSwap,
	)

}

func (p *Plugin) registerWorkflow(scope order.Action, templates ...string) {

	count := uint64(1) // Achtung, Zählung beginnt in diesem Automaten bei 1, da der Job.Task Zähler bereits im ersten Durchlauf erhöht wird

	for i := range templates {
		p.workflow.Set(int(scope), count, templates[i])
		count++
	}

}
