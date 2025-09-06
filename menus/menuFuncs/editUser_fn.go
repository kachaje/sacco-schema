package menufuncs

import (
	"sacco/parser"
)

func EditUser(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string = "Edit User\n\n00. Main Menu"

	_ = data

	return response
}
