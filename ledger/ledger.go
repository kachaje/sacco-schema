package ledger

import (
	"fmt"
	"net/http"
	"sacco/ledger/models"

	"github.com/gorilla/mux"
)

type LedgerEntry struct {
	ReferenceNumber string             `json:"referenceNumber"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	DebitCredit     models.DebitCredit `json:"debitCredit"`
	Amount          int                `json:"amount"`
	AccountId       int                `json:"accountId"`
	AccountType     models.AccountType `json:"accountType"`
}

type TransactionBodyType struct {
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	LedgerEntries []models.AccountEntry `json:"ledgerEntries"`
}

var SaveHandler func(query string) ([]map[string]any, error)

func GetAccountDirection(accountType models.AccountType, debitCredit models.DebitCredit, amount int) string {
	switch accountType {
	case models.ASSET:
		if debitCredit == models.DEBIT {
			return fmt.Sprintf(`balance = COALESCE(balance, 0) + %v`, amount)
		} else {
			return fmt.Sprintf(`balance = COALESCE(balance, 0) - %v`, amount)
		}
	case models.LIABILITY:
		if debitCredit == models.CREDIT {
			return fmt.Sprintf(`balance = COALESCE(balance, 0) + %v`, amount)
		} else {
			return fmt.Sprintf(`balance = COALESCE(balance, 0) - %v`, amount)
		}
	default:
		return ""
	}
}

func CreateEntryTransactions(entry LedgerEntry) error {
	amount := entry.Amount
	debitCredit := entry.DebitCredit
	name := entry.Name
	accountType := entry.AccountType
	referenceNumber := entry.ReferenceNumber
	description := entry.Description

	if SaveHandler != nil {
		query := fmt.Sprintf(`
INSERT INTO accountEntry (
	accountId, 
	referenceNumber, 
	name, 
	description, 
	debitCredit, 
	amount
) VALUES (
	(SELECT id FROM account WHERE accountType = '%s'),
	'%s', '%s', '%s', '%s', %v
)`, accountType, referenceNumber, name,
			description, debitCredit, amount)

		_, err := SaveHandler(query)
		if err != nil {
			return err
		}

	}

	return nil
}

func HandlePost(w http.ResponseWriter, r *http.Request) {

}

func HandleGet(w http.ResponseWriter, r *http.Request) {

}

func ledgerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandlePost(w, r)
		return
	case http.MethodGet:
		HandleGet(w, r)
		return
	default:
		http.Error(w, "Method Not Implemented", http.StatusNotImplemented)
		return
	}
}

func Main(saveFn func(query string) ([]map[string]any, error)) *mux.Router {
	SaveHandler = saveFn

	r := mux.NewRouter()

	r.HandleFunc("/api/transaction", ledgerHandler)

	return r
}
