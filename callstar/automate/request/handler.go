package request

import (
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/callstar/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	errInternalPacket = errors.New("object packet is nil")
	errInternalJob    = errors.New("job failed")
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// if respectLinkState {
	// 	if !p.linkState(addr) {
	// 		automate.NextAction(name, addr, shutdown, t)
	// 		return
	// 	}
	// }

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}

	switch state {

	case busy:
		if job == nil || !job.InProcess() {
			err = errInternalJob
			break
		}

		if job.Context == nil {
			err = errInternalPacket
			break
		}

		packet, ok := job.Context.(*ifc.LogicalPacket)
		if !ok {
			err = errInternalPacket
			break
		}

		switch packet.Name {

		case template.PacketRequest:

			driver := p.driver

			room := driver.DecodeString(job.Station, packet.Data()["Room"])

			data := string(packet.Data()["Data"])

			phrases := strings.Split(data, "\\")

			for i := range phrases {

				cmd := strings.Split(phrases[i], "_")
				if len(cmd) == 0 || len(cmd[0]) == 0 {
					continue
				}

				char := cmd[0][0]
				cmd[0] = cmd[0][1:]

				switch char {

				case 'A', 'G':
					p.handleCallPacket(addr, room, cmd, packet)

				case 'I':
					p.handleWakeupSet(addr, room, cmd, packet)

				case 'L':
					p.handleMessageLamp(addr, room, cmd, packet)

				case 'M':
					p.handleMinibar(addr, room, cmd, packet)

				case 'P', 'T', 'X':
					p.handlePostCharge(addr, room, cmd, packet)

				case 'Z':
					p.handleVoiceCharge(addr, room, cmd, packet)

				case 'S':
					p.handleRoomStatus(addr, room, cmd, packet)

				case 'V':
					p.handleVoiceMessage(addr, room, cmd, packet)

				// case 'W': // message waiting ?

				default:
					log.Warn().Msgf("%s addr '%s' unknown command character '%c'", name, addr, char)

				}

			}

		default:
			log.Warn().Msgf("%s addr '%s' unknown packet '%s'", name, addr, packet.Name)
			automate.NextAction(name, addr, shutdown, t)
			return
		}

		automate.NextAction(name, addr, success, t)
		return

	case success, shutdown:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		if err == errInternalPacket || err == errInternalJob {
			automate.NextAction(name, addr, shutdown, t)
			return
		}
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) getString(data string) string {
	str := strings.Trim(data, " ")
	return str
}

func (p *Plugin) getNumeric(data string) int {
	str := strings.Trim(data, " ")
	str = strings.TrimLeft(str, "0")
	ret := cast.ToInt(str)
	return ret
}

func (p *Plugin) getFloat64(data string) float64 {
	str := strings.Trim(data, " ")
	str = strings.TrimLeft(str, "0")
	ret := cast.ToFloat64(str)
	return ret
}

func (p *Plugin) getTime(data string) (time.Time, error) {
	var timeFormat string
	if len(data) == 4 {
		timeFormat = "1504"
	}
	timestamp, err := time.Parse(timeFormat, data)
	return timestamp, err
}
