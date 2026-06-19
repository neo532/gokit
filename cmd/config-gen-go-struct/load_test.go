package main

import (
	"os"
	"testing"
)

func TestGeneratedConfigLoad(t *testing.T) {
	data, err := os.ReadFile("testdata/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	// Using package path directly
	_ = data
}
