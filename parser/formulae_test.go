package parser_test

import (
	"encoding/json"
	"reflect"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestGetTokens(t *testing.T) {
	target := map[string]any{
		"op": "SUM",
		"terms": []any{
			"totalCostOfGoods",
			"employeesWages",
			"ownSalary",
			"transport",
			"loanInterest",
			"utilities",
			"rentals",
			"otherCosts",
		},
	}

	result := parser.GetTokens("SUM({{totalCostOfGoods}}, {{employeesWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})")

	if reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestResultFromFormulaeSUM(t *testing.T) {
	tokens := map[string]any{
		"op": "SUM",
		"terms": []any{
			"totalCostOfGoods",
			"employeesWages",
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
		"employeesWages":   "500000",
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
		"employeesWages":   "500000",
		"ownSalary":        "100000",
		"transport":        "50000",
		"loanInterest":     "0",
		"utilities":        "35000",
		"rentals":          "50000",
		"otherCosts":       "0",
		"totalIncome":      "2000000",
	}

	wf.FormulaFields["totalCosts"] = "SUM({{totalCostOfGoods}}, {{employeesWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})"
	wf.FormulaFields["netProfitLoss"] = "DIFF({{totalIncome}},{{totalCostOfGoods}}, {{employeesWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})"

	wait := make(chan bool, 1)

	wf.CalculateFormulae(wait)

	target := map[string]any{
		"employeesWages":   "500000",
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
