package request

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleGuestCheckDetails(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	// template.PacketGuestCheckDetails

	driver := p.driver
	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	station, _ := dispatcher.GetStationAddr(addr)

	GuestCheckDetails := record.GuestCheckDetails{
		Station: station,
	}

	if err := driver.UnmarshalPacket(packet, &GuestCheckDetails); err != nil {
		log.Error().Msgf("%s addr '%s' unmarshal packet '%s' failed, err=%s", name, addr, packet.Name, err)
	}

	answer, err := driver.MarshalPacket(GuestCheckDetails)

	if err != nil {
		log.Error().Msgf("%s addr '%s' marshal type '%T' failed, err=%s", name, addr, GuestCheckDetails, err)
	}

	answer.Set("AS", "UR")
	answer.Set("CT", "unsupported")

	err = p.sendGuestCheckDetailsAnswer(addr, action, packet.Tracking, answer)

	return nil
}

func (p *Plugin) sendGuestCheckDetailsAnswer(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketGuestCheckDetails, tracking, "CK", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err

}
