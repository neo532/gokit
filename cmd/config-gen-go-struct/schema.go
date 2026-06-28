package main

import (
	"fmt"
	"sort"
	"strings"
)

type nameGen struct {
	used       map[string]bool
	structDefs map[string][]fieldInfo // accumulated fields per struct type (union merge)
}

func newNameGen() *nameGen {
	return &nameGen{
		used:       map[string]bool{},
		structDefs: map[string][]fieldInfo{},
	}
}

// mergeStruct accumulates fields for a struct type, returning the union of all
// occurrences. When the same _-prefixed struct key appears in multiple places
// (e.g. database_user._database_conf and database_biz._database_conf), this
// ensures the generated struct has all fields from every occurrence.
func (g *nameGen) mergeStruct(structType string, fields []fieldInfo) []fieldInfo {
	existing, ok := g.structDefs[structType]
	if !ok {
		g.structDefs[structType] = fields
		g.used[structType] = true
		return fields
	}
	existingKeys := map[string]bool{}
	for _, f := range existing {
		existingKeys[f.Key] = true
	}
	for _, f := range fields {
		if !existingKeys[f.Key] {
			existing = append(existing, f)
		}
	}
	sort.Slice(existing, func(i, j int) bool { return existing[i].Key < existing[j].Key })
	g.structDefs[structType] = existing
	return existing
}

func (g *nameGen) unique(base string) string {
	if !g.used[base] {
		g.used[base] = true
		return base
	}
	for i := 2; ; i++ {
		c := fmt.Sprintf("%s%d", base, i)
		if !g.used[c] {
			g.used[c] = true
			return c
		}
	}
}

type fieldInfo struct {
	Name       string // Go field name (PascalCase)
	Key        string // original config key
	GoType     string // actual Go type for comment
	AtomicType string // atomic field type, e.g. "atomic.Int64", "atomic.Value"
	StructType string // struct type name (for nested structs)
	IsStruct   bool   // nested struct
	IsList     bool   // slice
	Children   []fieldInfo
}

func inferFields(data map[string]any, ng *nameGen, filePrefix string) []fieldInfo {
	return inferFieldsWithPrefix(data, ng, "", filePrefix)
}

func inferFieldsWithPrefix(data map[string]any, ng *nameGen, prefix, filePrefix string) []fieldInfo {
	var fields []fieldInfo
	for k, v := range data {
		f := fieldInfo{
			Name: toPascal(k),
			Key:  k,
		}
		switch val := v.(type) {
		case nil:
			f.GoType = "string"
			f.AtomicType = "atomic.Value"
		case string:
			f.GoType = "string"
			f.AtomicType = "atomic.Value"
		case int, int64:
			f.GoType = "int64"
			f.AtomicType = "atomic.Int64"
		case float64:
			f.GoType = "float64"
			f.AtomicType = "atomic.Value"
		case bool:
			f.GoType = "bool"
			f.AtomicType = "atomic.Bool"
		case map[string]any:
			f.IsStruct = true
			pascal := f.Name
			var childPrefix string
			if strings.HasPrefix(k, "_") {
				pascal = toPascal(strings.TrimPrefix(k, "_"))
				f.StructType = filePrefix + pascal + "Cfg"
				childPrefix = filePrefix + pascal
				children := inferFieldsWithPrefix(val, ng, childPrefix, filePrefix)
				f.Children = ng.mergeStruct(f.StructType, children)
			} else {
				f.StructType = ng.unique(prefix + pascal + "Cfg")
				childPrefix = prefix + pascal
				f.Children = inferFieldsWithPrefix(val, ng, childPrefix, filePrefix)
			}
			f.GoType = f.StructType
			f.AtomicType = f.StructType
		case []any:
			f.IsList = true
			elemType := inferSliceElementType(val)
			f.GoType = "[]" + elemType
			f.AtomicType = "atomic.Value"
		default:
			f.GoType = "any"
			f.AtomicType = "atomic.Value"
		}
		fields = append(fields, f)
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].Key < fields[j].Key })
	return fields
}

func inferSliceElementType(slice []any) string {
	if len(slice) == 0 {
		return "any"
	}
	types := map[string]int{}
	for _, v := range slice {
		switch v.(type) {
		case string:
			types["string"]++
		case int, int64:
			types["int64"]++
		case float64:
			types["float64"]++
		case bool:
			types["bool"]++
		default:
			types["any"]++
		}
	}
	if len(types) == 1 {
		for t := range types {
			return t
		}
	}
	return "any"
}
