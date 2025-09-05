package drawio2json_test

import (
	"path/filepath"
	"reflect"
	drawio2json "sacco/drawIo2Json"
	"testing"
)

func TestDefault(t *testing.T) {
	result, err := drawio2json.Main(filepath.Join(".", "fixtures", "diagram.xml"))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{
		"cells": map[string]any{
			"kVijt7gfVD9ZtySMmpSK-1": map[string]any{
				"parent": "1",
				"value":  "myReceiver",
			},
		},
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}
