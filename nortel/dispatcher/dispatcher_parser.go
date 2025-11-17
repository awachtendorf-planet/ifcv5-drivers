package nortel

import (
	"bytes"
)

func (d *Dispatcher) NeedMoreData(_ string, data []byte) bool {

	// quick and  dirty

	incomplete := 1

	// the vendor sends data packets beginning with 0d0a (wtf)
	// we recognise the beginning as garbage, but we must increase the 0d0a counter
	// otherwise we would not be able to correctly identify possible incomplete packages, because the parser dont recognise to wait for more data

	if len(data) > 2 && data[0] == 0x0d && data[1] == 0x0a {
		incomplete++
	}

	count := bytes.Count(data, []byte{0x0d, 0x0a})
	return count == incomplete

}
