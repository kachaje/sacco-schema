package tests

import (
	"encoding/json"
	"fmt"
	"sacco/reports"
	"testing"
)

func TestLoansReport(t *testing.T) {
	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	reports := reports.NewReports(db)

	result, err := reports.LoansReport("2025-11-01")
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}
