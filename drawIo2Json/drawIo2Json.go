package drawio2json

import (
	"github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"
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
