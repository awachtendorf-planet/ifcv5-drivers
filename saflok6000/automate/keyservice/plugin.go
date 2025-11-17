package keyservice

import (
	"fmt"
	"sync"

	saflok6000 "github.com/weareplanet/ifcv5-drivers/saflok6000/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state     *helper.State
	driver    *saflok6000.Dispatcher
	automate  *dispatcher.Automate
	job       chan order.DriverJob
	kill      chan struct{}
	waitGroup sync.WaitGroup
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *saflok6000.Dispatcher) *Plugin {
	return &Plugin{driver: parent}
}

const (
	retryDelay       = saflok6000.RetryDelay
	packetTimeout    = saflok6000.PacketTimeout
	nextActionDelay  = saflok6000.NextActionDelay
	keyAnswerTimeout = saflok6000.KeyAnswerTimeout
	maxError         = saflok6000.MaxKeyEncoderError
	maxNetworkError  = saflok6000.MaxNetworkError

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	answer      = automatestate.WaitForAnswer
	success     = automatestate.Success
	timeout     = automatestate.Timeout
	shutdown    = automatestate.Shutdown
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)

	p.state = helper.NewState()
	p.job = make(chan order.DriverJob)
	p.kill = make(chan struct{})

	config := dispatcher.AutomateConfig{
		Name:             fmt.Sprintf("%T", p),
		RetryDelay:       retryDelay,
		PacketTimeout:    packetTimeout,
		NextActionDelay:  nextActionDelay,
		MaxError:         maxError,
		MaxNetworkError:  maxNetworkError,
		HandleNextAction: p.handleNextAction,
		HandleError:      p.handleFailed,
	}

	p.automate = automate
	p.automate.Configure(config)
	p.registerStateMaschine()

	p.registerAction()
}

// Startup ...
func (p *Plugin) Startup() {
	if p.driver == nil {
		log.Error().Msgf("startup %T failed, err=%s", p, "driver not registered")
		return
	}
	log.Info().Msgf("startup %T", p)

	go func() {
		p.main()
	}()
}

// Cleanup ...
func (p *Plugin) Cleanup() {
	log.Info().Msgf("shutdown %T", p)
	close(p.kill)
	p.waitGroup.Wait()
}

func (p *Plugin) registerStateMaschine() {
	p.automate.RegisterRule(template.PacketAck, commandSent, p.automate.HandlePacketAcknowledge, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketNak, commandSent, p.automate.HandlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
}

func (p *Plugin) registerAction() {
	name := p.automate.Name()
	dispatcher := p.automate.Dispatcher()

	dispatcher.RegisterDriverJobHandler(name, p.job, order.KeyService)
}
