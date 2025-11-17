package robobar

import (
	"time"
)

// ...
const (
	RetryDelay      = 5 * time.Second
	PacketTimeout   = 8 * time.Second
	AliveTimeout    = 60 * time.Second
	AnswerTimeout   = 10 * time.Second
	PmsTimeout      = 10 * time.Second
	PmsSyncTimeout  = 30 * time.Second // database sync
	NextActionDelay = 0
	MaxError        = 9
	MaxNetworkError = 1000
)

func (d *Dispatcher) setDefaults() {}
