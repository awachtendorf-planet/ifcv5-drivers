package datachange

import (
	"github.com/weareplanet/ifcv5-drivers/centigram/template"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
)

func (p *Plugin) setWorkflow() {

	p.registerWorkflow(order.Checkin,
		template.PacketCheckIn,
		template.PacketChangeName,
		template.PacketChangeFCOS,
		template.PacketChangeVoiceMail,
	)

	p.registerWorkflow(order.DataChange,
		template.PacketRoomChange, // optional if room changed
		template.PacketChangeName,
		template.PacketChangeFCOS,
		template.PacketChangeVoiceMail,
	)

	p.registerWorkflow(order.Checkout,
		template.PacketCheckOut,
	)

	p.registerWorkflow(order.RoomStatus,
		template.PacketChangeVoiceMail,
	)

}

func (p *Plugin) registerWorkflow(scope order.Action, templates ...string) {

	count := uint64(1) // Achtung, Zählung beginnt in diesem Automaten bei 1, da der Job.Task Zähler bereits im ersten Durchlauf erhöht wird

	for i := range templates {
		p.workflow.Set(int(scope), count, templates[i])
		count++
	}

}
