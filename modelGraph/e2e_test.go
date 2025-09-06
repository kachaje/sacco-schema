package modelgraph_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	modelgraph "sacco/modelGraph"
	"slices"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestSchemaE2E(t *testing.T) {
	graphData := map[string]any{}
	modelsData := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "graph.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &graphData)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join("..", "schema", "models", "modelsData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &modelsData)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		t.Fatal(err)
	}

	for _, filename := range []string{
		filepath.Join("..", "schema", "schema.sql"),
		filepath.Join("..", "schema", "seed.sql"),
		filepath.Join("..", "schema", "triggers.sql"),
	} {
		content, err = os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		statement := string(content)

		_, err = db.Exec(statement)
		if err != nil {
			t.Fatal(err)
		}
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	tables := []string{}

	for rows.Next() {
		var table string

		err = rows.Scan(&table)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Contains([]string{"sqlite_sequence"}, table) {
			tables = append(tables, table)
		}
	}

	totalTables := 43

	for {
		if len(tables) < totalTables {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

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
				t.Fatalf("Test failed. Expected: %v; Actual: %v", targetCount, resultCount)
			}
		}
	}
}
