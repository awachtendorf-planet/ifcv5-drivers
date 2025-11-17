package vc3000

import (
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

var (
	interestedIn = map[string][]string{
		// "KR": {"KC", "KT", "K#", "KO", "SI", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS", "SR", "T1", "T2", "T3"},
		"KR": {"KC", "KT", "K#", "KO", "SI", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS", "T1", "T2", "T3"},
		"KD": {"KC", "SI", "GF", "GN", "RN", "G#", "GA", "GD", "ID", "WS"},
		"KZ": {"KC", "ID", "WS"},
	}
)

func (d *Dispatcher) initOverwrite() {
	d.MovesAllowed = d.isMovesAllowed
	d.LoginStation = d.loginStation // subscribe if station is ready
	d.GetDriverAddr = d.getDriverAddr
}

func (d *Dispatcher) isMovesAllowed(station uint64, deviceNumber int) bool {
	return false
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

func (d *Dispatcher) getDriverAddr(station uint64, _ order.QueueType, _ *order.Job) ([]string, bool) {

	// nur das erste linkup device zurück geben
	//
	// für SocketLayer gilt,
	// der eigentliche Kommunikations-Kanal wird über registrierte und freie Slot's gewählt
	// wenn wir mindestens ein linkup device haben, haben wir auch mindestens einen Slot
	// Slots werden bei Anfragen an das KES System vom Automaten belegt und anschließend wieder frei gegeben
	// wenn das Ifc ein Socket Client ist (default), dann werden benötigte Verbindungen ondemand aufgebaut

	addrs, exist := d.GetStationLinkUpAddrs(station)
	if !exist {
		return addrs, exist
	}

	if d.IsSerialDevice(addrs[0]) {
		return addrs[:1], exist
	}

	// remote handle entfernen, wird durch einen freien Slot neu gesetzt, reine Kosmetik für das log file
	// pcon-4711-local:47110815:1:192.168.0.205#65231 -> pcon-4711-local:47110815:1:192.168.0.205
	if addr, err := d.GetCustomerConnectionAddr(addrs[0]); err == nil {
		return []string{addr}, exist
	}

	return addrs[:1], exist
}
