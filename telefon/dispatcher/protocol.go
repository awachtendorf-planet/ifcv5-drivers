package telefon

import (
	"github.com/weareplanet/ifcv5-drivers/telefon/template"

	"github.com/spf13/cast"
)

func (d *Dispatcher) setProtocol(slot uint, protocol Protocol) {
	d.protocolGuard.Lock()
	d.protocol[slot] = protocol
	d.protocolGuard.Unlock()
}

func (d *Dispatcher) getProtocol(slot uint) Protocol {
	d.protocolGuard.RLock()
	protocol := d.protocol[slot]
	d.protocolGuard.RUnlock()
	return protocol
}

func (d *Dispatcher) clearProtocol(slot uint) {
	d.protocolGuard.Lock()
	delete(d.protocol, slot)
	d.protocolGuard.Unlock()
}

func (d *Dispatcher) NeedPolling(addr string) bool {

	data := d.getStationSetting(addr, "IgnorePolling", "false")
	value := cast.ToBool(data)
	if value {
		return false
	}

	slot := d.GetParserSlot(addr)
	protocol := d.getProtocol(slot)
	return len(protocol.Polling.Char) > 0 && protocol.Polling.Interval > 0
}

func (d *Dispatcher) GetPollInterval(addr string) int {
	slot := d.GetParserSlot(addr)
	protocol := d.getProtocol(slot)
	interval := protocol.Polling.Interval
	return interval
}

func (d *Dispatcher) sendENQReply(addr string, name string, slot uint) bool {

	data := d.getStationSetting(addr, "IgnoreENQ", "false")
	value := cast.ToBool(data)
	if value {
		return false
	}

	protocol := d.getProtocol(slot)
	return len(protocol.Reply.Enq) > 0
}

func (d *Dispatcher) needLRCCheck(addr string, name string, slot uint) bool {

	switch name {

	case template.Ack, template.Nak, template.Enq:
		return false

	}

	data := d.getStationSetting(addr, "IgnoreLRC", "false")
	value := cast.ToBool(data)
	if value {
		return false
	}

	protocol := d.getProtocol(slot)
	return protocol.LRC.Len > 0
}

func (d *Dispatcher) getLRCSizeAfterFraming(addr string, name string, slot uint) int {

	if !d.needLRCCheck(addr, name, slot) {
		return 0
	}

	if d.IsAcknowledgement(addr, name) {
		if protocol := d.getProtocol(slot); !protocol.LRC.Inside {
			return protocol.LRC.Len
		}
	}

	return 0
}

func (d *Dispatcher) getLRCSizeInsideFraming(addr string, name string, slot uint) int {

	if !d.needLRCCheck(addr, name, slot) {
		return 0
	}

	if protocol := d.getProtocol(slot); protocol.LRC.Inside {
		return protocol.LRC.Len
	}

	return 0

}

func (d *Dispatcher) getLRCMethode(addr string, name string, slot uint) string {
	protocol := d.getProtocol(slot)
	return protocol.LRC.Type
}

func (d *Dispatcher) getLRCSeed(addr string, name string, slot uint) byte {
	protocol := d.getProtocol(slot)
	return byte(protocol.LRC.Seed)
}
