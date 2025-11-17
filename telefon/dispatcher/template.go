package telefon

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	templatestate "github.com/weareplanet/ifcv5-drivers/telefon/state"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils/dynamicstruct"

	"github.com/weareplanet/ifcv5-drivers/telefon/template"

	"github.com/pkg/errors"
)

var (
	ovrCounter = 0
)

type Config struct {
	Template struct {
		Driver   string   `json:"driver"`             // must matched 'telefon'
		Vendor   string   `json:"vendor"`             // some useful vendor description
		Disabled bool     `json:"disabled"`           // set to true to remove the layout from the parser
		TryRun   bool     `json:"tryrun"`             // set to true to test the meta compiler
		Framing  Framing  `json:"framing,omitempty"`  // optional, defined the byte stream framing
		Protocol Protocol `json:"protocol,omitempty"` // optional, defined the byte stream framing
		Layout   []Layout `json:"layout"`             //
	} `json:"template"`
}

type Layout struct {
	internal bool
	outgoing bool
	Name     string  `json:"name,omitempty"`    // auto generated, overwrite only if absolutely necessary
	Hint     string  `json:"hint,omitempty"`    // auto generated, overwrite only if absolutely necessary
	Garbage  bool    `json:"garbage,omitempty"` // auto generated, overwrite only if absolutely necessary
	Rewind   int     `json:"rewind,omitempty"`  // auto generated, overwrite only if absolutely necessary
	Field    []Field `json:"field"`             // defined the parser commands
}

type Framing struct {
	Start string `json:"start,omitempty"` // meta compiler generate garbage filter and insert 'Start' to each field
	End   string `json:"end,omitempty"`   // meta compiler generate garbage filter and append 'overread until End' to each field
}

type Protocol struct {
	Ack     string  `json:"ack,omitempty"` // defines incoming and outgoing low level packet
	Nak     string  `json:"nak,omitempty"` // defines incoming and outgoing low level packet
	Enq     string  `json:"enq,omitempty"` // defines incoming and outgoing low level packet
	LRC     LRC     `json:"lrc,omitempty"`
	Reply   Reply   `json:"reply,omitempty"`
	Polling Polling `json:"polling,omitempty"`
}

type LRC struct {
	Type   string `json:"type,omitempty"`   // defines LRC method
	Len    int    `json:"len,omitempty"`    // length of LRC
	Seed   int    `json:"seed,omitempty"`   // seed of LRC
	Inside bool   `json:"inside,omitempty"` // before or after framing
}

type Reply struct {
	Enq string `json:"enq,omitempty"` // reply with enq byte or sequence
}

type Polling struct {
	Char     string `json:"char,omitempty"`     // send char byte or sequence as polling packet
	Interval int    `json:"interval,omitempty"` // polling interval in seconds, 0 = disable polling (internal min 3)
}

type Field struct {
	Name   string `json:"name,omitempty"`   // auto generated, overwrite only if absolutely necessary
	Type   string `json:"type,omitempty"`   // auto generated, overwrite only if absolutely necessary (int, byte, []byte)
	Len    int    `json:"len,omitempty"`    // 0 = overread, auto generated for 'Equal' value
	Equal  string `json:"equal,omitempty"`  // match a specifically byte or sequence of bytes (hexadecimal representation)
	Endian string `json:"endian,omitempty"` // defined the int representation (little, big)

	Overread     bool   `json:"overread,omitempty"`     // ! describes an unneeded field, can be combined with 'Len'
	Extension    bool   `json:"extension,omitempty"`    // ! extractable value
	DialedNumber bool   `json:"dialednumber,omitempty"` // ! extractable value
	Duration     bool   `json:"duration,omitempty"`     // ! extractable value, format eg "hh:mm:ss" or "mmmss"
	Units        bool   `json:"units,omitempty"`        // ! extractable value
	Amount       bool   `json:"amount,omitempty"`       // ! extractable value, optional format eg "2" (decimals)
	CallDate     bool   `json:"calldate,omitempty"`     // ! extractable value, format eg "2006/01/02" or "2006/01/02 15:04:05"
	CallTime     bool   `json:"calltime,omitempty"`     // ! extractable value, optional if 'CallDate' only match the day, format eg "15:04:05" or "3:04pm"
	CallType     bool   `json:"calltype,omitempty"`     // ! extractable value
	User         bool   `json:"user,omitempty"`         // ! extractable value
	Outlet       bool   `json:"outlet,omitempty"`       // ! extractable value
	RoomStatus   bool   `json:"roomstatus,omitempty"`   // ! extractable value
	Article      bool   `json:"article,omitempty"`      // ! extractable value
	Quantity     bool   `json:"quantity,omitempty"`     // ! extractable value
	Format       string `json:"format,omitempty"`       // ! available for CallDate, CallTime, Duration, Amount

}

