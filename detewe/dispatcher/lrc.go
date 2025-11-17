package detewe

import (
	"fmt"

	"github.com/spf13/cast"
)

func (d *Dispatcher) CalcLRCHighLow(data *[]byte) (byte, byte) {

	sum := 0x1000000

	for _, byt := range *data {

		sum -= cast.ToInt(byt)
	}

	letterH := fmt.Sprintf("%X", ((sum & 0xF0) >> 4))
	letterL := fmt.Sprintf("%X", (sum & 0x0F))

	return letterH[0], letterL[0]

}
