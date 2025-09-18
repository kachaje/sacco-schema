package tests

import (
	"sacco/menus"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestRegistrationSubMenu(t *testing.T) {
	demo := true

	m := menus.NewMenus(nil, &demo)

	session := parser.NewSession(nil, nil, nil, nil)

	session.GlobalIds = map[string]any{
		"memberId": map[string]any{
			"key":   "member.id",
			"value": "1",
		},
		"memberContactId": map[string]any{
			"key":   "member.memberContact.id",
			"value": "1",
		},
	}

	session.AddedModels = map[string]bool{
		"member":        true,
		"memberContact": true,
	}

	session.ActiveData = map[string]any{
		"member.createdAt":                              "2025-09-09 12:23:08",
		"member.dateJoined":                             "2025-09-09",
		"member.dateOfBirth":                            "1999-09-01",
		"member.firstName":                              "Mary",
		"member.gender":                                 "Female",
		"member.id":                                     1,
		"member.lastName":                               "Banda",
		"member.maritalStatus":                          "Single",
		"member.memberContact.createdAt":                "2025-09-09 12:23:33",
		"member.memberContact.homeDistrict":             "Karonga",
		"member.memberContact.homeTraditionalAuthority": "Kyungu",
		"member.memberContact.homeVillage":              "Songwe",
		"member.memberContact.id":                       1,
		"member.memberContact.memberId":                 1,
		"member.memberContact.postalAddress":            "P.O. Box 1",
		"member.memberContact.residentialAddress":       "Area 49",
		"member.memberContact.updatedAt":                "2025-09-09 12:23:33",
		"member.memberIdNumber":                         "KSM046018",
		"member.nationalIdentifier":                     "JDKD47483",
		"member.phoneNumber":                            "1234567890",
		"member.title":                                  "Miss",
		"member.updatedAt":                              "2025-09-09 12:23:08",
		"member.utilityBillNumber":                      "948476363",
		"member.utilityBillType":                        "ESCOM",
	}

	text := "1"

	result := m.LoadMenu(session.CurrentMenu, session, "", text, "")

	target := `
CON Choose Activity
1. Member Details (*)
2. Contact Details (*)
3. Beneficiaries
4. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}
