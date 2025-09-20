package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

var content []byte
var err error
var (
	data     = map[string]any{}
	loanData = map[string]any{}
)

func init() {
	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "memberLoanPayment.json"))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &loanData)
	if err != nil {
		panic(err)
	}
}

func setupLoanEnv() (*parser.WorkFlow, *string, error) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"))
	if err != nil {
		return nil, nil, err
	}

	activeData := map[string]any{}

	err = json.Unmarshal(content, &activeData)
	if err != nil {
		return nil, nil, err
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "loanNumber.txt"))
	if err != nil {
		return nil, nil, err
	}

	loanNumber := string(content)

	phoneNumber := "0999999999"

	session := parser.Session{
		ActiveData: activeData,
	}

	sessions := map[string]*parser.Session{
		phoneNumber: &session,
	}

	wf := parser.NewWorkflow(loanData, nil, nil, &phoneNumber, nil, nil, nil, sessions, nil)

	return wf, &loanNumber, nil
}

func TestLoadAjaxOptions(t *testing.T) {
	wf, loanNumber, err := setupLoanEnv()
	if err != nil {
		t.Fatal(err)
	}

	result, keys := wf.LoadAjaxOptions("memberLoanPaymentSchedule@loanNumber", *loanNumber, "dueDate", []string{
		"interest",
		"processingFee",
		"insurance",
		"instalment",
	})

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "sample.ajaxOptions.json"), payload, 0644)

		payload, _ = json.MarshalIndent(keys, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "sample.ajaxOptions.keys.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.ajaxOptions.json"))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatalf(`Test failed.
Expected: 
%v
Actual: 
%v`, target, result)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.ajaxOptions.keys.json"))
	if err != nil {
		t.Fatal(err)
	}

	targetKeys := []string{}

	err = json.Unmarshal(content, &targetKeys)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(targetKeys, keys) {
		t.Fatalf(`Test failed.
Expected: 
%v
Actual: 
%v`, targetKeys, keys)
	}
}

func TestAjaxOptions(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	result, err := wf.NextNode("")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}

func TestGetNode(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	result := wf.GetNode("enterLanguage")

	if result == nil {
		t.Fatal("Test failed")
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}
}

func TestInputIncluded(t *testing.T) {
	targetRoute := "enterOtherName"

	options := []any{
		map[string]any{
			"position": 1,
			"label": map[string]any{
				"en": "Yes",
				"ny": "Inde",
			},
			"nextScreen": targetRoute,
		},
		map[string]any{
			"position": 2,
			"label": map[string]any{
				"en": "No",
				"ny": "Ayi",
			},
			"nextScreen": "enterGender",
		},
	}

	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	defaultRoute := "enterAskOtherName"

	wf.CurrentScreen = defaultRoute

	result, nextRoute := wf.InputIncluded("3", options)

	if result {
		t.Fatalf("Test failed. Expected: false; Actual: %v", result)
	}
	if nextRoute != "" {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", defaultRoute, nextRoute)
	}

	result, nextRoute = wf.InputIncluded("1", options)

	if !result {
		t.Fatalf("Test failed. Expected: true; Actual: %v", result)
	}
	if nextRoute != targetRoute {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", targetRoute, nextRoute)
	}

	wf.CurrentScreen = defaultRoute

	result, nextRoute = wf.InputIncluded("2", options)

	if !result {
		t.Fatalf("Test failed. Expected: true; Actual: %v", result)
	}
	if nextRoute != "enterGender" {
		t.Fatalf("Test failed. Expected: enterGender; Actual: %s", nextRoute)
	}
}

func TestNodeOptions(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	result := wf.NodeOptions("enterLanguage")

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. English", "2. Chichewa"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}

	wf.CurrentLanguage = "2"

	result = wf.NodeOptions("enterMaritalStatus")

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. Inde", "2. Ayi", "3. Woferedwa", "4. Osudzulidwa"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}

	wf.CurrentLanguage = "en"

	result = wf.NodeOptions("enterMaritalStatus")

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. Married", "2. Single", "3. Widowed", "4. Divorced"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}
}

