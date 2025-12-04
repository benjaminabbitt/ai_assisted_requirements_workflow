# Standards Compliance Review Prompt

## Purpose

Review Go code for compliance with project technical standards, with special emphasis on the IoC pattern and the prohibition of business logic in production factories.

## AI vs. Human Responsibilities

### What AI Handles (Automated Compliance)
- Detecting pattern violations (missing constructors, business logic in factories)
- Identifying structural issues (wrong directory, missing documentation)
- Checking naming conventions and style compliance
- Flagging common anti-patterns
- Generating compliance reports with specific violations

### What Humans Handle (Software Engineering)
- **Refinement of generated code** - Improving clarity, maintainability, and elegance
- **Deduplication** - Identifying and consolidating repeated patterns across codebase
- **Software engineering decisions** - Architecture choices, abstraction levels, API design
- **Context and intent** - Understanding why code exists and whether it solves the right problem
- **Trade-off evaluation** - Balancing performance, maintainability, and complexity
- **Code review beyond compliance** - Logic correctness, algorithm efficiency, security implications
- **Mentoring and knowledge transfer** - Teaching patterns and explaining rationale to team members

**Key Principle:** AI ensures code follows the rules. Humans ensure code is good software engineering.

## Files in This Directory

- `prompt.md` - Minimal, token-efficient compliance review prompt for AI
- `sample-violations.go` - Example code with violations (for testing)
- `sample-correct.go` - Example code following standards

## Understanding the Standards

### Why These Patterns?

**Primary Constructor + Production Factory Pattern:**

The separation exists for one reason: **testability**.

- **Primary constructor** takes ALL dependencies as parameters → You can inject mocks in tests
- **Production factory** builds dependencies internally → Used only in production, excluded from coverage

**Why exclude production factories from coverage?**

They're infrastructure glue, not business logic. Testing them would require:
- Real database connections
- Real external services
- Real file systems

This defeats the purpose of unit testing. Instead:
- Business logic lives in services (100% testable via primary constructor)
- Infrastructure wiring lives in factories (excluded from coverage, tested via integration/E2E)

**Why prohibit business logic in production factories?**

If business logic exists in a production factory, it becomes **untestable**:

```go
// ❌ This business logic cannot be unit tested
func NewUserServiceForProduction(db *sql.DB) *UserService {
    repo := persistence.NewUserRepository(db)

    // Business rule: warn if > 1000 users
    if count, _ := repo.Count(); count > 1000 {  // Requires real DB!
        logger.Warn("High count")
    }

    return NewUserService(repo)
}
```

Moving it to a service method makes it testable:

```go
// ✅ Testable - no real DB needed
func TestUserService_CheckCapacity(t *testing.T) {
    mockRepo := mocks.NewUserRepository(t)
    mockRepo.EXPECT().Count().Return(1500, nil)

    service := NewUserService(mockRepo, logger)
    err := service.CheckCapacity()  // Fully testable

    assert.NoError(t, err)
}
```

**When configuration varies:**

If you need different implementations based on configuration, make it a **shared dependency**:

```go
// Configuration layer decides (SHOULD BE TESTED - contains decision logic)
func buildValidator(cfg Config) Validator {
    if cfg.StrictMode {
        return validation.NewStrictValidator()
    }
    return validation.NewLenientValidator()
}

// Test the decision logic
func TestBuildValidator(t *testing.T) {
    strictCfg := Config{StrictMode: true}
    assert.IsType(t, &validation.StrictValidator{}, buildValidator(strictCfg))

    lenientCfg := Config{StrictMode: false}
    assert.IsType(t, &validation.LenientValidator{}, buildValidator(lenientCfg))
}

// Factory receives pre-made decision (excluded from coverage)
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger, validator Validator) *UserService {
    repo := persistence.NewUserRepository(db)
    return NewUserService(repo, logger, validator)
}
```

This keeps decisions testable (simple unit tests, no infrastructure) and factories simple (pure wiring).

### The Coverage Strategy

**What we cover:**
- ✅ Business logic in services (via primary constructor + mocks)
- ✅ Domain rules and validations
- ✅ Use case orchestration
- ✅ Error handling paths
- ✅ **Configuration decision logic** (`buildValidator`, `buildProcessor`) - simple tests, no infrastructure

**What we exclude:**
- ❌ Production factories (`New*ForProduction`) - pure wiring, no decisions
- ❌ Container wiring (`NewContainer`) - infrastructure initialization
- ❌ Infrastructure initialization (database, logger, external service connections)

