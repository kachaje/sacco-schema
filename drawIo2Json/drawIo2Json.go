package drawio2json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/kachaje/sacco-schema/utils"

	"github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"
	"golang.org/x/net/html"
)

type MxCell struct {
	ID       string      `xml:"id,attr"`
	Value    string      `xml:"value,attr"`
	Parent   string      `xml:"parent,attr"`
	Style    string      `xml:"style,attr"`
	Geometry *MxGeometry `xml:"mxGeometry"`
}

type MxGeometry struct {
	X      string `xml:"x,attr"`
	Y      string `xml:"y,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}

type DiagramRoot struct {
	Cells []MxCell `xml:"mxCell"`
}

type MxGraphModel struct {
	Root DiagramRoot `xml:"root"`
}

func Main(filename, configsFolder, targetFolder string) error {
	data, err := D2J(filename)
	if err != nil {
		return err
	}

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configsFolder, "rawData.json"), payload, 0644)
	if err != nil {
		return err
	}

	modelsData, err := ExtractJsonModels(data)
	if err != nil {
		return err
	}

	payload, err = json.MarshalIndent(modelsData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configsFolder, "modelsData.json"), payload, 0644)
	if err != nil {
		return err
	}

	err = CreateYmlFiles(modelsData, targetFolder)
	if err != nil {
		return err
	}

	return nil
}

func D2J(filename string) (map[string]any, error) {
	diagram, err := xml.Parse(filename)
	if err != nil {
		return nil, err
	}

	diagramMap := map[string]any{}
	cellsMap := map[string]any{}

	for _, cell := range diagram.Diagram.MxGraphModel.Root.MxCells {
		cellData := map[string]any{
			"value":  cell.Value,
			"parent": cell.Parent,
		}
		cellsMap[cell.ID] = cellData
	}
	diagramMap["cells"] = cellsMap

	return diagramMap, nil
}

func ValueMapFromString(value string) (map[string]any, error) {
	reLt := regexp.MustCompile("\u003c")
	reGt := regexp.MustCompile("\u003e")

	value = reGt.ReplaceAllLiteralString(value, ">")
	value = reLt.ReplaceAllLiteralString(value, "<")

	doc, err := html.Parse(strings.NewReader(value))
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"model":   "",
		"fields":  map[string]any{},
		"parents": []string{},
	}

	i := 0
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "div":
				parts := strings.Split(n.FirstChild.Data, ":")

				data["model"] = parts[0]

				if len(parts) > 1 {
					val, err := strconv.Atoi(fmt.Sprintf("%v", parts[1]))
					if err == nil {
						data["totalLoops"] = val
						data["hasLoops"] = true
					}
				}
			case "td":
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					re := regexp.MustCompile(`([A-Za-z]+)\s*(\([^\)]+\))`)

					if re.MatchString(c.Data) {
						var field string
						attributes := map[string]any{"order": i}

						parts := re.FindAllStringSubmatch(c.Data, -1)

						if len(parts[0]) > 2 {
							field = parts[0][1]
							attrs := strings.Split(strings.TrimRight(strings.TrimLeft(parts[0][2], "("), ")"), ";")

							attributes["type"] = strings.ToLower(attrs[0])

							for _, v := range attrs[1:] {
								opt := strings.Split(v, ":")

								key := opt[0]

								if slices.Contains([]string{"true", "false"}, opt[1]) {
									vl, err := strconv.ParseBool(opt[1])
									if err == nil {
										attributes[key] = vl
									}
									if key == "many" {
										data["many"] = true
									}
								} else if key == "options" {
									attributes[key] = strings.Split(opt[1], ",")
								} else {
									if key == "default" && regexp.MustCompile(`@`).MatchString(fmt.Sprintf("%v", opt[1])) {
										attributes["dynamicDefault"] = opt[1]
									} else {
										attributes[key] = opt[1]
									}
								}
							}

						}

						data["fields"].(map[string]any)[field] = attributes

						if strings.HasSuffix(field, "Id") {
							model := strings.TrimRight(field, "Id")

							data["fields"].(map[string]any)[field].(map[string]any)["referenceTable"] = model

							data["parents"] = append(data["parents"].([]string), model)
						}

						if field == "id" {
							attributes["primaryKey"] = true
						}

						i++
					}
					f(c)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return data, nil
}

func ExtractJsonModels(rawData map[string]any) (map[string]any, error) {
	models := map[string]any{}

	if cells, ok := rawData["cells"].(map[string]any); ok {
		for _, row := range cells {
			if val, ok := row.(map[string]any); ok {
				if value, ok := val["value"]; ok {
					if vs, ok := value.(string); ok && strings.HasPrefix(vs, "<div") {
						modelData, err := ValueMapFromString(vs)
						if err == nil {
							if model, ok := modelData["model"].(string); ok {
								models[model] = modelData
							}
						}
					}
				}
			}
		}
	}

	return models, nil
}

func CreateYmlFiles(data map[string]any, targetFolder string) error {
	if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
		err := os.MkdirAll(targetFolder, 0755)
		if err != nil {
			return err
		}
	}

	for model, value := range data {
		if val, ok := value.(map[string]any); ok {
			if val["fields"] != nil {
				if fields, ok := val["fields"].(map[string]any); ok {
					keyOrder := map[int]string{}

					for k, v := range fields {
						if vv, ok := v.(map[string]any); ok && vv["order"] != nil {
							index, err := strconv.Atoi(fmt.Sprintf("%v", vv["order"]))
							if err == nil {
								keyOrder[index] = k
							}
						}
					}

					var content string

					for i := range len(keyOrder) {
						key := keyOrder[i]

						rowContent, err := utils.DumpYaml(map[string]any{
							key: fields[key].(map[string]any),
						})
						if err == nil {
							content = fmt.Sprintf("%s\n%s", content, *rowContent)
						}
					}

					err := os.WriteFile(filepath.Join(targetFolder, fmt.Sprintf("%s.yml", model)), []byte(content), 0644)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}

	return nil
}
