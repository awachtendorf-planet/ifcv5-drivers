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

	"github.com/weareplanet/ifcv5-drivers/detewe/template"

	// "github.com/weareplanet/ifcv5-drivers/detewe/automate/datachange"
	"github.com/weareplanet/ifcv5-drivers/detewe/automate/datachange"
	"github.com/weareplanet/ifcv5-drivers/detewe/automate/linkcontrol"
	"github.com/weareplanet/ifcv5-drivers/detewe/automate/request"
	"github.com/weareplanet/ifcv5-main/ifc/automate/inhousesync"

	drv "github.com/weareplanet/ifcv5-drivers/detewe/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/byteparserstaged"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:detewe:1"
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

	prepareForSecondStage := func(slot uint, in []byte) []byte {

		if slot == 1 && len(in) > 2 {

			return append(in[:len(in)-2], byte(0x03))
		}
		if len(in) > 12 {
			return append(in[12:], byte(0x03))
		} else {
			return []byte{}
		}
	}

	// new protocol parser
	parser := analyser.NewParserWithOption(128*5, "StagedPayload", prepareForSecondStage)

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
	parser.Handler.PostProcessOutgoing = dispatcher.PostProcessOutgoing
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
	parser.RegisterIncomingTemplate(1, template.PacketSerialFrameIn, &template.TplSerialFrameInPacket{}, false)

	parser.RegisterIncomingTemplate(2, template.PacketLoginAnswer, &template.TplTCPLoginPositiveAnswerPacket{}, false)
	parser.RegisterIncomingTemplate(2, template.PacketLoginAnswer, &template.TplTCPLoginPositiveAnswerLenPacket{}, false)
	parser.RegisterIncomingTemplate(2, template.PacketTCPFrameIn, &template.TplTCPAnswerFramePacket{}, false)

	// incoming - stagedTemplates
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG41, &template.TplTg41AnswerPacket{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG67, &template.TplTg67AnswerPacket{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG71, &template.TplTg71AnswerPacket{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG60, &template.TplTg60AnswerPacket{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG70, &template.TplTg70AnswerPacket{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG80, &template.TplTg80AnswerPacket{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketAnswerTG00, &template.TplTestMessageAnswerPacket{}, true)

	parser.RegisterIncomingTemplate(1, template.PacketTG10, &template.TplTg10Packet{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketTG20, &template.TplTg20Packet{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketTG40, &template.TplTg40Packet{}, true)
	parser.RegisterIncomingTemplate(1, template.PacketTG72, &template.TplTg72Packet{}, true)

	// incoming low level
	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.TplACK{}, false)
	parser.RegisterIncomingTemplate(1, template.PacketNak, &template.TplNAK{}, false)
	parser.RegisterIncomingTemplate(1, template.PacketEnq, &template.TplENQ{}, false)

	// incoming garbage
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbage_Framing_1{}, false, 1)

	parser.RegisterIncomingTemplateWithOption(2, template.PacketTCPGarbage, &template.TplTCPGarbagePacket{}, false, 11)
	parser.RegisterIncomingTemplateWithOption(2, template.PacketLoginGarbage, &template.TplTCPLoginGarbagePacket{}, false, 13)

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, false, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, false, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, false, 1)

	// outgoing
	parser.RegisterOutgoingTemplate(2, template.PacketTCPLogin, &template.TplTCPLoginPacket{})

	// outgoing low level
	parser.RegisterOutgoingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(1, template.PacketNak, &template.TplNAK{})

	// outgoing - stagedTemplates

	parser.RegisterOutgoingTemplate(1, template.PacketTG41, &template.TplTg41Packet{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG60, &template.TplTg60Packet{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG67, &template.TplTg67Packet{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG70, &template.TplTg70Packet{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG71, &template.TplTg71Packet{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG80, &template.TplTg80Packet{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG00, &template.TplTestMessagePacket{})

	parser.RegisterOutgoingTemplate(1, template.PacketTG10Answer, &template.TplTg10AnswerPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG20Answer, &template.TplTg20AnswerPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketTG40Answer, &template.TplTg40AnswerPacket{})

	parser.RegisterOutgoingTemplate(2, template.PacketTG41, &template.TplTg41Packet{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG60, &template.TplTg60Packet{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG67, &template.TplTg67Packet{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG70, &template.TplTg70Packet{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG71, &template.TplTg71Packet{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG80, &template.TplTg80Packet{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG00, &template.TplTestMessagePacket{})

	parser.RegisterOutgoingTemplate(2, template.PacketTG10Answer, &template.TplTg10AnswerPacket{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG20Answer, &template.TplTg20AnswerPacket{})
	parser.RegisterOutgoingTemplate(2, template.PacketTG40Answer, &template.TplTg40AnswerPacket{})

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
