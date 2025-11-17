package visionline

import (
	"github.com/pkg/errors"
	"github.com/weareplanet/ifcv5-main/log"
)

// ProcessIncomingLRC callback, check LRC
func (d *Dispatcher) ProcessIncomingLRC(addr string, data []byte, dataLength int, _ []byte, _ string) error {

	//dataLen := len(data)

	if dataLength < 3 { // quick and dirty, no framing packet (ack, nak, enq)
		return nil
	}

	if !d.IsAcknowledgement(addr, "") { // check for serial layer
		return nil
	}

	calculated := uint(0)
	for i := 0; i < dataLength-3; i++ {
		char := data[i]
		switch char {

		case 0x02:
			calculated = 0x02

		default:
			calculated += uint(char)

		}
	}

	calculated %= 256
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

	if len(*data) < 3 { // quick and dirty, no framing packet (ack, nak, enq)
		return
	}

	if !d.IsAcknowledgement(addr, name) { // check for serial layer
		return
	}

	calculated := uint(0)
	for _, char := range *data {
		switch char {

		case 0x02:
			calculated = 0x02

		case 0x03:
			calculated %= 256
			hb := d.getHighByte(byte(calculated))
			lb := d.getLowByte(byte(calculated))

			// quick and dirty, inject LRC into outgoing byte stream
			dataPtr := *data
			*data = append(dataPtr[:len(dataPtr)-1], hb) // first LRC byte
			*data = append(*data, lb)                    // second LRC byte
			*data = append(*data, 0x03)                  // append ETX, overwritten by first LRC Byte
			return

		default:
			calculated += uint(char)
		}
	}

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
