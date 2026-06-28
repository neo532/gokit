package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// findExistingDirs finds all directories that already have a wireProviderSet.go.
func findExistingDirs(root string) map[string]bool {
	dirs := map[string]bool{}
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == filename {
			dirs[filepath.Dir(path)] = true
		}
		return nil
	})
	return dirs
}

// findMatchDirs finds all directories that contain .go files with
// func New or type WireFieldsOf but don't have a wireProviderSet.go yet.
func findMatchDirs(root string) map[string]bool {
	dirs := map[string]bool{}
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		if d.Name() == filename || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		if hasNewOrWireFieldsOf(path) {
			dirs[filepath.Dir(path)] = true
		}
		return nil
	})
	return dirs
}

// mergeDirs unions two dir maps.
func mergeDirs(a, b map[string]bool) map[string]bool {
	m := map[string]bool{}
	for d := range a {
		m[d] = true
	}
	for d := range b {
		m[d] = true
	}
	return m
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
