package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"sacco/utils"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema/schema.sql
var schemaString string

//go:embed schema/seed.sql
var seedString string

//go:embed schema/rates.sql
var ratesString string

//go:embed schema/customTriggers/*
var customTriggerFiles embed.FS

//go:embed schema/triggers/*
var triggerFiles embed.FS

//go:embed schema/configs/models.yml
var modelTemplates string

var modelTemplatesData map[string]any

type Database struct {
	DbName        string
	DB            *sql.DB
	GenericModels map[string]*Model
	SkipFields    []string

	Mu *sync.Mutex
}

func init() {
	var err error

	modelTemplatesData, err = utils.LoadYaml(modelTemplates)
	if err != nil {
		log.Fatal(err)
	}
}

func NewDatabase(dbname string) *Database {
	if dbname != ":memory:" && !strings.HasSuffix(dbname, ".db") {
		dbname = fmt.Sprintf("%s.db", dbname)
	}

	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		log.Fatal(err)
	}

	instance := &Database{
		DbName:        dbname,
		DB:            db,
		GenericModels: map[string]*Model{},
		Mu:            &sync.Mutex{},
		SkipFields: []string{
			"active", "createdAt", "updatedAt", "dateJoined",
		},
	}

	err = instance.initDb()
	if err != nil {
		log.Fatal(err)
	}

	for table, row := range modelTemplatesData {
		if value, ok := row.(map[string]any); ok {
			val, ok := value["fields"].([]any)
			if ok {
				fields := []string{}

				for _, v := range val {
					kv, ok := v.(map[string]any)
					if ok {
						for key := range kv {
							fields = append(fields, key)
						}
					}
				}

				model, err := NewModel(instance.DB, table, fields)
				if err != nil {
					log.Printf("server.database.NewDatabase: %s", err.Error())
					continue
				}

				instance.GenericModels[table] = model
			}
		}
	}

	return instance
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) initDb() error {
	for _, statement := range []string{
		schemaString,
		seedString,
		ratesString,
	} {
		_, err := d.DB.Exec(statement)
		if err != nil {
			return err
		}
	}

	for _, files := range []embed.FS{triggerFiles, customTriggerFiles} {
		err := fs.WalkDir(files, ".", func(path string, folder fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if folder.IsDir() {
				return nil
			}

			content, err := fs.ReadFile(files, path)
			if err != nil {
				return err
			}

			_, err = d.DB.Exec(string(content))
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	for {
		rows, err := d.DB.QueryContext(context.Background(), "SELECT name FROM sqlite_master WHERE type='table'")
		if err == nil {
			count := 0

			for rows.Next() {
				count++
			}

			if count >= 43 {
				time.Sleep(1 * time.Second)

				break
			}
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

func (d *Database) GenericsSaveData(data map[string]any,
	model string,
	retries int,
) (*int64, error) {
	time.Sleep(time.Duration(retries) * time.Second)

	if !d.Mu.TryLock() {
		if retries < 5 {
			retries++

			return d.GenericsSaveData(data, model, retries)
		}

		return nil, fmt.Errorf("server.database.GenericsSaveData: failed to save due to lock error")
	}
	defer d.Mu.Unlock()

	if d.GenericModels[model] == nil {
		return nil, fmt.Errorf("server.database.GenericsSaveData: model %s does not exist", model)
	}

	var id *int64
	var err error

	id, err = d.GenericModels[model].AddRecord(data)
	if err != nil {
		if regexp.MustCompile("UNIQUE").MatchString(err.Error()) {
			if data["id"] != nil {
				if val, ok := data["id"].(int64); ok {
					id = &val
				} else if val, ok := data["id"].(int); ok {
					v := int64(val)
					id = &v
				} else if val, ok := data["id"].(float64); ok {
					v := int64(val)
					id = &v
				} else if val, ok := data["id"].(string); ok {
					v, err := strconv.ParseInt(val, 10, 64)
					if err == nil {
						id = &v
					}
				}
			}

			if id != nil {
				err = d.GenericModels[model].UpdateRecord(data, *id)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("no id found")
			}
		} else {
			return nil, err
		}
	}

	return id, nil
}

func (d *Database) MemberByPhoneNumber(phoneNumber string, skipFields []string) (map[string]any, error) {
	if skipFields == nil {
		skipFields = []string{
			"active", "createdAt", "updatedAt", "dateJoined",
		}
	}

	d.SkipFields = skipFields

	return d.FullMemberRecord(phoneNumber)
}

func (d *Database) SQLQuery(query string) ([]map[string]any, error) {
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

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
			if slices.Contains(d.SkipFields, col) {
				continue
			}

			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func (d *Database) ValidatePassword(username, password string) (*int64, *string, bool) {
	result, err := d.SQLQuery(fmt.Sprintf(`SELECT id, password, userRole FROM user WHERE username = "%v"`, username))
	if err == nil && len(result) > 0 {
		passHash := fmt.Sprintf("%v", result[0]["password"])

		if utils.CheckPasswordHash(password, passHash) {
			var id int64

			val, err := strconv.ParseInt(fmt.Sprintf("%v", result[0]["id"]), 10, 64)
			if err == nil {
				id = val
			}

			userRole := fmt.Sprintf("%v", result[0]["userRole"])

			return &id, &userRole, true
		}
	}

	return nil, nil, false
}

func (d *Database) UsernameFree(username string) bool {
	result, err := d.SQLQuery(fmt.Sprintf(`SELECT id, userRole FROM user WHERE username = "%v"`, username))
	if err == nil && len(result) > 0 {
		return false
	}

	return true
}
