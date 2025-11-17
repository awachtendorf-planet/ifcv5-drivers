package keyservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/weareplanet/ifcv5-drivers/vc3000/dispatcher"
	vc3000 "github.com/weareplanet/ifcv5-drivers/vc3000/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/vc3000/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	keyLifetime = 120 * time.Second
	slotTimeout = 15 * time.Second
	maxSlots    = 25
)

func (p *Plugin) main() {

	dispatcher := p.automate.Dispatcher()

	for {

		select {

		// shutdown
		case <-p.kill:
			return

		// action
		case order := <-p.job:
			encoder := 0
			if guest, ok := order.Job.Context.(*record.Guest); ok {

				if coder, exist := guest.GetGeneric(defines.EncoderNumber); exist {
					encoder = cast.ToInt(coder)
				}

				isSerialLayer := dispatcher.IsSerialDevice(order.Addr)
				if isSerialLayer {
					p.handleJob(order.Addr, encoder, order.Job)
					continue
				}

				slot, err := p.driver.GetSlot(order.Addr)
				if err == nil {
					order.Addr = slot
					p.handleJob(order.Addr, encoder, order.Job)
				} else {
					// Aktuell ist kein Slot frei.
					// Neuen Slot anfordern bzw. warten mit Timeout bis ein Slot frei wird.
					go p.waitForSlot(order, encoder)
				}
			} else {
				errStr := fmt.Sprintf("job context is not a guest record (%T)", order.Job.Context)
				p.finaliseJobWithoutSlot(order.Addr, order.Job, errStr)
			}

		}

	}

}

func (p *Plugin) waitForSlot(order order.DriverJob, encoder int) {

	dispatcher := p.automate.Dispatcher()
	name := p.automate.Name()
	driver := p.driver

	log.Info().Msgf("%s addr '%s' wait for available slot", name, order.Addr)

	handleFailed := func(err error) {
		if err != nil {
			log.Error().Msgf("%s", err)
		}
		errStr := "all connections are currently in use"
		p.handleFailed(order.Addr, order.Job, errStr)
		if err != nil {
			errStr = err.Error()
		}
		p.finaliseJobWithoutSlot(order.Addr, order.Job, errStr)
	}

	station, err := dispatcher.GetStationAddr(order.Addr)
	if err != nil {
		handleFailed(err)
		return
	}

	slotName := driver.GetSlotEventName(station)
	subscriber, err := dispatcher.RegisterEvents(slotName)
	if err != nil {
		handleFailed(err)
		return
	}

	defer func() {
		dispatcher.DeregisterEvents(subscriber)
	}()

	slotsCount := driver.GetSlotsCount(order.Addr)
	if slotsCount < maxSlots {
		err = dispatcher.RequestConnection(order.Addr)
		if err != nil {
			handleFailed(err)
			return
		}
	}

	start := time.Now()

	for {

		if time.Since(start) >= slotTimeout {
			handleFailed(nil)
			return
		}

		select {

		// shutdown
		case <-p.kill:
			handleFailed(nil)
			return

		// timeout
		case <-time.After(5 * time.Second):

		// event broker new slot available
		case sub := <-subscriber.GetMessages():
			if sub == nil {
				continue
			}
			event := sub.GetPayload()

			switch event.(type) {

			case *vc3000.SlotEvent:
				event := event.(*vc3000.SlotEvent)
				if event.Station == station {
					slot, err := p.driver.GetSlot(order.Addr)
					if err == nil {
						order.Addr = slot
						p.handleJob(order.Addr, encoder, order.Job)
						return
					}
				}

			}

		}
	}

}

func (p *Plugin) handleAnswer(addr string, job *order.Job, keyAnswer record.KeyAnswer) {
	if job == nil || job.IsDone() {
		return
	}
	job.Done()

	automate := p.automate
	dispatcher := automate.Dispatcher()

	keyAnswer.Station = job.Station
	dispatcher.CreatePmsJob(addr, job, keyAnswer)
}

func (p *Plugin) handleFailed(addr string, job *order.Job, reason string) {
	if job == nil || job.IsDone() {
		return
	}
	job.Done()

	automate := p.automate
	dispatcher := automate.Dispatcher()

	keyAnswer := record.KeyAnswer{
		Success:       false,
		Station:       job.Station,
		Message:       reason,
		ResponseCode:  0,
		EncoderNumber: p.getEncoder(addr),
	}

	dispatcher.CreatePmsJob(addr, job, keyAnswer)
}

