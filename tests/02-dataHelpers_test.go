package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/kachaje/sacco-schema/utils"
)

func TestFlattenMapIdMapOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data, true)

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "flatIdsMap.json"), payload, 0644)
	}

	target := map[string]any{}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "flatIdsMap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{
		"memberDependantId",
		"memberLoanWitnessId",
		"memberLoanSecurityId",
		"memberLoanLiabilityId",
		"memberLoanPaymentScheduleId",
		"memberLoanPaymentDetailId",
		"memberLoanSettlementId",
		"memberLoanPaymentId",
		"memberContributionScheduleId",
		"memberContributionDepositId",
		"memberLoanPaymentSchedule",
		"memberLoanTaxId",
	} {
		delete(target, key)
		delete(result, key)
	}

	if !reflect.DeepEqual(target, result) {
		diff := utils.GetMapDiff(target, result)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf("Test failed; Diff: %s", payload)
	}
}

func TestFlattenMapAllData(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data, false)

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"), payload, 0644)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	target, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatalf("Test failed; Expected: %s; Actual: %s", target, payload)
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
		"memberNomineeId":                    "member.memberDependant.id",
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
			"memberDependant": map[string]any{
				"id": 9,
			},
		},
	}

	if !reflect.DeepEqual(target, data) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, data)
	}
}
