package vc3000

import (
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PendingDelay       = 15 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 60 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)

func (d *Dispatcher) setDefaults() {

	d.InitScope("vc3000")

	d.SetMappedGenericField("WS", defines.WorkStation)
	d.SetMappedGenericField("KC", defines.EncoderNumber)
	d.SetMappedGenericField("KT", defines.KeyType)

	d.SetMappedGenericField("C", defines.KeyCount)
	d.SetMappedGenericField("A", defines.KeyOptions)
	d.SetMappedGenericField("1", defines.Track1)
	d.SetMappedGenericField("2", defines.Track2)
	d.SetMappedGenericField("I", defines.Track4)
	d.SetMappedGenericField("L", defines.AdditionalRooms)

	//d.SetMappedGenericField("ID", defines.UserID)

}
