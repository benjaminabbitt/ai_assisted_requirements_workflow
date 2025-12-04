# Developer Implementation - Conversational Agent

## Purpose

**Collaborate with developers** in implementing step definitions for approved Gherkin specifications while following TDD and IoC patterns.

**Mode:** CONVERSATIONAL - Interactive, iterative implementation with developer guidance

## Key Principle

This is **collaborative engineering**—developer and AI work together iteratively to produce secure, performant, maintainable code.

- ❌ **Not:** AI writes code → Developer reviews later (sequential)
- ✅ **Instead:** Developer + AI work together iteratively (collaborative)

## Files in This Directory

- **`prompt.md`** - Complete conversational implementation guide
  - Plan → Review → Execute model
  - Context file introspection (tech_standards.md, architecture.md, dependencies)
  - TDD cycle with interactive review
  - Engineering quality checklist (security, performance, consistency, architecture)

- **`example-conversations/`** - Demonstrates conversational flow
  - `01-basic-implementation.md` - Plan-review-execute with basic password reset
  - `02-security-review.md` - Finding and fixing security vulnerabilities interactively
  - More examples showing performance, refactoring, architecture scenarios

- **`sample-spec.feature`** - Example approved spec with `@pending` scenarios

## Process Model: Plan → Review → Execute

### Phase 1: Plan (AI Introspection)

AI reads context and introspects codebase **before writing code:**

```
AI: "Planning implementation..."

1. READ tech_standards.md (Go patterns, IoC rules, directory structure)
2. READ architecture.md (external dependencies, API constraints)
3. INTROSPECT dependencies (go.mod, package.json, etc.)
4. ANALYZE existing code (find reusable patterns, validators, helpers)
5. GENERATE implementation plan

AI: "Plan:
- Services needed: PasswordResetService (new)
- Dependencies: UserRepository (exists), EmailService (exists), TokenGenerator (check)
- Files: [list]
- Test approach: Mocks for unit tests, test container for BDD
Ready to proceed?"
```

### Phase 2: Review (Developer Validation)

```
Developer: "Check if TokenGenerator exists before we create it"
AI: "Found existing TokenGenerator in security package. Updated plan to reuse."
Developer: "Good. Proceed with implementation"
```

### Phase 3: Execute (Iterative Implementation)

```
AI: "Writing test..." → Developer: "Run it" → AI: "Fails correctly" →
AI: "Implementing..." → Developer: "Review for security" → AI: "Found 3 issues" →
Developer: "Fix them" → AI: "Fixed" → Developer: "Refactor this" → AI: "Refactored"
[Continue until production-ready]
```

## How It Works

### 1. Context Introspection

**AI automatically reads:**
- `tech_standards.md` - Language, framework, IoC patterns, directory structure
- `architecture.md` - External dependencies, system constraints, API contracts
- `go.mod` / `package.json` / etc. - Internal dependencies
- Existing codebase - Patterns, validators, helpers, test fixtures

**AI uses this to:**
- Follow correct patterns (IoC, error handling, testing)
- Reuse existing code (validators, services, repositories)
- Place files in correct locations
- Generate compliant code from the start

### 2. TDD Cycle

```
1. AI writes failing test (using primary constructor + mocks)
2. Developer verifies test fails correctly
3. AI implements minimal code
4. Developer reviews for engineering quality:
   - Security (input validation, authorization, injection prevention)
   - Performance (N+1 queries, data structures, algorithm efficiency)
   - Consistency (extract common logic, unified behavior)
   - Architecture (proper layers, dependency direction)
   - Maintainability (clarity, naming, complexity)
5. AI refactors based on feedback
6. Test passes
7. Repeat for next scenario
```

### 3. Engineering Review (Not Just Correctness)

**Developer ensures code is:**
- ✅ **Secure** (no vulnerabilities)
- ✅ **Performant** (efficient algorithms, correct data structures)
- ✅ **Consistent** (logic reused, not duplicated)
- ✅ **Maintainable** (clear intent, good abstractions)
- ✅ **Correct** (business logic matches spec)

**See `example-conversations/02-security-review.md` for how this catches vulnerabilities during implementation.**

## AI vs. Human Responsibilities

### AI Generates (Following Patterns)

**Tests:**
- Unit test skeletons with mocks
- Table-driven test structures
- Boundary condition tests

