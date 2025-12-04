# Testing Documentation

## Overview

This document outlines the comprehensive testing strategy for the project, covering Test-Driven Development (TDD), Behavior-Driven Development (BDD) with Cucumber, and all testing practices.

## Testing Philosophy

### Core Principles

1. **Test First**: Write tests before implementation (TDD)
2. **Living Documentation**: Tests serve as executable specifications
3. **Fast Feedback**: Quick test execution for rapid iteration
4. **Isolation**: Tests should be independent and repeatable
5. **Meaningful**: Each test validates specific behavior

### Testing Pyramid

```
        /\
       /  \
      / E2E\     10% - Full system tests (Cucumber)
     /──────\
    /  INT   \   20% - Integration tests
   /──────────\
  /    UNIT    \ 70% - Unit tests (TDD)
 /──────────────\
```

## Test-Driven Development (TDD)

### The TDD Cycle (Red-Green-Refactor)

```
┌─────────────┐
│  Write Test │ ──→ RED (Test fails)
└─────────────┘
       │
       ↓
┌─────────────┐
│  Write Code │ ──→ GREEN (Test passes)
└─────────────┘
       │
       ↓
┌─────────────┐
│  Refactor   │ ──→ CLEAN (Improve code)
└─────────────┘
       │
       ↓ (repeat)
```

### TDD Workflow

#### 1. RED: Write a Failing Test

```go
// internal/domain/services/user_service_test.go
func TestUserService_CreateUser_Success(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)
    service := services.NewUserService(mockRepo, mockLogger)

    email := "test@example.com"
    name := "Test User"

    expectedUser := &domain.User{
        ID:    "user-123",
        Email: email,
        Name:  name,
    }

    mockRepo.EXPECT().
        Create(mock.Anything, email, name).
        Return(expectedUser, nil)

    // Act
    user, err := service.CreateUser(context.Background(), email, name)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, email, user.Email)
    assert.Equal(t, name, user.Name)
}
```

Run test: `just test-unit` → Test fails (RED)

#### 2. GREEN: Write Minimal Code to Pass

```go
// internal/domain/services/user_service.go
type UserService struct {
    repo   domain.UserRepository
    logger Logger
}

func NewUserService(repo domain.UserRepository, logger Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    return s.repo.Create(ctx, email, name)
}
```

Run test: `just test-unit` → Test passes (GREEN)

#### 3. REFACTOR: Improve Code Quality

```go
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    // Add validation
    if err := s.validateEmail(email); err != nil {
        return nil, err
    }

    // Check uniqueness
    exists, err := s.repo.EmailExists(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("checking email existence: %w", err)
    }
    if exists {
        return nil, domain.ErrDuplicateEmail
    }

    // Create user
    user, err := s.repo.Create(ctx, email, name)
    if err != nil {
        return nil, fmt.Errorf("creating user: %w", err)
    }

    s.logger.Info("User created", "userID", user.ID, "email", email)
    return user, nil
}
```

Run test: `just test-unit` → Still passes (REFACTOR complete)

### TDD Best Practices

1. **One Test at a Time**: Focus on single behavior
2. **Small Steps**: Make incremental progress
3. **Test Names**: Should describe behavior, not implementation
4. **AAA Pattern**: Arrange, Act, Assert
5. **Fast Tests**: Unit tests should run in milliseconds

### Test Naming Convention

```go
// Pattern: Test{Component}_{Method}_{Scenario}
func TestUserService_CreateUser_Success(t *testing.T)
func TestUserService_CreateUser_DuplicateEmail(t *testing.T)
func TestUserService_CreateUser_InvalidEmail(t *testing.T)
```

## Unit Testing

### Structure

```
internal/
├── domain/
│   ├── services/
│   │   ├── user_service.go
│   │   └── user_service_test.go      # Unit tests
│   ├── entities/
│   │   ├── user.go
│   │   └── user_test.go              # Entity tests
```

### Example Unit Tests

#### Testing Domain Logic

