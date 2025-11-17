package callpacket

import (
	"fmt"
	"strings"
	"sync"

	nortel "github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/nortel/template"

	"github.com/pkg/errors"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver       *nortel.Dispatcher
	records      map[string]PbxRecord
	recordsGuard sync.RWMutex
}

// New return a new plugin
func New(parent *nortel.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
		make(map[string]PbxRecord),
		sync.RWMutex{},
	}

	p.Setup(rq.Config{
		Name:          fmt.Sprintf("%T", p),
		ProcessPacket: p.processPacket,
		Startup:       p.startup,
		Shutdown:      p.shutdown,
	})

	p.RegisterPacket(template.PacketCallPacket)

	return p
}

func (p *Plugin) startup() {

	p.loadPendingRecords()

}

func (p *Plugin) shutdown() {

	p.storePendingRecords()

}

func (p *Plugin) processPacket(addr string, _ *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketCallPacket:
		p.handleCallPacket(addr, packet)

	default:
		return errors.Errorf("no handler defined to process packet '%s'", packet.Name)

	}

	return nil
}

func (p *Plugin) getValue(packet *ifc.LogicalPacket, key string) ([]byte, bool) {

	if packet == nil {
		return []byte{}, false
	}

	data, exist := packet.Data()[key]
	return data, exist
}

func (p *Plugin) getValueAsString(packet *ifc.LogicalPacket, key string) string {
	data, _ := p.getValue(packet, key)
	value := string(data)
	value = strings.Trim(value, " ")
	return value
}

func (p *Plugin) getExtension(packet *ifc.LogicalPacket) string {
	ext := p.getValueAsString(packet, "Originator")
	ext = strings.TrimLeft(ext, "DN")
	return ext
}
