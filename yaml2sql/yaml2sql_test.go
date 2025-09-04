package yaml2sql_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/utils"
	"sacco/yaml2sql"
	"testing"
)

func TestYml2Sql(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "models", "loanServiceFee.yml"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := yaml2sql.Yml2Sql("loanServiceFee", string(content))
	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("Test failed")
	}

	content, err = os.ReadFile(filepath.Join(".", "fixtures", "models", "loanServiceFee.sql"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	fmt.Println(*result)

	if utils.CleanString(target) != utils.CleanString((*result)) {
		t.Fatal("Test failed")
	}
}
