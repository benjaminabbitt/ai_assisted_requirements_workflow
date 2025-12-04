# Workflow Stages

This document describes the complete end-to-end workflow from story creation through implementation, including all decision points and actor responsibilities.

---

## Overview

The AI-augmented requirements workflow has four main stages:

1. **Story Creation** (Product Owner) - Write story with acceptance criteria
2. **AI Analysis & Drafting** (AI Agent) - Analyze context, draft Gherkin, or escalate
3. **BO Review** (Business Owner) - Review and approve specification
4. **Implementation** (Developer) - Implement step definitions following TDD

**Optional Stage 0:** Requirements drafting assistance for complex features (conversational exploration with BO).

---

## Stage 0 (Optional): Requirements Drafting Assistance (BO: ~15-30 minutes, conversational)

**Agent:** `requirements-drafting-assistant` (Conversational)
**Prompt:** [`docs/prompts/requirements-drafting-assistant/`](../prompts/requirements-drafting-assistant/)

**NEW:** For complex or vague requirements, a conversational AI agent helps Business Owners articulate requirements before formal spec creation.

**When to use:**
- Requirement is vague or exploratory ("we need better security")
- Complex feature with many edge cases
- Need to validate against API constraints and business rules
- BO wants help thinking through scenarios

**How it works:**
1. BO describes requirement to conversational agent
2. Agent asks clarifying questions (references context files, API contracts)
3. Agent explores edge cases, security implications, technical constraints
4. Conversation builds up complete requirement incrementally
5. Agent produces structured requirement document

**Output:**
- Structured requirement with user story, acceptance criteria, edge cases, API dependencies
- Ready for ticket creation (PROJ-XXX)
- Feeds into Stage 1 (Story Creation)

**Key difference:** This is CONVERSATIONAL (back-and-forth dialogue) while the requirements-analyst agent (Stage 2) is ONE-SHOT (reads ticket, drafts spec).

---

## Stage 1: Story Creation (PO: ~10 minutes)

Product Owner creates a ticket with this structure:

```markdown
## User Story
As a [persona], I want [capability], so that [benefit].

## Acceptance Criteria
- [ ] [Verifiable outcome]
- [ ] [Verifiable outcome]

## Notes (optional)
[Any additional context]
```

**That's it.** No Gherkin drafting. No scheduling a Three Amigos meeting. No waiting for Dev and QA availability. AI handles the analysis and drafting.

**Branch name derived from ticket URL.** AI extracts the ticket ID (e.g., `PROJ-1234`) and uses it for branch naming (`spec/PROJ-1234` for specifications). No manual entry needed.

**Trigger:** Ticket labeled "ready-for-spec" → AI creates branch and begins review immediately.

---

## Stage 2: AI Review & Drafting (AI: seconds to minutes)

**Agent:** `requirements-analyst` (One-Shot)
**Prompt:** [`docs/prompts/requirements-analyst/`](../prompts/requirements-analyst/)

AI analyzes the story against all context files, package management, and API interfaces—doing in seconds what would take a human team an hour of discussion:

| Source | AI Asks... |
|--------|------------|
| Package management (`go.mod`, etc.) | What internal services does this depend on? |
| Internal API specs | What parameters are required? What are valid values? Any constraints? |
| `architecture.md` | What external dependencies? Any third-party APIs involved? |
| External API specs | What does the third-party API require? |
| `testing.md` | What steps already exist? What edge cases are typical for this pattern? |
| `business.md` | What business rules apply? Any compliance concerns? |
| `tech_standards.md` | How should step definitions be structured? Where do files go? |

**AI validates against actual contracts.** When reading API interfaces (internal or external), AI identifies missing required parameters, invalid values, and constraint violations—flagging them before the story reaches development.

**AI generates boundary conditions.** From API specs, AI derives boundary test cases: min/max values, empty inputs, null handling, format edge cases. QA defines the patterns in `testing.md`; AI applies them to each story's specific parameters.

**Fuzzing libraries augment edge case discovery.** Reference fuzzing libraries in `testing.md` (e.g., for string formats, numeric ranges, date handling). AI uses these to generate additional edge case scenarios beyond simple boundary testing.

Based on this analysis, AI determines confidence level and either drafts Gherkin or escalates.

### If Confident

AI creates:

- `.feature` file with scenarios tagged `@pending @story-{id}` (ready for implementation after merge)
- Skeleton step definitions (method signatures, no implementation)
- PR ready for BO review with confidence notes

### If Uncertain

AI posts to ticket:

- What it's uncertain about
- Specific questions for humans
- What context is missing

---

## Stage 3: BO Review (BO: ~5-15 minutes, async)

**Agent:** `bo-review` (One-Shot) - Optional guide for BO reviewers
**Prompt:** [`docs/prompts/bo-review/`](../prompts/bo-review/)

**Trigger:** Runs on Gherkin/Cucumber PRs (spec/ branches containing .feature files)

The PR is where Business Owner verifies AI's work. This is the quality gate—but it's a **review** gate, not a **creation** gate. The hard work is done.

