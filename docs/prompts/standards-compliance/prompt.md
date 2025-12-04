# Go Standards Compliance Review

Review Go code for compliance with project standards. Focus on IoC pattern violations, especially business logic in production factories.

## Critical Rules

### IoC Pattern (REQUIRED)

**Every service MUST have:**
1. Primary constructor taking ALL dependencies
2. Production factory (`New*ForProduction`) taking only shared dependencies
3. Coverage exclusion on production factory (`// coverage:ignore`)

**Production factories MUST NOT contain:**
- Conditional logic (if/else, switch)
- Loops (for, range)
- Business rule validation
- Data transformation/calculations
- Configuration decisions based on business rules
- External API calls
- Database queries

**Production factories ONLY:**
- Create instances (New* calls)
- Pass dependencies to constructors
- Simple configuration loading

### Patterns

**✅ CORRECT:**
```go
// Primary constructor - ALL dependencies
func NewUserService(repo UserRepository, logger Logger, validator Validator) *UserService {
    return &UserService{repo: repo, logger: logger, validator: validator}
}

// Production factory - simple wiring only
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
    repo := persistence.NewUserRepository(db)
    validator := validation.NewUserValidator()
    return NewUserService(repo, logger, validator)
}
```

**If configuration varies:** Pass as shared dependency
```go
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger, validator Validator) *UserService {
    repo := persistence.NewUserRepository(db)
    return NewUserService(repo, logger, validator)
}

// Configuration decision logic - SHOULD BE TESTED
func buildValidator(cfg Config) Validator {
    if cfg.StrictMode {
        return validation.NewStrictValidator()
    }
    return validation.NewLenientValidator()
}

func TestBuildValidator(t *testing.T) {
    assert.IsType(t, &StrictValidator{}, buildValidator(Config{StrictMode: true}))
    assert.IsType(t, &LenientValidator{}, buildValidator(Config{StrictMode: false}))
}
```

**❌ VIOLATIONS:**
```go
// Missing primary constructor
func NewOrderServiceForProduction(db *sql.DB) *OrderService {
    return &OrderService{repo: persistence.NewOrderRepository(db)}
}

// Business logic in factory
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
    repo := persistence.NewUserRepository(db)
    if count, _ := repo.Count(); count > 1000 { // ❌ Business logic
        logger.Warn("High count")
    }
    return NewUserService(repo, logger)
}

// Configuration decision in factory
func NewPaymentServiceForProduction(db *sql.DB, cfg Config) *PaymentService {
    var processor PaymentProcessor
    if cfg.StrictMode { // ❌ Business decision
        processor = NewStrictProcessor()
    }
    return NewPaymentService(processor)
}
```

## Other Standards

**Directory Structure:**
- Domain: `/internal/domain/`
- Application: `/internal/application/`
- Infrastructure: `/internal/infrastructure/`
- IoC: `/internal/ioc/`
- Features: `/features/`
- Step definitions: `/features/step_definitions/`

**Testing:**
- Tests use primary constructor with mocks (never production factory)
- Test naming: `Test{Type}_{Method}_{Scenario}`
- File naming: `{name}_test.go`

**Error Handling:**
- Domain errors as package variables
- Wrap with context: `fmt.Errorf("context: %w", err)`

**Code Style:**
- Exported functions need godoc comments
- Use `gofmt` formatting
- Package names: lowercase, no underscores

**Coverage:**
- Target: 80%+ (excluding pure wiring)
- Exclude with `// coverage:ignore`:
  - `*ForProduction` functions (pure wiring, no decisions)
  - Container wiring (`NewContainer`, `init*` methods)
  - Infrastructure initialization
- DO NOT exclude (must test):
  - Configuration decision logic (`build*` with conditionals)
  - Any function with business-relevant decisions

## Output Format

**VIOLATIONS FOUND:** [count]

### Critical Violations (must fix)
1. `[file:line]` - [violation]
   - **Issue:** [description]
   - **Standard:** [which rule]
   - **Fix:** [how to fix]

### Warnings (should fix)
[list]

### Suggestions (consider)
[list]

**COMPLIANCE SCORE:** [percentage]

## Checklist

Per file:
- [ ] Primary constructor exists taking ALL dependencies?
- [ ] Production factory named `New*ForProduction`?
- [ ] Production factory takes only shared dependencies?
- [ ] **Production factory contains ZERO business logic?**
- [ ] Production factory marked `// coverage:ignore`?
- [ ] Tests use primary constructor with mocks?
- [ ] Files in correct directory?
- [ ] Exported items have godoc?
