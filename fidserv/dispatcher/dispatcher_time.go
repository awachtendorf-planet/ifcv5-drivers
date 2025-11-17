package fidserv

import (
	"time"

	"github.com/pkg/errors"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
)

// ParseTime returns a time object from a logical packet DA/TI
func (d *Dispatcher) ParseTime(in *ifc.LogicalPacket) (time.Time, error) {

	da, _ := d.GetField(in, "DA")
	ti, _ := d.GetField(in, "TI")

	// IFC-733
	if len(ti) == 4 {
		ti = ti + "00"
	}

	return d.GetTime(da, ti)

}

// GetTime returns a time object from fias DA/TI
func (d *Dispatcher) GetTime(da string, ti string) (time.Time, error) {

	if len(da) == 6 && len(ti) == 6 {
		return time.Parse("060102150405", da+ti)
	}

	if len(da) == 6 && len(ti) == 0 {
		return time.Parse("060102", da)
	}

	if len(da) == 0 && len(ti) == 6 {
		return time.Parse("150405", ti)
	}

	if len(da) == 0 && len(ti) == 0 {
		return time.Time{}, errors.New("no time value provided")
	}

	return time.Time{}, errors.New("unknown time layout")

}
