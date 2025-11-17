package ilco

import (
	"time"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	PmsTimeout         = 10 * time.Second
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000

	RoomNumberLength = 5

	KeyDisplayLine  = byte(0x20)
	KeyRoomNumber   = byte(0x22)
	KeyCheckoutDate = byte(0x24)
	KeyCardType     = byte(0x25)
	KeyNumberKeys   = byte(0x26)
	KeyAreas        = byte(0x29)
	KeyAuthNo       = byte(0x2b)
	KeyFolioNo      = byte(0x2c)
	KeyTrack2       = byte(0x46)

	FieldSeparator = byte(0x1c)
)

func (d *Dispatcher) setDefaults() {

}
