package stringx

import (
	"math"
	"testing"
)

func TestItoa(t *testing.T) {
	if got := Itoa(123); got != "123" {
		t.Errorf("Itoa(123) = %q", got)
	}
	if got := Itoa(int64(-456)); got != "-456" {
		t.Errorf("Itoa(-456) = %q", got)
	}
	if got := Itoa(uint(789)); got != "789" {
		t.Errorf("Itoa(789) = %q", got)
	}
}

func TestAtoi(t *testing.T) {
	v8, err := Atoi[int8]("123")
	if err != nil || v8 != 123 {
		t.Errorf("Atoi[int8](123) = %d, %v", v8, err)
	}

	v64, err := Atoi[int64]("-456")
	if err != nil || v64 != -456 {
		t.Errorf("Atoi[int64](-456) = %d, %v", v64, err)
	}

	// overflow check
	_, err = Atoi[int8]("999")
	if err == nil {
		t.Error("Atoi[int8](999) should overflow")
	}

	// parse error
	_, err = Atoi[int]("abc")
	if err == nil {
		t.Error("Atoi[int](abc) expected error")
	}
}

func TestAtoui(t *testing.T) {
	v8, err := AtoUi[uint8]("123")
	if err != nil || v8 != 123 {
		t.Errorf("Atoui[uint8](123) = %d, %v", v8, err)
	}

	v64, err := AtoUi[uint64]("456")
	if err != nil || v64 != 456 {
		t.Errorf("Atoui[uint64](456) = %d, %v", v64, err)
	}

	// overflow
	_, err = AtoUi[uint8]("999")
	if err == nil {
		t.Error("Atoui[uint8](999) should overflow")
	}

	// negative
	_, err = AtoUi[uint]("-1")
	if err == nil {
		t.Error("Atoui[uint](-1) expected error")
	}
}

func TestFtoa(t *testing.T) {
	if got := Ftoa(3.14); got != "3.14" {
		t.Errorf("Ftoa(3.14) = %q", got)
	}
}

func TestAtof(t *testing.T) {
	v, err := Atof("3.14")
	if err != nil || v != 3.14 {
		t.Errorf("Atof(3.14) = %f, %v", v, err)
	}

	v, err = Atof("1.7976931348623157e+308")
	if err != nil || v != math.MaxFloat64 {
		t.Errorf("Atof(MaxFloat64) = %e, %v", v, err)
	}

	_, err = Atof("abc")
	if err == nil {
		t.Error("Atof(abc) expected error")
	}
}

func TestFtoaTrunc(t *testing.T) {
	if got := FtoaTrunc(3.14159); got != "3.14" {
		t.Errorf("Ftoa(3.14159) = %q", got)
	}
}
