package drawio2json_test

import (
	"encoding/json"
	"fmt"
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

func TestValueMapFromString(t *testing.T) {
	value := "\u003cdiv style=\"box-sizing:border-box;width:100%;background:#e4e4e4;padding:2px;\"\u003enotification\u003c/div\u003e\u003ctable style=\"width:100%;font-size:1em;\" cellpadding=\"2\" cellspacing=\"0\"\u003e\u003ctbody\u003e\u003ctr\u003e\u003ctd\u003ePK\u003cbr\u003eFK1\u003cbr\u003e\u003c/td\u003e\u003ctd\u003eid (INT)\u003cbr\u003ememberId (INT)\u003c/td\u003e\u003c/tr\u003e\u003ctr\u003e\u003ctd\u003e\u003cbr\u003e\u003c/td\u003e\u003ctd\u003edate (TEXT)\u003c/td\u003e\u003c/tr\u003e\u003ctr\u003e\u003ctd\u003e\u003c/td\u003e\u003ctd\u003emessage (REAL)\u003cbr\u003edelivered (INT)\u003cbr\u003eread (INT)\u003cbr\u003eactive (INT)\u003cbr\u003ecreatedAt (TIMESTAMP)\u003cbr\u003eupdatedAt (TIMESTAMP)\u003cbr\u003e\u003cbr\u003e\u003c/td\u003e\u003c/tr\u003e\u003c/tbody\u003e\u003c/table\u003e"

	result, err := drawio2json.ValueMapFromString(value)
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}
