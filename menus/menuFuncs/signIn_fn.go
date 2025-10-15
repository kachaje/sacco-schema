package menufuncs

import (
	"fmt"

	"github.com/kachaje/workflow-parser/parser"

	"github.com/google/uuid"
)

func SignIn(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string
	var phoneNumber, text, preferencesFolder, content string

	title := "Login\n\n"
	footer := "\n00. Main Menu\n"

	if data["phoneNumber"] != nil {
		if val, ok := data["phoneNumber"].(string); ok {
			phoneNumber = val
		}
	}
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
		return loadMenu("main", session, phoneNumber, "", preferencesFolder)
	}

	askUsername := func(msg string) string {
		return fmt.Sprintf("Username: %s\n", msg)
	}
	askPassword := func(msg string) string {
		return fmt.Sprintf("PIN Code: %s\n", msg)
	}

	switch session.LastPrompt {
	case "username":
		if text == "" {
			content = askUsername("(Required Field)")
		} else {
			session.Cache["username"] = text

			text = ""

			session.LastPrompt = "password"

			content = askPassword("")
		}
	case "password":
		if text == "" {
			content = askPassword("(Required Field)")
		} else {
			session.Cache["password"] = text

			text = ""

			if id, role, ok := DB.ValidatePassword(fmt.Sprintf("%v", session.Cache["username"]), fmt.Sprintf("%v", session.Cache["password"])); ok {
				token := uuid.NewString()
				session.SessionToken = &token
				session.SessionUserId = id
				session.SessionUserRole = role

				session.CurrentMenu = "main"

				username := fmt.Sprintf("%v", session.Cache["username"])

				session.SessionUser = &username

				session.Cache = map[string]any{}
				session.LastPrompt = ""

				return loadMenu("main", session, phoneNumber, text, preferencesFolder)
			} else {
				session.Cache = map[string]any{}
				session.LastPrompt = "username"

				content = askUsername("(Invalid credentials)")
			}
		}
	default:
		session.LastPrompt = "username"

		content = askUsername("")
	}

	response = fmt.Sprintf("%s%s%s", title, content, footer)

	return response
}
