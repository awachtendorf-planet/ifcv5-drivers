package fidserv

var (
	answerText = map[string]string{
		"AA": "virtual number already assigned",
		"AN": "virtual number not found",
		"BM": "balance mismatch",
		"BY": "device busy",
		"CD": "check-out date is not today",
		"CO": "posting denied because overwriting the credit limit is not allowed",
		"DE": "has been deleted",
		"DM": "sum of subtotals does not match total amount",
		"DN": "request denied",
		"FX": "guest not allowed this feature",
		"IA": "invalid account",
		"NA": "night audit",
		"NF": "feature not enabled or check-out not running",
		"NG": "guest not found",
		"NM": "message or locator not found",
		"NP": "posting denied for this guest (no post flag)",
		"NR": "no response",
		"OK": "command or request completed successfully",
		"RF": "referral",
		"RY": "retry",
		"SV": "wakeup has been sent to external system",
		"UR": "unprocessable request",
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
