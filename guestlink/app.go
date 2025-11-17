package main

import (
	"os"
	"time"

	"github.com/weareplanet/ifcv5-main/app"
	"github.com/weareplanet/ifcv5-main/config"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	pmsautomate "github.com/weareplanet/ifcv5-main/ifc/generic/pms/automate/base"
	pms "github.com/weareplanet/ifcv5-main/ifc/generic/pms/network"
	"github.com/weareplanet/ifcv5-main/ifc/noise/bridge"
	"github.com/weareplanet/ifcv5-main/ifc/noise/builder"
	noise "github.com/weareplanet/ifcv5-main/ifc/noise/const"
	"github.com/weareplanet/ifcv5-main/ifc/noise/plugin/access"
	"github.com/weareplanet/ifcv5-main/ifc/noise/plugin/api"
	devicemanager "github.com/weareplanet/ifcv5-main/ifc/noise/plugin/device/ifc"
	"github.com/weareplanet/ifcv5-main/ifc/noise/transport"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/noise/network"
	"github.com/weareplanet/ifcv5-main/slog"

	"github.com/weareplanet/ifcv5-drivers/guestlink/automate/datachange"
	"github.com/weareplanet/ifcv5-drivers/guestlink/automate/linkcontrol"
	requestitems "github.com/weareplanet/ifcv5-drivers/guestlink/automate/request_items"
	requestsimple "github.com/weareplanet/ifcv5-drivers/guestlink/automate/request_simple"
	"github.com/weareplanet/ifcv5-drivers/guestlink/automate/router"
	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	"github.com/weareplanet/ifcv5-drivers/guestlink/template/otrum"
	"github.com/weareplanet/ifcv5-drivers/guestlink/template/sonifi"
	"github.com/weareplanet/ifcv5-main/ifc/automate/inhousesync"

	drv "github.com/weareplanet/ifcv5-drivers/guestlink/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparserextended"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:guestlink:1"
)

// App represents a interface app
type App struct {
	start      time.Time
	parser     *analyser.GenericProtocol
	driver     *ifc.Network
	router     *analyser.Router
	dispatcher *drv.Dispatcher
	network    *network.Network
	logStream  string
	driverAddr string
}

