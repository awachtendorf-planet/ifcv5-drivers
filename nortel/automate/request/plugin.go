package request

import (
	"fmt"
	"strings"

	nortel "github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/nortel/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *nortel.Dispatcher
}

// New return a new plugin
func New(parent *nortel.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{
		Name:          fmt.Sprintf("%T", p),
		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(
		template.PacketRoomStatus,
		template.PacketMinibarItem,
		template.PacketMinibarTotal,
		template.PacketVoiceCount,
		template.PacketWakeupSet,
		template.PacketWakeupClear,
		template.PacketWakeupAnswer,
	)

	return p
}

func (p *Plugin) processPacket(addr string, _ *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketRoomStatus:
		p.handleRoomStatus(addr, packet)

	case template.PacketMinibarItem, template.PacketMinibarTotal:
		p.handleMinibar(addr, packet)

	case template.PacketVoiceCount:
		p.handleVoiceCount(addr, packet)

	case template.PacketWakeupSet, template.PacketWakeupClear, template.PacketWakeupAnswer:
		p.handleWakeup(addr, packet)

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

func (p *Plugin) getValueAsNumeric(packet *ifc.LogicalPacket, key string) int {
	data, _ := p.getValue(packet, key)
	value := string(data)
	value = strings.Trim(value, " ")
	value = strings.TrimLeft(value, "0")
	return cast.ToInt(value)
}

func (p *Plugin) getExtension(packet *ifc.LogicalPacket) string {
	ext := p.getValueAsString(packet, "Ext")
	return ext
}
