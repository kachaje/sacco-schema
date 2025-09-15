package utils

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func Map2Table(data map[string]any, selectFields []string) string {
	pattern := func(left bool) string {
		if left {
			return "%-13v"
		} else {
			return "%13v"
		}
	}
	line := strings.Repeat("-", 13)

	var rows = [][]string{
		{fmt.Sprintf(pattern(true), "")},
		{fmt.Sprintf(pattern(true), line)},
	}
	var result string

	fields := []string{}
	keys := []string{}

	for key, value := range data {
		keys = append(keys, key)

		if val, ok := value.(map[string]any); ok {
			for field := range val {
				if !slices.Contains(fields, field) {
					fields = append(fields, field)
				}
			}
		}
	}

	sort.Strings(keys)
	sort.Strings(fields)

	for i, key := range keys {
		value := data[key]

		if val, ok := value.(map[string]any); ok {
			row := []string{
				fmt.Sprintf(pattern(true), key),
			}

			for _, field := range fields {
				if selectFields != nil && !slices.Contains(selectFields, field) {
					continue
				}

				if i == 0 {
					rows[0] = append(rows[0], fmt.Sprintf(pattern(true), IdentifierToLabel(field)))
					rows[1] = append(rows[1], line)
				}

				v := val[field]

				if v == nil {
					v = 0
				}

				if regexp.MustCompile(`^[0-9\.\+e]+$`).MatchString(fmt.Sprintf("%v", v)) &&
					!regexp.MustCompile(`phone|bill`).MatchString(strings.ToLower(field)) {
					p := message.NewPrinter(language.English)

					var vn float64

					vr, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
					if err == nil {
						vn = vr
					}

					row = append(row, fmt.Sprintf(pattern(false), p.Sprintf("%0.2f", number.Decimal(vn))))
				} else {
					row = append(row, fmt.Sprintf(pattern(false), v))
				}
			}

			rows = append(rows, row)
		}
	}

	for _, row := range rows {
		result = fmt.Sprintf("%s%s\n", result, strings.Join(row, " | "))
	}

	return result
}
