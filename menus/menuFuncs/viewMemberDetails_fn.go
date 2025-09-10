package menufuncs

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sacco/parser"
	"strings"

	_ "embed"
)

//go:embed templates/member.template.json
var memberTemplateContent []byte

var memberTemplateData map[string]any

func init() {
	err := json.Unmarshal(memberTemplateContent, &memberTemplateData)
	if err != nil {
		log.Fatalf("menus.init: %s", err.Error())
	}
}

func ViewMemberDetails(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var preferredLanguage *string
	var response string
	var phoneNumber, text, preferencesFolder string

	if data["preferredLanguage"] != nil {
		if val, ok := data["preferredLanguage"].(*string); ok {
			preferredLanguage = val
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
	if data["memberTemplateData"] != nil {
		if val, ok := data["memberTemplateData"].(map[string]any); ok {
			memberTemplateData = val
		}
	}

	if session != nil {
		if strings.TrimSpace(text) == "99" {
			parentMenu := "main"

			if regexp.MustCompile(`\.\d+$`).MatchString(session.CurrentMenu) {
				parentMenu = regexp.MustCompile(`\.\d+$`).ReplaceAllLiteralString(session.CurrentMenu, "")
			}

			session.CurrentMenu = parentMenu
			text = ""
			return loadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)
		} else {
			data = LoadTemplateData(session.ActiveData, memberTemplateData)

			table := TabulateData(data)

			tableString := strings.Join(table, "\n")

			if preferredLanguage != nil && *preferredLanguage == "ny" {
				response = "CON Zambiri za Membala\n" +
					"\n" +
					fmt.Sprintf("%s\n", tableString) +
					"\n" +
					"99. Basi\n" +
					"00. Tiyambirenso"
			} else {
				response = "CON Member Details\n" +
					"\n" +
					fmt.Sprintf("%s\n", tableString) +
					"\n" +
					"99. Cancel\n" +
					"00. Main Menu"
			}
		}
	} else {
		response = "Member Details\n\n" +
			"00. Main Menu\n"
	}

	return response
}
