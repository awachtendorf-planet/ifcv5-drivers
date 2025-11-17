package dummy

import (
	"github.com/weareplanet/ifcv5-drivers/dummy/template"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
)

// GetAdditionalLRCSize return 1 if LRC is expected
func (d *Dispatcher) GetAdditionalLRCSize(addr string, name string) int {

	switch name {
	case template.PacketAck, template.PacketNak, template.PacketEnq:
		return 0
	}

	if d.IsAcknowledgement(addr, name) {
		return 1
	}

	return 0
}

// ProcessIncomingLRC callback, check LRC
func (d *Dispatcher) ProcessIncomingLRC(addr string, data []byte, dataLength int, LRC []byte, _ string) error {
	var err error
	if len(LRC) < 1 {
		err = errors.New("crc failed, no lrc provided")
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
				if calculated == LRC[0] {
					return nil
				}
				err = errors.Errorf("crc mismatch, received 0x%x, calculated: 0x%x", LRC[0], calculated)

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
			*data = append(*data, calculated)
			return
		default:
			calculated ^= char
		}
	}
}
