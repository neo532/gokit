package util

import (
	"fmt"
	"testing"
)

func TestCompareVersion(t *testing.T) {

	cv := func(v1, v2 string, r int) {
		if v, err := CompareVersion(v1, v2); err != nil || v != r {
			t.Errorf("%s: %s-%s is %d should %d,err %+v", t.Name(), v1, v2, v, r, err)
		}
	}

	cv("1.2", "1.2.3", Smaller)
	cv("1.4", "1.2.3", Larger)
	cv("1.2.3", "1.2.3", Equal)

	fmt.Println(t.Name())
}

func TestPageExec(t *testing.T) {

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	m := map[int][]int{
		1: []int{1, 2, 3},
		2: []int{4, 5, 6},
		3: []int{7, 8, 9},
		4: []int{10},
	}

	// [1 2 3] [4 5 6] [7 8 9] [10]
	err := PageExec(int64(len(arr)), 3, func(b, e int64, p int) (err error) {
		a := arr[b:e]
		if v, ok := m[p]; ok && !CompareSlice(a, v) {
			t.Errorf("%s is %+v should %+v", t.Name(), a, v)
		}
		return
	})
	if err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	fmt.Println(t.Name())
}
