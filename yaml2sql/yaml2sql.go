package yaml2sql

import (
	"fmt"
	"sacco/utils"
	"strconv"
	"strings"
)

func Main() {

}

func Yml2Sql(model, content string) (*string, error) {
	data, err := utils.LoadYaml(content)
	if err != nil {
		return nil, err
	}

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
			var extras string

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
					extras = strings.TrimSpace(fmt.Sprintf(`%s NOT NULL`, extras))
				}
			}
			if val["primaryKey"] != nil {
				if vPk, ok := val["primaryKey"].(bool); ok && vPk {
					extras = strings.TrimSpace(fmt.Sprintf(`%s PRIMARY KEY`, extras))
				}
			}
			if val["autoIncrement"] != nil {
				if vAutInc, ok := val["autoIncrement"].(bool); ok && vAutInc {
					extras = strings.TrimSpace(fmt.Sprintf(`%s AUTOINCREMENT`, extras))
				}
			}

			row := strings.TrimSpace(fmt.Sprintf(`%s %s %s,`, key, fieldType, extras))

			rows = append(rows, row)
		}
	}

	rows = append(rows, fmt.Sprintf(`active INTEGER DEFAULT 1,
createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

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
