#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  go build -o convert ./cmd/gen/*.go

elif [[ "$1" == "-lm" ]]; then

  ./convert -f designs/sacco.drawio

elif [[ "$1" == "-c" ]]; then

  rm -rf *.db loadModels dio2json convert

fi
