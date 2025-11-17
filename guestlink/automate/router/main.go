package router

import (
	"time"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	"github.com/weareplanet/ifcv5-main/log"
)

const (
	wildcardAddress  = "*"
	respectLinkState = true
)

func (p *Plugin) main() error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := p.driver
	name := automate.Name()

	incoming, err := dispatcher.RegisterPacketRoute(name, wildcardAddress,
		template.PacketVerify,
		template.PacketError,
	)
	if err != nil {
		return err
	}

	defer func() {
		dispatcher.DeregisterPacketRoute(name, wildcardAddress)
	}()

	for {

		select {

		// shutdown
		case <-p.kill:
			return nil

		// incoming logical packet
		case packet := <-incoming.C:

			if packet == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			packetName := packet.Name
			packetAddr := packet.Addr

			switch packet.Name {

			case template.PacketVerify, template.PacketError:

				transaction := driver.GetTransaction(packet)

				log.Info().Msgf("%s addr '%s' received packet '%s', transaction '%s'", name, packetAddr, packetName, transaction.Identifier)

				if owner, exist := driver.Owner(packetAddr, transaction); exist {

					if forward, exist := driver.GetRoute(owner, packetAddr); exist && forward != nil && forward.C != nil {

						log.Debug().Msgf("%s addr '%s' owner '%T' for transaction: '%s' found, dispatch packet '%s'", name, packetAddr, owner, transaction.Identifier, packetName)

						// defer func() {
						// 	err := recover()
						// 	if err != nil {
						// 		log.Error().Msgf("%s addr '%s' owner '%T' transaction: '%s' dispatch packet '%s' failed, err=%s", name, packetAddr, owner, transaction.Identifier, packetName, err)
						// 	}
						// }()

						select {

						case forward.C <- packet:

						case <-time.After(10 * time.Second):
							log.Error().Msgf("%s addr '%s' owner '%T' transaction: '%s' dispatch packet '%s' failed, err=%s", name, packetAddr, owner, transaction.Identifier, packetName, "channel timeout")
						}

					} else {
						log.Warn().Msgf("%s addr '%s' owner '%T' for transaction: '%s' found, but no route registered, packet '%s'", name, packetAddr, owner, transaction.Identifier, packetName)
					}

				} else {
					log.Warn().Msgf("%s addr '%s' no owner for transaction: '%s' found, packet '%s'", name, packetAddr, transaction.Identifier, packetName)
				}

			}

		}

	}

	return nil
}
