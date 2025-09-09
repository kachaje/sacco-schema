package menus_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/menus"
	"sacco/utils"
	"testing"
)

func TestResolveCacheData(t *testing.T) {
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "data.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "cacheQueries.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "targetData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.ResolveCacheData(data, cacheData)

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}
