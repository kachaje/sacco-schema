package menufuncs_test

import (
	"sacco/database"
	"sacco/menus"
	menufuncs "sacco/menus/menuFuncs"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestSignInFn(t *testing.T) {
	var text string

	menufuncs.DB = database.NewDatabase(":memory:")
	defer func() {
		menufuncs.DB.DB.Close()
	}()

	username := "testuser"
	password := "password"

	_, err := menufuncs.DB.GenericModels["user"].AddRecord(map[string]any{
		"name":     "Test User",
		"username": username,
		"password": password,
		"userRole": "Admin",
	})
	if err != nil {
		t.Fatal(err)
	}

	m := menus.NewMenus(nil, nil)

	session := parser.NewSession(nil, nil, nil, nil)

	result := m.LoadMenu("signIn", session, "", text, "")

	target := `
Welcome! Select Action

1. Sign In
2. Sign Up
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu("signIn", session, "", text, "")

	target = `
Login

Username: 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = username

	result = m.LoadMenu("signIn", session, "", text, "")

	target = `
Login

PIN Code: 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = password

	result = m.LoadMenu("signIn", session, "", text, "")

	target = `
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
		t.Fatal("Test failed")
	}
}