**Compare to traditional:** BO would attend a meeting, discuss requirements, then wait for someone to write Gherkin, then review. Here, BO reviews a complete draft at their convenience.

### PR Contents

- Feature file with scenarios
- Skeleton step definitions
- Context versions used (for traceability)
- AI confidence notes (if any inferences were made)
- API validation results
- Link to original ticket

### Reviewer Checklist

| Check | Looking For |
|-------|-------------|
| Intent match | Do scenarios capture what the story asks for? |
| Inference review | Are AI's assumptions acceptable? |
| Gap check | Any obvious missing scenarios? |
| Tag check | Is `@story-{id}` present? |

### Approval Requirements

| AI Confidence | Approvals Needed | Rationale |
|---------------|------------------|-----------|
| High (no flags) | 1 (must include BO) | Routine work, BO validates business intent |
| Medium (flagged inferences) | 1 (must include BO) | Verify assumptions are acceptable |
| After escalation | 2 (must include BO) | Higher scrutiny for complex cases |

**BO approval is mandatory.** Business Owner must approve all feature file changes to ensure acceptance criteria align with business requirements.

**On BO approval:** PR merges to main. Scenarios remain `@pending`. Spec work is complete. Developers can now begin implementation.

**Note:** For meta-instructions on running these agents with Claude and logging requirements, see `../CLAUDE.md`.

---

## Stage 4: Implementation

**Agent:** `developer-implementation` (Conversational) - Interactive implementation collaboration
**Prompt:** [`docs/prompts/developer-implementation/`](../prompts/developer-implementation/)

**Automated Check:** `standards-compliance` - Runs on PR create and code updates
**Prompt:** [`docs/prompts/standards-compliance/`](../prompts/standards-compliance/)

After spec PR merges, developer creates a new implementation branch:

1. Create branch `impl/{ticket-id}` from main (which now has the `@pending @story-{id}` scenarios)
2. Implement step definitions following TDD and IoC patterns
3. Run scenarios until passing
4. **Remove `@pending` tags** from scenarios for this story (keep `@story-{id}`)
5. Open implementation PR
   - **CI blocks merge if `@pending` tags present for this story**
   - **CI automatically runs `standards-compliance` agent** to check IoC patterns, tech standards adherence
   - Reviews code for compliance violations before human review
6. Merge after CI passes (including @pending check and standards compliance)
7. Demo/release to business for validation (tracked in deployment/release process)

**If implementation reveals spec issues:** Create an `amend/{ticket-id}` branch to fix the spec. BO must approve spec changes before they merge.

---

## Step Definition Lifecycle

This section addresses how step definitions evolve from AI-generated skeletons through implementation and maintenance.

### Skeleton Generation

When AI creates a PR, it generates skeleton step definitions based on `tech_standards.md` patterns:

```go
// Generated by AI - PROJ-1234
// TODO: Implement
func (s *AuthSteps) theUserEntersValidCredentials() error {
    panic("not implemented")
}
```

### Implementation

Developer implements the skeleton, potentially adjusting the signature:

| Change Type | Action Required | Who Does It |
|-------------|-----------------|-------------|
| Parameter added | Update feature file step text | Developer in same PR |
| Step renamed | Update feature file step text | Developer in same PR |
| Step split into multiple | Update feature file, may need new scenarios | Developer, BO re-approves |
| Step consolidated | Update all affected feature files | Developer, BO re-approves |

### Signature Change Protocol

When implementation requires changing a step signature:

- **Same PR rule:** Step definition changes and feature file updates must be in the same PR
- **BO re-approval:** If feature file changes, BO must approve again
- **Backward compatibility:** If step is used in multiple features, keep old signature as deprecated alias
- **Documentation:** Update `testing.md` with new step signature

### Shared Step Libraries

For enterprises with multiple teams sharing steps:

#### Structure

```
/shared-steps/
├── common/           # Cross-domain steps (login, navigation)
├── domain-a/         # Domain-specific shared steps
└── domain-b/
```

#### Governance

- **Shared steps require 2 approvals** from different teams
- **Breaking changes require deprecation period** (2 sprints minimum)
- **Shared step changes trigger CI across all consuming repos**

---

## AI Confidence & Escalation

### Confidence Levels

| Level | Symbol | Meaning | AI Action |
|-------|--------|---------|-----------|
| High | ✓ | All info in context files, existing patterns | Drafts PR directly |
| Medium | ⚠ | Made reasonable inferences from context | Drafts PR with flagged assumptions |
| Low | ✗ | Missing context or high uncertainty | Escalates, does not draft |

### Escalation Triggers

AI escalates when it cannot find information from code, context files, or API interfaces:

| Situation | Why AI Escalates |
|-----------|------------------|
| External dependency not in `architecture.md` | Can't assess third-party API constraints |
| API interface spec missing or inaccessible | Can't validate parameters and constraints |
| Required parameters unclear in spec | Can't determine valid scenarios |
| Security/auth changes | Requires explicit human judgment |
| No similar patterns in `testing.md` | Can't identify edge cases |
| Conflicting requirements | Can't resolve ambiguity |
| Compliance implications unclear | Risk too high to assume |

