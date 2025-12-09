package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/kachaje/sacco-schema/database"
)

func TestSchema(t *testing.T) {
	t.Parallel()
	db := database.NewDatabase(":memory:")
	defer db.Close()

	// Add new member
	statement := `INSERT INTO member (
		firstName,
		lastName,
		gender,
		phoneNumber,
		title,
		maritalStatus,
		dateOfBirth,
		nationalIdentifier,
		utilityBillType,
		utilityBillNumber
	) 
	VALUES (
		"Mary",
		"Banda",
		"Female",
		"0999888777",
		"Miss",
		"Single",
		"1999-09-01",
		"KJFFJ58584",
		"ESCOM",
		"949488473"
	)`

	result, err := db.DB.Exec(statement)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	if id <= 0 {
		t.Fatalf("Test failed. Got %v", id)
	}

	rows, err := db.DB.Query(fmt.Sprintf(`SELECT memberIdNumber FROM member WHERE id = %v`, id))
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var memberIdNumber string

	rows.Next()

	err = rows.Scan(&memberIdNumber)
	if err != nil {
		t.Fatal(err)
	}

	if !regexp.MustCompile(`^KSM\d{6}$`).MatchString(memberIdNumber) {
		t.Fatalf("Test failed. Got %v", memberIdNumber)
	}
}
