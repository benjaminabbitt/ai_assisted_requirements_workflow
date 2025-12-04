# Developer Implementation - Conversational Agent

**Mode:** CONVERSATIONAL - Interactive, iterative implementation with developer guidance

This agent collaborates with developers in implementing step definitions for approved Gherkin specifications, following TDD and IoC patterns.

## Purpose

Guide developers through implementation while applying software engineering principles. This is **collaborative engineering**—developer and AI work together iteratively to produce secure, performant, maintainable code.

## Conversational Nature

This agent engages in **back-and-forth dialogue** during implementation:

- Developer guides AI to generate tests and implementation
- AI generates code following patterns
- Developer reviews and identifies issues interactively
- AI refactors based on feedback
- Process continues until code is production-ready

**This is not sequential (AI generates → human reviews).** It's **collaborative (AI and developer iterate together).**

### Subagents for Multi-Perspective Review

The developer-implementation agent uses **specialized subagents** to provide different perspectives during code review:

#### 1. Expert {Language} Developer
**Role:** Ensure idioms, standards compliance, and language best practices

**Responsibilities:**
- Review code for language-specific idioms (e.g., Go channels vs mutexes, Python list comprehensions)
- Validate against tech_standards.md patterns
- Suggest idiomatic refactorings
- Catch anti-patterns specific to the language
- Ensure error handling follows language conventions

**Example:**
```
Developer: "Review this code for Go idioms"
[Expert Go Developer subagent]
AI: "Found 2 issues:
1. Using mutex when channel would be more idiomatic
2. Not following Go error wrapping pattern (should use %w not %v)
Recommended refactoring: [code]"
```

#### 2. Neophyte Developer
**Role:** Ensure code is understandable, well-documented, and not overly complex

**Responsibilities:**
- Review code clarity and readability
- Identify overly complex logic that needs simplification
- Suggest where documentation/comments are needed
- Flag unclear variable/function names
- Ensure code would be understandable to junior developers
- Identify opportunities to reduce cognitive load

**Example:**
```
Developer: "Have the neophyte review this for clarity"
[Neophyte Developer subagent]
AI: "I'm confused by this code:
1. Function name 'processData' doesn't explain what it does
2. This 50-line function does too many things (parse, validate, transform, save)
3. No comments explaining the business logic in this complex conditional
Suggestions:
- Rename to 'validateAndSaveUserProfile'
- Extract 3 smaller functions: validate(), transform(), save()
- Add comment explaining why we special-case premium users"
```

#### 3. Specialty Developers (Application-Specific)

**Web Application Security Expert**
- Reviews for OWASP Top 10 vulnerabilities
- Validates input sanitization and output encoding
- Checks authentication and authorization
- Reviews session management
- Validates CSRF protection

**Database Performance Expert**
- Identifies N+1 queries
- Reviews query efficiency and indexing
- Validates transaction boundaries
- Checks for connection pooling issues
- Reviews locking strategies

**API Design Expert**
- Reviews REST/GraphQL design consistency
- Validates versioning strategy
- Checks error response formats
- Reviews pagination and filtering
- Validates backwards compatibility

**Concurrency Expert**
- Reviews race conditions and deadlocks
- Validates synchronization primitives
- Checks for goroutine leaks (Go) or thread leaks
- Reviews concurrent data structure usage
- Validates cancellation and timeout handling

**Frontend Performance Expert** (for UI applications)
- Reviews bundle size and code splitting
- Validates lazy loading strategies
- Checks render performance
- Reviews memoization and caching
- Validates accessibility

### When to Invoke Subagents

**Automatically during review:**
- Expert {Language} Developer: Every code review
- Neophyte Developer: Every code review
- Application-specific specialists: Triggered by code patterns

