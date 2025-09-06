package modelgraph_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	modelgraph "sacco/modelGraph"
	"sacco/utils"
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

	if len(models) != len(data) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", len(data), len(models))
	}
}
