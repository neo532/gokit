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
	cv("1.2.3", "1.2.3 ", Equal)

	fmt.Println(t.Name())
}

func TestPageExec(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		pageSize int
		pages    map[int][]int
	}{
		{
			name:  "total zero",
			total: 0, pageSize: 3,
			pages: map[int][]int{},
		},
		{
			name:  "pageSize zero",
			total: 10, pageSize: 0,
			pages: map[int][]int{},
		},
		{
			name:  "one page",
			total: 5, pageSize: 5,
			pages: map[int][]int{
				1: {0, 1, 2, 3, 4},
			},
		},
		{
			name:  "pageSize larger than total",
			total: 3, pageSize: 10,
			pages: map[int][]int{
				1: {0, 1, 2},
			},
		},
		{
			name:  "exact multiple",
			total: 9, pageSize: 3,
			pages: map[int][]int{
				1: {0, 1, 2},
				2: {3, 4, 5},
				3: {6, 7, 8},
			},
		},
		{
			name:  "non exact multiple",
			total: 10, pageSize: 3,
			pages: map[int][]int{
				1: {0, 1, 2},
				2: {3, 4, 5},
				3: {6, 7, 8},
				4: {9},
			},
		},
		{
			name:  "single element pages",
			total: 3, pageSize: 1,
			pages: map[int][]int{
				1: {0},
				2: {1},
				3: {2},
			},
		},
		{
			name:  "pageSize two with odd total",
			total: 7, pageSize: 3,
			pages: map[int][]int{
				1: {0, 1, 2},
				2: {3, 4, 5},
				3: {6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PageExec(int64(tt.total), tt.pageSize, func(b, e int64, p int) error {
				expected, ok := tt.pages[p]
				if !ok {
					t.Errorf("unexpected page %d", p)
					return nil
				}
				var got []int
				for i := b; i < e; i++ {
					got = append(got, int(i))
				}
				if !CompareSlice(got, expected) {
					t.Errorf("page %d: got %+v, want %+v", p, got, expected)
				}
				return nil
			})
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}
