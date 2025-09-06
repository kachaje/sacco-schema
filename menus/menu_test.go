package menus_test

import (
	"sacco/menus"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestLanding(t *testing.T) {
	m := menus.NewMenus(nil, nil)

	session := parser.NewSession(nil, nil, nil)

	result := m.LoadMenu("main", session, "", "", "")

	target := `
Welcome! Select Action

1. Sign In
2. Sign Up
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}

func TestMainMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil)

	role := "admin"

	session.SessionUserRole = &role

	result := m.LoadMenu("main", session, "", "", "")

	target := `
CON Welcome to Kaso SACCO
1. Membership Application
2. Loans
3. Check Balance
4. Banking Details
5. Preferred Language
6. Administration
7. Exit
9. Set PhoneNumber
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}

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

	result := m.LoadMenu("main", session, "", "1", "")

	target := `
CON Choose Activity
1. Member Details
2. Contact Details
3. Next of Kin Details
4. Beneficiaries
5. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}

func TestBusinessMenu(t *testing.T) {
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
		"memberBusinessId": map[string]any{
			"key":   "member.memberLoan.0.memberBusiness.id",
			"value": "1",
		},
	}

	result := m.LoadMenu("business", session, "", "", "")

	target := `
CON Business
1. Business Details
2. Previous Year History
3. Next Year Projection
4. Business Summary

99. Cancel
00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}

func TestEmployementMenu(t *testing.T) {
	t.Skip()

	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil)

	result := m.LoadMenu("employment", session, "", "", "")

	target := `CON Employement
3. Employement Summary

99. Cancel
00. Main Menu`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

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

	result = m.LoadMenu("employment", session, "", "", "")

	target = `CON Employement
1. Employement Details
3. Employement Summary

99. Cancel
00. Main Menu`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	session.GlobalIds = map[string]any{
		"memberId": map[string]any{
			"key":   "member.id",
			"value": "1",
		},
		"memberLoanId": map[string]any{
			"key":   "member.memberLoan.0.id",
			"value": "1",
		},
		"memberOccupationId": map[string]any{
			"key":   "member.memberLoan.0.memberOccupation.id",
			"value": "1",
		},
	}
	session.AddedModels["memberOccupation"] = true

	result = m.LoadMenu("employment", session, "", "", "")

	target = `CON Employement
1. Employement Details (*)
2. Employement Verification
3. Employement Summary

99. Cancel
00. Main Menu`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}

func TestLoanMenu(t *testing.T) {
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
		"memberBusinessId": map[string]any{
			"key":   "member.memberLoan.0.memberBusiness.id",
			"value": "1",
		},
	}

	result := m.LoadMenu("loan", session, "", "", "")

	target := `
CON Loans
1. Loan Application
6. Employment Details
7. Business Details
8. Member Loans Summary

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	role := "admin"

	session.SessionUserRole = &role

	result = m.LoadMenu("loan", session, "", "", "")

	target = `
CON Loans
1. Loan Application
2. Loan Liability
3. Loan Security
4. Loan Witness
5. Loan Approvals
6. Employment Details
7. Business Details
8. Member Loans Summary

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}
