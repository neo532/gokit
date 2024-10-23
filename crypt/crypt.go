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
