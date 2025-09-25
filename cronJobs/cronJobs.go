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
	SUM(t.balance)/COUNT(t.id) AS average,
	(SUM(t.balance)/COUNT(t.id)) * COALESCE(
				(
					SELECT
						interestRate
					FROM
						savingsType
					WHERE
						savingsTypeName = s.savingsTypeName
						AND active = 1
				),
				0
			) / 4 AS interest
FROM memberSavingTransaction t
LEFT OUTER JOIN memberSaving s ON s.id = t.memberSavingId
WHERE t.savingsTypeName = 'Ordinary Deposit'
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

func (c *CronJobs) CalculateFixedDepositInterests(targetDate string) error {
	_, err := c.DB.SQLQuery(fmt.Sprintf(`
INSERT INTO memberSavingInterest (memberSavingId, description, amount, dueDate)
WITH
	savings AS (
		SELECT
			t.savingsTypeName,
			memberSavingId,
			STRFTIME ('%%Y', transactionDate) transactionYear,
			CONCAT (
				s.savingsTypeName, ' (',
				STRFTIME ('%%Y', transactionDate),
				'/',
				STRFTIME ('%%m', transactionDate),
				')'
			) AS description,
			SUM(t.balance) / COUNT(t.id) AS average,
			(SUM(t.balance) / COUNT(t.id)) * COALESCE(
				(
					SELECT
						interestRate
					FROM
						savingsType
					WHERE
						savingsTypeName = s.savingsTypeName
						AND active = 1
				),
				0
			) / 12 AS interest
		FROM
			memberSavingTransaction t
			LEFT OUTER JOIN memberSaving s ON s.id = t.memberSavingId
		WHERE
			t.savingsTypeName IN ('Fixed Deposit', '30 day Call Deposit')
		GROUP BY
			transactionYear,
			description,
			memberSavingId
	)
SELECT
	memberSavingId, description, interest, CURRENT_TIMESTAMP
FROM
	savings 
WHERE description = CONCAT (
				savingsTypeName, ' (',
				STRFTIME ('%%Y', '%s'),
				'/',
				STRFTIME ('%%m', '%s'),
				')'
			)
	`, targetDate, targetDate))
	if err != nil {
		return err
	}

	return nil
}
