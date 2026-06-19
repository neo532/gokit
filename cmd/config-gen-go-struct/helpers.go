package main

import "strings"

func toPascal(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, "")
}

func splitWords(s string) []string {
	var words []string
	var cur []rune
	for i, r := range s {
		switch {
		case r == '_' || r == '-' || r == '.':
			if len(cur) > 0 {
				words = append(words, string(cur))
				cur = nil
			}
		case r >= 'A' && r <= 'Z':
			if len(cur) > 0 && i > 0 {
				prev := rune(s[i-1])
				if prev >= 'a' && prev <= 'z' {
					words = append(words, string(cur))
					cur = nil
				}
			}
			cur = append(cur, r)
		default:
			cur = append(cur, r)
		}
	}
	if len(cur) > 0 {
		words = append(words, string(cur))
	}
	return words
}
