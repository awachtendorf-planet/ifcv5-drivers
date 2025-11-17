package linkcontrol

import (
	"fmt"
	"sync"

	saflok6000 "github.com/weareplanet/ifcv5-drivers/saflok6000/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/saflok6000/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/saflok6000/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *saflok6000.Dispatcher
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
func New(parent *saflok6000.Dispatcher) *Plugin {
	return &Plugin{driver: parent}

}

const (
	retryDelay      = saflok6000.RetryDelay
	packetTimeout   = saflok6000.PacketTimeout
	nextActionDelay = saflok6000.NextActionDelay
	aliveTimeout    = saflok6000.AliveTimeout
	maxError        = saflok6000.MaxError
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

	automate.RegisterRule(template.PacketAck, commandSent, p.handlePacketAcknowledge, dispatcher.StateAction{})
	automate.RegisterRule(template.PacketNak, commandSent, p.handlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})

	automate.RegisterRule(template.PacketAlive, linkDown, p.handlePacket, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketAlive, linkUp, p.handlePacket, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})

	automate.RegisterRule(template.PacketBeaconRequest, linkUp, p.sendBeaconResponse, dispatcher.StateAction{})

}
