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

	"github.com/weareplanet/ifcv5-drivers/nortel/template"

	"github.com/weareplanet/ifcv5-drivers/nortel/automate/callpacket"
	"github.com/weareplanet/ifcv5-drivers/nortel/automate/datachange"
	"github.com/weareplanet/ifcv5-drivers/nortel/automate/request"
	"github.com/weareplanet/ifcv5-main/ifc/automate/inhousesync"
	linkcontrol "github.com/weareplanet/ifcv5-main/ifc/automate/simplelinkcontrol"

	drv "github.com/weareplanet/ifcv5-drivers/nortel/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/byteparseradvanced"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:nortel:1"
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
	parser := analyser.NewParserWithOption(256)

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

	// set parser callback functions
	parser.Handler.GetIncomingSlot = dispatcher.GetIncomingParserSlot
	parser.Handler.GetOutgoingSlot = dispatcher.GetOutgoingParserSlot
	parser.Handler.GetAdditionalLRCSize = dispatcher.GetAdditionalLRCSize

	parser.Handler.ProcessIncomingBytes = dispatcher.ProcessIncomingBytes // remove MSB if present
	parser.Handler.ProcessOutgoingBytes = dispatcher.ProcessOutgoingBytes // framing or bgd device

	parser.Handler.ProcessIncomingLRC = dispatcher.ProcessIncomingLRC // verify incoming LRC
	parser.Handler.ProcessOutgoingLRC = dispatcher.ProcessOutgoingLRC // calculate outgoing LRC

	parser.Handler.ProcessIncomingPacket = dispatcher.ProcessIncomingPacket // send ACK/NAK

	// reduce byte stream log output
	parser.Handler.LogIncomingData = dispatcher.LogIncomingBytes
	parser.Handler.LogOutgoingData = dispatcher.LogOutgoingBytes

	// advanced parser function
	parser.Handler.NeedMoreData = dispatcher.NeedMoreData

	// dashboard log
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

	// slot 1 incoming = management device framed
	// slot 1 outgoing = management device framed and background terminal mode
	// slot 2 incoming = management device background terminal mode
	// slot 3 incoming = call logging

	// call logging slot with threshold for single line handling
	parser.RegisterThreshold(3, 2*time.Second)

	// incoming management device variable extension length
	parser.RegisterIncomingTemplate(1, template.PacketPolling, &template.TplPolling{})
	parser.RegisterIncomingTemplate(1, template.PacketPolling, &template.TplTest{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomStatus, &template.TplRoomStatusAdvanced{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomStatus, &template.TplRoomStatus{})
	parser.RegisterIncomingTemplate(1, template.PacketMinibarItem, &template.TplMinibarItem{})
	parser.RegisterIncomingTemplate(1, template.PacketMinibarTotal, &template.TplMinibarTotal{})
	parser.RegisterIncomingTemplate(1, template.PacketVoiceCount, &template.TplVoiceCount{})

	parser.RegisterIncomingTemplate(1, template.PacketWakeupSet, &template.TplWakeupSet4{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupSet, &template.TplWakeupSet3{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupClear, &template.TplWakeupClear{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupAnswer, &template.TplWakeupAnswer{})

	// fixed extension length, extension are always 7 characters, right justified and left blank padded, does not work with variable templates
	// sollte das im selben slot zu Problemen führen, dann einen eigen Slot verwenden
	// und über einen Config Parameter steuern, GetIncomingParserSlot entsprechend anpassen
	parser.RegisterIncomingTemplate(1, template.PacketRoomStatus, &template.TplRoomStatusAdvanced_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomStatus, &template.TplRoomStatus_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketMinibarItem, &template.TplMinibarItem_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketMinibarTotal, &template.TplMinibarTotal_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketVoiceCount, &template.TplVoiceCount_FixedLength{})

	parser.RegisterIncomingTemplate(1, template.PacketWakeupSet, &template.TplWakeupSet4_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupSet, &template.TplWakeupSet3_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupClear, &template.TplWakeupClear_FixedLength{})
	parser.RegisterIncomingTemplate(1, template.PacketWakeupAnswer, &template.TplWakeupAnswer_FixedLength{})

	// bg terminal
	parser.RegisterIncomingTemplate(2, template.PacketRoomStatus, &template.TplBgdRoomStatusAdvanced{})
	parser.RegisterIncomingTemplate(2, template.PacketRoomStatus, &template.TplBgdRoomStatus{})
	parser.RegisterIncomingTemplate(2, template.PacketMinibarItem, &template.TplBgdMinibarItem{})
	parser.RegisterIncomingTemplate(2, template.PacketMinibarTotal, &template.TplBgdMinibarTotal{})
	parser.RegisterIncomingTemplate(2, template.PacketVoiceCount, &template.TplBgdVoiceCount{})

	parser.RegisterIncomingTemplate(2, template.PacketWakeupSet, &template.TplBgdWakeupSet4{})
	parser.RegisterIncomingTemplate(2, template.PacketWakeupSet, &template.TplBgdWakeupSet3{})
	parser.RegisterIncomingTemplate(2, template.PacketWakeupClear, &template.TplBgdWakeupClear{})
	parser.RegisterIncomingTemplate(2, template.PacketWakeupAnswer, &template.TplBgdWakeupAnswer{})

	parser.RegisterIncomingTemplate(3, template.PacketCallPacket, &template.TplCallPacketN{})
	parser.RegisterIncomingTemplate(3, template.PacketCallPacket, &template.TplCallPacketS{})
	parser.RegisterIncomingTemplate(3, template.PacketCallPacket, &template.TplCallPacketE{})
	parser.RegisterIncomingTemplate(3, template.PacketCallPacket, &template.TplCallPacketX{})
	parser.RegisterIncomingTemplate(3, template.PacketCallPacket, &template.TplCallPacketL{})
	parser.RegisterIncomingTemplate(3, template.PacketCallPacketSingle, &template.TplCallPacketNSingle{})

	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorMnemonic{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorInput{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorTryAgain{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorNameBig{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorDuplicate{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorNoData{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorNoSetCPNDData{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorNoCustCPNDData{})
	parser.RegisterIncomingTemplate(1, template.PacketError, &template.TplErrorNoCPNDMemory{})

	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorMnemonic{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorInput{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorTryAgain{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorNameBig{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorDuplicate{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorNoData{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorNoSetCPNDData{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorNoCustCPNDData{})
	parser.RegisterIncomingTemplate(2, template.PacketError, &template.TplBgdErrorNoCPNDMemory{})

	// incoming low level
	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(1, template.PacketEnq, &template.TplENQ{})

	// incoming garbage
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbage_Framing_2{}, 1)
	parser.RegisterIncomingTemplate(1, template.PacketUnknown, &template.TplUnknownPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketTermination, &template.TplCR{})
	parser.RegisterIncomingTemplate(1, template.PacketTermination, &template.TplCRSequence{})

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// bgd garbage
	parser.RegisterIncomingTemplate(2, template.PacketUnknown, &template.TplUnknownPacket{})     // STX ... ETX
	parser.RegisterIncomingTemplate(2, template.PacketGarbage, &template.TplGarbage_Framing_1{}) // ... STX
	parser.RegisterIncomingTemplate(2, template.PacketGarbage, &template.TplGarbage_Framing_2{}) // STX ... STX

	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplBgdGarbage_Overread_1a{}, 2) // ovr .. ST
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplBgdGarbage_Overread_1b{}, 2) // ovr .. MB
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplBgdGarbage_Overread_1c{}, 2) // ovr .. WA
	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplBgdGarbage_Overread_1d{}, 2) // ovr .. IS
	parser.RegisterIncomingTemplate(2, template.PacketUnknownBgd, &template.TplBgdGarbage_Overread_2{})            // ovr .. CR LF

	// callpacket device garbage
	parser.RegisterIncomingTemplateAsGarbage(3, template.PacketGarbage, &template.TplCallPacketOvr{}) // ... 0d0a

	// callpacket device management garbage
	parser.RegisterIncomingTemplateAsGarbage(3, template.PacketGarbage, &template.TplUnknownPacket{}) // STX ... ETX
	parser.RegisterIncomingTemplateAsGarbage(3, template.PacketGarbage, &template.TplACK{})
	parser.RegisterIncomingTemplateAsGarbage(3, template.PacketGarbage, &template.TplNAK{})
	parser.RegisterIncomingTemplateAsGarbage(3, template.PacketGarbage, &template.TplENQ{})

	// outgoing
	parser.RegisterOutgoingTemplate(1, template.PacketCheckInExtension, &template.TplCheckinExtension{})
	parser.RegisterOutgoingTemplate(1, template.PacketCheckOutExtension, &template.TplCheckoutExtension{})
	parser.RegisterOutgoingTemplate(1, template.PacketDisplayName, &template.TplSetDisplayName{})
	parser.RegisterOutgoingTemplate(1, template.PacketRoomStatus, &template.TplSetRoomStatus{})
	parser.RegisterOutgoingTemplate(1, template.PacketLanguage, &template.TplSetLanguage{})
	parser.RegisterOutgoingTemplate(1, template.PacketMessageLamp, &template.TplSetMessageLamp{})
	parser.RegisterOutgoingTemplate(1, template.PacketDoNotDisturb, &template.TplSetDoNotDisturb{})
	parser.RegisterOutgoingTemplate(1, template.PacketVipState, &template.TplSetVipState{})
	parser.RegisterOutgoingTemplate(1, template.PacketClassOfService, &template.TplSetClassOfService{})
	// parser.RegisterOutgoingTemplate(1, template.PacketSetCCRS, &template.TplSetCCRS{})
	// parser.RegisterOutgoingTemplate(1, template.PacketSetECC1, &template.TplSetECC1{})
	// parser.RegisterOutgoingTemplate(1, template.PacketSetECC2, &template.TplSetECC2{})

	parser.RegisterOutgoingTemplate(1, template.PacketWakeupSet, &template.TplSetWakeup{})
	parser.RegisterOutgoingTemplate(1, template.PacketWakeupClear, &template.TplClearWakeup{})

	// outgoing low level
	parser.RegisterOutgoingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(1, template.PacketEnq, &template.TplENQ{})

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
	automate.AddPlugin(linkcontrol.New(nil))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(datachange.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(callpacket.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(request.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(inhousesync.New(7))

}
