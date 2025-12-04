package services

import (
	"context"
	"database/sql"
)

// ✅ CORRECT: UserService with proper IoC pattern
type UserService struct {
	repo      UserRepository
	logger    Logger
	validator Validator
}

// ✅ CORRECT: Primary constructor taking ALL dependencies
// Used in tests with mocks
func NewUserService(
	repo UserRepository,
	logger Logger,
	validator Validator,
) *UserService {
	return &UserService{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}
}

// ✅ CORRECT: Production factory
// - Takes only shared dependencies (db, logger)
// - Builds non-shared dependencies internally
// - Contains zero business logic
// - Calls primary constructor
// - Marked for coverage exclusion
// coverage:ignore
func NewUserServiceForProduction(db *sql.DB, logger Logger) *UserService {
	// Only dependency wiring - no business logic
	validator := validation.NewUserValidator()
	repo := persistence.NewUserRepository(db)

	return NewUserService(repo, logger, validator)
}

// ✅ CORRECT: Business logic in service method, not factory
func (s *UserService) CheckCapacity(ctx context.Context) error {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return fmt.Errorf("checking capacity: %w", err)
	}

	// Business logic belongs in service methods
	if count > 1000 {
		s.logger.Warn("High user count detected", "count", count)
	}

	return nil
}

// ✅ CORRECT: CreateUser implements business logic
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*User, error) {
	if err := s.validator.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := s.repo.Create(ctx, email, name)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	s.logger.Info("User created", "userID", user.ID, "email", email)
	return user, nil
}

// ✅ CORRECT: OrderService with proper pattern
type OrderService struct {
	repo       OrderRepository
	logger     Logger
	calculator PriceCalculator
}

// ✅ CORRECT: Primary constructor
func NewOrderService(
	repo OrderRepository,
	logger Logger,
	calculator PriceCalculator,
) *OrderService {
	return &OrderService{
		repo:       repo,
		logger:     logger,
		calculator: calculator,
	}
}

// ✅ CORRECT: Production factory with simple wiring
// coverage:ignore
func NewOrderServiceForProduction(db *sql.DB, logger Logger) *OrderService {
	repo := persistence.NewOrderRepository(db)
	calculator := pricing.NewCalculator()

	return NewOrderService(repo, logger, calculator)
}

// ✅ CORRECT: PaymentService - configuration passed as value, not logic in factory
type PaymentService struct {
	processor PaymentProcessor
	repo      PaymentRepository
	logger    Logger
}

func NewPaymentService(
	processor PaymentProcessor,
	repo PaymentRepository,
	logger Logger,
) *PaymentService {
	return &PaymentService{
		processor: processor,
		repo:      repo,
		logger:    logger,
	}
}

// ✅ CORRECT: No business decisions - builds what's configured
// Configuration decisions made elsewhere, factory just wires
// coverage:ignore
func NewPaymentServiceForProduction(db *sql.DB, logger Logger) *PaymentService {
	repo := persistence.NewPaymentRepository(db)
	processor := processors.NewProcessor() // Simple instantiation

	return NewPaymentService(processor, repo, logger)
}

// ✅ CORRECT: NotificationService - timeout passed as parameter
type NotificationService struct {
	sender  EmailSender
	logger  Logger
	timeout int
}

func NewNotificationService(
	sender EmailSender,
	logger Logger,
	timeout int,
) *NotificationService {
	return &NotificationService{
		sender:  sender,
		logger:  logger,
		timeout: timeout,
	}
}

// ✅ CORRECT: Timeout calculated in config layer, passed here
// coverage:ignore
func NewNotificationServiceForProduction(logger Logger, cfg Config) *NotificationService {
	sender := email.NewSMTPSender(cfg.SMTPHost)

	// ✅ CORRECT: Using pre-computed value from config, not calculating here
	return NewNotificationService(sender, logger, cfg.NotificationTimeout)
}

