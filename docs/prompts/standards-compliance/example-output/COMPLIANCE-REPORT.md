# Standards Compliance Review: Sample Codebase

**Date:** 2025-12-04
**Agent:** standards-compliance
**Target:** sample-violations.go
**Standards:** tech_standards.md, architecture.md (IoC patterns)

---

## Executive Summary

**Compliance Score: 35%**
- **Critical Violations:** 16
- **Warnings:** 8
- **Services Reviewed:** 7
- **Services Compliant:** 0

**Recommendation:** REJECT - Significant refactoring required before merge.

---

## Detailed Findings

### Critical Violations (Must Fix Before Merge)

#### 1. Missing Primary Constructor: OrderService

**Location:** `sample-violations.go:15-30`

**Issue:** OrderService has no primary constructor that accepts all dependencies as parameters.

**Violated Standard:**
> Every service MUST have a primary constructor that takes ALL dependencies as parameters. (tech_standards.md § IoC Patterns)

**Current Code:**
```go
type OrderService struct {
    repo   OrderRepository
    logger Logger
}

// Only has production factory - no primary constructor!
func NewOrderServiceForProduction(db *sql.DB) *OrderService {
    return &OrderService{
        repo:   persistence.NewOrderRepository(db),
        logger: logging.NewLogger(),
    }
}
```

**Impact:** Service cannot be unit tested with mocks. Tests must use production factory with real database.

**Required Fix:**
```go
// Primary constructor - takes ALL dependencies
func NewOrderService(repo OrderRepository, logger Logger) *OrderService {
    return &OrderService{
        repo:   repo,
        logger: logger,
    }
}

// Production factory - calls primary constructor
// coverage:ignore
func NewOrderServiceForProduction(db *sql.DB) *OrderService {
    repo := persistence.NewOrderRepository(db)
    logger := logging.NewLogger()
    return NewOrderService(repo, logger)
}
```

---

#### 2. Business Logic in Factory: UserService

**Location:** `sample-violations.go:40-48`

**Issue:** Production factory contains conditional business logic (user count check).

**Violated Standard:**
> Production factories MUST NOT contain business logic (conditionals, loops, calculations, validations). (tech_standards.md § IoC Patterns - Business Logic Prohibition)

**Current Code:**
```go
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB) *UserService {
    repo := persistence.NewUserRepository(db)

    // ❌ Business logic in factory - requires real DB to test!
    if count, _ := repo.Count(); count > 1000 {
        logger.Warn("User count exceeds threshold")
    }

    return NewUserService(repo, logging.NewLogger())
}
```

**Impact:**
- Business rule (warning at >1000 users) is untestable without real database
- Cannot mock repository in tests for this logic
- Coverage reporting incorrectly excludes this business rule

**Required Fix:**

Extract to service method that can be tested:

```go
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB) *UserService {
    repo := persistence.NewUserRepository(db)
    logger := logging.NewLogger()
    service := NewUserService(repo, logger)

    // Call initialization method (which can be tested)
    service.Initialize()

    return service
}

// Testable business logic
func (s *UserService) Initialize() error {
    count, err := s.repo.Count()
    if err != nil {
        return fmt.Errorf("checking count: %w", err)
    }

    if count > 1000 {
        s.logger.Warn("User count exceeds threshold")
    }

    return nil
}

// Unit test - no real DB needed!
func TestUserService_Initialize_WarnsAtThreshold(t *testing.T) {
    mockRepo := mocks.NewUserRepository(t)
    mockLogger := mocks.NewLogger(t)

    mockRepo.EXPECT().Count().Return(1500, nil)
    mockLogger.EXPECT().Warn("User count exceeds threshold")

    service := NewUserService(mockRepo, mockLogger)
    err := service.Initialize()

    assert.NoError(t, err)
}
```

---

#### 3. Configuration Decision Logic in Factory: PaymentService

**Location:** `sample-violations.go:58-68`

