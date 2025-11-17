package requestitems

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

type nextaction struct {
	identifier    string
	sequence      int
	Context       interface{}
	TemplateName  string
	CorrelationId string
	TrackingId    string
}

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state        *helper.State
	driver       *guestlink.Dispatcher
	automate     *dispatcher.Automate
	records      map[string][]nextaction
	recordsMutex sync.RWMutex
	done         chan string
	kill         chan struct{}
	waitGroup    sync.WaitGroup
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
	nextActionDelay = guestlink.NextActionDelay
	maxNetworkError = guestlink.MaxNetworkError
	maxError        = guestlink.MaxError
	pmsTimeOut      = guestlink.PmsTimeout
	pmsSyncTimeout  = guestlink.PmsSyncTimeout

	initial     = automatestate.Initial
	idle        = automatestate.Idle
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
	p.records = make(map[string][]nextaction)
	p.done = make(chan string, 128)
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
	p.waitGroup.Wait()
}

func (p *Plugin) registerStateMaschine() {
	p.automate.RegisterRule(template.PacketAck, commandSent, p.automate.HandlePacketAcknowledge, dispatcher.StateAction{})
	p.automate.RegisterRule(template.PacketNak, commandSent, p.automate.HandlePacketRefusion, dispatcher.StateAction{NextTimeout: retryDelay})
}

func (p *Plugin) getSlotName(addr string, scope string) string {
	name := fmt.Sprintf("%s#%s", addr, scope)
	return name
}

func (p *Plugin) clearRecords(addr string, scope string) {
	name := p.getSlotName(addr, scope)
	p.recordsMutex.Lock()
	delete(p.records, name)
	p.recordsMutex.Unlock()
}

func (p *Plugin) getNextRecord(addr string, scope string) (nextaction, bool) {
	name := p.getSlotName(addr, scope)
	p.recordsMutex.RLock()
	if records, exist := p.records[name]; exist {
		record := records[0]
		p.recordsMutex.RUnlock()
		return record, true

	}
	p.recordsMutex.RUnlock()
	return nextaction{}, false
}

func (p *Plugin) addNextRecord(addr string, scope string, identifier string, sequence *int, templateName string, ctx interface{}, correlationId string, trackingId string) {

	name := p.getSlotName(addr, scope)

	seq := *sequence

	record := nextaction{
		identifier:    identifier,
		sequence:      seq,
		Context:       ctx,
		TemplateName:  templateName,
		CorrelationId: correlationId,
		TrackingId:    trackingId,
	}

	*sequence++

	p.recordsMutex.Lock()
	p.records[name] = append(p.records[name], record)
	p.recordsMutex.Unlock()
}

func (p *Plugin) dropRecord(addr string, scope string) bool {
	name := p.getSlotName(addr, scope)
	p.recordsMutex.Lock()
	if records, exist := p.records[name]; exist {
		records = append(records[:0], records[0+1:]...)
		p.records[name] = records
		p.recordsMutex.Unlock()
		return len(records) > 0
	}
	p.recordsMutex.Unlock()
	return false
}
