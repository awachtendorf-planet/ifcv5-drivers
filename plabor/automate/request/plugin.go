package request

import (
	"fmt"
	"strings"
	"sync"

	plabor "github.com/weareplanet/ifcv5-drivers/plabor/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	"github.com/weareplanet/ifcv5-main/log"
	helper "github.com/weareplanet/ifcv5-main/utils/state"

	"github.com/weareplanet/ifcv5-drivers/plabor/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/plabor/template"

	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*dispatcher.Plugin
	state        *helper.State
	driver       *plabor.Dispatcher
	automate     *dispatcher.Automate
	records      map[string][]interface{}
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
func New(parent *plabor.Dispatcher) *Plugin {
	return &Plugin{driver: parent}
}

const (
	retryDelay      = plabor.RetryDelay
	packetTimeout   = plabor.PacketTimeout
	nextActionDelay = plabor.NextActionDelay
	maxError        = plabor.MaxError
	maxNetworkError = plabor.MaxNetworkError
	pmsTimeOut      = plabor.PmsTimeout

	initial     = automatestate.Initial
	idle        = automatestate.Idle
	busy        = automatestate.Busy
	packetSent  = automatestate.PacketSent
	commandSent = automatestate.CommandSent
	nextAction  = automatestate.NextAction
	success     = automatestate.Success
	shutdown    = automatestate.Shutdown
)

// Init ...
func (p *Plugin) Init(automate *dispatcher.Automate) {
	log.Info().Msgf("init %T", p)
	p.records = make(map[string][]interface{})
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

func (p *Plugin) getSlotName(addr string, context string) string {
	return fmt.Sprintf("%s#%s", addr, context)
}

func (p *Plugin) clearRecords(addr string, context string) {
	name := p.getSlotName(addr, context)
	p.recordsMutex.Lock()
	delete(p.records, name)
	p.recordsMutex.Unlock()
}

func (p *Plugin) getNextRecord(addr string, context string) interface{} {
	name := p.getSlotName(addr, context)
	p.recordsMutex.RLock()
	if records, exist := p.records[name]; exist {
		record := records[0]
		p.recordsMutex.RUnlock()
		return record

	}
	p.recordsMutex.RUnlock()
	return nil
}

func (p *Plugin) addNextRecord(addr string, context string, record interface{}) {
	name := p.getSlotName(addr, context)
	p.recordsMutex.Lock()
	p.records[name] = append(p.records[name], record)
	p.recordsMutex.Unlock()
}

func (p *Plugin) dropRecord(addr string, context string) bool {
	name := p.getSlotName(addr, context)
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

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		return err
	}

	err = p.automate.SendPacket(addr, packet, action)

	return err
}

func (p *Plugin) getField(packet *ifc.LogicalPacket, field string) string {

	if packet == nil {
		return ""
	}

	data := packet.Data()
	value := cast.ToString(data[field])
	value = strings.Trim(value, " ")

	return value
}
