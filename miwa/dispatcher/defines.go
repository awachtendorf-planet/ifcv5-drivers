package miwa

import (
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 30 * time.Second
	KeyAnswerTimeout   = 40 * time.Second
	PmsTimeout         = 10 * time.Second
	PmsSyncTimeout     = 30 * time.Second // database sync
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)

const (
	// payload object from templates
	payload = "Data"
)

func (d *Dispatcher) setDefaults() {

	d.InitScope("miwa")

	// DT = Checkout Time, mapped to GD
	if value, exist := d.GetMappedGuestField("GD"); exist {
		d.SetMappedGuestField("DT", value)
	}

	d.SetMappedGenericField("WS", defines.WorkStation)
	d.SetMappedGenericField("KC", defines.EncoderNumber)
	d.SetMappedGenericField("KT", defines.KeyType)
	d.SetMappedGenericField("K#", defines.KeyCount)
	d.SetMappedGenericField("KO", defines.KeyOptions)
	d.SetMappedGenericField("$1", defines.Track1)
	d.SetMappedGenericField("$2", defines.Track2)
	d.SetMappedGenericField("$3", defines.Track3)
	d.SetMappedGenericField("$4", defines.Track4)
	d.SetMappedGenericField("ID", defines.UserID)
	d.SetMappedGenericField("RO", defines.OldRoomName)
	d.SetMappedGenericField("RS", defines.RoomStatus)
	d.SetMappedGenericField("RN", defines.WakeExtension)
	d.SetMappedGenericField("SI", defines.AdditionalRooms)

	d.SetMappedGenericField("DA", defines.Timestamp)
	d.SetMappedGenericField("TI", defines.Timestamp)

}

var (
	errorText = map[string]string{
		"1": "In use initiated by terminal.",
		"2": "In use initiated by PC.",
		"3": "In use initiated by PMS.",
		"4": "Failure in terminal.",
		"5": "No response from terminal.",
		"6": "Failure in text (failure in text ID or terminal ID)",
		"7": "Failure in text integrity.",
	}
)

// GetErrorText returns a clear text
func (d *Dispatcher) GetErrorText(code string) string {
	text, exist := errorText[code]
	if !exist {
		return "unknown error response '" + code + "'"
	}
	return text
}

var (
	statusText = map[string]string{
		"0": "Ready",
		"1": "In operation initiated by terminal",
		"2": "In operation initiated by PC",
		"3": "In operation initiated by PMS",
		"4": "Failure",
	}
)

// GetKeyCreateResponse returns a clear text
func (d *Dispatcher) GetKeyCreateResponse(code string) string {
	text, exist := statusText[code]
	if !exist {
		return "unknown error response '" + code + "'"
	}
	return text
}

var (
	responseText = map[string]string{
		"0": "Normally completed",
		"1": "Failure in writing",
		"2": "Card being stuck",
		"3": "Time out",
		"4": "Other cards than guest or maintenance cards",
		"9": "Other failures",
	}
)

// GetResponseText returns a clear text
func (d *Dispatcher) GetKeyReadResponseText(code string) string {
	text, exist := responseText[code]
	if !exist {
		return "unknown error response '" + code + "'"
	}
	return text
}
