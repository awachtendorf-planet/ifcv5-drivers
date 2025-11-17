package request

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleRemoteCheckout(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketRemoteCheckOut

	driver := p.driver
	automate := p.automate
	name := automate.Name()

	dispatcher := automate.Dispatcher()
	station, _ := dispatcher.GetStationAddr(addr)

	type remoteCheckout struct {
		Station       uint64
		BalanceAmount float64 `fias:"BA"`
		ReservationID string  `fias:"G#"`
		Message       string  `fias:"CT"`
		RoomNumber    string  `fias:"RN"`
	}

	checkout := remoteCheckout{
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &checkout); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	answer, err := driver.MarshalPacket(checkout)

	if err != nil {
		log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, checkout, err)
	}

	answer.Set("AS", "UR")
	answer.Set("CT", "unsupported")

	err = p.sendRemoteCheckout(addr, action, packet.Tracking, answer)

	return nil
}

func (p *Plugin) sendRemoteCheckout(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketRemoteCheckOut, tracking, "XC", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