**Implementation:**
- Services with primary constructor + production factory
- Repository implementations using GORM
- Standard CRUD operations
- Boilerplate code

### Developer Applies Engineering

**Security:**
- Input validation, authorization checks
- Injection prevention (SQL, XSS, command)
- Sensitive data handling (passwords, PII, tokens)

**Performance:**
- Fix N+1 queries
- Choose correct data structures (maps vs slices)
- Optimize algorithms (time complexity)
- Add caching where appropriate

**Consistency & Reuse:**
- Extract common logic into shared functions
- Ensure logic isn't fragmented across codebase
- Prevent divergent behavior in different code paths
- Monitor for duplication across features

**Architecture:**
- Validate proper layer separation
- Ensure correct dependency direction
- Review abstractions for appropriateness
- Design for maintainability and change

## Example Workflow

```bash
# 1. Create implementation branch
git checkout -b impl/PROJ-1234

# 2. Start conversational session
# Developer: "Let's implement password reset from PROJ-1234"
# AI: [Plans, introspects, generates plan]
# Developer: "Review plan... good, proceed"
# AI: [Writes test] → Developer: "Run it" → AI: [Implements]
# Developer: "Review for security" → AI: [Finds issues] → Developer: "Fix"
# [Continue iterating]

# 3. Run tests during development
just test-unit          # Unit tests
just test-bdd           # BDD scenarios
just test-coverage      # Coverage check

# 4. Remove @pending tags when complete
# Edit: features/auth/password_reset.feature
# Change: @pending @story-PROJ-1234
# To:     @story-PROJ-1234

# 5. Open PR
git add .
git commit -m "Implement password reset (PROJ-1234)"
git push origin impl/PROJ-1234

# 6. CI validates
# - No @pending for PROJ-1234
# - Standards-compliance check
# - All tests pass
```

## When to Amend Specs

During implementation, you may discover spec issues.

### Requires Amendment

- Spec is technically infeasible
- Missing critical edge cases
- Step signature needs adjustment
- Business logic contradiction

**Process:**
1. Create `amend/PROJ-1234` branch
2. Update feature file, document why
3. Get BO re-approval
4. **While waiting:** Pull amendment into impl branch and continue working
   ```bash
   git pull origin amend/PROJ-1234  # Continue with updated spec
   ```

### No Amendment Needed

- Implementation details (data structures, algorithms)
- Performance optimizations
- Error message wording (if not specified)
- Internal refactoring

## Success Criteria

**Technical:**
- ✅ All scenarios pass
- ✅ Unit test coverage > 80%
- ✅ No business logic in production factories
- ✅ IoC pattern followed
- ✅ `@pending` tags removed for this story

**Engineering Quality:**
- ✅ No security vulnerabilities
- ✅ Performance optimized
- ✅ Logic reused, not duplicated
- ✅ Code maintainable
- ✅ Business logic correct

**Process:**
- ✅ TDD followed (test-first)
- ✅ Iterative refinement (multiple review cycles)
- ✅ Collaborative approach (developer + AI together)

## Example Conversations

See `example-conversations/` for complete implementation scenarios demonstrating:

1. **Basic Implementation** (`01-basic-implementation.md`)
   - Plan → Review → Execute cycle
   - Context introspection and dependency discovery
   - TDD with interactive review
   - Refactoring based on feedback

2. **Security Review** (`02-security-review.md`)
   - Finding authorization vulnerabilities
   - Preventing XSS injection
   - Fixing mass assignment issues
   - Adding comprehensive security test coverage

More examples showing:
- Performance optimization (N+1 queries, data structures)
- Refactoring duplication
- Architecture fixes
- Handling amended specs

## Integration with Workflow

```
BO approves spec →
  Merges with @pending @story-{id} tags →
    Dev creates impl/{ticket-id} branch →
      [This conversational agent guides implementation] →
        Dev + AI iterate: plan → review → test → implement → refactor →
          All scenarios pass →
            @pending removed for this story →
              Implementation PR →
                CI validates →
                  Merges to main
```

## References

- **Tech Standards:** See `../../tech_standards.md` for IoC patterns
- **Architecture:** See `../../architecture.md` for system design
- **Testing:** See `../../testing.md` for test patterns
- **Standards Compliance:** See `../standards-compliance/` for automated code review
