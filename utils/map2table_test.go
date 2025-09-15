package utils_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/utils"
	"testing"
)

func TestMap2Table(t *testing.T) {
	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "parser", "fixtures", "schedule.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result := utils.Map2Table(data, []string{"principal", "totalDue"})

	os.WriteFile(filepath.Join(".", "fixtures", "schedule.partial.txt"), []byte(result), 0644)

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "schedule.partial.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatalf(`
Test failed; 
Expected: 
%v
Actual: 
%v
`, target, result)
	}

	result = utils.Map2Table(data, nil)

	os.WriteFile(filepath.Join(".", "fixtures", "schedule.txt"), []byte(result), 0644)

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "schedule.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target = string(content)

	if utils.CleanString(target) != utils.CleanString(result) {
		t.Fatal("Test failed")
	}
}
