package menufuncs

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kachaje/sacco-schema/parser"

	_ "embed"
)

//go:embed templates/employmentSummary.template.json
var employmentSummaryContent []byte

var employmentSummaryData map[string]any

func init() {
	err := json.Unmarshal(employmentSummaryContent, &employmentSummaryData)
	if err != nil {
		log.Fatalf("menus.init: %s", err.Error())
	}
}

func EmploymentSummary(
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
			session.CurrentMenu = "loan"
			text = ""
			return loadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)
		} else {
			data = LoadTemplateData(session.ActiveData, employmentSummaryData, &RefDate)

			table := TabulateData(data)

			tableString := strings.Join(table, "\n")

			response = "CON Employement Summary\n" +
				"\n" +
				fmt.Sprintf("%s\n", tableString) +
				"\n" +
				"99. Cancel\n" +
				"00. Main Menu"
		}
	}

	return response
}
