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

func TestLoansReport(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)

	sampleScript, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.sql"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.DB.Exec(string(sampleScript))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rpt := reports.NewReports(db)

	result, err := rpt.LoansReport("2025-11-01")
	if err != nil {
		t.Fatal(err)
	}

	removeIds := func(result *reports.LoansReportData) {
		for i, row := range result.Data {
			newRow := reports.LoansReportRow{
				LastName:      row.LastName,
				FirstName:     row.FirstName,
				LoanAmount:    row.LoanAmount,
				LoanStartDate: row.LoanStartDate,
				LoanDueDate:   row.LoanDueDate,
				BalanceAmount: row.BalanceAmount,
			}

			result.Data[i] = newRow
		}
	}

	fixturesFile := filepath.Join(".", "fixtures", "loansReport.data.json")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(fixturesFile, payload, 0644)
	}

	target := reports.LoansReportData{}

	targetContent, err := os.ReadFile(fixturesFile)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(targetContent, &target)
	if err != nil {
		t.Fatal(err)
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
