# Business Owner Review Prompt

Review AI-drafted Gherkin specifications for business alignment before approval.

## Your Role

You are the Business Owner. Your job is to ensure AI-drafted specifications match business intent. You are NOT reviewing code quality or technical implementation - only business correctness.

## What to Review

### 1. Intent Match
Do the scenarios capture what the story requested?

**Check:**
- [ ] All acceptance criteria covered by scenarios
- [ ] Scenarios test the RIGHT behavior
- [ ] No scenarios testing WRONG behavior

### 2. Business Rules
Are business rules correctly applied?

**Check against business.md:**
- [ ] Compliance requirements addressed
- [ ] Rate limits correct
- [ ] Security policies followed
- [ ] Domain-specific rules applied

### 3. AI Assumptions (if flagged)
If AI marked MEDIUM confidence with assumptions:

**Review each assumption:**
- [ ] Assumption is acceptable → Approve
- [ ] Assumption is wrong → Request change
- [ ] Assumption reveals missing context → Update business.md

### 4. Completeness
Any obvious missing scenarios?

**Common gaps:**
- [ ] Missing error cases
- [ ] Missing edge cases user would encounter
- [ ] Missing business-critical validations

### 5. Boundary Conditions
Do boundary tests make business sense?

**Check:**
- [ ] Email/input validation boundaries reasonable
- [ ] Error messages match business tone
- [ ] Edge cases relevant to actual use

## What NOT to Review

❌ **Don't review:**
- Step definition implementation (developer's job)
- Technical feasibility (developer will flag)
- Test code quality (developer/QA job)
- Directory structure (enforced by standards)

✅ **Do review:**
- Business behavior described by scenarios
- Business rules application
- User-facing error messages
- Completeness from business perspective

## Decision Matrix

| Situation | Action |
|-----------|--------|
| Intent matches, rules correct | **Approve** |
| Minor wording issues | **Approve** with comment (dev can fix) |
| Wrong business logic | **Request changes** |
| Missing critical scenario | **Request changes** |
| AI assumption wrong | **Request changes** |
| Unsure about assumption | **Ask question** in PR comments |

## Review Checklist

Quick checklist for every spec PR:

- [ ] `@story-{id}` tag present and correct
- [ ] All acceptance criteria from story covered
- [ ] Business rules from business.md applied correctly
- [ ] If MEDIUM confidence: reviewed all assumptions
- [ ] Error messages are business-appropriate
- [ ] No obvious missing scenarios
- [ ] Boundary tests make business sense

## Approval Actions

### If Approved
1. Approve PR in GitHub/GitLab
2. PR merges automatically with `@pending` tags
3. Developer can now implement

### If Changes Needed
1. Leave comments on specific scenarios
2. Tag AI or PO to revise
3. Re-review when updated

### If Unsure
1. Ask questions in PR comments
2. Tag domain expert if needed
3. Discuss in team channel

## Common Review Scenarios

### Scenario 1: AI Generated Good Spec

```gherkin
@story-PROJ-123
Scenario: User submits valid order
  Given I am logged in as "customer@example.com"
  And I have items in my cart
  When I submit the order
  Then the order should be created
  And I should receive a confirmation email
```

**Review:**
- ✅ Covers acceptance criteria
- ✅ Matches business flow
- ✅ Confirmation email aligns with policy

**Action:** Approve

### Scenario 2: Missing Business Rule

```gherkin
@story-PROJ-456
Scenario: User places order
  When I place an order for $100
  Then the order should be created
```

**Review:**
- ❌ Missing: minimum order amount check (business rule BR-042: $25 minimum)
- ❌ Missing: payment method validation

**Action:** Request changes
**Comment:** "Missing BR-042 minimum order amount check and payment validation"

### Scenario 3: Wrong Error Message

```gherkin
Scenario: Invalid email format
  When I enter email "invalid"
  Then I should see "Error: bad email"
```

**Review:**
- ❌ Error message tone doesn't match brand (should be helpful, not terse)

**Action:** Request changes
**Comment:** "Error message should be: 'Please enter a valid email address'"

### Scenario 4: AI Assumption to Verify

**AI Note:** "Assumed email delivery is asynchronous - verify with BO"

**Review:**
- ✅ Correct - email is async via SendGrid

**Action:** Approve
**Comment:** "Confirmed: email delivery is asynchronous"

## Time Estimate

- **High confidence PR:** 5-10 minutes
- **Medium confidence PR:** 10-15 minutes (review assumptions)
- **Complex feature:** 15-20 minutes

Target: < 24 hours from PR creation to approval/feedback

## Output Format

### Approval Comment Template

```markdown
## Business Owner Review

✅ **APPROVED**

**Verified:**
- All acceptance criteria covered
- Business rules correctly applied
- No missing critical scenarios

**Notes:**
- [any clarifications or confirmations]

Merging to main. Ready for implementation.
```

### Request Changes Template

```markdown
## Business Owner Review

🔄 **CHANGES REQUESTED**

**Issues:**

1. **Missing scenario:** [description]
   - Required by: [business rule/acceptance criteria]
   - Should test: [expected behavior]

2. **Incorrect behavior:** [scenario reference]
   - Current: [what it tests]
   - Should be: [correct behavior]

3. **Wrong assumption:** [AI assumption]
   - Actual: [correct information]

**Once fixed:** Tag me for re-review
```

## Integration with Workflow

```
AI drafts specs → Opens PR in spec/ branch →
  BO reviews (this prompt) →
    If approved: PR merges with @pending tags →
      Developer implements
    If changes: AI/PO revises → BO re-reviews
```

Your approval is the gate between specification and implementation. Developers trust that approved specs represent correct business requirements.
