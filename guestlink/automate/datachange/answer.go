package datachange

import (
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleAnswerPositiv(addr string, job *order.Job) {

	if job == nil {
		return
	}

	if p.isInitResponse(job.Context) {
		return
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	log.Info().Msgf("%s addr '%s' action '%s' was successful", name, addr, job.Action.String())

	response := record.AsyncAnswer{
		Success:       true,
		CorrelationID: job.Initiator,
		Station:       job.Station,
		Action:        int(job.Action),
	}

	dispatcher.CreatePmsJob(addr, job, response)
}

func (p *Plugin) handleAnswerNegativ(addr string, job *order.Job, errorCode int) {
	driver := p.driver
	errorMessage := driver.GetAnswerText(errorCode)
	p.handleError(addr, job, errorMessage)
}

func (p *Plugin) handleError(addr string, job *order.Job, reason string) {

	if job == nil {
		return
	}

	if p.isInitResponse(job.Context) {
		return
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()
	name := automate.Name()

	log.Error().Msgf("%s addr '%s' action '%s' failed, err=%s", name, addr, job.Action.String(), reason)

	response := record.AsyncAnswer{
		Success:       false,
		CorrelationID: job.Initiator,
		Station:       job.Station,
		Action:        int(job.Action),
		ResponseCode:  0,
		ResponseText:  reason,
	}

	dispatcher.CreatePmsJob(addr, job, response)
}

func (p *Plugin) isInitResponse(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		if value, exist := guest.GetGeneric("i:Sync"); exist {
			if sync, ok := value.(bool); ok && sync {
				return true
			}
		}
	}
	return false
}

func (p *Plugin) isSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}
	return false
}
