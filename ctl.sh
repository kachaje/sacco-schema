#!/usr/bin/env bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if gotestsum is installed
check_gotestsum() {
    if command -v gotestsum > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

if [[ "$1" == "-b" ]]; then

  rm -rf convert

  go build -o convert ./cmd/gen/*.go

  go build -o cli cmd/wscli/*.go

  go build -o svr cmd/server/*.go

elif [[ "$1" == "-g" ]]; then

  go run cmd/gen/*.go -f designs/sacco.drawio

  npx prettier -w .

  npx sql-formatter database/schema/schema.sql -l sql --output database/schema/schema.sql

elif [[ "$1" == "-c" ]]; then

  rm -rf *.db* **/*.db* settings/ **/**/settings/ **/**/data/ **/**/*.db* **/tmp*/ *.out *.log **/*.log

elif [[ "$1" == "-rename" ]]; then

  i=0
  for f in $(ls); do 
    if [[ ! -d $f ]] && [[ ! $f =~ "01" ]]; then 
      ((i++))
      mv $f $(printf "%02d-$f" "$i")
    fi 
  done

elif [[ "$1" == "-t" ]]; then

  echo -e "${GREEN}Running tests...${NC}"
  
  # Use gotestsum if available for pytest-like progress output
  if check_gotestsum; then
    # gotestsum provides progress output similar to pytest
    # --format testname shows test names as they run (like pytest)
    # Use TEST_FORMAT env var to override: dots, testname, standard-verbose, etc.
    TEST_FORMAT="${TEST_FORMAT:-testname}"
    # Use script to preserve colors, then filter EMPTY lines
    script -q /dev/null bash -c "gotestsum --format \"$TEST_FORMAT\" -- --count=1 ./tests/..." 2>&1 | grep -v --color=always 'EMPTY' 
  else
    # Fallback to go test with verbose output
    echo -e "${YELLOW}Note: Install gotestsum for better progress output:${NC}"
    echo -e "${YELLOW}  go install gotest.tools/gotestsum@latest${NC}"
    echo ""
    # Use script to preserve colors, then filter EMPTY lines
    script -q /dev/null bash -c "go test -v ./tests/..." 2>&1 | grep -v --color=always 'EMPTY' 
  fi

elif [[ "$1" == "-l" ]]; then

  echo "Running go vet..."
  go vet ./...

  if command -v golangci-lint &> /dev/null; then
    echo ""
    echo "Running golangci-lint..."
    golangci-lint run
  else
    echo ""
    echo "golangci-lint not found. Install it with:"
    echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
  fi

fi
