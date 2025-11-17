package messerschmitt

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/pkg/errors"
)

var (
	interestedIn = map[string][]string{
		"KR": {"KC", "KT", "K#", "KO", "RN", "G#", "GS", "GA", "GD"},
		"KD": {"KC", "RN", "G#", "GS"},
	}
)

const (
	ifcType = "KES"
)

func (d *Dispatcher) initOverwrite() {
	d.Acknowledgement = d.acknowledgement
	d.PreCheckDriverJob = d.preCheck
	d.LoginStation = d.loginStation
}

func (d *Dispatcher) acknowledgement(addr string, packetName string) bool {

	if !d.IsSerialDevice(addr) {
		return false
	}

	return true
}

func (d *Dispatcher) preCheck(job *order.Job) error {

	if job == nil {
		return nil
	}

	switch job.Action {
	case order.KeyRequest, order.KeyDelete:
		// OK
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}

	guest, ok := job.Context.(*record.Guest)
	if !ok {
		return errors.Errorf("context '%T' not supported", job.Context)
	}

	if job.Action == order.KeyDelete && guest.Reservation.SharedInd {
		return errors.New("key delete with sharer flag is not allowed")
	}

	if keyType, err := d.getKeyType(guest); err != nil {
		return err
	} else if keyType == "0" && guest.Reservation.SharedInd {
		return errors.New("key request with key type new and sharer flag is not allowed")
	}

	if len(guest.Reservation.ReservationID) > 16 {
		return errors.Errorf("reservation id '%s' too long (max: %d)", guest.Reservation.ReservationID, 16)
	}

	if job.Action == order.KeyRequest {

		keyCount := d.getKeyCount(guest)

		if keyCount == 0 || keyCount > 9 {
			return errors.New("key count must be between 1 and 9")
		}

		if guest.Reservation.DepartureDate.IsZero() {
			return errors.New("departure date is not valid")
		}
	}

	// room, encoder check moved to automate, because we did not know yet if serial or socket communication

	return nil
}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	addr := d.GetVendorAddr(station, 0)

	ifcType := d.GetConfig(station, defines.IFCType, ifcType)

	subscribe := record.Subscribe{
		Station: station,
		IFCType: ifcType,
	}

	subscribe.MessageName = make(map[string][]string)

	for k, v := range interestedIn {
		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}
