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
)

var Delimiter = "=========="
var Grade = 2

func New(format string, args ...interface{}) error {

	return fmt.Errorf("[%s] %s",
		caller(2),
		fmt.Sprintf(format, args...),
	)
}

func Wrap(err error) error {

	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] %w",
		caller(2),
		err,
	)
}

func Wrapf(err error, format string, args ...interface{}) error {

	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] %s has error:%w",
		caller(2),
		fmt.Sprintf(format, args...),
		err,
	)
}

func WrapError(err, err1 error) error {

	if err == nil {
		return err1
	}
	return fmt.Errorf("[%s] %w %s %w",
		caller(2),
		err1,
		Delimiter,
		err,
	)
}

func WrapErrorf(err, err1 error, format string, args ...interface{}) error {

	if err == nil {
		return fmt.Errorf("[%s] %s has error:%w",
			caller(2),
			fmt.Sprintf(format, args...),
			err1,
		)
	}
	return fmt.Errorf("[%s] %s has error:%w %s %w",
		caller(2),
		fmt.Sprintf(format, args...),
		err1,
		Delimiter,
		err,
	)
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