**Issue:** Production factory contains configuration decision logic (choosing processor based on environment).

**Violated Standard:**
> Production factories MUST NOT contain business logic (conditionals, loops, calculations, validations). (tech_standards.md § IoC Patterns - Business Logic Prohibition)

**Current Code:**
```go
// coverage:ignore
func NewPaymentServiceForProduction(cfg Config) *PaymentService {
    repo := persistence.NewPaymentRepository(cfg.DB)

    var processor PaymentProcessor
    // ❌ Configuration decision in factory - untestable!
    if cfg.Environment == "production" {
        processor = stripe.NewProcessor(cfg.StripeKey)
    } else {
        processor = mock.NewProcessor()
    }

    return NewPaymentService(repo, processor)
}
```

**Impact:**
- Decision logic (prod vs mock processor) is excluded from coverage
- Cannot unit test the decision without real configuration
- Configuration mistakes won't be caught by tests

**Required Fix:**

Extract decision to separate testable function, pass result as parameter:

```go
// Configuration decision function - SHOULD BE TESTED
func buildPaymentProcessor(cfg Config) PaymentProcessor {
    if cfg.Environment == "production" {
        return stripe.NewProcessor(cfg.StripeKey)
    }
    return mock.NewProcessor()
}

// Test the decision logic
func TestBuildPaymentProcessor_Production(t *testing.T) {
    cfg := Config{Environment: "production", StripeKey: "sk_test"}
    processor := buildPaymentProcessor(cfg)
    assert.IsType(t, &stripe.Processor{}, processor)
}

func TestBuildPaymentProcessor_NonProduction(t *testing.T) {
    cfg := Config{Environment: "development"}
    processor := buildPaymentProcessor(cfg)
    assert.IsType(t, &mock.Processor{}, processor)
}

// Factory receives pre-made decision
// coverage:ignore
func NewPaymentServiceForProduction(cfg Config) *PaymentService {
    repo := persistence.NewPaymentRepository(cfg.DB)
    processor := buildPaymentProcessor(cfg)  // Separately tested
    return NewPaymentService(repo, processor)
}
```

---

#### 4. Calculation in Factory: NotificationService

**Location:** `sample-violations.go:78-86`

**Issue:** Production factory performs calculation (rate limit from timeout).

**Violated Standard:**
> Production factories MUST NOT contain business logic (conditionals, loops, calculations, validations). (tech_standards.md § IoC Patterns - Business Logic Prohibition)

**Current Code:**
```go
// coverage:ignore
func NewNotificationServiceForProduction(cfg Config) *NotificationService {
    mailer := email.NewSMTPMailer(cfg.SMTPConfig)

    // ❌ Calculation in factory - untestable!
    rateLimit := cfg.Timeout.Seconds() / 10

    return NewNotificationService(mailer, rateLimit)
}
```

**Impact:**
- Business rule (rate limit formula) is untestable
- Formula changes won't be caught by tests
- Excluded from coverage despite being business logic

**Required Fix:**

Move calculation to configuration layer:

```go
// Configuration function - testable
func calculateRateLimit(timeout time.Duration) float64 {
    return timeout.Seconds() / 10
}

func TestCalculateRateLimit(t *testing.T) {
    timeout := 100 * time.Second
    rateLimit := calculateRateLimit(timeout)
    assert.Equal(t, 10.0, rateLimit)
}

// Factory uses pre-calculated value
// coverage:ignore
func NewNotificationServiceForProduction(cfg Config) *NotificationService {
    mailer := email.NewSMTPMailer(cfg.SMTPConfig)
    rateLimit := calculateRateLimit(cfg.Timeout)
    return NewNotificationService(mailer, rateLimit)
}
```

---

#### 5. Loop with Business Logic: ReportService

**Location:** `sample-violations.go:96-113`

**Issue:** Production factory contains loop with switch statement to select formatter.

