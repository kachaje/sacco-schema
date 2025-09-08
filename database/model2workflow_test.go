package database_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/database"
	"sacco/utils"
	"testing"
)

func TestModel2Workflow(t *testing.T) {
	workingFolder := filepath.Join(".", "tmpM2WBasic")

	model := "member"
	srcFile := filepath.Join(".", "database", "schema", "configs", "models.yml")
	dstFile := filepath.Join(workingFolder, fmt.Sprintf("%s.yml", model))

	err := os.MkdirAll(workingFolder, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(workingFolder)
	}()

	content, err := os.ReadFile(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result, _, _, _, err := database.Main(model, dstFile, data)
	if err != nil {
		t.Fatal(err)
	}

	target, err := os.ReadFile(filepath.Join(".", "fixtures", "member.yml"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanString(*result) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}
