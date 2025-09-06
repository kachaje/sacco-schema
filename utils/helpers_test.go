package utils_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sacco/utils"
	"testing"
)

func TestLoadYaml(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "newMember.yml"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	refData, err := os.ReadFile(filepath.Join(".", "fixtures", "newMember.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(refData, &target)
	if err != nil {
		t.Fatal(err)
	}

	compareObjects := func(obj1, obj2 map[string]any) bool {
		if len(obj1) != len(obj2) {
			return false
		}

		for key, val1 := range obj1 {
			val2, exists := obj2[key]
			if !exists || fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
				return false
			}
		}

		return true
	}

	if !compareObjects(target, result) {
		t.Fatal("Test failed")
	}
}

func TestLockFile(t *testing.T) {
	rootFolder := filepath.Join(".", "tmpFileLock")

	os.MkdirAll(rootFolder, 0755)
	defer func() {
		os.RemoveAll(rootFolder)
	}()

	filename := filepath.Join(rootFolder, "lock.txt")

	err := os.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(filename)
	}()

	lockFilename, err := utils.LockFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(lockFilename)
	}()

	_, err = os.Stat(lockFilename)
	if os.IsNotExist(err) {
		t.Fatal("Test failed")
	}
}

func TestUnLockFile(t *testing.T) {
	rootFolder := filepath.Join(".", "tmpFileUnLock")

	os.MkdirAll(rootFolder, 0755)
	defer func() {
		os.RemoveAll(rootFolder)
	}()

	filename := filepath.Join(rootFolder, "lock.txt")
	lockFilename := fmt.Sprintf("%s.lock", filename)

	err := os.WriteFile(lockFilename, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(lockFilename)
	}()

	err = utils.UnLockFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(lockFilename)
	if !os.IsNotExist(err) {
		t.Fatal("Test failed")
	}
}

func TestFileLocked(t *testing.T) {
	rootFolder := filepath.Join(".", "tmpFileLocked")

	os.MkdirAll(rootFolder, 0755)
	defer func() {
		os.RemoveAll(rootFolder)
	}()

	filename := filepath.Join(rootFolder, "lock.txt")
	lockFilename := fmt.Sprintf("%s.lock", filename)

	err := os.WriteFile(lockFilename, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := os.Stat(lockFilename); !os.IsNotExist(err) {
			os.Remove(lockFilename)
		}
	}()

	locked := utils.FileLocked(filename)
	if !locked {
		t.Fatalf("Test failed. Expected: true; Actual: %v", locked)
	}

	err = os.Remove(lockFilename)
	if err != nil {
		t.Fatal(err)
	}

	locked = utils.FileLocked(filename)
	if locked {
		t.Fatalf("Test failed. Expected: false; Actual: %v", locked)
	}
}

func TestIdentifierToLabel(t *testing.T) {
	result := utils.IdentifierToLabel("thisIsAString")

	target := "This Is A String"

	if result != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, result)
	}
}

func TestIndex(t *testing.T) {
	numbers := []int{10, 20, 30, 40, 50}

	result := utils.Index(numbers, 30)

	if result != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %d", result)
	}

	result = utils.Index(numbers, 60)

	if result != -1 {
		t.Fatalf("Test failed. Expected: -1; Actual: %d", result)
	}
}

func TestUnpackData(t *testing.T) {
	data := map[string]any{}
	target := []map[string]any{}

	for i := range 4 {
		row := map[string]any{}

		for _, key := range []string{"id", "name", "value"} {
			label := fmt.Sprintf("%s%d", key, i+1)
			value := fmt.Sprintf("%s%d", key, i+1)

			data[label] = value
			row[key] = value
		}

		target = append(target, row)
	}

	result := utils.UnpackData(data)

	if len(result) != len(target) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", len(target), len(result))
	}

	if len(result[0]) != len(target[0]) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", len(target[0]), len(result[0]))
	}

	data = map[string]any{
		"id":    "1",
		"name":  "test",
		"value": "something",
	}
	target = []map[string]any{}

	target = append(target, data)

	result = utils.UnpackData(data)

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestGetSkippedRefIds(t *testing.T) {
	refData := []map[string]any{
		{
			"contact":    "P.O. Box 2",
			"id":         2,
			"memberId":   1,
			"name":       "Benefator 2",
			"percentage": 8,
		},
		{
			"contact":    "P.O. Box 3",
			"id":         3,
			"memberId":   1,
			"name":       "Benefator 3",
			"percentage": 5,
		},
		{
			"contact":    "P.O. Box 4",
			"id":         4,
			"memberId":   1,
			"name":       "Benefator 4",
			"percentage": 2,
		},
		{
			"contact":    "P.O. Box 1",
			"id":         1,
			"memberId":   1,
			"name":       "Benefator 1",
			"percentage": 10,
		},
	}
	data := []map[string]any{
		{
			"contact":    "P.O. Box 5678",
			"id":         2,
			"memberId":   1,
			"name":       "Benefator 2",
			"percentage": 25,
		},
		{
			"contact":    "P.O. Box 1234",
			"id":         1,
			"memberId":   1,
			"name":       "Benefator 1",
			"percentage": 35,
		},
	}

	result := utils.GetSkippedRefIds(data, refData)

	target := []map[string]any{
		{
			"contact":    "P.O. Box 3",
			"id":         3,
			"memberId":   1,
			"name":       "Benefator 3",
			"percentage": 5},
		{
			"contact":    "P.O. Box 4",
			"id":         4,
			"memberId":   1,
			"name":       "Benefator 4",
			"percentage": 2,
		},
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}
