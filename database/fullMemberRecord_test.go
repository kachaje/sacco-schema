package database_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/database"
	"sacco/utils"
	"testing"

	_ "embed"
)

//go:embed fixtures/sample.sql
var sampleScript string

func setupDb() (*database.Database, error) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)

	_, err := db.DB.Exec(sampleScript)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func deleteLoanNumber(target map[string]any) {
	if target["member"] != nil {
		if vm, ok := target["member"].(map[string]any); ok {
			if vm["memberLoan"] != nil {
				if vl, ok := vm["memberLoan"].(map[string]any); ok {
					if vl["1"] != nil {
						if v1, ok := vl["1"].(map[string]any); ok {
							if v1["loanNumber"] != nil {
								delete(v1, "loanNumber")
							}
						}
					}
				}
			}
		}
	}
}

func TestLoadModelChildren(t *testing.T) {
	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SkipFields = append(db.SkipFields, []string{"createdAt", "updatedAt", "loanNumber"}...)

	result, err := db.LoadModelChildren("member", 1)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	deleteLoanNumber(target)

	if !utils.MapsEqual(target["member"].(map[string]any), result) {
		t.Fatal("Test failed")
	}
}

func TestFullMemberRecord(t *testing.T) {
	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SkipFields = append(db.SkipFields, []string{"createdAt", "updatedAt", "loanNumber"}...)

	phoneNumber := "09999999999"

	result, err := db.FullMemberRecord(phoneNumber)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	deleteLoanNumber(target)

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}
