package datachange

import (
	"github.com/weareplanet/ifcv5-drivers/bartech/template"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
)

func (p *Plugin) setWorkflow() {

	p.registerWorkflow(order.Checkin,
		template.PacketCheckIn,
	)

	p.registerWorkflow(order.DataChange,
		template.PacketCheckIn,
	)

	p.registerWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.registerWorkflow(order.RoomStatus,
		template.PacketLockBar, // one of them is filtered out
		template.PacketUnlockBar,
	)

	p.registerWorkflow(order.NightAuditEnd,
		template.PacketEndOfDay,
	)

}

func (p *Plugin) registerWorkflow(scope order.Action, templates ...string) {

	count := uint64(1) // Achtung, Zählung beginnt in diesem Automaten bei 1, da der Job.Task Zähler bereits im ersten Durchlauf erhöht wird

	for i := range templates {
		p.workflow.Set(int(scope), count, templates[i])
		count++
	}

}
