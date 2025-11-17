package datachange

import (
	"fmt"
	"sync"

	centigram "github.com/weareplanet/ifcv5-drivers/centigram/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/centigram/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/centigram/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state     *helper.State
	driver    *centigram.Dispatcher
	automate  *dispatcher.Automate
	workflow  *utils.NNSMap
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
func New(parent *centigram.Dispatcher) *Plugin {
	return &Plugin{driver: parent}
}

const (
	retryDelay      = centigram.RetryDelay
	packetTimeout   = centigram.PacketTimeout
	nextActionDelay = centigram.NextActionDelay
	maxError        = maxRefusion + 1
	maxNetworkError = centigram.MaxNetworkError

	busy        = automatestate.Busy
	nextAction  = automatestate.NextAction
	commandSent = automatestate.CommandSent
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)

	p.state = helper.NewState()
	p.workflow = utils.NewNNSMap()
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
		HandleError:      p.handleError,
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

	go p.main()
}

// Cleanup ...
func (p *Plugin) Cleanup() {
	log.Info().Msgf("shutdown %T", p)
	close(p.kill)
	p.waitGroup.Wait()
}

func (p *Plugin) registerStateMaschine() {
	p.automate.RegisterRule(template.PacketAck, commandSent, p.automate.HandlePacketAcknowledge, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketNak, commandSent, p.handlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
	p.automate.RegisterRule(template.PacketRefusion, commandSent, p.handlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
}

func (p *Plugin) registerAction() {
	name := p.automate.Name()
	dispatcher := p.automate.Dispatcher()

	dispatcher.RegisterDriverJobHandler(name, p.job, order.ASW)
}
