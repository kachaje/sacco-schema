package menufuncs

import (
	"fmt"
	"regexp"
	"sacco/parser"
	"slices"
)

func SignUp(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string
	var phoneNumber, text, preferencesFolder string
	var content string

	title := "Member SignUp\n\n"
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
	askNewPassword := func(msg string) string {
		return fmt.Sprintf("PIN Code: %s\n", msg)
	}
	askConfirmPassword := func(msg string) string {
		return fmt.Sprintf("Confirm PIN: %s\n", msg)
	}
	askName := func(msg string) string {
		return fmt.Sprintf("What's your name? : %s\n", msg)
	}

	if session.LastPrompt == "username" &&
		slices.Contains([]string{"", "000"}, text) &&
		regexp.MustCompile(`^\d+$`).MatchString(phoneNumber) {
		if !DB.UsernameFree(phoneNumber) {
			session.LastPrompt = "username"

			content = askUsername(fmt.Sprintf("(%s already taken)", phoneNumber))
		} else {
			session.Cache["username"] = phoneNumber

			session.LastPrompt = "fullname"

			content = askName("")
		}
	} else {
		switch session.LastPrompt {
		case "username":
			if text == "" {
				content = askUsername("(Required Field)")
			} else {
				if !DB.UsernameFree(text) {
					session.LastPrompt = "username"

					content = askUsername(fmt.Sprintf("(%s already taken)", text))
				} else {
					session.Cache["username"] = text

					text = ""

					session.LastPrompt = "fullname"

					content = askName("")
				}
			}
		case "fullname":
			if text == "" {
				content = askName("(Required Field)")
			} else {
				session.Cache["name"] = text

				text = ""

				session.LastPrompt = "newPassword"

				content = askNewPassword("")
			}
		case "newPassword":
			if text == "" {
				content = askNewPassword("(Required Field)")
			} else {
				session.Cache["password"] = text

				text = ""

				session.LastPrompt = "confirmPassword"

				content = askConfirmPassword("")
			}
		case "confirmPassword":
			if text == "" {
				content = askConfirmPassword("(Required Field)")
			} else {
				if text != session.Cache["password"] {
					session.LastPrompt = "newPassword"

					content = askNewPassword("(password mismatch)")
				} else {
					_, err := DB.GenericModels["user"].AddRecord(map[string]any{
						"name":     session.Cache["name"],
						"username": session.Cache["username"],
						"password": session.Cache["password"],
						"userRole":     "Member",
					})
					if err != nil {
						content = err.Error()
					} else {
						session.Cache = map[string]string{}
						session.LastPrompt = ""

						content = "Welcome on board!\n"
					}
				}
			}
		default:
			session.LastPrompt = "username"

			content = askUsername("")
		}
	}

	response = fmt.Sprintf("%s%s%s", title, content, footer)

	return response
}
