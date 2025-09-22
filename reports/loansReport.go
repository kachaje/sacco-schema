package reports

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LoansReportRow struct {
	MemberIdNumber string  `json:"memberIdNumber"`
	LoanNumber     string  `json:"loanNumber"`
	LastName       string  `json:"lastName"`
	FirstName      string  `json:"firstName"`
	LoanAmount     float64 `json:"loanAmount"`
	LoanStartDate  string  `json:"loanStartDate"`
	LoanDueDate    string  `json:"loanDueDate"`
	BalanceAmount  float64 `json:"balanceAmount"`
}

type LoansReportData struct {
	TargetDate         string                    `json:"targetDate"`
	Data               map[string]LoansReportRow `json:"data"`
	TotalLoanAmount    float64                   `json:"totalLoanAmount"`
	TotalBalanceAmount float64                   `json:"totalBalanceAmount"`
}

func (r *Reports) LoansReport(targetDate string) (*LoansReportData, error) {
	report := LoansReportData{
		TargetDate:         targetDate,
		Data:               map[string]LoansReportRow{},
		TotalLoanAmount:    0,
		TotalBalanceAmount: 0,
	}

	rawData, err := r.DB.SQLQuery(fmt.Sprintf(`
WITH schedule AS (
	SELECT memberLoanId, amountRecommended, SUM(instalment) AS paid FROM memberLoanPaymentSchedule WHERE DATE(dueDate) <= DATE('%s') AND COALESCE(amountPaid, 0) > 0 GROUP BY memberLoanId
)
SELECT firstName, lastName, memberIdNumber, l.loanNumber, 
COALESCE(s.amountRecommended, 0) AS loanAmount, loanStartDate, loanDueDate, 
COALESCE((s.amountRecommended - COALESCE(s.paid, 0)), 0)  AS balanceAmount
FROM memberLoan l
LEFT OUTER JOIN member m ON m.id = l.memberId
LEFT OUTER JOIN schedule s ON s.memberLoanId = l.id`,
		targetDate))
	if err != nil {
		return nil, err
	}

	for i, data := range rawData {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		reportRow := LoansReportRow{}

		err = json.Unmarshal(jsonBytes, &reportRow)
		if err != nil {
			return nil, err
		}

		report.TotalBalanceAmount += reportRow.BalanceAmount
		report.TotalLoanAmount += reportRow.LoanAmount

		report.Data[fmt.Sprint(i+1)] = reportRow
	}

	return &report, nil
}

func (r *Reports) LoansReport2Table(data LoansReportData) (*string, error) {
	targetDate, err := time.Parse("2006-01-02", data.TargetDate)
	if err != nil {
		return nil, err
	}

	lines := []string{
		"Loans Report",
		targetDate.Format("02-Jan-2006"),
	}

	result := strings.Join(lines, "\n")

	return &result, nil
}
