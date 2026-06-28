package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const defaultFilename = "wireProviderSet.go"

var filename = defaultFilename

func main() {
	flag.StringVar(&filename, "filename", defaultFilename, "output filename")
	flag.Parse()

	root := "."
	if args := flag.Args(); len(args) > 0 {
		root = args[len(args)-1]
	}
	root, err := filepath.Abs(root)
	check(err)

	existingDirs := findExistingDirs(root)
	matchDirs := findMatchDirs(root)

	allDirs := mergeDirs(existingDirs, matchDirs)
	for _, d := range sortedKeys(allDirs) {
		os.MkdirAll(d, 0755)
		outputPath := filepath.Join(d, filename)

		info := collectDirInfo(d)
		gen := generateContent(filepath.Base(d), info)

		existing, _ := os.ReadFile(outputPath)
		if bytes.Equal(existing, gen) {
			continue
		}

		fmt.Println(d)
		os.WriteFile(outputPath, gen, 0644)
	}
}
