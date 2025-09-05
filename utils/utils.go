package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type DiffResult struct {
	Added   map[string]any
	Removed map[string]any
	Changed map[string]any
}

func CleanScript(content []byte) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func CleanString(content string) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func DumpYaml(data map[string]any) (*string, error) {
	var result string

	payload, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	result = string(payload)

	return &result, nil
}

func LoadYaml(yamlData string) (map[string]any, error) {
	var data map[string]any

	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MapsEqual(m1, m2 map[string]any) bool {
	if (m1 == nil) != (m2 == nil) {
		return false
	}
	if m1 == nil && m2 == nil {
		return true
	}
	if len(m1) != len(m2) {
		return false
	}

	for key, val1 := range m1 {
		val2, ok := m2[key]
		if !ok {
			return false
		}

		switch v1 := val1.(type) {
		case map[string]any:
			v2, ok := val2.(map[string]any)
			if !ok || !MapsEqual(v1, v2) {
				return false
			}
		default:
			if fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
				return false
			}
		}
	}
	return true
}

func GetMapDiff(map1, map2 map[string]any) DiffResult {
	diff := DiffResult{
		Added:   make(map[string]any),
		Removed: make(map[string]any),
		Changed: make(map[string]any),
	}

	for key, val1 := range map1 {
		if val2, ok := map2[key]; !ok {
			diff.Removed[key] = val1
		} else {
			if nestedMap1, isMap1 := val1.(map[string]any); isMap1 {
				if nestedMap2, isMap2 := val2.(map[string]any); isMap2 {
					nestedDiff := GetMapDiff(nestedMap1, nestedMap2)
					if len(nestedDiff.Added) > 0 || len(nestedDiff.Removed) > 0 || len(nestedDiff.Changed) > 0 {
						diff.Changed[key] = nestedDiff
					}
				} else {
					diff.Changed[key] = map[string]any{
						"old":     val1,
						"new":     val2,
						"oldType": reflect.TypeOf(val1).String(),
						"newType": reflect.TypeOf(val2).String(),
					}
				}
			} else if !reflect.DeepEqual(val1, val2) {
				diff.Changed[key] = map[string]any{
					"old":     val1,
					"new":     val2,
					"oldType": reflect.TypeOf(val1).String(),
					"newType": reflect.TypeOf(val2).String(),
				}
			}
		}
	}

	for key, val2 := range map2 {
		if _, ok := map1[key]; !ok {
			diff.Added[key] = val2
		}
	}

	return diff
}
