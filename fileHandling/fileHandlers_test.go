package filehandling_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sacco/database"
	filehandling "sacco/fileHandling"
	"sacco/parser"
	"sacco/utils"
	"strings"
	"testing"
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
	t.Skip()

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
		"memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		"memberDependant.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		"memberDependant.1efda9a6-84f4-11f0-8797-1e4d4999250c.json",
	} {
		content, err := os.ReadFile(filepath.Join(".", "fixtures", "cache", phoneNumber, file))
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

	target := map[string]any{
		"member": map[string]any{
			"dateOfBirth":   "1999-09-01",
			"fileNumber":    "",
			"firstName":     "Mary",
			"gender":        "Female",
			"id":            1,
			"lastName":      "Banda",
			"maritalStatus": "Single",
			"memberDependant": map[string]any{
				"1": map[string]any{
					"contact":    "0888777444",
					"id":         1,
					"memberId":   1,
					"name":       "John Phiri",
					"percentage": 10,
				},
				"2": map[string]any{
					"contact":    "07746635653",
					"id":         2,
					"memberId":   1,
					"name":       "Jean Banda",
					"percentage": 5,
				},
			},
			"memberContact": map[string]any{
				"email":                    any(nil),
				"homeDistrict":             "Lilongwe",
				"homeTraditionalAuthority": "Kalolo",
				"homeVillage":              "Kalulu",
				"id":                       1,
				"memberId":                 1,
				"postalAddress":            "P.O. Box 1",
				"residentialAddress":       "Area 49",
			},
			"memberNominee": map[string]any{
				"address":     "P.O. Box 1",
				"id":          1,
				"memberId":    1,
				"name":        "John Phiri",
				"phoneNumber": "0888444666",
			},
			"nationalIdentifier": "DHFYR8475",
			"oldFileNumber":      "",
			"otherName":          "",
			"phoneNumber":        "0999888777",
			"title":              "Miss",
			"utilityBillNumber":  "29383746",
			"utilityBillType":    "ESCOM",
		},
	}

	payloadResult, _ := json.MarshalIndent(result, "", "  ")
	payloadTarget, _ := json.MarshalIndent(target, "", "  ")

	if utils.CleanScript(payloadResult) != utils.CleanScript(payloadTarget) {
		t.Fatal("Test failed")
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
		t.Fatalf("Test failed")
	}
}

func TestArrayChildData(t *testing.T) {
	t.Skip()

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
		"contact1":    "P.O. Box 1234",
		"id1":         1,
		"memberId1":   1,
		"name1":       "Benefator 1",
		"percentage1": 35,
		"contact2":    "P.O. Box 5678",
		"id2":         2,
		"memberId2":   1,
		"name2":       "Benefator 2",
		"percentage2": 25,
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

func TestCacheDataByModel(t *testing.T) {
	t.Skip()

	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	folder := "tmp11"
	cacheFolder := filepath.Join(".", folder, "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", folder))
	}()

	for _, file := range []string{
		"memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		"memberOccupation.27395048-84f4-11f0-9d0e-1e4d4999250c.json",
		"memberDependant.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		"memberDependant.1efda9a6-84f4-11f0-8797-1e4d4999250c.json",
	} {
		src, err := os.Open(filepath.Join(sourceFolder, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(cacheFolder, phoneNumber, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			t.Fatal(err)
			continue
		}

		_, err = os.Stat(dst.Name())
		if os.IsNotExist(err) {
			t.Fatalf("Test failed. Failed to create %s", dst.Name())
		}
	}

	result, err := utils.CacheDataByModel("memberDependant", sessionFolder)
	if err != nil {
		t.Fatal(err)
	}

	target := []map[string]any{
		{
			"data": map[string]any{
				"contact":    "0888777444",
				"name":       "John Phiri",
				"percentage": 10.0,
			},
			"filename": "memberDependant.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		},
		{
			"data": map[string]any{
				"contact":    "07746635653",
				"name":       "Jean Banda",
				"percentage": 5.0,
			},
			"filename": "memberDependant.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		},
	}

	if !reflect.DeepEqual(result, target) {
		t.Fatal("Test failed")
	}

	result, err = utils.CacheDataByModel("memberContact", sessionFolder)
	if err != nil {
		t.Fatal(err)
	}

	target = []map[string]any{
		{
			"data": map[string]any{
				"homeDistrict":             "Lilongwe",
				"homeTraditionalAuthority": "Kalolo",
				"homeVillage":              "Kalulu",
				"phoneNumber":              "0999888777",
				"postalAddress":            "P.O. Box 1",
				"residentialAddress":       "Area 49",
			},
			"filename": "memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		},
	}

	if !reflect.DeepEqual(result, target) {
		t.Fatal("Test failed")
	}
}