func TestNextNode(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	result, err := wf.NextNode("")
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "initialScreen" {
		t.Fatalf("Test failed. Expected: 'initialScreen'; Actual: '%v'", wf.PreviousScreen)
	}

	result, err = wf.NextNode("3")
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "initialScreen" {
		t.Fatalf("Test failed. Expected: 'initialScreen'; Actual: '%v'", wf.PreviousScreen)
	}

	wf.CurrentScreen = "enterDateOfBirth"

	wf.NextNode("1999")

	if wf.CurrentScreen != "enterDateOfBirth" {
		t.Fatalf("Test failed. Expected: 'enterDateOfBirth'; Actual: '%v'", wf.CurrentScreen)
	}

	wf.NextNode("1999-09-01")

	if wf.CurrentScreen != "enterMaritalStatus" {
		t.Fatalf("Test failed. Expected: 'enterMaritalStatus'; Actual: '%v'", wf.CurrentScreen)
	}

	wf.CurrentScreen = "enterLanguage"

	result, err = wf.NextNode("1")
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"type", "text", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterFirstName" {
		t.Fatalf("Test failed. Expected: 'enterFirstName'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.PreviousScreen)
	}

	if wf.Data["dateOfBirth"] == nil || fmt.Sprintf("%v", wf.Data["dateOfBirth"]) != "1999-09-01" {
		t.Fatalf("Test failed. Expected: '1999-09-01'; Actual: %v", wf.Data["dateOfBirth"])
	}

	if wf.Data["language"] == nil || fmt.Sprintf("%v", wf.Data["language"]) != "1" {
		t.Fatalf("Test failed. Expected: '1'; Actual: %v", wf.Data["dateOfBirth"])
	}
}

func TestOptionValue(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	wf.CurrentLanguage = "2"

	options := []any{
		map[string]any{
			"position": 1,
			"code":     "y",
			"label": map[string]any{
				"en": "Yes",
				"ny": "Inde",
			},
			"nextScreen": "",
		},
		map[string]any{
			"position": 2,
			"label": map[string]any{
				"all": "No",
			},
			"nextScreen": "enterGender",
		},
	}

	result, code := wf.OptionValue(options, "2")

	if result != "No" {
		t.Fatalf("Test failed. Expected: No; Actual: %v", result)
	}

	if *code != "" {
		t.Fatalf("Test failed. Got: %v", *code)
	}

	result, code = wf.OptionValue(options, "1")

	if *code != "y" {
		t.Fatalf("Test failed. Expected: y; Actual: %v", *code)
	}

	if result != "Yes" {
		t.Fatalf("Test failed. Expected: Yes; Actual: %v", result)
	}
}

func TestResolveData(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	result := wf.ResolveData(map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	}, false)

	target := map[string]any{
		"language":      "English",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "No",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "Single",
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	for key, val := range target {
		if result[key] == nil || fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", val) {
			t.Fatalf("Test failed. Expected: %v; Actual: %v", val, result[key])
		}
	}
}

func TestLoadLabel(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	target := "Language"

	result := wf.LoadLabel("language")

	if result != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, result)
	}
}

