package reports

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sacco/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
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

func (r *Reports) ContributionsReport2Table(reportData ContributionReportData) (*string, error) {
	targetDate, err := time.Parse("2006-01-02", reportData.TargetDate)
	if err != nil {
		return nil, err
	}

	pattern := func(left bool, size int) string {
		if left {
			return "%-" + fmt.Sprint(size) + "v"
		} else {
			return "%" + fmt.Sprint(size) + "v"
		}
	}
	line := func(size int) string {
		return strings.Repeat("-", size)
	}

	size := 16

	lines := []string{
		utils.CenterString("Contributions Report", (size * 9)),
		utils.CenterString(targetDate.Format("2 January, 2006"), (size * 9)),
	}

	var rows = [][]string{
		{fmt.Sprintf(pattern(true, 3), line(3))},
		{fmt.Sprintf(pattern(true, 3), "")},
		{fmt.Sprintf(pattern(true, 3), line(3))},
	}

	data := map[string]any{}

	jsonBytes, err := json.Marshal(reportData.Data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	fields := []string{
		"1-memberName",
		"2-contributionId",
		"3-monthlyContribution",
		"4-memberTotal",
		"5-updatedOn",
		"6-memberId",
		"7-percentOfTotal",
	}

	for key := range data {
		keys = append(keys, key)
	}

	utils.SortSlice(keys)
	sort.Strings(fields)

	formatNumber := func(row []string, v any, decimals bool, localSize int) []string {
		p := message.NewPrinter(language.English)

		var vn float64

		vr, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if err == nil {
			vn = vr
		}

		if decimals {
			row = append(row, fmt.Sprintf(pattern(false, localSize), p.Sprintf("%0.2f", number.Decimal(vn))))
		} else {
			row = append(row, fmt.Sprintf(pattern(false, localSize), p.Sprintf("%f", number.Decimal(vn))))
		}

		return row
	}

	for i, key := range keys {
		value := data[key]

		if val, ok := value.(map[string]any); ok {
			row := []string{
				fmt.Sprintf(pattern(false, 3), key),
			}

			for j, f := range fields {
				parts := strings.Split(f, "-")

				if len(parts) < 2 {
					continue
				}

				localSize := size

				switch j {
				case 0:
					localSize = size + 12
				case 2:
					localSize = size + 5
				}

				field := parts[1]

				if i == 0 {
					rows[0] = append(rows[0], line(localSize))
					rows[1] = append(rows[1], fmt.Sprintf(pattern(true, localSize), utils.IdentifierToLabel(field)))
					rows[2] = append(rows[2], line(localSize))
				}

				v := val[field]

				if v == nil {
					v = 0
				}

				if regexp.MustCompile(`^[0-9\.\+e]+$`).MatchString(fmt.Sprintf("%v", v)) &&
					!regexp.MustCompile(NUMBER_FORMAT_ESCAPE).MatchString(strings.ToLower(field)) {
					if false {
						p := message.NewPrinter(language.English)

						var vn float64

						vr, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
						if err == nil {
							vn = vr
						}

						row = append(row, fmt.Sprintf(pattern(false, localSize), p.Sprintf("%0.2f", number.Decimal(vn))))
					}

					row = formatNumber(row, v, true, localSize)
				} else {
					if regexp.MustCompile(`Date$`).MatchString(fmt.Sprintf("%v", field)) {
						v = fmt.Sprintf("%v", v)[0:10]
					}

					row = append(row, fmt.Sprintf(pattern(true, localSize), v))
				}
			}

			rows = append(rows, row)
		}
	}

	row := []string{
		line(3),
	}
	for i := range fields {
		localSize := size

		switch i {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		row = append(row, line(localSize))
	}
	rows = append(rows, row)

	row = []string{
		fmt.Sprintf(pattern(true, 3), ""),
	}
	for i := range fields {
		localSize := size

		switch i {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		switch i {
		case 0:
			row = append(row, fmt.Sprintf(pattern(true, localSize), "Total Contributions:"))
		case 1:
			row = append(row, fmt.Sprintf(pattern(true, localSize), ""))
		case 3:
			row = formatNumber(row, reportData.TotalAmount, true, localSize)
		case 7:
			row = formatNumber(row, reportData.AveragePerMonth, true, localSize)
		default:
			row = append(row, fmt.Sprintf(pattern(true, localSize), ""))
		}
	}
	rows = append(rows, row)

	row = []string{
		line(3),
	}
	for i := range fields {
		localSize := size

		switch i {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		row = append(row, line(localSize))
	}
	rows = append(rows, row)

	row = []string{
		fmt.Sprintf(pattern(true, 3), ""),
	}
	for i := range fields {
		localSize := size

		switch i {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		switch i {
		case 0:
			row = append(row, fmt.Sprintf(pattern(true, localSize), "Avg. Monthly Contribution:"))
		case 2:
			row = formatNumber(row, reportData.AveragePerMonth, false, localSize)
		default:
			row = append(row, fmt.Sprintf(pattern(true, localSize), ""))
		}
	}
	rows = append(rows, row)

	row = []string{
		line(3),
	}
	for i := range fields {
		localSize := size

		switch i {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		row = append(row, line(localSize))
	}
	rows = append(rows, row)

	row = []string{
		fmt.Sprintf(pattern(true, 3), ""),
	}
	for i := range fields {
		localSize := size

		switch i {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		switch i {
		case 0:
			row = append(row, fmt.Sprintf(pattern(true, localSize), "Count of Members:"))
		case 5:
			row = formatNumber(row, len(reportData.Data), false, localSize)
		default:
			row = append(row, fmt.Sprintf(pattern(true, localSize), ""))
		}
	}
	rows = append(rows, row)

	row = []string{
		line(3),
	}
	for j := range fields {
		localSize := size

		switch j {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		row = append(row, line(localSize))
	}
	rows = append(rows, row)

	for j, row := range rows {
		localSize := size

		switch j {
		case 0:
			localSize = size + 12
		case 2:
			localSize = size + 5
		}

		if regexp.MustCompile(line(localSize)).MatchString(strings.Join(row, "")) {
			lines = append(lines, strings.Join(row, "-+-"))
		} else {
			lines = append(lines, strings.Join(row, " | "))
		}
	}

	result := strings.Join(lines, "\n")

	return &result, nil
}
