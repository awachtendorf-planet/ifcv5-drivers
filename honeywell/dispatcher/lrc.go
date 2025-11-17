package honeywell

import (
	"strconv"
)

// ProcessOutgoingLRC callback, calculate LRC
func (d *Dispatcher) ProcessOutgoingLRC(addr string, _ string, data *[]byte) {

	station, _ := d.GetStationAddr(addr)
	protocol := d.GetProtocolType(station)

	switch protocol {

	case HONEYWELL_PROTOCOL, ALERTON_PROTOCOL_2: // no LRC
		return

	case ALERTON_PROTOCOL_1: // sum room number, send digi one

		if data == nil || len(*data) != 5 {
			return
		}

		calculated := byte(0)

		for i := 1; i < 5; i++ {
			char := (*data)[i]
			calculated += char - 0x30
		}

		lrc := strconv.Itoa(int(calculated))

		*data = append(*data, lrc[len(lrc)-1])

		return

	}

}
