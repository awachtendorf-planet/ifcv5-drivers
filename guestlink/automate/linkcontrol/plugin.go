package linkcontrol

import (
	"fmt"
	"sync"

	guestlink "github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *guestlink.Dispatcher
	automate *dispatcher.Automate

	answer      map[string]bool
	answerGuard sync.RWMutex

	kill      chan struct{}
	waitGroup sync.WaitGroup
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *guestlink.Dispatcher) *Plugin {
	return &Plugin{driver: parent}

}

const (
	retryDelay      = guestlink.RetryDelay
	packetTimeout   = guestlink.PacketTimeout
	answerTimeout   = guestlink.AnswerTimeout
	nextActionDelay = guestlink.NextActionDelay
	aliveTimeout    = guestlink.AliveTimeout
	maxError        = guestlink.MaxError
	maxNetworkError = 2

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
	hello       = automatestate.NextAction

	linkDown = automatestate.LinkDown
	linkUp   = automatestate.LinkUp
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)

	p.answer = make(map[string]bool)
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

	automate.RegisterRule(template.PacketVerify, linkDown, p.handleCommandAcknowledge, dispatcher.StateAction{NextState: hello, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketError, linkDown, p.handleCommandRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
	automate.RegisterRule(template.PacketStart, linkDown, p.sendCommandAcknowledge, dispatcher.StateAction{NextState: hello, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketTest, linkDown, p.sendCommandAcknowledge, dispatcher.StateAction{NextState: hello, NextTimeout: nextActionDelay})

	automate.RegisterRule(template.PacketVerify, hello, p.handleCommandAcknowledge, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketError, hello, p.handleCommandRefusion, dispatcher.StateAction{NextState: linkUp, NextTimeout: nextActionDelay})

	automate.RegisterRule(template.PacketStart, linkUp, p.sendCommandAcknowledge, dispatcher.StateAction{NextTimeout: aliveTimeout})
	automate.RegisterRule(template.PacketTest, linkUp, p.sendCommandAcknowledge, dispatcher.StateAction{NextTimeout: aliveTimeout})
	automate.RegisterRule(template.PacketVerify, linkUp, p.handleCommandAcknowledge, dispatcher.StateAction{NextTimeout: aliveTimeout})
	automate.RegisterRule(template.PacketError, linkUp, p.handleCommandRefusion, dispatcher.StateAction{NextState: linkDown, NextTimeout: retryDelay})

}

func (p *Plugin) expectAnswer(addr string) {
	p.answerGuard.Lock()
	p.answer[addr] = true
	p.answerGuard.Unlock()
}

func (p *Plugin) receivedAnswer(addr string) {
	p.answerGuard.Lock()
	p.answer[addr] = false
	p.answerGuard.Unlock()
}

func (p *Plugin) isAnswerExpected(addr string) bool {
	p.answerGuard.RLock()
	state, exist := p.answer[addr]
	p.answerGuard.RUnlock()
	return exist && state
}

func (p *Plugin) clearAddress(addr string) {
	p.answerGuard.Lock()
	delete(p.answer, addr)
	p.answerGuard.Unlock()
}
