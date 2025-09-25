package cronjobs

import (
	"fmt"
	"log"
	"sacco/database"
)

type CronJobs struct {
	Jobs map[string]func(targetDate string) error
	DB   *database.Database
}

func NewCronJobs(db *database.Database) *CronJobs {
	jobs := &CronJobs{
		DB:   db,
		Jobs: map[string]func(targetDate string) error{},
	}

	jobs.Jobs["ordinaryDepositsInterest"] = jobs.CalculateOrdinaryDepositsInterest

	return jobs
}

func (c *CronJobs) RunCronJobs(targetDate string) error {
	for job, jobFn := range c.Jobs {
		log.Printf("Running job %s\n", job)

		if err := jobFn(targetDate); err != nil {
			return err
		}
	}

	return nil
}

func (c *CronJobs) CalculateOrdinaryDepositsInterest(targetDate string) error {
	_, err := c.DB.SQLQuery(fmt.Sprintf(`
INSERT INTO memberSavingInterest (memberSavingId, description, amount, dueDate)
WITH RECURSIVE savings AS ( SELECT 
	memberSavingId, 
	STRFTIME('%%Y', transactionDate) transactionYear,
	CONCAT(STRFTIME('%%Y', transactionDate), ' - ',
	CASE
        WHEN CAST(STRFTIME('%%m', transactionDate) AS INTEGER) BETWEEN 1 AND 3 THEN 'Q1'
        WHEN CAST(STRFTIME('%%m', transactionDate) AS INTEGER) BETWEEN 4 AND 6 THEN 'Q2'
        WHEN CAST(STRFTIME('%%m', transactionDate) AS INTEGER) BETWEEN 7 AND 9 THEN 'Q3'
        ELSE 'Q4'
    END) AS description,
	SUM(balance)/COUNT(id) AS average,
	(SUM(balance)/COUNT(id)) * 0.1 * 0.25 AS interest
FROM memberSavingTransaction 
WHERE savingsTypeName = 'Ordinary Deposit'
GROUP BY transactionYear, description, memberSavingId
) 
SELECT memberSavingId, description, interest, CURRENT_TIMESTAMP 
FROM savings 
WHERE description = CONCAT(STRFTIME('%%Y', '%s'), ' - ',
	CASE
        WHEN CAST(STRFTIME('%%m', '%s') AS INTEGER) BETWEEN 1 AND 3 THEN 'Q1'
        WHEN CAST(STRFTIME('%%m', '%s') AS INTEGER) BETWEEN 4 AND 6 THEN 'Q2'
        WHEN CAST(STRFTIME('%%m', '%s') AS INTEGER) BETWEEN 7 AND 9 THEN 'Q3'
        ELSE 'Q4'
  END)
	`, targetDate, targetDate, targetDate, targetDate))
	if err != nil {
		return err
	}

	return nil
}
