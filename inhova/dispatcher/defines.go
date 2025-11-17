package inhova

import (
	"time"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	KeyAnswerTimeout   = 30 * time.Second
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)
