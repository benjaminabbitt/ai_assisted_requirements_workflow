# Example Conversation: Basic Implementation with Plan-Review-Execute

**Scenario:** Implementing password reset feature (PROJ-1234)
**Duration:** ~45 minutes
**Outcome:** All scenarios passing, ready for PR

---

## Phase 1: Planning

**Developer:** Let's implement the password reset feature from PROJ-1234

**AI:** Planning implementation for PROJ-1234...

Reading context files:
- ✓ Read tech_standards.md
- ✓ Read architecture.md
- ✓ Read testing.md
- ✓ Introspected go.mod for internal dependencies
- ✓ Analyzed existing codebase patterns

**Implementation Plan:**

**Context Review:**
- **Tech Standards:** Go 1.21, Godog for BDD, primary constructor + production factory pattern
- **Architecture:** PostgreSQL database, SMTP email via SendGrid, Redis for rate limiting
- **Directory Structure:** Domain services in `internal/domain/services/`, step definitions in `features/step_definitions/`

**Internal Dependencies** (from go.mod):
- `github.com/company/auth/domain` - UserRepository available
- `github.com/company/notifications` - EmailService available
- `github.com/company/audit` - AuditLogger available
- `github.com/company/security` - Need to check for TokenGenerator

**Implementation Strategy:**
- **New Service:** PasswordResetService
- **Dependencies:**
  - UserRepository (exists - validates against user table)
  - EmailService (exists - sends via SendGrid)
  - TokenGenerator (need to verify if exists)
  - AuditLogger (exists - logs to audit_events table)

**Files to Create:**
1. `internal/domain/services/password_reset_service.go` - Main service
2. `internal/domain/services/password_reset_service_test.go` - Unit tests with mocks
3. `features/step_definitions/password_reset_steps.go` - BDD step definitions

**Test Approach:**
- Unit tests using primary constructor with mocks (testify/mock)
- BDD scenarios using test container with real database
- Coverage target: >80% excluding production factories

**Scenarios from Approved Spec:**
1. ✓ User requests password reset successfully
2. ✓ User requests reset with invalid email format
3. ✓ User requests reset with rate limit exceeded

Ready to proceed?

---

**Developer:** Good plan. Before we start, check if TokenGenerator exists in the security package

**AI:** Searching github.com/company/security package...

Found: `security/tokens/generator.go`
```go
type Generator interface {
    Generate() (string, error)
    GenerateWithExpiry(duration time.Duration) (string, time.Time, error)
}
```

Should we use this existing TokenGenerator instead of creating new?

---

**Developer:** Yes, use the existing one. Update your plan

**AI:** Updated Implementation Plan:

**Dependencies (Final):**
- ✓ UserRepository (auth/domain)
- ✓ EmailService (notifications)
- ✓ tokens.Generator (security) - **reusing existing**
- ✓ AuditLogger (audit)

**Files to Create:**
1. `internal/domain/services/password_reset_service.go`
2. `internal/domain/services/password_reset_service_test.go`
3. `features/step_definitions/password_reset_steps.go`

No new infrastructure components needed—all dependencies exist.

Proceed with implementation?

---

**Developer:** Yes. Start with TDD—write the unit test for the first scenario

---

## Phase 2: Execution - Scenario 1 (Success Case)

**AI:** Starting TDD cycle for Scenario 1: "User requests password reset successfully"

Writing unit test using primary constructor with mocks:

```go
// internal/domain/services/password_reset_service_test.go
package services_test

import (
    "context"
    "testing"

    "github.com/company/app/internal/domain/services"
    "github.com/company/app/internal/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

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

    user := &domain.User{
        ID:    "user-123",
        Email: "user@example.com",
    }

    mockUserRepo.EXPECT().
        FindByEmail(mock.Anything, "user@example.com").
        Return(user, nil)

    mockTokenGen.EXPECT().
        GenerateWithExpiry(24 * time.Hour).
        Return("reset-token-abc", time.Now().Add(24*time.Hour), nil)

    mockEmailService.EXPECT().
        SendPasswordResetEmail(mock.Anything, user, "reset-token-abc").
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

Test written. Should I run it to verify it fails correctly?

---

**Developer:** Yes, run the test

**AI:** Running: `go test ./internal/domain/services/...`

```
--- FAIL: TestPasswordResetService_RequestReset_Success (0.00s)
    password_reset_service_test.go:25:
        Error: undefined: services.NewPasswordResetService
FAIL
FAIL    github.com/company/app/internal/domain/services    0.002s
```

Test fails as expected—`NewPasswordResetService` doesn't exist yet.

Proceeding to implementation...

---

**AI:** Implementing PasswordResetService following tech_standards.md IoC pattern:

```go
// internal/domain/services/password_reset_service.go
package services

