package dbsync

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

	p.registerWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.registerWorkflow(order.WakeupRequest,
		template.PacketWakeupSet,
	)

	p.registerWorkflow(order.WakeupClear,
		template.PacketWakeupClear,
	)

	p.registerWorkflow(order.NightAuditStart,
		template.PacketSwapStart,
	)

	p.registerWorkflow(order.NightAuditEnd,
		template.PacketSwapEnd,
	)

}

func (p *Plugin) registerWorkflow(scope order.Action, templates ...string) {

	count := uint64(1) // Achtung, Zählung beginnt in diesem Automaten bei 1, da der Job.Task Zähler bereits im ersten Durchlauf erhöht wird

	for i := range templates {
		p.workflow.Set(int(scope), count, templates[i])
		count++
	}

}
