# Developer Implementation Prompt

## Purpose

Guide developers in implementing step definitions for approved Gherkin specifications while following TDD and IoC patterns.

## Files in This Directory

- `prompt.md` - Implementation guide with TDD cycle, IoC patterns, examples
- `sample-spec.feature` - Example approved spec with `@pending` scenarios

## How It Works

1. **Input:** Approved spec (merged to main, `@pending` tags)
2. **Implementation:** Developer writes tests first, implements services using primary constructors
3. **Testing:** Run scenarios until passing
4. **Output:** Implementation PR merged, ready for demo/release
5. **Validation:** After business validates in demo/release, change `scenarios remain `@pending`

## Key Principles

**TDD Cycle:**
```
Write failing test → Implement → Test passes → Refactor
```

**IoC Pattern:**
- Primary constructor in tests (with mocks)
- Production factory for wiring (excluded from coverage)
- NO business logic in factories

**Coverage:**
- Unit tests: 80%+ (business logic)
- Production factories: Excluded (`// coverage:ignore`)

## AI vs. Human Responsibilities

### AI Generates Code (Following Patterns)

**Tests:**
- Unit test skeletons following TDD patterns
- Mock setup code following IoC patterns
- Boundary condition tests
- Table-driven test structures

**Implementation:**
- Services following IoC pattern (primary constructor + production factory)
- Repository implementations using GORM
- Standard CRUD operations
- Boilerplate and repetitive code

### Developer Applies Software Engineering (Making It Production-Ready)

**For Tests:**
- **Refactor:** Extract helpers, improve naming, eliminate duplication
- **Maintain:** Clear assertions, good error messages, readable test flow
- **Optimize:** Parallelize tests, reduce fixture overhead
- **Validate:** Meaningful coverage, not just line count
- **Ensure consistency:** Test style matches team conventions

**For Implementation (Critical - This is Engineering, Not Brick-Laying):**
- **Security Review:**
  - Input validation on all external data
  - Authorization checks (who can access what)
  - SQL injection prevention (GORM helps, but review raw queries)
  - Sensitive data handling (passwords, PII, tokens)
  - OWASP Top 10 vulnerability checks
- **Performance Review:**
  - Identify N+1 query problems (common with ORMs)
  - **Verify correct data structures:**
    - Maps for lookups (O(1)), not slice iteration (O(n))
    - Sets for unique collections, not slice + duplicate checking
    - Consider time complexity of operations
  - Review algorithm efficiency
  - Optimize hot paths and expensive operations
  - Add appropriate caching
  - Reduce unnecessary allocations
- **Consistency & Reuse Review:**
  - **Extract common logic:** Validation, business rules, transformations should exist once, be reused everywhere
  - **Avoid fragmented logic:** Same email validation everywhere (not different rules in different code paths)
  - **Prevent divergence:** Logic shouldn't vary by endpoint, UI path, or AI reimplementation
  - **Refactor duplicates:** AI often generates similar-but-different logic across sessions (context window exhaustion)
  - **Ensure unified behavior:** User gets same result regardless of which code path executes
  - **Unify error handling patterns**
  - **Maintain architectural consistency**
- **Architecture Review:**
  - Validate proper layer separation (domain/application/infrastructure)
  - Ensure dependencies point correctly (no domain depending on infrastructure)
  - Review abstractions for appropriateness
  - Validate API design
- **Maintainability Review:**
  - Simplify overly complex logic
  - Extract functions when needed
  - Improve clarity with better names
  - Add comments for non-obvious decisions
- **Correctness Review:**
  - Verify edge cases handled properly
  - Ensure error paths are comprehensive
  - Validate business logic matches requirements

**Key principle:** AI generates code that compiles, follows patterns, and passes tests. Developers ensure code is secure, performant, consistent, maintainable, and architected correctly. **This is software engineering, not brick-laying.**

**Complete workflow:**
1. AI generates unit test with mocks
2. Developer runs test, verifies it fails correctly
3. AI generates implementation following patterns
4. **Developer reviews implementation for engineering quality:**
   - Security: "Any vulnerabilities? Input validation?"
   - Performance: "N+1 queries? Using map for lookups or iterating slice? O(n²) when could be O(n)?"
   - Data structures: "Is this the right data structure for this operation?"
   - Consistency: "Is this logic duplicated elsewhere? Should it be extracted and reused?"
   - Fragmentation: "Could user get different results through different code paths?"
   - Architecture: "Right layer? Proper abstractions?"
5. **Developer refactors BOTH implementation and test code**
6. Developer validates coverage is meaningful and complete
7. Developer ensures entire solution is production-ready

## Example Workflow

```bash
# 1. Create implementation branch
git checkout -b impl/PROJ-1234

# 2. Write unit test (failing)
# Edit: internal/domain/services/password_reset_service_test.go

# 3. Implement service
# Edit: internal/domain/services/password_reset_service.go

# 4. Run unit tests
just test-unit

# 5. Implement step definition
# Edit: features/step_definitions/password_reset_steps.go

# 6. Run BDD scenarios
just test-bdd

# 7. Open PR
git push origin impl/PROJ-1234

# 8. After PR merges, demo/release to business

# 9. After business validates, update tags: scenarios remain @pending
# Edit: features/auth/password_reset.feature
```

## Integration with Other Prompts

```
requirements-analyst → Drafts spec →
  bo-review → Approves (merges with @pending) →
    developer-implementation → Implements and merges →
      Demo/release and business validation →
        scenarios remain @pending (after business validates)
```

## What Success Looks Like

- ✅ All scenarios pass
- ✅ Unit test coverage > 80%
- ✅ No business logic in production factories
- ✅ Tests use primary constructors with mocks
- ✅ All production factories marked `// coverage:ignore`
- ✅ Scenarios remain `@pending` after implementation

## When to Escalate

- Spec is technically infeasible → Create `amend/{ticket-id}` branch
- Missing edge cases discovered → Add scenarios, get BO re-approval
- API doesn't match spec → Document issue, propose solution

## References

- IoC Patterns: `../../architecture.md § Inversion of Control`
- Testing Standards: `../../testing.md`
- Tech Standards: `../../tech_standards.md`
- Personas: `../personas.md` (Developer role)
