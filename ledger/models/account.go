package models

type AccountType string

const (
	ASSET     AccountType = "ASSET"
	EXPENSE   AccountType = "EXPENSE"
	LIABILITY AccountType = "LIABILITY"
	EQUITY    AccountType = "EQUITY"
	REVENUE   AccountType = "REVENUE"
)

type Account struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	AccountType AccountType    `json:"accountType"`
	Entries     []AccountEntry `json:"entries"`
	Balance     int            `json:"balance"`
}
