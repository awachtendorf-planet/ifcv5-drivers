package messerschmitt

import (
	"fmt"
)

var (
	answerText = map[int]string{
		0: "coding/reading ok",
		1: "timeout lapsed",
		2: "room not found",
		3: "wrong data in checkin telegram",
		4: "invalid expiration date",
		5: "number of cards or encoder number invalid",
		6: "already 2 separate checkin (max. 2)", // invalid data
	}
)

// GetAnswerText returns a clear text for fias
func (d *Dispatcher) GetAnswerText(answerStatus int) string {
	text, exist := answerText[answerStatus]
	if !exist {
		return fmt.Sprintf("unknown answer status: %d", answerStatus)
	}
	return text
}
