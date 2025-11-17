package dbsync

import (
	"fmt"

	informatel "github.com/weareplanet/ifcv5-drivers/informatel/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/informatel/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/informatel/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *informatel.Dispatcher
	automate *dispatcher.Automate
	kill     chan struct{}
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *informatel.Dispatcher) *Plugin {
	return &Plugin{driver: parent}
}

const (
	retryDelay      = informatel.RetryDelay
	packetTimeout   = informatel.PacketTimeout
	nextActionDelay = informatel.NextActionDelay
	maxError        = 3
	maxNetworkError = informatel.MaxNetworkError
	pmsTimeOut      = informatel.PmsTimeout
	pmsSyncTimeout  = informatel.PmsSyncTimeout

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown

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
	p.registerStateMaschine()
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
}

func (p *Plugin) registerStateMaschine() {
	p.automate.RegisterRule(template.PacketAck, commandSent, p.automate.HandlePacketAcknowledge, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketNak, commandSent, p.automate.HandlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
}
