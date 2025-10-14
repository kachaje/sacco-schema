package database

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/kachaje/sacco-schema/utils"

	"github.com/google/uuid"
)

type Model struct {
	ModelName string
	Fields    []string

	db *sql.DB
}

func NewModel(
	db *sql.DB,
	modelName string,
	fields []string,
) (*Model, error) {
	if modelName == "" {
		return nil, fmt.Errorf("missing required modelName")
	}

	m := &Model{
		db:        db,
		ModelName: modelName,
		Fields:    fields,
	}

	return m, nil
}

func (m *Model) AddRecord(data map[string]any) (*int64, error) {
	fields := []string{}
	values := []any{}
	markers := []string{}
	var id int64

	if data["id"] != nil {
		val, err := strconv.ParseInt(fmt.Sprintf("%v", data["id"]), 10, 64)
		if err == nil {
			id = val

			err := m.UpdateRecord(data, id)
			if err != nil && err.Error() == "no match found" {
			} else {
				return &id, err
			}
		}
	}

	if m.ModelName == "member" {
		if data["memberIdNumber"] == nil {
			memberIdNumber := strings.ToUpper(
				regexp.MustCompile(`[^A-Za-z0-9]`).
					ReplaceAllLiteralString(uuid.NewString(), ""),
			)

			data["memberIdNumber"] = memberIdNumber
			data["shortMemberId"] = memberIdNumber[:8]
		}

		data["dateJoined"] = time.Now().Format("2006-01-02")
	}

	for key, value := range data {
		if !slices.Contains(m.Fields, key) {
			continue
		}

		fields = append(fields, key)
		markers = append(markers, "?")

		if strings.ToLower(key) == "password" {
			password, err := utils.HashPassword(fmt.Sprintf("%v", value))
			if err != nil {
				return nil, err
			}

			values = append(values, password)
		} else {
			values = append(values, value)
		}
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		m.ModelName,
		strings.Join(fields, ", "),
		strings.Join(markers, ", "),
	)

	result, err := utils.QueryWithRetry(m.db, context.Background(), 0, query, values...)
	if err != nil {
		return nil, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return nil, err
	}

	return &id, nil
}

func (m *Model) UpdateRecord(data map[string]any, id int64) error {
	fields := []string{}
	values := []any{}

	for key, value := range data {
		if !slices.Contains(append(m.Fields, []string{"updatedAt"}...), key) {
			continue
		}

		fields = append(fields, fmt.Sprintf("%s = ?", key))

		if strings.ToLower(key) == "password" {
			password, err := utils.HashPassword(fmt.Sprintf("%v", value))
			if err != nil {
				return err
			}

			values = append(values, password)
		} else {
			values = append(values, value)
		}
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE %s SET %s WHERE id=?", m.ModelName, strings.Join(fields, ", "))

	result, err := utils.QueryWithRetry(
		m.db,
		context.Background(), 0,
		statement, values...,
	)
	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); err == nil && count <= 0 {
		return fmt.Errorf("no match found")
	}

	return nil
}

func (m *Model) loadRows(rows *sql.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]any, len(cols))
	scanArgs := make([]any, len(cols))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	results := []map[string]any{}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range cols {
			val := values[i]
			if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 && val != nil {
				if b, ok := val.([]byte); ok {
					rowMap[col] = string(b)
				} else {
					rowMap[col] = val
				}
			}
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func (m *Model) FetchById(id int64) (map[string]any, error) {
	rows, err := m.db.Query(fmt.Sprintf(`SELECT * FROM %s WHERE active=1 AND id=?`, m.ModelName), id)
	if err != nil {
		return nil, err
	}

	result, err := m.loadRows(rows)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return nil, nil
}

func (m *Model) FilterBy(whereStatement string) ([]map[string]any, error) {
	if !regexp.MustCompile("active").MatchString(whereStatement) {
		whereStatement = fmt.Sprintf("%s AND active=1", whereStatement)
	}

	rows, err := m.db.Query(fmt.Sprintf(`SELECT * FROM %s %s`, m.ModelName, whereStatement))
	if err != nil {
		return nil, err
	}

	return m.loadRows(rows)
}
