#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  go build -o loadModels ./cmd/gen/*.go

  go build -o dio2json ./cmd/d2j/*.go

fi
