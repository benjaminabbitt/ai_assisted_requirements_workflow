# Architecture

## Overview

This project follows a clean architecture approach using Go, with dependency injection through manual IoC patterns, comprehensive testing via TDD and BDD, and a dual-protocol API layer (gRPC with REST conversion).

## Technology Stack

- **Language**: Go 1.21+
- **API Protocol**: gRPC (primary) with REST conversion via sidecar
- **Build Tool**: Just (command runner)
- **Testing**:
  - Unit/Integration: Go's native testing + testify
  - BDD: Cucumber (godog)
- **Dependency Injection**: Manual IoC (no framework dependencies)

## Architecture Principles

### 1. Dependency Inversion
- High-level modules do not depend on low-level modules
- Both depend on abstractions (interfaces)
- Manual dependency injection without frameworks like fx

### 2. Clean Architecture Layers
```
┌─────────────────────────────────────┐
│         API Layer (gRPC)            │
│   ┌─────────────────────────────┐   │
│   │   REST Sidecar (grpc-gw)    │   │
│   └─────────────────────────────┘   │
├─────────────────────────────────────┤
│      Application Services           │
├─────────────────────────────────────┤
│       Domain/Business Logic         │
├─────────────────────────────────────┤
│      Repository Interfaces          │
├─────────────────────────────────────┤
│    Infrastructure (DB, External)    │
└─────────────────────────────────────┘
```

## Project Structure

```
.
├── cmd/
│   ├── server/              # Main gRPC server
│   └── gateway/             # REST gateway sidecar
├── internal/
│   ├── domain/              # Domain models and business logic
│   │   ├── entities/
│   │   ├── repositories/    # Repository interfaces
│   │   └── services/        # Domain services
│   ├── application/         # Application services/use cases
│   │   └── usecases/
│   ├── infrastructure/      # External dependencies
│   │   ├── persistence/     # Database implementations
│   │   ├── grpc/            # gRPC server implementations
│   │   └── http/            # HTTP handlers (if needed)
│   └── ioc/                 # Dependency injection container
│       ├── container.go     # Manual DI container
│       └── wire.go          # Wiring logic
├── pkg/                     # Public packages
├── api/
│   └── proto/               # Protocol buffer definitions
│       └── v1/
├── features/                # Cucumber feature files
│   └── step_definitions/    # Step definitions
├── test/
│   ├── integration/         # Integration tests
│   └── fixtures/            # Test fixtures
├── justfile                 # Build commands
└── go.mod
```

## Inversion of Control (IoC) Pattern

### Manual Dependency Injection

Instead of using frameworks like fx, we use a constructor-based pattern with production factory functions.

### Pattern: Primary Constructor + Production Factory

Each component follows this pattern:

1. **Primary constructor** - Takes ALL dependencies as parameters (enables testing)
2. **Production factory** - Builds non-shared dependencies internally, takes shared dependencies as arguments
3. **Coverage exclusion** - Production factory is excluded from test coverage
4. **No business logic** - Production factory MUST NOT contain any business logic, only dependency wiring

```go
// internal/domain/services/user_service.go

// UserService handles user-related business logic
type UserService struct {
    repo      domain.UserRepository
    logger    Logger
    validator Validator
    emailer   EmailService
}

// NewUserService is the primary constructor - takes ALL dependencies
// Use this in tests with mocks
func NewUserService(
    repo domain.UserRepository,
    logger Logger,
    validator Validator,
    emailer EmailService,
) *UserService {
    return &UserService{
        repo:      repo,
        logger:    logger,
        validator: validator,
        emailer:   emailer,
    }
}

// NewUserServiceForProduction builds non-shared dependencies internally
// Takes only shared dependencies (logger, db, etc.)
// Excluded from test coverage
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
    // Build non-shared dependencies internal to this service
    validator := validation.NewUserValidator()
    emailer := email.NewSMTPEmailService(logger)
    repo := persistence.NewUserRepository(db)

    // Call primary constructor
    return NewUserService(repo, logger, validator, emailer)
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    if err := s.validator.ValidateEmail(email); err != nil {
        return nil, err
    }

    user, err := s.repo.Create(ctx, email, name)
    if err != nil {
        return nil, fmt.Errorf("creating user: %w", err)
    }

    s.logger.Info("User created", "userID", user.ID)
    return user, nil
}
```