import (
    "context"
    "fmt"
    "time"

    "github.com/company/app/internal/domain"
    "github.com/company/audit"
    "github.com/company/notifications"
    "github.com/company/security/tokens"
)

type PasswordResetService struct {
    userRepo     domain.UserRepository
    emailService notifications.EmailService
    tokenGen     tokens.Generator
    auditLog     audit.Logger
}

// Primary constructor - takes ALL dependencies
func NewPasswordResetService(
    userRepo domain.UserRepository,
    emailService notifications.EmailService,
    tokenGen tokens.Generator,
    auditLog audit.Logger,
) *PasswordResetService {
    return &PasswordResetService{
        userRepo:     userRepo,
        emailService: emailService,
        tokenGen:     tokenGen,
        auditLog:     auditLog,
    }
}

// Production factory - simple wiring only, no business logic
// coverage:ignore
func NewPasswordResetServiceForProduction(
    db *sql.DB,
    emailCfg notifications.Config,
    logger log.Logger,
) *PasswordResetService {
    userRepo := persistence.NewUserRepository(db)
    emailService := notifications.NewSMTPService(emailCfg, logger)
    tokenGen := tokens.NewSecureGenerator()
    auditLog := audit.NewDatabaseLogger(db)

    return NewPasswordResetService(userRepo, emailService, tokenGen, auditLog)
}

func (s *PasswordResetService) RequestReset(ctx context.Context, email string) error {
    // Find user
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return fmt.Errorf("finding user: %w", err)
    }

    // Generate reset token
    token, expiry, err := s.tokenGen.GenerateWithExpiry(24 * time.Hour)
    if err != nil {
        return fmt.Errorf("generating token: %w", err)
    }

    // Store token (simplified - would save to DB in real implementation)
    // TODO: Save token and expiry to password_reset_tokens table

    // Send email
    if err := s.emailService.SendPasswordResetEmail(ctx, user, token); err != nil {
        return fmt.Errorf("sending email: %w", err)
    }

    // Audit log
    if err := s.auditLog.Log(ctx, "password_reset_requested", user.ID); err != nil {
        // Log error but don't fail - audit is important but not critical path
        // In production: log this error to monitoring system
    }

    return nil
}
```

Implementation complete. Running test...

---

**AI:** Running: `go test ./internal/domain/services/...`

```
--- PASS: TestPasswordResetService_RequestReset_Success (0.00s)
PASS
ok      github.com/company/app/internal/domain/services    0.003s
```

✓ Test passes!

---

**Developer:** Good. Now review the implementation for any issues

**AI:** Reviewing implementation against engineering checklist...

**Issues found:**

1. **Missing:** Token storage - TODO comment but no implementation
2. **Security:** No input validation on email parameter
3. **Architecture:** No rate limiting (architecture.md specifies 5 attempts/hour)
4. **Error handling:** Silent failure on audit logging could mask issues

Should I address these?

---

**Developer:** Yes, fix all four issues

**AI:** Fixing issues...

**1. Token storage** - Adding repository method and call:
```go
// Store token with expiry
if err := s.userRepo.SavePasswordResetToken(ctx, user.ID, token, expiry); err != nil {
    return fmt.Errorf("storing reset token: %w", err)
}
```

**2. Input validation** - Adding email format check:
```go
func (s *PasswordResetService) RequestReset(ctx context.Context, email string) error {
    // Validate email format
    if !isValidEmail(email) {
        return fmt.Errorf("invalid email format: %s", email)
    }

    // ... rest of implementation
}
```

**3. Rate limiting** - Adding check:
```go
// Check rate limit (5 attempts per hour per IP/user)
if err := s.rateLimiter.CheckLimit(ctx, email, 5, time.Hour); err != nil {
    return fmt.Errorf("rate limit exceeded: %w", err)
}
```

Wait—I need to add RateLimiter to dependencies. Should I update the service?

---

**Developer:** Yes, add RateLimiter dependency. Update the constructor and tests

**AI:** Updating service with RateLimiter dependency...

**Updated constructor:**
```go
type PasswordResetService struct {
    userRepo     domain.UserRepository
    emailService notifications.EmailService
    tokenGen     tokens.Generator
    auditLog     audit.Logger
    rateLimiter  RateLimiter  // NEW
}