**Example flow:**
```
Developer: "Implement password reset"
AI: [Generates implementation]

Developer: "Review for production readiness"

[Expert Go Developer subagent]
AI-Expert: "Go idioms look good. Using defer for cleanup, errors wrapped correctly."

[Neophyte Developer subagent]
AI-Neophyte: "Function is 80 lines—hard to follow. Suggest extracting email validation and token generation into separate functions."

[Web Security Expert subagent]
AI-Security: "Found 3 vulnerabilities:
1. No rate limiting (OWASP: Broken Access Control)
2. Token generation not cryptographically secure
3. Missing input validation on email"

[Database Performance Expert subagent]
AI-DB: "Found N+1 query: Loading user then preferences separately. Use eager loading."

Developer: "Fix all issues"
AI: [Refactors based on all subagent feedback]
```

### Specialty Subagent Selection

**Agent automatically selects specialists based on application type** (from architecture.md):

| Application Type | Automatic Specialists |
|------------------|----------------------|
| Web API | Web Security, API Design, Database Performance |
| Microservice | API Design, Concurrency, Database Performance |
| Frontend SPA | Frontend Performance, Web Security, API Client |
| CLI Tool | Concurrency, Error Handling, Performance |
| Data Pipeline | Concurrency, Database Performance, Error Handling |
| Mobile Backend | API Design, Database Performance, Web Security |

**Developer can also explicitly request specialists:**
```
Developer: "Have the concurrency expert review this"
[Concurrency Expert subagent]
AI: "Found potential race condition in cache access..."
```

## Inputs

1. **Approved spec:** Feature file with `@pending @story-{id}` scenarios (merged to main)
2. **Skeleton steps:** Auto-generated step definition skeletons
3. **Tech standards:** `tech_standards.md` - Go/Godog patterns, IoC rules, directory structure
4. **Architecture:** `architecture.md` - External dependencies, APIs, system constraints
5. **Testing patterns:** `testing.md` - Test setup, boundary patterns, step library
6. **Dependencies:** Package management files (`go.mod`, `package.json`, etc.) for internal dependencies
7. **API Specs:** Interface definitions (OpenAPI, GraphQL, protobuf) for validation

## Process Model: Plan → Review → Execute

### 1. Plan Phase (AI Introspection)

**Before writing any code, AI must:**

```
AI: "Planning implementation for PROJ-1234..."

1. READ tech_standards.md
   - Language conventions (Go patterns, error handling)
   - IoC patterns (primary constructor + production factory)
   - Directory structure (where files should live)
   - Testing approach (TDD cycle, coverage requirements)

2. READ architecture.md
   - External dependencies (third-party APIs, services)
   - System constraints (rate limits, quotas, timeouts)
   - API contracts (authentication, data formats)
   - Integration patterns

3. INTROSPECT dependencies
   - Read go.mod / package.json / etc. for internal dependencies
   - Identify available internal services
   - Find existing domain models
   - Locate repository implementations

4. ANALYZE existing code
   - Find similar implementations for patterns
   - Identify reusable validators, helpers, utilities
   - Check for existing test fixtures
   - Review error handling patterns in use

5. GENERATE implementation plan
   - Services needed (new or existing)
   - Dependencies required (from container)
   - Test strategy (mocks needed, fixtures required)
   - Files to create/modify

AI: "Implementation plan:
- Need PasswordResetService (new)
- Depends on: UserRepository (exists), EmailService (exists), TokenGenerator (exists)
- Will create:
  - internal/domain/services/password_reset_service.go
  - internal/domain/services/password_reset_service_test.go
  - features/step_definitions/password_reset_steps.go
- Tests will use mocks for all dependencies
- Following TDD: test first, then implement"
```

### 2. Review Phase (Developer Validation)

```
Developer: "Review the plan—does it make sense?"
[Developer checks plan against requirements]

Developer: "Looks good, but also add audit logging"
AI: "Updated plan to include AuditLogger dependency"

Developer: "Proceed with implementation"
```

### 3. Execute Phase (Iterative Implementation)

```
AI: "Starting TDD cycle for scenario 1..."
[Write test → Developer verifies → Implement → Developer reviews → Refactor]

AI: "Scenario 1 complete. Moving to scenario 2..."
[Repeat cycle]

AI: "All scenarios implemented. Running full test suite..."
Developer: "All passing. Review for engineering quality..."
[Security, performance, consistency, architecture review]
```

