package dictionary

import (
	"github.com/neo532/gokit/crypt/converter"
	"github.com/neo532/gokit/util"
	"math/big"
	"strings"
)

var _ converter.Converter = (*Dictionary)(nil)

type Dictionary struct {
	dict string
}

type opt func(o *Dictionary)

func WithDictionary(s string) opt {
	return func(o *Dictionary) {
		o.dict = s
	}
}

func New(opts ...opt) (os *Dictionary) {
	os = &Dictionary{
		dict: "abcdefghijklmnopqrstuvwxyz1234567890",
	}
	for _, fn := range opts {
		fn(os)
	}
	return os
}

func (o *Dictionary) String(i int64) (s string) {

	var str strings.Builder
	lenD := int64(len(o.dict))
	for {
		if i <= 0 {
			break
		}
		str.WriteString(string(o.dict[i%lenD]))
		i = i / lenD
	}
	return util.Reverse(str.String())
}

func (o *Dictionary) Int(s string) (i int64, err error) {
	lenD := len(o.dict)
	lenS := len(s)

	rst := big.NewInt(0)
	j := big.NewInt(1)
	var pos int
	for i := 0; i < lenS; i++ {
		pos = strings.Index(o.dict, string(s[i]))

		j = j.Exp(big.NewInt(int64(lenD)), big.NewInt(int64(lenS-i-1)), nil)
		j = j.Mul(big.NewInt(int64(pos)), j)
		rst = rst.Add(rst, j)
	}
	i = rst.Int64()
	return
}
