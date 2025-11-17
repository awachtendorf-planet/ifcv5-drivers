package main

import (
	"os"
	"time"

	"github.com/weareplanet/ifcv5-main/app"
	"github.com/weareplanet/ifcv5-main/config"
	"github.com/weareplanet/ifcv5-main/ifc/noise/builder"
	noise "github.com/weareplanet/ifcv5-main/ifc/noise/const"
	"github.com/weareplanet/ifcv5-main/ifc/noise/plugin/access"
	"github.com/weareplanet/ifcv5-main/ifc/noise/plugin/api"
	"github.com/weareplanet/ifcv5-main/ifc/noise/transport"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/noise/network"

	"github.com/robfig/cron/v3"
)

var (
	_ app.AppInterface = (*App)(nil)
)

// App represents a pcon app
type App struct {
	start   time.Time
	network *network.Network
	cron    *cron.Cron
}

// New returns a new pcon app
func New(configFile string) (*App, error) {

	cfg, err := config.Read(noise.Connector, configFile)
	if err != nil {
		return nil, err
	}

	if len(cfg.MetricAddr) > 0 {
		initMetric(cfg.MetricAddr, "pcon")
	}

	noiseConfig := builder.Config{
		Access: &access.Config{
			Alias: "pcon-default-local",
			Token: cfg.Token,
			Mode:  noise.Connector,
			Store: cfg.Store,
		},
		Api: &api.Config{
			Address: cfg.APIAddr,
			Mode:    noise.Connector,
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
	network, err := builder.NewConnectorNetwork(&noiseConfig)
	if err != nil {
		return nil, err
	}

	return &App{
		start:   time.Now(),
		network: network,
		cron:    cron.New(),
	}, nil
}

// Startup ...
func (a *App) Startup() error {

	// here we go
	log.Info().Msgf("startup pid: %d version '%s'", os.Getpid(), config.VersionString())

	if plugin, registered := a.network.Plugin(api.PluginID); registered {
		api := plugin.(*api.Plugin)
		api.RegisterControlHandler(a.controlHandler)
	}

	// start noise network
	go a.network.Listen()
	a.network.BlockUntilListening()
	a.network.SetReady()

	return nil
}

// InitLogRotate ...
func (a *App) InitLogRotate(trigger func()) {

	if trigger != nil && a.cron != nil {

		if _, err := a.cron.AddFunc("0 4 * * * ", trigger); err != nil {
			log.Error().Msgf("add cron job for log rotate failed, err=%s", err)
		} else {
			log.Debug().Msg("add cron job for log rotate")
		}

		a.cron.Start()
	}

}

// Close ...
func (a *App) Close() {

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

	// shutdown cron jobs
	if a.cron != nil {
		a.cron.Stop()
	}

	log.Info().Msgf("done, uptime %s", time.Since(a.start))
}
