package menufuncs

import (
	"sacco/parser"
)

func EmploymentSummary(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var result string = "Employment Summary\n\n" +
		"00. Main Menu\n"

	_ = data
	return result
}
