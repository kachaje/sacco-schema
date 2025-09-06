package utils_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sacco/utils"
	"sort"
	"testing"
)

func TestFlattenMapIdMapOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data, true)

	target := map[string]any{
		"memberBeneficiaryId": map[string]any{
			"key":   "member.memberBeneficiary.0.id",
			"value": "1",
		},
		"memberBusinessId": map[string]any{
			"key":   "member.memberLoan.0.memberBusiness.id",
			"value": "1",
		},
		"memberContactId": map[string]any{
			"key":   "member.memberContact.id",
			"value": "1",
		},
		"memberId": map[string]any{
			"key":   "member.id",
			"value": "1",
		},
		"memberLastYearBusinessHistoryId": map[string]any{
			"key":   "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id",
			"value": "1",
		},
		"memberLoanApprovalId": map[string]any{
			"key":   "member.memberLoan.0.memberLoanApproval.id",
			"value": "1",
		},
		"memberLoanId": map[string]any{
			"key":   "member.memberLoan.0.id",
			"value": "1",
		},
		"memberNextYearBusinessProjectionId": map[string]any{
			"key":   "member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id",
			"value": "1",
		},
		"memberNomineeId": map[string]any{
			"key":   "member.memberNominee.id",
			"value": "1",
		},
		"memberOccupationId": map[string]any{
			"key":   "member.memberLoan.0.memberOccupation.id",
			"value": "1",
		},
		"memberOccupationVerificationId": map[string]any{
			"key":   "member.memberLoan.0.memberOccupation.memberOccupationVerification.id",
			"value": "1",
		},
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestFlattenMapAllData(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data, false)

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	if utils.CleanScript(payload) != utils.CleanScript(content) {
		t.Fatal("Test failed")
	}
}

func TestSetNestedValue(t *testing.T) {
	rawData := map[string]any{
		"memberBeneficiaryId":                "member.memberBeneficiary.0.id",
		"memberBusinessId":                   "member.memberLoan.0.memberBusiness.id",
		"memberContactId":                    "member.memberContact.id",
		"memberId":                           "member.id",
		"memberLastYearBusinessHistoryId":    "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id",
		"memberLoanApprovalId":               "member.memberLoan.0.memberLoanApproval.id",
		"memberLoanId":                       "member.memberLoan.0.id",
		"memberNextYearBusinessProjectionId": "member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id",
		"memberNomineeId":                    "member.memberNominee.id",
		"memberOccupationId":                 "member.memberLoan.0.memberOccupation.id",
		"memberOccupationVerificationId":     "member.memberLoan.0.memberOccupation.memberOccupationVerification.id",
	}

	keys := []string{}

	for key := range rawData {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	data := map[string]any{}

	for i, key := range keys {
		value := rawData[key]

		utils.SetNestedValue(data, value.(string), i+1)
	}

	target := map[string]any{
		"member": map[string]any{
			"id": 4,
			"memberBeneficiary": map[string]any{
				"0": map[string]any{
					"id": 1,
				},
			},
			"memberContact": map[string]any{
				"id": 3,
			},
			"memberLoan": map[string]any{
				"0": map[string]any{
					"id": 7,
					"memberBusiness": map[string]any{
						"id": 2,
						"memberLastYearBusinessHistory": map[string]any{
							"0": map[string]any{
								"id": 5,
							},
						},
						"memberNextYearBusinessProjection": map[string]any{
							"0": map[string]any{
								"id": 8,
							},
						},
					},
					"memberLoanApproval": map[string]any{
						"id": 6,
					},
					"memberOccupation": map[string]any{
						"id": 10,
						"memberOccupationVerification": map[string]any{
							"id": 11,
						},
					},
				},
			},
			"memberNominee": map[string]any{
				"id": 9,
			},
		},
	}

	if !reflect.DeepEqual(target, data) {
		t.Fatal("Test failed")
	}
}
