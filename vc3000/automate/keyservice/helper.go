package keyservice

import (
	"fmt"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"

	"github.com/spf13/cast"
)

func (p *Plugin) reconstructKeyOptions(packet *ifc.LogicalPacket, station uint64) string {

	driver := p.driver
	dispatcher := p.automate.Dispatcher()

	var ko string
	if value, exist := driver.GetPayLoadField(packet, "A"); exist {
		ko = value
	}

	// did we force values
	if _, err := dispatcher.GetConfigIfExist(station, defines.AccessPoints); err == nil {
		return ko
	}

	// did we use special Radisson logic
	if driver.IsRadisson(station) {

		userType, _ := driver.GetPayLoadField(packet, "T")
		userGroup, _ := driver.GetPayLoadField(packet, "U")

		userType = dispatcher.GetPMSMapping("User Type", station, "T", userType)
		userGroup = dispatcher.GetPMSMapping("User Group", station, "U", userGroup)

		// user type 1.te stelle
		ko1 := cast.ToUint(userType)

		// user group 2.te stelle
		ko2 := cast.ToUint(userGroup)

		return fmt.Sprintf("%d%d", ko1, ko2)

	}

	// return A field from vendor
	return ko
}
