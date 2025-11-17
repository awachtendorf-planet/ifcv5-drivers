package vendor

import (
	"bytes"
	"encoding/xml"
	"time"
)

type xsdDate time.Time

func (t *xsdDate) UnmarshalText(text []byte) error {
	return _unmarshalTime(text, (*time.Time)(t), "060102")
}

func (t xsdDate) MarshalText() ([]byte, error) {
	return _marshalTime((time.Time)(t), "060102")
}

func (t xsdDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if (time.Time)(t).IsZero() {
		return nil
	}
	m, err := t.MarshalText()
	if err != nil {
		return err
	}
	return e.EncodeElement(m, start)
}

func (t xsdDate) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if (time.Time)(t).IsZero() {
		return xml.Attr{}, nil
	}
	m, err := t.MarshalText()
	return xml.Attr{Name: name, Value: string(m)}, err
}

func _unmarshalTime(text []byte, t *time.Time, format string) (err error) {
	s := string(bytes.TrimSpace(text))
	*t, err = time.Parse(format, s)
	if _, ok := err.(*time.ParseError); ok {
		*t, err = time.Parse(format+"Z07:00", s)
	}
	return err
}

func _marshalTime(t time.Time, format string) ([]byte, error) {
	//return []byte(t.Format(format + "Z07:00")), nil
	return []byte(t.Format(format)), nil
}

type xsdTime time.Time

func (t *xsdTime) UnmarshalText(text []byte) error {
	//return _unmarshalTime(text, (*time.Time)(t), "150405.999999999")
	return _unmarshalTime(text, (*time.Time)(t), "150405")
}

func (t xsdTime) MarshalText() ([]byte, error) {
	//return _marshalTime((time.Time)(t), "150405.999999999")
	return _marshalTime((time.Time)(t), "150405")
}

func (t xsdTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if (time.Time)(t).IsZero() {
		return nil
	}
	m, err := t.MarshalText()
	if err != nil {
		return err
	}
	return e.EncodeElement(m, start)
}

func (t xsdTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if (time.Time)(t).IsZero() {
		return xml.Attr{}, nil
	}
	m, err := t.MarshalText()
	return xml.Attr{Name: name, Value: string(m)}, err
}
