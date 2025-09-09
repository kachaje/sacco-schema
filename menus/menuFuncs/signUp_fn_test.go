package menufuncs_test

import (
	"sacco/database"
	"sacco/menus"
	menufuncs "sacco/menus/menuFuncs"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

func TestSignUpFn(t *testing.T) {
	var text string

	menufuncs.DB = database.NewDatabase(":memory:")
	defer func() {
		menufuncs.DB.DB.Close()
	}()

	username := "testuser"
	password := "password"

	m := menus.NewMenus(nil, nil)

	session := parser.NewSession(nil, nil, nil)

	result := m.LoadMenu("signUp", session, "", text, "")

	target := `
Welcome! Select Action

1. Sign In
2. Sign Up
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "2"

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

Username: (Required Field)

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = username

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

What's your name? : 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Test User"

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

PIN Code: 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = password

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

Confirm PIN: 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = password

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

Welcome on board!

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}
}
