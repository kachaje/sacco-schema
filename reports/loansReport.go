package reports

import (
	"encoding/json"
	"fmt"
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

var (
	NUMBER_FORMAT_ESCAPE = `phone|bill|serial|year|number`
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

func (r *Reports) LoansReport2Table(reportData LoansReportData) (*string, error) {
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
		utils.CenterString("Loans Report", (size * 9)),
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
		"1-lastName",
		"2-firstName",
		"3-memberIdNumber",
		"4-loanNumber",
		"5-loanAmount",
		"6-loanStartDate",
		"7-loanDueDate",
		"8-balanceAmount",
	}

	for key := range data {
		keys = append(keys, key)
	}

	utils.SortSlice(keys)
	sort.Strings(fields)

	parseMoney := func(row []string, v any) []string {
		p := message.NewPrinter(language.English)

		var vn float64

		vr, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if err == nil {
			vn = vr
		}

		row = append(row, fmt.Sprintf(pattern(false, size), p.Sprintf("%0.2f", number.Decimal(vn))))

		return row
	}

	for i, key := range keys {
		value := data[key]

		if val, ok := value.(map[string]any); ok {
			row := []string{
				fmt.Sprintf(pattern(true, 3), key),
			}

			for _, f := range fields {
				parts := strings.Split(f, "-")

				if len(parts) < 2 {
					continue
				}

				field := parts[1]

				if i == 0 {
					rows[0] = append(rows[0], line(size))
					rows[1] = append(rows[1], fmt.Sprintf(pattern(true, size), utils.IdentifierToLabel(field)))
					rows[2] = append(rows[2], line(size))
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

						row = append(row, fmt.Sprintf(pattern(false, size), p.Sprintf("%0.2f", number.Decimal(vn))))
					}

					row = parseMoney(row, v)
				} else {
					if regexp.MustCompile(`Date$`).MatchString(fmt.Sprintf("%v", field)) {
						v = fmt.Sprintf("%v", v)[0:10]
					}

					row = append(row, fmt.Sprintf(pattern(true, size), v))
				}
			}

			rows = append(rows, row)
		}
	}

	row := []string{
		line(3),
	}
	for range fields {
		row = append(row, line(size))
	}
	rows = append(rows, row)

	row = []string{
		fmt.Sprintf(pattern(true, 3), ""),
	}
	for i := range fields {
		switch i {
		case 0:
			row = append(row, fmt.Sprintf(pattern(true, size), "Total:"))
		case 4:
			row = parseMoney(row, reportData.TotalLoanAmount)
		case 7:
			row = parseMoney(row, reportData.TotalBalanceAmount)
		default:
			row = append(row, fmt.Sprintf(pattern(true, size), ""))
		}
	}
	rows = append(rows, row)

	row = []string{
		line(3),
	}
	for range fields {
		row = append(row, line(size))
	}
	rows = append(rows, row)

	row = []string{
		fmt.Sprintf(pattern(true, 3), ""),
	}
	for i := range fields {
		switch i {
		case 0:
			row = append(row, fmt.Sprintf(pattern(true, size), "Count:"))
		case 1:
			row = parseMoney(row, len(reportData.Data))
		default:
			row = append(row, fmt.Sprintf(pattern(true, size), ""))
		}
	}
	rows = append(rows, row)

	row = []string{
		line(3),
	}
	for range fields {
		row = append(row, line(size))
	}
	rows = append(rows, row)

	for _, row := range rows {
		if regexp.MustCompile(line(size)).MatchString(strings.Join(row, "")) {
			lines = append(lines, strings.Join(row, "-+-"))
		} else {
			lines = append(lines, strings.Join(row, " | "))
		}
	}

	result := strings.Join(lines, "\n")

	return &result, nil
}