### Container for Shared Dependencies

The container manages only shared dependencies (database, logger, config):

```go
// internal/ioc/container.go
type Container struct {
    // Shared infrastructure only
    db     *sql.DB
    logger Logger
    cfg    Config

    // Services (built via production factories)
    userService *services.UserService
}

// NewContainer builds the dependency graph
// Excluded from test coverage
func NewContainer(cfg Config) (*Container, error) {
    c := &Container{cfg: cfg}

    // Initialize shared infrastructure
    var err error
    c.db, err = initDatabase(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("initializing database: %w", err)
    }

    c.logger = initLogger(cfg.LogLevel)

    // Wire services using production factories
    c.userService = services.NewUserServiceForProduction(c.db, c.logger)

    return c, nil
}

func (c *Container) UserService() *services.UserService {
    return c.userService
}

func (c *Container) Close() error {
    if c.db != nil {
        return c.db.Close()
    }
    return nil
}
```

### Testing with Primary Constructors

In tests, use the primary constructor with mocks:

```go
// internal/domain/services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange - inject mocks via primary constructor
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)
    mockValidator := mocks.NewValidator(t)
    mockEmailer := mocks.NewEmailService(t)

    service := services.NewUserService(mockRepo, mockLogger, mockValidator, mockEmailer)

    mockValidator.EXPECT().ValidateEmail("test@example.com").Return(nil)
    mockRepo.EXPECT().Create(mock.Anything, "test@example.com", "Test User").Return(&domain.User{
        ID: "user-123",
        Email: "test@example.com",
    }, nil)

    // Act
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### Benefits

- **Testability:** Primary constructor accepts all dependencies, enabling full mock injection
- **No framework dependency:** Pure Go, no magic
- **Coverage exclusion:** Production factories are infrastructure glue, excluded from coverage metrics
- **No business logic in factories:** Production factories only wire dependencies, all business logic stays in testable components
- **Clear separation:** Shared vs. non-shared dependencies are explicit
- **Simple:** Easy to understand and debug, no reflection or code generation

## Testing Strategy

### Test-Driven Development (TDD)

1. **Write test first** - Define expected behavior
2. **Implement minimal code** - Make test pass
3. **Refactor** - Clean up while keeping tests green

```go
// Example: internal/domain/services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewUserRepository(t)
    service := NewUserService(mockRepo, logger)

    // Act
    user, err := service.CreateUser(ctx, "test@example.com")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Behavior-Driven Development (Cucumber/Godog)

```gherkin
# features/user_creation.feature
Feature: User Creation
  As a system administrator
  I want to create new users
  So that they can access the system

  Scenario: Create a valid user
    Given I am authenticated as an admin
    When I create a user with email "test@example.com"
    Then the user should be created successfully
    And I should receive a user ID
```

```go
// features/step_definitions/user_steps.go
func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Step(`^I create a user with email "([^"]*)"$`, createUser)
    ctx.Step(`^the user should be created successfully$`, userCreatedSuccessfully)
}
```

## gRPC + REST Sidecar Pattern

### Architecture

```
┌──────────┐         ┌──────────────┐         ┌──────────┐
│  Client  │ ──REST─→│   Gateway    │ ─gRPC──→│  Server  │
│ (HTTP)   │         │   Sidecar    │         │ (gRPC)   │
└──────────┘         └──────────────┘         └──────────┘
                     grpc-gateway
```

### Protocol Buffers

```proto
// api/proto/v1/user.proto
syntax = "proto3";

