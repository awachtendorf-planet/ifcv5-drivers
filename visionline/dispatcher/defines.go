package visionline

import (
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 60 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	PmsTimeout         = 10 * time.Second
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

	d.InitScope("visionline")

	//d.SetMappedGenericField("WS", defines.WorkStation)
	d.SetMappedGenericField("EA", defines.EncoderNumber)
	d.SetMappedGenericField("NC", defines.KeyCount)
	d.SetMappedGenericField("T1", defines.Track1)
	d.SetMappedGenericField("T2", defines.Track2)
	d.SetMappedGenericField("T3", defines.Track3)
	d.SetMappedGenericField("IO", defines.UserID)
	d.SetMappedGenericField("SR", defines.KeyID)
	d.SetMappedGenericField("MK", defines.KeyRingID)
	d.SetMappedGenericField("CR", defines.AdditionalRooms)
	d.SetMappedGenericField("OR", defines.OldRoomName)

}
