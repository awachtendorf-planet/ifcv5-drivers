package request

import (
	"fmt"
	"strings"
	"time"

	messerschmitt "github.com/weareplanet/ifcv5-drivers/messerschmitt/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/requestcomplex"

	"github.com/weareplanet/ifcv5-drivers/messerschmitt/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/messerschmitt/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *messerschmitt.Dispatcher
}

// New return a new plugin
func New(parent *messerschmitt.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,
		TemplateEnq: template.PacketEnq,
		TemplateEOT: template.PacketEOT,

		SendENQ: p.sendENQ,
		SendEOT: p.sendEOT,

		SendPacket:    p.send,
		ProcessPacket: p.processPacket,
	})

	p.RegisterWorkflow(template.PacketFunctionCards, "")
	p.RegisterWorkflow(template.PacketSpecialReaders, "")
	p.RegisterWorkflow(template.PacketVendingMachines, "")
	p.RegisterWorkflow(template.PacketRoomEquipment, "")

	p.RegisterWorkflow(template.PacketSynchronisation, template.PacketSynchronisation)

	return p
}

func (p *Plugin) sendENQ(addr string, _ *order.Job) bool {

	if dispatcher := p.GetDispatcher(); dispatcher != nil {
		if !dispatcher.IsSerialDevice(addr) {
			return false
		}
	}

	return true
}

func (p *Plugin) sendEOT(addr string, job *order.Job) bool {

	if dispatcher := p.GetDispatcher(); dispatcher != nil {
		if !dispatcher.IsSerialDevice(addr) {
			return false
		}
	}

	// if we reply with a packet thats ends of ETB (0x17) there is no need for EOT

	return true

}

func (p *Plugin) getKeyAsString(key string, packet *ifc.LogicalPacket) string {

	if packet == nil || len(key) == 0 {
		return ""
	}

	value := string(packet.Data()[key])
	value = strings.TrimLeft(value, " ")

	return value
}

func (p *Plugin) getKeyAsInt(key string, packet *ifc.LogicalPacket) int {

	if packet == nil || len(key) == 0 {
		return 0
	}

	value := string(packet.Data()[key])

	value = strings.TrimLeft(value, " ")
	value = strings.TrimLeft(value, "0")

	n := cast.ToInt(value)

	return n
}

func (p *Plugin) getRoom(addr string, packet *ifc.LogicalPacket) string {

	var room string

	if p.driver.IsSerialDevice(addr) {
		number := p.getKeyAsInt("Room", packet)
		if number == 0 {
			return ""
		}
		room = cast.ToString(number)
	} else {
		room = p.getKeyAsString("Room", packet)
	}

	return room

}

func (p *Plugin) getTime(addr string, packet *ifc.LogicalPacket) time.Time {

	d := p.getKeyAsString("Date", packet)
	t := p.getKeyAsString("Time", packet)

	var ts time.Time

	if len(d) == 6 && len(t) == 6 {
		if c, err := time.Parse("020106150405", d+t); err == nil { // ddmmyyhhmmss
			ts = c
		}
	} else if len(d) == 8 && len(t) == 6 {
		if c, err := time.Parse("02012006150405", d+t); err == nil { // ddmmyyyyhhmmss
			ts = c
		}
	}

	return ts
}

func (p *Plugin) processPacket(addr string, packet *ifc.LogicalPacket) error {

	switch packet.Name {

	case template.PacketSynchronisation:

	case template.PacketFunctionCards:
		return p.handleFunctionCards(addr, packet)

	case template.PacketSpecialReaders:
		return p.handleSpecialReaders(addr, packet)

	case template.PacketVendingMachines:
		return p.handleVendingMachines(addr, packet)

	case template.PacketRoomEquipment:
		return p.handleRoomEquipment(addr, packet)

	default:
		return errors.Errorf("no handler defined to process packet '%s'", packet.Name)

	}

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, _ *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
