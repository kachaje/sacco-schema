# Getting Started with sacco-schema

This guide will help you get up and running with sacco-schema in minutes. You'll learn how to install, configure, and use the SACCO management system.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.24.3 or later** installed ([Download Go](https://go.dev/dl/))
- **Git** for cloning the repository
- **SQLite** (included with Go, no separate installation needed)
- Basic understanding of YAML and SQL
- (Optional) **Draw.io** for database design visualization

## Installation

### Option 1: Clone from Repository

```bash
git clone https://github.com/kachaje/sacco-schema.git
cd sacco-schema
go mod download
```

### Option 2: Use as Dependency

If you're building on top of sacco-schema:

```bash
go get github.com/kachaje/sacco-schema
```

## Quick Start

### 1. Build the Application

Use the provided build script:

```bash
chmod +x ctl.sh
./ctl.sh -b
```

This creates three executables:
- `svr`: Server application
- `cli`: WebSocket client for testing
- `convert`: Code generation tool

### 2. Run the Server

Start the server with default settings:

```bash
./svr
```

Or with custom options:

```bash
./svr -p 8080 -n sacco.db -d -o
```

**Flags**:
- `-p`: Server port (default: auto-assigned free port)
- `-n`: Database name (default: `:memory:`)
- `-d`: Enable dev mode (shows dev menus)
- `-o`: Enable demo mode (bypasses authentication)

### 3. Access the Application

**Web Interface**:
Open your browser to `http://localhost:8080` (or the port shown in logs)

**USSD Simulation**:
Use the WebSocket client:

```bash
./cli -p 8080 -n 1234567890
```

**Direct USSD**:
Send HTTP POST to `/ussd`:
```bash
curl -X POST http://localhost:8080/ussd \
  -d "sessionId=test123" \
  -d "phoneNumber=1234567890" \
  -d "text="
```

## Configuration

### Database Setup

The database is automatically initialized on first run. The system:

1. Creates all tables from `database/schema/schema.sql`
2. Seeds initial data from `database/schema/seed.sql`
3. Loads rate configurations from `database/schema/rates.sql`
4. Creates triggers from `database/schema/triggers/`

**Using Persistent Database**:

```bash
./svr -n production.db
```

This creates `production.db` file that persists data.

**Using In-Memory Database**:

```bash
./svr -n :memory:
```

Data is lost when server stops (useful for testing).

### Menu Configuration

Menus are defined in `menus/menus/*.yml`:

```yaml
title: Welcome to Kaso SACCO
fields:
  "1":
    id: registration
    label:
      en: Membership Application
  "2":
    id: loan
    label:
      en: Loans
```

Edit these files to customize menu structure.

### Workflow Configuration

Workflows are in `menus/workflows/*.yml`. Each workflow defines:

- Screens and navigation
- Field validation
- Data collection
- Formula calculations

See [Workflow Design Documentation](./workflow-design.md) for details.

## Common Use Cases

### Use Case 1: Member Registration

**Via USSD**:

1. Start server: `./svr -o` (demo mode)
2. Connect client: `./cli -n 1234567890`
3. Navigate:
   - Select `1` (Membership Application)
   - Select `1` (Member Details)
   - Fill in member information
   - Submit with `0`

**Via Web**:

1. Open `http://localhost:8080`
2. Enter phone number
3. Follow same navigation as USSD

**Programmatically**:

```go
package main

import (
    "fmt"
    "github.com/kachaje/sacco-schema/database"
    "github.com/kachaje/sacco-schema/menus"
    "github.com/kachaje/workflow-parser/parser"
)

func main() {
    // Initialize database
    db := database.NewDatabase(":memory:")
    defer db.Close()
    
    // Initialize menus
    demoMode := true
    m := menus.NewMenus(nil, &demoMode)
    
    // Create session
    phoneNumber := "1234567890"
    session := parser.NewSession(nil, &phoneNumber, nil, nil)
    
    // Load main menu
    response := m.LoadMenu("main", session, phoneNumber, "", "")
    fmt.Println(response)
}
```

### Use Case 2: Loan Application

**Workflow**:

1. User selects "Loans" from main menu
2. Selects "Loan Application"
3. Enters loan details:
   - Loan purpose
   - Loan amount
   - Repayment period
   - Loan type
4. System generates payment schedule
5. User reviews and confirms
6. Loan saved to database

**Admin Approval**:

1. Admin selects "Loan Approvals"
2. Reviews loan application
3. Approves, denies, or partially approves
4. System updates loan status

### Use Case 3: Contribution Deposit

