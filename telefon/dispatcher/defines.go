package telefon

import (
	"time"
)

// ...
const (
	RetryDelay      = 5 * time.Second
	PacketTimeout   = 8 * time.Second
	NextActionDelay = 0
	MaxError        = 9

	MaxNetworkError = 1000
)

func (d *Dispatcher) setDefaults() {}
