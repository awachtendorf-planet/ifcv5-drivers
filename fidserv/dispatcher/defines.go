package fidserv

import (
	"fmt"
	"time"

	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/record"
)

// ...
const (
	RetryDelay         = 5 * time.Second
	PacketTimeout      = 8 * time.Second
	AliveTimeout       = 30 * time.Second
	KeyAnswerTimeout   = 60 * time.Second
	PmsTimeout         = 10 * time.Second
	PmsSyncTimeout     = 60 * time.Second // database sync
	NextActionDelay    = 0
	MaxError           = 9
	MaxKeyEncoderError = 3
	MaxNetworkError    = 1000
)

const (
	// payload object from templates
	payload = "Data"
)

func (d *Dispatcher) setDefaults() {

	d.InitScope("fias")

	// DT = Checkout Time, mapped to GD
	if value, exist := d.GetMappedGuestField("GD"); exist {
		d.SetMappedGuestField("DT", value)
	}

	d.SetMappedGenericField("WS", defines.WorkStation)
	d.SetMappedGenericField("KC", defines.EncoderNumber)
	d.SetMappedGenericField("KT", defines.KeyType)
	d.SetMappedGenericField("K#", defines.KeyCount)
	d.SetMappedGenericField("KO", defines.KeyOptions)
	d.SetMappedGenericField("$1", defines.Track1)
	d.SetMappedGenericField("$2", defines.Track2)
	d.SetMappedGenericField("$3", defines.Track3)
	d.SetMappedGenericField("$4", defines.Track4)
	d.SetMappedGenericField("ID", defines.UserID)
	d.SetMappedGenericField("RO", defines.OldRoomName)
	d.SetMappedGenericField("RS", defines.RoomStatus)
	d.SetMappedGenericField("RN", defines.WakeExtension)
	d.SetMappedGenericField("SI", defines.AdditionalRooms)
	d.SetMappedGenericField("RT", defines.RequestType)

	//d.SetMappedGenericField("SI", defines.SuiteInfo)

	d.SetMappedGenericField("DA", defines.Timestamp)
	d.SetMappedGenericField("TI", defines.Timestamp)

	d.buildPresets()
}

func (d *Dispatcher) buildPresets() {

	if !d.ExistModule(0, "udf") {

		records := []string{"GI", "GC", "KR", "KM", "PL"}

		udf := record.Mapping{
			Driver: d.Name(),
			Module: "udf",
		}

		for r := range records {
			for i := 0; i < 10; i++ {
				item := record.MappingItem{
					Key:   fmt.Sprintf("%sA%d", records[r], i),
					Value: fmt.Sprintf("UDF%d", i+1),
				}
				udf.Mapping = append(udf.Mapping, item)
			}
		}

		d.WriteModule(udf, false)
	}

	if !d.ExistModule(0, "paytvright") {

		items := []record.MappingItem{
			{Key: "0", Value: "TN"},
			{Key: "1", Value: "TX"},
			{Key: "2", Value: "TM"},
			{Key: "3", Value: "TU"},
		}

		mapping := record.Mapping{
			Driver:  d.Name(),
			Module:  "paytvright",
			Alias:   "TV",
			Default: "TU",
		}

		for i := range items {
			mapping.Mapping = append(mapping.Mapping, items[i])
		}

		d.WriteModule(mapping, false)
	}

	if !d.ExistModule(0, "minibarright") {

		items := []record.MappingItem{
			{Key: "0", Value: "ML"},
			{Key: "2", Value: "MU"},
		}

		mapping := record.Mapping{
			Driver:  d.Name(),
			Module:  "minibarright",
			Alias:   "MR",
			Default: "MN",
		}

		for i := range items {
			mapping.Mapping = append(mapping.Mapping, items[i])
		}

		d.WriteModule(mapping, false)
	}

	if !d.ExistModule(0, "videoright") {

		items := []record.MappingItem{
			{Key: "0", Value: "VN"},
			{Key: "1", Value: "VB"},
			{Key: "2", Value: "VA"},
		}

		mapping := record.Mapping{
			Driver:  d.Name(),
			Module:  "videoright",
			Alias:   "VR",
			Default: "VN",
		}

		for i := range items {
			mapping.Mapping = append(mapping.Mapping, items[i])
		}

		d.WriteModule(mapping, false)
	}
}
