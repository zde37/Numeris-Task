# Numeris-Task

Numeris-Task is a robust invoice management system built with Go, leveraging the power of the Gin web framework and PostgreSQL database. This project implements a comprehensive solution for creating, managing, and tracking invoices, users, and payment methods.

## Features

- User management
- Invoice creation and management
- Payment method handling
- Invoice activity tracking
- Detailed invoice retrieval
- Recent invoice and activity fetching

## Project Structure

The project follows a clean architecture pattern to ensure scalability and maintainability:

``` bash
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   ├── config.go
│   │   ├── config_test.go
│   │   ├── database.go
│   │   └── database_test.go
│   ├── controller
│   │   ├── handler.go
│   │   ├── handler_impl.go 
│   │   ├── handler_test.go
│   │   └── stress_test.go
│   ├── helpers
│   │   ├── helpers.go
│   │   └── helpers_test.go
│   ├── mock
│   │   ├── invoice_repo.go
│   │   ├── invoice_service.go
│   │   ├── user_repo.go
│   │   └── user_service.go
│   ├── models
│   │   ├── model.go
│   │   └── request.go
│   ├── repository
│   │   ├── invoice.go
│   │   ├── repo.go
│   │   ├── repo_test.go
│   │   └── user.go
│   └── service
│       ├── invoice.go
│       ├── invoice_test.go
│       ├── service.go
│       ├── user.go
│       └── user_test.go
└── migrations
    ├── 000001_init_schema.down.sql
    └── 000001_init_schema.up.sql
```

- `cmd/`: Contains the main application entry point.
- `internal/`: Houses the core application code.
  - `config/`: Configuration management.
  - `controller/`: HTTP request handlers.
  - `helpers/`: Helper functions.
  - `mocks/`: Contains mocked interfaces for testing.
  - `models/`: Data structures and domain models.
  - `repository/`: Database interaction layer.
  - `service/`: Business logic implementation.
- `migrations/`: Database migration files. 

## Clean Architecture

Numeris-Task adheres to clean architecture principles, separating concerns into distinct layers:

1. Presentation Layer (Controllers)
2. Business Logic Layer (Services)
3. Data Access Layer (Repositories)

This separation ensures that the codebase is modular, testable, and easy to maintain.

## Testing

The project includes both unit tests and stress tests to ensure reliability and performance.

### Running Unit Tests

There is no manual configuration needed for running tests. Ensure you have docker installed and running before running tests. This is important because some tests require a running PostgreSQL database which is provided by a docker test container. The containers are setup automatically.

To run unit tests, use the following command:
``` 
make test
```

If you get any error from the test containers, run this command:
``` bash
export TESTCONTAINERS_RYUK_DISABLED=true
```

### Running Stress Tests

Ensure the server is running first before running stress tests. Start the server with the following command:
```
make run
```

Stress tests can be executed using the following command:
```
make stress
```

## Documentation

The project is well-documented with clear comments explaining the purpose and functionality of each component.  

## Getting Started

1. Clone the repository

2. Set up the PostgreSQL database and update the configuration in `cmd/main.go`
```go 
config.Load(os.Getenv("ENVIRONMENT"), os.Getenv("HTTP_SERVER_ADDRESS"), os.Getenv("DSN"))
```

3. Start the server:
```
make run
```

## Build
Build and run the application:
```
make build-run
```