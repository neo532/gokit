package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fieldsOfInfo struct {
	structName string
	fields     []string
}

type dirInfo struct {
	providers []string
	fieldsOf  []fieldsOfInfo
}

// hasNewOrWireFieldsOf checks if a Go file contains func New or type WireFieldsOf.
func hasNewOrWireFieldsOf(path string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return false
	}
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil && strings.HasPrefix(d.Name.Name, "New") {
				return true
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if ok && strings.HasPrefix(ts.Name.Name, "WireFieldsOf") {
					return true
				}
			}
		}
	}
	return false
}

// collectDirInfo parses all .go files in a directory and extracts providers.
func collectDirInfo(dir string) dirInfo {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return dirInfo{}
	}

	var info dirInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if entry.Name() == filename || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil && strings.HasPrefix(d.Name.Name, "New") {
					info.providers = append(info.providers, d.Name.Name)
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok || !strings.HasPrefix(ts.Name.Name, "WireFieldsOf") {
						continue
					}
					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}
					var fields []string
					for _, fld := range st.Fields.List {
						for _, n := range fld.Names {
							if ast.IsExported(n.Name) {
								fields = append(fields, n.Name)
							}
						}
					}
					if len(fields) > 0 {
						info.fieldsOf = append(info.fieldsOf, fieldsOfInfo{
							structName: ts.Name.Name,
							fields:     fields,
						})
					}
				}
			}
		}
	}

	sort.Strings(info.providers)
	sort.Slice(info.fieldsOf, func(i, j int) bool {
		return info.fieldsOf[i].structName < info.fieldsOf[j].structName
	})
	return info
}
