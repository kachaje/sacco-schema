package modelgraph_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	modelgraph "sacco/modelGraph"
	"sacco/utils"
	"sort"
	"testing"
)

func TestCreateGraph(t *testing.T) {
	data := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "schema", "models", "modelsData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result, err := modelgraph.CreateGraph(data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "graph.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}

	models := map[string]bool{}
	for key, vp := range result {
		models[key] = true
		if vc, ok := vp.(map[string]any); ok {
			for _, v := range vc {
				if vi, ok := v.([]any); ok {
					for _, k := range vi {
						models[k.(string)] = true
					}
				} else if vi, ok := v.([]string); ok {
					for _, k := range vi {
						models[k] = true
					}
				}
			}
		}
	}

	keys := []string{}

	for key := range models {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	dataCount := len(keys)

	keys = []string{}

	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	refCount := len(keys)

	if refCount != dataCount {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", refCount, dataCount)
	}
}

func TestCreateModelQuerySimple(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"fields": map[string]any{
				"id": map[string]any{
					"autoIncrement": true,
					"order":         0,
					"primaryKey":    true,
					"type":          "int",
				},
				"name": map[string]any{
					"order": 3,
					"type":  "text",
				},
				"password": map[string]any{
					"order": 2,
					"type":  "text",
				},
				"userRole": map[string]any{
					"order": 4,
					"type":  "text",
				},
				"username": map[string]any{
					"order":  1,
					"type":   "text",
					"unique": true,
				},
			},
			"model":   "user",
			"parents": []map[string]any{},
		},
	}

	seed := map[string]any{
		"username": "sample",
		"password": "123456789",
		"name":     "Sample User",
		"userRole": "default",
	}

	result, err := modelgraph.CreateModelQuery("user", data, seed)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	target := fmt.Sprintf(`
INSERT INTO 
	user (username, password, name, userRole) 
VALUES 
	("%v", "%v", "%v", "%v");`,
		seed["username"], seed["password"], seed["name"], seed["userRole"],
	)

	if utils.CleanString(*result) != utils.CleanString(target) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, *result)
	}
}
