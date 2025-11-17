package guestlink

import (
	"time"
)

// ...
const (
	RetryDelay      = 5 * time.Second
	PacketTimeout   = 8 * time.Second
	AliveTimeout    = 90 * time.Second
	AnswerTimeout   = 15 * time.Second
	PmsTimeout      = 10 * time.Second
	PmsSyncTimeout  = 30 * time.Second // database sync
	NextActionDelay = 0
	MaxError        = 3
	MaxNetworkError = 1000
)

const (
	guestlink_tan   = "Transaction"
	guestlink_seq   = "Sequence"
	guestlink_error = "ErrorCode"
)

const (
	DrvGuestlink       = 1
	DrvEclipse         = 2
	DrvSonifi          = 3
	DrvMovielink       = 4
	DrvOnCommand       = 5
	DrvQuadriga        = 6
	DrvTripleGuest     = 7
	DrvMagiNet         = 8
	DrvMaginetEnhanced = 9
)

const (
	JustifiedLeft  = 0
	JustifiedRight = 1
)

func (d *Dispatcher) setDefaults() {
}