**Key distinction:**
- **Decision logic** = test it (even if it's configuration-related)
- **Pure wiring** = exclude it (no logic to test)

**Result:** 80%+ coverage of business logic + decision logic, 0% coverage of infrastructure glue.

## 🤖 ACTUAL SUBAGENT EXECUTION

The subagent performed a comprehensive standards compliance review and identified **24 violations** with a **35% compliance score**.

### Key Findings

**Critical Violations (16):**
- Missing primary constructor (OrderService)
- Business logic in 6 production factories
- All 7 factories missing `// coverage:ignore` markers
- Test using production factory instead of primary constructor

**Warning Violations (8):**
- Missing godoc comments on exported types and functions

**Result:** Code requires significant refactoring to meet standards.

### What the AI Caught

✅ **Pattern violations:**
- Detected missing primary constructor
- Identified business logic in factories (conditionals, loops, calculations, API calls)
- Found test using production factory (anti-pattern)
- Flagged all missing coverage exclusion markers

✅ **Specific issues per service:**
- OrderService: Direct instantiation, no primary constructor
- UserService: Count check logic in factory
- PaymentService: Configuration decision in factory
- NotificationService: Calculation in factory
- ReportService: Loop with switch in factory
- WeatherService: External API call in factory
- DataService: Data transformation in factory

✅ **Code quality:**
- Missing documentation on 8 exported items
- Provided specific fixes with code examples for each violation

### Compliance Score Breakdown

- **35% compliant** - Significant work needed
- **0 out of 7 services** fully compliant
- **16 critical** IoC violations
- **8 documentation** warnings

**Priority fixes:**
1. Add primary constructor to OrderService
2. Extract business logic from all factories → separate testable functions
3. Add `// coverage:ignore` to all 7 factories
4. Fix test to use primary constructor + mocks
5. Add godoc comments

### Full Output

See complete compliance report with:
- 24 specific violations with file:line references
- Issue descriptions and violated standards
- Concrete fix examples for each issue
- Per-service compliance checklist

This demonstrates the standards compliance prompt catches all IoC violations and provides actionable guidance for fixes.

---

## Quick Usage

### For AI Code Review

```bash
# Read the prompt and sample code
cat prompt.md

# Then provide code to review with:
"Review the following code against the standards defined in prompt.md:

[PASTE CODE HERE]
"
```

### For Claude Subagent

Use the Task tool to launch a review agent:

```
Review the code in sample-violations.go against the standards in prompt.md
```

## Example Execution with Claude Subagent

### Scenario: Review Code with Multiple Violations

**Input Code:** `sample-violations.go`

**Execution:**

Using Claude Code's Task tool with a general-purpose agent:

```markdown
**Prompt to Claude:**
"I need you to review sample-violations.go against our project standards.
Read docs/prompts/standards-compliance/prompt.md for the full standards,
then review docs/prompts/standards-compliance/sample-violations.go and
provide a detailed compliance report."
```

**Expected Output:**

```markdown
## Standards Compliance Review

**VIOLATIONS FOUND:** 4

### Critical Violations (must fix)

1. `sample-violations.go:25` - Business logic in production factory
   - **Issue:** Production factory `NewUserServiceForProduction` contains conditional
     logic checking user count (lines 25-28)
   - **Standard:** Production factories MUST contain zero business logic
     (tech_standards.md § Dependency Injection)
   - **Fix:** Remove the user count check. If this validation is needed, it belongs
     in a service method, not the factory.

2. `sample-violations.go:15` - Missing primary constructor
   - **Issue:** `OrderService` only has production factory, no primary constructor
     taking all dependencies
   - **Standard:** Every component requires primary constructor for testing
     (tech_standards.md § Dependency Injection)
   - **Fix:** Add `NewOrderService(repo OrderRepository, logger Logger, calculator
     PriceCalculator)` that takes all dependencies

3. `sample-violations.go:42` - Configuration decision in production factory
   - **Issue:** Production factory uses if/else to choose validator based on config
     (business rule)
   - **Standard:** Production factories cannot make business decisions
     (tech_standards.md § Dependency Injection)
   - **Fix:** Move validator selection logic out of factory. Create separate factories
     for different configurations or pass validator as dependency

4. `sample-violations.go:58` - Business calculation in production factory
   - **Issue:** Production factory calculates default timeout value based on
     environment
   - **Standard:** Production factories cannot contain calculations
     (tech_standards.md § Dependency Injection)
   - **Fix:** Move timeout calculation to configuration loading, pass computed value
     to factory

### Warnings (should fix)

1. `sample-violations.go:35` - Missing godoc comment
   - Exported function `ProcessPayment` lacks documentation

2. `sample-violations.go:72` - Incorrect test constructor usage
   - Test uses production factory instead of primary constructor with mocks

### Suggestions (consider)

1. `sample-violations.go:80` - Consider table-driven test
   - Multiple similar test cases could be consolidated

**COMPLIANCE SCORE:** 45%

**PRIORITY ACTIONS:**
1. Add primary constructors to all services
2. Remove all business logic from production factories
3. Add coverage exclusion markers to production factories
```

### Verification Steps

After fixes, re-run the review:

```markdown
**Prompt to Claude:**
"Review sample-correct.go against standards in prompt.md. This should have
100% compliance."
```

**Expected Output:**

```markdown
## Standards Compliance Review

**VIOLATIONS FOUND:** 0

**COMPLIANCE SCORE:** 100%

✅ All IoC patterns correct
✅ Primary constructors present
✅ Production factories contain zero business logic
✅ Proper coverage exclusion markers
✅ Tests use primary constructors with mocks
✅ Directory structure compliant
✅ Error handling follows standards
```

## Example Execution: CI/CD Integration

### GitHub Actions Workflow

```yaml
name: Standards Compliance

on: [pull_request]

jobs:
  compliance-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Review changed Go files
        run: |
          # Get changed .go files
          CHANGED_FILES=$(git diff --name-only origin/main...HEAD | grep '\.go$' || true)

          if [ -n "$CHANGED_FILES" ]; then
            # Use Claude API or custom linter with standards prompt
            for file in $CHANGED_FILES; do
              echo "Reviewing $file..."
              # Call your compliance checker
              ./scripts/check-compliance.sh "$file"
            done
          fi
```

## Example Execution: Manual Review Workflow

### Step 1: Developer Submits PR

```bash
# Developer creates PR with new service
git checkout -b feature/new-payment-service
# ... make changes ...
git push origin feature/new-payment-service
```

### Step 2: Tech Lead Reviews

```markdown
**Tech Lead Action:**
1. Open PR in browser
2. Copy service code
3. Paste into Claude with standards prompt
4. Review compliance report
5. Request changes or approve
```

### Step 3: Developer Fixes Violations

```bash
# Developer addresses feedback
# Example fix: Move business decision out, pass as dependency

# BEFORE (violation):
func NewPaymentServiceForProduction(db *sql.DB, logger Logger) *PaymentService {
    repo := persistence.NewPaymentRepository(db)

    // ❌ VIOLATION: Business logic - choosing processor based on db type
    var processor PaymentProcessor
    if db.Type() == "postgres" {
        processor = NewPostgresProcessor()
    } else {
        processor = NewMySQLProcessor()
    }

    return NewPaymentService(repo, logger, processor)
}

# AFTER (compliant):
// Factory now accepts processor as shared dependency
// coverage:ignore
func NewPaymentServiceForProduction(db *sql.DB, logger Logger, processor PaymentProcessor) *PaymentService {
    repo := persistence.NewPaymentRepository(db)
    // Simple wiring only - no business decisions
    return NewPaymentService(repo, logger, processor)
}

// Configuration layer makes the business decision (in container or main)
// This SHOULD be tested - it contains decision logic
func buildPaymentProcessor(db *sql.DB) PaymentProcessor {
    if db.Type() == "postgres" {
        return processors.NewPostgresProcessor()
    }
    return processors.NewMySQLProcessor()
}

// Test the decision
func TestBuildPaymentProcessor(t *testing.T) {
    postgresDB := &mockDB{dbType: "postgres"}
    assert.IsType(t, &processors.PostgresProcessor{}, buildPaymentProcessor(postgresDB))

    mysqlDB := &mockDB{dbType: "mysql"}
    assert.IsType(t, &processors.MySQLProcessor{}, buildPaymentProcessor(mysqlDB))
}

// Container wires it together
// coverage:ignore
func (c *Container) initPaymentService() {
    processor := buildPaymentProcessor(c.db)
    c.paymentService = NewPaymentServiceForProduction(c.db, c.logger, processor)
}
```

### Step 4: Re-review and Merge

```markdown
**Tech Lead Action:**
1. Re-run compliance check
2. Verify 100% compliance
3. Approve PR
4. Merge to main
```

## Common Violations and Fixes

### Violation 1: Business Logic in Production Factory

```go
// ❌ VIOLATION
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
    repo := persistence.NewUserRepository(db)

    // Business logic - checking count
    count, _ := repo.Count()
    if count > 1000 {
        logger.Warn("High user count")
    }

    return NewUserService(repo, logger)
}

// ✅ FIX
func NewUserService(repo UserRepository, logger Logger) *UserService {
    return &UserService{repo: repo, logger: logger}
}

// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
    repo := persistence.NewUserRepository(db) // Only wiring
    return NewUserService(repo, logger)
}

// Move business logic to service method
func (s *UserService) CheckCapacity() error {
    count, err := s.repo.Count()
    if err != nil {
        return err
    }
    if count > 1000 {
        s.logger.Warn("High user count")
    }
    return nil
}
```

### Violation 2: Missing Primary Constructor

```go
// ❌ VIOLATION - Only production factory exists
func NewOrderServiceForProduction(db *sql.DB) *OrderService {
    repo := persistence.NewOrderRepository(db)
    return &OrderService{repo: repo}
}

// ✅ FIX - Add primary constructor
func NewOrderService(repo OrderRepository) *OrderService {
    return &OrderService{repo: repo}
}

// coverage:ignore
func NewOrderServiceForProduction(db *sql.DB) *OrderService {
    repo := persistence.NewOrderRepository(db)
    return NewOrderService(repo)
}
```

### Violation 3: Configuration Decision in Factory

```go
// ❌ VIOLATION
func NewUserServiceForProduction(db *sql.DB, logger Logger, cfg Config) *UserService {
    repo := persistence.NewUserRepository(db)

    // ❌ VIOLATION: Business decision based on config
    var validator Validator
    if cfg.StrictMode {
        validator = validation.NewStrictValidator()
    } else {
        validator = validation.NewLenientValidator()
    }

    return NewUserService(repo, logger, validator)
}

// ✅ FIX - Pass validator as shared dependency
func NewUserService(repo UserRepository, logger Logger, validator Validator) *UserService {
    return &UserService{repo: repo, logger: logger, validator: validator}
}

// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger, validator Validator) *UserService {
    repo := persistence.NewUserRepository(db)
    // Simple wiring - validator already decided by caller
    return NewUserService(repo, logger, validator)
}

// Configuration layer makes the business decision (in container or main)
// This SHOULD be tested - it contains decision logic
func buildValidator(cfg Config) Validator {
    if cfg.StrictMode {
        return validation.NewStrictValidator()
    }
    return validation.NewLenientValidator()
}

// Test the decision
func TestBuildValidator(t *testing.T) {
    strictCfg := Config{StrictMode: true}
    assert.IsType(t, &validation.StrictValidator{}, buildValidator(strictCfg))

    lenientCfg := Config{StrictMode: false}
    assert.IsType(t, &validation.LenientValidator{}, buildValidator(lenientCfg))
}

// Container wires it together
// coverage:ignore
func (c *Container) initUserService() {
    validator := buildValidator(c.cfg)
    c.userService = NewUserServiceForProduction(c.db, c.logger, validator)
}
```

## Integration with Development Workflow

### Pre-Commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Check staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')

if [ -n "$STAGED_GO_FILES" ]; then
    echo "Running standards compliance check..."

    for file in $STAGED_GO_FILES; do
        # Quick check for common violations
        if grep -A 10 "ForProduction" "$file" | grep -E "if|for|switch" > /dev/null; then
            echo "⚠️  WARNING: Possible business logic in production factory: $file"
            echo "   Please review against docs/prompts/standards-compliance/prompt.md"
        fi
    done
fi
```

### IDE Integration

**VS Code Task:**

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Standards Compliance Check",
            "type": "shell",
            "command": "claude-code",
            "args": [
                "review",
                "${file}",
                "--standards",
                "docs/prompts/standards-compliance/prompt.md"
            ],
            "problemMatcher": [],
            "group": "test"
        }
    ]
}
```

## Success Metrics

Track these metrics over time:

- **Compliance Score:** Average across all reviews
- **Common Violations:** Most frequent issues
- **Time to Fix:** Average time from violation to resolution
- **Violation Trend:** Decreasing over time indicates learning

**Target Metrics:**
- Compliance Score: > 95%
- Business Logic in Factories: 0 violations
- Primary Constructor Coverage: 100%
- Review Time: < 10 minutes per file

## References

- Main Standards: `../../tech_standards.md`
- IoC Pattern: `../../architecture.md § Inversion of Control`
- Personas: `../personas.md`
