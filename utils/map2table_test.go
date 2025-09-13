package utils_test

import (
	"encoding/json"
	"fmt"
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

	result := utils.Map2Table(data)

	fmt.Println(result)
}
