package guestlink

import (
	"fmt"
	"time"

	"github.com/weareplanet/ifcv5-drivers/guestlink/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
	"github.com/weareplanet/ifcv5-main/utils/pad"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// ConstructPacket construct a logical packet
func (d *Dispatcher) ConstructPacket(addr string, packetName string, tracking string, context interface{}, transaction Transaction) (*ifc.LogicalPacket, error) {

	defer utils.TimeTrack(time.Now(), "construct packet '"+packetName+"'")

	// normalise device address
	addr = ifc.DeviceAddress(addr)

	station, _ := d.GetStationAddr(addr)

	packet := ifc.NewLogicalPacket(packetName, addr, tracking)

	// filter low level packets
	switch packetName {

	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return packet, nil

	}

	accountNumberLength := d.GetAccountNumberLength(station)
	accountNumberAlignment := d.GetAccountNumberAlignment(station)

	// set protocol transaction identifier and sequence number

	tan := []byte(transaction.Identifier)
	seq := []byte(fmt.Sprintf("%04d", transaction.Sequence))

	packet.Add(guestlink_tan, tan) // Mxxx
	packet.Add(guestlink_seq, seq) // xxxx

	switch packetName {

	case template.PacketVerify, template.PacketStart, template.PacketTest:

	case template.PacketError:

		errorCode, ok := context.(int)
		if !ok {
			errorCode = 0
		} else if errorCode > 99 || errorCode < 0 {
			errorCode = 0
		}

		err := []byte(fmt.Sprintf("%02d", errorCode))

		packet.Add(guestlink_error, err)

	case template.PacketHelo:

		packet.Add("ID", d.formatString(addr, "protel", 8, false))
		packet.Add("Version", d.formatInt(5, 4))
		packet.Add("STAT", d.formatBoolean(true))
		packet.Add("FUDP", d.formatBoolean(false))
		packet.Add("ROST", d.formatBoolean(false))
		packet.Add("MSGS", d.formatBoolean(true))
		packet.Add("DGMH", d.formatBoolean(false))
		packet.Add("MSGD", d.formatBoolean(true))
		packet.Add("PayTVChannel", d.formatBoolean(false))
		packet.Add("POST", d.formatBoolean(true))
		packet.Add("HSKP", d.formatBoolean(true))
		packet.Add("VariableMsgLength", d.formatBoolean(false))
		packet.Add("LocalFolioNumber", d.formatBoolean(true))
		packet.Add("PairedFolioNumber", d.formatBoolean(false))
		packet.Add("DecimalPosition", d.formatString(addr, d.getDecimalPoint(station), 1, false))
		packet.Add("CharacterSet", d.formatString(addr, "ISO88591", 8, false))
		packet.Add("FutureSpace", d.formatString(addr, "", 32, false))

	case template.PacketCheckIn, template.PacketCheckOut:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Guest'", context)
		}
		packet.Add("RoomNumber", d.formatString(addr, guest.Reservation.RoomNumber, 6, true))

	case template.PacketNameReply, template.PacketInfoReply:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Guest'", context)
		}

		packet.Add("RoomNumber", d.formatString(addr, guest.Reservation.RoomNumber, 6, true))
		packet.Add("AccountNumber", d.formatAccountNumber(guest.Reservation.ReservationID, accountNumberLength, accountNumberAlignment))
		packet.Add("GuestName", d.formatString(addr, d.getGuestName(guest), 20, true))
		packet.Add("MessageWaiting", d.formatBoolean(guest.Reservation.MessageLightStatus))

		if packetName == template.PacketNameReply {
			break
		}

		packet.Add("GroupName", d.formatString(addr, "", 5, false))
		packet.Add("BillView", d.formatBoolean(true))
		packet.Add("Checkout", d.formatBoolean(false))
		packet.Add("Welcome", d.formatString(addr, "", 1, false))
		packet.Add("Language", d.formatString(addr, d.getLanguage(station, guest.Language), 1, false))
		packet.Add("Rights", d.formatString(addr, d.getRights(station, guest.Rights.TV), 3, false))

	case template.PacketItemReply:

		item, ok := context.(record.BillPreviewItem)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected 'record.BillPreviewItem'", context)
		}

		packet.Add("Date", d.formatDate(item.Time)) // MMDD
		packet.Add("Description", d.formatString(addr, item.Description, 12, true))
		packet.Add("Indicator", d.formatString(addr, "", 2, false))
		packet.Add("Amount", d.formatAmount(item.Amount, 7))

	case template.PacketBalanceReply:

		balance, ok := context.(record.BillPreviewBalance)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected 'record.BillPreviewBalance'", context)
		}

		packet.Add("Amount", d.formatAmount(balance.Amount, 8))

	case template.PacketWakeupSet, template.PacketWakeupClear:

		generic, ok := context.(*record.Generic)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Generic'", context)
		}

		extension, _ := generic.Get(defines.WakeExtension)
		waketime, _ := generic.Get(defines.Timestamp)

		packet.Add("RoomNumber", d.formatString(addr, cast.ToString(extension), 6, true))
		packet.Add("Time", d.formatTime(cast.ToTime(waketime)))
		packet.Add("Date", d.formatDateWithYear(cast.ToTime(waketime)))

	case template.PacketGuestMessageStatus:

		guest, ok := context.(*record.Guest)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected '*record.Guest'", context)
		}

		packet.Add("RoomNumber", d.formatString(addr, guest.Reservation.RoomNumber, 6, true))
		packet.Add("StatusCode", d.formatBoolean(guest.Reservation.MessageLightStatus)) // Y/N, OnCommand 0/1

	case template.PacketGuestMessageHeader:

		msg, ok := context.(record.GuestMessage)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected 'record.GuestMessage'", context)
		}

		packet.Add("AccountNumber", d.formatAccountNumber(msg.ReservationID, accountNumberLength, accountNumberAlignment))
		packet.Add("MessageNumber", d.formatString(addr, msg.MessageID, 6, false))
		packet.Add("Date", d.formatDateWithYear(msg.Time))
		packet.Add("Time", d.formatTimeWithSeconds(msg.Time))
		packet.Add("ReceiverName", d.formatString(addr, d.getGuestName(msg), 24, true))
		packet.Add("ByPerson", d.formatBoolean(false))
		packet.Add("ByPhone", d.formatBoolean(false))
		packet.Add("PleaseCall", d.formatBoolean(false))
		packet.Add("Callback", d.formatBoolean(false))
		packet.Add("ReturnedCall", d.formatBoolean(false))
		packet.Add("Urgent", d.formatBoolean(false))

	case template.PacketGuestMessageCaller:

		msg, ok := context.(record.GuestMessage)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected 'record.GuestMessage'", context)
		}

		packet.Add("MessageNumber", d.formatString(addr, msg.MessageID, 6, false))
		packet.Add("CallerName", d.formatString(addr, msg.UserID, 24, true))
		packet.Add("CallerLocation", d.formatString(addr, "", 24, true))
		packet.Add("CallerPhoneNumber", d.formatString(addr, "", 24, true))

	case template.PacketGuestMessageText:

		msg, ok := context.(record.GuestMessage)
		if !ok {
			return packet, errors.Errorf("wrong context object '%T', expected 'record.GuestMessage'", context)
		}

		packet.Add("MessageNumber", d.formatString(addr, msg.MessageID, 6, false))
		packet.Add("MessageText", d.formatString(addr, msg.Text, 64, true))

	}

	return packet, nil
}

