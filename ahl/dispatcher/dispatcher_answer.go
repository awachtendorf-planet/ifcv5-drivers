package ahl

import (
	"fmt"
)

var (
	answerText = map[byte]string{
		'G': "Invalid GPIN or room extension (out of range)",
		'A': "Unavailable GPIN or room extension (already checked in for Check-in or already checked out)",
		'R': "Invalid Room Extension Number (out of range)",
		'U': "Unavailable Room Extension",
		'O': "Out of order Room Extension (to be cleaned ...)",
		'P': "Invalid or unavailable Guest Password Empty non DID GPIN File or DID GPIN File or both or Execute without attribution or condition (TCP/IP)",
		'F': "Forced non DID GPIN or DID GPIN",
		'V': "Out of order Voice Mail or link",
		'W': "Working Mail Box",
		'C': "Cancel forwarding for change from mono-occupation to multi-occupation",
		'B': "Forwarding on Paging or Voice Mail Box",
		'M': "Not consulted Message (in mail box or at message desk)",
		'L': "Busy Line with outgoing call",
		'X': "Exceeding Voice Mail Capacity",
		'N': "Name too short",
		'Q': "Unreachable database",
		'D': "Unreachable destination node",
		'H': "Handled correctly",
		'b': "Execute and attribute GPIN or password",
		'Y': "Out of facility range (mail box number ...)",
		'Z': "Already existing facility (mail box ...)",
		'I': "Not attributed facility (mail box ...)",
		'J': "Non available feature",
		'K': "Wrong message",
		'T': "Other",
		' ': "no error code reported",
	}
)

// GetAnswerText returns a clear text
func (d *Dispatcher) GetAnswerText(answerStatus byte) string {
	text, exist := answerText[answerStatus]
	if !exist {
		return fmt.Sprintf("unknown error code '%c'", answerStatus)
	}
	return text
}