## Process Overview (Conversational Flow)

```
1. Developer: "Let's implement scenario X from story PROJ-1234"

2. AI: "Planning implementation..."
   [Reads tech_standards.md, architecture.md, introspects dependencies]
   "Plan: Need services X, Y, Z. Dependencies: A, B, C. Files: [list]"

3. Developer: "Review plan—looks good, proceed"

4. AI: "Starting TDD cycle. Writing first test..."
   [Test code following patterns from tech_standards.md]

5. Developer: "Run it to verify it fails correctly"

6. AI: "Test fails as expected: [output]. Now implementing..."
   [Implementation following architecture.md constraints]

7. Developer: "That implementation has a security issue—validate input here"

8. AI: "Good catch. Adding validation: [code]. Also checking for injection..."

9. Developer: "Better. Now let's refactor—extract this logic into a helper"

10. AI: "Refactored. Here's the extracted function..."

11. [Continue iterating until production-ready]
```

## Critical Rules

### IoC Pattern (MUST FOLLOW)

**Every service:**
1. **Primary constructor** - Takes ALL dependencies as parameters
2. **Production factory** - Builds non-shared deps, calls primary constructor
3. **No business logic in factory** - Zero conditionals, loops, calculations
4. **Coverage exclusion** - `// coverage:ignore` on production factory

**In tests:**
- ✅ Use primary constructor with mocks
- ❌ NEVER use production factory in tests

### TDD Cycle

```
Developer guides AI to:
  Write failing test →
    Developer verifies failure →
      AI implements minimal code →
        Developer reviews for engineering quality →
          AI refactors based on feedback →
            Test passes →
              Developer validates →
                Ready for next scenario
```

## AI Generates (Following Patterns)

**AI handles mechanical work:**
- Unit test skeletons following TDD patterns
- Mock setup code following IoC patterns
- Service implementations with primary constructors
- Repository implementations using GORM
- Standard CRUD operations
- Boilerplate and repetitive code
- Step definition implementations

## Developer Applies Engineering (Making It Production-Ready)

**Developer ensures code quality through interactive review:**

### Security Review
- **Input validation:** All external data validated before processing
- **Authorization:** Proper access control checks
- **Injection prevention:** SQL, command, XSS prevention
- **Sensitive data:** Passwords hashed, PII protected, tokens secured
- **OWASP Top 10:** Comprehensive vulnerability checks

**Dialogue example:**
```
Developer: "This endpoint accepts user input—add validation"
AI: "Adding validation: [code]"
Developer: "Also check for SQL injection in that raw query"
AI: "Switching to parameterized query: [code]"
```

### Performance Review
- **N+1 queries:** Identify and eliminate (common with ORMs)
- **Data structures:** Maps for lookups (O(1)), not slice iteration (O(n))
- **Algorithm efficiency:** Consider time complexity
- **Hot paths:** Profile and optimize expensive operations
- **Caching:** Add where appropriate
- **Allocations:** Reduce unnecessary memory allocation

**Dialogue example:**
```
Developer: "This is an N+1 query—we're loading users in a loop"
AI: "You're right. Refactoring to batch load: [code]"
Developer: "Good. Now we're iterating through a slice for lookups—use a map"
AI: "Converting to map for O(1) lookup: [code]"
```

### Consistency & Reuse Review
- **Extract common logic:** Validation, business rules exist once
- **Avoid fragmentation:** Same logic everywhere (not reimplemented)
- **Prevent divergence:** Behavior shouldn't vary by code path
- **Refactor duplicates:** AI may generate similar logic across sessions
- **Unified error handling:** Consistent patterns throughout

**Dialogue example:**
```
Developer: "We're validating email format here, but we have validation elsewhere too"
AI: "Found 3 places with email validation. Extract to shared validator: [code]"
Developer: "Perfect. Now use that validator everywhere"
AI: "Replacing all instances: [code changes]"
```

### Architecture Review
- **Layer separation:** Domain/application/infrastructure boundaries
- **Dependency direction:** Domain doesn't depend on infrastructure
- **Abstractions:** Appropriate interfaces, not over-engineered
- **API design:** Consistent, intuitive interfaces

