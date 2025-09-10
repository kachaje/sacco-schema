package menufuncs

import (
	"sacco/parser"
)

func Landing(loadMenu func(
	menuName string, session *parser.Session,
	phoneNumber, text, preferencesFolder string,
) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response string
	var phoneNumber, text, preferencesFolder string

	if data["session"] != nil {
		if val, ok := data["session"].(*parser.Session); ok {
			session = val
		}
	}
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

	data = map[string]any{
		"phoneNumber":       phoneNumber,
		"session":           session,
		"preferencesFolder": preferencesFolder,
		"text":              text,
	}

	session.LastPrompt = ""
	session.Cache = map[string]any{}

	switch text {
	case "1":
		session.CurrentMenu = "signIn"
		data["text"] = ""
		return SignIn(loadMenu, data, session)
	case "2":
		session.CurrentMenu = "signUp"
		session.LastPrompt = "username"
		data["text"] = ""
		return SignUp(loadMenu, data, session)
	default:
		response = "Welcome! Select Action\n\n" +
			"1. Sign In\n" +
			"2. Sign Up\n"
	}

	return response
}
