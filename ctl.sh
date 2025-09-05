#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  go build -o loadModels ./cmd/gen/*.go

  go build -o dio2json ./cmd/d2j/*.go

elif [[ "$1" == "-dj" ]]; then

  ./dio2json -f designs/sacco.drawio > drawIo2Json/fixtures/data.json 

elif [[ "$1" == "-lm" ]]; then

  ./loadModels 

fi
