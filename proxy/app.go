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
)

var (
	_ app.AppInterface = (*App)(nil)
)

// App represents a proxy app
type App struct {
	start   time.Time
	network *network.Network
}

// New returns a new proxy app
func New(configFile string) (*App, error) {

	cfg, err := config.Read(noise.Proxy, configFile)
	if err != nil {
		return nil, err
	}

	if len(cfg.MetricAddr) > 0 {
		initMetric(cfg.MetricAddr, "proxy")
	}

	noiseConfig := builder.Config{
		Access: &access.Config{
			Alias:     "proxy",
			Token:     cfg.Token,
			Blacklist: cfg.Blacklist,
			Mode:      noise.Proxy,
			Store:     cfg.Store,
		},
		Api: &api.Config{
			Address: cfg.APIAddr,
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
	network, err := builder.NewProxyNetwork(&noiseConfig)
	if err != nil {
		return nil, err
	}

	return &App{
		start:   time.Now(),
		network: network,
	}, nil
}

// Startup ...
func (a *App) Startup() error {
	// here we go
	log.Info().Msgf("startup pid: %d version '%s'", os.Getpid(), config.VersionString())
	go a.network.Listen()
	a.network.BlockUntilListening()
	a.network.SetReady()
	return nil
}

// InitLogRotate ...
func (a *App) InitLogRotate(_ func()) {

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
	log.Info().Msgf("done, uptime %s", time.Since(a.start))
}
