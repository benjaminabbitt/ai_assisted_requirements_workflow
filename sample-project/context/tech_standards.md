# Technical Standards

---
version: 1.5.0
last-reviewed: 2025-12-04
reviewed-by: Tech Lead
changelog: |
  - Added GORM v2 as database ORM with comprehensive patterns and examples
  - Added guidance on AI-assisted test development with human software engineering refinement
  - Added comprehensive guidance on AI-assisted implementation requiring software engineering review
  - Emphasized security, performance, consistency, and architecture reviews as real engineering work
  - Added critical guidance on avoiding fragmented logic and ensuring reuse (consistency means unified behavior)
  - Added example showing fragmented email validation causing inconsistent user experience
  - Added emphasis on choosing correct data structures (maps vs slices, time complexity)
  - Added examples showing O(n²) to O(n) optimizations with proper data structure choice
---

## Language & Framework

- **Language:** Go 1.21+
- **Test Framework:** Godog (Cucumber for Go) - `github.com/cucumber/godog`
- **Build Tool:** Just (command runner) - `https://github.com/casey/just`
- **Dependency Management:** Go modules (`go.mod`)
- **Dependency Injection:** Manual IoC pattern (no framework - no fx, no wire)
- **Database ORM:** GORM v2 - `gorm.io/gorm`
- **Database Drivers:**
  - PostgreSQL: `gorm.io/driver/postgres`
  - SQLite (for testing): `gorm.io/driver/sqlite`

## Project Structure

### Directory Layout

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
├── features/                # Cucumber/Godog feature files
│   └── step_definitions/    # Step definitions
├── test/
│   ├── integration/         # Integration tests
│   └── fixtures/            # Test fixtures
├── context/                 # Context files for AI
│   ├── business.md
│   ├── architecture.md
│   ├── testing.md
│   └── tech_standards.md    # This file
├── justfile                 # Build commands
└── go.mod
```

### Feature File Location

**Path:** `/features/`

**Naming Convention:** `{domain}_{capability}.feature`

**Examples:**
- `features/auth/login.feature`
- `features/user/user_creation.feature`
- `features/editor/document_edit.feature`

### Step Definitions Location

**Path:** `/features/step_definitions/`

**Naming Convention:** `{domain}_steps.go`

**Examples:**
- `features/step_definitions/auth_steps.go`
- `features/step_definitions/user_steps.go`
- `features/step_definitions/common_steps.go`

## Step Definition Pattern

### Basic Structure

```go
// features/step_definitions/{domain}_steps.go
package step_definitions

import (
    "context"
    "github.com/cucumber/godog"
)

// DomainSteps holds the state and dependencies for domain-specific steps
type DomainSteps struct {
    world *World
}

// NewDomainSteps creates a new instance with the shared World context
func NewDomainSteps(world *World) *DomainSteps {
    return &DomainSteps{world: world}
}

// Register registers all step definitions for this domain
func (s *DomainSteps) Register(ctx *godog.ScenarioContext) {
    ctx.Step(`^step pattern with "([^"]*)" parameter$`, s.stepMethod)
}

// stepMethod implements the step logic
func (s *DomainSteps) stepMethod(param string) error {
    // Implementation
    return nil
}
```

### World Context (Shared State)

```go
// features/support/world.go
package support

import (
    "yourproject/internal/ioc"
    pb "yourproject/api/proto/v1"
    "google.golang.org/grpc"
)

// World holds shared state across steps within a scenario
type World struct {
    container      *ioc.Container
    grpcClient     pb.ServiceClient
    grpcConn       *grpc.ClientConn
    lastResponse   interface{}
    lastError      error
    authToken      string
    testData       map[string]interface{}
}

func NewWorld() *World {
    return &World{
        testData: make(map[string]interface{}),
    }
}

func (w *World) Reset() {
    w.lastResponse = nil
    w.lastError = nil
    w.testData = make(map[string]interface{})
}
```

### Hooks (Setup/Teardown)

```go
// features/support/hooks.go
package support

import (
    "context"
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

        return ctx, nil
    })

    ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
        // Teardown: Clean up
        if world.container != nil {
            world.container.Close()
        }
        world.Reset()
        return ctx, nil
    })
}
```

### Test Suite Runner

```go
// features/features_test.go
package features

import (
    "testing"
    "github.com/cucumber/godog"
    "yourproject/features/support"
    "yourproject/features/step_definitions"
)