func TestGetLabel(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	node, err := wf.NextNode("")
	if err != nil {
		t.Fatal(err)
	}

	if node == nil {
		t.Fatal("Test failed")
	}

	target := "Language: 1. English 2. Chichewa 99. Cancel"

	result := wf.GetLabel(node, wf.CurrentScreen)

	label := utils.CleanScript([]byte(result))

	if label != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, label)
	}

	wf.CurrentLanguage = "2"

	target = "Chiyankhulo: 1. English 2. Chichewa 99. Basi"

	result = wf.GetLabel(node, wf.CurrentScreen)

	label = utils.CleanScript([]byte(result))

	if label != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, label)
	}

	wf.Data = map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	}

	wf.CurrentScreen = "formSummary"

	node = wf.GetNode(wf.CurrentScreen)

	target = `Zomwe Mwalemba
- Chiyankhulo: English
- Dzina Loyamba: Mary
- Dzina La Abambo: Banda
- Dzina Lina?: No
- Tsiku Lobadwa: 1999-09-01
- Muli M'banja: Single

0. Zatheka
00. Tiyambirenso
98. Bwererani
99. Basi
`

	result = wf.GetLabel(node, wf.CurrentScreen)

	if utils.CleanScript([]byte(target)) != utils.CleanScript([]byte(result)) {
		t.Fatalf(`Test failed. 
Expected: %v; 
Actual: %v`, target, result)
	}
}

func TestGotoMenu(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	wf.Data = map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	}

	node, err := wf.NextNode("formSummary")
	if err != nil {
		t.Fatal(err)
	}

	if node == nil {
		t.Fatal("Test failed")
	}

	wf.NextNode("00")

	if len(wf.Data) != 0 {
		t.Fatalf("Test failed. Expected: 0; Actual: %v", len(wf.Data))
	}

	target := "initialScreen"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target, wf.CurrentScreen)
	}
}

func TestCancel(t *testing.T) {
	called := false

	wf := parser.NewWorkflow(data, func(m any, model, phoneNumber, preferenceFolder *string, saveFunc func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
		if m != nil {
			t.Fatalf("Test failed. Expected: nil; Actual: %v", m)
		}

		called = true
		return nil
	}, nil, nil, nil, nil, nil, nil, nil)

	wf.Data = map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	}

	node, err := wf.NextNode("formSummary")
	if err != nil {
		t.Fatal(err)
	}

	if node == nil {
		t.Fatal("Test failed")
	}

	node, err = wf.NextNode("99")
	if err != nil {
		t.Fatal(err)
	}

	if node != nil {
		t.Fatal("Test failed")
	}

	if len(wf.Data) != 0 {
		t.Fatalf("Test failed. Expected: 0; Actual: %v", len(wf.Data))
	}

	target := "initialScreen"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target, wf.CurrentScreen)
	}

	if called {
		t.Fatal("Test failed")
	}
}

func TestSubmit(t *testing.T) {
	called := false

	wf := parser.NewWorkflow(data, func(m any, model, phoneNumber, preferenceFolder *string, saveFunc func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
		if m == nil {
			t.Fatalf("Test failed")
		}

		val, ok := m.(map[string]any)
		if ok {
			if len(val) != 6 {
				t.Fatalf("Test failed. Expected: 6; Actual: %v", len(val))
			}

			called = true
		}
		return nil
	}, nil, nil, nil, nil, nil, nil, nil)

	wf.Data = map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	}

	wf.CurrentScreen = "formSummary"

	node, err := wf.NextNode("0")
	if err != nil {
		t.Fatal(err)
	}

	if node != nil {
		t.Fatal("Test failed")
	}

	if len(wf.Data) != 0 {
		t.Fatalf("Test failed. Expected: 0; Actual: %v", len(wf.Data))
	}

	target := "initialScreen"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", target, wf.CurrentScreen)
	}

	if !called {
		t.Fatal("Test failed")
	}
}