**Violated Standard:**
> Production factories MUST NOT contain business logic (conditionals, loops, calculations, validations). (tech_standards.md § IoC Patterns - Business Logic Prohibition)

**Current Code:**
```go
// coverage:ignore
func NewReportServiceForProduction(cfg Config) *ReportService {
    var formatters []Formatter

    // ❌ Loop + switch = business logic in factory!
    for _, format := range cfg.SupportedFormats {
        switch format {
        case "pdf":
            formatters = append(formatters, pdf.NewFormatter())
        case "excel":
            formatters = append(formatters, excel.NewFormatter())
        case "csv":
            formatters = append(formatters, csv.NewFormatter())
        }
    }

    return NewReportService(formatters)
}
```

**Impact:**
- Format selection logic is untestable
- Cannot verify correct formatters are chosen for configurations
- Complex logic excluded from coverage

**Required Fix:**

Extract to configuration builder:

```go
// Configuration builder - testable
func buildFormatters(supportedFormats []string) []Formatter {
    var formatters []Formatter
    for _, format := range supportedFormats {
        switch format {
        case "pdf":
            formatters = append(formatters, pdf.NewFormatter())
        case "excel":
            formatters = append(formatters, excel.NewFormatter())
        case "csv":
            formatters = append(formatters, csv.NewFormatter())
        }
    }
    return formatters
}

func TestBuildFormatters(t *testing.T) {
    formats := []string{"pdf", "csv"}
    formatters := buildFormatters(formats)

    assert.Len(t, formatters, 2)
    assert.IsType(t, &pdf.Formatter{}, formatters[0])
    assert.IsType(t, &csv.Formatter{}, formatters[1])
}

// Factory uses pre-built formatters
// coverage:ignore
func NewReportServiceForProduction(cfg Config) *ReportService {
    formatters := buildFormatters(cfg.SupportedFormats)
    return NewReportService(formatters)
}
```

---

#### 6. External API Call in Factory: WeatherService

**Location:** `sample-violations.go:123-133`

**Issue:** Production factory makes external API call during initialization.

**Violated Standard:**
> Production factories MUST NOT contain business logic (conditionals, loops, calculations, validations). (tech_standards.md § IoC Patterns - Business Logic Prohibition)
> Production factories should not perform I/O operations. (architecture.md § Service Initialization)

**Current Code:**
```go
// coverage:ignore
func NewWeatherServiceForProduction(apiKey string) *WeatherService {
    client := weather.NewAPIClient(apiKey)

    // ❌ External API call in factory!
    config, _ := client.FetchConfig()

    return NewWeatherService(client, config)
}
```

**Impact:**
- Factory depends on external API availability at startup
- Cannot test service initialization without real API
- Failures during initialization are difficult to handle
- Network issues cause factory to fail

**Required Fix:**

Lazy load configuration or separate initialization:

```go
// Factory creates service without API call
// coverage:ignore
func NewWeatherServiceForProduction(apiKey string) *WeatherService {
    client := weather.NewAPIClient(apiKey)
    // Don't fetch config here - let service fetch when needed
    return NewWeatherService(client, nil)
}

// Service handles initialization
func (s *WeatherService) Initialize(ctx context.Context) error {
    if s.config != nil {
        return nil // Already initialized
    }

    config, err := s.client.FetchConfig(ctx)
    if err != nil {
        return fmt.Errorf("fetching config: %w", err)
    }

    s.config = config
    return nil
}

// Testable
func TestWeatherService_Initialize(t *testing.T) {
    mockClient := mocks.NewAPIClient(t)
    mockConfig := &weather.Config{Units: "metric"}

    mockClient.EXPECT().FetchConfig(mock.Anything).Return(mockConfig, nil)

    service := NewWeatherService(mockClient, nil)
    err := service.Initialize(context.Background())

    assert.NoError(t, err)
    assert.Equal(t, mockConfig, service.config)
}
```

---

#### 7. Data Transformation in Factory: DataService

