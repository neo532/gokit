package stringx

import (
	"unicode"
)

// Reverse returns the string after reversing.
func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// UpperFirstChar returns the string,the first letter is upper.
func UpperFirstChar(str string) string {
	r := []rune(str)
	if len(r) == 0 {
		return ""
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
