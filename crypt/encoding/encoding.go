package encoding

type Encoding interface {
	Encode(origin []byte) (code string)
	Decode(code string) (origin []byte, err error)
}
