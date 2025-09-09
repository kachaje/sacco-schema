package menus_test

import (
	"fmt"
	"sacco/menus"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestRegistrationSubMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil)

	session.GlobalIds = map[string]any{
		"memberId": map[string]any{
			"key":   "member.id",
			"value": "1",
		},
		"memberLoanId": map[string]any{
			"key":   "member.memberLoan.0.id",
			"value": "1",
		},
	}

	text := "1"

	result := m.LoadMenu("main", session, "", text, "")

	target := `
CON Choose Activity
1. Member Details
2. Contact Details
3. Beneficiaries
4. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "3"

	result = m.LoadMenu(session.CurrentMenu, session, "", text, "")

	fmt.Println(result)
}
