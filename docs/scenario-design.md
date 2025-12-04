# Scenario Design Guidance

This document provides rules for well-structured Gherkin specifications that are maintainable, readable, and effective.

---

## Overview

Well-designed scenarios:
- Capture business requirements in readable format
- Are maintainable as the system evolves
- Reuse steps to reduce duplication
- Focus on behavior, not implementation details

---

## Feature File Organization

### One Capability Per Feature

Each `.feature` file should cover a single user capability. If you find yourself using "and" to describe the feature, split it.

- **Good:** `auth_login.feature`, `auth_password_reset.feature`
- **Bad:** `auth_login_and_password_reset.feature`

### Size Limits

| Metric | Limit | Rationale |
|--------|-------|-----------|
| Scenarios per feature | 15 max | Reduces merge conflicts, improves readability |
| Steps per scenario | 10 max | Keeps scenarios focused |
| Examples per outline | 10 max | Prevents data-driven bloat |

**Why limits matter:**
- Smaller feature files = fewer merge conflicts when multiple stories touch same domain
- Focused scenarios = easier to understand what's being tested
- Limited examples = data tables stay readable

---

## When to Use Scenario Outline

Use Scenario Outline when:

- Same behavior with different data (validation rules, boundary conditions)
- Data combinations matter (role × action permission matrix)
- Error messages vary by input

Do NOT use Scenario Outline when:

