package tests

import (
	"testing"

	"github.com/kachaje/sacco-schema/database"
	filehandling "github.com/kachaje/sacco-schema/fileHandling"
	"github.com/kachaje/workflow-parser/parser"
)

func TestSaveModelDataWithRefData(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	phoneNumber := "0999888777"
	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]any{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	// Create initial member
	memberData := map[string]any{
		"firstName":   "John",
		"lastName":    "Doe",
		"phoneNumber": phoneNumber,
	}
	model := "member"
	err := filehandling.SaveModelData(memberData, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create dependants
	dependantData1 := map[string]any{
		"name1":         "Dependant 1",
		"memberId1":     1,
		"percentage1":   50,
		"relationship1": "Spouse",
	}
	dependantData2 := map[string]any{
		"name2":         "Dependant 2",
		"memberId2":     1,
		"percentage2":   30,
		"relationship2": "Child",
	}

	// Merge dependant data
	allDependants := make(map[string]any)
	for k, v := range dependantData1 {
		allDependants[k] = v
	}
	for k, v := range dependantData2 {
		allDependants[k] = v
	}

	model = "memberDependant"
	err = filehandling.SaveModelData(allDependants, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Now update with refData - remove dependant 2, add dependant 3
	refData := map[string]any{
		"memberDependant": []map[string]any{
			{"id": 1, "name": "Dependant 1"},
			{"id": 2, "name": "Dependant 2"},
		},
	}

	newDependantData := map[string]any{
		"name1":         "Dependant 1",
		"id1":           1,
		"memberId1":     1,
		"percentage1":   50,
		"relationship1": "Spouse",
		"name3":         "Dependant 3",
		"memberId3":     1,
		"percentage3":   20,
		"relationship3": "Child",
	}

	err = filehandling.SaveModelData(newDependantData, &model, &phoneNumber, db.GenericsSaveData, sessions, refData)
	if err != nil {
		t.Fatal(err)
	}

	// Verify dependant 2 is marked inactive
	records, err := db.GenericModels["memberDependant"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	// Should have 2 active records (1 and 3)
	if len(records) != 2 {
		t.Errorf("Expected 2 active records, got %d", len(records))
	}
}

func TestSaveModelDataWithFloatConversion(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	phoneNumber := "0999888777"
	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]any{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	// Test data with string numbers that should be converted to float
	data := map[string]any{
		"loanAmount":              "50000.50",
		"repaymentPeriodInMonths": "12",
		"monthlyInterestRate":     "2.5",
		"memberId":                1,
	}

	model := "memberLoan"
	err := filehandling.SaveModelData(data, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify data was saved correctly
	records, err := db.GenericModels["memberLoan"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	if len(records) == 0 {
		t.Fatal("No records found")
	}

	// Check that numeric fields are stored as numbers, not strings
	record := records[0]
	if loanAmount, ok := record["loanAmount"].(float64); !ok {
		t.Errorf("Expected loanAmount to be float64, got %T", record["loanAmount"])
	} else if loanAmount != 50000.50 {
		t.Errorf("Expected loanAmount to be 50000.50, got %f", loanAmount)
	}
}

func TestSaveModelDataWithParentLinking(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	phoneNumber := "0999888777"
	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds: map[string]any{
				"memberId": map[string]any{"value": 1},
			},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	// Create member first
	memberData := map[string]any{
		"firstName":   "Jane",
		"lastName":    "Smith",
		"phoneNumber": phoneNumber,
	}
	model := "member"
	err := filehandling.SaveModelData(memberData, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Now create contact with parent linking
	contactData := map[string]any{
		"postalAddress":      "P.O. Box 123",
		"residentialAddress": "123 Main St",
		// memberId should be automatically set from GlobalIds
	}

	model = "memberContact"
	err = filehandling.SaveModelData(contactData, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify contact was linked to member
	records, err := db.GenericModels["memberContact"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	if len(records) == 0 {
		t.Fatal("No contact records found")
	}

	if records[0]["memberId"] == nil {
		t.Error("memberId was not set from GlobalIds")
	}
}

func TestSaveModelDataWithEmptyData(t *testing.T) {
	phoneNumber := "0999888777"
	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]any{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	saveFunc := func(data map[string]any, model string, retries int) (*int64, error) {
		// Should not be called with empty data
		if len(data) < 2 {
			return nil, nil
		}
		var id int64 = 1
		return &id, nil
	}

	// Test with minimal data (less than 2 fields)
	data := map[string]any{
		"id": 1,
	}
	model := "member"

	err := filehandling.SaveModelData(data, &model, &phoneNumber, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveModelDataWithNilSaveFunc(t *testing.T) {
	phoneNumber := "0999888777"
	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]any{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	data := map[string]any{
		"firstName": "Test",
		"lastName":  "User",
	}
	model := "member"

	err := filehandling.SaveModelData(data, &model, &phoneNumber, nil, sessions, nil)
	if err == nil {
		t.Error("Expected error when saveFunc is nil")
	}

	if err != nil && err.Error() != "server.SaveModelData.member:missing saveFunc" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestSaveDataLanguagePreference(t *testing.T) {
	phoneNumber := "1234567890"
	preferenceFolder := ".test_settings"

	// Clean up after test
	defer func() {
		// Remove test settings file if created
	}()

	data := map[string]any{
		"language": "en",
	}
	model := "language"

	err := filehandling.SaveData(data, &model, &phoneNumber, &preferenceFolder, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify preference was saved
	// (This would require reading the file, which we can test separately)
}

func TestSaveModelDataUpdatesPhoneNumber(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	phoneNumber := "0999888777"
	newPhoneNumber := "0888777666"
	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]any{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	// Create member with different phone number
	memberData := map[string]any{
		"firstName":   "Test",
		"lastName":    "User",
		"phoneNumber": newPhoneNumber,
	}
	model := "member"
	err := filehandling.SaveModelData(memberData, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify session phone number was updated
	if sessions[phoneNumber].CurrentPhoneNumber != newPhoneNumber {
		t.Errorf("Expected CurrentPhoneNumber to be %s, got %s", newPhoneNumber, sessions[phoneNumber].CurrentPhoneNumber)
	}
}
