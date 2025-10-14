package menufuncs

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/kachaje/sacco-schema/parser"

	_ "embed"
)

//go:embed templates/businessSummary.template.json
var businessSummaryTemplateContent []byte

var businessSummaryTemplateData map[string]any

func init() {
	err := json.Unmarshal(businessSummaryTemplateContent, &businessSummaryTemplateData)
	if err != nil {
		log.Fatalf("menus.init: %s", err.Error())
	}
}

func BusinessSummary(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var phoneNumber, text, preferencesFolder string
	var response string

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
			data = LoadTemplateData(session.ActiveData, businessSummaryTemplateData, &RefDate)

			table := TabulateData(data)

			tableString := strings.Join(table, "\n")

			response = "CON Business Summary\n" +
				"\n" +
				fmt.Sprintf("%s\n", tableString) +
				"\n" +
				"99. Cancel\n" +
				"00. Main Menu"
		}
	}

	return response
}
