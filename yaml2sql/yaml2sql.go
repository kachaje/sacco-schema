package yaml2sql

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sacco/utils"
	"strconv"
	"strings"
)

func Main(folder, targetFile *string) error {
	workingFolder := filepath.Join(".", "models")
	schemaFilename := filepath.Join(".", "schema.sql")

	if folder != nil {
		workingFolder = *folder
	}
	if targetFile != nil {
		schemaFilename = *targetFile
	}

	if _, err := os.Stat(workingFolder); os.IsNotExist(err) {
		return fmt.Errorf("folder %s not found", workingFolder)
	}

	result, err := LoadModels(workingFolder)
	if err != nil {
		return err
	}

	err = os.WriteFile(schemaFilename, []byte(*result), 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadModels(folder string) (*string, error) {
	models := []string{}

	err := filepath.WalkDir(folder, func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		basename := filepath.Base(file)

		if !strings.HasSuffix(basename, ".yml") {
			return nil
		}

		model := strings.TrimSuffix(basename, ".yml")

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		result, err := Yml2Sql(model, string(content))
		if err != nil {
			return err
		}

		models = append(models, *result)

		return nil
	})
	if err != nil {
		return nil, err
	}

	result := strings.Join(models, "\n\n")

	return &result, nil
}

func Yml2Sql(model, content string) (*string, error) {
	data, err := utils.LoadYaml(content)
	if err != nil {
		return nil, err
	}

	footers := []string{`active INTEGER DEFAULT 1,`,
		`createdAt TEXT DEFAULT CURRENT_TIMESTAMP,`,
		`updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,`}

	rows := []string{fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (`, model)}

	keysOrder := map[int]string{}

	for key, value := range data {
		if val, ok := value.(map[string]any); ok && val["order"] != nil {
			index, err := strconv.Atoi(fmt.Sprintf("%v", val["order"]))
			if err == nil {
				keysOrder[index] = key
			}
		}
	}

	for i := range len(keysOrder) {
		if _, ok := keysOrder[i]; !ok {
			continue
		}
		key := keysOrder[i]

		value := data[key]

		if val, ok := value.(map[string]any); ok {
			var fieldType string

			if vType, ok := val["type"]; ok {
				switch vType {
				case "int":
					fieldType = "INTEGER"
				case "real":
					fieldType = "REAL"
				default:
					fieldType = "TEXT"
				}
			}

			if val["required"] != nil {
				if vReq, ok := val["required"].(bool); ok && vReq {
					fieldType = fmt.Sprintf(`%s NOT NULL`, fieldType)
				}
			}
			if val["primaryKey"] != nil {
				if vPk, ok := val["primaryKey"].(bool); ok && vPk {
					fieldType = fmt.Sprintf(`%s PRIMARY KEY`, fieldType)
				}
			}
			if val["autoIncrement"] != nil {
				if vAutInc, ok := val["autoIncrement"].(bool); ok && vAutInc {
					fieldType = strings.TrimSpace(fmt.Sprintf(`%s AUTOINCREMENT`, fieldType))
				}
			}
			if val["referenceTable"] != nil {
				footer := fmt.Sprintf(`FOREIGN KEY (%s) REFERENCES %v (id) ON DELETE CASCADE,`, key, val["referenceTable"])

				footers = append(footers, footer)
			}

			row := fmt.Sprintf(`%s %s,`, key, strings.TrimSpace(fieldType))

			rows = append(rows, row)
		}
	}

	if len(footers) > 0 {
		lastFooter := strings.TrimSuffix(footers[len(footers)-1], ",")

		footers[len(footers)-1] = lastFooter
	}

	rows = append(rows, footers...)

	rows = append(rows, fmt.Sprintf(`);

CREATE TRIGGER IF NOT EXISTS %sUpdated AFTER
UPDATE ON %s FOR EACH ROW BEGIN
UPDATE %s
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;`, model, model, model))

	result := strings.Join(rows, "\n")

	return &result, nil
}
