package inhova

var (
	answerText = map[string]string{
		"E1": "The command cannot be accomplished because the user cannot be recognized, or the command is not applicable.",
		"E2": "Syntax Error.",
		"E3": "Card lodged in the encoder. The card does not have the required plastic size or opacity.",
		"E8": "The encoder has been waiting too long for the card.",
		"EA": "The encoder does not answer. Failure in the communications.",
		"ED": "Room is not occupied.",
		"EU": "The user does not exist or it has been deleted.",
		"EE": "General Reading error. The reading operation is not successful.",
		"EF": "General Encoding error. The encoding operation is not successful.",
	}
)

// GetAnswerText returns a clear text
func (d *Dispatcher) GetAnswerText(code string) string {
	text, exist := answerText[code]
	if !exist {
		return "unknown error response '" + code + "'"
	}
	return text
}
