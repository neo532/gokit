package main

import "strings"

func iniDecode(data []byte) (map[string]any, error) {
	sections := parseINI(string(data))
	result := make(map[string]any)
	for section, kv := range sections {
		if section == "" {
			for k, v := range kv {
				result[k] = v
			}
		} else {
			sub := make(map[string]any)
			for k, v := range kv {
				sub[k] = v
			}
			result[section] = sub
		}
	}
	return result, nil
}

func parseINI(raw string) map[string]map[string]string {
	sections := map[string]map[string]string{}
	current := ""
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			if end := strings.Index(line, "]"); end > 0 {
				current = line[1:end]
			}
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if sections[current] == nil {
			sections[current] = map[string]string{}
		}
		sections[current][key] = val
	}
	return sections
}
