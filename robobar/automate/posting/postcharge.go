package posting

import (
	"strings"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	"github.com/weareplanet/ifcv5-drivers/robobar/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type sale struct {
	Sequence    int       // posting sequence number
	Room        string    // room/cause
	Index       int       // article number
	Charge      float64   // amount
	Description string    // description
	Date        time.Time // posting date
}

func (p *Plugin) handlePostCharge(addr string, packet *ifc.LogicalPacket) error {

	automate := p.automate
	dispatcher := automate.Dispatcher()

	station, _ := dispatcher.GetStationAddr(addr)

	item, err := p.unmarshal(packet)
	if err != nil {
		return err
	}

	if len(item.Room) == 0 {
		return errors.New("empty room provided")
	}

	if item.Charge == 0 {
		return errors.New("zero charge amount")
	}

	charge := record.ChargePosting{}

	chargeItem := record.ChargeItem{
		ID:    1,
		Name:  item.Description,
		Value: item.Charge,
	}
	charge.References = append(charge.References, chargeItem)

	article := record.Article{
		Quantity: 1,
		Number:   item.Index,
	}
	charge.MinibarInfo.Items = append(charge.MinibarInfo.Items, article)

	posting := record.SimplePosting{
		Time:           item.Date,
		Station:        station,
		Context:        charge,
		PostingType:    1, // direct charge
		RoomName:       item.Room,
		SequenceNumber: cast.ToString(item.Sequence),
		TotalAmount:    item.Charge,
		OutletNumber:   item.Index, // totally confusing bullshit
	}

	dispatcher.CreatePmsJob(addr, packet, posting)

	p.driver.SetSequenceNumber(station, item.Sequence)

	return nil

}

func (p *Plugin) unmarshal(packet *ifc.LogicalPacket) (sale, error) {
	item := sale{}

	if packet == nil {
		return item, errors.New("empty logical packet")
	}

	if packet.Name != template.PacketSale {
		return item, errors.Errorf("wrong logical packet '%s', expected '%s'", packet.Name, template.PacketSale)
	}

	// * room
	// 2 index
	// 7 charge
	// 20 description
	// 10 date

	data := packet.Data()["Data"]
	if len(data) < 40 {
		return item, errors.New("logical packet payload too short")
	}

	tmp := string(packet.Data()["SequenceNumber"])
	tmp = strings.TrimLeft(tmp, " ")
	tmp = strings.TrimLeft(tmp, "0")
	item.Sequence = cast.ToInt(tmp)

	// room
	//tmp = p.driver.DecodeString(packet.Addr, data[:len(data)-39])
	tmp = string(data[:len(data)-39])
	tmp = strings.TrimLeft(tmp, " ")
	tmp = strings.TrimLeft(tmp, "0")
	item.Room = tmp

	// index
	data = data[len(data)-39:]
	if len(data) < 2 {
		return item, errors.New("logical packet field 'index' failed")
	}

	tmp = string(data[:2])
	tmp = strings.TrimLeft(tmp, " ")
	tmp = strings.TrimLeft(tmp, "0")
	item.Index = cast.ToInt(tmp)

	// charge
	data = data[2:]
	if len(data) < 7 {
		return item, errors.New("logical packet field 'charge' failed")
	}

	tmp = string(data[:7])
	tmp = strings.TrimLeft(tmp, " ")
	tmp = strings.TrimLeft(tmp, "0")
	item.Charge = cast.ToFloat64(tmp)

	// description
	data = data[7:]
	if len(data) < 20 {
		return item, errors.New("logical packet field 'description' failed")
	}

	tmp = string(data[:20])
	tmp = strings.Trim(tmp, " ")
	item.Description = tmp

	// date
	data = data[20:]
	if len(data) < 10 {
		return item, errors.New("logical packet field 'date' failed")
	}

	tmp = string(data[:10])
	if constructed, err := time.Parse("0601021504", tmp); err == nil { // YYMMDDhhmm
		item.Date = constructed
	}

	return item, nil

}
