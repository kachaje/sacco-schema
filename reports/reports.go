package reports

import "github.com/kachaje/sacco-schema/database"

type Reports struct {
	DB *database.Database
}

func NewReports(db *database.Database) *Reports {
	return &Reports{
		DB: db,
	}
}
