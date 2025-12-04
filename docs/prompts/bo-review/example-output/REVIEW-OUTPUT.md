# Business Owner Review: Password Reset Feature (PROJ-1234)

**Reviewer:** Business Owner (using bo-review prompt)
**PR:** Draft spec from requirements-analyst
**Review Date:** 2025-12-04
**Decision:** 🔴 **CHANGES REQUESTED**

---

## Executive Summary

The AI-drafted specification provides a solid foundation with comprehensive security coverage (token expiry, single-use tokens, rate limiting) and consistent audit logging. However, **7 critical issues** require correction before approval:

- **1 security vulnerability** (email enumeration attack)
- **2 missing business rules** (admin password complexity, email content validation)
- **2 UX/clarity issues** (error message presentation, rate limiting logic)
- **2 missing scenarios** (multiple active tokens, async email delivery)

The spec demonstrates the AI correctly captured most requirements, but business owner review caught critical security, compliance, and user experience gaps.

---

## 🔴 Critical Issues (Must Fix)

### Issue 1: Email Enumeration Security Vulnerability

**Location:** Scenario "User attempts reset for nonexistent email"

**Current behavior:**
```gherkin
Then I should see "If that email is registered, you will receive a reset link"
```

**Problem:** The spec contradicts itself. Scenario 1 shows users receiving an email immediately, but Scenario 2 claims the system shows a generic message. An attacker can distinguish registered vs. unregistered emails by:
1. Timing attacks (database lookup vs. no lookup)
2. Observing whether email actually arrives
3. Subsequent API responses that might leak user existence

**Business Impact:** Violates privacy requirements and exposes registered user emails to enumeration attacks.

**Required Fix:**
1. System MUST always return the same message for both registered and unregistered emails
2. System MUST delay response time to prevent timing attacks (simulate email send time)
3. Consider adding to business rules document: "BR-006: System must not reveal whether an email is registered"

**Specific change needed:**
```gherkin
Scenario: System prevents email enumeration attacks
  When I request a password reset for "nonexistent@example.com"
  Then I should see "If that email is registered, you will receive a reset link"
  And the response time should be consistent with successful requests
  And the password reset attempt should be logged with "email_not_found"

Scenario: System handles registered email consistently
  Given a user exists with email "user@example.com"
  When I request a password reset for "user@example.com"
  Then I should see "If that email is registered, you will receive a reset link"
  And I should receive a reset email
  And the password reset attempt should be logged
```

---

### Issue 2: Missing Business Rule - Admin Password Complexity

**Location:** Scenario "User sets password that fails complexity requirements"

**Current behavior:**
```gherkin
When I set a new password to "weak" using the reset token
Then I should see "Password must be at least 8 characters"
```

**Problem:** According to BR-003 in the requirements, administrator accounts require 12+ character passwords. The spec doesn't distinguish between regular users (8+ chars) and admins (12+ chars).

**Business Impact:** Compliance violation - admin accounts could be created with insufficient password strength.

**Required Fix:** Add scenario for admin password reset:

```gherkin
@security @admin
Scenario: Administrator password requires 12+ characters
  Given an administrator exists with email "admin@example.com"
  And I requested a password reset 1 hour ago
  When I set a new password to "Short123!" using the reset token
  Then I should see "Administrator passwords must be at least 12 characters"
  And my password should not be changed
  And the password reset attempt should be logged with "validation_failed"

@security @admin
Scenario: Administrator successfully resets password with valid complexity
  Given an administrator exists with email "admin@example.com"
  And I requested a password reset 1 hour ago
  When I set a new password to "AdminPass123!Secure" using the reset token
  Then I should be logged in
  And my old password should be invalidated
  And I should receive a confirmation email
```

**Question for Team:** Do we need to enforce the special character and number requirements mentioned in BR-003, or just the length? Spec currently only validates length.

---

### Issue 3: Missing Email Content/Link Validation

