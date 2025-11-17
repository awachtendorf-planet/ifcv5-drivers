package fidserv

import (
	"github.com/weareplanet/ifcv5-drivers/fidserv/template"
)

func (d *Dispatcher) filterField(packetName string, key string) bool {
	if key == "SF" {
		return true
	}

	// filter KZ
	if packetName == template.PacketKeyRead && !(key == "WS" || key == "KC" || key == "DA" || key == "TI") {
		return true
	}

	// filter RE
	if packetName == template.PacketRoomData && !(key == "RN" || key == "CS" || key == "DN" || key == "G#" || key == "ML" || key == "MR" || key == "TV" || key == "RS" || key == "DA" || key == "TI") {
		// RS is not spec conform, but we support RoomStatus from PMS
		return true
	}

	return false
}
