package ahl

import (
	"github.com/weareplanet/ifcv5-drivers/ahl/template"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
)

// ProcessIncomingLRC callback, check LRC
func (d *Dispatcher) ProcessIncomingLRC(addr string, data []byte, dataLength int, _ []byte, name string) error {

	if dataLength < 3 || name == template.PacketLinkStart || name == template.PacketLinkAlive || name == template.PacketLogOutput { // quick and dirty, no framing packet (ack, nak, enq)
		return nil
	}

	calculated := uint(0)
	for i := 0; i < dataLength-3; i++ {
		char := data[i]
		switch char {

		case 0x02:
			calculated = 0

		default:
			calculated ^= uint(char)

		}
	}

	hb := d.getHighByte(byte(calculated))
	lb := d.getLowByte(byte(calculated))

	if data[dataLength-3] == hb && data[dataLength-2] == lb {
		return nil
	}

	err := errors.Errorf("crc mismatch, received 0x%x 0x%x, calculated: 0x%x 0x%x", data[dataLength-3], data[dataLength-2], hb, lb)

	if station, e := d.GetStationAddr(addr); e == nil {
		if d.IsDebugMode(station) {
			log.Warn().Msgf("debug mode: %s", err.Error())
			return nil
		}
	}

	d.sendRefusion(addr)
	return err
}

// ProcessOutgoingLRC callback, calculate LRC
func (d *Dispatcher) ProcessOutgoingLRC(addr string, name string, data *[]byte) {

	dataLength := len(*data)

	if dataLength < 3 || name == template.PacketLinkStart || name == template.PacketLinkAlive { // quick and dirty, no framing packet (ack, nak, enq)
		return
	}

	calculated := uint(0)

	for i := 0; i < dataLength-3; i++ {
		char := (*data)[i]
		switch char {

		case 0x02:
			calculated = 0

		default:
			calculated ^= uint(char)

		}
	}

	hb := d.getHighByte(byte(calculated))
	lb := d.getLowByte(byte(calculated))

	// if (*data)[dataLength-3] != '?' || (*data)[dataLength-2] != '?' {
	// 	fmt.Println("arg, something wrong with the LRC offset")
	// }

	(*data)[dataLength-3] = hb
	(*data)[dataLength-2] = lb

}

func (d *Dispatcher) getHighByte(data byte) byte {
	b := data >> 4
	if b > 9 {
		b = b + byte('A'-10)
	} else {
		b = b + byte('0')
	}
	return b
}

func (d *Dispatcher) getLowByte(data byte) byte {
	b := data & 0x0f
	if b > 9 {
		b = b + byte('A'-10)
	} else {
		b = b + byte('0')
	}
	return b
}
