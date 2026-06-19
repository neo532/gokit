package filepath

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors a directory for file changes.
type Watcher struct {
	dir    string
	mu     sync.Mutex
	w      *fsnotify.Watcher
	closed chan struct{}
}

// New creates a Watcher for the given directory.
func New(dir string) *Watcher {
	return &Watcher{
		dir:    dir,
		closed: make(chan struct{}),
	}
}

// Read reads a file within the watched directory.
func (w *Watcher) Read(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(w.dir, name))
}

// Watch monitors the directory for file changes and calls fn with the filename
// and content. Before watching begins, fn is called once synchronously for every
// file in the directory. Subsequent calls happen asynchronously on writes/creates;
// errors from fn are ignored to keep watching.
func (w *Watcher) Watch(ctx context.Context, fn func(fileName string, data []byte) error) error {
	// Sync: deliver initial state for all files.
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(w.dir, e.Name()))
		if err != nil {
			continue
		}
		if err := fn(e.Name(), data); err != nil {
			continue
		}
	}

	// Async: watch for changes.
	return w.startWatch(ctx, fn)
}

func (w *Watcher) startWatch(ctx context.Context, fn func(fileName string, data []byte) error) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.w != nil {
		return nil
	}

	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := fw.Add(w.dir); err != nil {
		fw.Close()
		return err
	}

	w.w = fw
	go w.loop(ctx, fw, fn)
	return nil
}

func (w *Watcher) loop(ctx context.Context, fw *fsnotify.Watcher, fn func(fileName string, data []byte) error) {
	defer fw.Close()

	var timer *time.Timer
	var timerC <-chan time.Time
	var pending string

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.closed:
			return
		case event, ok := <-fw.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				pending = filepath.Base(event.Name)
				if timer == nil {
					timer = time.NewTimer(50 * time.Millisecond)
					timerC = timer.C
				} else {
					timer.Reset(50 * time.Millisecond)
				}
			}
		case <-timerC:
			timer = nil
			timerC = nil
			if pending == "" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(w.dir, pending))
			if err != nil {
				continue
			}
			_ = fn(pending, data)
			pending = ""
		case _, ok := <-fw.Errors:
			if !ok {
				return
			}
		}
	}
}

// Close stops the watcher.
func (w *Watcher) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.w == nil {
		return nil
	}
	close(w.closed)
	w.w.Close()
	return nil
}
