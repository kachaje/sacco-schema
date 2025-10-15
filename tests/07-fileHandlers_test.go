package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kachaje/sacco-schema/database"
	filehandling "github.com/kachaje/sacco-schema/fileHandling"
	"github.com/kachaje/sacco-schema/parser"
	"github.com/kachaje/utils/utils"
)

func TestSaveDataOne(t *testing.T) {
	phoneNumber := "0999888777"

	session := &parser.Session{
		AddedModels: map[string]bool{},
	}

	sessions := make(map[string]*parser.Session)

	sessions[phoneNumber] = session

	saveFunc := func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error) {
		var id int64 = 13

		return &id, nil
	}

	model := "member"

	data := map[string]any{
		"dateOfBirth":        "1999-09-01",
		"phoneNumber":        "09999999999",
		"fileNumber":         "",
		"firstName":          "Mary",
		"gender":             "Female",
		"id":                 1,
		"lastName":           "Banda",
		"maritalStatus":      "Single",
		"nationalIdentifier": "DHFYR8475",
		"oldFileNumber":      "",
		"otherName":          "",
		"title":              "Miss",
		"utilityBillNumber":  "29383746",
		"utilityBillType":    "ESCOM",
	}

	err := filehandling.SaveData(data, &model, &phoneNumber, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleBeneficiaries(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()
	}()

	data := map[string]any{
		"address1":      "P.O. Box 1234",
		"id1":           1,
		"memberId1":     1,
		"name1":         "Benefator 1",
		"phoneNumber1":  "0999888777",
		"percentage1":   35,
		"relationship1": "Spouse",
		"address2":      "P.O. Box 5678",
		"id2":           2,
		"memberId2":     1,
		"name2":         "Benefator 2",
		"phoneNumber2":  "0999777888",
		"percentage2":   25,
		"relationship2": "Child",
	}

	phoneNumber := "0999888777"

	sessions := map[string]*parser.Session{
		phoneNumber: {
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{"memberDependant": true},
		},
	}

	model := "memberDependant"

	err := filehandling.SaveModelData(data, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !sessions[phoneNumber].AddedModels["memberDependant"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v",
			sessions[phoneNumber].AddedModels["memberDependant"])
	}

	result, err := db.GenericModels["memberDependant"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}
}

func TestHandleMemberDetails(t *testing.T) {
	phoneNumber := "0999888777"

	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()
	}()

	data := map[string]any{
		"dateOfBirth":        "1999-09-01",
		"phoneNumber":        phoneNumber,
		"fileNumber":         "",
		"firstName":          "Mary",
		"gender":             "Female",
		"id":                 1,
		"lastName":           "Banda",
		"maritalStatus":      "Single",
		"nationalIdentifier": "DHFYR8475",
		"oldFileNumber":      "",
		"otherName":          "",
		"title":              "Miss",
		"utilityBillNumber":  "29383746",
		"utilityBillType":    "ESCOM",
	}

	sessions := map[string]*parser.Session{
		phoneNumber: {
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

	model := "member"

	err := filehandling.SaveModelData(data, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"memberContact.json",
		"memberDependant.json",
	} {
		content, err := os.ReadFile(filepath.Join(".", "fixtures", file))
		if err != nil {
			t.Fatal(err)
			continue
		}

		model := strings.Split(filepath.Base(file), ".")[0]

		if model == "memberDependant" {
			data := []map[string]any{}

			err = json.Unmarshal(content, &data)
			if err != nil {
				t.Fatal(err)
			}

			for _, row := range data {
				row["memberId"] = 1

				err = filehandling.SaveModelData(row, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
				if err != nil {
					t.Fatal(err)
				}
			}
		} else {
			data := map[string]any{}

			err = json.Unmarshal(content, &data)
			if err != nil {
				t.Fatal(err)
			}

			data["memberId"] = 1

			err = filehandling.SaveModelData(data, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	result, err := db.MemberByPhoneNumber(phoneNumber, nil)
	if err != nil {
		t.Fatal(err)
	}

	delete(result["member"].(map[string]any), "memberIdNumber")

	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.MapsEqual(target, result) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}

func TestChildNestedModel(t *testing.T) {
	phoneNumber := "0999888777"

	session := &parser.Session{
		AddedModels: map[string]bool{},
		GlobalIds: map[string]any{
			"memberId":           16,
			"memberLoanId":       13,
			"memberOccupationId": 1,
		},
	}

	sessions := make(map[string]*parser.Session)

	sessions[phoneNumber] = session

	count := 0

	saveFunc := func(
		data map[string]any,
		model string,
		retries int,
	) (*int64, error) {
		count++

		var id int64 = int64(count)

		data["id"] = id

		return &id, nil
	}

	model := "memberOccupation"

	sessions[phoneNumber].AddedModels["member"] = true

	data := map[string]any{
		"employerAddress":        "Kanengo",
		"employerName":           "SOBO",
		"employerPhone":          "01282373737",
		"grossPay":               100000,
		"highestQualification":   "Secondary",
		"jobTitle":               "Driver",
		"netPay":                 90000,
		"periodEmployedInMonths": "36",
	}

	err := filehandling.SaveModelData(data, &model, &phoneNumber, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("Test failed. Expected: 1; Actual: %v", count)
	}

	target := map[string]any{
		"memberId":           16,
		"memberLoanId":       13,
		"memberOccupationId": 1,
	}

	if !reflect.DeepEqual(target, session.GlobalIds) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, session.GlobalIds)
	}
}

func TestArrayChildData(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()
	}()

	phoneNumber := "0999888777"

	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]any{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{"memberDependant": true},
		},
	}

	data := map[string]any{
		"dateOfBirth":        "1999-09-01",
		"phoneNumber":        "09999999999",
		"fileNumber":         "",
		"firstName":          "Mary",
		"gender":             "Female",
		"id":                 1,
		"lastName":           "Banda",
		"maritalStatus":      "Single",
		"nationalIdentifier": "DHFYR8475",
		"oldFileNumber":      "",
		"otherName":          "",
		"title":              "Miss",
		"utilityBillNumber":  "29383746",
		"utilityBillType":    "ESCOM",
	}

	model := "member"

	err := filehandling.SaveModelData(data, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	data = map[string]any{
		"address1":      "P.O. Box 1234",
		"id1":           1,
		"memberId1":     1,
		"name1":         "Benefator 1",
		"percentage1":   35,
		"phoneNumber1":  "0999888777",
		"relationship1": "Spouse",
		"address2":      "P.O. Box 5678",
		"id2":           2,
		"memberId2":     1,
		"name2":         "Benefator 2",
		"percentage2":   25,
		"phoneNumber2":  "0888999777",
		"relationship2": "Child",
	}

	model = "memberDependant"

	err = filehandling.SaveModelData(data, &model, &phoneNumber, db.GenericsSaveData, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !sessions[phoneNumber].AddedModels["memberDependant"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v",
			sessions[phoneNumber].AddedModels["memberDependant"])
	}

	result, err := db.GenericModels["memberDependant"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}
}