// ✅ CORRECT: ReportService - generators built by factory function
type ReportService struct {
	generators []ReportGenerator
	logger     Logger
}

func NewReportService(generators []ReportGenerator, logger Logger) *ReportService {
	return &ReportService{generators: generators, logger: logger}
}

// ✅ CORRECT: Calls separate factory function to build generators
// coverage:ignore
func NewReportServiceForProduction(logger Logger, cfg Config) *ReportService {
	// ✅ Delegates complex building to helper function
	generators := buildGenerators(cfg.EnabledReports)

	return NewReportService(generators, logger)
}

// ✅ CORRECT: Helper function for complex building (also excluded from coverage)
// coverage:ignore
func buildGenerators(reportTypes []string) []ReportGenerator {
	var generators []ReportGenerator

	for _, reportType := range reportTypes {
		switch reportType {
		case "sales":
			generators = append(generators, reporting.NewSalesGenerator())
		case "inventory":
			generators = append(generators, reporting.NewInventoryGenerator())
		}
	}

	return generators
}

// ✅ CORRECT: Test using primary constructor with mocks
func TestUserService_CreateUser(t *testing.T) {
	// ✅ Arrange - inject mocks via primary constructor
	mockRepo := mocks.NewUserRepository(t)
	mockLogger := mocks.NewLogger(t)
	mockValidator := mocks.NewValidator(t)

	service := NewUserService(mockRepo, mockLogger, mockValidator)

	mockValidator.EXPECT().ValidateEmail("test@example.com").Return(nil)
	mockRepo.EXPECT().
		Create(mock.Anything, "test@example.com", "Test User").
		Return(&User{ID: "user-123", Email: "test@example.com"}, nil)

	// ✅ Act
	user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")

	// ✅ Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)

	mockRepo.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
}

// ✅ CORRECT: Table-driven test for multiple scenarios
func TestUserService_ValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email format",
			email:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewUserRepository(t)
			mockLogger := mocks.NewLogger(t)
			mockValidator := mocks.NewValidator(t)

			service := NewUserService(mockRepo, mockLogger, mockValidator)

			if tt.wantErr {
				mockValidator.EXPECT().ValidateEmail(tt.email).Return(ErrInvalidEmail)
			} else {
				mockValidator.EXPECT().ValidateEmail(tt.email).Return(nil)
				mockRepo.EXPECT().
					Create(mock.Anything, tt.email, mock.Anything).
					Return(&User{}, nil)
			}

			_, err := service.CreateUser(context.Background(), tt.email, "Test")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ✅ CORRECT: Exported function with godoc comment
// ProcessPayment handles payment processing for an order
func ProcessPayment(ctx context.Context, amount float64) error {
	// Implementation
	return nil
}

// ✅ CORRECT: Domain errors defined as package variables
var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email already exists")
)

// ✅ CORRECT: Container manages shared dependencies only
type Container struct {
	db     *sql.DB
	logger Logger
	cfg    Config

	userService  *UserService
	orderService *OrderService
}

// ✅ CORRECT: Container uses production factories
// coverage:ignore
func NewContainer(cfg Config) (*Container, error) {
	c := &Container{cfg: cfg}

	var err error
	c.db, err = initDatabase(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	c.logger = initLogger(cfg.LogLevel)

	// ✅ Wire services using production factories
	c.userService = NewUserServiceForProduction(c.db, c.logger)
	c.orderService = NewOrderServiceForProduction(c.db, c.logger)

	return c, nil
}

func (c *Container) UserService() *UserService {
	return c.userService
}

func (c *Container) OrderService() *OrderService {
	return c.orderService
}

func (c *Container) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// ✅ CORRECT: Test container for integration tests
// coverage:ignore
func NewTestContainer() (*Container, error) {
	cfg := Config{
		DatabaseURL: "postgres://test:test@localhost:5432/test_db",
		LogLevel:    "debug",
	}
	return NewContainer(cfg)
}