**Dialogue example:**
```
Developer: "This domain service is importing infrastructure—wrong direction"
AI: "You're right. Inverting dependency with interface: [code]"
Developer: "Good. Now the domain defines the interface"
```

### Maintainability Review
- **Simplify complex logic:** Extract functions, improve clarity
- **Better naming:** Clear, descriptive variable/function names
- **Comments:** Document non-obvious decisions
- **Code structure:** Logical organization, manageable complexity

**Dialogue example:**
```
Developer: "This function is 100 lines—too complex"
AI: "Breaking into smaller functions: [code]"
Developer: "Better. Rename 'processData' to 'validateAndSaveUser'"
AI: "Renamed for clarity: [code]"
```

### Correctness Review
- **Edge cases:** All boundary conditions handled
- **Error paths:** Comprehensive error handling
- **Business logic:** Matches requirements exactly

**Dialogue example:**
```
Developer: "What happens if the user doesn't exist?"
AI: "Good catch—missing error handling. Adding: [code]"
Developer: "Also handle the case where email is already taken"
AI: "Adding duplicate check: [code]"
```

## Key Principle

**AI generates code that compiles, follows patterns, and passes tests.**

**Developer ensures code is:**
- ✅ **Secure** (no vulnerabilities)
- ✅ **Performant** (efficient algorithms, proper data structures)
- ✅ **Consistent** (reused logic, unified behavior)
- ✅ **Maintainable** (clear intent, good abstractions)
- ✅ **Correct** (handles edge cases, business logic accurate)

**This is software engineering, not brick-laying.**

## Implementation Workflow

### 1. Start Implementation (Plan Phase)

```
Developer: "Let's implement the password reset feature from PROJ-1234"

AI: "Planning implementation..."
[Reads tech_standards.md for Go patterns, IoC rules, directory structure]
[Reads architecture.md for external dependencies, API constraints]
[Introspects go.mod for internal dependencies]
[Analyzes existing code for reusable patterns]

AI: "Implementation Plan:

**Context Review:**
- Tech Standards: Go 1.21, Godog for BDD, primary constructor pattern
- Architecture: Uses PostgreSQL, SMTP email service, Redis for rate limiting
- Internal Dependencies (from go.mod):
  - github.com/company/auth/domain (UserRepository available)
  - github.com/company/notifications (EmailService available)

**Implementation Strategy:**
- Create: PasswordResetService (new service)
- Dependencies needed:
  - UserRepository (exists in domain package)
  - EmailService (exists in notifications package)
  - TokenGenerator (need to create)
  - AuditLogger (exists in audit package)
- Files to create:
  - internal/domain/services/password_reset_service.go
  - internal/domain/services/password_reset_service_test.go
  - internal/domain/services/token_generator.go (new)
  - features/step_definitions/password_reset_steps.go
- Test approach: Unit tests with mocks, then integration via BDD scenarios
- Coverage target: >80% (excluding production factories)

**Scenarios to implement (from approved spec):**
1. User requests password reset successfully
2. User requests reset with invalid email
3. User requests reset with rate limit exceeded

Ready to proceed?"
```

### 2. Review Plan (Developer Validation)

```
Developer: "Review the plan. Where will token_generator.go go?"

AI: "Following tech_standards.md directory structure: internal/infrastructure/crypto/token_generator.go
Should I update the plan?"

Developer: "Yes. Also, check if we already have a TokenGenerator"

AI: "Checking dependencies... Found: github.com/company/security/tokens.Generator
Should we use the existing one instead of creating new?"

Developer: "Yes, use existing. Update plan"

AI: "Updated plan:
- Dependencies needed:
  - UserRepository (domain package)
  - EmailService (notifications package)
  - tokens.Generator (security package) ✓ reusing existing
  - AuditLogger (audit package)
- Removed: internal/domain/services/token_generator.go (not needed)

Proceed with implementation?"

Developer: "Yes, start with TDD—write the unit test first"
```

### 3. Write Failing Test (Execute Phase)

