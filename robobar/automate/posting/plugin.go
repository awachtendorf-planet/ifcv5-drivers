package posting

import (
	"fmt"
	"sync"

	robobar "github.com/weareplanet/ifcv5-drivers/robobar/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/robobar/dispatcher"
	//"github.com/weareplanet/ifcv5-drivers/robobar/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state     *helper.State
	driver    *robobar.Dispatcher
	automate  *dispatcher.Automate
	done      chan string
	kill      chan struct{}
	waitGroup sync.WaitGroup
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *robobar.Dispatcher) *Plugin {
	return &Plugin{driver: parent}
}

const (
	retryDelay      = robobar.RetryDelay
	packetTimeout   = robobar.PacketTimeout
	nextActionDelay = robobar.NextActionDelay
	maxNetworkError = robobar.MaxNetworkError
	maxError        = robobar.MaxError
	pmsTimeOut      = robobar.PmsTimeout

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	commandSent = automatestate.CommandSent
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)
	p.done = make(chan string, 128)
	p.kill = make(chan struct{})
	p.state = helper.NewState()

	config := dispatcher.AutomateConfig{
		Name:             fmt.Sprintf("%T", p),
		RetryDelay:       retryDelay,
		PacketTimeout:    packetTimeout,
		NextActionDelay:  nextActionDelay,
		MaxError:         maxError,
		MaxNetworkError:  maxNetworkError,
		HandleNextAction: p.handleNextAction,
	}

	p.automate = automate
	p.automate.Configure(config)
	//p.registerStateMaschine()
}

// Startup ...
func (p *Plugin) Startup() {
	if p.driver == nil {
		log.Error().Msgf("startup %T failed, err=%s", p, "driver not registered")
		return
	}
	log.Info().Msgf("startup %T", p)

	go func() {
		if err := p.main(); err != nil {
			log.Error().Msgf("startup %T failed, err=%s", p, err)
		}
	}()
}

// Cleanup ...
func (p *Plugin) Cleanup() {
	log.Info().Msgf("shutdown %T", p)
	close(p.kill)
	p.waitGroup.Wait()
}

// func (p *Plugin) registerStateMaschine() {
// 	p.automate.RegisterRule(template.PacketAck, commandSent, p.automate.HandlePacketAcknowledge, dispatcher.StateAction{})
// 	p.automate.RegisterRule(template.PacketNak, commandSent, p.automate.HandlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
// }
