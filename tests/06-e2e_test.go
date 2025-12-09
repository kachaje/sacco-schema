package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kachaje/sacco-schema/database"
	modelgraph "github.com/kachaje/workflow-parser/modelGraph"

	_ "modernc.org/sqlite"
)

func TestSchemaE2E(t *testing.T) {
	graphData := map[string]any{}
	modelsData := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "schema", "configs", "graph.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &graphData)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join("..", "database", "schema", "configs", "modelsData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &modelsData)
	if err != nil {
		t.Fatal(err)
	}

	// Use database package which already has embedded schema files
	dbInstance := database.NewDatabase(":memory:")
	defer dbInstance.Close()

	db := dbInstance.DB

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		t.Fatal(err)
	}

	// Tables are created synchronously during database initialization,
	// so no polling is needed

	records := map[string][]int64{}

	var addRecord func(model string, data map[string]any) (*int64, error)

	addRecord = func(model string, data map[string]any) (*int64, error) {
		query, err := modelgraph.CreateModelQuery(model, modelsData, nil)
		if err != nil {
			t.Fatal(err)
		}

		result, err := db.Exec(*query)
		if err != nil {
			t.Fatal(err)
		}

		parentId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		if records[model] == nil {
			records[model] = []int64{}
		}
		records[model] = append(records[model], parentId)

		if len(data) > 0 {
			for _, value := range data {
				models := []string{}

				if val, ok := value.([]any); ok {
					for _, v := range val {
						models = append(models, fmt.Sprintf("%v", v))
					}
				} else if val, ok := value.([]string); ok {
					models = append(models, val...)
				}

				for _, key := range models {
					_, err = addRecord(key, map[string]any{})
					if err != nil {
						return nil, err
					}
				}
			}
		}

		return &parentId, nil
	}

	for model, value := range graphData {
		if strings.HasSuffix(model, "IdsCache") {
			continue
		}

		if val, ok := value.(map[string]any); ok {
			_, err := addRecord(model, val)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(records, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "records.json"), payload, 0644)
	}

	targetRecords := map[string]any{}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "records.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetRecords)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range targetRecords {
		if records[key] == nil {
			t.Fatalf("Test failed. Missing %v", key)
		}

		targetCount := 0
		if val, ok := value.([]any); ok {
			targetCount = len(val)
		} else if val, ok := value.([]string); ok {
			targetCount = len(val)
		}

		if child, ok := records[key]; ok {
			resultCount := len(child)

			if resultCount != targetCount {
				t.Fatalf("Test failed. %s Expected: %v; Actual: %v", key, targetCount, resultCount)
			}
		}
	}
}
