package detewe

var (
	answerText = map[string]string{
		"411": "Data were not entered or are not available",
		"413": "Registered/deleted as cross-linked",
		"602": "Not configurable",
		"603": "Participant not available",
		"672": "Programming call forwarding not possible",
		"673": "Participant or destination not available",
		"702": "Participant busy",
		"703": "Participant did not answer",
		"704": "Participant not available",
		"710": "No more wake-up orders available, all deleted",
		"711": "New alarm time programmed",
		"712": "Alarm time not programmed",
		"715": "Wake-up order found",
		"716": "No wake-up order found",
		"802": "Entry not possible",
		"803": "Participant not available",
	}
)

// GetAnswerText returns a clear text for fias
func (d *Dispatcher) GetAnswerText(answerStatus string) string {
	text, exist := answerText[answerStatus]
	if !exist {
		return answerStatus
	}
	return text
}
