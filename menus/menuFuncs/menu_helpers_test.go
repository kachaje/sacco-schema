package menufuncs_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	menufuncs "sacco/menus/menuFuncs"
	"sacco/utils"
	"strings"
	"testing"
)

func TestLoadLoanApplicationForm(t *testing.T) {
	data := map[string]any{}
	templateData := map[string]any{}
	targetData := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "utils", "fixtures", "sample.flatmap.json"))
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

	result := menufuncs.LoadLoanApplicationForm(data, templateData)

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))

	content, err = os.ReadFile(filepath.Join("..", "..", "utils", "fixtures", "loanApplication.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(targetData, result) {
		t.Fatal("Test failed")
	}
}

func TestLoadTemplateData(t *testing.T) {
	t.Skip()

	data := map[string]any{}
	templateData := map[string]any{}
	targetData := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "member.json"))
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

	delete(templateData, "1. OFFICIAL DETAILS")

	content, err = os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "member.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	result := menufuncs.LoadTemplateData(data, templateData)

	if !reflect.DeepEqual(targetData, result) {
		t.Fatal("Test failed")
	}
}

func TestTabulateData(t *testing.T) {
	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "member.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join("..", "..", "database", "fixtures", "member.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	result := menufuncs.TabulateData(data)

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(strings.Join(result, "\n"))

		os.WriteFile(filepath.Join("..", "..", "database", "models", "fixtures", "member.txt"), []byte(strings.Join(result, "\n")), 0644)
	}

	if utils.CleanString(target) != utils.CleanString(strings.Join(result, "\n")) {
		t.Fatal("Test failed")
	}
}
