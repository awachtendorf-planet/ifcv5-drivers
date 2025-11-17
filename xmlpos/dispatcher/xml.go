package xmlpos

import (
	"bytes"
	"fmt"

	"encoding/xml"

	vendor "github.com/weareplanet/ifcv5-drivers/xmlpos/record"
	ifc "github.com/weareplanet/ifcv5-main/ifc/base"

	"github.com/weareplanet/ifcv5-drivers/xmlpos/template"

	"github.com/pkg/errors"
	"github.com/pschou/go-xmltree"
)

func (d *Dispatcher) initHandler() {

	d.handler[template.PacketLinkDescription] = func() interface{} { return &vendor.LinkDescription{} }
	d.handler[template.PacketLinkAlive] = func() interface{} { return &vendor.LinkAlive{} }
	d.handler[template.PacketLinkStart] = func() interface{} { return &vendor.LinkStart{} }

	d.handler[template.PacketPostInquiry] = func() interface{} { return &vendor.PostInquiry{} }
	d.handler[template.PacketPostRequest] = func() interface{} { return &vendor.PostRequest{} }

}

// unmarshal decode the incoming data
func (d *Dispatcher) unmarshal(in *ifc.LogicalPacket) error {

	if in == nil {
		return nil
	}

	payload := in.Data()["Payload"]

	reader := bytes.NewReader(payload)

	root, err := xmltree.Parse(reader)

	if err != nil {
		return err
	}

	v, exist := d.handler[root.Name.Local]

	if !exist {
		return errors.Errorf("no data type for '%s' registered", root.Name.Local)
	}

	data := v()

	if err = xml.Unmarshal(payload, &data); err != nil {
		return err
	}

	in.Name = root.Name.Local
	in.Context = data

	return nil

}

// marshal encode the outgoing data
func (d *Dispatcher) marshal(out *ifc.LogicalPacket) error {

	if out == nil || out.Context == nil {
		return nil
	}

	data, err := xml.Marshal(out.Context)

	if err == nil {
		data, err = d.encode(out.Addr, data)
		out.Add("Payload", data)
	}

	return err

}

// encode and entity handling
func (d *Dispatcher) encode(addr string, data []byte) ([]byte, error) {

	station, _ := d.GetStationAddr(addr)

	encoding := d.GetEncodingByStation(station)

	if len(encoding) == 0 {
		return data, nil
	}

	enc, err := d.Encode(data, encoding)

	if err == nil && len(enc) > 0 {

		// xml compatible entity, ä -> &#xe4;

		if d.entityEncoding(station) {

			for i := 0; i < len(enc); i++ {

				char := enc[i]
				if char <= 127 {
					continue
				}

				entity := fmt.Sprintf("&#x%x;", char)

				padding := len(entity)

				enc = append(enc[:i+padding], enc[i+1:]...) // expand the slice

				for p := 0; p < padding; p++ { // insert entity
					enc[i+p] = entity[p]
				}

				i = i + padding - 1 // move next index

			}
		}

		return enc, nil

	}

	return data, err
}
