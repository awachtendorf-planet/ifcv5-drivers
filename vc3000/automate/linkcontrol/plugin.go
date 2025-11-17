package linkcontrol

import (
	"fmt"
	"sync"

	vc3000 "github.com/weareplanet/ifcv5-drivers/vc3000/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/vc3000/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/vc3000/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *vc3000.Dispatcher
	automate *dispatcher.Automate

	kill      chan struct{}
	waitGroup sync.WaitGroup
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *vc3000.Dispatcher) *Plugin {
	return &Plugin{driver: parent}

}

const (
	retryDelay      = vc3000.RetryDelay
	packetTimeout   = vc3000.PacketTimeout
	nextActionDelay = vc3000.NextActionDelay
	pendingDelay    = vc3000.PendingDelay
	aliveTimeout    = vc3000.AliveTimeout
	maxError        = vc3000.MaxError
	maxNetworkError = 2

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown

	linkDown = automatestate.LinkDown
	linkUp   = automatestate.LinkUp
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)

	p.kill = make(chan struct{})
	p.state = helper.NewState()

	config := dispatcher.AutomateConfig{
		Name:            fmt.Sprintf("%T", p),
		RetryDelay:      retryDelay,
		PacketTimeout:   packetTimeout,
		NextActionDelay: nextActionDelay,
		MaxError:        maxError,
		MaxNetworkError: maxNetworkError,
		LinkControl:     true,
	}

	p.automate = automate
	p.automate.Configure(config)
	p.registerStateMaschine()
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
	automate := p.automate

	automate.RegisterRule(template.PacketRegister, linkDown, p.handleRegister, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketRegisterAck, linkDown, p.handleRegisterAck, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketRegisterNak, linkDown, p.handleRegisterNak, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketUnregister, linkUp, p.handleUnregister, dispatcher.StateAction{NextState: linkDown, NextTimeout: pendingDelay})

}
