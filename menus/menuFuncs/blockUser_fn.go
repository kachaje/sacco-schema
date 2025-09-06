package menufuncs

import (
	"sacco/server/parser"
)

func BlockUser(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string = "Block User\n\n00. Main Menu"

	_ = data

	return response
}
