package keyservice

import (
	"strings"
	"time"

	messerschmitt "github.com/weareplanet/ifcv5-drivers/messerschmitt/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	ks "github.com/weareplanet/ifcv5-main/ifc/automate/keyservice"

	"github.com/weareplanet/ifcv5-drivers/messerschmitt/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/messerschmitt/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	keyLifetime = 120 * time.Second
)

// Plugin ...
type Plugin struct {
	*ks.Plugin
	driver *messerschmitt.Dispatcher
}

// New return a new plugin
func New(parent *messerschmitt.Dispatcher) *Plugin {

	p := &Plugin{
		ks.New(),
		parent,
	}

	p.Setup(ks.Config{

		TemplateAck: template.PacketAck,
		TemplateNak: template.PacketNak,
		TemplateEnq: template.PacketEnq,
		TemplateEOT: template.PacketEOT,

		InitHandler:        p.init,
		SendPacket:         p.send,
		PreCheck:           p.preCheck,
		GetAnswerTimeout:   p.getAnswerTimeout,
		WaitForReplyPacket: p.waitForReplyPacket,
		SendENQ:            p.sendENQ,
		SendEOT:            p.sendEOT,
	})

	return p
}

func (p *Plugin) init() {

	p.setWorkflow()

	p.RegisterRule(template.PacketCommandAcknowledge, automatestate.WaitForAnswer, p.handleReply, dispatcher.StateAction{NextState: automatestate.NextRecord})

}

func (p *Plugin) setWorkflow() {

	p.RegisterWorkflow(order.KeyRequest,
		template.PacketKeyRequest,
	)

	p.RegisterWorkflow(order.KeyDelete,
		template.PacketKeyDelete,
	)

}

func (p *Plugin) preCheck(addr string, job *order.Job) (bool, error) {

	if job != nil && job.Timestamp > 0 {
		timestamp := time.Unix(0, job.Timestamp)
		since := time.Since(timestamp)
		if since > keyLifetime {
			return false, errors.Errorf("key service deadline exceeded (%s)", since)
		}
	}

	guest, ok := job.Context.(*record.Guest)
	if !ok {
		return false, errors.Errorf("context '%T' not supported", job.Context)
	}

	encoder := p.driver.GetEncoder(guest)
	isSerial := p.isSerial(addr)

	if isSerial {
		if encoder > 15 {
			return false, errors.Errorf("encoder number must be less than or equal to 15")
		}
	} else {
		if encoder > 32 {
			return false, errors.Errorf("encoder number must be less than or equal to 32")
		}
	}

	room := guest.Reservation.RoomNumber

	if isSerial {
		if len(room) > 4 {
			return false, errors.Errorf("room number '%s' too long (max: %d)", guest.Reservation.RoomNumber, 4)
		}
		if len(room) > 0 && room[0] == '0' {
			return false, errors.Errorf("room number '%s' has leading zero", guest.Reservation.RoomNumber)
		}
		if number := cast.ToInt(room); number == 0 {
			return false, errors.Errorf("room number '%s' is not numeric", guest.Reservation.RoomNumber)
		}
	} else {
		room = p.driver.EncodeString(job.Station, room)
		if len(room) > 10 {
			return false, errors.Errorf("room name '%s' too long (max: %d)", guest.Reservation.RoomNumber, 10)
		}
	}

	return true, nil
}

func (p *Plugin) waitForReplyPacket(addr string, job *order.Job) bool {

	if job.Action == order.KeyRequest {
		return true
	}

	return false
}

func (p *Plugin) sendENQ(addr string, _ *order.Job) bool {

	if !p.isSerial(addr) {
		return false
	}

	return true
}

func (p *Plugin) sendEOT(addr string, _ *order.Job) bool {

	if !p.isSerial(addr) {
		return false
	}

	// if we reply with a packet thats ends of ETB (0x17) there is no need for EOT

	return true

}

func (p *Plugin) isSerial(addr string) bool {

	if dispatcher := p.GetDispatcher(); dispatcher != nil {
		if dispatcher.IsSerialDevice(addr) {
			return true
		}
	}

	return false

}

func (p *Plugin) getAnswerTimeout(addr string, job *order.Job) time.Duration {

	answerTimeout := messerschmitt.KeyAnswerTimeout

	if timeout, ok := p.GetSetting(addr, defines.KeyAnswerTimeout, uint(0)).(uint64); ok && timeout > 0 {
		if duration := cast.ToDuration(time.Duration(timeout) * time.Second); duration > 0 {
			answerTimeout = duration
		}
	}

	keyCount := p.driver.GetKeyCount(job)

	return answerTimeout * time.Duration(keyCount)
}

func (p *Plugin) getKeyAsString(key string, packet *ifc.LogicalPacket) string {

	if packet == nil || len(key) == 0 {
		return ""
	}

	value := string(packet.Data()[key])
	value = strings.TrimLeft(value, " ")

	return value
}

func (p *Plugin) getKeyAsInt(key string, packet *ifc.LogicalPacket) int {

	if packet == nil || len(key) == 0 {
		return 0
	}

	value := string(packet.Data()[key])

	value = strings.TrimLeft(value, " ")
	value = strings.TrimLeft(value, "0")

	n := cast.ToInt(value)

	return n
}

func (p *Plugin) handleReply(addr string, packet *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

	if guest, ok := job.Context.(*record.Guest); ok {

		isSerial := p.isSerial(addr)

		// check answer encoder

		encoderAnswer := p.getKeyAsString("Encoder", packet)
		encoder := p.driver.FormatEncoder(guest, isSerial)

		if encoder != encoderAnswer {

			name := p.GetName()
			log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' for encoder: %s, expect: %s", name, addr, packet.Name, encoderAnswer, encoder)
			p.CalculateNextTimeout(addr, action, job)

			return nil
		}

		// check answer room

		/*

			station, _ := p.driver.GetStationAddr(addr)

			room := p.driver.GetRoom(station, guest, isSerial)

			roomAnswer := p.getKeyAsString("Room", packet)

			if room != roomAnswer {

				name := p.GetName()
				log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' for room: %s, expect: %s", name, addr, packet.Name, roomAnswer, room)
				p.CalculateNextTimeout(addr, action, job)

				return nil
			}
		*/

	}

	status := p.getKeyAsInt("Status", packet)

	switch status {

	case 0:

		track := p.getKeyAsString("TrackData", packet)
		if track == "5100000000000000000000000000000000000" {
			track = ""
		}
		p.HandleSuccess(addr, job, track)

	default:

		answer := p.driver.GetAnswerText(status)
		p.HandleError(addr, job, answer)

	}

	p.ChangeState(addr, action.NextState)

	return nil
}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, packetName string, tracking string, context interface{}, _ *order.Job) error {

	packet, err := p.driver.ConstructPacket(addr, packetName, tracking, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}