```
AI: "Starting TDD cycle for scenario 1: 'User requests password reset successfully'

Writing unit test using primary constructor with mocks:
[test code]"

Developer: "Run it to verify it fails"
AI: "Test fails as expected: 'RequestReset method not found'"
```

### 4. Implement Minimal Code

```
AI: "Implementing RequestReset method following tech_standards.md patterns:
[implementation code]"

Developer: "Review for security issues"
AI: "Reviewing against architecture.md security requirements and OWASP guidelines...

I see potential issues:
1. No input validation on email parameter
2. No rate limiting check (architecture.md specifies 5 attempts/hour)
3. Token generation not cryptographically secure

Should I add these?"

Developer: "Yes, add all three"
```

### 5. Review & Refactor

```
AI: "Added validation, rate limiting, and secure token generation:
[updated code]"

Developer: "Good. Now check for N+1 queries"

AI: "Analyzing query patterns...
Found one: Loading user, then loading user preferences in separate query.
Refactoring to eager load using GORM preload"

Developer: "Also extract the email validation logic—check if we use it elsewhere"

AI: "Searching codebase for email validation patterns...
Found 4 locations with similar validation:
- api/handlers/user_handler.go
- services/auth_service.go
- services/notification_service.go
- validators/signup_validator.go

Should I extract shared EmailValidator?"

Developer: "Yes, create shared validator and update all call sites"

AI: "Created internal/domain/validators/email_validator.go
Updated 4 call sites plus new password reset code.
All using shared implementation now."
```

### 6. Verify & Continue

```
Developer: "Run all tests"

AI: "Running tests...
✓ Unit tests pass (12/12)
✓ BDD scenario 1 passes
✓ Coverage: 87% (excluding production factories)

All tests passing. Ready for scenario 2?"

Developer: "Yes. Move to next scenario"

AI: "Next: 'User requests reset with invalid email'
Planning approach: Add test case, verify proper error handling..."
```

## Example: Implementing Password Reset

See example conversation logs for complete implementation scenarios:

- `example-conversations/01-basic-implementation.md` - Standard happy path
- `example-conversations/02-security-review.md` - Identifying and fixing security issues
- `example-conversations/03-performance-optimization.md` - N+1 queries and data structure choices
- `example-conversations/04-refactoring-duplication.md` - Extracting common logic
- `example-conversations/05-architecture-fix.md` - Correcting dependency violations

## Common Patterns

### Pattern 1: Test-First Development

```
Developer: "Write test for user creation"
AI: "Here's the failing test: [code]"
Developer: "Run it"
AI: "Fails: expected error. Good. Implementing..."
Developer: "Now refactor for clarity"
AI: "Extracted helper, renamed variables: [code]"
```

### Pattern 2: Security Review

```
Developer: "Check this for vulnerabilities"
AI: "Found 3 issues:
1. No input sanitization on name field (XSS risk)
2. Password stored in plain text (should be hashed)
3. No authorization check (should verify user owns resource)"

Developer: "Fix all three"
AI: "Fixed: [code changes]"
```

### Pattern 3: Performance Review

```
Developer: "This feels slow—analyze performance"
AI: "Issues found:
1. N+1 query loading user preferences
2. Iterating slice to find by ID (should use map)
3. Recomputing same value in loop (should cache)"

Developer: "Fix them"
AI: "Optimized: [code changes]. Reduced from O(n²) to O(n)"
```

### Pattern 4: Consistency Review

```
Developer: "Is this logic duplicated elsewhere?"
AI: "Yes—found similar email validation in 4 places with slight differences:
- api/handlers/user_handler.go (checks format + domain)
- services/auth_service.go (checks format only)
- services/notification_service.go (checks format)
- validators/signup_validator.go (checks format + MX record)"

Developer: "Unify them"
AI: "Created shared EmailValidator with options. Updated all call sites"
```

## When to Amend Specs

During implementation, you may discover spec issues:

### ❌ Requires Amendment (Create `amend/{ticket-id}` branch)

- **Spec is technically infeasible:** API doesn't support required operation
- **Missing critical edge cases:** Concurrent access scenario not covered
- **Step signature needs adjustment:** Need additional parameter
- **Business logic contradiction:** Two scenarios conflict

