package linkcontrol

import (
	"fmt"
	"sync"
	"time"

	robobar "github.com/weareplanet/ifcv5-drivers/robobar/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/robobar/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/robobar/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *robobar.Dispatcher
	automate *dispatcher.Automate

	lastRestartMessage map[string]time.Time
	lastRestartGuard   sync.RWMutex
	kill               chan struct{}
	waitGroup          sync.WaitGroup
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *robobar.Dispatcher) *Plugin {
	return &Plugin{
		driver:             parent,
		lastRestartMessage: make(map[string]time.Time),
	}

}

const (
	retryDelay      = robobar.RetryDelay
	packetTimeout   = robobar.PacketTimeout
	nextActionDelay = robobar.NextActionDelay
	aliveTimeout    = robobar.AliveTimeout
	pmsTimeOut      = robobar.PmsTimeout
	pmsSyncTimeout  = robobar.PmsSyncTimeout
	maxError        = robobar.MaxError
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

	resyncStart      = automatestate.ResyncStart
	resyncRecord     = automatestate.ResyncRecord
	resyncNextRecord = automatestate.ResyncNextRecord
	resyncEnd        = automatestate.ResyncEnd
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

func (p *Plugin) setLastRestartMessage(addr string) {
	p.lastRestartGuard.Lock()
	p.lastRestartMessage[addr] = time.Now()
	p.lastRestartGuard.Unlock()
}

func (p *Plugin) getLastRestartMessage(addr string) time.Time {
	p.lastRestartGuard.RLock()
	lastTimestamp, exist := p.lastRestartMessage[addr]
	p.lastRestartGuard.RUnlock()
	if !exist {
		return time.Unix(0, 0)
	}
	return lastTimestamp
}

func (p *Plugin) clearLastRestartMessage(addr string) {
	p.lastRestartGuard.Lock()
	delete(p.lastRestartMessage, addr)
	p.lastRestartGuard.Unlock()
}

func (p *Plugin) registerStateMaschine() {
	automate := p.automate

	automate.RegisterRule(template.PacketAck, commandSent, p.handlePacketAcknowledge, dispatcher.StateAction{NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketNak, commandSent, p.handlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})

	automate.RegisterRule(template.PacketRestart, linkDown, p.sendStartup, dispatcher.StateAction{NextState: resyncStart, NextTimeout: nextActionDelay})
	automate.RegisterRule(template.PacketRestart, linkUp, p.sendStartup, dispatcher.StateAction{NextState: resyncStart, NextTimeout: nextActionDelay})

}
