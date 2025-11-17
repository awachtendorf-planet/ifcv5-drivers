package ahl

import (
	"time"
)

// ...
const (
	RetryDelay      = 5 * time.Second
	PacketTimeout   = 8 * time.Second
	AnswerTimeout   = 8 * time.Second
	AliveTimeout    = 30 * time.Second
	PmsTimeout      = 10 * time.Second
	NextActionDelay = 0
	MaxError        = 3
	MaxNetworkError = 1000
)

const (
	AHL4400_5 = 1
	AHL4400_8 = 2
)
