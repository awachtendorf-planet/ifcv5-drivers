package request

import (
	"fmt"

	syncstate "github.com/weareplanet/ifcv5-drivers/informatel/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

func (p *Plugin) handleDBSyncReq(addr string, packet *ifc.LogicalPacket) error {

	name := p.GetName()

	broker := p.driver.Broker()
	if broker != nil {
		broker.Broadcast(syncstate.NewEvent(addr, syncstate.Start, nil), syncstate.Start.String())
		return nil
	}

	return fmt.Errorf("%s addr '%s' dbsync broker not found", name, addr)

}
