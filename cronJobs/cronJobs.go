package cronjobs

import (
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

	jobs.Jobs["savingsInterests"] = jobs.CalculateSavingsInterests

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

func (c *CronJobs) CalculateSavingsInterests(targetDate string) error {
	_, err := c.DB.SQLQuery(`
	
	`)
	if err != nil {
		return err
	}

	return nil
}
