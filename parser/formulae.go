package parser

import (
	"fmt"
	"math"
	"regexp"
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
		for _, key := range terms {
			if data[key] == nil {
				continue
			}

			startDate, err := time.Parse("2006-01-02", fmt.Sprintf("%v", data[key]))
			if err != nil {
				return nil, err
			}

			duration := time.Since(startDate)

			result = math.Round(duration.Abs().Hours() / (365 * 24))

			break
		}
	}

	return &result, nil
}