func TestFeatures(t *testing.T) {
    suite := godog.TestSuite{
        ScenarioInitializer: InitializeScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features"},
            TestingT: t,
            Tags:     "@pending", // Run pending scenarios (implemented tests)
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
    step_definitions.NewAuthSteps(world).Register(ctx)
    step_definitions.NewUserSteps(world).Register(ctx)
    step_definitions.NewCommonSteps(world).Register(ctx)
}
```

## Dependency Injection Pattern

### Pattern: Primary Constructor + Production Factory

**Key Principles:**
1. **Primary constructor** - Takes ALL dependencies (for testing)
2. **Production factory** - Builds non-shared dependencies internally, takes only shared ones
3. **No business logic** - Production factories MUST NOT contain any business logic, only dependency wiring
4. **Coverage exclusion** - Production factories are excluded from test coverage

### Service Example

```go
// internal/domain/services/user_service.go
package services

import (
    "context"
    "gorm.io/gorm"
    "yourproject/internal/domain"
)

type UserService struct {
    repo      domain.UserRepository
    logger    Logger
    validator Validator
}

// NewUserService is the PRIMARY CONSTRUCTOR
// Takes ALL dependencies - use this in tests
func NewUserService(
    repo domain.UserRepository,
    logger Logger,
    validator Validator,
) *UserService {
    return &UserService{
        repo:      repo,
        logger:    logger,
        validator: validator,
    }
}

// NewUserServiceForProduction is the PRODUCTION FACTORY
// Builds non-shared dependencies internally
// Takes only shared dependencies (db, logger, config)
// EXCLUDED FROM TEST COVERAGE
// coverage:ignore
func NewUserServiceForProduction(db *gorm.DB, logger Logger) *UserService {
    // Build non-shared dependencies inside this function
    validator := validation.NewUserValidator()
    repo := persistence.NewUserRepository(db)

    // Call primary constructor
    return NewUserService(repo, logger, validator)
}
```

### Container for Shared Dependencies Only

The container manages shared infrastructure and calls production factories:

```go
// internal/ioc/container.go
package ioc

import (
    "gorm.io/gorm"
    "yourproject/internal/domain/services"
)

type Container struct {
    // Shared infrastructure only
    db     *gorm.DB
    logger Logger
    cfg    Config

    // Services (built via production factories)
    userService *services.UserService
}

// NewContainer builds the dependency graph
// EXCLUDED FROM TEST COVERAGE
// coverage:ignore
func NewContainer(cfg Config) (*Container, error) {
    c := &Container{cfg: cfg}

    // Initialize shared infrastructure
    var err error
    c.db, err = initDatabase(cfg.DatabaseURL)
    if err != nil {
        return nil, err
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
        sqlDB, err := c.db.DB()
        if err != nil {
            return err
        }
        return sqlDB.Close()
    }
    return nil
}
```

### Testing with Primary Constructors

In tests, always use the primary constructor with mocks:

```go
// internal/domain/services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange - inject ALL dependencies as mocks
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)
    mockValidator := mocks.NewValidator(t)

    // Use primary constructor
    service := services.NewUserService(mockRepo, mockLogger, mockValidator)

    // Set expectations
    mockValidator.EXPECT().ValidateEmail("test@example.com").Return(nil)
    mockRepo.EXPECT().Create(mock.Anything, "test@example.com", "Test").Return(
        &domain.User{ID: "123", Email: "test@example.com"}, nil,
    )

    // Act
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "123", user.ID)
}
```

### Test Container for BDD

For integration/BDD tests that need real infrastructure:

```go
// internal/ioc/test_container.go
package ioc

// NewTestContainer creates container with test configuration
// EXCLUDED FROM TEST COVERAGE
func NewTestContainer() (*Container, error) {
    cfg := Config{
        DatabaseURL: "postgres://test:test@localhost:5432/test_db",
        LogLevel:    "debug",
    }
    return NewContainer(cfg)
}
```

### Coverage Exclusion

Add coverage exclusion comments to production factories:

```go
// NewUserServiceForProduction ...
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
    // ...
}

// NewContainer ...
// coverage:ignore
func NewContainer(cfg Config) (*Container, error) {
    // ...
}
```

Or use build tags to exclude from coverage:

```go
//go:build !test

package ioc

// Container wiring excluded from test coverage
```

### Benefits

- **100% testable:** Primary constructor enables full dependency injection
- **No framework:** Pure Go, no reflection, no magic
- **No business logic in factories:** Production factories only wire dependencies, all business logic stays in testable components
- **Clear separation:** Shared vs. component-specific dependencies
- **Coverage accuracy:** Production wiring excluded, business logic fully covered
- **Simple:** Easy to understand, easy to debug

## Database Access with GORM

### Domain Entity

Define domain entities as plain Go structs:

```go
// internal/domain/entities/user.go
package entities

