package drawio2json_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	drawio2json "sacco/drawIo2Json"
	"sacco/utils"
	"testing"
)

func TestValueMapFromString(t *testing.T) {
	value := "\u003cdiv style=\"box-sizing:border-box;width:100%;background:#e4e4e4;padding:2px;\"\u003enotification\u003c/div\u003e\u003ctable style=\"width:100%;font-size:1em;\" cellpadding=\"2\" cellspacing=\"0\"\u003e\u003ctbody\u003e\u003ctr\u003e\u003ctd\u003ePK\u003cbr\u003eFK1\u003cbr\u003e\u003c/td\u003e\u003ctd\u003eid (INT;autoIncrement:true)\u003cbr\u003ememberId (INT;many:true)\u003c/td\u003e\u003c/tr\u003e\u003ctr\u003e\u003ctd\u003e\u003cbr\u003e\u003c/td\u003e\u003ctd\u003edate (TEXT;default:CURRENT_TIMESTAMP)\u003c/td\u003e\u003c/tr\u003e\u003ctr\u003e\u003ctd\u003e\u003c/td\u003e\u003ctd\u003emessage (TEXT)\u003cbr\u003emsgDelivered (TEXT;options:Yes,No;default:No;optional:true)\u003cbr\u003emsgRead (TEXT;options:Yes,No;default:No;optional:true)\u003cbr\u003e\u003cbr\u003e\u003c/td\u003e\u003c/tr\u003e\u003c/tbody\u003e\u003c/table\u003e"

	result, err := drawio2json.ValueMapFromString(value)
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{
		"fields": map[string]any{
			"date": map[string]any{
				"default": "CURRENT_TIMESTAMP",
				"order":   2,
				"type":    "text",
			},
			"id": map[string]any{
				"autoIncrement": true,
				"order":         0,
				"primaryKey":    true,
				"type":          "int",
			},
			"memberId": map[string]any{
				"many":           true,
				"order":          1,
				"referenceTable": "member",
				"type":           "int",
			},
			"message": map[string]any{
				"order": 3,
				"type":  "text",
			},
			"msgDelivered": map[string]any{
				"default":  "No",
				"optional": true,
				"options": []string{
					"Yes",
					"No",
				},
				"order": 4,
				"type":  "text",
			},
			"msgRead": map[string]any{
				"default":  "No",
				"optional": true,
				"options": []string{
					"Yes",
					"No",
				},
				"order": 5,
				"type":  "text",
			},
		},
		"many":    true,
		"model":   "notification",
		"parents": []string{"member"},
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestD2J(t *testing.T) {
	result, err := drawio2json.D2J(filepath.Join(".", "fixtures", "diagram.xml"))
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

func TestExtractJsonModels(t *testing.T) {
	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "schema", "configs", "rawData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	result, err := drawio2json.ExtractJsonModels(data)
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	content, err = os.ReadFile(filepath.Join("..", "database", "schema", "configs", "modelsData.json"))
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

func TestCreateYmlFiles(t *testing.T) {
	folder := filepath.Join(".", "tmp")

	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "schema", "configs", "modelsData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	err = drawio2json.CreateYmlFiles(data, folder)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			os.RemoveAll(folder)
		}
	}()

	for key := range data {
		if _, err := os.Stat(filepath.Join(folder, fmt.Sprintf("%s.yml", key))); os.IsNotExist(err) {
			t.Fatalf("Test failed on %s", key)
		}
	}
}

func TestDefault(t *testing.T) {
	folder := filepath.Join(".", "tmpMain")

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		t.Fatal(err)
	}

	srcFile, err := os.Open(filepath.Join("..", "database", "schema", "configs", "rawData.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filepath.Join(folder, "rawData.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	err = drawio2json.Main(filepath.Join("..", "designs", "sacco.drawio"), folder, folder)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			os.RemoveAll(folder)
		}
	}()

	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "schema", "configs", "modelsData.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	for key := range data {
		if _, err := os.Stat(filepath.Join(folder, fmt.Sprintf("%s.yml", key))); os.IsNotExist(err) {
			t.Fatalf("Test failed on %s", key)
		}
	}
}