func (p *Plugin) cancelJob(addr string, job *order.Job, reason string) {
	if job == nil {
		return
	}
	if !job.IsDone() {
		p.handleFailed(addr, job, reason)
	}
	p.finaliseJob(addr, job, reason)
}

func (p *Plugin) finaliseJob(addr string, job *order.Job, reason string) {
	if job == nil || job.IsRemoved() {
		return
	}
	automate := p.automate
	dispatcher := automate.Dispatcher()

	automate.LogFinaliseJob(addr, job, reason)
	dispatcher.RemoveJob(job)

	driver := p.driver
	if driver != nil {
		encoder := p.getEncoder(addr)
		if encoder >= 0 {
			// pcon-local:1234567890:2:192.168.61.82#57103:1 -> pcon-local:1234567890:2:192.168.61.82#57103
			addr = strings.TrimSuffix(addr, fmt.Sprintf(":%d", encoder))
		}
		isSerialLayer := dispatcher.IsSerialDevice(addr)
		if !isSerialLayer {
			driver.FreeSlot(addr)
		}
	}
}

func (p *Plugin) finaliseJobWithoutSlot(addr string, job *order.Job, reason string) {
	if job == nil {
		return
	}
	automate := p.automate
	dispatcher := automate.Dispatcher()

	automate.LogFinaliseJob(addr, job, reason)
	dispatcher.RemoveJob(job)
}

func (p *Plugin) linkState(addr string) bool {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	// normalise device address
	deviceAddr := ifc.DeviceAddress(addr)
	linkState := dispatcher.GetLinkState(deviceAddr)

	if linkState != automatestate.LinkUp {
		log.Info().Msgf("%s addr '%s' canceled, because of link state '%s'", name, addr, linkState.String())
		return false
	}
	return true
}

func (p *Plugin) preCheck(addr string, job *order.Job) error {

	if encoder := p.getEncoder(addr); encoder > 99 {
		return errors.New("encoder number is greater than 99")
	}

	if job == nil {
		return errors.New("job nil object")
	}

	if job.Timestamp > 0 {
		timestamp := time.Unix(0, job.Timestamp)
		since := time.Since(timestamp)
		if since > keyLifetime {
			return errors.Errorf("key service deadline exceeded (%s)", since)
		}
	}

	switch job.Action {

	case order.KeyRequest, order.KeyChange:

		is2800 := p.driver.Is2800(addr)

		if is2800 {

			if job.Action == order.KeyChange {
				return errors.New("2800 protocol did not support key change")
			}

			if job.Action == order.KeyRequest {
				if g, ok := job.Context.(*record.Guest); ok {
					keyType, _ := g.GetGeneric(defines.KeyType)
					if keyType == "D" {
						return errors.New("2800 protocol did not support key type duplicate")
					}
				}
			}
		}

		break

	case order.KeyRead, order.KeyDelete:
		break

	}

	return nil
}

func (p *Plugin) getEncoder(addr string) int {
	s := strings.Split(addr, ":")
	if len(s) < 5 {
		return -1
	}
	encoder := cast.ToInt(s[4])
	return encoder
}

func (p *Plugin) getKeyAnswerTimeout(addr string) time.Duration {
	automate := p.automate
	if timeout, ok := automate.GetSetting(addr, defines.KeyAnswerTimeout, uint(0)).(uint64); ok && timeout > 0 {
		if duration := cast.ToDuration(time.Duration(timeout) * time.Second); duration > 0 {
			return duration
		}
	}

	return keyAnswerTimeout
}

func (p *Plugin) getSource(packet *ifc.LogicalPacket) uint16 {
	if packet == nil {
		return 0
	}
	driver := p.driver
	raw, _ := driver.GetField(packet, "Source")
	data := cast.ToString(raw)
	data = strings.TrimLeft(data, "0")
	source := cast.ToUint16(data)
	return source
}

func (p *Plugin) getDestination(packet *ifc.LogicalPacket) uint16 {
	if packet == nil {
		return 0
	}
	driver := p.driver
	raw, _ := driver.GetField(packet, "Destination")
	data := cast.ToString(raw)
	data = strings.TrimLeft(data, "0")
	dst := cast.ToUint16(data)
	return dst
}

