package simphony

import (
	"fmt"
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	order "github.com/weareplanet/ifcv5-main/ifc/job"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

var (
	interestedIn = map[string][]string{}
)

const (
	ifcType = "MSC"
)

func (d *Dispatcher) initOverwrite() {

	// always use ack because of protocol design
	d.Acknowledgement = func(addr string, packetName string) bool { return d.IsSerialDevice(addr) }

	// register async response packets
	// d.ConfigureESBHandler = d.configureESBHandler

	// pre-check pms -> driver logic
	d.PreCheckDriverJob = d.PreCheck

	// subscribe if station is ready
	d.LoginStation = d.loginStation

}

func (d *Dispatcher) PreCheck(job *order.Job) error {
	if job == nil {
		return nil
	}
	return nil
}

func (d *Dispatcher) loginStation(station uint64) {

	if token := d.GetConfig(station, defines.ESBToken, ""); len(token) == 0 {
		return
	}

	// TODO: remove in production
	if d.IsDebugMode(station) {
		return
	}

	addr := d.GetVendorAddr(station, 0)

	ifcType := d.GetConfig(station, defines.IFCType, ifcType)

	subscribe := record.Subscribe{
		Station: station,
		IFCType: ifcType,
	}

	subscribe.MessageName = make(map[string][]string)

	for k, v := range interestedIn {
		subscribe.MessageName[k] = v
	}

	d.CreatePmsJob(addr, nil, subscribe)
}

func (d *Dispatcher) CalcChecksum(addr string, in *ifc.LogicalPacket) []byte {

	// Checksum
	if d.IsSerialDevice(addr) {

		dataSet := in.Data()
		byts := []byte{}

		// switch in.Name {
		// case template.PacketCIReqM:
		// 	byts = append(byts, []byte(" 1")...)
		// case template.PacketOCPReqM:
		// 	byts = append(byts, []byte(" 2")...)
		// }

		for name, data := range dataSet {
			if name != "SOH_" && name != "Checksum" && name != "EOT_" {
				byts = append(byts, data...)
			}
		}

		var sum uint16
		for _, b := range byts {
			sum = sum + uint16(b)
		}

		return []byte(fmt.Sprintf("%04X", sum))
	}
	return []byte{}

}

func (d *Dispatcher) CheckChecksum(in *ifc.LogicalPacket) bool {

	// Checksum

	dataSet := in.Data()
	olChecksum := dataSet["Checksum"]

	sum := d.CalcChecksum(in.Addr, in)

	return strings.ToUpper(string(olChecksum)) == string(sum)
}