import (
    "time"
    "gorm.io/gorm"
)

// User represents a user in the system
type User struct {
    ID        string         `gorm:"primaryKey;type:varchar(36)"`
    Email     string         `gorm:"uniqueIndex;not null;type:varchar(255)"`
    Name      string         `gorm:"not null;type:varchar(255)"`
    Status    string         `gorm:"not null;type:varchar(20);default:'active'"`
    CreatedAt time.Time      `gorm:"not null"`
    UpdatedAt time.Time      `gorm:"not null"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
    return "users"
}
```

### Repository Interface (Domain Layer)

Define repository interfaces in the domain layer:

```go
// internal/domain/repositories/user_repository.go
package repositories

import (
    "context"
    "yourproject/internal/domain/entities"
)

// UserRepository defines operations for user persistence
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id string) (*entities.User, error)
    FindByEmail(ctx context.Context, email string) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*entities.User, error)
    Count(ctx context.Context) (int64, error)
}
```

### Repository Implementation (Infrastructure Layer)

Implement repository using GORM in infrastructure layer:

```go
// internal/infrastructure/persistence/user_repository.go
package persistence

import (
    "context"
    "errors"
    "yourproject/internal/domain/entities"
    "yourproject/internal/domain/repositories"
    "gorm.io/gorm"
)

// UserRepositoryImpl implements UserRepository using GORM
type UserRepositoryImpl struct {
    db *gorm.DB
}

// NewUserRepository creates a new UserRepository implementation
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &UserRepositoryImpl{db: db}
}

// Create inserts a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
    result := r.db.WithContext(ctx).Create(user)
    return result.Error
}

// FindByID retrieves a user by ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.User, error) {
    var user entities.User
    result := r.db.WithContext(ctx).First(&user, "id = ?", id)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, repositories.ErrUserNotFound
    }

    return &user, result.Error
}

// FindByEmail retrieves a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
    var user entities.User
    result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, repositories.ErrUserNotFound
    }

    return &user, result.Error
}

// Update modifies an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
    result := r.db.WithContext(ctx).Save(user)
    return result.Error
}

// Delete soft-deletes a user
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
    result := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id)
    if result.RowsAffected == 0 {
        return repositories.ErrUserNotFound
    }
    return result.Error
}

// List retrieves paginated users
func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
    var users []*entities.User
    result := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&users)

    return users, result.Error
}

// Count returns total number of users
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
    var count int64
    result := r.db.WithContext(ctx).Model(&entities.User{}).Count(&count)
    return count, result.Error
}
```

### Service Using Repository (Following IoC Pattern)

```go
// internal/domain/services/user_service.go
package services

import (
    "context"
    "yourproject/internal/domain/entities"
    "yourproject/internal/domain/repositories"
)

type UserService struct {
    repo      repositories.UserRepository
    logger    Logger
    validator Validator
}

// NewUserService is the PRIMARY CONSTRUCTOR
func NewUserService(
    repo repositories.UserRepository,
    logger Logger,
    validator Validator,
) *UserService {
    return &UserService{
        repo:      repo,
        logger:    logger,
        validator: validator,
    }
}

// NewUserServiceForProduction is the PRODUCTION FACTORY
// coverage:ignore
func NewUserServiceForProduction(db *gorm.DB, logger Logger) *UserService {
    repo := persistence.NewUserRepository(db)
    validator := validation.NewUserValidator()
    return NewUserService(repo, logger, validator)
}

// CreateUser validates and creates a new user
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*entities.User, error) {
    // Validate input
    if err := s.validator.ValidateEmail(email); err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }

    // Check for duplicates
    existing, err := s.repo.FindByEmail(ctx, email)
    if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
        return nil, fmt.Errorf("checking duplicate: %w", err)
    }
    if existing != nil {
        return nil, repositories.ErrUserAlreadyExists
    }

    // Create user
    user := &entities.User{
        ID:     uuid.New().String(),
        Email:  email,
        Name:   name,
        Status: "active",
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("creating user: %w", err)
    }

    s.logger.Info("User created", "id", user.ID, "email", user.Email)
    return user, nil
}
```

### Testing with Mock Repository

```go
// internal/domain/services/user_service_test.go
package services_test

import (
    "context"
    "testing"
    "yourproject/internal/domain/entities"
    "yourproject/internal/domain/services"
    "yourproject/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser_Success(t *testing.T) {
    // Arrange - use mocks, not real GORM
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)
    mockValidator := mocks.NewValidator(t)

    // Use primary constructor
    service := services.NewUserService(mockRepo, mockLogger, mockValidator)

    email := "test@example.com"
    name := "Test User"

    // Set expectations
    mockValidator.EXPECT().ValidateEmail(email).Return(nil)
    mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(nil, repositories.ErrUserNotFound)
    mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
        return u.Email == email && u.Name == name
    })).Return(nil)
    mockLogger.EXPECT().Info("User created", "id", mock.Anything, "email", email)

    // Act
    user, err := service.CreateUser(context.Background(), email, name)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, email, user.Email)
    assert.Equal(t, name, user.Name)
    assert.NotEmpty(t, user.ID)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)
    mockValidator := mocks.NewValidator(t)

    service := services.NewUserService(mockRepo, mockLogger, mockValidator)

    email := "existing@example.com"
    existingUser := &entities.User{ID: "123", Email: email}

    mockValidator.EXPECT().ValidateEmail(email).Return(nil)
    mockRepo.EXPECT().FindByEmail(mock.Anything, email).Return(existingUser, nil)

    // Act
    user, err := service.CreateUser(context.Background(), email, "New User")

    // Assert
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.ErrorIs(t, err, repositories.ErrUserAlreadyExists)
}
```

### Database Initialization in Container

```go
// internal/ioc/container.go
package ioc

import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "yourproject/internal/domain/entities"
)

// initDatabase initializes GORM connection
// coverage:ignore
func initDatabase(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        // Configure GORM
        PrepareStmt: true,
    })
    if err != nil {
        return nil, fmt.Errorf("opening database: %w", err)
    }

    // Run migrations (development only)
    if err := db.AutoMigrate(
        &entities.User{},
        // Add other entities here
    ); err != nil {
        return nil, fmt.Errorf("running migrations: %w", err)
    }

    return db, nil
}
```

### Test Database Setup (SQLite)

```go
// internal/ioc/test_container.go
package ioc

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "yourproject/internal/domain/entities"
)

// NewTestContainer creates container with in-memory SQLite
// coverage:ignore
func NewTestContainer() (*Container, error) {
    // Use in-memory SQLite for tests
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Run migrations
    if err := db.AutoMigrate(&entities.User{}); err != nil {
        return nil, err
    }

    logger := &testLogger{}

    c := &Container{
        db:     db,
        logger: logger,
    }

    // Wire services
    c.userService = services.NewUserServiceForProduction(db, logger)

    return c, nil
}
```

### Transaction Support

For operations requiring transactions:

```go
// internal/domain/services/user_service.go

// TransferOwnership transfers ownership from one user to another (requires transaction)
func (s *UserService) TransferOwnership(ctx context.Context, fromUserID, toUserID string) error {
    // Get GORM DB from repository (if needed)
    // For complex transactions, consider adding BeginTx/Commit/Rollback to repository interface

    return s.repo.WithTransaction(ctx, func(txRepo repositories.UserRepository) error {
        // All operations in this function use the transaction

        fromUser, err := txRepo.FindByID(ctx, fromUserID)
        if err != nil {
            return fmt.Errorf("finding from user: %w", err)
        }

        toUser, err := txRepo.FindByID(ctx, toUserID)
        if err != nil {
            return fmt.Errorf("finding to user: %w", err)
        }

        // Transfer logic...
        fromUser.Status = "inactive"
        toUser.Status = "owner"

        if err := txRepo.Update(ctx, fromUser); err != nil {
            return fmt.Errorf("updating from user: %w", err)
        }

        if err := txRepo.Update(ctx, toUser); err != nil {
            return fmt.Errorf("updating to user: %w", err)
        }

        return nil
    })
}
```

### GORM Best Practices

1. **Always use context:** Pass context to all GORM operations via `WithContext(ctx)`
2. **Define custom errors:** Map `gorm.ErrRecordNotFound` to domain-specific errors
3. **Use struct tags:** Define all constraints in entity tags (indexes, types, constraints)
4. **Migrations:** Use `AutoMigrate` in development, dedicated migration tool in production
5. **Connection pooling:** Configure via `gorm.Config` and underlying `database/sql` settings
6. **Soft deletes:** Use `gorm.DeletedAt` for entities that should be soft-deleted
7. **Preloading:** Use `Preload()` for loading associations (be mindful of N+1 queries)
8. **Raw SQL (when needed):** Use `db.Raw()` for complex queries, but prefer GORM query builder when possible

## API Protocols

### gRPC (Primary)

- **Server Port:** `:50051`
- **Protocol Buffers:** `/api/proto/v1/*.proto`
- **Generated Code:** Auto-generated via `protoc`

### REST (Via Gateway Sidecar)

- **Gateway Port:** `:8080`
- **Implementation:** grpc-gateway
- **Deployment:** Sidecar container proxies to gRPC

## Build System (Just)

### Running Tests

```bash
# Run unit tests
just test-unit

# Run integration tests
just test-integration

# Run BDD/Cucumber tests
just test-bdd

# Run all tests
just test
```

### Justfile Commands

The `justfile` at the project root defines all build commands. Key test-related commands:

```justfile
# Run cucumber/BDD tests
test-bdd:
    go test -tags=bdd ./features/...

# Run specific feature
test-feature feature:
    go test -tags=bdd ./features/ -godog.feature={{feature}}

# Run scenarios with specific tag
test-tag tag:
    go test -tags=bdd ./features/ -godog.tags={{tag}}
```

## Code Style & Conventions

### General Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Exported functions/types require godoc comments

### Package Naming

- Package names: lowercase, no underscores (`package user`, not `package userService`)
- Directory names match package names
- Internal packages: `/internal/` (not importable outside project)
- Public packages: `/pkg/` (importable by external projects)

### Interface Naming

- Interfaces: Noun or adjective (e.g., `Reader`, `Closer`, `UserRepository`)
- Single-method interfaces: `-er` suffix (e.g., `Handler`, `Validator`)

### Error Handling

```go
// Define domain errors as package-level variables
var (
    ErrNotFound       = errors.New("not found")
    ErrInvalidInput   = errors.New("invalid input")
    ErrUnauthorized   = errors.New("unauthorized")
)

// Wrap errors with context
if err != nil {
    return fmt.Errorf("creating user: %w", err)
}

// Check error types
if errors.Is(err, domain.ErrNotFound) {
    // Handle not found
}
```

### Testing Conventions

#### Unit Tests

- File naming: `{name}_test.go` alongside source
- Test naming: `Test{Type}_{Method}_{Scenario}`
- Use table-driven tests for multiple cases
- Mock external dependencies

```go
func TestUserService_CreateUser_Success(t *testing.T) {
    // Arrange
    // Act
    // Assert
}
```

#### AI-Assisted Test Development

**AI generates unit tests following TDD patterns, but developers apply software engineering principles:**

**AI's role:**
- Generate test skeleton following TDD cycle
- Apply mocking patterns from context files
- Follow naming conventions
- Generate boundary condition tests
- Create table-driven test structures

**Developer's role (software engineering):**
- **Maintainability:** Refactor test code for clarity, extract test helpers, eliminate duplication
- **Consistency:** Ensure test style matches team conventions, naming coherent across suite
- **Performance:** Optimize slow tests, parallelize where appropriate, reduce fixture overhead
- **Completeness:** Review edge cases, add missing scenarios, validate meaningful coverage
- **Quality:** Ensure assertions are meaningful, error messages clear, setup/teardown proper

**Example refinement workflow:**
```go
// 1. AI generates initial test
func TestUserService_CreateUser_Success(t *testing.T) {
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)
    mockValidator := mocks.NewValidator(t)
    service := services.NewUserService(mockRepo, mockLogger, mockValidator)
    mockValidator.EXPECT().ValidateEmail("test@example.com").Return(nil)
    mockRepo.EXPECT().FindByEmail(mock.Anything, "test@example.com").Return(nil, repositories.ErrUserNotFound)
    mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test")
    assert.NoError(t, err)
    assert.NotNil(t, user)
}

// 2. Developer refactors for maintainability
func TestUserService_CreateUser_Success(t *testing.T) {
    // Arrange
    svc, mocks := setupUserServiceTest(t)
    email := "test@example.com"

    mocks.validator.EXPECT().ValidateEmail(email).Return(nil)
    mocks.repo.EXPECT().FindByEmail(mock.Anything, email).Return(nil, repositories.ErrUserNotFound)
    mocks.repo.EXPECT().Create(mock.Anything, matchUser(email, "Test")).Return(nil)

    // Act
    user, err := svc.CreateUser(context.Background(), email, "Test")

    // Assert
    require.NoError(t, err)
    assert.Equal(t, email, user.Email)
    assert.NotEmpty(t, user.ID, "user ID should be generated")
}

// Helper function extracted by developer
func setupUserServiceTest(t *testing.T) (*services.UserService, *userServiceMocks) {
    t.Helper()
    mocks := &userServiceMocks{
        repo:      mocks.NewUserRepository(t),
        logger:    mocks.NewLogger(t),
        validator: mocks.NewValidator(t),
    }
    service := services.NewUserService(mocks.repo, mocks.logger, mocks.validator)
    return service, mocks
}
```

**Key principle:** AI generates syntactically correct tests following patterns. Developers ensure tests form a coherent, maintainable suite that serves the team long-term.

#### AI-Assisted Implementation Development

**AI generates implementation code following patterns, but developers apply software engineering to production code:**

**AI's role:**
- Generate code following TDD cycle
- Apply IoC patterns (primary constructor + production factory)
- Follow naming conventions and directory structure
- Implement standard CRUD operations
- Generate boilerplate and repetitive code

**Developer's role (software engineering):**
- **Security:** Review for vulnerabilities
  - Input validation on all external data
  - SQL injection prevention (parameterized queries/ORM)
  - Authentication/authorization checks
  - Sensitive data handling (PII, credentials)
  - OWASP Top 10 considerations
- **Performance:** Optimize implementations
  - Identify N+1 query problems
  - Reduce unnecessary allocations
  - **Choose correct data structures for the operation:**
    - Maps for lookups/membership (O(1)), not slice iteration (O(n))
    - Sets (map[T]struct{}) for unique collections, not slice + dedup
    - Slices for ordered collections, arrays for fixed size
    - Consider time complexity: O(1) vs O(n) vs O(n²)
  - Choose efficient algorithms
  - Add caching where appropriate
  - Profile and optimize hot paths
- **Consistency & Reuse:** Maintain codebase coherence and avoid fragmented logic
  - **Extract and reuse common logic:** Validation, business rules, transformations should exist in one place
  - **Avoid fragmentation:** Same logic must produce same results regardless of code path
  - **Prevent divergent implementations:** Email validation shouldn't differ between API and UI paths
  - **Refactor duplicates:** AI often generates similar-but-different logic across sessions (context window exhaustion)
  - **Ensure naming conventions consistent throughout**
  - **Unify error handling patterns**
  - **Maintain architectural consistency**
- **Architecture:** Ensure proper design
  - Validate proper layer separation (domain/application/infrastructure)
  - Review abstractions for appropriateness
  - Ensure dependencies point in correct direction
  - Validate API design and contracts
- **Maintainability:** Keep code readable and manageable
  - Simplify complex logic
  - Extract functions when needed
  - Add clarifying comments for non-obvious decisions
  - Keep functions small and focused
- **Correctness:** Validate business logic
  - Verify edge cases are handled
  - Ensure error paths are comprehensive
  - Validate against business rules
  - Test with real-world scenarios

**Common issues requiring human intervention:**

| Issue | AI Generates | Developer Fixes |
|-------|--------------|-----------------|
| **Context exhaustion** | Inconsistent naming, duplicate logic across files | Refactor for consistency, extract common code |
| **Fragmented logic** | Same validation implemented differently in multiple places | Extract to single reusable validator, ensure uniform behavior |
| **Security gaps** | Basic validation, no depth | Add input sanitization, check authorization, handle sensitive data properly |
| **Performance issues** | Naive implementations (N+1 queries, wrong data structures, inefficient loops) | Optimize with batching, caching, correct data structures, better algorithms |
| **Over-engineering** | Complex abstractions for simple problems | Simplify to appropriate level |
| **Under-engineering** | Repeated code, no abstractions | Extract common patterns, create reusable components |
| **Poor error handling** | Generic errors, incomplete error paths | Add context to errors, handle all failure modes |

**Example 1: Security and Performance Issues**

```go
// AI-generated: Syntactically correct, follows patterns, but has issues
func (s *UserService) GetUserWithOrders(ctx context.Context, userID string) (*entities.User, error) {
    // ⚠️ Security: No authorization check - anyone can view any user
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    // ⚠️ Performance: N+1 query - fetches orders one at a time
    for i := range user.Orders {
        order, _ := s.orderRepo.FindByID(ctx, user.Orders[i].ID)
        user.Orders[i] = *order
    }

    return user, nil
}

// Developer-refined: Secure, performant, maintainable
func (s *UserService) GetUserWithOrders(ctx context.Context, userID string, requestorID string) (*entities.User, error) {
    // Security: Verify authorization
    if !s.authz.CanViewUser(ctx, requestorID, userID) {
        return nil, domain.ErrUnauthorized
    }

    // Fetch user
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("finding user %s: %w", userID, err)
    }

    // Performance: Batch load orders (single query)
    orderIDs := make([]string, len(user.Orders))
    for i, order := range user.Orders {
        orderIDs[i] = order.ID
    }

    orders, err := s.orderRepo.FindByIDs(ctx, orderIDs)
    if err != nil {
        return nil, fmt.Errorf("loading orders for user %s: %w", userID, err)
    }

    // Map orders back to user
    orderMap := make(map[string]*entities.Order)
    for _, order := range orders {
        orderMap[order.ID] = order
    }

    for i := range user.Orders {
        if order, exists := orderMap[user.Orders[i].ID]; exists {
            user.Orders[i] = *order
        }
    }

    return user, nil
}
```

**Example 2: Wrong Data Structure (Performance Issue)**

```go
// ⚠️ Problem: AI uses slice iteration for membership checks (O(n) when should be O(1))

// AI-generated: Works, but inefficient
func (s *UserService) HasPermissions(ctx context.Context, userID string, requiredPerms []string) (bool, error) {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return false, err
    }

    // ⚠️ Performance: O(n*m) - checking each required perm against user perms with nested loops
    for _, required := range requiredPerms {
        found := false
        for _, userPerm := range user.Permissions {
            if userPerm == required {
                found = true
                break
            }
        }
        if !found {
            return false, nil
        }
    }

    return true, nil
}

// Developer-refined: Use map for O(1) lookups
func (s *UserService) HasPermissions(ctx context.Context, userID string, requiredPerms []string) (bool, error) {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return false, err
    }

    // Build map once: O(n)
    userPermSet := make(map[string]struct{}, len(user.Permissions))
    for _, perm := range user.Permissions {
        userPermSet[perm] = struct{}{}
    }

    // Check each required perm: O(m) with O(1) lookups
    for _, required := range requiredPerms {
        if _, exists := userPermSet[required]; !exists {
            return false, nil
        }
    }

    return true, nil
}

// Performance improvement: O(n*m) → O(n+m)
// For 100 user perms and 10 required perms:
//   Before: Up to 1000 comparisons
//   After: 110 operations (100 to build map + 10 lookups)
```

**Another common pattern: Filtering duplicates**

```go
// ⚠️ AI-generated: O(n²) duplicate check
func RemoveDuplicates(items []string) []string {
    var result []string
    for _, item := range items {
        // O(n) check for each item = O(n²) total
        found := false
        for _, existing := range result {
            if existing == item {
                found = true
                break
            }
        }
        if !found {
            result = append(result, item)
        }
    }
    return result
}

// ✅ Developer-refined: O(n) with set
func RemoveDuplicates(items []string) []string {
    seen := make(map[string]struct{}, len(items))
    result := make([]string, 0, len(items))

    for _, item := range items {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }

    return result
}

// Performance: O(n²) → O(n)
// For 1000 items: ~500,000 ops → ~1000 ops
```

**Example 3: Fragmented Logic (Critical Consistency Issue)**

```go
// ⚠️ Problem: AI generates email validation in multiple places with subtle differences

// In UserService (AI session 1)
func (s *UserService) CreateUser(ctx context.Context, email string) error {
    // Basic regex check
    if !regexp.MustCompile(`^[a-z0-9]+@[a-z0-9]+\.[a-z]+$`).MatchString(email) {
        return errors.New("invalid email")
    }
    // ... create user
}

// In ProfileService (AI session 2, different context)
func (s *ProfileService) UpdateEmail(ctx context.Context, userID, email string) error {
    // Different regex - allows uppercase and more complex domains
    if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
        return errors.New("invalid email format")
    }
    // ... update email
}

// In AuthService (AI session 3, copy-pasted from somewhere)
func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
    // Yet another validation - allows plus addressing
    if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
        return errors.New("email is invalid")
    }
    // ... send reset
}

// ⚠️ Result: User creates account with "John+test@example.com"
//            CreateUser: REJECTS (no + in regex)
//            User tries UpdateEmail with same address: ACCEPTS (allows +)
//            User confused: "Why did it fail before but work now?"
//            ResetPassword: ACCEPTS (minimal check)
//
// Same email, three different validation results - fragmented logic!

// ✅ Developer solution: Extract to single reusable validator

// internal/domain/validators/email_validator.go
type EmailValidator struct {
    allowPlusAddressing bool
    allowUppercase      bool
}

func NewEmailValidator() *EmailValidator {
    return &EmailValidator{
        allowPlusAddressing: true,  // Business decision: allow plus addressing
        allowUppercase:      true,  // Business decision: case insensitive
    }
}

func (v *EmailValidator) Validate(email string) error {
    email = strings.ToLower(strings.TrimSpace(email))

    // Single regex used everywhere - one source of truth
    pattern := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
    if !regexp.MustCompile(pattern).MatchString(email) {
        return domain.ErrInvalidEmail
    }

    // Additional checks (optional, but consistent everywhere)
    if len(email) > 254 {
        return domain.ErrEmailTooLong
    }

    return nil
}

// Now use everywhere - consistent behavior
func (s *UserService) CreateUser(ctx context.Context, email string) error {
    if err := s.emailValidator.Validate(email); err != nil {
        return err
    }
    // ... create user
}

func (s *ProfileService) UpdateEmail(ctx context.Context, userID, email string) error {
    if err := s.emailValidator.Validate(email); err != nil {
        return err
    }
    // ... update email
}

func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
    if err := s.emailValidator.Validate(email); err != nil {
        return err
    }
    // ... send reset
}

// ✅ Result: Same validation everywhere, consistent user experience
//            If validation rules change, update one place
//            Tests validate the single validator, not each service
```

**Key principle:** AI generates code that compiles and passes tests. Developers ensure code is secure, **uses correct data structures**, performant, **has unified logic (not fragmented)**, maintainable, and correctly implements business requirements. **This is software engineering, not brick-laying.**

#### BDD Step Naming

- Use clear, business-readable language
- Parameterize with regex capture groups
- Reuse existing steps when possible

```go
ctx.Step(`^I create a user with email "([^"]*)"$`, s.iCreateUserWithEmail)
ctx.Step(`^the user should be created successfully$`, s.userCreatedSuccessfully)
```

## Dependency Management

### Adding Dependencies

```bash
# Add a new dependency
go get github.com/example/package

