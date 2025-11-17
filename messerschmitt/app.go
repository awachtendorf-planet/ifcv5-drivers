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

	"github.com/weareplanet/ifcv5-drivers/messerschmitt/template"

	"github.com/weareplanet/ifcv5-drivers/messerschmitt/automate/keyservice"
	"github.com/weareplanet/ifcv5-drivers/messerschmitt/automate/request"
	linkcontrol "github.com/weareplanet/ifcv5-main/ifc/automate/simplelinkcontrol"

	drv "github.com/weareplanet/ifcv5-drivers/messerschmitt/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparserextended"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:messerschmitt:1"
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
	registerAutomates(dispatcher)

	// can not directly created by dispatcher because of package cycle
	pms.NewNetwork(dispatcher.Dispatcher, nil)

	// register dispatcher for API callbacks
	plugin, registered = network.Plugin(api.PluginID)
	if registered {
		api := plugin.(*api.Plugin)
		api.RegisterDispatcher(dispatcher.Dispatcher)
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

	// SLOT 1 = serial
	// SLOT 2 = socket

	// Serial

	// incoming
	parser.RegisterIncomingTemplate(1, template.PacketSynchronisation, &template.SerialSynchronisation1{})
	parser.RegisterIncomingTemplate(1, template.PacketSynchronisation, &template.SerialSynchronisation2{})
	parser.RegisterIncomingTemplate(1, template.PacketCommandAcknowledge, &template.SerialCommandAcknowledge1{})
	parser.RegisterIncomingTemplate(1, template.PacketCommandAcknowledge, &template.SerialCommandAcknowledge2{})
	parser.RegisterIncomingTemplate(1, template.PacketFunctionCards, &template.SerialFunctionCards1{})
	parser.RegisterIncomingTemplate(1, template.PacketFunctionCards, &template.SerialFunctionCards2{})
	parser.RegisterIncomingTemplate(1, template.PacketSpecialReaders, &template.SerialSpecialReaders1{})
	parser.RegisterIncomingTemplate(1, template.PacketSpecialReaders, &template.SerialSpecialReaders2{})
	parser.RegisterIncomingTemplate(1, template.PacketVendingMachines, &template.SerialVendingMachines1{})
	parser.RegisterIncomingTemplate(1, template.PacketVendingMachines, &template.SerialVendingMachines2{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomEquipment, &template.SerialRoomEquipment1{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomEquipment, &template.SerialRoomEquipment2{})

	// incoming low level
	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.SerialACK0{})
	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.SerialACK1{})
	parser.RegisterIncomingTemplate(1, template.PacketNak, &template.SerialNAK{})
	parser.RegisterIncomingTemplate(1, template.PacketEnq, &template.SerialENQ{})
	parser.RegisterIncomingTemplate(1, template.PacketEOT, &template.SerialEOT{})
	parser.RegisterIncomingTemplate(1, template.PacketWACK, &template.SerialWACK{})
	parser.RegisterIncomingTemplate(1, template.PacketTTD, &template.SerialTTD{})

	// incoming garbage
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.SerialGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.SerialGarbage_Framing_2{}, 1)
	parser.RegisterIncomingTemplate(1, template.PacketUnknown, &template.SerialUnknownPacket1{}) // STX ... ETX
	parser.RegisterIncomingTemplate(1, template.PacketUnknown, &template.SerialUnknownPacket2{}) // STX ... ETB

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.SerialGarbage_ACK0{}, 2)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.SerialGarbage_ACK1{}, 2)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.SerialGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.SerialGarbage_ENQ{}, 1)

	// outgoing
	parser.RegisterOutgoingTemplate(1, template.PacketSynchronisation, &template.SerialSynchronisationReply{})
	parser.RegisterOutgoingTemplate(1, template.PacketKeyRequest, &template.SerialKeyCreate{})
	parser.RegisterOutgoingTemplate(1, template.PacketKeyDelete, &template.SerialKeyDelete{})

	// outgoing low level
	parser.RegisterOutgoingTemplate(1, template.PacketAck0, &template.SerialACK0{})
	parser.RegisterOutgoingTemplate(1, template.PacketAck1, &template.SerialACK1{})
	parser.RegisterOutgoingTemplate(1, template.PacketAck, &template.SerialACK1{})
	parser.RegisterOutgoingTemplate(1, template.PacketNak, &template.SerialNAK{})
	parser.RegisterOutgoingTemplate(1, template.PacketEnq, &template.SerialENQ{})
	parser.RegisterOutgoingTemplate(1, template.PacketEOT, &template.SerialEOT{})

	// Socket

	// incoming
	parser.RegisterIncomingTemplate(2, template.PacketSynchronisation, &template.SocketSynchronisation{})
	parser.RegisterIncomingTemplate(2, template.PacketCommandAcknowledge, &template.SocketCommandAcknowledge{})
	parser.RegisterIncomingTemplate(2, template.PacketFunctionCards, &template.SocketFunctionCards{})
	parser.RegisterIncomingTemplate(2, template.PacketSpecialReaders, &template.SocketSpecialReaders{})

	// outgoing
	parser.RegisterOutgoingTemplate(2, template.PacketSynchronisation, &template.SocketSynchronisationReply{})
	parser.RegisterOutgoingTemplate(2, template.PacketKeyRequest, &template.SocketKeyCreate{})
	parser.RegisterOutgoingTemplate(2, template.PacketKeyDelete, &template.SocketKeyDelete{})

	// sort templates
	parser.Ready()
}

func registerAutomates(parent *drv.Dispatcher) {

	if parent == nil || parent.Dispatcher == nil {
		return
	}

	parent.Dispatcher.AddAutomates(

		// driver automates
		linkcontrol.New(nil),
		keyservice.New(parent),
		request.New(parent),

		// pms automates
		pmsautomate.NewKeyService(),
		pmsautomate.NewRoomStatus(),
		pmsautomate.NewPosting(),
		pmsautomate.NewSysAdmin(),
		pmsautomate.NewTelemetry(),
	)

}