```go
// internal/domain/entities/user_test.go
func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
        errType error
    }{
        {
            name: "valid user",
            user: &User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: &User{
                Email: "invalid-email",
                Name:  "Test User",
            },
            wantErr: true,
            errType: ErrInvalidEmail,
        },
        {
            name: "empty name",
            user: &User{
                Email: "test@example.com",
                Name:  "",
            },
            wantErr: true,
            errType: ErrEmptyName,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if tt.wantErr {
                assert.Error(t, err)
                assert.ErrorIs(t, err, tt.errType)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### Testing with Mocks

```go
// internal/application/usecases/create_user_test.go
func TestCreateUserUseCase_Execute(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewUserRepository(t)
    mockEventBus := mocks.NewEventBus(t)
    useCase := usecases.NewCreateUserUseCase(mockRepo, mockEventBus)

    input := &usecases.CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    }

    expectedUser := &domain.User{
        ID:    "user-123",
        Email: input.Email,
        Name:  input.Name,
    }

    mockRepo.EXPECT().
        Create(mock.Anything, input.Email, input.Name).
        Return(expectedUser, nil)

    mockEventBus.EXPECT().
        Publish(mock.Anything, mock.MatchedBy(func(evt domain.UserCreated) bool {
            return evt.UserID == expectedUser.ID
        })).
        Return(nil)

    // Act
    output, err := useCase.Execute(context.Background(), input)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser.ID, output.UserID)
    mockRepo.AssertExpectations(t)
    mockEventBus.AssertExpectations(t)
}
```

### Running Unit Tests

```bash
# Run all unit tests
just test-unit

# Run with coverage
go test -short -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run specific package
go test ./internal/domain/services/...

# Run specific test
go test -run TestUserService_CreateUser ./internal/domain/services/
```

## Integration Testing

### Purpose
Test interactions between components (service + database, service + external API)

### Structure

```
test/
├── integration/
│   ├── user_repository_test.go
│   ├── user_service_integration_test.go
│   └── testcontainers/
│       └── postgres.go
```

### Example Integration Test

```go
// test/integration/user_repository_test.go
// +build integration

func TestUserRepository_Create_Integration(t *testing.T) {
    // Arrange - Start test database
    ctx := context.Background()
    container := testcontainers.StartPostgres(t)
    defer container.Terminate(ctx)

    db := container.DB()
    repo := persistence.NewUserRepository(db)

    // Act
    user, err := repo.Create(ctx, "test@example.com", "Test User")

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)

    // Verify in database
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "test@example.com").Scan(&count)
    assert.NoError(t, err)
    assert.Equal(t, 1, count)
}
```

### Running Integration Tests

```bash
# Run integration tests only
just test-integration

# With verbose output
go test -v -tags=integration ./test/integration/...
```

## Behavior-Driven Development (BDD) with Cucumber

### Purpose
- Executable specifications
- Living documentation
- Collaboration between business and technical teams
- End-to-end testing

### Structure

```
features/
├── user_management.feature
├── authentication.feature
├── step_definitions/
│   ├── user_steps.go
│   ├── auth_steps.go
│   └── common_steps.go
└── support/
    ├── hooks.go
    └── world.go
```

### Feature File Example

```gherkin
# features/user_management.feature
Feature: User Management
  As a system administrator
  I want to manage users
  So that I can control system access

  Background:
    Given the system is initialized
    And I am authenticated as an admin

  Scenario: Create a new user
    When I create a user with the following details:
      | email           | name      |
      | john@example.com| John Doe  |
    Then the user should be created successfully
    And the user should have status "Active"
    And I should receive a user ID

  Scenario: Prevent duplicate email
    Given a user exists with email "existing@example.com"
    When I create a user with email "existing@example.com"
    Then the creation should fail
    And I should receive error "email already exists"

  Scenario: Invalid email format
    When I create a user with email "invalid-email"
    Then the creation should fail
    And I should receive error "invalid email format"

  Scenario Outline: Create multiple users
    When I create a user with email "<email>" and name "<name>"
    Then the user should be created successfully

    Examples:
      | email              | name        |
      | alice@example.com  | Alice Smith |
      | bob@example.com    | Bob Jones   |
      | carol@example.com  | Carol White |
