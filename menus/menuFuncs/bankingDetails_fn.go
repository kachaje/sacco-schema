package menufuncs

import (
	"fmt"
	"sacco/parser"
)

func BankingDetails(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var preferredLanguage *string
	var response, text string

	if data["preferredLanguage"] != nil {
		if val, ok := data["preferredLanguage"].(*string); ok {
			preferredLanguage = val
		}
	}
	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}

	firstLine := "CON Banking Details\n"
	lastLine := "00. Main Menu\n"
	name := "Name"
	number := "Number"
	branch := "Branch"

	if preferredLanguage != nil && *preferredLanguage == "ny" {
		firstLine = "CON Matumizidwe\n"
		lastLine = "0. Bwererani Pofikira"
		name = "Dzina"
		number = "Nambala"
		branch = "Buranchi"
	}

	switch text {
	case "1":
		response = "CON National Bank of Malawi\n" +
			fmt.Sprintf("%-8s: Kaso SACCO\n", name) +
			fmt.Sprintf("%-8s: 1006857589\n", number) +
			fmt.Sprintf("%-8s: Lilongwe\n", branch) +
			"\n99. Cancel\n" +
			lastLine
	case "2":
		response = "CON Airtel Money\n" +
			fmt.Sprintf("%-8s: Kaso SACCO\n", name) +
			fmt.Sprintf("%-8s: 0985 242 629\n", number) +
			"\n99. Cancel\n" +
			lastLine
	default:
		response = firstLine +
			"1. National Bank\n" +
			"2. Airtel Money\n" +
			"\n" +
			lastLine
	}

	return response
}
