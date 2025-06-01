package crypt

type Crypt interface {
	Encrypt(origin []byte) (encrpy string, err error)
	Decrypt(encrpy string) (origin []byte, err error)
}

type Encoding interface {
	Encode(origin []byte) (code string)
	Decode(code string) (origin []byte, err error)
}

type Compressor interface {
	String(i int64) (s string)
	Int(s string) (i int64, err error)
}

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