// New returns a new interface app
func New(configFile string) (*App, error) {

	cfg, err := config.Read(noise.Ifc, configFile)
	if err != nil {
		return nil, err
	}

	if len(cfg.DriverAddr) == 0 {
		cfg.DriverAddr = driverAddr
	}

	if len(cfg.MetricAddr) > 0 {
		ifc.InitMetric(cfg.MetricAddr, cfg.DriverAddr)
	}

	if len(cfg.ConfigPath) > 0 {
		dispatcher.ConfigPath = cfg.ConfigPath
	}

	noiseConfig := builder.Config{
		Access: &access.Config{
			Alias: "",
			Token: cfg.Token,
			Mode:  noise.Ifc,
		},
		Api: &api.Config{
			Address: cfg.APIAddr,
			Mode:    noise.Ifc,
			Store:   cfg.Store,
		},
		Address: cfg.ListenAddr,
		Peers:   cfg.Peers,
		NAT:     cfg.NAT,
		TLS: &transport.TLS{
			Certificate:        cfg.TLS.Certificate,
			Key:                cfg.TLS.Key,
			InsecureSkipVerify: cfg.TLS.Insecure,
			ServerName:         cfg.TLS.ServerName,
			CAs:                cfg.TLS.RootCAs,
		},
	}

	// new noise network
	network, err := builder.NewIfcNetwork(&noiseConfig)
	if err != nil {
		return nil, err
	}

	// new protocol parser
	parser := analyser.NewParser()

	// new packet router
	router := analyser.NewRouter()

	// new ifc network <-> noise bridge
	bridge := bridge.NewBridge()

	// register bridge devicemanager
	plugin, registered := network.Plugin(devicemanager.PluginID)
	if registered {
		dm := plugin.(*devicemanager.Plugin)
		bridge.RegisterDeviceManager(dm)
	}

	// new ifc network
	driver := ifc.NewNetwork(cfg.DriverAddr, nil)
	driver.RegisterParser(parser)
	driver.RegisterRouter(router)
	driver.RegisterBridge(bridge)

	// new driver dispatcher
	dispatcher := drv.NewDispatcher(driver)

	// config debug settings
	if cfg.Debug.DisableDriverJobStorage {
		dispatcher.DisableDriverJobStorage()
	}

	if cfg.Debug.DisablePMSJobStorage {
		dispatcher.DisablePMSJobStorage()
	}

	if cfg.Debug.StoreJobsOnShutdownOnly {
		dispatcher.StoreJobsOnShutdown()
	}

	// set callbacks for incoming/outgoing packets (encoding/lrc/dump)
	parser.Handler.GetSlot = dispatcher.GetParserSlot
	parser.Handler.GetAdditionalLRCSize = dispatcher.GetAdditionalLRCSize
	parser.Handler.ProcessIncomingLRC = dispatcher.ProcessIncomingLRC
	parser.Handler.ProcessOutgoingLRC = dispatcher.ProcessOutgoingLRC
	parser.Handler.ProcessIncomingPacket = dispatcher.ProcessIncomingPacket
	parser.Handler.LogIncomingData = dispatcher.LogDriverIncomingBytesSimple
	parser.Handler.LogOutgoingData = dispatcher.LogDriverOutgoingBytesSimple
	parser.Handler.LogIncomingPacket = dispatcher.LogDriverIncomingPacket
	parser.Handler.LogOutgoingPacket = dispatcher.LogDriverOutgoingPacket

	// register automates
	registerDriverAutomates(dispatcher)
	registerPMSAutomates(dispatcher.Dispatcher)

	// can not directly created by dispatcher because of package cycle
	pms.NewNetwork(dispatcher.Dispatcher, nil)

	// register dispatcher for API callbacks
	plugin, registered = network.Plugin(api.PluginID)
	if registered {
		api := plugin.(*api.Plugin)
		api.RegisterDispatcher(dispatcher.Dispatcher)
		api.RegisterControlHandler(dispatcher.ControlHandler)
	}

	return &App{
		start:      time.Now(),
		parser:     parser,
		driver:     driver,
		router:     router,
		dispatcher: dispatcher,
		network:    network,
		logStream:  cfg.LogStream,
		driverAddr: cfg.DriverAddr,
	}, nil
}

// Startup ...
func (a *App) Startup() error {

	if len(a.logStream) > 0 {
		slog.Connect(a.logStream)
		slog.Meta(a.driverAddr, config.VersionString())
	}

	log.Info().Msgf("startup pid: %d as '%s' version '%s'", os.Getpid(), a.driverAddr, config.VersionString())
	slog.Info(slog.Ctx{Origin: "ifc"}).Msg("interface startup")

	// here we go
	registerTemplates(a.parser)
	go a.network.Listen()
	a.network.BlockUntilListening()
	a.dispatcher.Startup()
	a.network.SetReady()

	return nil
}

// InitLogRotate ...
func (a *App) InitLogRotate(_ func()) {

}

// Close ...
func (a *App) Close() {

	slog.Info(slog.Ctx{Origin: "ifc"}).Msg("interface shutdown")

	// be friendly, say goodbye
	plugin, registered := a.network.Plugin(access.PluginID)
	if registered {
		access := plugin.(*access.Plugin)
		access.Shutdown()
	}

	// shutdown noise network
	if a.network != nil {
		a.network.Close()
	}

	// shutdown driver dispatcher
	if a.dispatcher != nil {
		a.dispatcher.Close()
	}

	// shutdown ifc network
	if a.driver != nil {
		a.driver.Close()
	}

	// shutdown packet router
	if a.router != nil {
		a.router.Close()
	}

	log.Info().Msgf("done, uptime %s", time.Since(a.start))

}

func registerTemplates(parser *analyser.GenericProtocol) {
	if parser == nil {
		return
	}

	// vorbereitung für
	// maginet enhanced, 6stellige AccountNumber, Post Description 20stellig
	// oncommand abweichendes GuestMessage handling
	registerOtrumTemplates(parser)

	// sort templates
	parser.Ready()
}

