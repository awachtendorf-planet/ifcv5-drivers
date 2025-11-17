package dbsync

import (
	"fmt"

	mitel "github.com/weareplanet/ifcv5-drivers/mitel/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
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
	workflow *utils.NNSMap
	kill     chan struct{}
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
	retryDelay      = mitel.SwapRetryDelay
	packetTimeout   = mitel.PacketTimeout
	nextActionDelay = mitel.NextActionDelay
	maxError        = 5
	maxNetworkError = mitel.MaxNetworkError
	pmsTimeOut      = mitel.PmsTimeout
	pmsSyncTimeout  = mitel.PmsSyncTimeout

	busy        = automatestate.Busy
	commandSent = automatestate.CommandSent
	nextAction  = automatestate.NextAction
	nextRecord  = automatestate.NextRecord
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)
	p.kill = make(chan struct{})
	p.state = helper.NewState()
	p.workflow = utils.NewNNSMap()

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
