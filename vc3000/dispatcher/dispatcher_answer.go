package vc3000

import (
	"fmt"
)

var (
	answerText = map[byte]string{
		'0': "OK",
		'1': "unspecified error",
		'2': "illegal device address (dd value)",
		'3': "illegal command code (ff value)",
		'4': "room is in use",
		'5': "device busy",
		'6': "cannot make more cards for room",
		'7': "key code already in use",
		'8': "device time-out",
		'9': "room not occupied",
		'R': "error with R or L Field (room name or room list)",
		'D': "error with D field (checkin time)",
		'O': "error with O field (checkout time)",
		'T': "error with T field (user type)",
		'U': "error with U field (user group)",
		'J': "error with J field (card type) ",
		'S': "error with S field (card serial number)",
		'V': "error with V field (unique code)",
		0:   "no error code reported",
	}
)

// GetAnswerText returns a clear text for fias
func (d *Dispatcher) GetAnswerText(answerStatus byte) string {
	text, exist := answerText[answerStatus]
	if !exist {
		return fmt.Sprintf("unknown error code '%c'", answerStatus)
	}
	return text
}
