package marshaler

type Marshaler interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
}

var registeredMarshalers = make(map[string]Marshaler)

func RegisterMarshaler(marshaler Marshaler) {
	if marshaler == nil {
		panic("cannot register a nil Marshaler")
	}
	if marshaler.Name() == "" {
		panic("cannot register Marshaler with empty string result for Name()")
	}
	registeredMarshalers[marshaler.Name()] = marshaler
}

func GetMarshaler(contentSubtype string) Marshaler {
	return registeredMarshalers[contentSubtype]
}
