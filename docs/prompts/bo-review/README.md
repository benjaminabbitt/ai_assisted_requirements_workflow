# Business Owner Review Prompt

## Purpose

Guide Business Owners in reviewing AI-drafted Gherkin specifications for business correctness before approval.

## Files in This Directory

- `prompt.md` - BO review checklist and guidelines
- `sample-pr.md` - Example PR from requirements-analyst (PROJ-1234)

## How It Works

1. **AI drafts spec** → Opens PR in spec/ branch
2. **BO reviews** (using this prompt as guide)
3. **BO approves or requests changes**
4. **If approved** → PR merges with `@pending` tags, dev can implement

## What BO Reviews

✅ **Business correctness:**
- Do scenarios match story intent?
- Are business rules applied correctly?
- Are error messages appropriate?
- Any missing critical scenarios?

❌ **NOT reviewing:**
- Code quality (dev's job)
- Technical feasibility (dev will flag)
- Directory structure (automated)

## Example Execution with Claude Subagent

### Scenario: Review Password Reset PR

**Input:** `sample-pr.md` (AI-drafted spec for PROJ-1234)

**Executing BO review...**

---

## 🤖 ACTUAL SUBAGENT EXECUTION

The subagent performed a thorough business review and **REQUESTED CHANGES** with 7 critical issues identified.

### Key Findings

**🔴 Critical Issues (7):**
1. **Security vulnerability** - Email enumeration attack (reveals registered emails)
2. **Missing business rule** - Admin password requires 12+ chars (BR-003 exception)
3. **Missing validation** - Email content/link not verified
4. **UX issue** - Error messages shown one-at-a-time (poor experience)
5. **Ambiguous logic** - Rate limiting scenario unclear (5th or 6th request fails?)
6. **Missing scenario** - Multiple active tokens not addressed
7. **Inconsistency** - Async email delivery assumption needs clarification

**✅ What Worked Well:**
- Comprehensive security coverage (expiry, single-use, rate limiting)
- Consistent audit logging across all scenarios
- Good use of tags and organization
- Boundary condition testing present

### Decision: CHANGES REQUESTED

**Rationale:** Critical business correctness issues affecting security, compliance, and UX. Foundation is solid, but targeted corrections needed.

### Full Review Output

See complete BO review with:
- Detailed issue descriptions
- Specific code changes required
- Questions for clarification
- Positive observations

The review demonstrates that AI catches most requirements but BO review is essential for:
- Security implications (email enumeration)
- Business rule exceptions (admin passwords)
- UX considerations (error message presentation)
- Real-world edge cases (multiple active tokens)

**Time taken:** ~15 minutes for thorough review (MEDIUM confidence PR)

---

