package tests

import (
	"testing"

	"github.com/kachaje/sacco-schema/database"
	"github.com/kachaje/sacco-schema/menus"
	menufuncs "github.com/kachaje/sacco-schema/menus/menuFuncs"
	"github.com/kachaje/utils/utils"
	"github.com/kachaje/workflow-parser/parser"
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

	session := parser.NewSession(nil, nil, nil, nil)

	result := m.LoadMenu("signUp", session, "", text, "")

	target := `
Welcome! Select Action

1. Sign In
2. Sign Up
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}

	text = "2"

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

Username: (Required Field)

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}

	text = username

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

What's your name? : 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}

	text = "Test User"

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

PIN Code: 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}

	text = password

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

Confirm PIN: 

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}

	text = password

	result = m.LoadMenu("signUp", session, "", text, "")

	target = `
Member SignUp

Welcome on board!

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}
