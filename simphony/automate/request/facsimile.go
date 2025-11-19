package request

import (
	"errors"

	"github.com/weareplanet/ifcv5-drivers/simphony/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
)

func (p *Plugin) handleFascimile(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {
	var err error

	retransmit := setLength((p.getField(packet, "MessageRetransmitFlag", true)), 1)
	sourceID := p.getField(packet, "SourceID", true)
	sequenceNumber := p.getField(packet, "SequenceNumber", true)

	data := p.getField(packet, "Data", true)

	response := ifc.NewLogicalPacket(template.PacketCheckFacsimileResponse, addr, packet.Tracking)

	successMessage := setLength(data, 70)

	response.Add("SourceID", []byte(setLength(sourceID, 25)))
	response.Add("SequenceNumber", []byte(sequenceNumber))
	response.Add("MessageRetransmitFlag", []byte(retransmit))

	response.Add("Text", []byte(successMessage))
	chksum := p.driver.CalcChecksum(response.Addr, response)
	response.Add("Checksum", chksum)

	if len(chksum) > 0 {
		if p.driver.CheckChecksum(packet) {
			err = p.SendPacket(addr, response, action)
			return err

		}
	} else {
		err = p.SendPacket(addr, response, action)
		return err
	}

	return errors.New("Checksum Mismatch")

}
