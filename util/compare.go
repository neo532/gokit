package util

/*
 * @abstract compare
 * @author liuxiaofeng
 * @mail neo532@126.com
 * @date 2024-10-23
 */

import (
	"strconv"
	"strings"
)

const (
	// Larger is value of compare.It means ver1 is larger than ver2.
	Larger = 1
	// Smaller is value of compare.It means ver1 is smaller than ver2.
	Smaller = -1
	// Equal is value of compare.It means ver1 is equal with ver2.
	Equal = 0
	// Error is value of compare.It means having a error.
	Error = -2
)

// CompareVersion returns the num after comparing two versions.
//
//	1: ver1 > ver2, Larger
//	0: ver1 = ver2, Equal
//
// -1: ver1 < ver2, Smaller
// -2: has error,	Error
func CompareVersion(ver1, ver2 string) (r int) {

	v1 := strings.Split(strings.Trim(ver1, "."), ".")
	v2 := strings.Split(strings.Trim(ver2, "."), ".")
	v1Len := len(v1)
	v2Len := len(v2)

	lMin := v1Len
	r = Smaller

	diff := v1Len - v2Len
	if diff > 0 {
		lMin = v2Len
		r = Larger
	}

	var err error
	var v1i, v2i int
	for i := 0; i < lMin; i++ {
		if v1i, err = strconv.Atoi(v1[i]); err != nil {
			return Error
		}
		if v2i, err = strconv.Atoi(v2[i]); err != nil {
			return Error
		}
		if v1i < v2i {
			return Smaller
		}
		if v1i > v2i {
			return Larger
		}
	}

	if diff == 0 {
		return Equal
	}

	return
}

func CompareSlice[T comparable](v1, v2 []T) (b bool) {

	l := len(v1)
	if l != len(v2) {
		return
	}
	if l == 0 {
		return true
	}

	for i := 0; i < l; i++ {
		if v1[i] != v2[i] {
			return
		}
	}
	return true
}
