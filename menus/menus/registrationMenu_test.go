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
	}

	session.AddedModels = map[string]bool{
		"member":          true,
		"memberContact":   true,
		"memberDependant": true,
		"memberNominee":   true,
	}

	session.ActiveData = map[string]any{
		"member.createdAt":                              "2025-09-09 11:49:53",
		"member.dateJoined":                             "2025-09-09",
		"member.dateOfBirth":                            "1999-09-01",
		"member.firstName":                              "Mary",
		"member.gender":                                 "Female",
		"member.id":                                     1,
		"member.lastName":                               "Banda",
		"member.maritalStatus":                          "Single",
		"member.memberContact.createdAt":                "2025-09-09 11:50:36",
		"member.memberContact.homeDistrict":             "Karonga",
		"member.memberContact.homeTraditionalAuthority": "Kyungu",
		"member.memberContact.homeVillage":              "Songwe",
		"member.memberContact.id":                       1,
		"member.memberContact.memberId":                 1,
		"member.memberContact.postalAddress":            "P.O. Box 1",
		"member.memberContact.residentialAddress":       "Area 49",
		"member.memberContact.updatedAt":                "2025-09-09 11:50:36",
		"member.memberDependant.1.createdAt":            "2025-09-09 11:52:48",
		"member.memberDependant.1.id":                   1,
		"member.memberDependant.1.isNominee":            "Yes",
		"member.memberDependant.1.memberId":             1,
		"member.memberDependant.1.name":                 "John Phiri",
		"member.memberDependant.1.phoneNumber":          "0999888777",
		"member.memberDependant.1.relationship":         "Spouse",
		"member.memberDependant.1.updatedAt":            "2025-09-09 11:52:48",
		"member.memberIdNumber":                         "KSM143579",
		"member.memberNominee.createdAt":                "2025-09-09 11:52:48",
		"member.memberNominee.id":                       1,
		"member.memberNominee.isNominee":                "Yes",
		"member.memberNominee.memberId":                 1,
		"member.memberNominee.name":                     "John Phiri",
		"member.memberNominee.phoneNumber":              "0999888777",
		"member.memberNominee.relationship":             "Spouse",
		"member.memberNominee.updatedAt":                "2025-09-09 11:52:48",
		"member.nationalIdentifier":                     "KDJD47483",
		"member.phoneNumber":                            "1234567890",
		"member.title":                                  "Miss",
		"member.updatedAt":                              "2025-09-09 11:49:53",
		"member.utilityBillNumber":                      "93844763",
		"member.utilityBillType":                        "ESCOM",
	}

	text := "1"

	result := m.LoadMenu(session.CurrentMenu, session, "", text, "")

	fmt.Println(result)
	
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
