package request

import (
	"errors"
	"fmt"
	"strings"

	"github.com/weareplanet/ifcv5-drivers/micros/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	// "github.com/spf13/cast"
)

func (p *Plugin) handlePostingInquiry(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction) error {

	dispatcherObj := p.GetDispatcher()

	station, _ := dispatcherObj.GetStationAddr(addr)

	workstation := p.getField(packet, "WorkStation", false)
	requestMessage := p.getField(packet, "AccountID", true)

	employee := p.getField(packet, "EmployeeNumber", true)
	waiter := p.getField(packet, "TransEmployeeNumber", true)
	cashier := p.getField(packet, "ChkEmployeeNumber", true)

	var userID string

	isCOResponse := true

	response := ifc.NewLogicalPacket(template.PacketInquiryResponse, addr, packet.Tracking)

	if len(employee) > 0 {
		userID = employee
		isCOResponse = false
	} else if len(waiter) > 0 {

		userID = waiter
	} else if len(cashier) > 0 {

		userID = cashier
	}

	if isCOResponse {

		response.Add("MessageType", []byte(" 2"))
	} else {

		response.Add("MessageType", []byte(" 1"))
	}

	inquiry := record.PostingInquiry{
		Station:        station,
		WorkStation:    workstation,
		Inquiry:        requestMessage,
		UserID:         userID,
		MaximumMatches: 8,
		RequestType:    "0",
	}

	response.Add("WorkStation", []byte(workstation))

	reply, err := p.PmsRequest(addr, inquiry, packet)
	if err != nil {
		response.Add("InformationMessages", []byte(setLength("/"+err.Error(), 16)))
		chksum := p.driver.CalcChecksum(response.Addr, response)
		response.Add("Checksum", chksum)
		err := p.SendPacket(addr, response, action)
		return err
	}
	if reply == nil {
		response.Add("InformationMessages", []byte(setLength("PMS did not answer", 16)))
		chksum := p.driver.CalcChecksum(response.Addr, response)
		response.Add("Checksum", chksum)
		err := p.SendPacket(addr, response, action)
		return err
	}

	var guests []*record.Guest

	switch reply.(type) {
	case *record.Guest:
		guest := reply.(*record.Guest)
		guests = append(guests, guest)

	case []*record.Guest:
		guests = reply.([]*record.Guest)
	}

	if len(guests) == 0 {
		response.Add("InformationMessages", []byte(setLength("No match", 16)))
		chksum := p.driver.CalcChecksum(response.Addr, response)
		response.Add("Checksum", chksum)
		err := p.SendPacket(addr, response, action)
		return err
	}

	maxGuests := int(8)

	var names string

	for index, gast := range guests {

		if index >= maxGuests {
			break
		}

		if isCOResponse && index == 0 {

			message := "?"

			names = names + setLength(message, 16)

			name := fmt.Sprintf("%s %s", gast.Reservation.RoomNumber, gast.DisplayName)

			names = names + setLength(name, 16)

		} else if index == 0 {

			name := fmt.Sprintf("%s %s", gast.Reservation.RoomNumber, gast.DisplayName)

			names = names + setLength(name, 16)
		} else {

			name := fmt.Sprintf("%s %s", gast.Reservation.RoomNumber, gast.DisplayName)

			names = names + setLength(name, 16)
		}

	}
	response.Add("InformationMessages", []byte(names))

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

func setLength(message string, length int) string {
	if len(message) < length {
		message += strings.Repeat(" ", length-len(message))
	} else if len(message) > length {
		message = message[:length]
	}

	return message
}
