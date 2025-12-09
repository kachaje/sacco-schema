package menufuncs

import (
	"fmt"

	"github.com/kachaje/workflow-parser/parser"
)

func ChangePassword(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response, text, content string

	title := "Change Password\n\n"
	footer := "\n00. Main Menu\n"

	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}

	currentPassword := func(msg string) string {
		return fmt.Sprintf("Current Password: %s\n", msg)
	}
	newPassword := func(msg string) string {
		return fmt.Sprintf("New Password: %s\n", msg)
	}
	confirmPassword := func(msg string) string {
		return fmt.Sprintf("Confirm Password: %s\n", msg)
	}

	switch session.LastPrompt {
	case "currentPassword":
		if text == "" {
			content = currentPassword("(Required Field)")
		} else {
			session.Cache["currentPassword"] = text

			session.LastPrompt = "newPassword"

			content = newPassword("")
		}
	case "newPassword":
		if text == "" {
			content = newPassword("(Required Field)")
		} else {
			session.Cache["newPassword"] = text

			session.LastPrompt = "confirmPassword"

			content = confirmPassword("")
		}
	case "confirmPassword":
		if text == "" {
			content = confirmPassword("(Required Field)")
		} else {
			session.Cache["confirmPassword"] = text

			if session.Cache["newPassword"] != session.Cache["confirmPassword"] {
				session.LastPrompt = "newPassword"

				content = newPassword("(Password Mismatch!)")
			} else {
				if id, _, ok := DB.ValidatePassword(*session.SessionUser, fmt.Sprintf("%v", session.Cache["currentPassword"])); ok {
					err := DB.GenericModels["user"].UpdateRecord(map[string]any{
						"password": session.Cache["newPassword"],
					}, *id)
					if err != nil {
						content = fmt.Sprintf("ERROR: %s\n", err.Error())
					} else {
						session.Cache = map[string]any{}
						session.LastPrompt = ""

						content = "Password Changed!\n"
					}
				} else {
					session.LastPrompt = "currentPassword"

					content = currentPassword("(Invalid credentials)")
				}
			}
		}
	default:
		session.LastPrompt = "currentPassword"

		content = currentPassword("")
	}

	response = fmt.Sprintf("%s%s%s", title, content, footer)

	return response
}
