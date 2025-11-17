package request

import (
	"fmt"

	dummy "github.com/weareplanet/ifcv5-drivers/dummy/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	//"github.com/weareplanet/ifcv5-main/ifc/record"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/dummy/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/dummy/template"
	// "github.com/pkg/errors"
	// "github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *dummy.Dispatcher
}

// New return a new plugin
func New(parent *dummy.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name: fmt.Sprintf("%T", p),

		// register Ack/Nak packets if needed
		// TemplateAck: template.PacketAck,
		// TemplateNak: template.PacketNak,

		// register init handler if needed
		// InitHandler: p.init,

		// register pre/post handler
		// PreHandler:    p.preHandler,
		// PostHandler:   p.postHandler,

		// register incoming packet handler
		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(template.PacketRequest1, template.PacketRequest2)

	return p
}

// func (p *Plugin) init() {
// 	// register some state handler if needed
// 	// p.RegisterHandler(rq.NextAction, p.handleNextAction)
// 	// p.RegisterHandler(rq.NextRecord, p.handleNextRecord)
// }

// func (p *Plugin) preHandler(addr string) {
// 	// called before the state maschine start
// }

// func (p *Plugin) postHandler(addr string) {
// 	// called after the state maschine finished
// }

// func (p *Plugin) handleNextAction(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {
// 	action.NextState = rq.NextRecord
// 	return nil
// }

// func (p *Plugin) handleNextRecord(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {
// 	action.NextState = rq.Success
// 	return nil
// }

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	// switch packet.Name {

	// case "packet_x":

	// case "packet_y":

	// default:
	// 	return errors.Errorf("no handler defined to process packet '%s'", packet.Name)

	// }

	/*
		// send a request to the pms

		// inquiry := record.PostingInquiry{
		// 	//Time:           time.Now(),
		// 	//Station:        station,
		// 	Inquiry:        "xyz",
		// 	MaximumMatches: 6,
		// 	RequestType:    "2",
		// }

		// reply, err := p.PmsRequest(addr, inquiry, packet)
		// if err != nil {
		// 	return err
		// }
		// do something with the reply packet

	*/

	// default: action.NextState == success, action.NextTimeout == nextActionDelay
	// change state to something action.NextState = rq.NextAction. make sure there is a state handler registered (p.RegisterHandler(rq.NextAction, p.handleNextAction))
	// return nil

	// return err shutdown the automate state maschine with failed state
	// return errors.New("something went wrong")

	// return nil shutdown the automate state maschine with success state
	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
