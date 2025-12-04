# Requirements Analyst Prompt

## Purpose

AI agent that analyzes user stories and drafts Gherkin specifications with boundary conditions and step definitions.

## Files in This Directory

- `prompt.md` - Minimal, token-efficient requirements analyst prompt
- `sample-story.md` - Example user story (password reset)
- `sample-context/` - Simplified context files for demonstration
  - `testing.md` - Step library and boundary patterns
  - `business.md` - Business rules and compliance

## How It Works

1. **Input:** User story with acceptance criteria
2. **Analysis:** AI reads context files, validates against API specs, identifies patterns
3. **Output:** Either:
   - **Confident:** Draft Gherkin + skeleton step definitions + PR
   - **Uncertain:** Escalation with specific questions
4. **After Escalation Resolution:** Codify answers into context files
   - Business rules → `business.md`
   - API behavior → `architecture.md`
   - Test patterns → `testing.md`
   - Tech standards → `tech_standards.md`
   - Creates feedback loop: fewer future escalations on similar issues

## Example Execution with Claude Subagent

### Scenario: Password Reset Feature

**Input:** `sample-story.md` (PROJ-1234)

**Execution captured below...**

---

## 🤖 ACTUAL SUBAGENT EXECUTION

The subagent successfully analyzed the story and produced a **MEDIUM confidence** specification with comprehensive boundary conditions and test scenarios.

### Key Highlights

**✅ What the AI Generated:**
- 11 test scenarios covering happy path, edge cases, boundaries, security, and concurrency
- Boundary condition scenarios for email validation and password complexity
- Rate limiting and token expiry scenarios
- 20+ step definitions (mix of reused and new)
- Complete PR description with assumptions flagged

**⚠️ Flagged Issues:**
- Missing architecture.md (SendGrid not documented)
- Missing API specifications for validation
- Made 6 assumptions requiring BO verification

**Result:** Ready for BO review with clear documentation of assumptions and missing context.

---

### Full Output

See the complete generated specification in `example-output/` directory:
- `password_reset.feature` - Complete Gherkin feature file
- `password_reset_steps.go` - Skeleton step definitions
- `pr-description.md` - Full PR description with metadata

### Output Summary

**Confidence:** MEDIUM (⚠️)
- Clear business rules from business.md
- Existing step library available for reuse
- Missing external API documentation (SendGrid)
- Missing internal API specifications

**11 Scenarios Generated:**
1. Happy path: successful reset request
2. Validation: email exists check
3. Happy path: complete password reset
4. Edge case: token expiry after 24 hours
5. Security: single-use token prevention
6. Security: rate limiting (5 requests/15 min)
7. Boundary: email validation (empty, invalid, too long)
8. Boundary: password complexity (length, uppercase, number, special char)
9. Edge case: token exactly at 24-hour boundary
10. Concurrency: simultaneous token usage attempts
11. Security: audit logging verification

**Boundary Conditions Applied:**
- Email: 5 boundary cases (empty, invalid formats, max length)
- Password: 5 complexity violations
- Time: Expiry at boundary and past boundary
- Concurrency: Simultaneous access handling

**Steps Reused from testing.md:**
- `Given a user exists with email` ✅
- `Then I should receive an email` ✅
- Authentication steps ✅

**New Steps Created:** 20+ password-reset-specific steps

---