# Update dependencies
go get -u ./...

# Tidy dependencies
go mod tidy
```

### Required Dependencies

```go
// go.mod
require (
    // Testing
    github.com/cucumber/godog v0.12.0
    github.com/stretchr/testify v1.8.0

    // gRPC
    google.golang.org/grpc v1.50.0
    google.golang.org/protobuf v1.28.0

    // Database
    gorm.io/gorm v1.25.0
    gorm.io/driver/postgres v1.5.0
    gorm.io/driver/sqlite v1.5.0  // For testing
)
```

## Context File Locations

AI reads these files for context when generating specifications:

| File | Path | Purpose |
|------|------|---------|
| Business Context | `/context/business.md` | Domain knowledge, personas, business rules |
| Architecture Context | `/context/architecture.md` | External dependencies, system architecture |
| Testing Context | `/context/testing.md` | Step library, boundary patterns, edge cases |
| Tech Standards | `/context/tech_standards.md` | This file - language, framework, patterns |

## API Interface Specs

### Internal APIs

- **Format:** Protocol Buffers (`.proto`)
- **Location:** `/api/proto/v1/`
- **Discovery:** AI reads these directly for parameter validation

### External APIs

- **Documented in:** `/context/architecture.md`
- **Local copies (if available):** `/api/specs/external/`

## CI/CD Integration

### GitHub Actions (Recommended)

```yaml
name: BDD Tests

