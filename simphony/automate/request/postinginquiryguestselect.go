package request

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/weareplanet/ifcv5-drivers/simphony/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handlePostingInquiry(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	dispatcherObj := p.GetDispatcher()

	station, _ := dispatcherObj.GetStationAddr(addr)

	requestInquiry := p.getField(packet, "GuestID", true)
	retransmit := setLength((p.getField(packet, "MessageRetransmitFlag", true)), 1)
	sequenceNumber := p.getField(packet, "SequenceNumber", true)
	paymentType := p.getField(packet, "PaymentType", true)
	sourceID := p.getField(packet, "SourceID", true)

	workStation := p.getWorkStation(sourceID, station)

	inquiry := record.PostingInquiry{
		Station:        station,
		Inquiry:        requestInquiry,
		MaximumMatches: 8,
		RequestType:    "0",
		WorkStation:    workStation,
		SequenceNumber: sequenceNumber,
		PaymentMethod:  paymentType,
	}

	//sequenceNumber = p.incrementSequenceNumber(sequenceNumber)

	reply, err := p.PmsRequest(addr, inquiry, packet)
	if err != nil {

		response := p.buildErrorResponse(err.Error(), retransmit, addr, packet.Tracking, sourceID, sequenceNumber)
		err := p.SendPacket(addr, response, action)
		return err
	}
	if reply == nil {
		response := p.buildErrorResponse("no reply received", retransmit, addr, packet.Tracking, sourceID, sequenceNumber)
		err := p.SendPacket(addr, response, action)
		return err
	}

	var guests []*record.Guest

	switch replyRecord := reply.(type) {
	case *record.Guest:
		guests = append(guests, replyRecord)

	case []*record.Guest:
		guests = replyRecord

	case *record.PostingResponse:
		response := p.buildErrorResponse(replyRecord.ResponseText, retransmit, addr, packet.Tracking, sourceID, sequenceNumber)
		err := p.SendPacket(addr, response, action)
		return err
	}

	var names string

	for idx, gast := range guests {

		name := setLength(fmt.Sprintf("%s / %s, %s", gast.Reservation.RoomNumber, gast.LastName, gast.FirstName), 40)

		if idx < len(guests)-1 {
			name = name + string(rune(0x1c))
		}

		names = names + name
	}

	if len(names) == 0 {
		response := p.buildErrorResponse("No guests found", retransmit, addr, packet.Tracking, sourceID, sequenceNumber)
		err := p.SendPacket(addr, response, action)
		return err
	}

	response := ifc.NewLogicalPacket(template.PacketInquiryResponse, addr, packet.Tracking)

	response.Add("SourceID", []byte(setLength(sourceID, 25)))
	response.Add("SequenceNumber", []byte(sequenceNumber))
	response.Add("MessageRetransmitFlag", []byte(retransmit))

	guestID := requestInquiry

	response.Add("GuestID", []byte(guestID))
	response.Add("ListSize", []byte(strconv.Itoa(len(guests))))
	response.Add("GuestList", []byte(names))

	chksum := p.driver.CalcChecksum(response.Addr, response)
	response.Add("Checksum", chksum)

	if p.driver.CheckChecksum(packet) {
		log.Warn().Msgf("checksum mismatch")
	}

	err = p.SendPacket(addr, response, action)
	return err

}

func (p *Plugin) getWorkStation(sourceId string, station uint64) string {
	length := p.driver.GetWorkStationLength(station)
	return sourceId[:length]
}

func setLength(message string, length int) string {
	if len(message) < length {
		message += strings.Repeat(" ", length-len(message))
	} else if len(message) > length {
		message = message[:length]
	}

	return message
}

func (p *Plugin) buildErrorResponse(message, retransmit, addr, tracking, sourceID, sequenceNumber string) *ifc.LogicalPacket {
	if len(retransmit) == 0 {
		retransmit = " "
	}

	response := ifc.NewLogicalPacket(template.PacketChargePostingAck, addr, tracking)

	response.Add("SourceID", []byte(setLength(sourceID, 25)))
	response.Add("MessageRetransmitFlag", []byte(retransmit))
	response.Add("SequenceNumber", []byte(sequenceNumber))
	response.Add("MessageStatus", []byte("N"))
	response.Add("Status", []byte("E"))
	response.Add("Message", []byte(setLength(message, 30)))

	chksum := p.driver.CalcChecksum(response.Addr, response)
	response.Add("Checksum", chksum)

	return response
}
