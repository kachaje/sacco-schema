package tests

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	_ "embed"

	_ "modernc.org/sqlite"
)

func TestSchema(t *testing.T) {
	var schemaStmt string
	var seedStmt string
	var triggersStmt string

	content, err := os.ReadFile(filepath.Join("..", "database", "schema", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}

	schemaStmt = string(content)

	content, err = os.ReadFile(filepath.Join("..", "database", "schema", "seed.sql"))
	if err != nil {
		t.Fatal(err)
	}

	seedStmt = string(content)

	content, err = os.ReadFile(filepath.Join("..", "database", "schema", "triggers", "triggers.sql"))
	if err != nil {
		t.Fatal(err)
	}

	triggersStmt = string(content)

	dbname := ":memory:"
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for _, statement := range []string{"PRAGMA journal_mode=WAL", schemaStmt, seedStmt, triggersStmt} {
		_, err = db.Exec(statement)
		if err != nil {
			t.Fatal(err)
		}
	}

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

	result, err := db.Exec(statement)
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

	rows, err := db.Query(fmt.Sprintf(`SELECT memberIdNumber FROM member WHERE id = %v`, id))
	if err != nil {
		t.Fatal(err)
	}

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
