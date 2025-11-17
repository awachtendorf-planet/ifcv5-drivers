package keyservice

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"

	ks "github.com/weareplanet/ifcv5-main/ifc/automate/keyservice"

	miwa "github.com/weareplanet/ifcv5-drivers/miwa/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/miwa/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	keyLifetime = 120 * time.Second
)

// Plugin ...
type Plugin struct {
	*ks.Plugin
	driver *miwa.Dispatcher
}

// New return a new plugin
func New(parent *miwa.Dispatcher) *Plugin {

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

	p.RegisterRule(template.PacketKeyCreateResult, automatestate.WaitForAnswer, p.handleKeyAnswer, dispatcher.StateAction{NextState: automatestate.NextRecord})

	p.RegisterRule(template.PacketKeyReadAnswer, automatestate.WaitForAnswer, p.handleKeyAnswer, dispatcher.StateAction{NextState: automatestate.NextRecord})

	p.RegisterRule(template.PacketError, automatestate.WaitForAnswer, p.handleKeyAnswer, dispatcher.StateAction{NextState: automatestate.NextRecord})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.KeyRequest,
		template.PacketKeyCreate,
	)

	p.RegisterWorkflow(order.KeyRead,
		template.PacketKeyRead,
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

func (p *Plugin) handleKeyAnswer(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	switch in.Name {

	case template.PacketKeyReadAnswer:
		result := p.getKeyAsString("Result", in)

		if result == "0" {
			ID := p.getKeyAsString("ID", in)
			keyInfo := p.getKeyAsString("KeyCardInformation", in)

			posInfo := p.getKeyAsString("POSInfo", in)

			p.HandleSuccess(addr, job, ID, keyInfo, posInfo)

			p.ChangeState(addr, action.NextState)

			return nil
		}

		reason := p.driver.GetKeyReadResponseText(result)

		p.HandleError(addr, job, reason)

		p.ChangeState(addr, action.NextState)

		return nil

	case template.PacketKeyCreateResult:

		result := p.getKeyAsString("Result", in)
		ID := p.getKeyAsString("ID", in)

		if result == "0" {
			p.HandleSuccess(addr, job, ID)

			p.ChangeState(addr, action.NextState)

			return nil
		}

		reason := p.driver.GetKeyCreateResponse(result)

		p.HandleError(addr, job, reason)

		p.ChangeState(addr, action.NextState)

		return nil

	case template.PacketError:
		errCode := p.getKeyAsString("Error", in)

		reason := p.driver.GetErrorText(errCode)

		p.HandleError(addr, job, reason)

		p.ChangeState(addr, action.NextState)
		return nil
	}

	p.CalculateNextTimeout(addr, action, job)

	return nil
}

func (p *Plugin) getKeyAsString(key string, in *ifc.LogicalPacket) string {
	returnString := cast.ToString(in.Data()[key])
	returnString = strings.Trim(returnString, " ")
	return returnString
}

// func (p *Plugin) getKeyAsInt(key string, in *ifc.LogicalPacket) int {
// 	data := in.Data()[key]
// 	value := string(data)
// 	value = strings.TrimLeft(value, " ")
// 	value = strings.TrimLeft(value, "0")

// 	n := cast.ToInt(value)
// 	return n
// }

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, job *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
