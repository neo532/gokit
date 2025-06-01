package xml

import (
	"encoding/xml"

	"github.com/neo532/gokit/crypt/marshaler"
)

func init() {
	marshaler.RegisterMarshaler(NewXml())
}

type opt func(cc *Xml)

// Codec is a Codec implementation with xml.
type Xml struct {
}

func NewXml(opts ...opt) (cc *Xml) {
	cc = &Xml{}
	for _, o := range opts {
		o(cc)
	}
	return
}

func (cc *Xml) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (cc *Xml) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (cc *Xml) Name() string {
	return "xml"
}
