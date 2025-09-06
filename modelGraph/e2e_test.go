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

	records := map[string]int64{}

	addRecord := func(model string, data map[string]any) (*int64, error) {
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
			t.Fatal(err)
		}

		records[model] = parentId

		fmt.Println(model, parentId)

		if len(data) > 0 {

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
}
