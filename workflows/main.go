package main

import (
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sacco/database"
	"sacco/utils"
	"sort"
	"strings"
)

func buildFuncs() {
	folder := filepath.Join("..", "menus", "menufuncs")

	rows := []string{}

	err := filepath.WalkDir(folder, func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		basename := filepath.Base(file)

		if !strings.HasSuffix(basename, "_fn.go") {
			return nil
		}

		fnName := strings.TrimSuffix(basename, "_fn.go")

		cFnName := utils.CapitalizeFirstLetter(fnName)

		row := fmt.Sprintf(`FunctionsMap["%s"] = %s`, fnName, cFnName)

		rows = append(rows, row)

		return nil
	})
	if err != nil {
		panic(err)
	}

	content := fmt.Sprintf(`package menufuncs

import (
	"sacco/database"
	"sacco/parser"
	"time"
)

var (
	DB       *database.Database
	Sessions = map[string]*parser.Session{}

	WorkflowsData = map[string]map[string]any{}

	DemoMode bool

	FunctionsMap = map[string]func(
		func(
			string, *parser.Session,
			string, string, string,
		) string,
		map[string]any,
		*parser.Session,
	) string{}

	ReRouteRemaps = map[string]any{}

	RefDate = time.Now().Format("2006-01-02")
)

func init() {
%s
}
`, strings.Join(rows, "\n"))

	filename := filepath.Join("..", "menus", "menufuncs", "menufuncs.go")

	rawData, err := format.Source([]byte(content))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, rawData, 0644)
	if err != nil {
		panic(err)
	}
}

func buildWorkflows() {
	workingFolder := filepath.Join("..", "menus", "workflows")

	_, err := os.Stat(workingFolder)
	if !os.IsNotExist(err) {
		os.RemoveAll(workingFolder)
	}

	err = os.MkdirAll(workingFolder, 0755)
	if err != nil {
		log.Panic(err)
	}

	content, err := os.ReadFile(filepath.Join("..", "database", "schema", "configs", "models.yml"))
	if err != nil {
		log.Panic(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		log.Panic(err)
	}

	relationships := map[string]any{}
	parentModels := []string{}
	floatFields := map[string]bool{}

	for model := range data {
		targetFile := filepath.Join(workingFolder, fmt.Sprintf("%s.yml", model))

		_, row, floats, parents, err := database.Main(model, targetFile, data)
		if err != nil {
			log.Panic(err)
		}

		if len(parents) > 0 {
			values := []string{}

			for _, key := range parents {
				values = append(values, fmt.Sprintf(`"%s",`, key))
			}

			entry := fmt.Sprintf(`"%s": {
			%s
		},`, model, strings.Join(values, "\n"))

			parentModels = append(parentModels, entry)
		}

		for key := range floats {
			if !floatFields[key] {
				floatFields[key] = true
			}
		}

		if len(row) > 0 {
			relationships[model] = row
		}
	}

	script := []string{}

	arraysGroup := []string{}
	singlesGroup := []string{}

	for key, value := range relationships {
		model := utils.CapitalizeFirstLetter(key)

		if val, ok := value.(map[string][]string); ok {
			if len(val["hasMany"]) > 0 {
				rows := []string{}

				for _, v := range val["hasMany"] {
					rows = append(rows, fmt.Sprintf(`"%s"`, v))
				}

				row := strings.TrimSpace(fmt.Sprintf(`%sArrayChildren = []string{
				%s,
				}`, model, strings.Join(rows, ",\n")))

				if len(row) > 0 {
					script = append(script, row)

					arraysGroup = append(arraysGroup, fmt.Sprintf(`"%sArrayChildren": %sArrayChildren,`, model, model))
				}
			}
			if len(val["hasOne"]) > 0 {
				rows := []string{}

				for _, v := range val["hasOne"] {
					rows = append(rows, fmt.Sprintf(`"%s"`, v))
				}

				row := strings.TrimSpace(fmt.Sprintf(`%sSingleChildren = []string{
				%s,
				}
				`, model, strings.Join(rows, ",\n")))

				if len(row) > 0 {
					script = append(script, row)

					singlesGroup = append(singlesGroup, fmt.Sprintf(`"%sSingleChildren": %sSingleChildren,`, model, model))
				}
			}
		}
	}

	floatKeys := []string{}

	for key := range floatFields {
		floatKeys = append(floatKeys, fmt.Sprintf(`"%s",`, key))
	}

	sort.Strings(script)
	sort.Strings(singlesGroup)
	sort.Strings(arraysGroup)
	sort.Strings(floatKeys)
	sort.Strings(parentModels)

	targetName := filepath.Join("..", "database", "models.go")

	content, err = format.Source(fmt.Appendf(nil, `package database
	
	var (
	%s
	SingleChildren = map[string][]string{
		%s
	}
	ArrayChildren = map[string][]string{
		%s
	}
	FloatFields = []string{
		%s
	}
	ParentModels = map[string][]string{
		%s
	}
	)`,
		strings.Join(script, "\n"),
		strings.Join(singlesGroup, "\n"),
		strings.Join(arraysGroup, "\n"),
		strings.Join(floatKeys, "\n"),
		strings.Join(parentModels, "\n"),
	))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(targetName, content, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	buildFuncs()

	buildWorkflows()
}
