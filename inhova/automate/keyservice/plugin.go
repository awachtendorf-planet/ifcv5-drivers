package keyservice

import (
	"strings"
	"time"

	inhova "github.com/weareplanet/ifcv5-drivers/inhova/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	ks "github.com/weareplanet/ifcv5-main/ifc/automate/keyservice"

	"github.com/weareplanet/ifcv5-drivers/inhova/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/inhova/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	keyLifetime = 120 * time.Second
)

// Plugin ...
type Plugin struct {
	*ks.Plugin
	driver *inhova.Dispatcher
}

// New return a new plugin
func New(parent *inhova.Dispatcher) *Plugin {

	p := &Plugin{
		ks.New(),
		parent,
	}

	p.Setup(ks.Config{

		TemplateAck:      template.PacketAck,
		TemplateNak:      template.PacketNak,
		GetAnswerTimeout: p.getAnswerTimeout,

		InitHandler: p.init,
		SendPacket:  p.send,
		PreCheck:    p.preCheck,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	p.RegisterRule(template.PacketCodeCardAnswer, automatestate.WaitForAnswer, p.handleKeyAnswer, dispatcher.StateAction{NextState: automatestate.Success})
	p.RegisterRule(template.PacketError, automatestate.WaitForAnswer, p.handleKeyAnswer, dispatcher.StateAction{NextState: automatestate.Success})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.KeyRequest,
		template.PacketKeyRequest,
	)

	p.RegisterWorkflow(order.KeyDelete,
		template.PacketKeyDelete,
	)

}

func (p *Plugin) preCheck(_ string, job *order.Job) (bool, error) {

	if job != nil && job.Timestamp > 0 {
		timestamp := time.Unix(0, job.Timestamp)
		since := time.Since(timestamp)
		if since > keyLifetime {
			return false, errors.Errorf("key service deadline exceeded (%s)", since)
		}
	}

	return true, nil
}

func (p *Plugin) getAnswerTimeout(addr string, job *order.Job) time.Duration {

	answerTimeout := inhova.KeyAnswerTimeout

	if timeout, ok := p.GetSetting(addr, defines.KeyAnswerTimeout, uint(0)).(uint64); ok && timeout > 0 {
		if duration := cast.ToDuration(time.Duration(timeout) * time.Second); duration > 0 {
			answerTimeout = duration
		}
	}

	keyCount := p.driver.GetKeyCount(job)

	return answerTimeout * time.Duration(keyCount)
}

func (p *Plugin) getKeyAsString(key string, in *ifc.LogicalPacket) string {
	returnString := cast.ToString(in.Data()[key])
	returnString = strings.Trim(returnString, " ")
	return returnString
}

func (p *Plugin) getKeyAsInt(key string, in *ifc.LogicalPacket) int {
	data := in.Data()[key]
	value := string(data)
	value = strings.TrimLeft(value, " ")
	value = strings.TrimLeft(value, "0")

	n := cast.ToInt(value)
	return n
}

func (p *Plugin) handleKeyAnswer(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	name := p.GetName()

	encoderAddr := p.getKeyAsInt("Encoder", in)
	encoder := p.GetEncoder(addr)

	if encoderAddr != encoder {
		log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' for encoder: %d, expect: %d", name, addr, in.Name, encoderAddr, encoder)
		p.CalculateNextTimeout(addr, action, job)
		return nil
	}

	switch in.Name {

	case template.PacketCodeCardAnswer:

		responceAction := p.getKeyAsString("Action", in) // I/G/O
		keyType, _ := p.driver.GetKeyType(job)           // no error expected, because of dispatcher pre-check

		if keyType != responceAction {
			log.Debug().Msgf("%s addr '%s' mismatching key type '%s', expect '%s'", name, addr, responceAction, keyType)
			break
		}

		cardID := p.getKeyAsString("CardID", in)
		p.HandleSuccess(addr, job, cardID)

		p.ChangeState(addr, action.NextState)
		return nil

	case template.PacketError:

		errorCode := p.getKeyAsString("Error", in)
		answer := p.driver.GetAnswerText("E" + errorCode)

		p.HandleError(addr, job, answer)

		p.ChangeState(addr, action.NextState)
		return nil

	}

	p.CalculateNextTimeout(addr, action, job)

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
