package errorx

/*
 * @abstract errorx
 * @mail neo532@126.com
 */

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var (
	DelimiterError    = "==="
	DelimiterPosition = "---"
	Grade             = 2
)

func New(format string, args ...interface{}) error {

	return fmt.Errorf("%s%s%s",
		caller(2),
		DelimiterPosition,
		fmt.Sprintf(format, args...),
	)
}

func Wrap(err error) error {

	if err == nil {
		return nil
	}
	return fmt.Errorf("%s%s%w",
		caller(2),
		DelimiterPosition,
		err,
	)
}

func Wrapf(err error, format string, args ...interface{}) error {

	if err == nil {
		return nil
	}
	return fmt.Errorf("%w%s%s%s%s",
		err,
		DelimiterError,
		caller(2),
		DelimiterPosition,
		fmt.Sprintf(format, args...),
	)
}

func WrapError(err, err1 error) error {

	if err == nil {
		return err1
	}
	return fmt.Errorf("%w%s%s%s%w",
		err,
		DelimiterError,
		caller(2),
		DelimiterPosition,
		err1,
	)
}

func WrapErrorf(err, err1 error, format string, args ...interface{}) error {

	if err == nil {
		return Wrapf(err1, format, args...)
	}
	return fmt.Errorf("%w%s%w%s%s%s%s",
		err,
		DelimiterError,
		err1,
		DelimiterError,
		caller(2),
		DelimiterPosition,
		fmt.Sprintf(format, args...),
	)
}

func IsCauseBy(err, err1 error) (b bool) {

	if err == nil || err1 == nil {
		return
	}
	if strings.SplitN(err.Error(), DelimiterError, 2)[0] == err1.Error() {
		return true
	}
	return
}

func CausePurly(err error) (ep error) {

	if err == nil {
		return
	}
	e := strings.SplitN(err.Error(), DelimiterError, 2)[0]
	es := strings.Split(e, DelimiterPosition)
	if l := len(es); l > 0 {
		return fmt.Errorf(es[l-1])
	}
	return fmt.Errorf(e)
}

func caller(depth int) (r string) {

	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return
	}

	r = ":" + strconv.Itoa(line)

	d := string(os.PathSeparator)
	var l int
	for i := len(file) - 1; i >= 0; i-- {
		f := string(file[i])
		if d == f {
			l++
			if l >= Grade {
				break
			}
		}
		r = f + r
	}
	return
}