```

### Step Definitions

```go
// features/step_definitions/user_steps.go
package step_definitions

import (
    "context"
    "github.com/cucumber/godog"
)

type UserSteps struct {
    world *World
}

func NewUserSteps(world *World) *UserSteps {
    return &UserSteps{world: world}
}

func (s *UserSteps) Register(ctx *godog.ScenarioContext) {
    ctx.Step(`^I create a user with the following details:$`, s.iCreateUserWithDetails)
    ctx.Step(`^I create a user with email "([^"]*)" and name "([^"]*)"$`, s.iCreateUserWithEmailAndName)
    ctx.Step(`^I create a user with email "([^"]*)"$`, s.iCreateUserWithEmail)
    ctx.Step(`^the user should be created successfully$`, s.userShouldBeCreatedSuccessfully)
    ctx.Step(`^the user should have status "([^"]*)"$`, s.userShouldHaveStatus)
    ctx.Step(`^I should receive a user ID$`, s.iShouldReceiveUserID)
    ctx.Step(`^a user exists with email "([^"]*)"$`, s.userExistsWithEmail)
    ctx.Step(`^the creation should fail$`, s.creationShouldFail)
    ctx.Step(`^I should receive error "([^"]*)"$`, s.iShouldReceiveError)
}

func (s *UserSteps) iCreateUserWithEmailAndName(email, name string) error {
    req := &pb.CreateUserRequest{
        Email: email,
        Name:  name,
    }

    resp, err := s.world.grpcClient.CreateUser(context.Background(), req)
    s.world.lastResponse = resp
    s.world.lastError = err

    return nil
}

func (s *UserSteps) userShouldBeCreatedSuccessfully() error {
    if s.world.lastError != nil {
        return fmt.Errorf("expected success but got error: %v", s.world.lastError)
    }
    if s.world.lastResponse == nil {
        return fmt.Errorf("expected response but got nil")
    }
    return nil
}

func (s *UserSteps) iShouldReceiveUserID() error {
    resp, ok := s.world.lastResponse.(*pb.CreateUserResponse)
    if !ok {
        return fmt.Errorf("unexpected response type")
    }
    if resp.UserId == "" {
        return fmt.Errorf("expected user ID but got empty string")
    }
    return nil
}
```

### World Context (Shared State)

```go
// features/support/world.go
package support

type World struct {
    container      *ioc.Container
    grpcClient     pb.UserServiceClient
    grpcConn       *grpc.ClientConn
    lastResponse   interface{}
    lastError      error
    authToken      string
}

func NewWorld() *World {
    return &World{}
}

func (w *World) Reset() {
    w.lastResponse = nil
    w.lastError = nil
}
```

### Hooks (Setup/Teardown)

```go
// features/support/hooks.go
package support

import (
    "github.com/cucumber/godog"
)

func InitializeHooks(ctx *godog.ScenarioContext, world *World) {
    ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
        // Setup: Initialize test environment
        container, err := ioc.NewTestContainer()
        if err != nil {
            return ctx, err
        }
        world.container = container

        // Start gRPC test server
        grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
        if err != nil {
            return ctx, err
        }
        world.grpcConn = grpcConn
        world.grpcClient = pb.NewUserServiceClient(grpcConn)

        return ctx, nil
    })

    ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
        // Teardown: Clean up
        if world.grpcConn != nil {
            world.grpcConn.Close()
        }
        if world.container != nil {
            world.container.Close()
        }
        world.Reset()
        return ctx, nil
    })
}
```

### Test Runner

```go
// features/features_test.go
package features

