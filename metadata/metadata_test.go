package metadata

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		input    []map[string][]string
		expected Metadata
	}{
		{
			name:     "empty input",
			input:    []map[string][]string{},
			expected: Metadata{},
		},
		{
			name: "single map",
			input: []map[string][]string{
				{"key1": {"value1"}, "key2": {"value2"}},
			},
			expected: Metadata{
				"key1": {"value1"},
				"key2": {"value2"},
			},
		},
		{
			name: "multiple maps",
			input: []map[string][]string{
				{"key1": {"value1"}},
				{"key2": {"value2"}},
			},
			expected: Metadata{
				"key1": {"value1"},
				"key2": {"value2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.input...)
			if len(result) != len(tt.expected) {
				t.Errorf("New() = %v, want %v", result, tt.expected)
			}
			for k, v := range tt.expected {
				if got := result[k]; len(got) != len(v) {
					t.Errorf("New()[%s] = %v, want %v", k, got, v)
				}
			}
		})
	}
}

func TestMetadata_Add(t *testing.T) {
	md := Metadata{}

	// Test adding single value
	md.Add("key1", "value1")
	if got := md["key1"]; len(got) != 1 || got[0] != "value1" {
		t.Errorf("Add() = %v, want [value1]", got)
	}

	// Test adding multiple values
	md.Add("key1", "value2")
	if got := md["key1"]; len(got) != 2 || got[1] != "value2" {
		t.Errorf("Add() = %v, want [value1 value2]", got)
	}

	// Test empty key
	md.Add("", "value3")
	if _, exists := md[""]; exists {
		t.Error("Add() should not add empty key")
	}

	// Test case insensitivity
	md.Add("KEY1", "value3")
	if got := md["key1"]; len(got) != 3 || got[2] != "value3" {
		t.Errorf("Add() = %v, want [value1 value2 value3]", got)
	}
}

func TestMetadata_Get(t *testing.T) {
	md := Metadata{
		"key1": {"value1", "value2"},
		"key2": {"value3"},
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"existing key", "key1", "value1"},
		{"existing key uppercase", "KEY1", "value1"},
		{"non-existing key", "key3", ""},
		{"empty key", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := md.Get(tt.key); got != tt.expected {
				t.Errorf("Get() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMetadata_Set(t *testing.T) {
	md := Metadata{}

	// Test setting new value
	md.Set("key1", "value1")
	if got := md["key1"]; len(got) != 1 || got[0] != "value1" {
		t.Errorf("Set() = %v, want [value1]", got)
	}

	// Test overwriting existing value
	md.Set("key1", "value2")
	if got := md["key1"]; len(got) != 1 || got[0] != "value2" {
		t.Errorf("Set() = %v, want [value2]", got)
	}

	// Test empty key and value
	md.Set("", "value3")
	md.Set("key2", "")
	if _, exists := md[""]; exists {
		t.Error("Set() should not set empty key")
	}
	if _, exists := md["key2"]; exists {
		t.Error("Set() should not set empty value")
	}
}

func TestMetadata_Values(t *testing.T) {
	md := Metadata{
		"key1": {"value1", "value2"},
		"key2": {"value3"},
	}

	tests := []struct {
		name     string
		key      string
		expected []string
	}{
		{"existing key", "key1", []string{"value1", "value2"}},
		{"existing key uppercase", "KEY1", []string{"value1", "value2"}},
		{"non-existing key", "key3", nil},
		{"empty key", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := md.Values(tt.key)
			if len(got) != len(tt.expected) {
				t.Errorf("Values() = %v, want %v", got, tt.expected)
			}
			for i, v := range tt.expected {
				if got[i] != v {
					t.Errorf("Values()[%d] = %v, want %v", i, got[i], v)
				}
			}
		})
	}
}

func TestMetadata_Clone(t *testing.T) {
	original := Metadata{
		"key1": {"value1", "value2"},
		"key2": {"value3"},
	}

	clone := original.Clone()

	// Test if clone is successful
	if len(clone) != len(original) {
		t.Errorf("Clone() = %v, want %v", clone, original)
	}

	// Test if modifying clone affects original data
	clone["key1"][0] = "newvalue"
	if original["key1"][0] == "newvalue" {
		t.Error("Clone() should create a deep copy")
	}
}

func TestContextFunctions(t *testing.T) {
	md := Metadata{
		"key1": {"value1"},
		"key2": {"value2"},
	}

	// Test server context
	ctx := context.Background()
	serverCtx := NewServerContext(ctx, md)
	if got, ok := FromServerContext(serverCtx); !ok || len(got) != len(md) {
		t.Errorf("FromServerContext() = %v, %v, want %v, true", got, ok, md)
	}

	// Test client context
	clientCtx := NewClientContext(ctx, md)
	if got, ok := FromClientContext(clientCtx); !ok || len(got) != len(md) {
		t.Errorf("FromClientContext() = %v, %v, want %v, true", got, ok, md)
	}

	// Test AppendToClientContext
	appendedCtx := AppendToClientContext(ctx, "key3", "value3", "key4", "value4")
	if got, ok := FromClientContext(appendedCtx); !ok || len(got) != 2 {
		t.Errorf("AppendToClientContext() = %v, %v, want 2 values, true", got, ok)
	}

	// Test MergeToClientContext
	newMd := Metadata{"key5": {"value5"}}
	mergedCtx := MergeToClientContext(ctx, newMd)
	if got, ok := FromClientContext(mergedCtx); !ok || len(got) != 1 {
		t.Errorf("MergeToClientContext() = %v, %v, want 1 value, true", got, ok)
	}
}
