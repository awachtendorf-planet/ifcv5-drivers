package request

import (
	"fmt"

	syncstate "github.com/weareplanet/ifcv5-drivers/caracas/state"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

func (p *Plugin) handleDBSyncReq(addr string, packet *ifc.LogicalPacket) error {

	name := p.GetName()

	syncType := 0
	reqType := p.getField(packet, "Status", false)

	switch reqType {

	case "1":
		syncType = 0b00001
	case "4":
		syncType = 0b01000
	case "5":
		syncType = 0b00100
	case "6":
		syncType = 0b10000
	case "0":
		syncType = 0b11101
	}

	broker := p.driver.Broker()
	if broker != nil {
		broker.Broadcast(syncstate.NewEvent(addr, syncType, syncstate.Start, nil), syncstate.Start.String())
		return nil
	}

	return fmt.Errorf("%s addr '%s' dbsync broker not found", name, addr)

}
