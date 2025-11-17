package datachange

import (
	"fmt"
	"sync"

	guestlink "github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state     *helper.State
	driver    *guestlink.Dispatcher
	automate  *dispatcher.Automate
	job       chan order.DriverJob
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
	maxError        = guestlink.MaxError
	maxNetworkError = guestlink.MaxNetworkError

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	answer      = automatestate.WaitForAnswer
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	timeout     = automatestate.Timeout
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)

	p.state = helper.NewState()
	p.job = make(chan order.DriverJob)
	p.kill = make(chan struct{})

	config := dispatcher.AutomateConfig{
		Name:             fmt.Sprintf("%T", p),
		RetryDelay:       retryDelay,
		PacketTimeout:    packetTimeout,
		NextActionDelay:  nextActionDelay,
		MaxError:         maxError,
		MaxNetworkError:  maxNetworkError,
		HandleNextAction: p.handleNextAction,
		HandleError:      p.handleError,
	}

	p.automate = automate
	p.automate.Configure(config)
	p.registerStateMaschine()

	p.registerAction()
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

	p.automate.RegisterRule(template.PacketAck, commandSent, p.automate.HandlePacketAcknowledge, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketNak, commandSent, p.automate.HandlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})

	p.automate.RegisterRule(template.PacketVerify, answer, p.handleCommandAcknowledge, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketVerify, commandSent, p.handleCommandAcknowledge, dispatcher.StateAction{})

	p.automate.RegisterRule(template.PacketError, answer, p.handleCommandRefusion, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketError, commandSent, p.handleCommandRefusion, dispatcher.StateAction{})

	// register automate packet route for router automate to dispatch VER/ERR packets
	p.automate.PacketRouteRegistered = func(addr string, route *ifc.Route) {
		if len(addr) == 0 || route == nil {
			return
		}
		p.driver.RegisterRoute(p, addr, route)
	}

}

func (p *Plugin) registerAction() {
	name := p.automate.Name()
	dispatcher := p.automate.Dispatcher()

	dispatcher.RegisterDriverJobHandler(name, p.job, order.ASW)
}
