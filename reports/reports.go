package reports

import "sacco/database"

type Reports struct {
	DB *database.Database
}

func NewReports(dbname string) *Reports {
	return &Reports{
		DB: database.NewDatabase(dbname),
	}
}