func (d *Dispatcher) formatAccountNumber(data string, width int, alignment int) []byte {
	if len(data) > width {
		data = data[0:width]
	}
	if len(data) == width {
		return []byte(data)
	}
	switch alignment {
	case JustifiedRight:
		data = pad.Left(data, width, " ")
	default:
		data = pad.Right(data, width, " ")
	}
	return []byte(data)
}

func (d *Dispatcher) formatString(addr string, data string, width int, encode bool) []byte {
	if encode {
		data = d.encode(addr, data)
	}
	if len(data) > width {
		data = data[0:width]
	}
	if len(data) == width {
		return []byte(data)
	}
	data = pad.Right(data, width, " ")
	return []byte(data)
}

func (d *Dispatcher) formatInt(value int, width int) []byte {
	data := cast.ToString(value)
	if len(data) > width {
		data = data[0:width]
	}
	if len(data) == width {
		return []byte(data)
	}
	data = pad.Left(data, width, " ")
	return []byte(data)
}

func (d *Dispatcher) formatBoolean(b bool) []byte {
	if b {
		return []byte("Y")
	}
	return []byte("N")
}

func (d *Dispatcher) formatDate(t time.Time) []byte {
	value := t.Format("0102") // MMDD
	return []byte(value)
}

func (d *Dispatcher) formatDateWithYear(t time.Time) []byte {
	value := t.Format("010206") // MMDDYY
	return []byte(value)
}

func (d *Dispatcher) formatTime(t time.Time) []byte {
	value := t.Format("1504") // HHMM
	return []byte(value)
}

func (d *Dispatcher) formatTimeWithSeconds(t time.Time) []byte {
	value := t.Format("150405") // HHMMSS
	return []byte(value)
}

func (d *Dispatcher) formatAmount(amount float64, width int) []byte {
	data := cast.ToString(amount)
	data = pad.Left(data, width, " ")
	return []byte(data)
}

func (d *Dispatcher) encode(addr string, data string) string {
	encoding := d.GetEncoding(addr)
	if len(encoding) == 0 {
		return data
	}
	if enc, err := d.Encode([]byte(data), encoding); err == nil && len(enc) > 0 {
		return string(enc)
	} else if err != nil && len(encoding) > 0 {
		log.Warn().Msgf("%s", err)
	}
	return data
}

func (d *Dispatcher) getGuestName(context interface{}) string {

	var name string

	switch context.(type) {

	case *record.Guest:

		guest := context.(*record.Guest)
		name = guest.DisplayName
		if len(name) == 0 && len(guest.LastName) > 0 && len(guest.FirstName) > 0 {
			name = fmt.Sprintf("%s %s", guest.LastName, guest.FirstName)
		}

	case record.GuestMessage:
		msg := context.(record.GuestMessage)
		name = msg.DisplayName
		if len(name) == 0 && len(msg.LastName) > 0 && len(msg.FirstName) > 0 {
			name = fmt.Sprintf("%s %s", msg.LastName, msg.FirstName)
		}
	}

	return name

}
func (d *Dispatcher) getLanguage(station uint64, key string) string {
	if match, err := d.GetMapping("languagecode", station, key, false); err == nil {
		return match
	}
	return ""
}

func (d *Dispatcher) getRights(station uint64, key int) string {
	// 0 = off
	// 1 = no adult but other movies
	// 2 = no pay movies
	// 3 = all

	if match, err := d.GetMapping("paytvright", station, cast.ToString(key), false); err == nil {
		return match
	}
	return ""
}

func (d *Dispatcher) getDecimalPoint(station uint64) string {
	decimals := d.GetConfig(station, defines.Decimals, "2")
	return decimals
}
