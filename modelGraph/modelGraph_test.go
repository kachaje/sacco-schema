package modelgraph_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	modelgraph "sacco/modelGraph"
	"testing"
)

func TestCreateGraph(t *testing.T) {
	data := map[string]any{}

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

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(payload))
}
