package reports

import "sacco/database"

type Reports struct {
	DB *database.Database
}

func NewReports(db *database.Database) *Reports {
	return &Reports{
		DB: db,
	}
}
