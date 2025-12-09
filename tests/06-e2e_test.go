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

	records := map[string][]int64{}

	var addRecord func(model string, data map[string]any, parentModel string, parentId int64) (*int64, error)

	addRecord = func(model string, data map[string]any, parentModel string, parentId int64) (*int64, error) {
		if len(data) > 0 {
			id, err := dbInstance.GenericModels[model].AddRecord(data)
			if err != nil {
				return nil, err
			}
			if records[model] == nil {
				records[model] = []int64{}
			}
			records[model] = append(records[model], *id)
			return id, nil
		}

		query, err := modelgraph.CreateModelQuery(model, modelsData, nil)
		if err != nil {
			t.Fatal(err)
		}

		result, err := db.Exec(*query)
		if err != nil {
			t.Fatal(err)
		}

		recordId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		if records[model] == nil {
			records[model] = []int64{}
		}
		records[model] = append(records[model], recordId)

		if graphVal, ok := graphData[model].(map[string]any); ok {
			var childModels []string
			if singleChildren, ok := graphVal["singleChildren"].([]any); ok {
				for _, v := range singleChildren {
					childModels = append(childModels, fmt.Sprintf("%v", v))
				}
			} else if singleChildren, ok := graphVal["singleChildren"].([]string); ok {
				childModels = append(childModels, singleChildren...)
			}
			if arrayChildren, ok := graphVal["arrayChildren"].([]any); ok {
				for _, v := range arrayChildren {
					childModels = append(childModels, fmt.Sprintf("%v", v))
				}
			} else if arrayChildren, ok := graphVal["arrayChildren"].([]string); ok {
				childModels = append(childModels, arrayChildren...)
			}

			for _, key := range childModels {
				childData := map[string]any{}
				if modelData, ok := modelsData[key].(map[string]any); ok {
					if fields, ok := modelData["fields"].(map[string]any); ok {
						for fieldName, fieldDef := range fields {
							if field, ok := fieldDef.(map[string]any); ok {
								if refTable, ok := field["referenceTable"].(string); ok && refTable == model {
									childData[fieldName] = recordId
								} else if field["autoIncrement"] == nil && field["optional"] == nil && field["default"] == nil && field["dynamicDefault"] == nil {
									if options, ok := field["options"].([]any); ok && len(options) > 0 {
										childData[fieldName] = fmt.Sprintf("%v", options[0])
									} else if options, ok := field["options"].([]string); ok && len(options) > 0 {
										childData[fieldName] = options[0]
									} else if fieldType, ok := field["type"].(string); ok {
										if fieldType == "text" {
											childData[fieldName] = "test"
										} else if fieldType == "int" || fieldType == "integer" {
											childData[fieldName] = 10
										} else if fieldType == "real" || fieldType == "float" {
											childData[fieldName] = 10.0
										} else {
											childData[fieldName] = "test"
										}
									}
								} else if field["dynamicDefault"] != nil && field["optional"] == nil {
									if fieldType, ok := field["type"].(string); ok {
										if fieldType == "text" {
											if strings.Contains(strings.ToLower(fieldName), "number") {
												childData[fieldName] = fmt.Sprintf("LOAN-%d", recordId)
											} else {
												childData[fieldName] = fmt.Sprintf("test-%d", recordId)
											}
										} else if fieldType == "int" || fieldType == "integer" {
											childData[fieldName] = int(recordId)
										} else if fieldType == "real" || fieldType == "float" {
											childData[fieldName] = float64(recordId)
										} else {
											childData[fieldName] = fmt.Sprintf("test-%d", recordId)
										}
									}
								}
							}
						}
					}
				}
				_, err = addRecord(key, childData, model, recordId)
				if err != nil {
					return nil, err
				}
			}
		}

		return &recordId, nil
	}

	for model, value := range graphData {
		if strings.HasSuffix(model, "IdsCache") {
			continue
		}

		if _, ok := value.(map[string]any); ok {
			_, err := addRecord(model, map[string]any{}, "", 0)
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
