package tests

import (
	"sacco/menus"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestLanding(t *testing.T) {
	m := menus.NewMenus(nil, nil)

	session := parser.NewSession(nil, nil, nil, nil)

	result := m.LoadMenu("main", session, "", "", "")

	target := `
Welcome! Select Action

1. Sign In
2. Sign Up
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}

func TestMainMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil, nil)

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
7. Set PhoneNumber
8. Exit
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}

func TestBusinessMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil, nil)

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
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}

func TestEmployementMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil, nil)

	result := m.LoadMenu("employment", session, "", "", "")

	target := `CON Employement
3. Employement Summary

99. Cancel
00. Main Menu`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
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
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
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
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}

func TestLoanMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil, nil)

	session.GlobalIds = map[string]any{
		"memberId": map[string]any{
			"key":   "member.id",
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
7. Loan Repayment
8. Employment Details
9. Business Details
10. Member Loans Summary

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}

	role := "admin"

	session.SessionUserRole = &role

	session.GlobalIds["memberLoanId"] = map[string]any{
		"key":   "memberLoan.id",
		"value": "1",
	}

	session.GlobalIds["memberLoanApprovalId"] = map[string]any{
		"key":   "memberLoanApproval.id",
		"value": "1",
	}

	session.GlobalIds["memberLoanPaymentScheduleId"] = map[string]any{
		"key":   "memberLoanPaymentSchedule.id",
		"value": "1",
	}

	result = m.LoadMenu("loan", session, "", "", "")

	target = `
CON Loans
1. Loan Application
2. Loan Liability
3. Loan Security
4. Loan Witness
5. Loan Approvals
6. Loan Verification
7. Loan Repayment
8. Employment Details
9. Business Details
10. Member Loans Summary

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf(`Test failed. 
Expected: 
%s 
Actual: 
%s`, target, result)
	}
}