### Escalation Format

When AI escalates, it posts to the ticket:

```markdown
## AI Review - Escalation Required

**Confidence:** LOW - Cannot proceed without human input

**Blocking Questions:**
1. [Specific question]
2. [Specific question]

**Missing Context:**
- [What's not documented in context files]
- [Recommended context file to update after resolution]

**Recommended Action:**
After answering these questions, consider codifying answers into context files:
- Business rule clarifications → business.md
- External API behavior → architecture.md
- Edge case patterns → testing.md
- Technical standards → tech_standards.md

**Once resolved:** Answer questions in ticket, AI will re-review.
```

**Critical: Codify escalation answers into context.** When AI escalates due to missing information, the human answers should typically be added to the appropriate context file. This is not always necessary (one-off exceptions, story-specific details), but frequently is (general patterns, API behaviors, business rules).

This creates a feedback loop that reduces future escalations:
1. AI escalates with missing context
2. Experts answer (applying judgment to difficult problem)
3. Answers codified into context file
4. Future similar stories don't escalate—AI has the information

**Who updates context files after escalation:**
- Business rules → Product Owner updates `business.md`
- External APIs, technical constraints → Architect/Tech Lead updates `architecture.md`
- Test patterns, edge cases → QA Lead updates `testing.md`
- Technical standards → Tech Lead updates `tech_standards.md`

### Conflict Resolution Process

When business rules conflict with technical constraints:

- **AI documents both constraints** in escalation
- **Product Owner and Tech Lead** must both respond
- **Resolution recorded** in ticket and relevant context file
- **AI re-analyzes** with updated context

**After resolution:** The conflict resolution decision should be documented in the appropriate context file so AI (and future humans) understand how to handle similar situations.

---

## Workflow Diagram

**The key: AI handles the toil, humans provide judgment. Requirements keep pace with AI-assisted development.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│   PO writes story + acceptance criteria          ← Human: 10 min        │
│   (ticket ID in URL becomes branch name)                                │
│                           │                                             │
│                           ▼                                             │
│   AI creates spec/{ticket-id} branch                                    │
│   AI reviews context files + reads API interfaces ← AI: seconds         │
│   AI validates params, generates boundary cases                         │
│                           │                                             │
│              ┌────────────┴────────────┐                                │
│              ▼                         ▼                                │
│         CONFIDENT                  UNCERTAIN                            │
│       (most stories)            (some stories)                          │
│              │                         │                                │
│              ▼                         ▼                                │
│     AI drafts Gherkin            AI escalates with                      │
│     + boundary scenarios         specific questions   ← AI: seconds     │
│     Opens PR for BO review       (missing specs, etc.)                  │
│              │                         │                                │
│              ▼                         ▼                                │
│     BO approves PR ←───────────── Humans discuss,     ← Human: 5-15 min │
│              │                     then AI drafts       (async)         │
│              ▼                                                          │
│     SPEC PR MERGES TO MAIN (with @pending tags)                         │
│              │                                                          │
│              ▼                                                          │
│     Dev creates impl/{ticket-id} branch                                 │
│     Implements step definitions                                         │
│     IMPL PR MERGES TO MAIN                                              │
│              │                                                          │
│              ▼                                                          │
│     Demo/release to business                                            │
│     Business validates (tracked in deployment)                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

TAGS ARE THE SOURCE OF TRUTH:
  @pending = awaiting implementation (removed when implementation merges)
  CI blocks implementation PR if @pending present for that story

TRADITIONAL: 30-60 min meeting + writing time + calendar coordination
THIS WORKFLOW: 5-15 min async review (for confident stories)

BONUS:
- AI catches missing required params, type mismatches, constraint violations
- AI generates boundary conditions from API specs + QA patterns
- Fuzzing library references enable thorough edge case coverage
- Issues surface BEFORE development, not during implementation

Speed parity achieved: Requirements no longer bottleneck AI-assisted development.
```

---

## Speed Targets

| Path | Time to PR Ready | Human Time Required |
|------|------------------|---------------------|
| AI Confident | Same day | ~5 min review |
| AI Confident + Inferences | Same day | ~15 min review |
| AI Escalates | 2-3 days | Discussion required |

**Goal:** Most stories need only a quick PR review—no meetings, no back-and-forth, no waiting for schedules to align. The proportion depends on how complete your context files and API specs are.

**Compare to traditional:** A Three Amigos session takes 30-60 minutes per story, requires synchronous availability of 3+ people, and still needs someone to write the Gherkin afterward. This workflow reduces human time per story from hours to minutes.

---

## Next Steps

- **For story creation:** See [Scenario Design Guide](scenario-design.md)
- **For context file maintenance:** See [Context Files Guide](context-files.md)
- **For implementation:** See [Developer Guide](developer-implementation.md) prompt
- **For metrics:** See [Metrics Guide](metrics.md)
