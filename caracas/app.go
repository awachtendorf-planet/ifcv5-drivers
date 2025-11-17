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

	"github.com/weareplanet/ifcv5-drivers/caracas/template"

	"github.com/weareplanet/ifcv5-drivers/caracas/automate/datachange"

	"github.com/weareplanet/ifcv5-drivers/caracas/automate/dbsync"
	"github.com/weareplanet/ifcv5-drivers/caracas/automate/linkcontrol"
	"github.com/weareplanet/ifcv5-drivers/caracas/automate/request"
	"github.com/weareplanet/ifcv5-main/ifc/automate/inhousesync"

	drv "github.com/weareplanet/ifcv5-drivers/caracas/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparser"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:caracas:1"
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

	// new ifc network &lt;-> noise bridge
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
	parser.Handler.GetAdditionalLRCSize = dispatcher.GetAdditionalLRCSize
	parser.Handler.ProcessIncomingLRC = dispatcher.ProcessIncomingLRC
	parser.Handler.ProcessOutgoingLRC = dispatcher.ProcessOutgoingLRC
	parser.Handler.ProcessIncomingPacket = dispatcher.ProcessIncomingPacket
	parser.Handler.ProcessOutgoingPacket = dispatcher.ProcessOutgoingPacket
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

	// incoming
	// HERE WE GO
	parser.RegisterIncomingTemplate(template.PacketCallCharge, &template.CallChargePacket{})
	parser.RegisterIncomingTemplate(template.PacketMinibarCharge, &template.MinibarPostingPacket{})
	parser.RegisterIncomingTemplate(template.PacketMinibarChargeAdvanced, &template.AdvancedMinibarPostingPacket{})
	parser.RegisterIncomingTemplate(template.PacketMinibarChargeComplete, &template.CompleteMinibarPostingPacket{})

	parser.RegisterIncomingTemplate(template.PacketRoomstatus, &template.RoomstatusPacket{})
	parser.RegisterIncomingTemplate(template.PacketRoomstatusAdvanced, &template.AdvancedRoomstatusPacket{})

	parser.RegisterIncomingTemplate(template.PacketVoicemail, &template.VoiceMessagePacket{})
	parser.RegisterIncomingTemplate(template.PacketVoicemailAdvanced, &template.AdvancedVoiceMessagePacket{})

	parser.RegisterIncomingTemplate(template.PacketWakeupOrderIncoming, &template.IncomingWakeupOrderPacket{})
	parser.RegisterIncomingTemplate(template.PacketWakeupResult, &template.WakeupAttemptPacket{})

	parser.RegisterIncomingTemplate(template.PacketDBSyncRequest, &template.DBResyncRequestPacket{})

	parser.RegisterIncomingTemplate(template.PacketDBSyncRequest, &template.DBResyncRequestPacket{})

	// incoming low level
	parser.RegisterIncomingTemplate(template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(template.PacketEnq, &template.TplENQ{})
	parser.RegisterIncomingTemplate(template.PacketSyn, &template.TplSYN{})

	parser.RegisterIncomingTemplate(template.PacketBind, &template.BindPacket{})
	parser.RegisterIncomingTemplate(template.PacketDatetime, &template.DatetimePacket{})

	// incoming garbage
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbage, &template.TplGarbage_Framing_2{}, 1)
	parser.RegisterIncomingTemplate(template.PacketUnknown, &template.TplUnknownPacket{})

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// outgoing
	// HERE WE GO
	parser.RegisterOutgoingTemplate(template.PacketCheckIn, &template.CheckInPacket{})
	parser.RegisterOutgoingTemplate(template.PacketCheckOut, &template.CheckOutPacket{})
	parser.RegisterOutgoingTemplate(template.PacketUpdateCheckInName, &template.UpdateCheckInNamePacket{})
	parser.RegisterOutgoingTemplate(template.PacketDataChangeAdvanced, &template.AdvancedDataChangePacket{})

	parser.RegisterOutgoingTemplate(template.PacketMessageLight, &template.MessageLightPacket{})
	parser.RegisterOutgoingTemplate(template.PacketMessageLightAlternative, &template.AlternativeMessageLightPacket{})

	parser.RegisterOutgoingTemplate(template.PacketWakeupOrder, &template.WakeupOrderPacket{})

	parser.RegisterOutgoingTemplate(template.PacketDND, &template.DNDBRKPacket{})
	parser.RegisterOutgoingTemplate(template.PacketClassOfService, &template.DNDBRKPacket{})

	// outgoing low level
	parser.RegisterOutgoingTemplate(template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(template.PacketEot, &template.TplEOT{})

	parser.RegisterOutgoingTemplate(template.PacketBind, &template.BindPacket{})
	parser.RegisterOutgoingTemplate(template.PacketDatetime, &template.DatetimePacket{})

	// sort templates
	parser.Ready()
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
	automate.AddPlugin(pmsautomate.NewPosting())

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
	automate.AddPlugin(dbsync.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(datachange.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(request.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(inhousesync.New(7))

}
