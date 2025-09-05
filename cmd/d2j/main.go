package main

import (
	"flag"
	"log"
	"path/filepath"
	drawio2json "sacco/drawIo2Json"
)

func main() {
	var filename, targetFolder string

	flag.StringVar(&filename, "f", filename, "draw.io file to process")
	flag.StringVar(&targetFolder, "t", targetFolder, "target destination folder")

	flag.Parse()

	if filename == "" {
		flag.Usage()
		return
	}

	if targetFolder == "" {
		targetFolder = filepath.Join(".", "models")
	}

	err := drawio2json.Main(filename, targetFolder)
	if err != nil {
		log.Fatal(err)
	}
}
