package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	drawio2json "github.com/kachaje/workflow-parser/drawIo2Json"
	modelgraph "github.com/kachaje/workflow-parser/modelGraph"
	"github.com/kachaje/workflow-parser/yaml2sql"
)

func main() {
	var filename, targetFolder, configsFolder, schemaFilename string

	flag.StringVar(&filename, "f", filename, "draw.io file to process")
	flag.StringVar(&targetFolder, "t", targetFolder, "target destination folder")
	flag.StringVar(&configsFolder, "c", configsFolder, "configs destination folder")
	flag.StringVar(&schemaFilename, "s", schemaFilename, "schema filename")

	flag.Parse()

	if filename == "" {
		flag.Usage()
		return
	}

	if schemaFilename == "" {
		schemaFilename = filepath.Join(".", "database", "schema", "schema.sql")
	}

	if configsFolder == "" {
		configsFolder = filepath.Join(".", "database", "schema", "configs")
	}

	if targetFolder == "" {
		targetFolder = filepath.Join(".", "database", "schema", "models")
	}

	_, err := os.Stat(targetFolder)
	if !os.IsNotExist(err) {
		os.RemoveAll(targetFolder)
	}

	os.MkdirAll(targetFolder, 0755)

	err = drawio2json.Main(filename, configsFolder, targetFolder)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml2sql.Main(&targetFolder, &schemaFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = modelgraph.Main(&configsFolder)
	if err != nil {
		log.Fatal(err)
	}
}
