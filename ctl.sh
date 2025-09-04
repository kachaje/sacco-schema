#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  go build -o loadModels ./cmd/gen/*.go

fi
