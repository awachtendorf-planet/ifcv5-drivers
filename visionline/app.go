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

	"github.com/weareplanet/ifcv5-drivers/visionline/automate/keyservice"
	"github.com/weareplanet/ifcv5-drivers/visionline/automate/linkcontrol"
	"github.com/weareplanet/ifcv5-drivers/visionline/template"

	drv "github.com/weareplanet/ifcv5-drivers/visionline/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparserextended"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:visionline:1"
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
	parser.Handler.ProcessOutgoingPacket = dispatcher.ProcessOutgoingPacket
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

	// incoming cmd socket
	parser.RegisterIncomingTemplate(1, template.PacketCodeCard, &template.TplCodeCardSocket{})
	parser.RegisterIncomingTemplate(1, template.PacketCheckout, &template.TplCheckoutSocket{})
	parser.RegisterIncomingTemplate(1, template.PacketReadKey, &template.TplReadCardSocket{})
	parser.RegisterIncomingTemplate(1, template.PacketAlive, &template.TplAliveSocket{})
	parser.RegisterIncomingTemplate(1, template.PacketCardUpdate, &template.TplCardUpdateSocket{})

	// unknown packet socket
	parser.RegisterIncomingTemplate(1, template.PacketGeneric, &template.TplGenericPacketSocket{})

	// incoming garbage socket
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbageSocket{}, 2)
	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(1, template.PacketEnq, &template.TplENQ{})

	// outgoing socket
	parser.RegisterOutgoingTemplate(1, template.PacketCodeCard, &template.TplCodeCardSocket{})
	parser.RegisterOutgoingTemplate(1, template.PacketCheckout, &template.TplCheckoutSocket{})
	parser.RegisterOutgoingTemplate(1, template.PacketReadKey, &template.TplReadCardSocket{})
	parser.RegisterOutgoingTemplate(1, template.PacketAlive, &template.TplAliveSocket{})
	parser.RegisterOutgoingTemplate(1, template.PacketCardUpdate, &template.TplCardUpdateSocket{})

	// low level incoming serial
	parser.RegisterIncomingTemplate(2, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(2, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(2, template.PacketEnq, &template.TplENQ{})

	// incoming serial
	parser.RegisterIncomingTemplate(2, template.PacketCodeCard, &template.TplCodeCardSerial{})
	parser.RegisterIncomingTemplate(2, template.PacketCheckout, &template.TplCheckoutSerial{})
	parser.RegisterIncomingTemplate(2, template.PacketReadKey, &template.TplReadCardSerial{})
	parser.RegisterIncomingTemplate(2, template.PacketAlive, &template.TplAliveSerial{})
	parser.RegisterIncomingTemplate(2, template.PacketCardUpdate, &template.TplCardUpdateSerial{})

	// unknown packet serial
	parser.RegisterIncomingTemplate(2, template.PacketGeneric, &template.TplGenericPacketSerial{})

	// incoming garbage serial
	parser.RegisterIncomingTemplate(2, template.PacketUnknown, &template.TplUnknownPacketSerial{})
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplGarbageSerial_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplGarbageSerial_2{}, 1)

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// low level outgoing
	parser.RegisterOutgoingTemplate(2, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(2, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(2, template.PacketEnq, &template.TplENQ{})

	// outgoing serial
	parser.RegisterOutgoingTemplate(2, template.PacketCodeCard, &template.TplCodeCardSerial{})
	parser.RegisterOutgoingTemplate(2, template.PacketCheckout, &template.TplCheckoutSerial{})
	parser.RegisterOutgoingTemplate(2, template.PacketReadKey, &template.TplReadCardSerial{})
	parser.RegisterOutgoingTemplate(2, template.PacketAlive, &template.TplAliveSerial{})
	parser.RegisterOutgoingTemplate(2, template.PacketCardUpdate, &template.TplCardUpdateSerial{})

	// sort templates
	parser.Ready()
}

func registerPMSAutomates(parent *dispatcher.Dispatcher) {

	if parent == nil {
		return
	}

	automate := dispatcher.NewAutomate(parent)
	automate.AddPlugin(pmsautomate.NewKeyService())

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
	automate.AddPlugin(keyservice.New(parent))

}
