package datachange

import (
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/ticker"

	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
)

func (p *Plugin) handleNextAction(addr string, t *ticker.ResetTicker, job *order.Job) {

	automate := p.automate
	name := automate.Name()

	if !p.linkState(addr) {
		automate.NextAction(name, addr, shutdown, t)
		return
	}

	var err error
	state := automate.GetState(addr)
	action := dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: nextAction}

	lastTask := 0

	individualRoomStatusPackets, ok := automate.GetSetting(addr, "IndividualRoomStatusPackets", false).(bool)
	if !ok {
		individualRoomStatusPackets = false
	}

	aswWithoutRoomStatus := true

	if job.Action == order.Checkin || job.Action == order.Checkout || job.Action == order.DataChange || job.Action == order.RoomStatus {
		if individualRoomStatusPackets {
			lastTask = 7
		} else {
			lastTask = 1
		}
	}

	switch state {

	case busy:

		if job == nil || !job.InProcess() { // order was canceled externally
			err = defines.ErrorJobFailed
			break
		}

		if job.Context == nil && job.Action != order.NightAuditStart && job.Action != order.NightAuditEnd {
			err = defines.ErrorJobContext
			break
		}

		driver := p.driver
		requested := true

		// room status skip task 0 (asw, wake, eod)
		if job.Action == order.RoomStatus && job.Task == 0 {
			job.Task = 1
		}

		// send only ASW, without rights/roomstatus
		if aswWithoutRoomStatus {
			if job.Task == 1 && (job.Action == order.Checkin || job.Action == order.Checkout || job.Action == order.DataChange) {
				automate.NextAction(name, addr, success, t)
				return
			}
		}

		// send only ASW if CO, b/c it can be a CO from room-move and we dont have the data from the remaining room/guest
		if job.Task == 1 && job.Action == order.Checkout {
			automate.NextAction(name, addr, success, t)
			return
		}

		// skip task 1 if we use individual room status packets
		if job.Task == 1 && individualRoomStatusPackets {
			job.Task = 2
		}

		switch job.Task {

		case 0: // checkin, checkout, datachange, roomstatus, wake, eod

			switch job.Action {

			// guest asw
			case order.Checkin:
				if requested = driver.IsRequested(addr, "GI"); !requested {
					break
				}
				err = p.sendCheckIn(addr, &action, job.Initiator, job.Context)

			case order.Checkout:
				if requested = driver.IsRequested(addr, "GO"); !requested {
					break
				}
				err = p.sendCheckOut(addr, &action, job.Initiator, job.Context)

			case order.DataChange:
				if requested = driver.IsRequested(addr, "GC"); !requested {
					break
				}
				err = p.sendDataChange(addr, &action, job.Initiator, job.Context)

			// room status
			case order.RoomStatus:
				p.nextTask(name, addr, job, lastTask, t)
				return

			// wakeup
			case order.WakeupRequest:
				if requested = driver.IsRequested(addr, "WR"); !requested {
					break
				}
				action = dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}
				err = p.sendWakeupRequest(addr, &action, job.Initiator, job.Context)

			case order.WakeupClear:
				if requested = driver.IsRequested(addr, "WC"); !requested {
					break
				}
				action = dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}
				err = p.sendWakeupClear(addr, &action, job.Initiator, job.Context)

			// night audit
			case order.NightAuditStart:
				if requested = driver.IsRequested(addr, "NS"); !requested {
					break
				}
				action = dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}
				err = p.sendNightAuditStart(addr, &action, job.Initiator, job.Context)

			case order.NightAuditEnd:
				if requested = driver.IsRequested(addr, "NE"); !requested {
					break
				}
				action = dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}
				err = p.sendNightAuditEnd(addr, &action, job.Initiator, job.Context)

			// guest message online
			case order.GuestMessageOnline:
				if requested = driver.IsRequested(addr, "XL"); !requested {
					break
				}

				switch job.Context.(type) {

				case *record.Guest:

					guest := job.Context.(*record.Guest)
					if message, exist := guest.GetGeneric(defines.GuestMessage); exist {

						switch message.(type) {

						case record.GuestMessage:
							data := message.(record.GuestMessage)
							msg := record.NewGeneric()
							msg.Set("G#", data.ReservationID)
							msg.Set("RN", data.RoomNumber)
							msg.Set("MI", data.MessageID)
							msg.Set("MT", data.Text)

							action = dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}
							err = p.sendGuestMessageOnline(addr, &action, job.Initiator, msg)

						default:

							log.Error().Msgf("%s addr '%s' canceled, because of action '%s' has no valid message object, expected 'record.GuestMessage', got '%T'", name, addr, job.Action.String(), message)
							automate.NextAction(name, addr, success, t)
							return
						}
					}

				default:

					log.Error().Msgf("%s addr '%s' canceled, because of action '%s' has no valid context object, expected '*record.Guest', got '%T'", name, addr, job.Action.String(), job.Context)
					automate.NextAction(name, addr, success, t)
					return
				}

			// guest message delete
			case order.GuestMessageDelete:
				if requested = driver.IsRequested(addr, "XD"); !requested {
					break
				}

				switch job.Context.(type) {

				case *record.Guest:

					guest := job.Context.(*record.Guest)
					if message, exist := guest.GetGeneric(defines.GuestMessage); exist {

						switch message.(type) {

						case record.GuestMessage:
							data := message.(record.GuestMessage)
							msg := record.NewGeneric()
							msg.Set("G#", data.ReservationID)
							msg.Set("RN", data.RoomNumber)
							msg.Set("MI", data.MessageID)

							action = dispatcher.StateAction{NextTimeout: nextActionDelay, CurrentState: state, NextState: success}
							err = p.sendGuestMessageDelete(addr, &action, job.Initiator, msg)

						default:

							log.Error().Msgf("%s addr '%s' canceled, because of action '%s' has no valid message object, expected 'record.GuestMessage', got '%T'", name, addr, job.Action.String(), message)
							automate.NextAction(name, addr, success, t)
							return
						}
					}

				default:

					log.Error().Msgf("%s addr '%s' canceled, because of action '%s' has no valid context object, expected '*record.Guest', got '%T'", name, addr, job.Action.String(), job.Context)
					automate.NextAction(name, addr, success, t)
					return
				}

			default:
				automate.NextAction(name, addr, success, t)
				return
			}

		case 1: // room equipment status, one packet
			if requested = driver.IsRequested(addr, "RE"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context)

		case 2: // switch auth, class of service
			if requested = driver.IsRequestedField(addr, "RE", "CS"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context, "RN", "CS")

		case 3: // switch DND
			if requested = driver.IsRequestedField(addr, "RE", "DN"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context, "RN", "DN")

		case 4: // switch message lamp
			if requested = driver.IsRequestedField(addr, "RE", "ML"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context, "RN", "ML", "G#")

		case 5: // switch minibar rights
			if requested = driver.IsRequestedField(addr, "RE", "MR"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context, "RN", "MR")

		case 6: // switch TV rights
			if requested = driver.IsRequestedField(addr, "RE", "TV"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context, "RN", "TV")

		case 7: // room status
			if requested = driver.IsRequestedField(addr, "RE", "RS"); !requested {
				break
			}
			err = p.sendRoomData(addr, &action, job.Initiator, job.Context, "RN", "RS")

		default: // done
			automate.NextAction(name, addr, success, t)
			return

		}

		if !requested {
			p.nextTask(name, addr, job, lastTask, t)
			return
		}

	case nextAction:
		p.nextTask(name, addr, job, lastTask, t)
		return

	case success, shutdown:
		t.Tick()
		return

	default:
		action = dispatcher.StateAction{NextTimeout: retryDelay}
	}

	if err != nil {
		log.Error().Msgf("%s %s", name, err)
		if err == defines.ErrorJobContext || err == defines.ErrorJobFailed {
			automate.NextAction(name, addr, shutdown, t)
			return
		}
		automate.ChangeState(name, addr, state)
	}

	automate.SetNextTimeout(addr, action, err, t)
}

