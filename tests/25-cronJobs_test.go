package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	cronjobs "sacco/cronJobs"
	"sacco/database"
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

	result, err := db.GenericModels["memberSavingInterest"].FilterBy("WHERE active = 1")
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}
