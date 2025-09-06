package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sacco/database"
	"sacco/utils"
	"testing"

	_ "modernc.org/sqlite"
)

var (
	tableName = "person"
)

func setupDb(dbname string) (*sql.DB, *database.Model, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return nil, nil, err
	}

	sqlStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstName TEXT,
		lastName TEXT,
		gender TEXT,
		height REAL,
		weight REAL,
		active INTEGER DEFAULT 1,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP
	);`, tableName)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, nil, err
	}

	fields := []string{"firstName", "lastName", "gender", "height", "weight"}

	model, err := database.NewModel(db, tableName, fields)
	if err != nil {
		return nil, nil, err
	}

	return db, model, nil
}

func TestNewModel(t *testing.T) {
	dbname := "test.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	_, _ = db, model
}

func TestAddRecord(t *testing.T) {
	dbname := "testAdd.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	data := map[string]any{
		"firstName": "Mary",
		"lastName":  "Banda",
		"gender":    "Female",
		"height":    168.0,
		"weight":    62.0,
	}

	mid, err := model.AddRecord(data)
	if err != nil {
		t.Fatal(err)
	}

	if mid == nil {
		t.Fatal("Test failed. Got nil id")
	}

	row := db.QueryRow(fmt.Sprintf(`SELECT
		id,
		firstName,
		lastName,
		gender,
		height,
		weight
	FROM %s WHERE id=?`, tableName), *mid)

	var id int64
	var weight, height float64
	var firstName,
		lastName,
		gender string

	err = row.Scan(&id,
		&firstName,
		&lastName,
		&gender,
		&height,
		&weight,
	)
	if err != nil {
		t.Fatal(err)
	}

	if firstName != data["firstName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["firstName"], firstName)
	}
	if lastName != data["lastName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["lastName"], lastName)
	}
	if gender != data["gender"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["gender"], gender)
	}
	if height != data["height"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", data["height"], height)
	}
	if weight != data["weight"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", data["weight"], weight)
	}
}

func TestUpdateRecord(t *testing.T) {
	dbname := "testUpdate.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	data := map[string]any{
		"firstName": "John",
		"lastName":  "Phiri",
		"gender":    "Male",
		"height":    172.0,
		"weight":    95.0,
	}

	result, err := utils.QueryWithRetry(
		db,
		context.Background(), 0,
		fmt.Sprintf(`INSERT INTO %s (
			firstName,
			lastName,
			gender,
			height,
			weight
		) VALUES (
		 	?, ?, ?, ?, ?
		)`, tableName), "Mary", "Banda", "Female", 162.0, 72.0,
	)
	if err != nil {
		t.Fatal(err)
	}

	var id int64

	if id, err = result.LastInsertId(); err != nil {
		t.Fatal(err)
	}

	err = model.UpdateRecord(data, id)
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRow(fmt.Sprintf(`SELECT
		id,
		firstName,
		lastName,
		gender,
		height,
		weight
	FROM %s WHERE id=?`, tableName), id)

	var weight, height float64
	var firstName,
		lastName,
		gender string

	err = row.Scan(&id,
		&firstName,
		&lastName,
		&gender,
		&height,
		&weight,
	)
	if err != nil {
		t.Fatal(err)
	}

	if firstName != data["firstName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["firstName"], firstName)
	}
	if lastName != data["lastName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["lastName"], lastName)
	}
	if gender != data["gender"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["gender"], gender)
	}
	if height != data["height"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", data["height"], height)
	}
	if weight != data["weight"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", data["weight"], weight)
	}
}

func TestFetchById(t *testing.T) {
	dbname := "testFetchById.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	target := map[string]any{
		"firstName": "John",
		"lastName":  "Phiri",
		"gender":    "Male",
		"height":    172.0,
		"weight":    95.0,
	}

	result, err := utils.QueryWithRetry(
		db,
		context.Background(), 0,
		fmt.Sprintf(`INSERT INTO %s (
			firstName,
			lastName,
			gender,
			height,
			weight
		) VALUES (
		 	?, ?, ?, ?, ?
		)`, tableName),
		target["firstName"],
		target["lastName"],
		target["gender"],
		target["height"],
		target["weight"],
	)
	if err != nil {
		t.Fatal(err)
	}

	var id int64

	if id, err = result.LastInsertId(); err != nil {
		t.Fatal(err)
	}

	data, err := model.FetchById(id)
	if err != nil {
		t.Fatal(err)
	}

	if target["firstName"].(string) != data["firstName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target["firstName"], data["firstName"])
	}
	if target["lastName"].(string) != data["lastName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target["lastName"], data["lastName"])
	}
	if target["gender"].(string) != data["gender"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target["gender"], data["gender"])
	}
	if target["height"].(float64) != data["height"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target["height"], data["height"])
	}
	if target["weight"].(float64) != data["weight"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target["weight"], data["weight"])
	}
}

func TestFilterBy(t *testing.T) {
	dbname := "testFilterBy.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	target := map[string]any{
		"firstName": "John",
		"lastName":  "Phiri",
		"gender":    "Male",
		"height":    172.0,
		"weight":    95.0,
	}

	_, err = utils.QueryWithRetry(
		db,
		context.Background(), 0,
		fmt.Sprintf(`INSERT INTO %s (
			firstName,
			lastName,
			gender,
			height,
			weight
		) VALUES (
		 	?, ?, ?, ?, ?
		)`, tableName),
		target["firstName"],
		target["lastName"],
		target["gender"],
		target["height"],
		target["weight"],
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := model.FilterBy(`WHERE firstName = "Mary"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatal("Test failed. Expected no results")
	}

	result, err = model.FilterBy(`WHERE firstName = "John"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) <= 0 {
		t.Fatal("Test failed. No matches found")
	}

	data := result[0]

	if target["firstName"].(string) != data["firstName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target["firstName"], data["firstName"])
	}
	if target["lastName"].(string) != data["lastName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target["lastName"], data["lastName"])
	}
	if target["gender"].(string) != data["gender"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target["gender"], data["gender"])
	}
	if target["height"].(float64) != data["height"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target["height"], data["height"])
	}
	if target["weight"].(float64) != data["weight"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target["weight"], data["weight"])
	}
}
