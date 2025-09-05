#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  rm -rf convert

  go build -o convert ./cmd/gen/*.go

elif [[ "$1" == "-g" ]]; then

  ./convert -f designs/sacco.drawio

  npx prettier -w .

  npx sql-formatter schema/schema.sql -l sql --output schema/schema.sql

elif [[ "$1" == "-c" ]]; then

  rm -rf *.db **/*.db

fi
