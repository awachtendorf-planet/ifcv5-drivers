package dbsync

import (
	"fmt"

	definity "github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/definity/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/definity/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *definity.Dispatcher
	automate *dispatcher.Automate
	kill     chan struct{}
}

// ...
var (
	PluginID                            = (*Plugin)(nil)
	_        dispatcher.PluginInterface = (*Plugin)(nil)
)

// New return a new plugin
func New(parent *definity.Dispatcher) *Plugin {
	return &Plugin{driver: parent}
}

const (
	retryDelay      = definity.RetryDelay
	packetTimeout   = definity.PacketTimeout
	nextActionDelay = definity.NextActionDelay
	maxError        = 3
	maxNetworkError = definity.MaxNetworkError
	pmsTimeOut      = definity.PmsTimeout
	pmsSyncTimeout  = definity.PmsSyncTimeout
	answerTimeout   = definity.AnswerTimeout

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	answer      = automatestate.WaitForAnswer
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
	p.automate.RegisterRule(template.PacketRoomDataImageSwap, answer, p.handleReply, dispatcher.StateAction{NextTimeout: nextActionDelay, NextState: resyncNextRecord})
}
