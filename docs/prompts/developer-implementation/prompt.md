# Developer Implementation Prompt

Implement step definitions for approved Gherkin specifications following TDD and IoC patterns.

## Inputs

1. **Approved spec:** Feature file with `@pending` scenarios (merged to main)
2. **Skeleton steps:** Auto-generated step definition skeletons
3. **Tech standards:** `tech_standards.md` - Go/Godog patterns, IoC rules
4. **Architecture:** `architecture.md` - External dependencies, APIs

## Process

1. Create branch: `impl/{ticket-id}`
2. Review spec scenarios (`@pending` tags)
3. Implement step definitions using **primary constructors + mocks**
4. Write unit tests first (TDD)
5. Implement business logic
6. Run scenarios until passing
7. Change `scenarios remain `@pending`
8. Open implementation PR

## Critical Rules

### IoC Pattern (MUST FOLLOW)

**Every service:**
1. **Primary constructor** - Takes ALL dependencies
2. **Production factory** - Builds non-shared deps, calls primary constructor
3. **No business logic in factory** - Zero conditionals, loops, calculations
4. **Coverage exclusion** - `// coverage:ignore` on production factory

**In tests:**
- ✅ Use primary constructor with mocks
- ❌ NEVER use production factory in tests

### TDD Cycle

```
Write failing test → Implement minimal code → Test passes → Refactor
```

## Implementation Checklist

Per step definition:
- [ ] Step implementation uses services via primary constructor
- [ ] Unit tests use mocks, not real infrastructure
- [ ] No business logic in production factories
- [ ] Production factories marked `// coverage:ignore`
- [ ] Error handling with proper wrapping
- [ ] Scenario runs and passes

## Example: Implementing Password Reset Steps

### Step 1: Review Approved Spec

```gherkin
@pending @story-PROJ-1234
Scenario: User requests password reset successfully
  Given a user exists with email "user@example.com"
  When I request a password reset for "user@example.com"
  Then I should receive a reset email
  And the password reset attempt should be logged
```

### Step 2: Write Failing Unit Test (TDD)

```go
// internal/domain/services/password_reset_service_test.go
func TestPasswordResetService_RequestReset_Success(t *testing.T) {
    // Arrange - use primary constructor with mocks
    mockUserRepo := mocks.NewUserRepository(t)
    mockEmailService := mocks.NewEmailService(t)
    mockTokenGen := mocks.NewTokenGenerator(t)
    mockAuditLog := mocks.NewAuditLogger(t)

    service := services.NewPasswordResetService(
        mockUserRepo,
        mockEmailService,
        mockTokenGen,
        mockAuditLog,
    )

    user := &domain.User{ID: "user-123", Email: "user@example.com"}
    mockUserRepo.EXPECT().
        FindByEmail(mock.Anything, "user@example.com").
        Return(user, nil)

    mockTokenGen.EXPECT().
        Generate().
        Return("reset-token-123", nil)

    mockEmailService.EXPECT().
        SendPasswordResetEmail(mock.Anything, user, "reset-token-123").
        Return(nil)

    mockAuditLog.EXPECT().
        Log(mock.Anything, "password_reset_requested", user.ID).
        Return(nil)

    // Act
    err := service.RequestReset(context.Background(), "user@example.com")

    // Assert
    assert.NoError(t, err)
}
```

### Step 3: Implement Service (Minimal)

```go
// internal/domain/services/password_reset_service.go
type PasswordResetService struct {
    userRepo     domain.UserRepository
    emailService EmailService
    tokenGen     TokenGenerator
    auditLog     AuditLogger
}

// Primary constructor - takes ALL dependencies
func NewPasswordResetService(
    userRepo domain.UserRepository,
    emailService EmailService,
    tokenGen TokenGenerator,
    auditLog AuditLogger,
) *PasswordResetService {
    return &PasswordResetService{
        userRepo:     userRepo,
        emailService: emailService,
        tokenGen:     tokenGen,
        auditLog:     auditLog,
    }
}

// Production factory - simple wiring only
// coverage:ignore
func NewPasswordResetServiceForProduction(
    db *sql.DB,
    emailCfg EmailConfig,
    logger Logger,
) *PasswordResetService {
    userRepo := persistence.NewUserRepository(db)
    emailService := email.NewSMTPService(emailCfg, logger)
    tokenGen := crypto.NewSecureTokenGenerator()
    auditLog := audit.NewDatabaseLogger(db)

    return NewPasswordResetService(userRepo, emailService, tokenGen, auditLog)
}

func (s *PasswordResetService) RequestReset(ctx context.Context, email string) error {
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return fmt.Errorf("finding user: %w", err)
    }

    token, err := s.tokenGen.Generate()
    if err != nil {
        return fmt.Errorf("generating token: %w", err)
    }

    if err := s.emailService.SendPasswordResetEmail(ctx, user, token); err != nil {
        return fmt.Errorf("sending email: %w", err)
    }

    if err := s.auditLog.Log(ctx, "password_reset_requested", user.ID); err != nil {
        // Log but don't fail - audit is important but not critical path
        // In production, you'd log this error
    }

    return nil
}
```