**Location:** Scenario "User requests password reset successfully"

**Current step:**
```gherkin
Then I should receive a reset email
And the reset link should expire in 24 hours
```

**Problem:** Spec doesn't verify the email contains the actual reset link or has appropriate content. In BDD, we verify behavior, not just that "an email was sent."

**Business Impact:** If email service sends an email without the link, or with wrong content, tests pass but feature is broken.

**Required Fix:** Add explicit verification:

```gherkin
Then I should receive a reset email
And the email should contain a valid reset link
And the reset link should include the token
And the reset link should expire in 24 hours
And the email should include security notice "If you did not request this, please ignore"
```

**Question for Dev Team:** Should the reset link be a deep link to the app, or a web URL? This affects step implementation.

---

### Issue 4: Error Message UX Issue

**Location:** Scenario "User sets password that fails complexity requirements"

**Current behavior:**
```gherkin
When I set a new password to "weak" using the reset token
Then I should see "Password must be at least 8 characters"
And I should see "Password must contain uppercase letters"
And I should see "Password must contain lowercase letters"
```

**Problem:** Three separate "Then I should see" steps implies messages are shown one at a time. This creates poor UX - user fixes one issue, submits, gets another error, repeats.

**Business Impact:** Frustrating user experience, increased support tickets, higher abandonment rate.

**Required Fix:** Show all errors at once:

```gherkin
When I set a new password to "weak" using the reset token
Then I should see all password requirements:
  | requirement                               |
  | Password must be at least 8 characters    |
  | Password must contain uppercase letters   |
  | Password must contain lowercase letters   |
```

Or if that's too complex for first pass:

```gherkin
Then I should see "Password must meet the following requirements: at least 8 characters, uppercase letters, lowercase letters"
```

---

### Issue 5: Ambiguous Rate Limiting Logic

**Location:** Scenario "System enforces rate limiting"

**Current behavior:**
```gherkin
Given I have requested password resets 5 times in the last hour
When I request a password reset for "user@example.com"
Then I should see "Too many password reset attempts"
```

**Problem:** Ambiguous whether the 6th request fails or the 5th request fails. "Given I have requested 5 times" + "When I request" = is this the 6th attempt? Or did the 5th attempt already fail?

**Business Impact:** Unclear requirement leads to implementation disagreement and potential security gap.

**Required Fix:** Make it explicit:

```gherkin
Scenario: System blocks 6th password reset attempt within 1 hour
  Given a user exists with email "user@example.com"
  And I have successfully requested password resets 5 times in the last hour
  When I request a 6th password reset for "user@example.com"
  Then I should see "Too many password reset attempts. Please try again in 1 hour."
  And no email should be sent
  And the password reset attempt should be logged with "rate_limit_exceeded"
```

**Question:** Should the error message tell users how long to wait? Current message doesn't specify the 1-hour window.

---

### Issue 6: Missing Scenario - Multiple Active Reset Tokens

**Problem:** What happens if user requests reset multiple times before using any token? Current spec shows:
- Scenario 3: "I requested a password reset 1 hour ago" (singular)
- But no coverage for: user requests reset, then requests again, then tries to use first token

**Business Impact:** Security risk if old tokens remain valid, or UX issue if all tokens are invalidated and user is confused.

**Required Scenarios:**

```gherkin
@security @edge-case
Scenario: New reset request invalidates previous tokens
  Given a user exists with email "user@example.com"
  And I requested a password reset 1 hour ago
  And I saved the first reset token
  When I request a new password reset for "user@example.com"
  And I try to use the first reset token to set password to "NewPass123!"
  Then I should see "This reset link has expired or is invalid"
  And my password should not be changed

@security @edge-case
Scenario: Only most recent reset token is valid
  Given a user exists with email "user@example.com"
  And I requested a password reset 2 hours ago
  And I requested a new password reset 1 hour ago
  And I saved the second reset token
  When I set a new password to "NewPass123!" using the second token
  Then I should be logged in
  And my old password should be invalidated
```

