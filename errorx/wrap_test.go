package errorx

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func line() (line int) {
	_, _, line, _ = runtime.Caller(1)
	return
}

func TestNew(t *testing.T) {

	err := New("format args %d", 1)

	r := fmt.Sprintf("errorx/wrap_test.go:%d---format args 1", line()-2)
	if err.Error() != r {
		t.Errorf("%s has err", t.Name())
	}
}

func TestWrap(t *testing.T) {

	err := errors.New("err")
	err = Wrap(err)

	r := fmt.Sprintf("errorx/wrap_test.go:%d---err", line()-2)
	if err.Error() != r {
		t.Errorf("%s has err", t.Name())
	}
}

func TestWrapf(t *testing.T) {

	err := errors.New("err")
	err = Wrapf(err, "format args %d", 1)

	r := fmt.Sprintf("err===errorx/wrap_test.go:%d---format args 1", line()-2)
	if err.Error() != r {
		t.Errorf("%s has err", t.Name())
	}
}

func TestWrapError(t *testing.T) {

	err := errors.New("err")
	err = WrapError(err, errors.New("err1"))

	r := fmt.Sprintf("err===errorx/wrap_test.go:%d---err1", line()-2)
	if err.Error() != r {
		t.Errorf("%s has err", t.Name())
	}
}

func TestWrapErrorf(t *testing.T) {

	err := errors.New("err")
	err = WrapErrorf(err, errors.New("err1"), "format args %d", 1)

	r := fmt.Sprintf("err===err1===errorx/wrap_test.go:%d---format args 1", line()-2)
	if err.Error() != r {
		t.Errorf("%s has err", t.Name())
	}
}

func TestIsCauseBy(t *testing.T) {

	err := errors.New("err")
	err1 := WrapErrorf(err, errors.New("err1"), "format args %d", 1)

	if !IsCauseBy(err1, err) {
		t.Errorf("%s has err", t.Name())
	}
}

func TestCausePurly(t *testing.T) {

	err := New("err")
	err1 := WrapErrorf(err, errors.New("err1"), "format args %d", 1)

	if CausePurly(err1).Error() != "err" {
		t.Errorf("%s has err", t.Name())
	}
}
