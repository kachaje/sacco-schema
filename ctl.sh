#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  go build -o loadModels ./cmd/gen/*.go

  go build -o dio2json ./cmd/d2j/*.go

elif [[ "$1" == "-dj" ]]; then

  ./dio2json -f designs/sacco.drawio

elif [[ "$1" == "-lm" ]]; then

  ./loadModels 

elif [[ "$1" == "-c" ]]; then

  rm -rf *.db loadModels dio2json

fi