func NewPasswordResetService(
    userRepo domain.UserRepository,
    emailService notifications.EmailService,
    tokenGen tokens.Generator,
    auditLog audit.Logger,
    rateLimiter RateLimiter,  // NEW parameter
) *PasswordResetService {
    return &PasswordResetService{
        userRepo:     userRepo,
        emailService: emailService,
        tokenGen:     tokenGen,
        auditLog:     auditLog,
        rateLimiter:  rateLimiter,
    }
}
```

**Updated test:**
```go
func TestPasswordResetService_RequestReset_Success(t *testing.T) {
    mockUserRepo := mocks.NewUserRepository(t)
    mockEmailService := mocks.NewEmailService(t)
    mockTokenGen := mocks.NewTokenGenerator(t)
    mockAuditLog := mocks.NewAuditLogger(t)
    mockRateLimiter := mocks.NewRateLimiter(t)  // NEW

    service := services.NewPasswordResetService(
        mockUserRepo,
        mockEmailService,
        mockTokenGen,
        mockAuditLog,
        mockRateLimiter,  // NEW
    )

    // Add rate limiter expectation
    mockRateLimiter.EXPECT().
        CheckLimit(mock.Anything, "user@example.com", 5, time.Hour).
        Return(nil)

    // ... rest of test
}
```

**4. Audit logging** - Now returning error:
```go
// Audit log - fail if audit fails (changed from silent)
if err := s.auditLog.Log(ctx, "password_reset_requested", user.ID); err != nil {
    return fmt.Errorf("audit logging failed: %w", err)
}
```

All issues addressed. Running tests...

```
--- PASS: TestPasswordResetService_RequestReset_Success (0.00s)
PASS
ok      github.com/company/app/internal/domain/services    0.004s
```

✓ Tests still pass with all fixes applied!

---

**Developer:** Excellent. Now implement the BDD step definition

**AI:** Implementing step definition for "When I request a password reset for {string}":

```go
// features/step_definitions/password_reset_steps.go
package step_definitions

import (
    "context"
    "fmt"

    "github.com/cucumber/godog"
)

type PasswordResetSteps struct {
    world *World
}

func NewPasswordResetSteps(world *World) *PasswordResetSteps {
    return &PasswordResetSteps{world: world}
}

func (s *PasswordResetSteps) RegisterSteps(sc *godog.ScenarioContext) {
    sc.Step(`^I request a password reset for "([^"]*)"$`, s.requestPasswordResetFor)
    sc.Step(`^I should receive a reset email$`, s.shouldReceiveResetEmail)
    sc.Step(`^the password reset attempt should be logged$`, s.resetAttemptShouldBeLogged)
}

func (s *PasswordResetSteps) requestPasswordResetFor(email string) error {
    // Use service from World (initialized with test container)
    err := s.world.PasswordResetService.RequestReset(
        context.Background(),
        email,
    )

    s.world.lastError = err
    return nil // Don't return error here - check in "Then" steps
}

func (s *PasswordResetSteps) shouldReceiveResetEmail() error {
    // Check test email inbox
    emails := s.world.TestMailbox.GetEmails()
    if len(emails) == 0 {
        return fmt.Errorf("expected reset email but none were sent")
    }

    lastEmail := emails[len(emails)-1]
    if !strings.Contains(lastEmail.Subject, "Password Reset") {
        return fmt.Errorf("expected password reset email, got subject: %s", lastEmail.Subject)
    }

    return nil
}

func (s *PasswordResetSteps) resetAttemptShouldBeLogged() error {
    // Check audit log
    logs := s.world.TestAuditLog.GetLogs()
    for _, log := range logs {
        if log.Event == "password_reset_requested" {
            return nil // Found the log entry
        }
    }

    return fmt.Errorf("password reset attempt not found in audit log")
}
```

Step definitions complete. Running BDD scenario...

---

**AI:** Running: `go test -tags=bdd ./features/ -godog.tags='@story-PROJ-1234'`

```
Feature: Password Reset
  @pending @story-PROJ-1234
  Scenario: User requests password reset successfully
    Given a user exists with email "user@example.com"
    When I request a password reset for "user@example.com"
    Then I should receive a reset email
    And the password reset attempt should be logged

1 scenarios (1 passed)
4 steps (4 passed)
```

✓ Scenario 1 passes!

**Coverage check:**
```
go test -cover ./internal/domain/services/
ok      github.com/company/app/internal/domain/services    coverage: 89.2% of statements
```

✓ Coverage: 89.2% (target: >80%)

Scenario 1 complete. Ready for scenario 2?

---

**Developer:** Yes, continue with scenario 2

**AI:** Moving to Scenario 2: "User requests reset with invalid email format"

This should trigger validation error. Writing test...

[Implementation continues for remaining scenarios...]

---

## Summary

**Time:** 45 minutes
**Scenarios implemented:** 3/3
**Tests:** All passing
**Coverage:** 89.2%
**Issues found during review:** 4 (all fixed)
**Refactorings:** 1 (added RateLimiter dependency)

**Key collaboration points:**
- Developer caught missing TokenGenerator check in plan
- Developer requested dependency review before implementation
- Developer guided security and architecture fixes
- AI proactively identified issues during review

**Ready for:** Remove @pending tags, open implementation PR
