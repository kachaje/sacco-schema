package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	drawio2json "sacco/drawIo2Json"
)

func main() {
	var filename string

	flag.StringVar(&filename, "f", filename, "draw.io file to process")

	flag.Parse()

	if filename == "" {
		flag.Usage()
		return
	}

	result, err := drawio2json.Main(filename)
	if err != nil {
		log.Fatal(err)
	}

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(payload))
}