func TestNavNext(t *testing.T) {
	called := false

	wf := parser.NewWorkflow(data, func(m any, model, phoneNumber, preferenceFolder *string, saveFunc func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
		if m == nil {
			t.Fatalf("Test failed")
		}

		val, ok := m.(map[string]any)
		if ok {
			if len(val) != 6 {
				t.Fatalf("Test failed. Expected: 6; Actual: %v", len(val))
			}

			called = true
		}
		return nil
	}, nil, nil, nil, nil, nil, nil, nil)

	target := `Language: 
1. English
2. Chichewa
99. Cancel`

	result := wf.NavNext("")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	target = "enterLanguage"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	target = `Dzina Loyamba: 
00. Tiyambirenso
98. Bwererani
99. Basi`

	result = wf.NavNext("2")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	if wf.CurrentLanguage != "2" {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", wf.CurrentLanguage)
	}

	target = "enterFirstName"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	target = `Dzina La Abambo: 
00. Tiyambirenso
98. Bwererani
99. Basi`

	result = wf.NavNext("Mary")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	target = "enterLastName"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	target = `Dzina Lina?: 
1. Inde
2. Ayi
00. Tiyambirenso
98. Bwererani
99. Basi`

	result = wf.NavNext("Banda")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	target = "enterAskOtherName"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	target = `Tsiku Lobadwa: 
00. Tiyambirenso
98. Bwererani
99. Basi`

	result = wf.NavNext("2")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	target = "enterDateOfBirth"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	target = `Muli M'banja: 
1. Inde
2. Ayi
3. Woferedwa
4. Osudzulidwa
00. Tiyambirenso
98. Bwererani
99. Basi`

	result = wf.NavNext("1999-09-01")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	target = "enterMaritalStatus"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	target = `Zomwe Mwalemba
- Chiyankhulo: Chichewa
- Dzina Loyamba: Mary
- Dzina La Abambo: Banda
- Dzina Lina?: No
- Tsiku Lobadwa: 1999-09-01
- Muli M'banja: Single

0. Zatheka
00. Tiyambirenso
98. Bwererani
99. Basi`

	result = wf.NavNext("2")

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`Test failed.
Expected: %s
Actual: %s`, target, result)
	}

	target = "formSummary"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	wf.NavNext("0")

	if !called {
		t.Fatal("Test failed")
	}
}

func TestBack(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	wf.NavNext("")

	wf.NavNext("1")

	target := "enterFirstName"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}

	node, err := wf.NextNode("98")
	if err != nil {
		t.Fatal(err)
	}

	if node == nil {
		t.Fatal("Test failed")
	}

	target = "enterLanguage"

	if wf.CurrentScreen != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, wf.CurrentScreen)
	}
}

func TestEvalCondition(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	data = map[string]any{
		"loanStatus": "APPROVED",
	}

	result := wf.EvalCondition("loanStatus=APPROVED", data)

	if !result {
		t.Fatalf("Test failed. Expecting: true; Actual: %v", result)
	}

	result = wf.EvalCondition("loanStatus=IN[REJECTED,PARTIAL-APPROVAL]", data)

	if result {
		t.Fatalf("Test failed. Expecting: false; Actual: %v", result)
	}

	data = map[string]any{
		"loanStatus": "PARTIAL-APPROVAL",
	}

	result = wf.EvalCondition("loanStatus=IN[REJECTED,PARTIAL-APPROVAL]", data)

	if !result {
		t.Fatalf("Test failed. Expecting: true; Actual: %v", result)
	}
}

func TestLoadDynaDefault(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil)

	data = map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result := wf.LoadDynaDefault("memberLoan@loanNumber", data)

	if result == nil {
		t.Fatal("Test failed")
	}

	if os.Getenv("DEBUG") == "true" {
		payload := fmt.Appendf(nil, "%v", result)

		os.WriteFile(filepath.Join(".", "fixtures", "loanNumber.txt"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "loanNumber.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if result.(string) != target {
		t.Fatalf("Test failed. Expecting: %s; Actual: %v", target, result)
	}

	result = wf.LoadDynaDefault("memberLoan@loanAmount", data)

	if result == nil {
		t.Fatal("Test failed")
	}

	target = "200000"

	if fmt.Sprintf("%v", result) != target {
		t.Fatalf("Test failed. Expecting: %s; Actual: %v", target, result)
	}

	result = wf.LoadDynaDefault("memberLoan@nonExistent", data)

	if result != nil {
		t.Fatal("Test failed")
	}
}
