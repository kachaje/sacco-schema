package menufuncs

import (
	"sacco/server/parser"
)

func MemberLoansSummary(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string = "Member Loans Summary\n\n00. Main Menu"

	_ = data

	return response
}