on: [push, pull_request]

jobs:
  bdd-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run BDD tests
        run: just test-bdd
```

## Development Workflow

1. **Feature Development** (TDD + BDD):
   - Write Cucumber scenarios (in spec/ branch)
   - Write unit tests (failing)
   - Implement feature
   - Implement step definitions
   - Scenarios remain `@pending` (business validation tracked in deployment)

2. **Running Tests:**
   ```bash
   just test-unit        # Unit tests
   just test-bdd         # BDD tests
   just test             # All tests
   ```

3. **Code Quality:**
   ```bash
   just lint             # Run linter
   just fmt              # Format code
   ```

## Security & Best Practices

### Input Validation

- Validate all external inputs in domain layer
- Use strong typing (avoid `interface{}` where possible)
- Sanitize inputs before database queries

### Authentication

- Use JWT tokens with proper expiration
- Store secrets in environment variables, not code
- Test authentication in merged scenarios only (avoid credentials in unmerged specs)

### Database Access

- **ORM:** Use GORM for database operations (prevents SQL injection, provides type safety)
- **Repository Pattern:** Implement repositories that wrap GORM for testability
- **Transactions:** Use GORM transactions for multi-step operations
- **Migrations:** Use GORM AutoMigrate for schema management (development) or dedicated migration tool (production)

## Version Control

### Branch Naming

Follow Workflow.md conventions:

- Specification: `spec/{ticket-id}`
- Implementation: `impl/{ticket-id}`
- Amendment: `amend/{ticket-id}`
- Hotfix: `hotfix-spec/{ticket-id}`

### Commit Messages

- Use conventional commits format
- Include ticket ID when applicable
- Format: `type(scope): description`

Examples:
- `feat(auth): add user login endpoint`
- `test(auth): add BDD scenarios for login PROJ-1234`
- `fix(user): prevent duplicate email creation`

## References

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Godog Documentation](https://github.com/cucumber/godog)
- [Just Command Runner](https://github.com/casey/just)
- [gRPC-Go](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
