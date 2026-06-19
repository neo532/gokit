package main

import "encoding/json"

func jsonDecode(data []byte) (map[string]any, error) {
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return v, nil
}
