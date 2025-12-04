# Sample Approved Spec for Implementation

@pending @story-PROJ-1234
Feature: Password Reset
  As a registered user
  I want to reset my password
  So that I can regain access to my account

  @happy-path @security
  Scenario: User requests password reset successfully
    Given a user exists with email "user@example.com"
    When I request a password reset for "user@example.com"
    Then I should receive a reset email
    And the reset link should expire in 24 hours
    And the password reset attempt should be logged

  @validation @security
  Scenario: System returns generic message regardless of email existence
    When I request a password reset for "nonexistent@example.com"
    Then I should see "If that email is registered, you will receive a password reset link"
    And the password reset attempt should be logged
    And no email should be sent

  @happy-path
  Scenario: User completes password reset with valid token
    Given a user exists with email "user@example.com"
    And I requested a password reset 1 hour ago
    When I set a new password to "NewPass123!" using the reset token
    Then I should be logged in
    And my old password should be invalidated
    And I should receive a confirmation email
    And the password reset attempt should be logged