func (p *Plugin) nextTask(name string, addr string, job *order.Job, lastTask int, t *ticker.ResetTicker) {
	automate := p.automate
	if job == nil {
		automate.NextAction(name, addr, shutdown, t)
		return
	}
	if job.Task >= lastTask {
		automate.NextAction(name, addr, success, t)
	} else {
		job.Task++
		automate.NextAction(name, addr, busy, t)
	}
}

func (p *Plugin) logOutgoingPacket(packet *ifc.LogicalPacket) {
	name := p.automate.Name()
	p.driver.LogOutgoingPacket(name, packet)
}

func (p *Plugin) sendCheckIn(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketCheckIn, tracking, "GI", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendCheckOut(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketCheckOut, tracking, "GO", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendDataChange(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketDataChange, tracking, "GC", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendRoomData(addr string, action *dispatcher.StateAction, tracking string, context interface{}, fields ...string) error {
	packet := p.driver.ConstructPacket(addr, template.PacketRoomData, tracking, "RE", context, fields...)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendWakeupRequest(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketWakeupRequest, tracking, "WR", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendWakeupClear(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketWakeupClear, tracking, "WC", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendNightAuditStart(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketNightAuditStart, tracking, "NS", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendNightAuditEnd(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketNightAuditEnd, tracking, "NE", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendGuestMessageOnline(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketGuestMessageOnline, tracking, "XL", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}

func (p *Plugin) sendGuestMessageDelete(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {
	packet := p.driver.ConstructPacket(addr, template.PacketGuestMessageDelete, tracking, "XD", context)
	p.logOutgoingPacket(packet)
	err := p.automate.SendPacket(addr, packet, action)
	return err
}
