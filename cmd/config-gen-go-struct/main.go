package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	pkg := flag.String("pkg", "config", "package name")
	typeName := flag.String("type", "Config", "top-level struct name")
	split := flag.Bool("split", false, "split output into one file per top-level section")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: config-gen-go-struct [-pkg name] [-type name] [-split] <input.yaml|json|ini>... [output.go|output-dir]\n")
		os.Exit(1)
	}

	// Detect multi-file mode: all args are config files (end with .yaml/.yml/.json/.ini)
	allConfig := true
	for _, a := range args {
		ext := filepath.Ext(a)
		if ext != ".yaml" && ext != ".yml" && ext != ".json" && ext != ".ini" {
			allConfig = false
			break
		}
	}

	if allConfig && len(args) >= 2 {
		// Multi-file mode: each input generates its own .go file
		var entries []unifiedEntry
		for _, input := range args {
			generateOne(*pkg, input, *split, *typeName)
			pascalBase := toPascal(stripExt(filepath.Base(input)))
			entries = append(entries, unifiedEntry{
				structName: *typeName + pascalBase,
				loadName:   "Load" + pascalBase,
				yamlFile:   filepath.Base(input),
			})
		}
		// Generate unified config.go aggregating all sections
		code := generateUnified(*pkg, *typeName, entries)
		out := filepath.Join(filepath.Dir(args[0]), "config.go")
		if err := os.WriteFile(out, []byte(code), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "generated %s\n", out)
		return
	}

	// Single-file mode (original behavior with 1 or 2 args)
	input := args[0]
	raw, err := readConfig(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", input, err)
		os.Exit(1)
	}

	ng := newNameGen()
	fields := inferFields(raw, ng, "")
	format := detectFormat(input)

	// Auto-split when output path doesn't end with .go
	autoSplit := false
	if !*split && len(args) > 1 {
		outPath := args[1]
		if !strings.HasSuffix(outPath, ".go") {
			autoSplit = true
		}
	}

	if *split || autoSplit {
		outDir := "."
		baseName := stripExt(filepath.Base(input))
		if len(args) > 1 {
			outPath := args[1]
			if fi, err := os.Stat(outPath); err == nil && fi.IsDir() {
				outDir = outPath
			} else {
				outDir = filepath.Dir(outPath)
				baseName = stripExt(filepath.Base(outPath))
			}
		}
		files := generateSplit(*pkg, *typeName, fields, filepath.Base(input), format, baseName)
		for name, code := range files {
			path := filepath.Join(outDir, name)
			if err := os.WriteFile(path, []byte(code), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "generated %s\n", path)
		}
		return
	}

	code := generate(*pkg, *typeName, "Load", fields, filepath.Base(input), format)

	var out string
	if len(args) > 1 {
		out = args[1]
	}
	if out == "" {
		out = changeExt(input, ".go")
	}
	if err := os.WriteFile(out, []byte(code), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "generated %s\n", out)
}

func generateOne(pkg, input string, split bool, typeSuffix string) {
	raw, err := readConfig(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", input, err)
		os.Exit(1)
	}

	base := stripExt(filepath.Base(input))
	pascalBase := toPascal(base)
	format := detectFormat(input)
	out := changeExt(input, ".go")

	ng := newNameGen()
	fields := inferFields(raw, ng, pascalBase)

	if split {
		files := generateSplit(pkg, typeSuffix+pascalBase, fields, filepath.Base(input), format, base)
		for name, code := range files {
			path := filepath.Join(filepath.Dir(input), name)
			if err := os.WriteFile(path, []byte(code), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "generated %s\n", path)
		}
		return
	}

	typeName := typeSuffix + pascalBase
	loadName := "Load" + pascalBase
	code := generate(pkg, typeName, loadName, fields, filepath.Base(input), format)
	if err := os.WriteFile(out, []byte(code), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "generated %s\n", out)
}

func detectFormat(path string) string {
	switch filepath.Ext(path) {
	case ".json":
		return "json"
	case ".ini":
		return "ini"
	default:
		return "yaml"
	}
}

func changeExt(path, newExt string) string {
	base := stripExt(path)
	if base == path {
		return path + newExt
	}
	return base + newExt
}

func stripExt(path string) string {
	for _, ext := range []string{".yaml", ".yml", ".json", ".ini"} {
		if strings.HasSuffix(path, ext) {
			return strings.TrimSuffix(path, ext)
		}
	}
	return path
}
