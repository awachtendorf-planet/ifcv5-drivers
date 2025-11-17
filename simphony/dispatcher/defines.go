package simphony

import (
	"time"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 30 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	PmsTimeout         = 100 * time.Second
	PmsSyncTimeout     = 30 * time.Second // database sync
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)
