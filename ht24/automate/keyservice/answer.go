package keyservice

import (
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

func (p *Plugin) handleSuccess(addr string, job *order.Job, keyId string) {
	if job == nil || job.IsDone() {
		return
	}
	job.Done()

	keyAnswer := record.KeyAnswer{
		Success:       true,
		Track2:        keyId,
		EncoderNumber: p.getEncoder(addr),
	}

	automate := p.automate
	dispatcher := automate.Dispatcher()

	keyAnswer.Station = job.Station
	dispatcher.CreatePmsJob(addr, job, keyAnswer)
}

func (p *Plugin) handleError(addr string, job *order.Job, reason string) {
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