func registerOtrumTemplates(parser *analyser.GenericProtocol) {
	if parser == nil {
		return
	}

	// low level incoming
	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(1, template.PacketEnq, &template.TplENQ{})

	// incoming
	parser.RegisterIncomingTemplate(1, template.PacketVerify, &otrum.TplVerifyPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &otrum.TplErrorPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketStart, &otrum.TplStartPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketTest, &otrum.TplTestPacket{})

	//
	parser.RegisterIncomingTemplate(1, template.PacketInit, &otrum.TplInitPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomStatus, &otrum.TplRoomStatusPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketGuestMessageRead, &otrum.TplGuestMessageReadPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupResult, &otrum.TplWakeupResult{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupSet, &otrum.TplWakeupSetPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupClear, &otrum.TplWakeupClearPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketDisplayRequest, &otrum.TplDisplayRequest{})
	parser.RegisterIncomingTemplate(1, template.PacketLookupRequest, &otrum.TplLookupRequest{})
	parser.RegisterIncomingTemplate(1, template.PacketStatusRequest, &otrum.TplStatusRequest{})
	parser.RegisterIncomingTemplate(1, template.PacketGuestMessageRequest, &otrum.TplGuestMessageRequestPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketPostCharge, &otrum.TplPostChargePacket{})
	parser.RegisterIncomingTemplate(1, template.PacketPostCharge, &sonifi.TplPostChargeEnhancedPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketCheckoutRequest, &otrum.TplCheckoutRequest{})
	parser.RegisterIncomingTemplate(1, template.PacketCheckoutRequest, &sonifi.TplCheckoutEnhancedPacket{})

	// unknown or unsupported command
	parser.RegisterIncomingTemplate(1, template.PacketUnknownCommand, &template.TplUnknownCommand{})

	// incoming garbage
	parser.RegisterIncomingTemplate(1, template.PacketUnknown, &template.TplUnknownPacket{})
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbage_Framing_2{}, 1)

	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// low level outgoing
	parser.RegisterOutgoingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(1, template.PacketEnq, &template.TplENQ{})

	// outgoing
	parser.RegisterOutgoingTemplate(1, template.PacketVerify, &otrum.TplVerifyPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketError, &otrum.TplErrorPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketStart, &otrum.TplStartPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketTest, &otrum.TplTestPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketHelo, &sonifi.TplNameHeloPacket{})

	//
	parser.RegisterOutgoingTemplate(1, template.PacketCheckIn, &otrum.TplCheckinPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketCheckOut, &otrum.TplCheckoutPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketWakeupSet, &otrum.TplWakeupSetPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketWakeupClear, &otrum.TplWakeupClearPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketGuestMessageStatus, &otrum.TplGuestMessageStatusPacket{})

	parser.RegisterOutgoingTemplate(1, template.PacketNameReply, &otrum.TplNamePacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketInfoReply, &otrum.TplInfoPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketItemReply, &otrum.TplItemPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketBalanceReply, &otrum.TplBalancePacket{})

	parser.RegisterOutgoingTemplate(1, template.PacketGuestMessageHeader, &otrum.TplGuestMessageHeaderPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketGuestMessageCaller, &otrum.TplGuestMessageCallerPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketGuestMessageText, &otrum.TplGuestMessageTextPacket{})
}

func registerPMSAutomates(parent *dispatcher.Dispatcher) {

	if parent == nil {
		return
	}

	automate := dispatcher.NewAutomate(parent)
	automate.AddPlugin(pmsautomate.NewASW())

	automate = dispatcher.NewAutomate(parent)
	automate.AddPlugin(pmsautomate.NewWakeup())

	automate = dispatcher.NewAutomate(parent)
	automate.AddPlugin(pmsautomate.NewRoomStatus())

	automate = dispatcher.NewAutomate(parent)
	automate.AddPlugin(pmsautomate.NewSysAdmin())

	automate = dispatcher.NewAutomate(parent)
	automate.AddPlugin(pmsautomate.NewTelemetry())
}

func registerDriverAutomates(parent *drv.Dispatcher) {

	if parent == nil {
		return
	}

	automate := dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(router.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(linkcontrol.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(datachange.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(requestsimple.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(requestitems.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(inhousesync.New(7))

}
