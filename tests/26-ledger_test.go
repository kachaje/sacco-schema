package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kachaje/sacco-schema/ledger"
	"github.com/kachaje/sacco-schema/ledger/models"
	"github.com/kachaje/utils/utils"
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

	queryFn := func(query string) ([]map[string]any, error) {
		result = append(result, query)

		return nil, nil
	}

	ledger.QueryHandler = queryFn

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

func TestHandlePost(t *testing.T) {
	var result []string

	queryFn := func(query string) ([]map[string]any, error) {
		result = append(result, strings.TrimSpace(query))

		return nil, nil
	}

	ledger.QueryHandler = queryFn

	payload, err := json.Marshal(ledgerData)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	http.HandlerFunc(ledger.HandlePost).ServeHTTP(rr, req)

	fixturesFile := filepath.Join(".", "fixtures", "ledger.post.sql")

	if os.Getenv("DEBUG") == "true" {
		payload = []byte(strings.Join(result, ";\n"))

		if err := os.WriteFile(fixturesFile, payload, 0644); err != nil {
			t.Fatal(err)
		}
	}

	target := []string{}

	content, err := os.ReadFile(fixturesFile)
	if err != nil {
		t.Fatal(err)
	}

	for line := range strings.SplitSeq(string(content), ";") {
		if len(strings.TrimSpace(line)) > 0 {
			target = append(target, line)
		}
	}

	for i := range result {
		if strings.TrimSpace(target[i]) != strings.TrimSpace(result[i]) {
			t.Fatalf(`Test failed.
Expected:
%s
Actual:
%s`, target[i], result[i])
		}
	}
}

func TestHandleGet(t *testing.T) {
	var result []string

	queryFn := func(query string) ([]map[string]any, error) {
		result = append(result, strings.TrimSpace(query))

		return nil, nil
	}

	ledger.QueryHandler = queryFn

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(ledger.HandleGet).ServeHTTP(rr, req)

	req = httptest.NewRequest(
		http.MethodGet,
		"/api?startDate=2025-01-01&endDate=2025-12-31",
		nil,
	)

	http.HandlerFunc(ledger.HandleGet).ServeHTTP(rr, req)

	fixturesFile := filepath.Join(".", "fixtures", "ledger.get.sql")

	if os.Getenv("DEBUG") == "true" {
		payload := []byte(strings.Join(result, ";\n"))

		if err := os.WriteFile(fixturesFile, payload, 0644); err != nil {
			t.Fatal(err)
		}
	}

	target := []string{}

	content, err := os.ReadFile(fixturesFile)
	if err != nil {
		t.Fatal(err)
	}

	for line := range strings.SplitSeq(string(content), ";") {
		if len(strings.TrimSpace(line)) > 0 {
			target = append(target, line)
		}
	}

	for i := range result {
		if strings.TrimSpace(target[i]) != strings.TrimSpace(result[i]) {
			t.Fatalf(`Test failed.
Expected:
%s
Actual:
%s`, target[i], result[i])
		}
	}
}
