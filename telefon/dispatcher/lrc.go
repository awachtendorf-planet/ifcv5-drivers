package telefon

import (
	"bytes"

	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
)

// GetAdditionalLRCSize ...
func (d *Dispatcher) GetAdditionalLRCSize(addr string, name string, slot uint) int {
	return d.getLRCSizeAfterFraming(addr, name, slot)
}

// ProcessIncomingLRC callback, check LRC
func (d *Dispatcher) ProcessIncomingLRC(addr string, data []byte, dataLength int, LRC []byte, name string, slot uint) error {

	if !d.needLRCCheck(addr, name, slot) {
		return nil
	}

	var err error

	if len(data) < 3 {
		err = errors.New("crc failed, no input data provided")
	}

	// currently works only if STX and ETX are a single byte and not a sequence of bytes

	if err == nil {

		method := d.getLRCMethode(addr, name, slot)

		if len(method) == 0 {
			return nil
		}

		var payload []byte
		var payloadLength int

		seed := d.getLRCSeed(addr, name, slot)

		if LRC == nil || len(LRC) < 1 {

			// LRC inside framing
			// data = STX ... LRC ETX

			checkLength := d.getLRCSizeInsideFraming(addr, name, slot)
			if checkLength == 0 {
				return nil
			}

			payload = data[1:]                           // len of STX
			payloadLength = dataLength - 2 - checkLength // -2 = len of STX and len of ETX

			if payloadLength > 0 {
				if LRC == nil {
					LRC = []byte{}
				}
				LRC = append(LRC, data[payloadLength+1:payloadLength+1+checkLength]...)
			}

		} else {

			// LRC after framing
			// data = STX ... ETX LRC

			payload = data[1:]                        // len of STX
			payloadLength = dataLength - 1 - len(LRC) // -1 = len of STX

		}

		if handler := d.getLRCHandler(method); handler != nil {

			var calculated []byte

			if calculated, err = handler(seed, payload, payloadLength); err == nil {
				if LRC != nil && bytes.Equal(calculated, LRC) {
					return nil
				}
				err = errors.Errorf("crc mismatch, received 0x%x, calculated: 0x%x", LRC, calculated)
			}

		} else {
			err = errors.Errorf("lrc handler '%s' not found", method)
		}

	}

	if station, e := d.GetStationAddr(addr); e == nil {
		if d.IsDebugMode(station) {
			log.Warn().Msgf("debug mode: %s", err.Error())
			return nil
		}
	}

	d.sendRefusion(addr, slot)

	if err == nil {
		err = errors.New("crc mismatch, lrc handler failed")
	}

	return err
}
