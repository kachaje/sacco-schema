package parser

import (
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func GetTokens(query string) map[string]any {
	result := map[string]any{}

	re := regexp.MustCompile(`^([A-z_]+)`)

	op := re.FindAllString(query, -1)[0]

	result["op"] = op

	re = regexp.MustCompile(`([A-Za-z_]+)`)

	result["terms"] = []string{}

	for _, term := range re.FindAllString(query, -1) {
		if term != op {
			result["terms"] = append(result["terms"].([]string), term)
		}
	}

	return result
}

func ResultFromFormulae(tokens, data map[string]any) (*float64, error) {
	var result float64
	var op string
	var terms []string

	if tokens["op"] == nil {
		return nil, fmt.Errorf("missing required op token")
	}
	if tokens["terms"] == nil {
		return nil, fmt.Errorf("missing required terms token")
	}

	if val, ok := tokens["op"].(string); ok {
		op = val
	}
	if val, ok := tokens["terms"].([]any); ok {
		for _, term := range val {
			if val, ok := term.(string); ok {
				terms = append(terms, val)
			}
		}
	} else if val, ok := tokens["terms"].([]string); ok {
		terms = append(terms, val...)
	}

	switch strings.ToUpper(op) {
	case "SUM":
		for _, term := range terms {
			if data[term] != nil {
				val, err := strconv.ParseFloat(fmt.Sprintf("%v", data[term]), 64)
				if err == nil {
					result += val
				}
			}
		}
	case "DIFF":
		for i, term := range terms {
			if data[term] != nil {
				val, err := strconv.ParseFloat(fmt.Sprintf("%v", data[term]), 64)
				if err == nil {
					if i == 0 {
						result += val
					} else {
						result -= val
					}
				}
			}
		}
	case "DATE_DIFF_YEARS":
		var startDate time.Time
		var refDate time.Time
		var err error

		if data["startDate"] == nil {
			return nil, fmt.Errorf("missing required startDate input")
		}
		if data["refDate"] == nil {
			return nil, fmt.Errorf("missing required refDate input")
		}

		startDate, err = time.Parse("2006-01-02", fmt.Sprintf("%v", data["startDate"]))
		if err != nil {
			return nil, err
		}
		refDate, err = time.Parse("2006-01-02", fmt.Sprintf("%v", data["refDate"]))
		if err != nil {
			return nil, err
		}

		duration := refDate.Sub(startDate)

		result = math.Round(duration.Abs().Hours() / (365 * 24))

	case "DIV":
		var numerator float64
		var denominator float64

		if len(terms) > 1 {
			for key, value := range data {
				rawKey := regexp.MustCompile(`\d+$`).ReplaceAllLiteralString(key, "")

				switch rawKey {
				case terms[0]:
					val, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
					if err == nil {
						numerator = val
					} else {
						return nil, err
					}
				case terms[1]:
					val, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
					if err == nil {
						denominator = val
					} else {
						return nil, err
					}
				}
			}
		}

		if denominator > 0 {
			result = numerator / denominator
		}
	}

	return &result, nil
}

func GetScheduleParams(query string) map[string]any {
	var result = map[string]any{}

	query = regexp.MustCompile(`\{|\}`).ReplaceAllLiteralString(query, "")

	re := regexp.MustCompile(`^([A-Za-z_]+)\(([^,]+),([^,]+),\[([^\]]+)\],\[([^\]]+)\]\)$`)

	if re.MatchString(query) {
		var op, amount, duration string
		var oneTimeRates []string
		var recurringRates []string

		parts := re.FindAllStringSubmatch(query, -1)[0]

		if len(parts) == 6 {
			op = parts[1]
			amount = parts[2]
			duration = parts[3]

			for key := range strings.SplitSeq(parts[4], ",") {
				oneTimeRates = append(oneTimeRates, key)
			}

			for key := range strings.SplitSeq(parts[5], ",") {
				recurringRates = append(recurringRates, key)
			}

			result = map[string]any{
				"op":             op,
				"amount":         amount,
				"duration":       duration,
				"oneTimeRates":   oneTimeRates,
				"recurringRates": recurringRates,
			}
		}
	}

	return result
}

func GenerateSchedule(query string, data map[string]any) (map[string]any, error) {
	var schedule = map[string]any{}
	var op, amount, duration string
	var oneTimeRates []string
	var recurringRates []string

	tokens := GetScheduleParams(query)

	if tokens == nil {
		return nil, fmt.Errorf("invalid query")
	}

	for _, key := range []string{
		"amount", "duration", "op",
		"oneTimeRates", "recurringRates",
	} {
		if tokens[key] == nil {
			return nil, fmt.Errorf("missing required %s in query", key)
		}

		if slices.Contains([]string{"oneTimeRates", "recurringRates"}, key) {
			arrVal := []string{}

			if val, ok := tokens[key].([]string); ok {
				arrVal = val
			} else if val, ok := tokens[key].([]any); ok {
				for _, k := range val {
					arrVal = append(arrVal, fmt.Sprintf("%v", k))
				}
			} else {
				return nil, fmt.Errorf("failed to parse required %s from query", key)
			}

			switch key {
			case "oneTimeRates":
				oneTimeRates = arrVal
			case "recurringRates":
				recurringRates = arrVal
			}
		} else if val, ok := tokens[key].(string); ok {
			switch key {
			case "op":
				op = val
			case "amount":
				amount = val
			case "duration":
				duration = val
			}
		} else {
			return nil, fmt.Errorf("failed to parse required %s from query", key)
		}
	}

	switch strings.ToUpper(op) {
	case "REDUCING_SCHEDULE":
		keys := []string{amount, duration}
		keys = append(keys, oneTimeRates...)
		keys = append(keys, recurringRates...)

		values := map[string]float64{}

		for _, key := range keys {
			if data[key] == nil {
				return nil, fmt.Errorf("missing required %s value", key)
			} else {
				var val float64
				var err error

				if val, err = strconv.ParseFloat(fmt.Sprintf("%v", data[key]), 64); err != nil {
					return nil, err
				} else {
					values[key] = val
				}
			}
		}

		baseAmount := values[amount]
		period := values[duration]

		installmentPerMonth := baseAmount / period

		for i := range int(period) {
			principal := baseAmount - (float64(i) * installmentPerMonth)

			row := map[string]float64{
				"principal":   principal,
				"installment": installmentPerMonth,
				"totalDue":    installmentPerMonth,
			}

			if i == 0 {
				for _, key := range oneTimeRates {
					localKey := regexp.MustCompile(`Rate$`).ReplaceAllLiteralString(key, "")
					localKey = regexp.MustCompile(`^monthlyI`).ReplaceAllLiteralString(localKey, "i")

					row[localKey] = principal * values[key]

					row["totalDue"] += row[localKey]
				}
			}

			for _, key := range recurringRates {
				localKey := regexp.MustCompile(`Rate$`).ReplaceAllLiteralString(key, "")
				localKey = regexp.MustCompile(`^monthlyI`).ReplaceAllLiteralString(localKey, "i")

				row[localKey] = principal * values[key]

				row["totalDue"] += row[localKey]
			}

			schedule[fmt.Sprintf("Month %v", i+1)] = row
		}
	}

	return schedule, nil
}