func (f *Field) String() string {
	return fmt.Sprintf("name: %s, equal: %s, len: %d", f.Name, f.Equal, f.Len)
}

func (d *Dispatcher) handleTemplate(file string, data *[]byte) bool {

	if data == nil || d.parser == nil {
		return false
	}

	var config Config

	// TODO: check if file is json or yaml etc ...

	if err := json.Unmarshal(*data, &config); err != nil {
		return false
	}

	if !strings.EqualFold(config.Template.Driver, d.Name()) {
		if len(config.Template.Driver) > 0 {
			log.Warn().Msgf("%T skip config file '%s', mismatched driver name '%s' ", d, file, config.Template.Driver)
		}
		return false
	}

	if len(config.Template.Vendor) == 0 {
		log.Error().Msgf("%T template field 'vendor' is empty, ignore template", d)
		return true
	}

	tryRun := config.Template.TryRun

	slot := d.getSlot(config.Template.Vendor)
	if slot == 0 && !tryRun {
		slot = d.newSlot(config.Template.Vendor)
	}

	if tryRun {
		log.Info().Str("vendor", config.Template.Vendor).Msgf("%T try config file '%s'", d, file)
	} else {
		log.Info().Str("vendor", config.Template.Vendor).Uint("slot", slot).Msgf("%T process config file '%s'", d, file)
	}

	err := d.constructTemplate(slot, config, tryRun)
	if err != nil {
		log.Error().Msgf("%T %s", d, err)
	}

	return true

}

