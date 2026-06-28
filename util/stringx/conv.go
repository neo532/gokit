package stringx

import (
	"math"
	"strconv"
	"unsafe"
)

// Atoi parses a decimal string to a signed integer type.
func Atoi[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s string) (T, error) {
	var zero T
	v, err := strconv.ParseInt(s, 10, int(unsafe.Sizeof(zero)*8))
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

// AtoUi parses a decimal string to an unsigned integer type.
func AtoUi[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](s string) (T, error) {
	var zero T
	v, err := strconv.ParseUint(s, 10, int(unsafe.Sizeof(zero)*8))
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

// Atof parses a string to float64.
func Atof(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Ftoa converts float64 to decimal string.
// Precision rounds to the given decimal places; without it, trailing zeros are stripped.
func Ftoa(f float64, prec ...int) string {
	dec := -1
	if len(prec) > 0 {
		dec = prec[0]
	}
	return strconv.FormatFloat(f, 'f', dec, 64)
}

// FtoaTrunc converts float64 to decimal string, truncating without rounding.
func FtoaTrunc(num float64, prec ...int) string {
	d := float64(100)
	if len(prec) > 0 {
		d = math.Pow10(prec[0])
	}
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}

// Itoa converts any integer to a decimal string.
func Itoa[T ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64,
](v T) string {
	return strconv.FormatInt(int64(v), 10)
}
