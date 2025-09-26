package models

type DebitCredit string

const (
	DEBIT  DebitCredit = "DEBIT"
	CREDIT DebitCredit = "CREDIT"
)

type AccountEntry struct {
	ID              int         `json:"id"`
	ReferenceNumber string      `json:"referenceNumber"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	DebitCredit     DebitCredit `json:"debitCredit"`
	Amount          int         `json:"amount"`
	AccountId       int         `json:"accountId"`
}
