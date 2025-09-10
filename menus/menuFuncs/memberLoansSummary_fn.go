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

//go:embed templates/loanApplication.template.json
var loanTemplateContent []byte

var loanTemplateData map[string]any

func init() {
	err := json.Unmarshal(loanTemplateContent, &loanTemplateData)
	if err != nil {
		log.Fatalf("menus.init: %s", err.Error())
	}
}

func MemberLoansSummary(
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
			data = LoadTemplateData(session.ActiveData, loanTemplateData)

			table := TabulateData(data)

			tableString := strings.Join(table, "\n")

			response = "CON Loan Application Summary\n" +
				"\n" +
				fmt.Sprintf("%s\n", tableString) +
				"\n" +
				"99. Cancel\n" +
				"00. Main Menu"
		}
	}

	return response
}
