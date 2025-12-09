package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/sacco-schema/database"
	"github.com/kachaje/utils/utils"

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

func cleanMember(target map[string]any) {
	if target["member"] != nil {
		if vm, ok := target["member"].(map[string]any); ok {
			if vm["memberContribution"] != nil {
				delete(vm, "memberContribution")
			}

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
								if vc, ok := v1["memberLoanDisbursement"].(map[string]any); ok {
									if vc["date"] != nil {
										delete(vc, "date")
									}
									if vc["description"] != nil {
										delete(vc, "description")
									}
								}
							}

							if v1["memberLoanPayment"] != nil {
								if vc, ok := v1["memberLoanPayment"].(map[string]any); ok {
									if vc["date"] != nil {
										delete(vc, "date")
									}
									if vc["description"] != nil {
										delete(vc, "description")
									}
									if vc["loanNumber"] != nil {
										delete(vc, "loanNumber")
									}
								}
							}

							if v1["memberLoanPaymentSchedule"] != nil {
								delete(v1, "memberLoanPaymentSchedule")
							}

							for _, model := range []string{
								"memberLoanProcessingFee",
								"memberLoanTax",
							} {
								if v1[model] != nil {
									if vc, ok := v1[model].(map[string]any); ok {
										for _, v := range vc {
											if vs, ok := v.(map[string]any); ok {
												if vs["date"] != nil {
													delete(vs, "date")
												}
												if vs["description"] != nil {
													delete(vs, "description")
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
		}
	}
}

func TestLoadModelChildren(t *testing.T) {
	t.Parallel()
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

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.json"))
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

	cleanMember(map[string]any{
		"member": result,
	})
	cleanMember(target)

	if !utils.MapsEqual(target["member"].(map[string]any), result) {
		diff := utils.GetMapDiff(target["member"].(map[string]any), result)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf("Test failed; Diff: %s", payload)
	}
}

func TestFullMemberRecord(t *testing.T) {
	t.Parallel()
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

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		_ = os.WriteFile(filepath.Join(".", "fixtures", "sample.data.json"), payload, 0644)
	}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.json"))
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

	cleanMember(result)
	cleanMember(target)

	if !utils.MapsEqual(target, result) {
		diff := utils.GetMapDiff(target, result)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf("Test failed; Diff: %s", payload)
	}
}
