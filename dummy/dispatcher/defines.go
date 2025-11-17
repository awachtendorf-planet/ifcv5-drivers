package dummy

import (
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 30 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	PmsTimeout         = 10 * time.Second
	PmsSyncTimeout     = 30 * time.Second // database sync
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)

const (
	// payload object from templates
	payload = "Data"
)

func (d *Dispatcher) setDefaults() {

	d.InitScope("dummy")

	// DT = Checkout Time, mapped to GD
	if value, exist := d.GetMappedGuestField("GD"); exist {
		d.SetMappedGuestField("DT", value)
	}

	d.SetMappedGenericField("WS", defines.WorkStation)
	d.SetMappedGenericField("KC", defines.EncoderNumber)
	d.SetMappedGenericField("KT", defines.KeyType)
	d.SetMappedGenericField("K#", defines.KeyCount)
	d.SetMappedGenericField("KO", defines.KeyOptions)
	d.SetMappedGenericField("$1", defines.Track1)
	d.SetMappedGenericField("$2", defines.Track2)
	d.SetMappedGenericField("$3", defines.Track3)
	d.SetMappedGenericField("$4", defines.Track4)
	d.SetMappedGenericField("ID", defines.UserID)
	d.SetMappedGenericField("RO", defines.OldRoomName)
	d.SetMappedGenericField("RS", defines.RoomStatus)
	d.SetMappedGenericField("RN", defines.WakeExtension)
	d.SetMappedGenericField("SI", defines.AdditionalRooms)

	d.SetMappedGenericField("DA", defines.Timestamp)
	d.SetMappedGenericField("TI", defines.Timestamp)

}
