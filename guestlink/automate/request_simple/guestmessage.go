package requestsimple

import (
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleGuestMessageRead(addr string, packet *ifc.LogicalPacket, retry bool) (int, error) {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	accountNumber, err := p.getAccountNumber(addr, packet)
	if err != nil {
		return unknownAccount, err
	}

	if retry {
		return 0, nil
	}

	station, _ := dispatcher.GetStationAddr(addr)

	messageID := p.getField(packet, "MessageNumber", true)
	messageID = strings.TrimLeft(messageID, "0")

	guestMessage := record.GuestMessageDelete{
		Station:       station,
		ReservationID: accountNumber,
		MessageID:     messageID,
	}

	dispatcher.CreatePmsJob(addr, packet, guestMessage)

	return 0, nil

}
