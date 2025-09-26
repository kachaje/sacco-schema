package tests

import (
	"encoding/json"
	"sacco/ledger"
	"sacco/ledger/models"
	"sacco/utils"
	"testing"
)

var (
	ledgerData = map[string]any{
		"description": "Lots of groceries",
		"name":        "Albertson's transaction",
		"ledgerEntries": []map[string]any{
			{
				"referenceNumber": "1172",
				"amount":          1234,
				"accountType":     "ASSET",
				"debitCredit":     "DEBIT",
				"name":            "Some ledger entry",
			},
			{
				"referenceNumber": "1172",
				"amount":          1234,
				"accountType":     "ASSET",
				"debitCredit":     "CREDIT",
				"name":            "Some ledger entry",
			},
		},
	}
)

func TestGetAccountDirection(t *testing.T) {
	result := ledger.GetAccountDirection(models.ASSET, models.DEBIT, 1000)

	target := `balance = COALESCE(balance, 0) + 1000`

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed. Expected: '%s'; Actual: '%s'`, target, result)
	}

	result = ledger.GetAccountDirection(models.ASSET, models.CREDIT, 1000)

	target = `balance = COALESCE(balance, 0) - 1000`

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed. Expected: '%s'; Actual: '%s'`, target, result)
	}

	result = ledger.GetAccountDirection(models.LIABILITY, models.CREDIT, 1000)

	target = `balance = COALESCE(balance, 0) + 1000`

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed. Expected: '%s'; Actual: '%s'`, target, result)
	}

	result = ledger.GetAccountDirection(models.LIABILITY, models.DEBIT, 1000)

	target = `balance = COALESCE(balance, 0) - 1000`

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed. Expected: '%s'; Actual: '%s'`, target, result)
	}
}

func TestCreateEntryTransactions(t *testing.T) {
	var result []string

	saveFn := func(query string) ([]map[string]any, error) {
		result = append(result, query)

		return nil, nil
	}

	ledger.SaveHandler = saveFn

	data := ledgerData["ledgerEntries"].([]map[string]any)[0]

	entry := ledger.LedgerEntry{}

	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(payload, &entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.Description = ledgerData["description"].(string)

	err = ledger.CreateEntryTransactions(entry)
	if err != nil {
		t.Fatal(err)
	}

	target := `
INSERT INTO accountEntry (
	accountId, 
	referenceNumber, 
	name, 
	description, 
	debitCredit, 
	amount
) VALUES (
	(SELECT id FROM account WHERE accountType = 'ASSET'),
	'1172', 'Some ledger entry', 'Lots of groceries', 'DEBIT', 1234
)`

	if len(result) < 2 {
		t.Fatal("Test failed")
	}

	if utils.CleanString(target) != utils.CleanString(result[0]) {
		t.Fatalf(`Test failed.
Expected:
%s
Actual:
%s`, target, result[0])
	}

	target = `UPDATE account SET balance = COALESCE(balance, 0) + 1234 WHERE id = (SELECT id FROM account WHERE accountType = 'ASSET')`

	if utils.CleanString(target) != utils.CleanString(result[1]) {
		t.Fatalf(`Test failed.
Expected:
%s
Actual:
%s`, target, result[1])
	}
}
