package slicex

import (
	"strings"
	"testing"
)

func TestRange(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var got []string
		Range([]string{"a", "b", "c"}, func(v string) { got = append(got, v) })
		if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
			t.Errorf("got %v, want [a b c]", got)
		}
	})

	t.Run("nil", func(t *testing.T) {
		Range(nil, func(v int) { t.Error("should not be called") })
	})

	t.Run("nested", func(t *testing.T) {
		var b strings.Builder
		for _, s := range [][]string{{"a", "b"}, {"x", "y"}} {
			Range(s, func(v string) { b.WriteString(v) })
		}
		if b.String() != "abxy" {
			t.Errorf("got %q, want %q", b.String(), "abxy")
		}
	})
}
