package main

import "gopkg.in/yaml.v3"

func yamlDecode(data []byte) (map[string]any, error) {
	var v map[string]any
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return v, nil
}
