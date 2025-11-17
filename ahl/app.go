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

	"github.com/weareplanet/ifcv5-drivers/ahl/automate/datachange"
	"github.com/weareplanet/ifcv5-drivers/ahl/automate/linkcontrol"
	"github.com/weareplanet/ifcv5-drivers/ahl/automate/request"
	"github.com/weareplanet/ifcv5-drivers/ahl/template"
	"github.com/weareplanet/ifcv5-main/ifc/automate/inhousesync"

	drv "github.com/weareplanet/ifcv5-drivers/ahl/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparserextended"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:ahl:1"
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
	parser.Handler.ProcessIncomingLRC = dispatcher.ProcessIncomingLRC
	parser.Handler.ProcessOutgoingLRC = dispatcher.ProcessOutgoingLRC
	parser.Handler.ProcessIncomingPacket = dispatcher.ProcessIncomingPacket
	parser.Handler.LogIncomingData = dispatcher.LogDriverIncomingBytes
	parser.Handler.LogOutgoingData = dispatcher.LogDriverOutgoingBytes
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

	// incoming 4400 - 5
	parser.RegisterIncomingTemplate(drv.AHL4400_5, template.PacketReply, &template.TplReply_4400_5{})
	parser.RegisterIncomingTemplate(drv.AHL4400_5, template.PacketCallPacket, &template.TplCallPacket_4400_5{})
	parser.RegisterIncomingTemplate(drv.AHL4400_5, template.PacketRoomStatus, &template.TplRoomStatus_4400_5{})
	parser.RegisterIncomingTemplate(drv.AHL4400_5, template.PacketWakeupEvent, &template.TplWakeupEvent_4400_5{})
	parser.RegisterIncomingTemplate(drv.AHL4400_5, template.PacketDataTransfer, &template.TplDataTransfer_4400_5{})
	parser.RegisterIncomingTemplate(drv.AHL4400_5, template.PacketVoiceMailEvent, &template.TplVoiceMailEvent_4400_5{})

	// incoming 4400 - 8
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketReply, &template.TplReply_4400_8{})
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketCallPacket, &template.TplCallPacket_4400_8{})
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketCallPacketExtended, &template.TplCallPacket_4400_8_Extended{})
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketRoomStatus, &template.TplRoomStatus_4400_8{})
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketWakeupEvent, &template.TplWakeupEvent_4400_8{})
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketDataTransfer, &template.TplDataTransfer_4400_8{})
	parser.RegisterIncomingTemplate(drv.AHL4400_8, template.PacketVoiceMailEvent, &template.TplVoiceMailEvent_4400_8{})

	// outgoing 4400 - 5
	parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketCheckIn, &template.TplCheckIn_4400_5{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketCheckOut, &template.TplCheckOut_4400_5{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketDataChange, &template.TplDataChange_4400_5{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketWakeupSet, &template.TplDataChange_4400_5{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketWakeupClear, &template.TplDataChange_4400_5{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketRoomStatus, &template.TplDataChange_4400_5{})
	// parser.RegisterOutgoingTemplate(drv.AHL4400_5, template.PacketVoiceMailEvent, &template.TplVoiceMailEvent_4400_5{})

	// outgoing 4400 - 8
	parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketCheckIn, &template.TplCheckIn_4400_8{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketCheckOut, &template.TplCheckOut_4400_8{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketDataChange, &template.TplDataChange_4400_8{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketWakeupSet, &template.TplDataChange_4400_8{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketWakeupClear, &template.TplDataChange_4400_8{})
	parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketRoomStatus, &template.TplDataChange_4400_8{})
	// parser.RegisterOutgoingTemplate(drv.AHL4400_8, template.PacketVoiceMailEvent, &template.TplVoiceMailEvent_4400_8{})

	registerDefaultTemplate(parser, drv.AHL4400_5)
	registerDefaultTemplate(parser, drv.AHL4400_8)

	// sort templates
	parser.Ready()
}

func registerDefaultTemplate(parser *analyser.GenericProtocol, slot uint) {

	if parser == nil {
		return
	}

	// low level incoming
	parser.RegisterIncomingTemplate(slot, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(slot, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(slot, template.PacketEnq, &template.TplENQ{})

	parser.RegisterIncomingTemplate(slot, template.PacketLinkAlive, &template.TplLinkAliveIncoming{})

	// incoming garbage, protocol printer output
	parser.RegisterIncomingTemplate(slot, template.PacketUnknown, &template.TplUnknownPacket{})
	parser.RegisterIncomingTemplate(slot, template.PacketLogOutput, &template.TplLogOutputIncoming1{})
	parser.RegisterIncomingTemplate(slot, template.PacketLogOutput, &template.TplLogOutputIncoming2{})

	parser.RegisterIncomingTemplateWithOption(slot, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(slot, template.PacketGarbage, &template.TplGarbage_Framing_2{}, 1)

	parser.RegisterIncomingTemplateWithOption(slot, template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(slot, template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(slot, template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// low level outgoing
	parser.RegisterOutgoingTemplate(slot, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(slot, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(slot, template.PacketEnq, &template.TplENQ{})

	parser.RegisterOutgoingTemplate(slot, template.PacketLinkStart, &template.TplLinkStartOutgoing{})
	parser.RegisterOutgoingTemplate(slot, template.PacketLinkAlive, &template.TplLinkAliveOutgoing{})
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
	automate.AddPlugin(pmsautomate.NewPosting())

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
	automate.AddPlugin(linkcontrol.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(datachange.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(request.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(inhousesync.New(7))

}
