package filepath

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "config.yaml"), "hello")

	w := New(dir)
	data, err := w.Read("config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q, want %q", string(data), "hello")
	}
}

func TestWatchInitialCall(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.yaml"), "va")
	writeFile(t, filepath.Join(dir, "b.yaml"), "vb")

	w := New(dir)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	got := make(map[string]string, 2)
	if err := w.Watch(ctx, func(name string, data []byte) error {
		got[name] = string(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if got["a.yaml"] != "va" {
		t.Fatalf("a.yaml: got %q, want %q", got["a.yaml"], "va")
	}
	if got["b.yaml"] != "vb" {
		t.Fatalf("b.yaml: got %q, want %q", got["b.yaml"], "vb")
	}
}

func TestWatchFileChange(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "c.yaml"), "before")

	w := New(dir)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	changed := make(chan string, 2)
	if err := w.Watch(ctx, func(name string, data []byte) error {
		changed <- name
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// Drain initial event.
	<-changed

	time.Sleep(100 * time.Millisecond)
	writeFile(t, filepath.Join(dir, "c.yaml"), "after")

	select {
	case name := <-changed:
		if name != "c.yaml" {
			t.Fatalf("got filename %q, want %q", name, "c.yaml")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for change event")
	}
}

func TestWatchNewFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "existing.yaml"), "old")

	w := New(dir)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	changed := make(chan string, 2)
	if err := w.Watch(ctx, func(name string, data []byte) error {
		changed <- name
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// Drain initial event.
	<-changed

	time.Sleep(100 * time.Millisecond)
	writeFile(t, filepath.Join(dir, "new.yaml"), "new content")

	select {
	case name := <-changed:
		if name != "new.yaml" {
			t.Fatalf("got filename %q, want %q", name, "new.yaml")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for new file event")
	}
}

func TestClose(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "test.yaml"), "test")

	w := New(dir)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := w.Watch(ctx, func(name string, data []byte) error { return nil }); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
