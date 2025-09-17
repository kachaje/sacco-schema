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

							if v1["memberLoanApproval"] != nil {
								if vl, ok := v1["memberLoanApproval"].(map[string]any); ok {
									if vl["loanNumber"] != nil {
										delete(vl, "loanNumber")
									}
									if vl["dateOfApproval"] != nil {
										delete(vl, "dateOfApproval")
									}
									if vl["memberLoanVerification"] != nil {
										delete(vl, "memberLoanVerification")
									}
								}
							}

							if v1["memberLoanDisbursement"] != nil {
								if vl, ok := v1["memberLoanDisbursement"].(map[string]any); ok {
									if vl["date"] != nil {
										delete(vl, "date")
									}
									if vl["description"] != nil {
										delete(vl, "description")
									}
								}
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

	delete(result, "dateJoined")
	delete(result, "memberIdNumber")

	delete(target["member"].(map[string]any), "dateJoined")
	delete(target["member"].(map[string]any), "memberIdNumber")

	deleteLoanNumber(map[string]any{
		"member": result,
	})
	deleteLoanNumber(target)

	if !utils.MapsEqual(target["member"].(map[string]any), result) {
		diff := utils.GetMapDiff(target["member"].(map[string]any), result)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf("Test failed; Diff: %s", payload)
	}
}

func TestFullMemberRecord(t *testing.T) {
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

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	delete(result["member"].(map[string]any), "dateJoined")
	delete(result["member"].(map[string]any), "memberIdNumber")

	delete(target["member"].(map[string]any), "dateJoined")
	delete(target["member"].(map[string]any), "memberIdNumber")

	deleteLoanNumber(result)
	deleteLoanNumber(target)

	if !utils.MapsEqual(target, result) {
		diff := utils.GetMapDiff(target, result)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf("Test failed; Diff: %s", payload)
	}
}
