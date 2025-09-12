package database_test

import (
	"fmt"
	"sacco/database"
	"slices"
	"testing"
)

func TestDatabase(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	rows, err := db.DB.QueryContext(t.Context(), `SELECT name FROM sqlite_schema WHERE type="table"`)
	if err != nil {
		t.Fatal(err)
	}

	tables := []string{}

	for rows.Next() {
		var result string

		if err := rows.Scan(&result); err != nil {
			t.Fatal(err)
		}

		tables = append(tables, result)
	}

	for _, table := range []string{"account", "sqlite_sequence", "accountJournal", "accountStatement", "accountTransaction", "dividends", "insuranceProvider", "loanType", "member", "memberBusiness", "memberBusinessVerification", "memberContact", "memberDependant", "memberIdsCache", "memberLastYearBusinessHistory", "memberLoan", "memberLoanApproval", "memberLoanDisbursement", "memberLoanInsurance", "memberLoanLiability", "memberLoanPaymentSchedule", "memberLoanProcessingFee", "memberLoanReceipt", "memberLoanSecurity", "memberLoanTax", "memberLoanWitness", "memberNextYearBusinessProjection", "memberOccupation", "memberOccupationVerification", "memberSaving", "memberSavingDeposit", "memberSavingInterest", "memberSavingWithdrawal", "memberSavingsIdsCache", "memberShares", "memberSharesIdsCache", "notification", "savingsRate", "savingsType", "sharesDepositReceipt", "sharesDepositWithdraw", "user", "userRole"} {
		if !slices.Contains(tables, table) {
			t.Fatalf("Test failed. Missing: %s", table)
		}
	}
}

func TestGenericsSaveData(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	data := map[string]any{
		"yearsInBusiness": 1,
		"businessNature":  "Vendor",
		"businessName":    "Vendors Galore",
		"tradingArea":     "Mtandire",
		"memberLoanId":    1,
	}

	mid, err := db.GenericsSaveData(data, "memberBusiness", 0)
	if err != nil {
		t.Fatal(err)
	}

	if mid == nil {
		t.Fatal("Test failed. Got nil id")
	}

	{
		result, err := db.GenericModels["memberBusiness"].FetchById(*mid)
		if err != nil {
			t.Fatal(err)
		}

		if result == nil {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[key] == nil {
				t.Fatal("Test failed")
			}

			if fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", value) {
				t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[key])
			}
		}
	}
}

func TestGenericModel(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	data := map[string]any{
		"yearsInBusiness": 1,
		"businessNature":  "Vendor",
		"businessName":    "Vendors Galore",
		"tradingArea":     "Mtandire",
		"memberLoanId":    1,
	}

	mid, err := db.GenericModels["memberBusiness"].AddRecord(data)
	if err != nil {
		t.Fatal(err)
	}

	if mid == nil {
		t.Fatal("Test failed. Got nil id")
	}

	{
		result, err := db.GenericModels["memberBusiness"].FetchById(*mid)
		if err != nil {
			t.Fatal(err)
		}

		if result == nil {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[key] == nil {
				t.Fatal("Test failed")
			}

			if fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", value) {
				t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[key])
			}
		}
	}

	{
		result, err := db.GenericModels["memberBusiness"].FilterBy(`WHERE businessNature="Vendor"`)
		if err != nil {
			t.Fatal(err)
		}

		if len(result) <= 0 {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[0][key] == nil {
				t.Fatal("Test failed")
			}

			if fmt.Sprintf("%v", result[0][key]) != fmt.Sprintf("%v", value) {
				t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[0][key])
			}
		}
	}

	{
		err = db.GenericModels["memberBusiness"].UpdateRecord(map[string]any{
			"businessNature": "Taxi",
		}, *mid)
		if err != nil {
			t.Fatal(err)
		}

		result, err := db.GenericModels["memberBusiness"].FetchById(*mid)
		if err != nil {
			t.Fatal(err)
		}

		if result == nil {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[key] == nil {
				t.Fatal("Test failed")
			}

			if key == "businessNature" {
				if result[key].(string) != "Taxi" {
					t.Fatalf("Test failed. Expected: Taxi; Actual: %v", result[key])
				}
			} else {
				if fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", value) {
					t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[key])
				}
			}
		}
	}
}
