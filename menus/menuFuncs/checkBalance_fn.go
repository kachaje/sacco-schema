package menufuncs

import (
	"github.com/kachaje/sacco-schema/parser"
)

func CheckBalance(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var result string = "Check Balance\n\n" +
		"00. Main Menu\n"

	_ = data

	return result
}
