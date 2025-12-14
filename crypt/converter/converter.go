package converter

type Converter interface {
	String(i int64) (s string)
	Int(s string) (i int64, err error)
}
