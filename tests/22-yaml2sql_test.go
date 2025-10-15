package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kachaje/sacco-schema/yaml2sql"
	"github.com/kachaje/utils/utils"
)

func TestYml2Sql(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "models", "productRate.yml"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := yaml2sql.Yml2Sql("productRate", string(content))
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "models", "productRate.sql"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if utils.CleanString(target) != utils.CleanString((*result)) {
		t.Fatalf("Test failed; Expected: %s; Actual: %s", target, *result)
	}
}

func TestLoadModels(t *testing.T) {
	result, err := yaml2sql.LoadModels(filepath.Join(".", "fixtures", "models"))
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	if utils.CleanString(target) != utils.CleanString((*result)) {
		t.Fatalf("Test failed; Expected: %s; Actual: %s", target, *result)
	}
}
