package main

import (
	"flag"
	"log"
	"path/filepath"
	drawio2json "sacco/drawIo2Json"
	"sacco/yaml2sql"
)

func main() {
	var filename, targetFolder, schemaFilename string

	flag.StringVar(&filename, "f", filename, "draw.io file to process")
	flag.StringVar(&targetFolder, "t", targetFolder, "target destination folder")
	flag.StringVar(&schemaFilename, "s", schemaFilename, "schema filename")

	flag.Parse()

	if filename == "" {
		flag.Usage()
		return
	}

	if schemaFilename == "" {
		schemaFilename = filepath.Join(".", "schema", "schema.sql")
	}

	if targetFolder == "" {
		targetFolder = filepath.Join(".", "schema", "models")
	}

	err := drawio2json.Main(filename, targetFolder)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml2sql.Main(&targetFolder, &schemaFilename)
	if err != nil {
		log.Fatal(err)
	}
}
