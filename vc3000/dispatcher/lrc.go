package vc3000

import (
	"github.com/pkg/errors"
	"github.com/weareplanet/ifcv5-drivers/vc3000/template"
	"github.com/weareplanet/ifcv5-main/log"
)

// GetAdditionalLRCSize return 2 if LRC is expected
func (d *Dispatcher) GetAdditionalLRCSize(addr string, name string) int {
	if d.IsAcknowledgement(addr, name) {
		switch name {
		case template.PacketAck, template.PacketNak, template.PacketEnq:
			return 0
		default:
			return 2
		}
	}
	return 0
}

// ProcessIncomingLRC callback, check LRC
func (d *Dispatcher) ProcessIncomingLRC(addr string, data []byte, dataLength int, LRC []byte, _ string) error {

	var err error
	if len(LRC) < 1 {
		err = errors.New("crc failed, no lrc provided")
	} else if len(LRC) != 2 {
		err = errors.New("crc failed, lrc size expected 2 bytes")
	} else if len(data) < 3 {
		err = errors.New("crc failed, no input data provided")
	}

	if err == nil {

		calculated := byte(0)
		for _, char := range data {
			switch char {

			case 0x02:
				calculated = 0

			case 0x03:
				calculated ^= char

				hb := d.getHighByte(calculated)
				lb := d.getLowByte(calculated)

				if hb == LRC[0] && lb == LRC[1] {
					return nil
				}

				err = errors.Errorf("crc mismatch, received 0x%x 0x%x, calculated: 0x%x 0x%x", LRC[0], LRC[1], hb, lb)

				if station, e := d.GetStationAddr(addr); e == nil {
					if d.IsDebugMode(station) {
						log.Warn().Msgf("debug mode: %s", err.Error())
						return nil
					}
				}

				break

			default:
				calculated ^= char
			}
		}

	}

	d.sendRefusion(addr)
	if err == nil {
		err = errors.New("crc mismatch, packet framing failed")
	}

	return err
}

// ProcessOutgoingLRC callback, calculate LRC
func (d *Dispatcher) ProcessOutgoingLRC(addr string, name string, data *[]byte) {

	if len(*data) < 3 { // quick and dirty, no framing packet (ack, nack, enq)
		return
	}

	if d.GetAdditionalLRCSize(addr, name) == 0 {
		return
	}

	calculated := byte(0)
	for _, char := range *data {
		switch char {

		case 0x02:
			calculated = 0

		case 0x03:
			calculated ^= char

			char = d.getHighByte(calculated)
			*data = append(*data, char)

			char = d.getLowByte(calculated)
			*data = append(*data, char)

			return

		default:
			calculated ^= char
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
