package reports

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type ContributionReportRow struct {
	MemberName          string  `json:"memberName"`
	ContributionId      string  `json:"contributionId"`
	MonthlyContribution float64 `json:"monthlyContribution"`
	MemberTotal         float64 `json:"memberTotal"`
	UpdatedOn           string  `json:"updatedOn"`
	MemberId            string  `json:"memberId"`
	PercentOfTotal      float64 `json:"percentOfTotal"`
	MinDate             string  `json:"minDate"`
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

	rawData, err := r.DB.SQLQuery(fmt.Sprintf(`
WITH schedule AS (
	SELECT 
		CONCAT(m.lastName, ', ', m.firstname) AS memberName, 
		(
			SELECT contributionNumber 
			FROM memberContribution 
			WHERE id = s.memberContributionId
		) AS contributionId,
		s.expectedAmount AS monthlyContribution,
		SUM(s.paidAmount) AS memberTotal,
		MAX(s.dueDate) AS updatedOn,
		m.memberIdNumber AS memberId,
		MIN(s.dueDate) AS minDate
	FROM memberContributionSchedule s 
	LEFT OUTER JOIN memberContribution c ON c.id = s.memberContributionId
	LEFT OUTER JOIN member m ON m.id = c.memberId
	WHERE DATE(s.dueDate) <= DATE('%s')
	GROUP BY m.id
)
	SELECT 
		memberName, 
		contributionId, 
		monthlyContribution, 
		memberTotal, 
		updatedOn, 
		memberId, 
		memberTotal / (SELECT SUM(memberTotal) FROM schedule) AS percentOfTotal,
		minDate
	FROM schedule
`, targetDate))
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", targetDate)
	if err != nil {
		return nil, err
	}

	startDate := endDate

	for i, data := range rawData {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		reportRow := ContributionReportRow{}

		err = json.Unmarshal(jsonBytes, &reportRow)
		if err != nil {
			return nil, err
		}

		refDate, err := time.Parse("2006-01-02", reportRow.MinDate)
		if err == nil {
			if refDate.Before(startDate) {
				startDate = refDate
			}
		}

		report.TotalAmount += reportRow.MemberTotal
		report.MembersCount++

		report.Data[fmt.Sprint(i+1)] = reportRow
	}

	duration := math.Floor(endDate.Sub(startDate).Hours() / (24 * 30))

	averagePerMonth := report.TotalAmount / (duration * float64(report.MembersCount))

	report.AveragePerMonth = averagePerMonth

	return &report, nil
}
