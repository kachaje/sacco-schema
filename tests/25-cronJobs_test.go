package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	cronjobs "sacco/cronJobs"
	"sacco/database"
	"sacco/utils"
	"testing"
)

func TestCalculateOrdinaryDepositsInterest(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "savings.sql"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.DB.Exec(string(content))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	jobs := cronjobs.NewCronJobs(db)

	err = jobs.CalculateOrdinaryDepositsInterest("2025-08-30")
	if err != nil {
		t.Fatal(err)
	}

	err = jobs.CalculateOrdinaryDepositsInterest("2025-12-31")
	if err != nil {
		t.Fatal(err)
	}

	err = jobs.CalculateOrdinaryDepositsInterest("2026-03-30")
	if err != nil {
		t.Fatal(err)
	}

	err = jobs.CalculateOrdinaryDepositsInterest("2026-05-30")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.GenericModels["memberSavingInterest"].FilterBy("WHERE active = 1")
	if err != nil {
		t.Fatal(err)
	}

	result := map[string]any{}

	for _, row := range rows {
		for _, key := range []string{"createdAt", "updatedAt", "date", "dueDate"} {
			delete(row, key)
		}

		result[fmt.Sprintf("%v", row["id"])] = row
	}

	fixturesFile := filepath.Join(".", "fixtures", "savings.data.json")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(fixturesFile, payload, 0644)
	}

	target := map[string]any{}

	content, err = os.ReadFile(fixturesFile)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.MapsEqual(target, result) {
		diff := utils.GetMapDiff(target, result)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf(`Test failed.
DiffMap:
%s`, payload)
	}
}
