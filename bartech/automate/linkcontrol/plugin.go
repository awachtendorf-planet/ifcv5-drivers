package linkcontrol

import (
	"fmt"
	"sync"

	bartech "github.com/weareplanet/ifcv5-drivers/bartech/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/bartech/dispatcher"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state    *helper.State
	driver   *bartech.Dispatcher
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
func New(parent *bartech.Dispatcher) *Plugin {
	return &Plugin{driver: parent}

}

const (
	maxError        = bartech.MaxError
	maxNetworkError = 2

	initial  = automatestate.Initial
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
		MaxError:        maxError,
		MaxNetworkError: maxNetworkError,
	}

	p.automate = automate
	p.automate.Configure(config)
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
