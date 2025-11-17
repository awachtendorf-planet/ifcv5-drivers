package requestsimple

import (
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	databasesyncstate "github.com/weareplanet/ifcv5-main/ifc/generic/state/databasesync"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
)

func (p *Plugin) handleInitRequest(addr string, _ *ifc.LogicalPacket, retry bool) (int, error) {

	if retry {
		return 0, nil
	}

	// db sync still running
	if p.state.Exist("db:" + addr) {
		automate := p.automate
		name := automate.Name()
		log.Warn().Msgf("%s addr '%s' init request still in progress", name, addr)
		return 0, nil
	}

	go p.requestDBSync(addr)

	return 0, nil

}

func (p *Plugin) requestDBSync(addr string) {

	p.state.Register("db:" + addr)
	p.waitGroup.Add(1)
	defer func() {
		p.state.Remove("db:" + addr)
		p.waitGroup.Done()
	}()

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	log.Info().Msgf("%s addr '%s' prepare database swap", name, addr)

	station, err := dispatcher.GetStationAddr(addr)
	if err != nil {
		log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
		return
	}

	// prepare cache
	if err = dispatcher.DatabaseSyncPrepare(addr); err != nil {
		log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
		return
	}

	// register ready event
	subscriber, err := dispatcher.RegisterEvents(
		databasesyncstate.DatabaseSyncEvent,
	)
	if err != nil {
		log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, err)
		return
	}

	// cleanup
	defer func() {
		dispatcher.DeregisterEvents(subscriber)
		dispatcher.DatabaseSyncCleanup(addr)
	}()

	// request databasesync from pms

	databaseSync := record.DatabaseSync{
		Station:           station,
		RequestedMessages: 2, // CO only
	}

	if _, pmsErr, sendErr := automate.PmsRequest(station, databaseSync, pmsTimeOut, "", ""); pmsErr != nil || sendErr != nil {
		if pmsErr == nil {
			pmsErr = sendErr
		}

		if !dispatcher.IsDebugMode(station) {
			log.Error().Msgf("%s addr '%s' failed, err=%s", name, addr, pmsErr)
			return
		}
	}

	start := time.Now()

	ready := false

	for {

		// recalculate the remaining time for each run
		offset := time.Since(start)
		timeout := pmsSyncTimeout - offset

		select {

		case sub := <-subscriber.GetMessages(): // event broker, ready event
			if sub == nil {
				continue
			}
			event := sub.GetPayload()

			switch event.(type) {

			case *databasesyncstate.Event:
				event := event.(*databasesyncstate.Event)
				if event.Station == station { // databasesync was fulfilled by pms
					ready = true
				}
			}

		case <-p.kill: // shutdown
			return

		case <-time.After(timeout): // timeout
			log.Error().Msgf("%s addr '%s' timeout, err=%s", name, addr, "pms did not response")
			return

		}

		if ready {
			break
		}
	}

	accountNumberLength := p.driver.GetAccountNumberLength(station)

	// create CO jobs for datachange automate
	for {

		job := dispatcher.GetSyncRecord(addr)
		if job != nil && job.Action == order.Checkout {
			if guest, ok := job.Context.(*record.Guest); ok && !guest.Reservation.SharedInd {
				if len(guest.Reservation.RoomNumber) <= 6 && len(guest.Reservation.ReservationID) <= accountNumberLength {
					guest.SetGeneric("i:Sync", true)
					dispatcher.CreateDriverJob(job.Station, job.Action, job.Context, job.Initiator)
				}
			}
		}

		if !dispatcher.GetNextSyncRecord(addr) {
			break
		}

	}

}
