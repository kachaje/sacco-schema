# sacco-schema

[![Go](https://github.com/kachaje/sacco-schema/actions/workflows/main.yml/badge.svg)](https://github.com/kachaje/sacco-schema/actions/workflows/main.yml)

A comprehensive SACCO (Savings and Credit Cooperative Organization) management system built in Go, providing complete solutions for managing members, loans, contributions, savings, and financial transactions through USSD and Web interfaces.

## Features

- **Member Management**: Complete member registration and profile management
- **Loan Processing**: Loan applications, approvals, disbursements, and repayments
- **Contributions**: Deposit and withdrawal workflows with schedule management
- **Accounting**: Double-entry ledger system for financial transactions
- **Reports**: Comprehensive reports for loans and contributions
- **Scheduled Jobs**: Automated interest calculations and dividend distributions
- **Multi-Interface**: USSD, WebSocket, and HTTP API support
- **Workflow-Driven**: Dynamic form-driven workflows using workflow-parser

## Documentation

- **[Getting Started](./docs/getting-started.md)**: Installation, setup, and quick start guide
- **[Architecture](./docs/architecture.md)**: System architecture, components, and design patterns
- **[Workflow Design](./docs/workflow-design.md)**: Workflow configuration patterns and menu system design

## Quick Start

```bash
# Clone the repository
git clone https://github.com/kachaje/sacco-schema.git
cd sacco-schema

# Build the application
./ctl.sh -b

# Run the server
./svr -o  # Demo mode (bypasses authentication)
```

See [Getting Started Guide](./docs/getting-started.md) for detailed instructions.

## License

Private
