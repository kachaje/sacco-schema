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

func TestLoadModelChildren(t *testing.T) {
	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.SkipFields = append(db.SkipFields, []string{"createdAt", "updatedAt"}...)

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

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestFullMemberRecord(t *testing.T) {
	t.Skip()

	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	phoneNumber := "09999999999"

	result, err := db.FullMemberRecord(phoneNumber)
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	target, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}
