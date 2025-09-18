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

elif [[ "$1" == "-rename" ]]; then

  i=0
  for f in $(ls); do 
    if [[ ! -d $f ]] && [[ ! $f =~ "01" ]]; then 
      ((i++))
      mv $f $(printf "%02d-$f" "$i")
    fi 
  done

elif [[ "$1" == "-t" ]]; then

  pushd tests 2>&1 >/dev/null

  rm -rf tests.log

  for t in $(ls *.go); do
    echo $t 
    go test $t 2>&1 > tests.log
    if [[ $? -ne 0 ]]; then 
      break
    fi
  done

  popd  2>&1 >/dev/null

fi
