package guestlink

import (
	"fmt"
)

var (
	answerText = map[int]string{
		0:  "Undefined Error",
		1:  "Unknown Command Verb",
		2:  "Unknown Room Number",
		3:  "Room Unoccupied",
		4:  "Unknown Account Number",
		5:  "Account Number not Checked into Room",
		6:  "Invalid Method of Payment",
		7:  "Account Balance Changed",
		8:  "Unknown Maid Code",
		9:  "Night Audit in Progress",
		10: "Locked Folio",
		11: "Guest Message Not Found",
		12: "Guest Message Cannot Be Delivered",
	}
)

// GetAnswerText returns a clear text for fias
func (d *Dispatcher) GetAnswerText(errorCode int) string {
	text, exist := answerText[errorCode]
	if !exist {
		return fmt.Sprintf("unknown error code '%d'", errorCode)
	}
	return text
}
