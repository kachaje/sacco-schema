package reports

type ContributionReportRow struct {
	MemberName           string  `json:"memberName"`
	ContributionId       string  `json:"contributionId"`
	MonthlyContribution  float64 `json:"monthlyContribution"`
	TotalAmountPerMember float64 `json:"totalAmountPerMember"`
	UpdatedOn            string  `json:"updatedOn"`
	MemberId             string  `json:"memberId"`
	PercentOfTotal       float64 `json:"percentOfTotal"`
}

type ContributionReportData struct {
	TargetDate      string                           `json:"targetDate"`
	Data            map[string]ContributionReportRow `json:"data"`
	TotalAmount     float64                          `json:"totalAmount"`
	AveragePerMonth float64                          `json:"averagePerMonth"`
	MembersCount    int                              `json:"membersCount"`
}

func (r *Reports) ContributionsReport(targetDate string) (*ContributionReportData, error) {
	report := ContributionReportData{
		TargetDate:      targetDate,
		Data:            map[string]ContributionReportRow{},
		TotalAmount:     0,
		AveragePerMonth: 0,
		MembersCount:    0,
	}

	return &report, nil
}
