package marshaler

import (
	"github.com/neo532/gokit/crypt"
)

var registeredMarshalers = make(map[string]crypt.Marshaler)

func RegisterMarshaler(marshaler crypt.Marshaler) {
	if marshaler == nil {
		panic("cannot register a nil Marshaler")
	}
	if marshaler.Name() == "" {
		panic("cannot register Marshaler with empty string result for Name()")
	}
	registeredMarshalers[marshaler.Name()] = marshaler
}

func GetMarshaler(contentSubtype string) crypt.Marshaler {
	return registeredMarshalers[contentSubtype]
}
