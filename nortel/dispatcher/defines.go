package nortel

import (
	"time"
)

// ...
const (
	RetryDelay     = 5 * time.Second
	PmsTimeout     = 10 * time.Second
	PmsSyncTimeout = 30 * time.Second

	packetTimeout       = 8 * time.Second
	bgdTimeout          = 30 * time.Second
	packetAnswerTimeout = 1500 * time.Millisecond

	NextActionDelay = 0
	MaxError        = 9
	MaxNetworkError = 1000
)

func (d *Dispatcher) setDefaults() {
}
