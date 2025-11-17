package fidserv

import (
	"strings"

	"github.com/weareplanet/ifcv5-main/log"
)

const (
	ctxLinkRecord = "LR"
)

func (d *Dispatcher) setDefaultRecord(dev string) {
	// set default for link records
	d.records.Set(dev, "LS", "DATI")
	d.records.Set(dev, "LA", "DATI")
	d.records.Set(dev, "DR", "DATI")
	d.records.Set(dev, "DS", "DATI")
	d.records.Set(dev, "DE", "DATI")

	// set default for generic answer records

	//d.records.Set(dev, "NS", "DATI")
	//d.records.Set(dev, "NE", "DATI")

	d.records.Set(dev, "PA", "ASCTP#RNWSC#G#GNIDSODATI")
	d.records.Set(dev, "PL", "G#GNP#RNWSBAC#CLDAG+GAGDGFGGGLGTGVIDNPPMSOTI")
	d.records.Set(dev, "CK", "ASC#CTDATISO")

	// d.records.Set(dev, "XB", "ASBAG#RNDATI")
	// d.records.Set(dev, "XI", "BDBIDCG#RNF#FDDATI")
	// d.records.Set(dev, "XC", "ASBACTG#RNDATI")
	// d.records.Set(dev, "XT", "G#RNDATIMIMT")
}

// NewConnection initialised a new connection
func (d *Dispatcher) NewConnection(addr string) {
	dev, err := d.GetDeviceAddr(addr)
	if err != nil {
		log.Error().Msgf("%s", err)
		return
	}
	d.records.Delete(dev)
	d.setDefaultRecord(dev)

	// load last stored link records
	settings := d.LoadSetting(dev, ctxLinkRecord)
	for i, key := range settings {
		d.records.Set(dev, i, key)
	}
}

// ClearLinkDescription clear the current link discription
func (d *Dispatcher) ClearLinkDescription(addr string) {
	dev, err := d.GetDeviceAddr(addr)
	if err != nil {
		return
	}
	d.records.Delete(dev)
	d.setDefaultRecord(dev)
	d.DeleteSetting(dev, ctxLinkRecord, "")
}

// SetLinkRecord set a new link record
func (d *Dispatcher) SetLinkRecord(addr string, key string, value string) {
	dev, err := d.GetDeviceAddr(addr)
	if err != nil {
		log.Warn().Msgf("%s", err)
		return
	}
	if len(value)%2 != 0 {
		log.Warn().Msgf("%T link record '%s' with fields '%s' failed, odd length", d, key, value)
		return
	}
	d.records.Set(dev, key, value)
	d.WriteSetting(dev, ctxLinkRecord, key, value)

}

// GetLinkRecord return a link record
func (d *Dispatcher) GetLinkRecord(addr string, key string) []string {
	var record = []string{}
	dev, err := d.GetDeviceAddr(addr)
	if err != nil {
		log.Warn().Msgf("%s", err)
		return record
	}
	if records, exist := d.records.Get(dev, key); exist {
		for i := 0; i < len(records); i = i + 2 {
			record = append(record, records[i:i+2])
		}
	}
	return record
}

// IsRequested return true if a record id is requested
func (d *Dispatcher) IsRequested(addr string, key string) bool {
	dev, err := d.GetDeviceAddr(addr)
	if err != nil {
		log.Warn().Msgf("%s", err)
		return false
	}
	_, exist := d.records.Get(dev, key)
	return exist
}

// IsRequestedField return true if any of the fields are requested by record id
// or return true if none fields are provided, and the record id is requested
func (d *Dispatcher) IsRequestedField(addr string, key string, fields ...string) bool {
	dev, err := d.GetDeviceAddr(addr)
	if err != nil {
		log.Warn().Msgf("%s", err)
		return false
	}
	keys, exist := d.records.Get(dev, key)
	if !exist {
		return false
	}
	if len(fields) == 0 {
		return true
	}
	for i := range fields {
		if index := strings.Index(keys, fields[i]); index > -1 {
			if index%2 == 0 {
				return true
			}
		}
	}
	return false
}