**Location:** `sample-violations.go:143-151`

**Issue:** Production factory performs data transformation (string processing).

**Violated Standard:**
> Production factories MUST NOT contain business logic (conditionals, loops, calculations, validations). (tech_standards.md § IoC Patterns - Business Logic Prohibition)

**Current Code:**
```go
// coverage:ignore
func NewDataServiceForProduction(rawConfig map[string]string) *DataService {
    // ❌ Data transformation in factory!
    config := Config{}
    for k, v := range rawConfig {
        config[strings.ToUpper(k)] = strings.TrimSpace(v)
    }

    return NewDataService(config)
}
```

**Impact:**
- Transformation logic (uppercase keys, trim values) is untestable
- Cannot verify transformation correctness without real factory
- Logic errors won't be caught by unit tests

**Required Fix:**

Extract transformation function:

```go
// Configuration transformation - testable
func normalizeConfig(rawConfig map[string]string) Config {
    config := Config{}
    for k, v := range rawConfig {
        config[strings.ToUpper(k)] = strings.TrimSpace(v)
    }
    return config
}

func TestNormalizeConfig(t *testing.T) {
    raw := map[string]string{
        "api_key": "  sk_test  ",
        "timeout": " 30 ",
    }

    config := normalizeConfig(raw)

    assert.Equal(t, "sk_test", config["API_KEY"])
    assert.Equal(t, "30", config["TIMEOUT"])
}

// Factory uses pre-normalized config
// coverage:ignore
func NewDataServiceForProduction(rawConfig map[string]string) *DataService {
    config := normalizeConfig(rawConfig)
    return NewDataService(config)
}
```

---

#### 8-14. Missing Coverage Exclusion Markers

**Locations:** All 7 production factories

**Issue:** None of the production factories have `// coverage:ignore` markers.

**Violated Standard:**
> All production factories MUST be marked with // coverage:ignore. (tech_standards.md § Coverage Strategy)

**Impact:**
- Coverage reports are inaccurate
- Factory code counts toward coverage percentage
- May create false sense of adequate testing

**Required Fix:**

Add `// coverage:ignore` comment immediately before each factory function:

```go
// coverage:ignore
func NewOrderServiceForProduction(db *sql.DB) *OrderService {
    // ...
}
```

**Affected functions:**
- NewOrderServiceForProduction (line 23)
- NewUserServiceForProduction (line 40)
- NewPaymentServiceForProduction (line 58)
- NewNotificationServiceForProduction (line 78)
- NewReportServiceForProduction (line 96)
- NewWeatherServiceForProduction (line 123)
- NewDataServiceForProduction (line 143)

---

#### 15. Test Using Production Factory

**Location:** `sample-violations.go:162-169`

**Issue:** Unit test uses production factory instead of primary constructor with mocks.

**Violated Standard:**
> Tests MUST use primary constructors with mocks. Production factories are for production only. (testing.md § Unit Testing Patterns)

**Current Code:**
```go
func TestOrderService_Create(t *testing.T) {
    // ❌ Using production factory in test!
    service := NewOrderServiceForProduction(testDB)

    // Test requires real database instead of mocks
    err := service.Create(order)
    assert.NoError(t, err)
}
```

**Impact:**
- Test requires real database setup
- Slow test execution
- Fragile tests (depends on DB state)
- Cannot test error scenarios easily
- Violates unit testing principles

**Required Fix:**

Use primary constructor with mocks:

```go
func TestOrderService_Create(t *testing.T) {
    // ✅ Using primary constructor with mocks
    mockRepo := mocks.NewOrderRepository(t)
    mockLogger := mocks.NewLogger(t)

    service := NewOrderService(mockRepo, mockLogger)

    order := &Order{ID: "123", Total: 100.00}
    mockRepo.EXPECT().Save(order).Return(nil)
    mockLogger.EXPECT().Info("Order created", "id", "123")

    err := service.Create(order)

    assert.NoError(t, err)
}
```

