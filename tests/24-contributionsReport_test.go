package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/database"
	"sacco/reports"
	"testing"
)

func TestContributions(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "contributionsSample.sql"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.DB.Exec(string(content))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rpt := reports.NewReports(db)

	result, err := rpt.ContributionsReport("2025-11-01")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}
