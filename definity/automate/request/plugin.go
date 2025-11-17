package request

import (
	"fmt"
	"strings"

	definity "github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/definity/template"

	"github.com/pkg/errors"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *definity.Dispatcher
}

// New return a new plugin
func New(parent *definity.Dispatcher) *Plugin {

	rq.FixSerialError = true

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,

		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(
		template.PacketMessageLampOn,
		template.PacketMessageLampOff,
		template.PacketSetRestriction,
		template.PacketHouseKeeperRoom,
		template.PacketHouseKeeperStation,
	)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	var err error

	switch packet.Name {

	case template.PacketMessageLampOn, template.PacketMessageLampOff:

		err = p.handleMessageLamp(addr, action, packet)

	case template.PacketHouseKeeperRoom, template.PacketHouseKeeperStation:

		err = p.handleHouseKeeperStatus(addr, action, packet)

	case template.PacketSetRestriction:

		err = errors.Errorf("no idea how we should process the packet '%s', because the values are not compatible", packet.Name)

	default:

		err = errors.Errorf("no handler defined to process packet '%s'", packet.Name)

	}

	return err

}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		return err
	}

	// create a object copy, because of answertimeout are set by underlaying base function
	//sendAction := &dispatcher.StateAction{NextTimeout: action.NextTimeout, CurrentState: action.CurrentState, NextState: action.NextState}

	//err = p.SendPacket(addr, packet, sendAction)
	err = p.SendPacket(addr, packet, action)

	return err
}

// todo: create dispatcher function for transparent mode

func (p *Plugin) getString(packet *ifc.LogicalPacket, field string) string {

	if packet == nil || len(field) == 0 {
		return ""
	}

	value := string(packet.Data()[field])
	value = strings.TrimLeft(value, "0")

	return value
}

func (p *Plugin) getRoom(packet *ifc.LogicalPacket) string {
	return p.getString(packet, "RSN")
}

func (p *Plugin) getProc(packet *ifc.LogicalPacket) byte {

	if packet == nil {
		return 0
	}

	proc := string(packet.Data()["PROC"])

	if len(proc) > 0 {
		return proc[0]
	}

	return 0

}
