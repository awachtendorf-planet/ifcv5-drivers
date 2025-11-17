package visionline

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

var (
	interestedIn = map[string][]string{
		"KR": {"KC", "KT", "K#", "KO", "GS", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS", "SI", "SR", "MK", "T1", "T2", "T3"},
		"KM": {"KC", "KT", "K#", "KO", "GS", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS", "SI", "SR", "MK", "T1", "T2", "T3", "RO"},
		"KD": {"KC" /*             */, "GS", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS", "SI"},
		"KZ": {"KC" /*             */, "G#", "GD", "ID", "WS"},
	}
)

func (d *Dispatcher) initOverwrite() {
	d.MovesAllowed = d.isMovesAllowed
	d.LoginStation = d.loginStation // subscribe if station is ready
	//d.Acknowledgement = func(addr string) bool { return true }
}

func (d *Dispatcher) isMovesAllowed(station uint64, deviceNumber int) bool {
	return true
}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	addr := d.GetVendorAddr(station, 0)

	ifcType := d.GetConfig(station, defines.IFCType, "KES")

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
