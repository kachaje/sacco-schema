package menufuncs_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	menufuncs "sacco/menus/menuFuncs"
	"sacco/utils"
	"strings"
	"testing"
)

var (
	refDate = "2025-09-11"
)

func TestLoadGroupMembers(t *testing.T) {
	data := map[string]any{}
	targetData := []map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.LoadGroupMembers(data, "memberDependant")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join("..", "..", "database", "fixtures", "memberDependants.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "memberDependants.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(targetData, result) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", targetData, result)
	}
}

func TestResolveNestedQuery(t *testing.T) {
	data := map[string]any{
		"member.memberLoan.3.memberLoanWitness.5.name": "Mary Banda",
	}

	result := menufuncs.ResolveNestedQuery(data, "member.memberLoan.0.memberLoanWitness.0.name")

	target := "member.memberLoan.3.memberLoanWitness.5.name"

	if target != result {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", target, result)
	}
}

func TestLoadTemplateData(t *testing.T) {
	data := map[string]any{}
	templateData := map[string]any{}
	targetData := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "templates", "member.template.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &templateData)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.LoadTemplateData(data, templateData, &refDate)

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join("..", "fixtures", "member.template.output.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join("..", "fixtures", "member.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(targetData, result) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", targetData, result)
	}
}

func TestTabulateData(t *testing.T) {
	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "fixtures", "member.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.TabulateData(data)

	if os.Getenv("DEBUG") == "true" {
		payload := []byte(strings.Join(result, "\n"))

		os.WriteFile(filepath.Join("..", "fixtures", "member.txt"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join("..", "fixtures", "member.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if utils.CleanString(target) != utils.CleanString(strings.Join(result, "\n")) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}

func TestLoadLoanApplicationForm(t *testing.T) {
	data := map[string]any{}
	templateData := map[string]any{}
	targetData := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "templates", "loanApplication.template.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &templateData)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.LoadTemplateData(data, templateData, &refDate)

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join("..", "fixtures", "loanApplication.template.output.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join("..", "fixtures", "loanApplication.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.MapsEqual(targetData, result) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", targetData, result)
	}
}

func TestBusinessSummary(t *testing.T) {
	data := map[string]any{}
	templateData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "templates", "businessSummary.template.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &templateData)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.LoadTemplateData(data, templateData, &refDate)

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join("..", "fixtures", "businessSummary.template.output.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join("..", "fixtures", "businessSummary.template.output.json"))
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

func TestTabulateBusinessSummary(t *testing.T) {
	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "fixtures", "businessSummary.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.TabulateData(data)

	if os.Getenv("DEBUG") == "true" {
		payload := []byte(strings.Join(result, "\n"))

		os.WriteFile(filepath.Join("..", "fixtures", "businessSummary.txt"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join("..", "fixtures", "businessSummary.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if utils.CleanString(target) != utils.CleanString(strings.Join(result, "\n")) {
		t.Fatalf("Test failed; Expected: %#v; Actual: %#v", target, result)
	}
}
