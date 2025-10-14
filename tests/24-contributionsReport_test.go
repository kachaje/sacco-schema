package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kachaje/sacco-schema/database"
	"github.com/kachaje/sacco-schema/reports"
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

func TestContributionsReport2Table(t *testing.T) {
	reportData := reports.ContributionReportData{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "contributions.data.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &reportData)
	if err != nil {
		t.Fatal(err)
	}

	rpt := reports.Reports{}

	result, err := rpt.ContributionsReport2Table(reportData, nil)
	if err != nil {
		t.Fatal(err)
	}

	fixturesFile := filepath.Join(".", "fixtures", "contributions.data.txt")

	if os.Getenv("DEBUG") == "true" {
		payload := []byte(string(*result))

		os.WriteFile(fixturesFile, payload, 0644)
	}

	content, err = os.ReadFile(fixturesFile)
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if target != *result {
		t.Fatalf(`Test failed.
Expected:
%v
Actual:
%v`, target, *result)
	}
}
