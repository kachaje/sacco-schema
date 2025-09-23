package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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

	result, err := rpt.ContributionsReport("2026-10-01")
	if err != nil {
		t.Fatal(err)
	}

	fixturesFile := filepath.Join(".", "fixtures", "contributions.data.json")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(fixturesFile, payload, 0644)
	}

	target := reports.ContributionReportData{}

	targetContent, err := os.ReadFile(fixturesFile)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(targetContent, &target)
	if err != nil {
		t.Fatal(err)
	}

	removeIds := func(result *reports.ContributionReportData) {
		for i, row := range result.Data {
			newRow := reports.ContributionReportRow{
				MemberName:          row.MemberName,
				MonthlyContribution: row.MonthlyContribution,
				MemberTotal:         row.MemberTotal,
				UpdatedOn:           row.UpdatedOn,
				PercentOfTotal:      row.PercentOfTotal,
			}

			result.Data[i] = newRow
		}
	}

	removeIds(result)
	removeIds(&target)

	if !reflect.DeepEqual(&target, result) {
		resultContent, _ := json.MarshalIndent(result, "", "  ")

		t.Fatalf(`Test failed.
Expected:
%s
Actual:
%s`, targetContent, resultContent)
	}
}
