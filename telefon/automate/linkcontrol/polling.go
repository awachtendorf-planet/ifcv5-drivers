package linkcontrol

import (
	"time"

	"github.com/weareplanet/ifcv5-drivers/telefon/template"
	"github.com/weareplanet/ifcv5-main/log"
)

const (
	minPollInterval = 3 * time.Second
)

func (p *Plugin) polling(addr string, kill chan bool) {

	automate := p.automate
	name := automate.Name()

	if p.state.Exist(addr) {
		log.Warn().Msgf("%s addr '%s' polling handler still exist", name, addr)
		return
	}

	p.waitGroup.Add(1)
	p.state.Register(addr)

	defer func() {
		log.Info().Msgf("%s stop polling handler addr '%s'", name, addr)
		p.state.Remove(addr)
		p.waitGroup.Done()
	}()

	log.Info().Msgf("%s start polling handler addr '%s'", name, addr)

	driver := p.driver
	dispatcher := automate.Dispatcher()

	if driver == nil || dispatcher == nil {
		return
	}

	timeout := 1 * time.Second // first run

	for {

		select {

		case <-kill:
			return

		case <-p.kill:
			return

		case <-time.After(timeout):

			if state := dispatcher.GetLinkState(addr); state == linkUp {
				if err := driver.SendPacket(addr, template.Polling); err == nil {
					dispatcher.SetAlive(addr)
				}
			}

		}

		if interval := driver.GetPollInterval(addr); interval > 0 {
			timeout = time.Duration(interval) * time.Second
		}

		if timeout < minPollInterval {
			timeout = minPollInterval
		}

	}
}
