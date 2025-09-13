package utils

import "encoding/json"

func Map2Table(data map[string]any) string {
	var result string

	payload, _ := json.MarshalIndent(data, "", "  ")

	result = string(payload)

	return result
}
