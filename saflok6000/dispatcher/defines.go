package saflok6000

import (
	"time"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 30 * time.Second
	KeyAnswerTimeout   = 30 * time.Second
	PmsTimeout         = 10 * time.Second
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)

func (d *Dispatcher) setDefaults() {

}
