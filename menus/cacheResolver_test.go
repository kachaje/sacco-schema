package menus_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/menus"
	"sacco/utils"
	"testing"
)

func TestResolveCacheDataArray(t *testing.T) {
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "arrayCacheQueries.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "arrayTargetData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.ResolveCacheData(data, "member.memberDependant.0.")

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestResolveCacheDataFlat(t *testing.T) {
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "flatCacheQueries.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.ResolveCacheData(data, "member.")

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "flatTargetData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestResolveCacheDataNestedL1(t *testing.T) {
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedL1CacheQueries.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.ResolveCacheData(data, "member.memberLoan.")

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedL1TargetData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestResolveCacheDataNestedL2(t *testing.T) {
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "fixtures", "sample.flatmap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedL2CacheQueries.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedL2TargetData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &target)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.ResolveCacheData(data, "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.")

	if !utils.MapsEqual(target, result) {
		t.Fatal("Test failed")
	}
}
