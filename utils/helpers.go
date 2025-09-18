package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

type DiffResult struct {
	Added   map[string]any
	Removed map[string]any
	Changed map[string]any
}

func CleanScript(content []byte) string {
	s := string(content)

	return CleanString(s)
}

func CleanString(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return s
}

func DumpYaml(data map[string]any) (*string, error) {
	var result string

	payload, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	result = string(payload)

	return &result, nil
}

func LoadYaml(yamlData string) (map[string]any, error) {
	var data map[string]any

	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WaitForPort(host string, port string, timeout time.Duration, retryInterval time.Duration, debug bool) error {
	address := net.JoinHostPort(host, port)
	startTime := time.Now()

	for {
		conn, err := net.DialTimeout("tcp", address, retryInterval)
		if err == nil {
			conn.Close()
			if debug {
				fmt.Printf("Port %s on %s is open.\n", port, host)
			}
			return nil
		}

		if time.Since(startTime) >= timeout {
			return fmt.Errorf("timeout waiting for port %s on %s: %w", port, host, err)
		}

		if debug {
			fmt.Printf("Waiting for port %s on %s... Retrying in %v\n", port, host, retryInterval)
		}

		time.Sleep(retryInterval)
	}
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "0.0.0.0:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func LockFile(filename string) (string, error) {
	lockFilename := fmt.Sprintf("%s.lock", filename)

	return lockFilename, os.WriteFile(lockFilename, []byte{}, 0644)
}

func UnLockFile(filename string) error {
	lockFilename := fmt.Sprintf("%s.lock", filename)

	return os.Remove(lockFilename)
}

func FileLocked(filename string) bool {
	lockFilename := fmt.Sprintf("%s.lock", filename)

	_, err := os.Stat(lockFilename)

	return !os.IsNotExist(err)
}

func QueryWithRetry(db *sql.DB, ctx context.Context, retries int, query string, args ...any) (sql.Result, error) {
	time.Sleep(time.Duration(retries) * time.Second)

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		if regexp.MustCompile("SQL logic error: no such table").MatchString(err.Error()) {
			if retries < 3 {
				retries++

				log.Printf("models.QueryWithRetry.retry: %d\n", retries)

				return QueryWithRetry(db, ctx, retries, query, args...)
			}
		}
		payload, _ := json.MarshalIndent(args, "", "  ")

		return nil, fmt.Errorf(`utils.QueryWithRetry.1: %s
query: %s
data: %s`, err.Error(), query, payload)
	}

	return result, nil
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func IdentifierToLabel(identifier string) string {
	re := regexp.MustCompile("([A-Z][a-z]*)")

	parts := re.FindAllString(CapitalizeFirstLetter(identifier), -1)

	return strings.Join(parts, " ")
}

func Index[T comparable](s []T, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CenterString(s string, width int) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= width {
		return s
	}

	padding := width - runeCount
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return strings.Repeat(" ", leftPadding) + s + strings.Repeat(" ", rightPadding)
}

func CacheDataByModel(filterModel, sessionFolder string) ([]map[string]any, error) {
	result := []map[string]any{}

	err := filepath.WalkDir(sessionFolder, func(fullpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		filename := filepath.Base(fullpath)

		re := regexp.MustCompile(`\.[a-z0-9-]+\.json$`)

		if !re.MatchString(filename) {
			return nil
		}

		model := re.ReplaceAllLiteralString(filename, "")

		if model != filterModel {
			return nil
		}

		content, err := os.ReadFile(fullpath)
		if err != nil {
			return err
		}

		if data := map[string]any{}; json.Unmarshal(content, &data) == nil {
			result = append(result, map[string]any{
				"data":     data,
				"filename": filename,
			})
		} else if data := []map[string]any{}; json.Unmarshal(content, &data) == nil {
			rows := []map[string]any{}

			for _, row := range data {
				rows = append(rows, map[string]any{
					"data":     row,
					"filename": filename,
				})
			}

			result = append(result, rows...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UnpackData(data map[string]any) []map[string]any {
	result := []map[string]any{}
	rows := map[string]map[string]any{}

	for key, value := range data {
		re := regexp.MustCompile(`^(.+)(\d+)$`)

		if re.MatchString(key) {
			parts := re.FindAllStringSubmatch(key, -1)[0]

			field := parts[1]
			index := parts[2]

			if rows[index] == nil {
				rows[index] = map[string]any{}
			}

			rows[index][field] = value
		} else {
			if len(rows) == 0 {
				rows["1"] = data
				break
			}
		}
	}

	for _, row := range rows {
		result = append(result, row)
	}

	return result
}

func GetSkippedRefIds(data, refData []map[string]any) []map[string]any {
	result := []map[string]any{}

	for _, row := range refData {
		if row["id"] != nil {
			id := fmt.Sprintf("%v", row["id"])
			found := false

			for _, child := range data {
				if child["id"] != nil && id == fmt.Sprintf("%v", child["id"]) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, row)
			}
		}
	}

	return result
}

func CacheFile(filename string, data any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if FileLocked(filename) {
		if retries < 5 {
			retries++

			CacheFile(filename, data, retries)
			return
		}
	}
	_, err := LockFile(filename)
	if err != nil {
		log.Printf("server.Cachefile.1: %s", err.Error())
		retries = 0

		CacheFile(filename, data, retries)
		return
	}
	defer func() {
		err := UnLockFile(filename)
		if err != nil {
			log.Printf("server.Cachefile.2: %s", err.Error())
		}
	}()

	payload, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		err = os.WriteFile(filename, payload, 0644)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(err)
	}
}

func GetMapDiff(map1, map2 map[string]any) DiffResult {
	diff := DiffResult{
		Added:   make(map[string]any),
		Removed: make(map[string]any),
		Changed: make(map[string]any),
	}

	for key, val1 := range map1 {
		if val2, ok := map2[key]; !ok {
			diff.Removed[key] = val1
		} else {
			if nestedMap1, isMap1 := val1.(map[string]any); isMap1 {
				if nestedMap2, isMap2 := val2.(map[string]any); isMap2 {
					nestedDiff := GetMapDiff(nestedMap1, nestedMap2)
					if len(nestedDiff.Added) > 0 || len(nestedDiff.Removed) > 0 || len(nestedDiff.Changed) > 0 {
						diff.Changed[key] = nestedDiff
					}
				} else {
					diff.Changed[key] = map[string]any{
						"old":     val1,
						"new":     val2,
						"oldType": reflect.TypeOf(val1).String(),
						"newType": reflect.TypeOf(val2).String(),
					}
				}
			} else if !reflect.DeepEqual(val1, val2) {
				diff.Changed[key] = map[string]any{
					"old":     val1,
					"new":     val2,
					"oldType": reflect.TypeOf(val1).String(),
					"newType": reflect.TypeOf(val2).String(),
				}
			}
		}
	}

	for key, val2 := range map2 {
		if _, ok := map1[key]; !ok {
			diff.Added[key] = val2
		}
	}

	return diff
}

func MapsEqual(m1, m2 map[string]any) bool {
	if (m1 == nil) != (m2 == nil) {
		return false
	}
	if m1 == nil && m2 == nil {
		return true
	}
	if len(m1) != len(m2) {
		return false
	}

	for key, val1 := range m1 {
		val2, ok := m2[key]
		if !ok {
			return false
		}

		switch v1 := val1.(type) {
		case map[string]any:
			v2, ok := val2.(map[string]any)
			if !ok || !MapsEqual(v1, v2) {
				return false
			}
		default:
			if fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
				return false
			}
		}
	}
	return true
}

func extractNumber(s string) (int, string, bool) {
	var numStr string
	for i, r := range s {
		if r >= '0' && r <= '9' {
			numStr += string(r)
		} else {
			if numStr != "" {
				num, err := strconv.Atoi(numStr)
				if err == nil {
					return num, s[i:], true
				}
			}
			return 0, s, false
		}
	}
	if numStr != "" {
		num, err := strconv.Atoi(numStr)
		if err == nil {
			return num, "", true
		}
	}
	return 0, s, false
}

func SortSlice(data []string) {
	sort.Slice(data, func(i, j int) bool {
		s1 := data[i]
		s2 := data[j]

		parts1 := strings.FieldsFunc(s1, func(r rune) bool { return r == '_' || r == '.' || (r < '0' || r > '9') })
		parts2 := strings.FieldsFunc(s2, func(r rune) bool { return r == '_' || r == '.' || (r < '0' || r > '9') })

		for k := 0; k < len(parts1) && k < len(parts2); k++ {
			num1, _, isNum1 := extractNumber(parts1[k])
			num2, _, isNum2 := extractNumber(parts2[k])

			if isNum1 && isNum2 {
				if num1 != num2 {
					return num1 < num2
				}
			} else if parts1[k] != parts2[k] {
				return parts1[k] < parts2[k]
			}
		}

		return len(s1) < len(s2)
	})
}
