package services

import (
	"context"
	"database/sql"
	"fmt"
)

// ❌ VIOLATION 2: Missing primary constructor
// OrderService only has production factory, no primary constructor for testing
type OrderService struct {
	repo       OrderRepository
	logger     Logger
	calculator PriceCalculator
}

// ❌ VIOLATION 2: This should call a primary constructor
func NewOrderServiceForProduction(db *sql.DB, logger Logger) *OrderService {
	repo := persistence.NewOrderRepository(db)
	calculator := pricing.NewCalculator()

	// Direct instantiation instead of calling primary constructor
	return &OrderService{
		repo:       repo,
		logger:     logger,
		calculator: calculator,
	}
}

// UserService has correct structure but violation in production factory
type UserService struct {
	repo      UserRepository
	logger    Logger
	validator Validator
}

// ✅ CORRECT: Primary constructor exists
func NewUserService(repo UserRepository, logger Logger, validator Validator) *UserService {
	return &UserService{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}
}

// ❌ VIOLATION 1: Business logic in production factory
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
	validator := validation.NewUserValidator()
	repo := persistence.NewUserRepository(db)

	// ❌ VIOLATION: Checking user count is business logic
	count, _ := repo.Count(context.Background())
	if count > 1000 {
		logger.Warn("High user count detected", "count", count)
	}

	// ❌ VIOLATION: Also missing coverage exclusion marker
	return NewUserService(repo, logger, validator)
}

// PaymentService demonstrates configuration decision violation
type PaymentService struct {
	processor PaymentProcessor
	repo      PaymentRepository
	logger    Logger
}

func NewPaymentService(processor PaymentProcessor, repo PaymentRepository, logger Logger) *PaymentService {
	return &PaymentService{
		processor: processor,
		repo:      repo,
		logger:    logger,
	}
}

// ❌ VIOLATION 3: Configuration decision (business rule) in production factory
func NewPaymentServiceForProduction(db *sql.DB, logger Logger, cfg Config) *PaymentService {
	repo := persistence.NewPaymentRepository(db)

	// ❌ VIOLATION: Choosing processor based on business requirement
	var processor PaymentProcessor
	if cfg.StrictMode {
		processor = processors.NewStrictProcessor(cfg.Timeout)
	} else {
		processor = processors.NewFastProcessor()
	}

	return NewPaymentService(processor, repo, logger)
}

// NotificationService demonstrates calculation violation
type NotificationService struct {
	sender  EmailSender
	logger  Logger
	timeout int
}

func NewNotificationService(sender EmailSender, logger Logger, timeout int) *NotificationService {
	return &NotificationService{
		sender:  sender,
		logger:  logger,
		timeout: timeout,
	}
}

// ❌ VIOLATION 4: Calculation in production factory
func NewNotificationServiceForProduction(logger Logger, cfg Config) *NotificationService {
	sender := email.NewSMTPSender(cfg.SMTPHost)

	// ❌ VIOLATION: Calculating timeout based on environment
	timeout := 30
	if cfg.Environment == "production" {
		timeout = timeout * 2 // Calculation is business logic
	}

	return NewNotificationService(sender, logger, timeout)
}

// ❌ VIOLATION: Test using production factory instead of primary constructor
func TestUserService_CreateUser(t *testing.T) {
	db := setupTestDB()
	logger := setupTestLogger()

	// ❌ VIOLATION: Should use primary constructor with mocks
	service := NewUserServiceForProduction(db, logger)

	user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")

	assert.NoError(t, err)
	assert.NotNil(t, user)
}

// ❌ VIOLATION: Missing godoc comment on exported function
func ProcessPayment(ctx context.Context, amount float64) error {
	// Implementation
	return nil
}

// ❌ VIOLATION: Loop in production factory
type ReportService struct {
	generators []ReportGenerator
	logger     Logger
}

func NewReportService(generators []ReportGenerator, logger Logger) *ReportService {
	return &ReportService{generators: generators, logger: logger}
}

func NewReportServiceForProduction(logger Logger, cfg Config) *ReportService {
	var generators []ReportGenerator

	// ❌ VIOLATION: Loop with conditional logic in factory
	for _, reportType := range cfg.EnabledReports {
		switch reportType {
		case "sales":
			generators = append(generators, reporting.NewSalesGenerator())
		case "inventory":
			generators = append(generators, reporting.NewInventoryGenerator())
		}
	}

	return NewReportService(generators, logger)
}

// ❌ VIOLATION: External API call in production factory
type WeatherService struct {
	client APIClient
	cache  Cache
	logger Logger
}

func NewWeatherService(client APIClient, cache Cache, logger Logger) *WeatherService {
	return &WeatherService{client: client, cache: cache, logger: logger}
}

func NewWeatherServiceForProduction(logger Logger) *WeatherService {
	client := api.NewClient("https://api.weather.com")
	cache := cache.NewRedisCache()

	// ❌ VIOLATION: Making API call to check service health
	if !client.HealthCheck() {
		logger.Error("Weather API unavailable")
	}

	return NewWeatherService(client, cache, logger)
}

// ❌ VIOLATION: Data transformation in production factory
type DataService struct {
	transformer DataTransformer
	repo        DataRepository
}

func NewDataService(transformer DataTransformer, repo DataRepository) *DataService {
	return &DataService{transformer: transformer, repo: repo}
}

func NewDataServiceForProduction(db *sql.DB, format string) *DataService {
	repo := persistence.NewDataRepository(db)

	// ❌ VIOLATION: Transforming data based on format
	transformer := transformers.NewTransformer()
	if format == "json" {
		transformer.SetFormat("application/json")
	} else if format == "xml" {
		transformer.SetFormat("application/xml")
	}

	return NewDataService(transformer, repo)
}