**Workflow**:

1. User selects contribution menu
2. Selects "Deposit"
3. Enters contribution number (or selects from list)
4. Enters amount
5. Confirms transaction
6. System updates contribution balance

### Use Case 4: Running Reports

**Loans Report**:

```go
import "github.com/kachaje/sacco-schema/reports"

r := reports.NewReports(db)
report, err := r.LoansReport("2024-01-31")
if err != nil {
    log.Fatal(err)
}

table, err := r.LoansReport2Table(*report)
if err != nil {
    log.Fatal(err)
}

fmt.Println(*table)
```

**Contributions Report**:

```go
report, err := r.ContributionsReport("2024-01-31")
if err != nil {
    log.Fatal(err)
}

table, err := r.ContributionsReport2Table(*report, []string{
    "memberName",
    "contributionId",
    "memberTotal",
})
fmt.Println(*table)
```

## Command-Line Tools

### Code Generation Tool (`cmd/gen`)

Generate database schema from Draw.io diagram:

```bash
go run cmd/gen/main.go -f designs/sacco.drawio
```

Or use the script:

```bash
./ctl.sh -g
```

**Process**:
1. Parses Draw.io XML
2. Extracts model definitions
3. Generates YAML model files
4. Creates SQL schema
5. Generates workflow configurations

**Output Files**:
- `database/schema/models/*.yml`: Model definitions
- `database/schema/schema.sql`: SQL schema
- `database/schema/configs/*.json`: Configuration files

### WebSocket Client (`cmd/wscli`)

Test USSD workflows via WebSocket:

```bash
./cli -p 8080 -n 1234567890
```

**Options**:
- `-p`: Server port (default: 8080)
- `-n`: Phone number to use (default: 1234567890)
- `-s`: Silent mode (suppress prompts)

### Server (`cmd/server`)

Main server application:

```bash
./svr [options]
```

**Options**:
- `-p`: Port number (default: auto)
- `-n`: Database name (default: `:memory:`)
- `-d`: Dev mode (shows dev menus)
- `-o`: Demo mode (bypasses auth)

## API Examples

### USSD Endpoint

**POST** `/ussd`

**Parameters**:
- `sessionId`: Session identifier
- `phoneNumber`: User's phone number
- `text`: User input
- `serviceCode`: USSD service code (optional)

**Example**:

```bash
curl -X POST http://localhost:8080/ussd \
  -d "sessionId=abc123" \
  -d "phoneNumber=1234567890" \
  -d "text=1"
```

**Response**: Plain text USSD format

```
CON Welcome to Kaso SACCO
1. Membership Application
2. Loans
3. Check Balance
...
```

### WebSocket Endpoint

**GET** `/ws?phoneNumber=1234567890&sessionId=abc123`

**Protocol**: WebSocket text messages

**Example**:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?phoneNumber=1234567890');
ws.onmessage = (event) => {
    console.log('Response:', event.data);
};
ws.send('1'); // Select menu item 1
```

### Cron Jobs Endpoint

**POST** `/cron/jobs`

**Body**:

```json
{
  "targetDate": "2024-01-31",
  "profit": 100000
}
```

**Example**:

```bash
curl -X POST http://localhost:8080/cron/jobs \
  -H "Content-Type: application/json" \
  -d '{"targetDate": "2024-01-31", "profit": 100000}'
```

**Response**: `Done` on success

### Ledger API

**POST** `/api/transaction`

**Body**:

```json
{
  "name": "Loan Disbursement",
  "description": "Disburse loan LN001",
  "ledgerEntries": [
    {
      "referenceNumber": "LN001",
      "name": "Loan Account",
      "description": "Principal",
      "debitCredit": "DEBIT",
      "amount": 50000,
      "accountId": 1,
      "accountType": "ASSET"
    },
    {
      "referenceNumber": "LN001",
      "name": "Member Savings",
      "description": "Deposit",
      "debitCredit": "CREDIT",
      "amount": 50000,
      "accountId": 2,
      "accountType": "LIABILITY"
    }
  ]
}
```

**GET** `/api/transaction?startDate=2024-01-01&endDate=2024-01-31`

Returns account balances for date range.

## Database Operations

### Direct Database Access

```go
import "github.com/kachaje/sacco-schema/database"

db := database.NewDatabase("sacco.db")
defer db.Close()

// Generic save
data := map[string]any{
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "1234567890",
}
id, err := db.GenericsSaveData(data, "member", 0)

// Query
results, err := db.SQLQuery("SELECT * FROM member WHERE active=1")

