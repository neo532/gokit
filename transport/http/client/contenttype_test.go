package client

import (
	"testing"
)

func TestContentSubtype(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "json",
			input: "application/json",
			want:  "json",
		},
		{
			name:  "json with charset",
			input: "application/json;charset=utf-8",
			want:  "json",
		},
		{
			name:  "xml",
			input: "application/xml",
			want:  "xml",
		},
		{
			name:  "form urlencoded",
			input: "application/x-www-form-urlencoded;charset=utf-8",
			want:  "x-www-form-urlencoded",
		},
		{
			name:  "multipart without application prefix",
			input: "multipart/form-data",
			want:  "json",
		},
		{
			name:  "empty string",
			input: "",
			want:  "json",
		},
		{
			name:  "case insensitive",
			input: "Application/JSON;Charset=UTF-8",
			want:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContentSubtype(tt.input)
			if got != tt.want {
				t.Errorf("ContentSubtype(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
