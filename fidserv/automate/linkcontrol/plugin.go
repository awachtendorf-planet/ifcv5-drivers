package linkcontrol

import (
	"fmt"
	"sync"

	fidserv "github.com/weareplanet/ifcv5-drivers/fidserv/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/fidserv/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *fidserv.Dispatcher
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
func New(parent *fidserv.Dispatcher) *Plugin {
	return &Plugin{driver: parent}

}

const (
	retryDelay      = fidserv.RetryDelay
	packetTimeout   = fidserv.PacketTimeout
	nextActionDelay = fidserv.NextActionDelay
	aliveTimeout    = fidserv.AliveTimeout
	maxError        = fidserv.MaxError
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
	pending := p.state.Length()
	if pending > 0 {
		log.Debug().Msgf("shutdown %T has %d pending handler", p, pending)
	}
	close(p.kill)
	p.waitGroup.Wait()
}

func (p *Plugin) registerStateMaschine() {
	automate := p.automate

	automate.RegisterRule(template.PacketAck, commandSent, p.handlePacketAcknowledge, dispatcher.StateAction{})
	automate.RegisterRule(template.PacketNak, commandSent, p.handlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})

	automate.RegisterRule(template.PacketLinkStart, linkDown, p.sendLinkStart, dispatcher.StateAction{NextState: linkDown, NextTimeout: retryDelay})
	automate.RegisterRule(template.PacketLinkAlive, linkDown, p.sendLinkAlive, dispatcher.StateAction{NextState: linkUp, NextTimeout: aliveTimeout})
	automate.RegisterRule(template.PacketLinkDescription, linkDown, p.handleLinkDescription, dispatcher.StateAction{NextTimeout: retryDelay})
	automate.RegisterRule(template.PacketLinkRecord, linkDown, p.handleLinkRecord, dispatcher.StateAction{NextTimeout: retryDelay})

	automate.RegisterRule(template.PacketLinkStart, linkUp, p.sendLinkAlive, dispatcher.StateAction{NextState: linkUp, NextTimeout: aliveTimeout})
	automate.RegisterRule(template.PacketLinkAlive, linkUp, automate.HandlePacket, dispatcher.StateAction{NextState: linkUp, NextTimeout: aliveTimeout})
	automate.RegisterRule(template.PacketLinkDescription, linkUp, p.handleLinkDescription, dispatcher.StateAction{NextState: linkDown, NextTimeout: retryDelay})
	automate.RegisterRule(template.PacketLinkRecord, linkUp, p.handleLinkRecord, dispatcher.StateAction{NextTimeout: aliveTimeout})

	automate.RegisterRule(template.PacketLinkEnd, linkUp, automate.HandlePacket, dispatcher.StateAction{NextState: linkDown, NextTimeout: retryDelay})
}
