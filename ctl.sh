#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

  rm -rf convert

  go build -o convert ./cmd/gen/*.go

  go build -o cli cmd/wscli/*.go

  go build -o svr cmd/server/*.go

elif [[ "$1" == "-g" ]]; then

  go run cmd/gen/*.go -f designs/sacco.drawio

  pushd menus/workflows 2>&1 >/dev/null

  go run *.go

  popd 2>&1 >/dev/null

  npx prettier -w .

  npx sql-formatter database/schema/schema.sql -l sql --output database/schema/schema.sql

elif [[ "$1" == "-c" ]]; then

  rm -rf *.db* **/*.db* settings/ **/**/settings/ **/**/data/ **/**/*.db* **/tmp*/ *.out

fi