import (
    "testing"
    "github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
    suite := godog.TestSuite{
        ScenarioInitializer: InitializeScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features"},
            TestingT: t,
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    world := support.NewWorld()

    // Register hooks
    support.InitializeHooks(ctx, world)

    // Register step definitions
    step_definitions.NewUserSteps(world).Register(ctx)
    step_definitions.NewAuthSteps(world).Register(ctx)
    step_definitions.NewCommonSteps(world).Register(ctx)
}
```

### Running BDD Tests

```bash
# Run all cucumber tests
just test-bdd

# Run specific feature
go test -tags=bdd ./features/ -godog.feature=user_management.feature

# Run with specific format
go test -tags=bdd ./features/ -godog.format=pretty
go test -tags=bdd ./features/ -godog.format=cucumber:report.json
```

## Test Organization

### Directory Structure

```
.
├── internal/
│   └── [package]/
│       ├── service.go
│       └── service_test.go          # Unit tests alongside code
├── test/
│   ├── integration/
│   │   └── *_test.go                # Integration tests
│   ├── fixtures/                     # Test data
│   │   ├── users.json
│   │   └── test_data.sql
│   └── mocks/                        # Generated mocks
│       └── *.go
└── features/
    ├── *.feature                     # BDD scenarios
    ├── step_definitions/
    └── support/
```

## Mocking Strategy

### Generate Mocks

```go
//go:generate mockery --name=UserRepository --output=../../test/mocks
type UserRepository interface {
    Create(ctx context.Context, email, name string) (*User, error)
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

Generate: `go generate ./...`

### Using Mocks

```go
import "yourproject/test/mocks"

mockRepo := mocks.NewUserRepository(t)
mockRepo.EXPECT().
    FindByEmail(mock.Anything, "test@example.com").
    Return(nil, domain.ErrNotFound)
```

## Test Coverage

### Measuring Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage by function
go tool cover -func=coverage.out

# View HTML report
go tool cover -html=coverage.out -o coverage.html

# Check coverage threshold
go test -cover ./... | grep -E 'coverage: [0-9]+\.[0-9]+%' | awk '{if ($2 < 80.0) exit 1}'
```

### Coverage Targets
- **Unit Tests**: 80%+ coverage
- **Integration Tests**: Critical paths
- **BDD Tests**: All user-facing features

## Continuous Integration

### Test Pipeline (GitHub Actions Example)

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: just test-unit

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: just test-integration

  bdd-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: just test-bdd

  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v3
```

## Testing Best Practices

### General

1. **Independent Tests**: No test dependencies
2. **Deterministic**: Same input = same output
3. **Fast Execution**: Unit tests < 100ms each
4. **Clear Naming**: Test name describes behavior
5. **Single Assertion Focus**: Test one thing at a time

### TDD Specific

1. **Red First**: Always see test fail before implementing
2. **Small Steps**: Incremental progress
3. **Refactor Fearlessly**: Green tests enable safe refactoring
4. **Test Behavior, Not Implementation**: Test what, not how

### BDD Specific

1. **Business Language**: Use domain terms
2. **Declarative Scenarios**: Focus on what, not how
3. **Living Documentation**: Keep features up-to-date
4. **Collaboration**: Product + Dev + QA review features

## Tools and Libraries

### Testing Frameworks
- **Go testing**: Built-in testing package
- **testify**: Assertions and mocking (`github.com/stretchr/testify`)
- **godog**: Cucumber for Go (`github.com/cucumber/godog`)

### Mocking
- **mockery**: Mock generation (`github.com/vektra/mockery`)
- **testify/mock**: Manual mocking

### Test Containers
- **testcontainers-go**: Docker containers for integration tests

### Coverage
- **go tool cover**: Built-in coverage tool
- **codecov**: Coverage reporting service

## Troubleshooting

### Common Issues

**Tests are slow**
- Solution: Use `-short` flag, parallelize tests with `t.Parallel()`

**Flaky tests**
- Solution: Remove time dependencies, use deterministic test data

**Mock setup is complex**
- Solution: Create test helpers, use table-driven tests

**Cucumber steps are brittle**
- Solution: Use higher-level steps, avoid implementation details

## References

- [Test-Driven Development by Kent Beck](https://www.amazon.com/Test-Driven-Development-Kent-Beck/dp/0321146530)
- [The Cucumber Book](https://pragprog.com/titles/hwcuc2/the-cucumber-book-second-edition/)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Godog Documentation](https://github.com/cucumber/godog)
