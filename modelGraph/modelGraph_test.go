package modelgraph_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	modelgraph "sacco/modelGraph"
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

	if false {
		payload, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Fatal(err)
		}

		os.WriteFile(filepath.Join(".", "fixtures", "graph.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "graph.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

}
