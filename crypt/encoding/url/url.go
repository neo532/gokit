package url

import (
	"encoding/base64"

	"github.com/neo532/gokit/crypt"
)

var _ crypt.Encoding = (*Url)(nil)

type Url struct {
}

func New() *Url {
	return &Url{}
}

func (o *Url) Encode(origin []byte) (code string) {
	return base64.URLEncoding.EncodeToString(origin)
}

func (o *Url) Decode(code string) (origin []byte, err error) {
	return base64.URLEncoding.DecodeString(code)
}
