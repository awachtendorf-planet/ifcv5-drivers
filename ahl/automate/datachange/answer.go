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

func (p *Plugin) handleAnswerNegativ(addr string, job *order.Job, reason byte) {
	driver := p.driver
	errorMessage := driver.GetAnswerText(reason)
	p.handleError(addr, job, errorMessage)
}

func (p *Plugin) handleError(addr string, job *order.Job, reason string) {

	if job == nil {
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
