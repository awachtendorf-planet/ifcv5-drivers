package fidserv

import (
	"strings"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/slog"
)

var (
	records = []string{"GI", "GO", "GC", "RE", "WR", "WA", "WC", "KR", "KD", "KA", "KM", "KZ", "XL" /*"XM", "XT",*/, "XD", "NS", "NE"}
)

func (d *Dispatcher) initOverwrite() {
	d.MovesAllowed = d.isMovesAllowed
	d.HandleLinkState = d.handleLinkState
	d.LoginStation = d.loginStation
}

func (d *Dispatcher) isMovesAllowed(station uint64, deviceNumber int) bool {
	addr := d.GetVendorAddr(station, deviceNumber)
	return d.IsRequestedField(addr, "GC", "RO")
}

func (d *Dispatcher) loginStation(station uint64) {
	addrs, _ := d.GetStationLinkUpAddrs(station)
	for i := range addrs {
		d.handleLinkState(addrs[i], true)
	}
}

func (d *Dispatcher) handleLinkState(addr string, state bool) {

	station, err := d.GetStationAddr(addr)
	if err != nil || !state {
		return
	}

	ifcType := d.GetConfig(station, defines.IFCType, "")

	subscribe := record.Subscribe{
		Station: station,
		IFCType: ifcType,
	}

	subscribe.MessageName = make(map[string][]string)

	var guestMessageHandling = -1 // handling not required

	if d.IsRequested(addr, "XL") && d.IsRequestedField(addr, "RE", "ML") {
		guestMessageHandling = d.getGuestMessageHandling(station) // IFCDEV-19 sinnloses Bullshit Ticket für ranzige Vendors
	}

	for i := range records {

		// Hätte uns nur einer gewarnt das uns das auf die Füße fallen wird.

		/*
			if records[i] == "XL" { // IFC-121 sinnloser Bullshit Filter für ranziges PMS

				if d.IsRequested(addr, "XL") && d.IsRequested(addr, "XM") && d.IsRequested(addr, "XT") && d.IsRequested(addr, "XD") {
					continue // drop XL subscription
				}

			}
		*/

		switch guestMessageHandling {

		case 0: // use RE, discard XL

			if records[i] == "XL" {
				continue
			}

		}

		if fields, exist := d.handleRecord(addr, station, records[i]); exist {

			if guestMessageHandling == 1 && records[i] == "RE" { // discard ML from RE

				for s := range fields {

					if fields[s] == "ML" {
						fields = append(fields[:s], fields[s+1:]...)
						break
					}

				}

			}

			subscribe.MessageName[records[i]] = fields
		}

	}

	log.Debug().Msgf("link records addr: '%s' %s", addr, subscribe.MessageName)

	d.CreatePmsJob(addr, nil, subscribe)

	d.logLinkRecords(addr)
}

func (d *Dispatcher) handleRecord(addr string, station uint64, key string) ([]string, bool) {

	fields := []string{}

	if d.IsRequested(addr, key) {
		fields = d.getFields(addr, station, key)
		return fields, true
	}

	return fields, false
}

func (d *Dispatcher) getFields(addr string, station uint64, key string) []string {

	fields := []string{}
	records := d.GetLinkRecord(addr, key)

	for i := range records {
		if mapped, err := d.GetMapping("udf", station, key+records[i], false); err == nil {
			fields = append(fields, mapped)
		} else {
			fields = append(fields, records[i])
		}
	}

	return fields
}

func (d *Dispatcher) logLinkRecords(addr string) {

	if !slog.IsEnabled() {
		return
	}

	station, err := d.GetStationAddr(addr)
	if err != nil {
		return
	}

	var dev string
	var records []string

	if dev, err = d.GetDeviceAddr(addr); err != nil {
		return
	}

	// rebuild link records from stored greetings, because we need all records
	// the records defined above only describe the direction to pms

	if records, err = d.GetSettingKeys(dev, ctxLinkRecord); err != nil {
		return
	}

	linkRecords := make(map[string]string)

	for i := range records {

		key := records[i]
		fields := d.GetLinkRecord(addr, key)
		linkRecords[key] = strings.Join(fields, ",")

	}

	devNumber, _ := d.GetDeviceNumber(addr)
	peer, _ := d.GetPeerAddr(addr)

	slog.Info(slog.Ctx{Origin: "vendor", Station: station}).
		Uint16("device", devNumber).
		Str("peer", peer).
		Interface("records", linkRecords).
		Msg("link records")
}
