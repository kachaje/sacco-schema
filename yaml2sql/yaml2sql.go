package yaml2sql

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sacco/utils"
	"strconv"
	"strings"
)

func Main(folder, targetFile *string) error {
	workingFolder := filepath.Join(".", "database", "schema", "models")
	schemaFilename := filepath.Join(".", "database", "schema", "schema.sql")

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

			fieldData := fieldType

			if val["scheduleFormula"] != nil {
			} else if val["primaryKey"] != nil {
				if vPk, ok := val["primaryKey"].(bool); ok && vPk {
					fieldData = fmt.Sprintf(`%s PRIMARY KEY`, fieldData)
				}
			} else if val["default"] != nil {
				if regexp.MustCompile(`@`).MatchString(fmt.Sprintf("%v", val["default"])) {
				} else if fmt.Sprintf(`%v`, val["default"]) == "CURRENT_USER" {
				} else {
					fieldData = fmt.Sprintf(`%s DEFAULT %v`, fieldData, val["default"])
				}
			} else if val["optional"] == nil {
				fieldData = fmt.Sprintf(`%s NOT NULL`, fieldData)
			}
			if val["autoIncrement"] != nil {
				if vAutInc, ok := val["autoIncrement"].(bool); ok && vAutInc {
					fieldData = strings.TrimSpace(fmt.Sprintf(`%s AUTOINCREMENT`, fieldData))
				}
			}
			if val["options"] != nil {
				if v, ok := val["options"].([]any); ok {
					opts := []string{}

					for _, k := range v {
						opts = append(opts, fmt.Sprintf(`'%v'`, k))
					}

					fieldData = strings.TrimSpace(fmt.Sprintf(`%s CHECK (%s IN (%s))`, fieldData, key, strings.Join(opts, ", ")))
				}
			}
			if val["unique"] != nil {
				fieldData = fmt.Sprintf(`%s UNIQUE`, fieldData)
			}
			if val["referenceTable"] != nil {
				footer := fmt.Sprintf(`FOREIGN KEY (%s) REFERENCES %v (id) ON DELETE CASCADE,`, key, val["referenceTable"])

				footers = append(footers, footer)
			}

			row := fmt.Sprintf(`%s %s,`, key, strings.TrimSpace(fieldData))

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
  updatedAt = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;`, model, model, model))

	result := strings.Join(rows, "\n")

	return &result, nil
}
