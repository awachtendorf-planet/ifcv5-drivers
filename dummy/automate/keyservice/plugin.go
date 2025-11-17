package keyservice

import (
	"time"

	dummy "github.com/weareplanet/ifcv5-drivers/dummy/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	ks "github.com/weareplanet/ifcv5-main/ifc/automate/keyservice"

	"github.com/weareplanet/ifcv5-drivers/dummy/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/dummy/template"

	"github.com/pkg/errors"
)

const (
	keyLifetime = 120 * time.Second
)

// Plugin ...
type Plugin struct {
	*ks.Plugin
	driver *dummy.Dispatcher
}

// New return a new plugin
func New(parent *dummy.Dispatcher) *Plugin {

	p := &Plugin{
		ks.New(),
		parent,
	}

	p.Setup(ks.Config{

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,

		InitHandler: p.init,
		SendPacket:  p.send,
		PreCheck:    p.preCheck,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	p.RegisterRule(template.PacketEnq, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{NextState: automatestate.NextRecord})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.KeyRequest,
		template.PacketKeyRequest,
	)

	p.RegisterWorkflow(order.KeyDelete,
		template.PacketKeyDelete,
	)

}

func (p *Plugin) preCheck(addr string, job *order.Job) (bool, error) {

	if job != nil && job.Timestamp > 0 {
		timestamp := time.Unix(0, job.Timestamp)
		since := time.Since(timestamp)
		if since > keyLifetime {
			return false, errors.Errorf("key service deadline exceeded (%s)", since)
		}
	}

	return true, nil
}

func (p *Plugin) handleReply(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	// success -> create pms reply
	p.HandleSuccess(addr, job, "track id")
	p.ChangeState(addr, automatestate.Success)

	// error -> create pms reply
	// p.HandleError(addr, job, "error reason")
	// p.ChangeState(addr, automatestate.Success)

	// ignore packet, re-calculate answer timeout
	// p.CalculateNextTimeout(addr, action, job)

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