func (d *Dispatcher) constructTemplate(slot uint, config Config, tryRun bool) (err error) {

	if !tryRun {
		d.clearProtocol(slot)
		d.parser.Lock()
		d.parser.ClearIncomingTemplates(slot)
		d.parser.ClearOutgoingTemplates(slot)
	}

	defer func() {
		if !tryRun {

			d.parser.Ready()
			d.parser.Unlock()

			broker := d.Broker()
			if broker != nil {
				broker.Broadcast(templatestate.NewEvent(slot, templatestate.Changed), templatestate.Changed.String())
			}

		}
	}()

	if config.Template.Disabled {
		//return errors.Errorf("template '%s' is disabled", config.Template.Vendor)
		return nil
	}

	// framing
	// construct incoming garbage filter

	framing := config.Template.Framing
	if len(framing.End) > 0 && len(framing.Start) > 0 {

		stx := Field{Equal: framing.Start}
		ovr := Field{Overread: true}
		etx := Field{Equal: framing.End}

		layout := Layout{Name: template.UnknownPacket, Garbage: false, internal: true}

		layout.Field = append(layout.Field, stx, ovr, etx)
		config.Template.Layout = append(config.Template.Layout, layout)
	}

	if len(framing.Start) > 0 {

		ovr := Field{Overread: true}
		etx := Field{Equal: framing.Start}

		layout := Layout{Name: template.Garbage, Hint: "overread until framing start", Garbage: true, internal: true, Rewind: len(framing.Start) / 2}

		layout.Field = append(layout.Field, ovr, etx)
		config.Template.Layout = append(config.Template.Layout, layout)
	}

	if len(framing.End) > 0 {

		ovr := Field{Overread: true}
		etx := Field{Equal: framing.End}

		layout := Layout{Name: template.Garbage, Hint: "overread until framing end", Garbage: true, internal: true}

		layout.Field = append(layout.Field, ovr, etx)
		config.Template.Layout = append(config.Template.Layout, layout)
	}

	// protocol

	protocol := config.Template.Protocol

	// LRC pre-checks

	if len(protocol.LRC.Type) > 0 {
		if handler := d.getLRCHandler(protocol.LRC.Type); handler == nil {
			return errors.Errorf("template '%s' LRC type '%s' not supported", config.Template.Vendor, protocol.LRC.Type)
		}
		if protocol.LRC.Len < 1 {
			return errors.Errorf("template '%s' LRC type '%s' defined, but no length specified", config.Template.Vendor, protocol.LRC.Type)
		}
		if len(framing.Start) > 2 {
			return errors.Errorf("template '%s' LRC type '%s' defined, but currently a framing start character with more than one byte is not supported", config.Template.Vendor, protocol.LRC.Type)
		}
		if len(framing.End) > 2 {
			return errors.Errorf("template '%s' LRC type '%s' defined, but currently a framing end character with more than one byte is not supported", config.Template.Vendor, protocol.LRC.Type)
		}
	}

	// construct outgoing/incoming low level pakets

	if len(protocol.Ack) > 0 {

		layout := d.constructLowLevel(template.Ack, protocol.Ack)

		layout.outgoing = false
		config.Template.Layout = append(config.Template.Layout, layout)

		layout.outgoing = true
		config.Template.Layout = append(config.Template.Layout, layout)

	}

	if len(protocol.Nak) > 0 {

		layout := d.constructLowLevel(template.Nak, protocol.Nak)

		layout.outgoing = false
		config.Template.Layout = append(config.Template.Layout, layout)

		layout.outgoing = true
		config.Template.Layout = append(config.Template.Layout, layout)

	}

	if len(protocol.Enq) > 0 {

		layout := d.constructLowLevel(template.Enq, protocol.Enq)

		layout.outgoing = false
		config.Template.Layout = append(config.Template.Layout, layout)

		layout.outgoing = true
		config.Template.Layout = append(config.Template.Layout, layout)

	}

	if len(protocol.Reply.Enq) > 0 {
		layout := d.constructLowLevel(template.EnqReply, protocol.Reply.Enq)

		layout.outgoing = true
		config.Template.Layout = append(config.Template.Layout, layout)
	}

	if len(protocol.Polling.Char) > 0 && protocol.Polling.Interval > 0 {
		layout := d.constructLowLevel(template.Polling, protocol.Polling.Char)

		layout.outgoing = true
		config.Template.Layout = append(config.Template.Layout, layout)
	}

	// construct layouts

	for i := range config.Template.Layout {

		layout := config.Template.Layout[i]

		if len(layout.Name) == 0 {
			if layout.Garbage {
				layout.Name = template.Garbage
			} else {
				layout.Name = template.Undefined
			}
		}

		if len(layout.Hint) == 0 {
			layout.Hint = fmt.Sprintf("%s - template: %d", config.Template.Vendor, i+1)
		}

		ovrCounter = 0
		failed := false
		lastLenUndefined := false
		seen := make(map[string]bool)

		tpl := dynamicstruct.NewStruct()

		if !layout.internal {
			if len(framing.Start) > 0 {
				stx := Field{
					Equal: framing.Start,
				}
				layout.Field = append(layout.Field, Field{}) // make sure slice has enough capacity
				copy(layout.Field[1:], layout.Field[0:])     // shift the slice
				layout.Field[0] = stx                        // insert at front
			}
			if len(framing.End) > 0 {
				// ovr := Field{
				// 	Overread: true,
				// }
				etx := Field{
					Equal: framing.End,
				}
				//layout.Field = append(layout.Field, ovr, etx)
				layout.Field = append(layout.Field, etx)
			}
		}

		for x := range layout.Field {

			field := layout.Field[x]

			field.Name = strings.Replace(field.Name, " ", "", -1)

			for i := range field.Name {
				char := field.Name[i]
				if (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || (char >= 48 && char <= 57) || char == 95 { // A-Z, a-z, 0-9, _
					continue
				}
				log.Warn().
					Interface("layout", layout).
					Interface("field", field).
					Msgf("%T replace illegal character '%s' in field name '%s'", d, string(char), field.Name)

				field.Name = strings.Replace(field.Name, string(char), "_", -1) // looks like a dirty hack
			}

			meta, e := d.constructFieldData(&field)
			if e != nil {
				err = errors.Errorf("layout '%s' construct field '%s' failed, err=%s", layout.Name, field.String(), e)
				failed = true
				break
			}

			if seen[field.Name] {
				err = errors.Errorf("layout '%s' field name '%s' already used", layout.Name, field.Name)
				failed = true
				break
			}

			seen[field.Name] = true

			if lastLenUndefined && field.Len == 0 {
				err = errors.Errorf("layout '%s' construct field '%s' failed, err=%s", layout.Name, field.String(), "consecutive variable lengths are not allowed")
				failed = true
				break
			}

			if lastLenUndefined && len(field.Equal) == 0 {
				err = errors.Errorf("layout '%s' construct field '%s' failed, err=%s", layout.Name, field.String(), "variable length must be followed by a fixed value")
				failed = true
				break
			}

			lastLenUndefined = field.Len == 0
			extractable := !strings.HasSuffix(field.Name, "_")

			// fmt.Printf("name: %s, type: %s, meta: %s, extractable: %t \n", field.Name, field.Type, meta, extractable)

			fieldType := interface{}(nil)

			switch field.Type {

			case "int":
				fieldType = 0

			case "byte":
				fieldType = byte(0)

			case "[]byte":
				fieldType = []byte{}

			}

			// construct bytesparser struct
			if extractable && len(field.Format) > 0 {
				tpl = tpl.AddField(field.Name, fieldType, `byte:"`+meta+`" format:"`+field.Format+`"`)
			} else {
				tpl = tpl.AddField(field.Name, fieldType, `byte:"`+meta+`"`)
			}

			// detect packet type
			if layout.Name == template.Undefined {
				if field.DialedNumber || field.Duration || field.Units || field.CallDate || field.CallTime {
					layout.Name = template.CallPacket
				} else if field.RoomStatus {
					layout.Name = template.RoomStatus
				} else if field.Article {
					layout.Name = template.Posting
				}
			}

		}

		if failed {
			return
		}

		defer func() {
			if e := recover(); e != nil {

				err = errors.Errorf("layout '%s' construct template failed, err=%s", layout.Name, fmt.Sprint(e))

				if !tryRun {
					d.parser.ClearIncomingTemplates(slot)
					d.parser.ClearOutgoingTemplates(slot)
				}

			}
		}()

		if !tryRun {

			template := tpl.Build().New()

			// fmt.Printf("%T \n", template)

			if layout.outgoing {
				d.parser.RegisterOutgoingTemplate(slot, layout.Name, template)
			} else {
				d.parser.RegisterIncomingTemplate(slot, layout.Name, layout.Hint, template, layout.Garbage, layout.Rewind)
			}

		}

	}

	if !tryRun {
		d.setProtocol(slot, protocol)
	}

	return
}

func (d *Dispatcher) constructLowLevel(name string, equal string) Layout {

	layout := Layout{}

	if len(name) > 0 && len(equal) > 0 {

		pkt := Field{
			Name:  name + "_",
			Equal: equal,
		}

		layout.Name = name
		layout.internal = true

		layout.Field = append(layout.Field, pkt)

	}

	return layout
}

func (d *Dispatcher) constructFieldData(field *Field) (string, error) {

	var meta strings.Builder

	field.Name = strings.Trim(field.Name, " ")
	field.Type = strings.Trim(field.Type, " ")
	field.Equal = strings.Trim(field.Equal, " ")
	field.Endian = strings.Trim(field.Endian, " ")

	if strings.HasSuffix(field.Name, "_") {
		field.Name = strings.ToUpper(field.Name)
	} else {
		field.Name = strings.Title(field.Name)
	}

	if len(field.Equal) > 0 {
		if (len(field.Equal) % 2) == 1 {
			return "", errors.Errorf("field 'equal' has wrong length: %d, must be a multiple of 2", len(field.Equal))
		}
		if field.Len == 0 {
			field.Len = len(field.Equal) / 2
		} else if field.Len != len(field.Equal)/2 {
			return "", errors.Errorf("field 'len' has wrong length: %d, should be: %d", field.Len, len(field.Equal)/2)
		}

		if len(field.Name) == 0 {
			ovrCounter++
			field.Name = fmt.Sprintf("OVR_AUTOGEN_%d_", ovrCounter)
		}
	}

	if field.Overread {
		if len(field.Name) == 0 {
			ovrCounter++
			field.Name = fmt.Sprintf("OVR_AUTOGEN_%d_", ovrCounter)
		}
	}

	meta.WriteString("len:")
	if field.Len == 0 {
		meta.WriteString("*")
	} else {
		meta.WriteString(fmt.Sprintf("%d", field.Len))
	}

	if len(field.Equal) > 0 {

		// bytesparser pre check
		src := []byte(field.Equal)
		dst := make([]byte, hex.DecodedLen(len(src)))
		if _, err := hex.Decode(dst, src); err != nil {
			return "", err
		}

		meta.WriteString(",equal:0x" + field.Equal)
	}

	if len(field.Endian) > 0 {
		field.Endian = strings.ToLower(field.Endian)
		if field.Endian == "little" || field.Endian == "big" {
			meta.WriteString(",endian:" + field.Endian)
		} else {
			return "", errors.Errorf("field 'endian' has wrong value '%s', possible values 'little' or 'big' ", field.Endian)
		}
		field.Type = "int"
	}

	field.Type = strings.ToLower(field.Type)
	if len(field.Type) == 0 {
		if field.Len == 1 {
			field.Type = "int"
		} else {
			field.Type = "[]byte"
		}
	}

	switch field.Type {

	case "int", "byte", "[]byte":

	default:
		return "", errors.Errorf("field 'type' has wrong value '%s'", field.Type)

	}

	if field.Extension {
		field.Name = template.Extension
	} else if field.DialedNumber {
		field.Name = template.DialedNumber
	} else if field.Duration {
		field.Name = template.Duration
	} else if field.Units {
		field.Name = template.Units
	} else if field.Amount {
		field.Name = template.Amount
	} else if field.CallDate {
		field.Name = template.CallDate
	} else if field.CallTime {
		field.Name = template.CallTime
	} else if field.CallType {
		field.Name = template.CallType
	} else if field.RoomStatus {
		field.Name = template.RoomStatus
	}

	if len(field.Name) == 0 {
		return "", errors.New("field 'name' is empty")
	}

	return meta.String(), nil

}
