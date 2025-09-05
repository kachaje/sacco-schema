package drawio2json

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

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

func Main(filename string) (map[string]any, error) {
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
		"model":  "",
		"fields": map[string]any{},
	}

	i := 0
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "div":
				data["model"] = n.FirstChild.Data
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
								} else if key == "options" {
									attributes[key] = strings.Split(opt[1], ",")
								} else {
									attributes[key] = opt[1]
								}
							}

						}

						data["fields"].(map[string]any)[field] = attributes

						if strings.HasSuffix(field, "Id") {
							model := strings.TrimRight(field, "Id")

							data["fields"].(map[string]any)[field].(map[string]any)["referenceTable"] = model
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
