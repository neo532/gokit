package std

import (
	"encoding/base64"
	"github.com/neo532/gokit/crypt/encoding"
)

var _ encoding.Encoding = (*Std)(nil)

type Std struct {
}

func New() *Std {
	return &Std{}
}

func (o *Std) Encode(origin []byte) (code string) {
	return base64.StdEncoding.EncodeToString(origin)
}

func (o *Std) Decode(code string) (origin []byte, err error) {
	return base64.StdEncoding.DecodeString(code)
}
