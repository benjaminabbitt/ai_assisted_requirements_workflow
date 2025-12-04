# Testing Context (Sample)

## Step Library

### Authentication Steps

| Step | Parameters | Notes |
|------|------------|-------|
| `Given a user exists with email "([^"]*)"` | email | Creates test user |
| `Given I am authenticated as "([^"]*)"` | email | Sets up auth context |
| `When I log in with email "([^"]*)" and password "([^"]*)"` | email, password | Full login flow |
| `Then I should be logged in` | none | Verifies auth state |

### Email Steps

| Step | Parameters | Notes |
|------|------------|-------|
| `Then I should receive an email at "([^"]*)"` | email | Checks test mailbox |
| `And the email should contain "([^"]*)"` | text | Validates email content |

## Boundary Condition Patterns

| Data Type | Boundary Conditions | Example Values |
|-----------|--------------------| ---------------|
| Email | valid, invalid, empty, max length | `"user@example.com"`, `"invalid"`, `""`, `"a"*255+"@example.com"` |
| String | empty, whitespace, max length | `""`, `"   "`, `"a"*255` |
| Password | min length, max length, complexity | `"short"`, `"Pass123!"`, `"a"*100` |

## Edge Case Patterns

| Pattern | When to Apply | Example |
|---------|---------------|---------|
| Token expiry | Time-based features | Test at expiry boundary (24h, 24h+1min) |
| Rate limiting | High-frequency actions | Test at limit, over limit |
| Concurrent access | State-changing operations | Two simultaneous requests |
