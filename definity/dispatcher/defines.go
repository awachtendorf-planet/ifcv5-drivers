package definity

import (
	"time"
)

// ...
const (
	RetryDelay      = 5 * time.Second
	PacketTimeout   = 8 * time.Second
	AnswerTimeout   = 5 * time.Second
	AliveTimeout    = 30 * time.Second
	PmsTimeout      = 10 * time.Second
	PmsSyncTimeout  = 30 * time.Second // database sync
	NextActionDelay = 0
	MaxError        = 9
	MaxNetworkError = 1000
)

const (
	ASCII_MODE       = 1
	TRANSPARENT_MODE = 2
	STANDARD_FORMAT  = 1 // 5/15 digit
	EXTENDED_FORMAT  = 2 // 7/30 digit
)

const (
	DEFAULT_ENCODING       = "ISO 8859-1" // ascii character set
	DEFAULT_LANGUAGE       = "20"         // 2 digit, 20 us
	DEFAULT_RESTRICT_LEVEL = "02"         // 2 digit, 02 station to station, 04 total restriction ?
	DEFAULT_PASSWORD       = "    "       // 4 digit, blank for default
)
