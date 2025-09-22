package reports

type LoansReportRow struct {
	MemberId      string  `json:"memberId"`
	LastName      string  `json:"lastName"`
	FirstName     string  `json:"firstName"`
	LoanAmount    float64 `json:"loanAmount"`
	LoanStartDate string  `json:"loanStartDate"`
	LoanDueDate   string  `json:"loanDueDate"`
	BalanceAmount float64 `json:"balanceAmount"`
}

type LoansReportData struct {
	Data               []LoansReportRow `json:"data"`
	TotalLoanAmount    float64          `json:"totalLoanAmount"`
	TotalBalanceAmount float64          `json:"totalBalanceAmount"`
}

func (r *Reports) LoansReport(targetDate string) LoansReportData {
	var report LoansReportData

	return report
}
