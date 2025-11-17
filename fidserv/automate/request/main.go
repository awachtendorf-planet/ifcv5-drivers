package request

import (
	"fmt"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/packetqueue"
	automatestate "github.com/weareplanet/ifcv5-main/ifc/generic/state/automate"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/spf13/cast"
)

const (
	wildcardAddress  = "*"
	respectLinkState = true
)

func (p *Plugin) main() error {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	driver := p.driver
	name := automate.Name()

	route, err := dispatcher.RegisterPacketRoute(name, wildcardAddress,
		template.PacketPostingSimple,
		template.PacketPostingRequest,
		template.PacketRoomData,
		template.PacketWakeupRequest,
		template.PacketWakeupClear,
		template.PacketWakeupAnswer,
		template.PacketGuestMessageRequest,
		template.PacketGuestMessageDelete,
		template.PacketGuestBillRequest,
		template.PacketRemoteCheckOut,
		template.PacketGuestCheckDetails,
	)
	if err != nil {
		return err
	}

	defer func() {
		dispatcher.DeregisterPacketRoute(name, wildcardAddress)
	}()

	queue := packetqueue.NewPacketQueue()

	for {

		var nextPacket = (*ifc.LogicalPacket)(nil)

		select {

		// shutdown
		case <-p.kill:
			return nil

		// packet completed
		case packetAddr := <-p.done:
			if queueAddr, err := p.getQueueAddr(packetAddr); err == nil {
				if next := queue.Remove(queueAddr); next != nil {
					nextPacket, _ = next.(*ifc.LogicalPacket)
				} else {
					log.Debug().Msgf("%s addr '%s' packet queue '%s' is empty", name, packetAddr, queueAddr)
				}
			}

		// incoming logical packet
		case packet := <-route.C:
			if packet == nil {
				return nil
			}

			payload, _ := driver.GetPayload(packet)
			log.Info().Msgf("%s addr '%s' received packet '%s' %s", name, packet.Addr, packet.Name, payload)

			if reason, success := p.preCheck(packet); !success {
				log.Info().Msgf("%s addr '%s' packet '%s' discarded, because of %s", name, packet.Addr, packet.Name, reason)
				continue
			}

			if queueAddr, err := p.getQueueAddr(packet.Addr); err == nil {
				if next := queue.Append(queueAddr, packet); next != nil {
					nextPacket, _ = next.(*ifc.LogicalPacket)
				} else {
					pending := queue.Length(queueAddr)
					if pending > 1 {
						log.Debug().Msgf("%s addr '%s' packet queue '%s' still in progress (%d pending packets)", name, packet.Addr, queueAddr, pending)
					} else {
						log.Debug().Msgf("%s addr '%s' packet queue '%s' still in progress (%d pending packet)", name, packet.Addr, queueAddr, pending)
					}
				}
			}

		}

		if nextPacket != nil {
			p.handleRequest(nextPacket.Addr, nextPacket)
		}

	}

	return nil
}

func (p *Plugin) preCheck(packet *ifc.LogicalPacket) (string, bool) {
	if packet == nil {
		return "packet nil object", false
	}

	driver := p.driver

	switch packet.Name {

	case template.PacketGuestBillRequest:
		if !driver.IsRequested(packet.Addr, "XI") && !driver.IsRequested(packet.Addr, "XB") {
			return "XI/XB not requested", false
		}

	case template.PacketGuestMessageRequest:
		if !driver.IsRequested(packet.Addr, "XT") {
			return "XT not requested", false
		}

	case template.PacketRemoteCheckOut:
		if !driver.IsRequested(packet.Addr, "XC") {
			return "XC not requested", false
		}

	case template.PacketGuestCheckDetails:
		if !driver.IsRequested(packet.Addr, "CK") {
			return "CK not requested", false
		}

	}

	return "", true
}

func (p *Plugin) getQueueAddr(addr string) (string, error) {
	station, err := p.getStationAddr(addr)
	if err == nil {
		addr = cast.ToString(station)
		return addr, nil
	}
	return "", err
}

func (p *Plugin) getStationAddr(addr string) (uint64, error) {
	automate := p.automate
	dispatcher := automate.Dispatcher()
	station, err := dispatcher.GetStationAddr(addr)
	return station, err
}

func (p *Plugin) finalisePacket(packet *ifc.LogicalPacket) {
	if packet == nil {
		return
	}
	p.done <- packet.Addr
}

func (p *Plugin) linkState(addr string) bool {

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	linkState := dispatcher.GetLinkState(addr)
	if linkState != automatestate.LinkUp {
		log.Info().Msgf("%s addr '%s' canceled, because of link state '%s'", name, addr, linkState.String())
		return false
	}
	return true
}

func (p *Plugin) createJob(addr string, packet *ifc.LogicalPacket) *order.Job {
	job := order.NewJob(order.Internal, order.Ifc, packet)
	id := packet.ID
	station, _ := p.getStationAddr(addr)
	jobAddr := fmt.Sprintf("%d:%d:%d", station, 0, 0)
	job.Station = station
	job.SetQueue(uint64(id), jobAddr)
	job.Process()
	return job
}

func (p *Plugin) logOutgoingPacket(packet *ifc.LogicalPacket) {
	name := p.automate.Name()
	p.driver.LogOutgoingPacket(name, packet)
}

func (p *Plugin) handleRequest(addr string, packet *ifc.LogicalPacket) {
	automate := p.automate
	name := automate.Name()

	if p.state.Exist(addr) {
		log.Warn().Msgf("%s addr '%s' request handler still exist", name, addr)
		return
	}

	p.state.Register(addr)
	p.waitGroup.Add(1)

	go func(addr string, packet *ifc.LogicalPacket) {

		job := p.createJob(addr, packet)

		err := automate.StateMaschine(
			// device address
			addr,
			// Job
			job,
			// preProcess
			nil,
			// postProcess
			nil,
			// shutdown chan
			p.kill,
			// templates
		)

		automate.LogFinaliseJob(addr, job, automate.ErrorString(err))

		p.state.Remove(addr)

		p.finalisePacket(packet)

		if err == nil {
			dispatcher := automate.Dispatcher()
			dispatcher.SetAlive(addr)
		}

		p.waitGroup.Done()

	}(addr, packet)
}
