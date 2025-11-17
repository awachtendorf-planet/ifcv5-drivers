package mitel

import (
	"time"
)

// ...
const (
	RetryDelay      = 5 * time.Second
	SwapRetryDelay  = 1 * time.Second
	PacketTimeout   = 8 * time.Second
	AliveTimeout    = 10 * time.Second
	PmsTimeout      = 10 * time.Second
	PmsSyncTimeout  = 30 * time.Second // database sync
	NextActionDelay = 0
	MaxError        = 3
	MaxNetworkError = 1000
)

func (d *Dispatcher) setDefaults() {}
