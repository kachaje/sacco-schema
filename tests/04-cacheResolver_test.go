package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/sacco-schema/menus"
	"github.com/kachaje/utils/utils"
)

func TestResolveCacheDataArray(t *testing.T) {
	t.Parallel()
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"))
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

	result := menus.ResolveCacheData(data, "member.memberDependant.0.")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "arrayTargetData.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "arrayTargetData.json"))
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

func TestResolveCacheDataFlat(t *testing.T) {
	t.Parallel()
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"))
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

	delete(result, "memberIdNumber")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "flatTargetData.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "flatTargetData.json"))
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

func TestResolveCacheDataNestedL1(t *testing.T) {
	t.Parallel()
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"))
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

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "nestedL1TargetData.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedL1TargetData.json"))
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

func TestResolveCacheDataNestedL2(t *testing.T) {
	t.Parallel()
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "sample.data.flatmap.json"))
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

	result := menus.ResolveCacheData(data, "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "nestedL2TargetData.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedL2TargetData.json"))
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

func TestNestedChild(t *testing.T) {
	t.Parallel()
	data := map[string]any{}
	cacheData := map[string]any{}
	target := map[string]any{}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "nestedChildSourceData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedChildQueries.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &cacheData)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.ResolveCacheData(data, "member.memberLoan.0.memberLoanLiability.0.")

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(result, "", "  ")

		os.WriteFile(filepath.Join(".", "fixtures", "nestedChildTargetData.json"), payload, 0644)
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "nestedChildTargetData.json"))
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
