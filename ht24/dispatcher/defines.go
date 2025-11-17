package ht24

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

	RoomNumberLength = 7
)

func (d *Dispatcher) setDefaults() {

}
