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

	"github.com/weareplanet/ifcv5-drivers/tiger/template"

	"github.com/weareplanet/ifcv5-drivers/tiger/automate/datachange"
	"github.com/weareplanet/ifcv5-drivers/tiger/automate/request"
	"github.com/weareplanet/ifcv5-main/ifc/automate/inhousesync"
	linkcontrol "github.com/weareplanet/ifcv5-main/ifc/automate/simplelinkcontrol"

	drv "github.com/weareplanet/ifcv5-drivers/tiger/dispatcher"
	analyser "github.com/weareplanet/ifcv5-main/ifc/generic/parser/bytesparserextended"
)

var (
	_ app.AppInterface = (*App)(nil)
)

var (
	driverAddr = "ifc:tiger:1"
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
	parser.RegisterIncomingTemplate(3, template.PacketRoomBasedCall, &template.TplRoomBasedCallDetailPacket{})

	parser.RegisterIncomingTemplate(3, template.PacketMiniBarBilling, &template.TplMiniBarBillingPacket{})
	parser.RegisterIncomingTemplate(3, template.PacketOtherChargePosting, &template.TplOtherChargingPacket{})

	parser.RegisterIncomingTemplate(3, template.PacketRoomBasedCall, &template.TplRoomBasedCallDetailPacketk{})

	parser.RegisterIncomingTemplate(3, template.PacketMiniBarBilling, &template.TplMiniBarBillingPacketk{})
	parser.RegisterIncomingTemplate(3, template.PacketOtherChargePosting, &template.TplOtherChargingPacketk{})

	parser.RegisterIncomingTemplate(3, template.PacketRoomstatus, &template.TplRoomStatusPacket{})

	parser.RegisterIncomingTemplate(3, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacket{})
	parser.RegisterIncomingTemplate(3, template.PacketMessageWaitingReservation, &template.TplMessageWaitingReservationPacket{})
	// Opera

	parser.RegisterIncomingTemplate(0, template.PacketRoomBasedCall, &template.TplRoomBasedCallDetailPacketOpera{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomBasedCall, &template.TplRoomBasedCallDetailPacketOpera{})
	parser.RegisterIncomingTemplate(2, template.PacketRoomBasedCall, &template.TplRoomBasedCallDetailPacketOpera{})

	parser.RegisterIncomingTemplate(0, template.PacketMiniBarBilling, &template.TplMiniBarBillingPacket{})
	parser.RegisterIncomingTemplate(1, template.PacketMiniBarBilling, &template.TplMiniBarBillingPacket{})
	parser.RegisterIncomingTemplate(2, template.PacketMiniBarBilling, &template.TplMiniBarBillingPacket{})

	parser.RegisterIncomingTemplate(0, template.PacketRoomstatus, &template.TplRoomStatusPacketOpera{})
	parser.RegisterIncomingTemplate(1, template.PacketRoomstatus, &template.TplRoomStatusPacketOpera{})
	parser.RegisterIncomingTemplate(2, template.PacketRoomstatus, &template.TplRoomStatusPacketOpera{})

	parser.RegisterIncomingTemplate(0, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacketOpera{})
	parser.RegisterIncomingTemplate(1, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacketOpera{})
	parser.RegisterIncomingTemplate(2, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacketOpera{})

	// incoming low level
	parser.RegisterIncomingTemplate(0, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(0, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(0, template.PacketEnq, &template.TplENQ{})

	parser.RegisterIncomingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(1, template.PacketEnq, &template.TplENQ{})

	parser.RegisterIncomingTemplate(2, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(2, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(2, template.PacketEnq, &template.TplENQ{})

	parser.RegisterIncomingTemplate(3, template.PacketAck, &template.TplACK{})
	parser.RegisterIncomingTemplate(3, template.PacketNak, &template.TplNAK{})
	parser.RegisterIncomingTemplate(3, template.PacketEnq, &template.TplENQ{})

	// incoming garbage
	parser.RegisterIncomingTemplateWithOption(0, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplate(0, template.PacketUnknown, &template.TplUnknownPacket{})

	parser.RegisterIncomingTemplateWithOption(1, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplate(1, template.PacketUnknown, &template.TplUnknownPacket{})

	parser.RegisterIncomingTemplateWithOption(2, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplate(2, template.PacketUnknown, &template.TplUnknownPacket{})

	parser.RegisterIncomingTemplateWithOption(3, template.PacketGarbage, &template.TplGarbage_Framing_1{}, 1)
	parser.RegisterIncomingTemplate(3, template.PacketUnknown, &template.TplUnknownPacket{})

	// incomig low level garbage
	parser.RegisterIncomingTemplateWithOption(0, template.PacketGarbageLowLevel, &template.TplGarbage_ACK{}, 1)
	parser.RegisterIncomingTemplateWithOption(0, template.PacketGarbageLowLevel, &template.TplGarbage_NAK{}, 1)
	parser.RegisterIncomingTemplateWithOption(0, template.PacketGarbageLowLevel, &template.TplGarbage_ENQ{}, 1)

	// outgoing
	// HERE WE GO
	parser.RegisterOutgoingTemplate(3, template.PacketRoomBasedCheckIn, &template.TplCheckInOutRoomBasedPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketRoomBasedCheckOut, &template.TplCheckInOutRoomBasedPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketAdditionalGuest, &template.TplAdditionalGuestPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketWakeUpSet, &template.TplWakeUpPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketWakeUpClear, &template.TplWakeUpPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketRoomstatus, &template.TplRoomStatusPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketRoomtransfer, &template.TplRoomTransferPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketInformationUpdate, &template.TplInformationUpdatePacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketVideoRights, &template.TplVideoRightsPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketDND, &template.TplDoNotDisturbPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketMessageWaitingReservation, &template.TplMessageWaitingReservationPacket{})
	parser.RegisterOutgoingTemplate(3, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacket{})
	// Opera
	parser.RegisterOutgoingTemplate(0, template.PacketRoomBasedCheckIn, &template.TplCheckInRoomBasedPacketOpera0{})
	parser.RegisterOutgoingTemplate(1, template.PacketRoomBasedCheckIn, &template.TplCheckInRoomBasedPacketOpera12{})
	parser.RegisterOutgoingTemplate(2, template.PacketRoomBasedCheckIn, &template.TplCheckInRoomBasedPacketOpera12{})

	parser.RegisterOutgoingTemplate(0, template.PacketRoomBasedCheckOut, &template.TplCheckOutRoomBasedPacketOpera0{})
	parser.RegisterOutgoingTemplate(1, template.PacketRoomBasedCheckOut, &template.TplCheckOutRoomBasedPacketOpera12{})
	parser.RegisterOutgoingTemplate(2, template.PacketRoomBasedCheckOut, &template.TplCheckOutRoomBasedPacketOpera12{})

	parser.RegisterOutgoingTemplate(0, template.PacketWakeUpSet, &template.TplWakeUpPacketOpera{})
	parser.RegisterOutgoingTemplate(1, template.PacketWakeUpSet, &template.TplWakeUpPacketOpera{})
	parser.RegisterOutgoingTemplate(2, template.PacketWakeUpSet, &template.TplWakeUpPacketOpera{})
	parser.RegisterOutgoingTemplate(0, template.PacketWakeUpClear, &template.TplWakeUpPacketOpera{})
	parser.RegisterOutgoingTemplate(1, template.PacketWakeUpClear, &template.TplWakeUpPacketOpera{})
	parser.RegisterOutgoingTemplate(2, template.PacketWakeUpClear, &template.TplWakeUpPacketOpera{})

	parser.RegisterOutgoingTemplate(0, template.PacketRoomtransfer, &template.TplRoomTransferPacketOpera01{})
	parser.RegisterOutgoingTemplate(1, template.PacketRoomtransfer, &template.TplRoomTransferPacketOpera01{})
	parser.RegisterOutgoingTemplate(2, template.PacketRoomtransfer, &template.TplRoomTransferPacketOpera2{})

	parser.RegisterOutgoingTemplate(0, template.PacketInformationUpdate, &template.TplInformationUpdatePacketOpera0{})
	parser.RegisterOutgoingTemplate(1, template.PacketInformationUpdate, &template.TplInformationUpdatePacketOpera12{})
	parser.RegisterOutgoingTemplate(2, template.PacketInformationUpdate, &template.TplInformationUpdatePacketOpera12{})

	parser.RegisterOutgoingTemplate(0, template.PacketDND, &template.TplDoNotDisturbPacket{})
	parser.RegisterOutgoingTemplate(1, template.PacketDND, &template.TplDoNotDisturbPacket{})
	parser.RegisterOutgoingTemplate(2, template.PacketDND, &template.TplDoNotDisturbPacket{})

	parser.RegisterOutgoingTemplate(0, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacketOpera{})
	parser.RegisterOutgoingTemplate(1, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacketOpera{})
	parser.RegisterOutgoingTemplate(2, template.PacketMessageWaitingGuest, &template.TplMessageWaitingGuestPacketOpera{})

	// outgoing low level
	parser.RegisterOutgoingTemplate(0, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(0, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(0, template.PacketEnq, &template.TplENQ{})

	parser.RegisterOutgoingTemplate(1, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(1, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(1, template.PacketEnq, &template.TplENQ{})

	parser.RegisterOutgoingTemplate(2, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(2, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(2, template.PacketEnq, &template.TplENQ{})

	parser.RegisterOutgoingTemplate(3, template.PacketAck, &template.TplACK{})
	parser.RegisterOutgoingTemplate(3, template.PacketNak, &template.TplNAK{})
	parser.RegisterOutgoingTemplate(3, template.PacketEnq, &template.TplENQ{})
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
	automate.AddPlugin(linkcontrol.New(nil))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(datachange.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(request.New(parent))

	automate = dispatcher.NewAutomate(parent.Dispatcher)
	automate.AddPlugin(inhousesync.New(7))

}
