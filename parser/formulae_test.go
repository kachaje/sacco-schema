package parser_test

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sacco/parser"
	"sacco/utils"
	"testing"
	"time"
)

func TestGetTokens(t *testing.T) {
	target := map[string]any{
		"op": "SUM",
		"terms": []any{
			"totalCostOfGoods",
			"employeeWages",
			"ownSalary",
			"transport",
			"loanInterest",
			"utilities",
			"rentals",
			"otherCosts",
		},
	}

	result := parser.GetTokens("SUM({{totalCostOfGoods}}, {{employeeWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})")

	if reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestResultFromFormulaeSUM(t *testing.T) {
	tokens := map[string]any{
		"op": "SUM",
		"terms": []any{
			"totalCostOfGoods",
			"employeeWages",
			"ownSalary",
			"transport",
			"loanInterest",
			"utilities",
			"rentals",
			"otherCosts",
		},
	}
	data := map[string]any{
		"totalCostOfGoods": "1000000",
		"employeeWages":    "500000",
		"ownSalary":        "100000",
		"transport":        "50000",
		"loanInterest":     "0",
		"utilities":        "35000",
		"rentals":          "50000",
		"otherCosts":       "0",
	}

	result, err := parser.ResultFromFormulae(tokens, data)
	if err != nil {
		t.Fatal(err)
	}

	target := 1735000.00

	if *result != target {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target, *result)
	}
}

func TestResultFromFormulaeDIFF(t *testing.T) {
	tokens := map[string]any{
		"op": "DIFF",
		"terms": []any{
			"totalIncome",
			"totalCostOfGoods",
		},
	}
	data := map[string]any{
		"totalCostOfGoods": "1735000.00",
		"totalIncome":      "2000000",
	}

	result, err := parser.ResultFromFormulae(tokens, data)
	if err != nil {
		t.Fatal(err)
	}

	target := 265000.0

	if *result != target {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target, *result)
	}
}

func TestCalculateFormulae(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	wf.Data = map[string]any{
		"totalCostOfGoods": "1000000",
		"employeeWages":    "500000",
		"ownSalary":        "100000",
		"transport":        "50000",
		"loanInterest":     "0",
		"utilities":        "35000",
		"rentals":          "50000",
		"otherCosts":       "0",
		"totalIncome":      "2000000",
	}

	wf.FormulaFields["totalCosts"] = "SUM({{totalCostOfGoods}}, {{employeeWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})"
	wf.FormulaFields["netProfitLoss"] = "DIFF({{totalIncome}},{{totalCostOfGoods}}, {{employeeWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})"

	wait := make(chan bool, 1)

	wf.CalculateFormulae(wait)

	target := map[string]any{
		"employeeWages":    "500000",
		"loanInterest":     "0",
		"netProfitLoss":    "265000.00",
		"otherCosts":       "0",
		"ownSalary":        "100000",
		"rentals":          "50000",
		"totalCostOfGoods": "1000000",
		"totalCosts":       "1735000.00",
		"totalIncome":      "2000000",
		"transport":        "50000",
		"utilities":        "35000",
	}

	if !reflect.DeepEqual(wf.Data, target) {
		diff := utils.GetMapDiff(wf.Data, target)

		payload, _ := json.MarshalIndent(diff, "", "  ")

		t.Fatalf("Test failed. Diff: %s\n", payload)
	}
}

func TestDATE_DIFF_YEARS(t *testing.T) {
	tokens := parser.GetTokens("DATE_DIFF_YEARS({{TODAY}}-{{dateOfBirth}})")

	data := map[string]any{
		"startDate": "1999-09-01",
		"refDate":   "2025-09-11",
	}

	tm, err := time.Parse("2006-01-02", "1999-09-01")
	if err != nil {
		t.Fatal(err)
	}

	d := time.Since(tm)

	target := math.Round(d.Abs().Hours() / (365 * 24))

	result, err := parser.ResultFromFormulae(tokens, data)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	if *result != target {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target, *result)
	}
}

func TestDIV(t *testing.T) {
	tokens := parser.GetTokens("DIV({{loanAmount}},{{repaymentPeriodInMonths}})")

	data := map[string]any{
		"repaymentPeriodInMonths1": 12,
		"loanAmount":               120000,
	}

	result, err := parser.ResultFromFormulae(tokens, data)
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	var target float64 = 10000

	if *result != target {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target, *result)
	}
}

func TestGetScheduleParams(t *testing.T) {
	tokens := parser.GetScheduleParams("REDUCING_SCHEDULE({{loanAmount}},{{repaymentPeriodInMonths}},[{{processingFeeRate}}],[{{monthlyInterestRate}},{{monthlyInsuranceRate}}])")

	target := map[string]any{
		"amount":   "loanAmount",
		"duration": "repaymentPeriodInMonths",
		"oneTimeRates": []string{
			"processingFeeRate",
		},
		"op": "REDUCING_SCHEDULE",
		"recurringRates": []string{
			"monthlyInterestRate",
			"monthlyInsuranceRate",
		},
	}

	if !reflect.DeepEqual(tokens, target) {
		t.Fatal("Test failed")
	}
}

func TestGenerateSchedule(t *testing.T) {
	tokens := parser.GetScheduleParams("REDUCING_SCHEDULE({{loanAmount}},{{repaymentPeriodInMonths}},[{{processingFeeRate}}],[{{monthlyInterestRate}},{{monthlyInsuranceRate}}])")

	data := map[string]any{
		"loanAmount":              200000,
		"repaymentPeriodInMonths": 6,
		"processingFeeRate":       0.05,
		"monthlyInterestRate":     0.05,
		"monthlyInsuranceRate":    0.15,
	}

	_ = data

	// result, err := parser.GenerateSchedule(tokens, data)

	fmt.Printf("%#v\n", tokens)
}