**Question:** Should we invalidate all previous tokens, or allow 1 active token per user at a time?

---

### Issue 7: Async Email Delivery Assumption

**Location:** Multiple scenarios reference "I should receive a reset email"

**Problem:** Email delivery is typically async. Step "Then I should receive a reset email" implies immediate delivery within the scenario execution. In production:
- SMTP may queue emails
- Rate limiting may delay delivery
- Transactional email services (SendGrid, etc.) may take seconds

**Business Impact:** Tests may pass in local environment with mock email service but fail in staging/production with real email provider.

**Required Clarification:** Add note to spec or adjust language:

**Option 1** - Make it explicit in spec:
```gherkin
# Note: "should receive" means email is queued for delivery, not necessarily delivered
Then I should receive a reset email  # Email is queued
```

**Option 2** - Adjust language to be implementation-agnostic:
```gherkin
Then a reset email should be sent to "user@example.com"
```

**Question for Dev Team:** How should step definitions handle email verification? Mock service in BDD tests? Check outbox queue?

---

## ✅ What Worked Well

1. **Comprehensive Security Coverage:**
   - Token expiry (24 hours)
   - Single-use tokens (Scenario 5)
   - Rate limiting (Scenario 6)
   - Audit logging on all paths

2. **Good Tag Organization:**
   - `@happy-path` for core flows
   - `@security` for security-critical scenarios
   - `@validation` for input validation
   - `@edge-case` for boundary conditions

3. **Consistent Audit Logging:**
   - Every scenario includes audit log verification
   - Shows understanding of compliance requirements

4. **Boundary Condition Testing:**
   - Expired tokens (Scenario 4)
   - Already-used tokens (Scenario 5)
   - Rate limiting (Scenario 6)

5. **Clear Given-When-Then Structure:**
   - Easy to understand
   - No ambiguous language (except issues noted above)

---

## Questions for Requirements Analyst / Dev Team

1. **Email content:** Should reset link be app deep link or web URL?
2. **Rate limiting window:** Should error message tell users how long to wait?
3. **Admin passwords:** Enforce special char/number requirements or just length?
4. **Multiple tokens:** Invalidate all previous tokens, or allow 1 active token?
5. **Email delivery:** How should BDD steps verify email in tests vs. production?
6. **Password validation:** Should we validate on submission or show live feedback as user types?

---

## Recommended Next Steps

1. **Requirements Analyst:** Address 7 critical issues above, update spec
2. **Security Review:** Verify email enumeration fix with security team
3. **Compliance Check:** Confirm admin password requirements with compliance officer
4. **UX Review:** Validate error message approach with design team
5. **BO Re-Review:** Once updated, I'll review again (should be quick second pass)

---

## Decision Rationale

**Why CHANGES REQUESTED instead of APPROVED:**

The foundation is solid (8/10 scenarios are well-written), but the 7 issues above include:
- 1 security vulnerability that could expose user data
- 2 compliance gaps that could cause audit failures
- 2 UX issues that could frustrate users
- 2 missing edge cases that could cause production bugs

These are not minor wording issues - they represent gaps in business logic that would cause real problems in production. The AI did excellent work capturing the core requirements, but BO review is essential for catching:
- Security implications (email enumeration)
- Business rule exceptions (admin passwords)
- Real-world edge cases (multiple tokens)
- User experience considerations (error messages)

**Estimated time to fix:** 30-45 minutes for Requirements Analyst to address issues and update spec.

**Estimated time for re-review:** 10-15 minutes (focused review of changes only)

---

## Metadata

- **Confidence Level:** HIGH (spec is well-structured, issues are clear and fixable)
- **Risk Level:** MEDIUM (security issue is serious, but easy to fix)
- **Review Time:** ~15 minutes
- **AI Performance:** 85% (captured most requirements correctly, missed security and edge cases)
