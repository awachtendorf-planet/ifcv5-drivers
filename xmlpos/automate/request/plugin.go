package request

import (
	"fmt"
	"strings"
	"time"

	xmlpos "github.com/weareplanet/ifcv5-drivers/xmlpos/dispatcher"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"

	rq "github.com/weareplanet/ifcv5-main/ifc/automate/request"

	"github.com/weareplanet/ifcv5-drivers/xmlpos/dispatcher"
	"github.com/weareplanet/ifcv5-drivers/xmlpos/template"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Plugin ...
type Plugin struct {
	*rq.Plugin
	driver *xmlpos.Dispatcher
}

// New return a new plugin
func New(parent *xmlpos.Dispatcher) *Plugin {

	p := &Plugin{
		rq.New(),
		parent,
	}

	p.Setup(rq.Config{

		Name:          fmt.Sprintf("%T", p),
		TemplateAck:   template.PacketAck,
		TemplateNak:   template.PacketNak,
		ProcessPacket: p.processPacket,
	})

	p.RegisterPacket(template.PacketPostInquiry, template.PacketPostRequest)

	return p
}

func (p *Plugin) processPacket(addr string, action *dispatcher.StateAction, packet *ifc.LogicalPacket) error {

	var err error

	switch packet.Name {

	case template.PacketPostInquiry:
		err = p.handlePostInquiry(addr, action, packet)

	case template.PacketPostRequest:
		err = p.handlePostRequest(addr, action, packet)

	default:
		err = errors.Errorf("no handler defined to process packet '%s'", packet.Name)

	}

	return err

}

func (p *Plugin) send(addr string, action *dispatcher.StateAction, tracking string, context interface{}) error {

	packet, err := p.driver.ConstructPacket(addr, template.PacketFramed, tracking, context)
	if err != nil {
		return err
	}

	err = p.SendPacket(addr, packet, action)

	return err
}

func (p *Plugin) getTime(d time.Time, t time.Time) time.Time {

	if d.IsZero() || t.IsZero() {
		return time.Now()
	}

	return t.AddDate(d.Year(), int(d.Month())-1, d.Day()-1)

}

func (p *Plugin) getRequestType(requestType int) int {

	var rt record.RequestType

	switch requestType {

	case 2: // track2
		rt = record.RequestByTrack2

	case 4: // room number
		rt = record.RequestByRoomNumber

	case 8: // guest name
		rt = record.RequestByLastName

	case 6: // track2 or room number
		rt = record.RequestByTrack2 | record.RequestByRoomNumber

	case 10: // track2 or guest name
		rt = record.RequestByTrack2 | record.RequestByLastName

	case 12: // room number or guest name
		rt = record.RequestByRoomNumber | record.RequestByLastName

	case 14: // track2 or guest name or room number
		rt = record.RequestByTrack2 | record.RequestByLastName | record.RequestByRoomNumber

	}

	return int(rt)
}

func (p *Plugin) toInt(v string) int {
	v = strings.TrimSpace(v)
	v = strings.TrimLeft(v, "0")
	return cast.ToInt(v)
}

func (p *Plugin) toString(v int) string {
	return cast.ToString(v)
}

/*
CO Posting denied because overwriting the CreditLimit is not allowed
DM Sum of Subtotals does not match TotalAmount
FX Posting denied because the guest is not allowed to use posting features
IA Invalid Account information
NA Posting Denied - Night-Audit procedure is currently running
NG Guest/Room not found
NP Posting denied because NoPost is set
OK Posted successfully
UR Unprocessable request (generic error)
*/