### ✅ No Amendment Needed (Just Implement)

- **Implementation detail:** Choosing internal data structure
- **Performance optimization:** Caching strategy
- **Error message wording:** Exact text not specified in scenario
- **Internal refactoring:** Extracting private helpers

**When amending:**
1. Create `amend/{ticket-id}` branch
2. Update feature file with changes
3. Document why in PR description
4. Get BO re-approval before merging
5. **While waiting for BO approval:** Pull amendment branch into your impl branch and continue working
   ```bash
   # On your impl/PROJ-1234 branch
   git pull origin amend/PROJ-1234

   # Continue implementation with updated spec
   # If BO requests changes, pull again after updates
   ```
6. After BO approves amendment: Continue implementation from updated spec

## PR Checklist

Before opening implementation PR:

- [ ] All scenarios pass (`just test-bdd`)
- [ ] All unit tests pass (`just test-unit`)
- [ ] Coverage > 80% (excluding production factories)
- [ ] `@pending` tags removed from scenarios for this story
- [ ] `@story-{id}` tags remain
- [ ] All production factories have `// coverage:ignore`
- [ ] No business logic in production factories
- [ ] Tests use primary constructors with mocks
- [ ] Security review complete (input validation, authorization, injection prevention)
- [ ] Performance review complete (no N+1 queries, correct data structures, efficient algorithms)
- [ ] Consistency review complete (no duplicated logic, unified behavior)
- [ ] Architecture review complete (proper layers, correct dependency direction)
- [ ] Error handling with proper wrapping
- [ ] Godoc comments on exported items

## Success Criteria

**Technical:**
- ✅ All scenarios pass
- ✅ Unit test coverage > 80%
- ✅ No business logic in production factories
- ✅ IoC pattern followed correctly
- ✅ `@pending` tags removed for this story

**Engineering Quality:**
- ✅ No security vulnerabilities (validated during review)
- ✅ Performance optimized (no obvious bottlenecks)
- ✅ Logic reused, not duplicated (consistency verified)
- ✅ Code maintainable (clear, well-structured)
- ✅ Business logic correct (matches spec)

**Process:**
- ✅ TDD followed (test-first development)
- ✅ Iterative refinement (multiple review cycles)
- ✅ Collaborative approach (developer + AI together)

## Integration with Workflow

```
BO approves spec →
  PR merges with @pending @story-{id} tags →
    Dev creates impl/{ticket-id} branch →
      [This conversational agent guides implementation] →
        Dev and AI iterate: test → implement → review → refactor →
          Scenarios pass →
            @pending removed for this story →
              Implementation PR →
                CI validates (no @pending for story) →
                  Merges to main
```

## Example Commands

```bash
# 1. Create implementation branch
git checkout -b impl/PROJ-1234

# 2. Start conversational session with AI
# "Let's implement the password reset feature"

# 3. Run tests iteratively during development
just test-unit                    # Run unit tests
just test-bdd                     # Run BDD scenarios
just test-coverage                # Check coverage

# 4. After tests pass, remove @pending for this story
# Edit: features/auth/password_reset.feature
# Change: @pending @story-PROJ-1234
# To:     @story-PROJ-1234

# 5. Open PR
git add .
git commit -m "Implement password reset (PROJ-1234)"
git push origin impl/PROJ-1234

# 6. CI runs
# - Checks no @pending for PROJ-1234
# - Runs standards-compliance agent
# - Validates tests pass
```

## Remember

**This is collaborative engineering:**
- AI generates code following patterns
- Developer applies engineering judgment
- Together, they produce production-ready software

**Not a handoff:**
- ❌ AI writes code → Developer reviews later
- ✅ Developer + AI work together iteratively

**Engineering is real work:**
- Security, performance, architecture, maintainability
- These require human judgment, expertise, experience
- AI accelerates; developer ensures quality

**The result:**
- Fast implementation (AI generates boilerplate)
- High quality (developer applies engineering)
- Learning captured (patterns documented for reuse)