func (p *Plugin) handleJob(addr string, encoder int, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	// addr:encoder
	// pcon-local:1234567890:2:192.168.61.82#57103 -> pcon-local:1234567890:2:192.168.61.82#57103:1
	addr = addr + ":" + cast.ToString(encoder)

	// handler already running
	if p.state.Exist(addr) {
		log.Error().Msgf("%s addr '%s' handler still exist", name, addr)
		p.cancelJob(addr, job, "handler still exist")
		return
	}

	// link state
	if !p.linkState(addr) {
		p.cancelJob(addr, job, "link state")
		return
	}

	// pre-check
	if err := p.preCheck(addr, job); err != nil {
		reason := err.Error()
		if job == nil {
			log.Info().Msgf("%s addr '%s' canceled, because of %s", name, addr, reason)
		} else {
			log.Info().Msgf("%s addr '%s' job-id: %d canceled, because of %s", name, addr, job.GetQueue().Id, reason)
		}
		p.cancelJob(addr, job, reason)
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	// run state maschine
	go func(addr string, job *order.Job) {

		automate := p.automate
		name := automate.Name()

		// automate pro encoder
		handleKeyAnswer := func(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {

			process := func(addr string, in *ifc.LogicalPacket, action *dispatcher.StateAction, job *order.Job) error {
				err := p.handleKeyAnswer(addr, in, action, job)
				return err
			}

			isSerialLayer := automate.Dispatcher().IsSerialDevice(addr)

			if !isSerialLayer {

				switch job.Action {

				case order.KeyRequest, order.KeyChange, order.KeyRead:

					if in.Name != template.PacketCodeCard {
						log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' because we expect '%s'", name, addr, in.Name, template.PacketCodeCard)
						break
					}

					encoder := p.getEncoder(addr)
					if source := p.getSource(in); int(source) != encoder {
						log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' for encoder: %d, expect: %d", name, addr, in.Name, source, encoder)
						break
					}

					err := process(addr, in, action, job)
					return err

				case order.KeyDelete:

					useRmtForCheckout := p.driver.UseRmtCommandForKeyDelete(addr)

					if useRmtForCheckout && in.Name != template.PacketCodeCard {
						log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' because we expect '%s'", name, addr, in.Name, template.PacketCodeCard)
						break
					}

					if !useRmtForCheckout && in.Name != template.PacketCheckout {
						log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' because we expect '%s'", name, addr, in.Name, template.PacketCheckout)
						break
					}

					err := process(addr, in, action, job)
					return err

				}

			} else {

				encoder := p.getEncoder(addr)
				source := p.getSource(in)

				if int(source) == encoder {
					err := process(addr, in, action, job)
					return err
				}

				log.Debug().Msgf("%s addr '%s' ignore incoming packet '%s' for encoder: %d, expect: %d", name, addr, in.Name, source, encoder)

			}

			// diff bilden vom letzten packet versand, und verbleibendes coder timeout neu setzen
			nextTimeout := p.getKeyAnswerTimeout(addr)
			if job.Timestamp != 0 {
				diff := time.Since(time.Unix(0, job.Timestamp))
				if diff < nextTimeout {
					nextTimeout = nextTimeout - diff
				}
			}

			action.NextTimeout = nextTimeout

			return nil

		}

		// socket layer
		automate.RegisterRule(template.PacketCodeCard, answer, handleKeyAnswer, dispatcher.StateAction{})
		automate.RegisterRule(template.PacketCheckout, answer, handleKeyAnswer, dispatcher.StateAction{})

		// serial layer
		automate.RegisterRule(template.PacketAnswer, answer, handleKeyAnswer, dispatcher.StateAction{})
		automate.RegisterRule(template.PacketAnswerData, answer, handleKeyAnswer, dispatcher.StateAction{})

		err := automate.StateMaschine(
			addr,
			job,
			nil,
			nil,
			p.kill,
			template.PacketCodeCard,
			template.PacketCheckout,
			template.PacketAnswer,
			template.PacketAnswerData,
		)

		p.state.Remove(addr)

		if err != nil {
			p.cancelJob(addr, job, err.Error())
		} else {
			dispatcher := automate.Dispatcher()
			dispatcher.SetAlive(addr)
			p.finaliseJob(addr, job, "")
		}

		p.waitGroup.Done()

	}(addr, job)

}
