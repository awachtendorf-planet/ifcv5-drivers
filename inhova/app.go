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

	"github.com/weareplanet/ifcv5-drivers/inhova/template"

	"github.com/weareplanet/ifcv5-drivers/inhova/automate/keyservice"

	linkcontrol "github.com/weareplanet/ifcv5-main/ifc/automate/simplelinkcontrol"

	drv "github.com/weareplanet/ifcv5-drivers/inhova/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparser"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:inhova:1"
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
	parser.Handler.GetAdditionalLRCSize = dispatcher.GetAdditionalLRCSize
	parser.Handler.ProcessIncomingLRC = dispatcher.ProcessIncomingLRC
	parser.Handler.ProcessOutgoingLRC = dispatcher.ProcessOutgoingLRC
	parser.Handler.ProcessIncomingPacket = dispatcher.ProcessIncomingPacket
	parser.Handler.LogIncomingData = dispatcher.LogDriverIncomingBytes
	parser.Handler.LogOutgoingData = dispatcher.LogDriverOutgoingBytes
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

	// incoming
	parser.RegisterIncomingTemplate(template.PacketCodeCardAnswer, &template.TplCodeAnswer{})
	parser.RegisterIncomingTemplate(template.PacketCodeCardAnswer, &template.TplCodeAnswerCO{})
	parser.RegisterIncomingTemplate(template.PacketError, &template.TplErrorPacket{})

	// incoming low level
	parser.RegisterIncomingTemplate(template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(template.PacketEnq, &template.TplENQ{})

	// incoming garbage
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbage, &template.TplGarbage_Framing_2{}, 1)
	parser.RegisterIncomingTemplate(template.PacketUnknown, &template.TplUnknownPacket{})

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// outgoing
	parser.RegisterOutgoingTemplate(template.PacketKeyRequest, &template.TplCodeCard{})
	parser.RegisterOutgoingTemplate(template.PacketKeyDelete, &template.TplDeleteCard{})

	// outgoing low level
	parser.RegisterOutgoingTemplate(template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(template.PacketNak, &template.TplNAK{})

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

		// pms automates
		pmsautomate.NewASW(),
		pmsautomate.NewKeyService(),
		pmsautomate.NewSysAdmin(),
		pmsautomate.NewTelemetry(),
	)

}
