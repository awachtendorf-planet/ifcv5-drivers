package ilco

// var (
// 	answerText = map[string]string{
// 		"ES": "Syntax error",
// 		"NC": "The encoder does not answer. Failure in the communications, switched off, do not exist...",
// 		"NF": "Data files in the PC interface are damaged or not found",
// 		"OV": "The encoder has not already accomplished the previous task",
// 		"EP": "Error in magnetic track, card inserted wrongly or without magnetic stripe",
// 		"EF": "The card has been encoded by another system or the magnetic strip may be damaged",
// 		"EN": "The card has been encoded with a too low magnetic level due to dust in the reader magnetic head or low quality card",
// 		"ET": "Stuck card. The card does not have the required physical size",
// 		"TD": "Unknown card",
// 		"ED": "Timeout error. Operation cancelled",
// 		"EA": "Requested copy corresponds to a room checked-out",
// 		"OS": "Room out of service",
// 		"EO": "The requested card is being recorded by other station",
// 		"EE": "Task canceled",
// 	}
// )

// // GetAnswerText returns a clear text for fias
// func (d *Dispatcher) GetAnswerText(answerStatus string) string {
// 	text, exist := answerText[answerStatus]
// 	if !exist {
// 		return answerStatus
// 	}
// 	return text
// }
