package robobar

/*
import (
	"github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/spf13/cast"
)

func (d *Dispatcher) databaseSyncPreProcess(station uint64, context interface{}, action order.Action) *order.Job {

	if action == order.Checkout {
		if guest, ok := context.(*record.Guest); ok {
			guest.Rights.Minibar = 0 // lock bar
			job := &order.Job{
				Station: station,
				Action:  order.RoomStatus,
				Context: guest,
			}
			return job
		}
	}

	return nil
}

func (d *Dispatcher) databaseSyncPostProcess(station uint64, context interface{}, action order.Action) *order.Job {

	if action == order.Checkin {

		if guest, ok := context.(*record.Guest); ok {

			data := d.GetConfig(station, "UnlockBar", "true")
			if cast.ToBool(data) {
				guest.Rights.Minibar = 2 // unlock bar
			} else {
				guest.Rights.Minibar = 0 // lock bar
			}

			job := &order.Job{
				Station: station,
				Action:  order.RoomStatus,
				Context: guest,
			}
			return job
		}
	}

	return nil
}
*/
