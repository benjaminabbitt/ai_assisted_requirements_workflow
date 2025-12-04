# Sample PR for BO Review

**PR #42:** [SPEC] PROJ-1234: Password Reset Feature
**Branch:** `spec/PROJ-1234`
**Author:** AI Requirements Analyst
**Labels:** `specification`, `ai-generated`, `needs-bo-approval`
**Status:** Awaiting Business Owner approval

---

## Specification for PROJ-1234

**Ticket:** https://jira.example.com/browse/PROJ-1234

**Confidence:** MEDIUM

### Context Versions Used
- business.md: v1.0
- architecture.md: NOT AVAILABLE
- testing.md: v1.0
- tech_standards.md: NOT AVAILABLE

### Scenarios Added
- User requests password reset successfully
- System validates email exists before sending reset
- Reset email contains valid token
- User completes password reset with valid token
- Old password is invalidated after reset
- Reset token expires after 24 hours
- Single-use token prevents reuse
- Rate limiting blocks excessive reset requests
- Email validation boundary conditions
- Password complexity validation

### Assumptions (MEDIUM Confidence)
- SendGrid is configured via environment variables - **verify with BO**
- Internal auth service has endpoints for password reset flow - **verify with BO**
- Token generation uses cryptographically secure random - **verify with BO**
- Email delivery is asynchronous - **verify with BO**

---

## Files Changed

### `features/auth/password_reset.feature`

```gherkin
@story-PROJ-1234
Feature: Password Reset
  As a registered user
  I want to reset my password
  So that I can regain access to my account if I forget my password

  Background:
    Given the system is initialized
    And the user database is clean

  @happy-path @security
  Scenario: User requests password reset successfully
    Given a user exists with email "user@example.com"
    When I request a password reset for "user@example.com"
    Then I should receive a reset email
    And the reset link should expire in 24 hours
    And the password reset attempt should be logged

  @validation
  Scenario: System validates email exists before sending reset
    When I request a password reset for "nonexistent@example.com"
    Then I should see "Email not found"
    And the password reset attempt should be logged

  @happy-path
  Scenario: User completes password reset with valid token
    Given a user exists with email "user@example.com"
    And I requested a password reset 1 hour ago
    When I set a new password to "NewPass123!" using the reset token
    Then I should be logged in
    And my old password should be invalidated
    And I should receive a confirmation email
    And the password reset attempt should be logged

  @security @edge-case
  Scenario: Reset token expires after 24 hours
    Given a user exists with email "user@example.com"
    And I requested a password reset 24 hours and 1 minute ago
    When I attempt to use the reset token
    Then I should see "Reset link expired"
    And the password reset attempt should be logged

  @security @edge-case
  Scenario: Single-use token prevents reuse
    Given a user exists with email "user@example.com"
    And I requested a password reset 1 hour ago
    And I set a new password to "NewPass123!" using the reset token
    When I attempt to use the reset token again
    Then I should see "Reset link already used"
    And the password reset attempt should be logged

  @security @rate-limiting
  Scenario: Rate limiting blocks excessive reset requests
    Given a user exists with email "user@example.com"
    When I request a password reset 5 times within 15 minutes for "user@example.com"
    Then the 5th request should succeed
    When I request a password reset for "user@example.com"
    Then I should see "Too many reset attempts. Please try again later."
    And the password reset attempt should be logged

  @boundary @validation
  Scenario Outline: Email validation boundary conditions
    When I request a password reset for "<email>"
    Then I should see "<error>"
    And the password reset attempt should be logged

    Examples:
      | email                        | error                |
      |                              | Email required       |
      | invalid                      | Invalid email format |
      | test@                        | Invalid email format |
      | @example.com                 | Invalid email format |

  @boundary @validation
  Scenario Outline: Password complexity validation
    Given a user exists with email "user@example.com"
    And I requested a password reset 1 hour ago
    When I set a new password to "<password>" using the reset token
    Then I should see "<error>"

    Examples:
      | password    | error                                          |
      | short       | Password must be at least 8 characters         |
      | nouppercase | Password must contain at least 1 uppercase     |
      | NOLOWERCASE | Password must contain at least 1 lowercase     |
      | NoNumber!   | Password must contain at least 1 number        |
      | NoSpecial1  | Password must contain at least 1 special char  |
```

### `features/step_definitions/password_reset_steps.go`

*[Skeleton step definitions - implementation by developer]*

---

## Business Owner Review Questions

1. **Assumptions:** Are the 4 assumptions listed above correct?
   - SendGrid configuration
   - Auth service endpoints
   - Token generation method
   - Email delivery async

2. **Business Rules:** Do the scenarios correctly implement:
   - BR-002: 24-hour token expiry
   - BR-003: Password complexity requirements
   - BR-004: Rate limiting (5 requests/15 min)
   - BR-001: Audit logging

3. **Error Messages:** Are the error messages appropriate?
   - "Email not found"
   - "Reset link expired"
   - "Reset link already used"
   - "Too many reset attempts. Please try again later."
   - Password complexity error messages

4. **Completeness:** Any missing scenarios from a business perspective?

---

## Developer Notes

Implementation will follow after BO approval. Skeleton step definitions provided. Will use primary constructor pattern per tech_standards.md (when available).

---

**Ready for Business Owner review** 🔍
