package caracas

import (
	"time"
	// "github.com/weareplanet/ifcv5-main/ifc/defines"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 30 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	PmsTimeout         = 10 * time.Second
	PmsSyncTimeout     = 30 * time.Second // database sync
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)
