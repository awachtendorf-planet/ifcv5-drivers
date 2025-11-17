package linkcontrol

import (
	"fmt"
	"sync"

	mitel "github.com/weareplanet/ifcv5-drivers/mitel/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/mitel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/mitel/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *mitel.Dispatcher
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
func New(parent *mitel.Dispatcher) *Plugin {
	return &Plugin{driver: parent}

}

const (
	retryDelay      = mitel.RetryDelay
	packetTimeout   = mitel.PacketTimeout
	nextActionDelay = mitel.NextActionDelay
	aliveTimeout    = mitel.AliveTimeout
	pmsSyncTimeout  = mitel.PmsSyncTimeout
	maxError        = mitel.MaxError
	maxNetworkError = 2

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
	nextAction  = automatestate.NextAction
	nextRecord  = automatestate.NextRecord

	linkDown = automatestate.LinkDown
	linkUp   = automatestate.LinkUp

	resyncStart  = automatestate.ResyncStart
	resyncRecord = automatestate.ResyncRecord
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
		HandleError:     p.handleError,
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

	automate.RegisterRule(template.PacketAck, commandSent, p.handlePacketAcknowledge, dispatcher.StateAction{NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketNak, commandSent, p.handlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
	automate.RegisterRule(template.PacketEnq, commandSent, p.handlePacketEnquire, dispatcher.StateAction{NextTimeout: retryDelay})

	automate.RegisterRule(template.PacketEnq, linkUp, p.setAlive, dispatcher.StateAction{NextTimeout: aliveTimeout})

	automate.RegisterRule(template.PacketSwapEnquiry, commandSent, p.handleSwapEnquiry, dispatcher.StateAction{NextState: resyncStart, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketSwapEnquiry, linkDown, p.handleSwapEnquiry, dispatcher.StateAction{NextState: resyncStart, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketSwapEnquiry, linkUp, p.handleSwapEnquiry, dispatcher.StateAction{NextState: resyncStart, NextTimeout: nextActionDelay})

}