// Model operations
model := db.GenericModels["member"]
record, err := model.FetchById(*id)
```

### Model CRUD

```go
// Create
data := map[string]any{
    "firstName": "Jane",
    "lastName": "Smith",
}
id, err := model.AddRecord(data)

// Update
updateData := map[string]any{
    "firstName": "Jane Updated",
}
err := model.UpdateRecord(updateData, *id)

// Read
record, err := model.FetchById(*id)

// Query
records, err := model.FilterBy("WHERE firstName LIKE 'J%'")
```

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Specific Test

```bash
go test ./tests/12-menu_test.go
```

### Run Tests with Script

```bash
./ctl.sh -t
```

### Test Coverage

```bash
go test -cover ./...
```

## Troubleshooting

### Server Won't Start

**Issue**: Port already in use
- **Solution**: Use `-p` flag to specify different port
- **Solution**: Kill process using the port

**Issue**: Database file locked
- **Solution**: Close other connections
- **Solution**: Use `:memory:` for testing

### Workflows Not Loading

**Issue**: YAML syntax error
- **Solution**: Validate YAML syntax
- **Solution**: Check workflow file exists in `menus/workflows/`

**Issue**: Menu not found
- **Solution**: Verify menu file in `menus/menus/`
- **Solution**: Check menu ID matches

### Data Not Saving

**Issue**: Missing saveFunc
- **Solution**: Ensure database initialized
- **Solution**: Check model exists in database

**Issue**: Validation errors
- **Solution**: Check validation rules in workflow
- **Solution**: Verify required fields provided

### Session Issues

**Issue**: ActiveData not updating
- **Solution**: Call `session.RefreshSession()`
- **Solution**: Check database queries return data

**Issue**: Cache not populating
- **Solution**: Verify cache query format
- **Solution**: Check session has ActiveData

### WebSocket Connection Fails

**Issue**: Connection refused
- **Solution**: Verify server running
- **Solution**: Check port number matches

**Issue**: Messages not received
- **Solution**: Check server logs
- **Solution**: Verify WebSocket upgrade successful

## Development Workflow

### 1. Modify Database Schema

1. Edit Draw.io diagram (`designs/sacco.drawio`)
2. Run code generator: `./ctl.sh -g`
3. Review generated SQL: `database/schema/schema.sql`
4. Test with in-memory database

### 2. Add New Workflow

1. Create workflow YAML: `menus/workflows/myWorkflow.yml`
2. Add menu item referencing workflow
3. Test workflow via client
4. Add tests

### 3. Add New Menu Function

1. Create function: `menus/menuFuncs/myFunction_fn.go`
2. Register in `menufuncs.FunctionsMap`
3. Add menu item with `function: myFunction`
4. Test and add tests

### 4. Modify Reports

1. Edit report function in `reports/`
2. Test with sample data
3. Verify table formatting
4. Add edge case tests

## Best Practices

### 1. Database

- Use transactions for multi-step operations
- Always check for errors
- Use parameterized queries
- Close connections properly

### 2. Workflows

- Set `order` field for all screens
- Use validation rules for user input
- Mark optional fields explicitly
- Use cache queries for editing

### 3. Menus

- Use descriptive menu IDs
- Add role-based access where needed
- Group related items together
- Provide clear labels

### 4. Testing

- Test with in-memory database
- Use table-driven tests
- Test edge cases
- Test error scenarios

### 5. Performance

- Use indexes for frequent queries
- Cache menu configurations
- Limit session data size
- Use pagination for large datasets

## Next Steps

Now that you've learned the basics:

1. **Read the Architecture Documentation** - Understand system internals
   - See [Architecture Documentation](./architecture.md)

2. **Learn Workflow Design** - Deep dive into workflow patterns
   - See [Workflow Design Documentation](./workflow-design.md)

3. **Explore Examples** - Review test fixtures
   - See `tests/fixtures/` directory

4. **Build Your Features** - Add custom workflows and functions
   - Start with simple workflows
   - Add menu functions for complex logic
   - Integrate with external systems

## Additional Resources

- **workflow-parser Documentation**: See `workflow-parser/docs/`
- **Utils Documentation**: See `utils/docs/`
- **Test Fixtures**: See `tests/fixtures/`
- **Schema Files**: See `database/schema/`

## Summary

You've learned:

- ✅ How to install and run sacco-schema
- ✅ How to configure database and menus
- ✅ How to use common workflows
- ✅ How to access APIs
- ✅ How to develop new features
- ✅ Troubleshooting tips

You're now ready to build and deploy your SACCO management system!