### Step 4: Implement Step Definition

```go
// features/step_definitions/password_reset_steps.go
func (s *PasswordResetSteps) requestPasswordResetFor(email string) error {
    // Use service from World (initialized in hooks with real container)
    err := s.world.PasswordResetService.RequestReset(
        context.Background(),
        email,
    )

    s.world.lastError = err
    return nil // Return nil - we check error in "Then" steps
}

func (s *PasswordResetSteps) shouldReceiveResetEmail() error {
    // Check test email inbox
    emails := s.world.TestMailbox.GetEmails()
    if len(emails) == 0 {
        return fmt.Errorf("expected reset email but none were sent")
    }

    lastEmail := emails[len(emails)-1]
    if !strings.Contains(lastEmail.Subject, "Password Reset") {
        return fmt.Errorf("expected password reset email, got: %s", lastEmail.Subject)
    }

    return nil
}
```

### Step 5: Run Scenario

```bash
just test-bdd
# OR
go test -tags=bdd ./features/ -godog.tags=@pending
```

### Step 6: Update Tags When Passing

```gherkin
@ story-PROJ-1234  # Changed from @pending
Scenario: User requests password reset successfully
  Given a user exists with email "user@example.com"
  When I request a password reset for "user@example.com"
  Then I should receive a reset email
  And the password reset attempt should be logged
```

## Common Patterns

### Pattern 1: Using Existing Services

```go
// features/support/world.go
type World struct {
    container            *ioc.Container
    PasswordResetService *services.PasswordResetService
    lastError            error
    testData             map[string]interface{}
}

// features/support/hooks.go
func InitializeHooks(ctx *godog.ScenarioContext, world *World) {
    ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
        // Create test container with real services but test infrastructure
        container, err := ioc.NewTestContainer()
        if err != nil {
            return ctx, err
        }
        world.container = container
        world.PasswordResetService = container.PasswordResetService()

        return ctx, nil
    })
}
```

### Pattern 2: Error Handling in Steps

```go
func (s *Steps) whenIDoSomething() error {
    // Capture error, don't return it immediately
    err := s.world.Service.DoSomething()
    s.world.lastError = err
    return nil
}

func (s *Steps) thenItShouldSucceed() error {
    // Check captured error in assertion step
    if s.world.lastError != nil {
        return fmt.Errorf("expected success but got: %w", s.world.lastError)
    }
    return nil
}

func (s *Steps) thenIShouldSee(message string) error {
    if s.world.lastError == nil {
        return fmt.Errorf("expected error but operation succeeded")
    }
    if !strings.Contains(s.world.lastError.Error(), message) {
        return fmt.Errorf("expected error containing '%s', got: %s",
            message, s.world.lastError.Error())
    }
    return nil
}
```

### Pattern 3: Test Data Setup

```go
func (s *Steps) userExistsWithEmail(email string) error {
    user := &domain.User{
        ID:    uuid.New().String(),
        Email: email,
    }

    err := s.world.container.UserRepository().Create(context.Background(), user)
    if err != nil {
        return fmt.Errorf("creating test user: %w", err)
    }

    // Store for later steps
    s.world.testData["currentUser"] = user
    return nil
}
```

## When to Amend Specs

Create `amend/{ticket-id}` branch if you discover:

❌ **Spec is technically infeasible**
- Example: API doesn't support required operation
- Action: Document why, propose alternative

❌ **Missing edge cases found during implementation**
- Example: Concurrent access scenario not covered
- Action: Add scenario, get BO re-approval

❌ **Step signature needs adjustment**
- Example: Need additional parameter for clarity
- Action: Update feature file + step definition in same PR, get BO approval

✅ **Implementation detail (NO amendment needed)**
- Example: Choosing internal data structure
- Action: Just implement

## PR Checklist

Before opening implementation PR:

- [ ] Implementation complete (scenarios remain `@pending`)
- [ ] All scenarios pass (`just test-bdd`)
- [ ] Unit tests pass (`just test-unit`)
- [ ] Coverage > 80% (excluding production factories)
- [ ] All production factories have `// coverage:ignore`
- [ ] No business logic in production factories
- [ ] Tests use primary constructors with mocks
- [ ] Error handling with proper wrapping
- [ ] Godoc comments on exported items

## Integration with Workflow

```
BO approves spec → PR merges with @pending tags →
  Dev creates impl/{ticket-id} branch →
    Implements step definitions (this prompt) →
      Scenarios pass →
        Changes scenarios remain @pending →
          Implementation PR → Merges to main
```

Your implementation makes the specification executable. The approved spec is the contract; your code fulfills it.
