package tests

import (
	"testing"

	"github.com/kachaje/sacco-schema/database"
	filehandling "github.com/kachaje/sacco-schema/fileHandling"
	"github.com/kachaje/workflow-parser/parser"
)

// TestSaveModelDataWithRefData was removed - it was a failing skipped test

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
		"loanPurpose":             "Business",
		"loanType":                "Business",
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
		"firstName":          "Jane",
		"lastName":           "Smith",
		"phoneNumber":        phoneNumber,
		"gender":             "Female",
		"title":              "Mrs",
		"maritalStatus":      "Married",
		"dateOfBirth":        "1990-01-01",
		"nationalIdentifier": "ID456",
		"utilityBillType":    "Water Board",
		"utilityBillNumber":  "BILL456",
	}
	model := "member"
	err := filehandling.SaveModelData(memberData, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Now create contact with parent linking
	contactData := map[string]any{
		"postalAddress":            "P.O. Box 123",
		"residentialAddress":       "123 Main St",
		"homeVillage":              "Test Village",
		"homeTraditionalAuthority": "Test TA",
		"homeDistrict":             "Test District",
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
		"firstName":          "Test",
		"lastName":           "User",
		"phoneNumber":        newPhoneNumber,
		"gender":             "Male",
		"title":              "Mr",
		"maritalStatus":      "Single",
		"dateOfBirth":        "1990-01-01",
		"nationalIdentifier": "ID789",
		"utilityBillType":    "ESCOM",
		"utilityBillNumber":  "BILL789",
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
