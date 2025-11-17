package robobar

import (
	"fmt"
)

var (
	answerText = map[byte]string{
		'V': "invalid roomid or message type",
		'N': "room has no bar",
		'O': "room is not occupied",
		'I': "room is occupied",
		' ': "no error code reported",
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