package api.v1;

import "google/api/annotations.proto";

service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
    option (google.api.http) = {
      post: "/v1/users"
      body: "*"
    };
  }
}
```

### Deployment

Both containers run side-by-side:
- **gRPC Server**: Listens on `:50051`
- **Gateway Sidecar**: Listens on `:8080`, proxies to `:50051`

Benefits:
- gRPC for internal/service-to-service communication
- REST for external/client-facing API
- Single source of truth (protobuf)
- Automatic OpenAPI/Swagger generation

## Build System (Just)

### Justfile Commands

```justfile
# justfile

# Run all tests
test:
    go test ./...

# Run unit tests only
test-unit:
    go test -short ./...

# Run integration tests
test-integration:
    go test -tags=integration ./test/integration/...

# Run cucumber/BDD tests
test-bdd:
    go test -tags=bdd ./features/...

# Generate protobuf code
proto-gen:
    protoc -I api/proto/v1 \
        --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
        api/proto/v1/*.proto

# Build gRPC server
build-server:
    go build -o bin/server cmd/server/main.go

# Build gateway sidecar
build-gateway:
    go build -o bin/gateway cmd/gateway/main.go

# Build all
build: proto-gen build-server build-gateway

# Run gRPC server
run-server:
    go run cmd/server/main.go

# Run gateway sidecar
run-gateway:
    go run cmd/gateway/main.go

# Run both (in background)
run:
    just run-server &
    just run-gateway

# Clean build artifacts
clean:
    rm -rf bin/
    go clean

# Install dependencies
deps:
    go mod download
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Run linter
lint:
    golangci-lint run ./...

# Format code
fmt:
    go fmt ./...
```

## Development Workflow

1. **Feature Development** (TDD + BDD):
   ```bash
   just test-unit        # Ensure existing tests pass
   # Write Cucumber feature
   # Write step definitions (failing)
   # Write unit tests (failing)
   # Implement feature
   just test             # All tests pass
   ```

2. **Build**:
   ```bash
   just proto-gen        # Generate protobuf code
   just build            # Compile binaries
   ```

3. **Run**:
   ```bash
   just run              # Start both server and gateway
   ```

4. **Quality Checks**:
   ```bash
   just lint             # Check code quality
   just fmt              # Format code
   ```

## Deployment Considerations

### Container Setup

```yaml
# docker-compose.yml
version: '3.8'
services:
  grpc-server:
    build:
      context: .
      dockerfile: Dockerfile.server
    ports:
      - "50051:50051"

  rest-gateway:
    build:
      context: .
      dockerfile: Dockerfile.gateway
    ports:
      - "8080:8080"
    depends_on:
      - grpc-server
    environment:
      GRPC_SERVER_ADDR: "grpc-server:50051"
```

### Configuration

Use environment variables or config files:
- `DATABASE_URL`
- `GRPC_PORT` (default: 50051)
- `HTTP_PORT` (default: 8080)
- `LOG_LEVEL` (debug, info, warn, error)

## Testing Pyramid

```
        /\
       /  \      E2E (Cucumber)
      /────\
     /      \    Integration Tests
    /────────\
   /          \  Unit Tests (TDD)
  /____________\
```

- **Unit Tests**: 70% - Test individual components in isolation
- **Integration Tests**: 20% - Test component interactions
- **E2E Tests**: 10% - Test complete user flows with Cucumber

## Key Design Decisions

1. **No IoC Framework**: Simplicity, no magic, full control
2. **Manual DI Container**: Explicit wiring, easy debugging
3. **gRPC Primary**: Performance, type safety, streaming support
4. **REST via Sidecar**: Client compatibility without code duplication
5. **Just over Make**: Simpler syntax, better cross-platform support
6. **Cucumber for BDD**: Living documentation, stakeholder collaboration

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Godog (Cucumber for Go)](https://github.com/cucumber/godog)
- [Just Command Runner](https://github.com/casey/just)