- Workflow differs between cases (use separate scenarios)
- Only 2-3 examples exist (inline is clearer)
- Examples require different setup (Background won't work)

**Example - Good use of Scenario Outline:**
```gherkin
Scenario Outline: Email validation
  Given I am on the registration page
  When I enter email "<email>"
  Then I should see "<result>"

  Examples:
    | email | result |
    | valid@example.com | Accepted |
    | invalid | Invalid email format |
    | | Email required |
    | a@b.co | Accepted |
```

**Example - Bad use of Scenario Outline:**
```gherkin
# DON'T DO THIS - workflows are different
Scenario Outline: User actions
  Given I am logged in as "<role>"
  When I perform "<action>"
  Then I should see "<result>"

  Examples:
    | role | action | result |
    | admin | delete user | User deleted |
    | user | view profile | Profile shown |
```

**Why bad:** Admin deleting a user is a completely different workflow from viewing a profile. These should be separate scenarios.

---

## Background Usage

Use Background for:

- Common preconditions shared by ALL scenarios in the feature
- User authentication when all scenarios need a logged-in user
- Data setup that doesn't vary between scenarios

Do NOT use Background when:

- Some scenarios need different setup
- Background would exceed 5 steps
- Setup is complex enough to obscure scenario intent

**Example - Good Background:**
```gherkin
Background:
  Given a user exists with email "user@example.com"
  And I am logged in as "user@example.com"

Scenario: View profile
  When I navigate to my profile
  Then I should see my email "user@example.com"

Scenario: Update profile
  When I update my name to "John Doe"
  Then my profile should show "John Doe"
```

**Example - Bad Background:**
```gherkin
# DON'T DO THIS - not all scenarios need admin
Background:
  Given I am logged in as admin
  And there are 100 users in the system
  And I navigate to the users page
  And I sort by email ascending

Scenario: Delete user
  # ... only this scenario needs admin
```

---

## Step Writing Rules

### Given-When-Then Discipline

| Keyword | Purpose | Verb Tense |
|---------|---------|------------|
| Given | Establish preconditions | Past/present state |
| When | Action under test | Present tense |
| Then | Assert outcomes | Present/future state |
| And/But | Continue previous keyword | Match previous |

**Examples:**
- **Given** a user exists with email "user@example.com" _(precondition)_
- **When** I request a password reset for "user@example.com" _(action)_
- **Then** I should receive a reset email _(assertion)_
- **And** the password reset attempt should be logged _(additional assertion)_

### Step Reuse Priority

1. **First:** Reuse existing step from `testing.md` exactly as written
2. **Second:** Parameterize existing step if minor variation needed
3. **Third:** Create new step only if no existing step fits

New steps must be documented in `testing.md` before PR merge.

**Example - Step reuse:**
```gherkin
# Existing step in testing.md:
# Given a user exists with email {string}

Scenario: Password reset
  Given a user exists with email "user@example.com"  # Reused step
  When I request a password reset for "user@example.com"
  Then I should receive a reset email
```

### Avoid Implementation Details

Focus on business behavior, not technical implementation.

**Bad** _(implementation details)_:
```gherkin
Given the database contains a user record with id=123
When I execute a SELECT query for id=123
Then the response should contain JSON with email field
```

**Good** _(business behavior)_:
```gherkin
Given a user exists with email "user@example.com"
When I request user details for "user@example.com"
Then I should see the user's email
```

### Use Domain Language

Use terms from `business.md`, not technical jargon.

**Bad:**
```gherkin
When I POST to /api/v1/auth/reset with payload {"email":"user@example.com"}
```

**Good:**
```gherkin
When I request a password reset for "user@example.com"
```

---

## Scenario Examples

### Well-Written Scenario

```gherkin
@story-PROJ-1234
Scenario: User requests password reset successfully
  Given a user exists with email "user@example.com"
  When I request a password reset for "user@example.com"
  Then I should receive a reset email
  And the password reset attempt should be logged
  And the reset token should expire in 24 hours
```

**Why good:**
- Clear business behavior (no implementation details)
- Uses domain language
- Focused (5 steps, single behavior)
- Tagged with story ID
- Reuses existing steps

### Common Mistakes

#### Mistake 1: Too Many Steps

```gherkin
# DON'T DO THIS - too many steps
Scenario: Complex workflow
  Given a user exists
  And I am logged in
  And I navigate to settings
  And I click on profile tab
  And I scroll to password section
  And I click change password
  When I enter old password "old123"
  And I enter new password "new456"
  And I confirm new password "new456"
  And I click save
  Then I should see success message
  And password should be updated
  And old password should not work
  And new password should work
```

**Fix:** Break into multiple focused scenarios or extract setup into Background.

#### Mistake 2: Implementation Leaking

```gherkin
# DON'T DO THIS - too technical
Scenario: Database update
  Given the users table has a record with id=123
  When I execute UPDATE users SET password_hash = 'xyz' WHERE id=123
  Then the database should contain the new hash
```

**Fix:** Use business language.

#### Mistake 3: Multiple Behaviors

```gherkin
# DON'T DO THIS - testing multiple things
Scenario: User management
  Given I am an admin
  When I create a new user
  Then the user should exist
  When I delete the user
  Then the user should not exist
```

**Fix:** Split into two scenarios (create, delete).

---

## Boundary Conditions & Edge Cases

AI generates boundary conditions based on patterns in `testing.md`. Ensure your testing.md includes patterns for:

- **Strings:** empty, whitespace, max length, unicode, special chars
- **Numbers:** zero, negative, min, max, overflow
- **Dates:** epoch, far future, leap year, timezone edges
- **Arrays:** empty, single item, max items, duplicates
- **Emails:** valid, invalid, max length, special domains

**Example - Boundary scenarios:**
```gherkin
Scenario Outline: Email validation boundary conditions
  When I register with email "<email>"
  Then I should see "<result>"

  Examples:
    | email | result |
    | valid@example.com | Success |
    | | Email required |
    | invalid | Invalid format |
    | a@b.co | Success |
    | very-long-email-address-that-exceeds-max-length@example.com | Email too long |
    | user+tag@example.com | Success |
```

---

## Tags

### Required Tags

- `@story-{id}` - Links to originating ticket (REQUIRED on all scenarios)

### Optional Tags

- `@smoke` - Include in fast feedback suite
- `@regression` - Include in full regression suite
- `@pending` - Awaiting implementation (CI skips these)
- `@wip` - Work in progress (skip in CI)
- `@hotfix` - Emergency addition (expedited review)

**Example:**
```gherkin
@story-PROJ-1234 @smoke
Scenario: User login with valid credentials
  Given a user exists with email "user@example.com" and password "pass123"
  When I log in with email "user@example.com" and password "pass123"
  Then I should be logged in
```

---

## Anti-Patterns to Avoid

### 1. Conjunctive Steps (And-ing)

**Bad:**
```gherkin
Given a user exists and is logged in and has admin privileges
```

**Good:**
```gherkin
Given a user exists with email "admin@example.com"
And the user has role "admin"
And I am logged in as "admin@example.com"
```

### 2. Incidental Details

**Bad:**
```gherkin
Given I navigate to https://example.com/login
And I wait for the page to load
And I fill in the username field with "user@example.com"
And I fill in the password field with "pass123"
And I check the "Remember me" checkbox
When I click the "Login" button with id="login-btn"
```

**Good:**
```gherkin
When I log in with email "user@example.com" and password "pass123"
```

### 3. Assertion in Given

**Bad:**
```gherkin
Given a user should exist with email "user@example.com"
```

**Good:**
```gherkin
Given a user exists with email "user@example.com"
```

_"Should" is for assertions (Then), not preconditions (Given)._

---

## Documentation

Document new steps in `testing.md`:

```markdown
## Step Library

### Authentication Steps

| Step | Parameters | Notes |
|------|------------|-------|
| `Given a user exists with email {string}` | email | Creates test user in database |
| `Given I am logged in as {string}` | email | Authenticates user and sets session |
| `When I log in with email {string} and password {string}` | email, password | Submits login form |
| `Then I should be logged in` | none | Verifies session exists |
```

---

## Next Steps

- **For step library patterns:** See [Context Files Guide - testing.md](context-files.md)
- **For workflow integration:** See [Workflow Guide](workflow.md)
- **For examples:** See [Sample Project](../sample-project/)
