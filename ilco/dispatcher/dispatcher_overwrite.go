package ilco

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/weareplanet/ifcv5-drivers/ilco/template"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

var (
	interestedIn = map[string][]string{
		"KR": {"KT", "K#", "KO", "KC", "SI", "GF", "GN", "GDN", "RN", "G#", "GD", "T2"},
	}
)

const (
	ifcType = "KES"
)

func (d *Dispatcher) initOverwrite() {

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.PreCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

	// always true beside LinkAlive packet because fucking ahl protocol logic
	d.Acknowledgement = d.acknowledgement

	d.ConfigChanged = d.configChanged

	d.GetDriverAddr = d.getDriverAddr

}

func (d *Dispatcher) acknowledgement(addr string, packetName string) bool {
	if packetName == template.PacketSof {
		return false
	}
	return true
}

func (d *Dispatcher) PreCheck(job *order.Job) error {
	if job == nil {
		return nil
	}
	switch job.Action {
	case order.KeyRequest:
		// OK
	default:
		return errors.Errorf("vendor does not support this command (%s)", job.Action.String())
	}
	guest, ok := job.Context.(*record.Guest)
	if !ok {
		return errors.Errorf("context '%T' not supported", job.Context)
	}

	reservationNumber := guest.Reservation.ReservationID

	if _, err := strconv.Atoi(reservationNumber); err != nil {
		return errors.Errorf("reservation id '%s' is not numeric", reservationNumber)
	}

	mainRoom := guest.Reservation.RoomNumber

	if _, err := strconv.Atoi(mainRoom); err != nil {
		return errors.Errorf("room name '%s' is not numeric", mainRoom)
	}

	if len(mainRoom) > RoomNumberLength {
		return errors.Errorf("room name '%s' too long (max: %d)", mainRoom, RoomNumberLength)
	} else if len(mainRoom) == 0 {
		return errors.New("empty room name")
	}

	keyType := d.GetKeyType(guest)
	if keyType == "0" && guest.Reservation.SharedInd {
		return errors.New("key request with key type new and sharer flag is not allowed")
	}

	return nil
}

func (d *Dispatcher) configChanged(station uint64) {

	oldStateT2 := d.sendTrack2Data[station]

	d.sendTrack2Data[station] = d.SendTrack2(station)

	if d.IsReady() && d.IsStationActive(station) || oldStateT2 != d.sendTrack2Data[station] { //Track2
		d.loginStation(station)
	}

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

		if k == "KR" {

			if !d.SendTrack2(station) {

				for index, element := range v {

					if element == "T2" {

						v = append(v[:index], v[index+1:]...)

						break

					}
				}
			}

		}
		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}

func (d *Dispatcher) getDriverAddr(station uint64, _ order.QueueType, job *order.Job) ([]string, bool) {

	var addrSlice []string

	withGateway := d.IsGatewayMode(station)

	addrs, exist := d.GetStationLinkUpAddrs(station)
	if !exist {
		return addrSlice, exist
	}
	if withGateway {

		addrSlice = addrs

	} else {

		var encoderNumber int
		if guest, ok := job.Context.(*record.Guest); ok {
			if tmp, exist := guest.GetGeneric(defines.EncoderNumber); exist {
				encoderNumber = cast.ToInt(tmp)
			}
		}

		for _, addr := range addrs {

			mightBeEncoder, err := d.GetDeviceNumber(addr)
			if err == nil && cast.ToInt(mightBeEncoder) == encoderNumber {

				addrSlice = []string{addr}
				break
			}
		}

	}

	return addrSlice, len(addrSlice) > 0

}
