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
)

// CompareVersion returns the num after comparing two versions.
//
//	1: ver1 > ver2, Larger
//	0: ver1 = ver2, Equal
//
// -1: ver1 < ver2, Smaller
func CompareVersion(ver1, ver2 string) (int, error) {
	v1 := strings.Split(strings.Trim(strings.TrimSpace(ver1), "."), ".")
	v2 := strings.Split(strings.Trim(strings.TrimSpace(ver2), "."), ".")

	v1Len := len(v1)
	v2Len := len(v2)

	maxLen := max(v1Len, v2Len)

	for i := 0; i < maxLen; i++ {
		var v1i, v2i int
		if i < v1Len {
			n, err := strconv.Atoi(v1[i])
			if err != nil {
				return 0, err
			}
			v1i = n
		}
		if i < v2Len {
			n, err := strconv.Atoi(v2[i])
			if err != nil {
				return 0, err
			}
			v2i = n
		}
		if v1i < v2i {
			return Smaller, nil
		}
		if v1i > v2i {
			return Larger, nil
		}
	}
	return Equal, nil
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