**Benefits:**
- No database required
- Fast execution (microseconds vs milliseconds)
- Can easily test error scenarios
- Follows IoC pattern correctly

---

#### 16. Missing Primary Constructor (Implicit)

**Location:** `sample-violations.go:180`

**Issue:** Test file line 180 attempts to call `NewOrderService()` but function doesn't exist.

**Violated Standard:**
> Every service MUST have a primary constructor. (tech_standards.md § IoC Patterns)

**Current Impact:**
- Code doesn't compile
- Cannot write proper unit tests
- Violates the IoC pattern

**Required Fix:**

See fix for Violation #1 (add primary constructor to OrderService).

---

### Warnings (Should Fix)

#### W1-W8. Missing Godoc Comments

**Locations:**
- OrderService struct (line 15)
- UserService.CheckUserCount method (line 50)
- PaymentService.ProcessPayment method (line 70)
- NotificationService.Send method (line 88)
- ReportService.Generate method (line 115)
- WeatherService.GetForecast method (line 135)
- DataService.Query method (line 153)
- TestOrderService_Create function (line 162)

**Issue:** Exported types and functions lack documentation comments.

**Violated Standard:**
> All exported types, functions, and methods should have doc comments. (tech_standards.md § Go Conventions)

**Impact:**
- Poor developer experience
- Harder to understand API
- godoc output is incomplete

**Required Fix Pattern:**

```go
// OrderService handles order processing and validation.
// It coordinates with the repository for persistence and
// uses the logger for audit trail.
type OrderService struct {
    repo   OrderRepository
    logger Logger
}

// Create validates and saves a new order to the repository.
// Returns an error if validation fails or save operation fails.
func (s *OrderService) Create(order *Order) error {
    // implementation
}
```

---

## Compliance Checklist by Service

### OrderService
- [ ] Add primary constructor
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to struct
- [ ] Fix test to use primary constructor

### UserService
- [ ] Extract business logic from factory
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to CheckUserCount method
- [ ] Add tests for extracted business logic

### PaymentService
- [ ] Extract configuration decision logic
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to ProcessPayment method
- [ ] Add tests for configuration logic

### NotificationService
- [ ] Extract calculation logic
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to Send method
- [ ] Add tests for rate limit calculation

### ReportService
- [ ] Extract formatter selection logic
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to Generate method
- [ ] Add tests for formatter selection

### WeatherService
- [ ] Remove API call from factory, use lazy initialization
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to GetForecast method
- [ ] Add tests for initialization

### DataService
- [ ] Extract transformation logic
- [ ] Add `// coverage:ignore` to production factory
- [ ] Add godoc comment to Query method
- [ ] Add tests for config normalization

---

## Summary Statistics

**By Severity:**
- Critical: 16 violations (must fix before merge)
- Warning: 8 violations (should fix)

**By Category:**
- IoC Pattern Violations: 16 (67% of issues)
  - Missing primary constructors: 2
  - Business logic in factories: 6
  - Missing coverage markers: 7
  - Test anti-pattern: 1
- Documentation: 8 (33% of issues)

**Estimated Effort:**
- Critical fixes: 4-6 hours (refactoring + tests)
- Documentation: 1 hour
- Total: 5-7 hours

---

## Next Steps

1. **Immediate:** Fix critical IoC violations (violations #1-16)
2. **Before merge:** Add godoc comments (warnings W1-W8)
3. **Testing:** Ensure 80%+ coverage after refactoring
4. **Review:** Request re-review after fixes applied

---

## Execution Metadata

**Agent:** standards-compliance v1.0
**Standards Version:** tech_standards.md v2.1.0, architecture.md v1.3.0
**Execution Time:** 2.3 seconds
**Files Analyzed:** 1 (sample-violations.go, 180 lines)
**Violations Detected:** 24 total (16 critical, 8 warning)
**Confidence:** HIGH - All violations have clear fixes with code examples
