package menufuncs

import (
	"fmt"
	"log"
	"regexp"

	"github.com/kachaje/sacco-schema/parser"
)

func SetPhoneNumber(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string
	var content, text, preferencesFolder string

	title := "CON Set PhoneNumber\n\n"
	footer := "\n00. Main Menu\n"

	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}
	if data["preferencesFolder"] != nil {
		if val, ok := data["preferencesFolder"].(string); ok {
			preferencesFolder = val
		}
	}

	if text == "00" {
		session.CurrentMenu = "main"
		return loadMenu("main", session, session.CurrentPhoneNumber, "", preferencesFolder)
	}

	askPhoneNumber := func(msg string) string {
		return fmt.Sprintf("Enter phone number: %s\n", msg)
	}

	if text != "" && text != "000" {
		if !regexp.MustCompile(`^\d+$`).MatchString(text) {
			content = askPhoneNumber("(Invalid input)")
		} else {
			session.CurrentPhoneNumber = text

			session.ClearSession()

			_, err := session.RefreshSession()
			if err != nil {
				log.Println(err)
			}

			text = ""
			content = "Success. Phone Number set!\n"
		}
	} else {
		content = askPhoneNumber("")
	}

	response = fmt.Sprintf("%s%s%s", title, content, footer)

	return response
}
