package dictionary

import "testing"

func TestDictionary_DefaultDictionaryRoundTrip(t *testing.T) {
	d := New()

	cases := []int64{0, 1, 2, 35, 36, 37, 1296, 1234567890123456}
	for _, num := range cases {
		str := d.String(num)
		got, err := d.Int(str)
		if err != nil {
			t.Fatalf("Int(%q) returned unexpected error: %v", str, err)
		}
		if got != num {
			t.Fatalf("round-trip mismatch for %d: got %d (encoded as %q)", num, got, str)
		}
	}
}

func TestDictionary_CustomDictionaryRoundTrip(t *testing.T) {
	d := New(WithDictionary("abc"))

	cases := []int64{0, 1, 2, 3, 5, 8, 26}
	for _, num := range cases {
		str := d.String(num)
		got, err := d.Int(str)
		if err != nil {
			t.Fatalf("Int(%q) returned unexpected error: %v", str, err)
		}
		if got != num {
			t.Fatalf("round-trip mismatch for %d: got %d (encoded as %q)", num, got, str)
		}
	}
}
