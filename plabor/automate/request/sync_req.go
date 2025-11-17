package request

import (
	syncstate "github.com/weareplanet/ifcv5-drivers/plabor/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

func (p *Plugin) handleDBSyncReq(addr string, packet *ifc.LogicalPacket) {

	broker := p.driver.Broker()
	if broker != nil {
		broker.Broadcast(syncstate.NewEvent(addr, syncstate.Start, nil), syncstate.Start.String())
	}

}
