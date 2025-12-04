# Roles & Responsibilities

This document defines the roles in the AI-augmented requirements workflow, their responsibilities, required access, and how they interact with the system.

---

## Overview

The AI-augmented workflow distributes work across specialized roles, each with specific responsibilities and access requirements. **Humans make decisions, AI implements them, humans verify.**

---

## 1. Product Owner (PO)

### Primary Responsibility
Write user stories with clear acceptance criteria in the ticketing system.

### What They Decide
- **What to build:** Features, priorities, business value
- **Success criteria:** What "done" means for each story
- **User experience:** How features should behave from user perspective

### Required Access
- **Ticketing System:** Full access (create, edit, update, close tickets)
- **Source Control:** Read-only access to context files (via MCP or web interface)
- **Context Files:** Reference `business.md` for domain terms, personas, business rules

### Workflow Interaction
1. **Create story ticket** with:
   - Clear title and description
   - User story format: "As a [persona], I want [action], so that [benefit]"
   - Acceptance criteria as checklist or Given/When/Then statements
   - Labels, component, priority
2. **Trigger AI analysis** (automatic via webhook or manual with label/comment)
3. **Monitor for escalations** - AI may post questions to ticket if context is missing
4. **Handoff to BO** - Once spec PR is created, BO takes over review

### Training Requirements
- **Duration:** 2 hours
- **Topics:**
  - Writing effective user stories
  - Clear acceptance criteria
  - Using ticketing system effectively
  - When to use requirements-drafting-assistant for complex features
  - How to interpret AI escalations

### Success Metrics
- Stories have clear, testable acceptance criteria
- Low AI escalation rate on well-written stories
- BO can review specs without needing PO clarification

---

## 2. Business Owner (BO)

### Primary Responsibility
Review and approve all specification changes. Ensure AI-drafted Gherkin scenarios match business intent.

### What They Decide
- **Specification correctness:** Does this spec match what we want to build?
- **Business logic:** Are business rules applied correctly?
- **Completeness:** Are critical scenarios missing?
- **User experience validation:** Are error messages appropriate?

### What They DO NOT Decide
- Implementation details (that's Developer responsibility)
- Technical architecture (that's Tech Lead responsibility)
- Testing approach (that's QA Lead responsibility)

### Required Access

**Ticketing System (via MCP):**
- Read tickets, comments, related tickets
- Write comments (for clarifications, questions)
- No need to update ticket status (optional)

**Source Control (via MCP - Read-Only):**
- Read `.feature` files (specifications)
- Read context files (`business.md`, `architecture.md`, `testing.md`, `tech_standards.md`)
- **No access to code files** (security: no credentials, secrets, implementation)
- **No direct write access** - all changes via PR approval

**Source Control (Write via PR Approval):**
- Approve/reject PRs for `.feature` file changes (enforced via CODEOWNERS)
- Request changes with comments
- Approve amendments to specifications

**Context Files:**
- **Own:** `business.md` (business rules, personas, compliance)
- **Review:** Other context files to understand constraints

### Workflow Interaction

**For new specifications:**
1. **Receive PR notification** (spec/PROJ-1234) with AI-drafted Gherkin scenarios
2. **Read AI analysis** in PR description (confidence level, assumptions, boundary conditions)
3. **Review scenarios** using bo-review agent guidance or manually:
   - Do scenarios match story intent?
   - Are business rules applied correctly?
   - Are error messages appropriate?
   - Are critical scenarios missing?
4. **Approve or request changes:**
   - ✅ **Approve:** Spec merges to main with `@pending` tags, implementation can begin
   - ⚠️ **Request changes:** Add comments, AI or human updates spec, re-review
5. **Update business.md** if AI escalated due to missing business rules

**For specification amendments (during implementation):**
1. Developer discovers edge case during implementation
2. Developer creates amendment PR (amend/PROJ-1234) with new scenarios
3. BO reviews amendment using same process
4. BO approves → Developer can pull amendment into their impl branch and continue

### Training Requirements
- **Duration:** 1 hour
- **Topics:**
  - Gherkin basics (Given/When/Then, Scenario, Scenario Outline)
  - How to review specs for business correctness
  - Using bo-review agent for guidance
  - How to request changes in PR
  - When to update business.md based on escalations

### Success Metrics
- Spec approval time < 15 minutes (async, no meetings required)
- Low rate of spec amendments during implementation (means specs were complete)
- Developers don't need BO clarification during implementation

### Access Requirements (MCP)
See [MCP Integration Requirements](mcp-integration-requirements.md#business-owner-role) for technical details.

---

## 3. Developer

### Primary Responsibility
Implement step definitions and service code following TDD. Ensure code quality, security, performance, and maintainability.

### What They Decide
- **Implementation approach:** How to implement the behavior
- **Code architecture:** Class design, abstractions, patterns
- **Technical tradeoffs:** Performance vs readability, simple vs extensible
- **Refactoring:** When and how to improve code structure

### What They DO NOT Decide
- **Requirements:** What to build (that's PO/BO responsibility)
- **Technical standards:** Coding patterns, IoC approach (that's Tech Lead responsibility)
- **Testing patterns:** Boundary conditions, step library structure (that's QA Lead responsibility)

### Required Access

**Full Repository Access:**
- Read and write all code files
- Create branches, commits, PRs
- Merge PRs after approval

**Ticketing System:**
- Read tickets
- Update ticket status
- Add comments

**Context Files:**
- **Read:** All context files to understand domain, architecture, patterns
- **Write:** Via PR (not primary owner, but can suggest updates)

**CI/CD:**
- Trigger builds locally and in CI
- View test results, logs
- Debug failures

### Workflow Interaction

**Implementation flow:**
1. **After spec approval:** Spec is merged to main with `@pending @story-PROJ-1234` tags
2. **Create implementation branch:** `impl/PROJ-1234`
3. **Collaborate with developer-implementation agent:**
   - **Plan:** Agent reads tech_standards.md, architecture.md, introspects dependencies
   - **Review:** Developer validates plan, provides feedback
   - **Execute:** Iterative TDD implementation with agent
4. **Follow TDD cycle:**
   - Write test using primary constructor + mocks
   - Implement minimal code to pass
   - Refactor for quality
5. **Apply software engineering:**
   - Review for security (injection, auth, data exposure)
   - Review for performance (N+1 queries, algorithm efficiency, data structure choice)
   - Review for consistency (reuse logic, don't fragment)
   - Review for maintainability (clear abstractions, readable code)
6. **Run scenarios:** `godog run -t @story-PROJ-1234`
7. **When passing:** Remove `@pending` tags for this story (keep `@story-PROJ-1234`)
8. **Run standards-compliance:** Fix violations before opening PR
9. **Open implementation PR:** impl/PROJ-1234 → main
10. **CI validates:** Blocks merge if `@pending` tags present for this story
11. **Merge:** Implementation complete

**Amendment flow (edge case discovered during implementation):**
1. **Discover edge case** not covered in spec
2. **Create amendment branch:** `amend/PROJ-1234`
3. **Add scenarios** with `@pending @story-PROJ-1234` tags
4. **Create amendment PR** → main
5. **Tag BO for review**
6. **Meanwhile, pull amendment into impl branch:** `git pull origin amend/PROJ-1234`
7. **Continue implementation** for new scenarios
8. **Wait for BO approval** before merging impl PR (depends on amendment merge)

### Training Requirements
- **Duration:** 4 hours
- **Topics:**
  - IoC patterns (primary constructor + production factory)
  - Godog/BDD testing framework
  - TDD cycle and red-green-refactor
  - Using developer-implementation agent effectively
  - Running standards-compliance before PR
  - Understanding tech_standards.md patterns

### Success Metrics
- Implementation time: Target <2 hours per story
- Standards compliance score: >95%
- Low rate of defects found in code review
- Scenarios pass on first CI run (TDD worked)

---

## 4. QA Lead

### Primary Responsibility
Maintain testing patterns, boundary conditions, and step library in `testing.md`. Ensure consistent edge case coverage.

### What They Decide
- **Testing patterns:** What boundary conditions to test for each data type
- **Step library:** Reusable Gherkin steps and their signatures
- **Edge case coverage:** What scenarios should be generated for common patterns
- **Fuzzing patterns:** What fuzzing libraries to reference for comprehensive testing

### What They DO NOT Decide
- **Business rules:** What the system should do (that's BO responsibility)
- **Implementation:** How tests are implemented (that's Developer responsibility)
- **Technical standards:** Testing framework, IoC patterns (that's Tech Lead responsibility)

### Required Access

**Source Control:**
- Read/write `testing.md` (owns this file)
- Read `.feature` files to ensure step reuse
- Read step definition code to validate implementation

**Context Files:**
- **Own:** `testing.md` (step library, boundary patterns, fuzzing refs)
- **Read:** `architecture.md` (API constraints inform boundary conditions)

### Workflow Interaction

**Pattern maintenance:**
1. **Review AI-drafted specs** for consistent edge case coverage
2. **Update testing.md** when new patterns emerge:
   - New data types requiring boundary conditions
   - New reusable steps developers create
   - New fuzzing libraries or tools
3. **Monitor AI escalations** - if AI doesn't know what edge cases to test, add pattern to testing.md
4. **Collaborate with developers** on step definition implementation

**Quality assurance:**
1. **Review implementation PRs** to ensure:
   - Step definitions follow patterns in testing.md
   - Boundary conditions are tested
   - Edge cases are covered
2. **Maintain step library** - refactor duplicated steps, consolidate patterns
3. **Track coverage** - ensure critical paths have scenario coverage

### Training Requirements
- **Duration:** 4 hours
- **Topics:**
  - Gherkin step library design
  - Boundary condition patterns (strings, numbers, dates, arrays)
  - Fuzzing approaches and libraries
  - Structure and governance of testing.md
  - Reviewing specs for edge case completeness

### Success Metrics
- Low variance in edge case coverage across stories (patterns are reused)
- High step reuse rate (new steps only when truly needed)
- AI-generated specs include comprehensive boundary conditions (testing.md is effective)

### Context File Ownership
Assigned via CODEOWNERS:
```
# CODEOWNERS
testing.md @qa-lead
```

---

## 5. Tech Lead

### Primary Responsibility
Maintain technical standards, architecture constraints, and ensure code quality across the team.

### What They Decide
- **Technical standards:** Language conventions, IoC patterns, directory structure
- **Architecture constraints:** External dependencies, system boundaries, API contracts
- **Code quality gates:** What violations block PRs, compliance thresholds
- **Tooling choices:** Testing frameworks, CI/CD setup, linters

### What They DO NOT Decide
- **Business requirements:** What to build (that's PO/BO responsibility)
- **Individual implementation:** How developers implement specific features (that's Developer responsibility)

### Required Access

**Source Control:**
- Full access (read/write all files)
- Code review rights on all PRs
- Can block PRs that violate standards

**Context Files:**
- **Own:** `tech_standards.md`, `architecture.md`
- **Review:** All context files for consistency

**CI/CD:**
- Configure pipelines, gates, compliance checks
- Integrate standards-compliance agent
- Set up branch protection rules

### Workflow Interaction

**Standards maintenance:**
1. **Define patterns** in `tech_standards.md`:
   - IoC patterns (primary constructor + production factory)
   - Directory structure
   - Naming conventions
   - Coverage strategy
2. **Update architecture.md** when:
   - New external dependencies added
   - API contracts change
   - System constraints change (rate limits, quotas)
3. **Review standards-compliance reports** in PRs
4. **Enforce compliance gates** in CI

**Code quality oversight:**
1. **Review implementation PRs** for architectural concerns:
   - Does this fit the system design?
   - Are abstractions appropriate?
   - Are new dependencies justified?
2. **Monitor escalations** - if AI escalates due to missing architecture docs, update architecture.md
3. **Codify escalation resolutions:**
   - Technical constraint clarifications → `tech_standards.md`
   - API behavior patterns → `architecture.md`
   - Security/performance patterns → relevant context file

**Context file governance:**
1. **Track escalation patterns** - which context files are frequently missing information?
2. **Version context files** - bump version when significant changes made
3. **Monitor context file freshness** - ensure docs stay current

### Training Requirements
- **Duration:** Minimal (defines the standards)
- **Topics:**
  - How to structure tech_standards.md effectively
  - Using standards-compliance agent in CI
  - Monitoring escalation patterns to improve context files

### Success Metrics
- Low AI escalation rate (context files are complete)
- High standards compliance scores across PRs (>95%)
- Consistent code patterns across team (standards are clear)
- Context files stay current (no stale documentation)

### Context File Ownership
Assigned via CODEOWNERS:
```
# CODEOWNERS
tech_standards.md @tech-lead
architecture.md @tech-lead
```

---

## 6. AI Requirements Agents (System Role)

### Requirements-Drafting-Assistant

**Purpose:** Help Business Owners articulate requirements through conversational exploration.

**Mode:** Conversational

**Required Access (via MCP):**
- Ticketing system: Read tickets, comments, related tickets, search
- Source control: Read `.feature` files, context files (no code access)

**Credentials:** Service account (shared with requirements-analyst)

**Use cases:**
- Explore complex or vague requirements with BO
- Pull historical context from tickets
- Reference existing feature files for consistency
- Identify edge cases and dependencies

### Requirements-Analyst

**Purpose:** Analyze story tickets and draft Gherkin specifications automatically.

**Mode:** One-shot

**Required Access (via MCP):**
- Ticketing system: Read tickets, comments, related tickets, search
- Source control: Read `.feature` files, context files, API specs (no code access)
- Source control (via git API): Create branches, commit files, open PRs

**Credentials:** Service account (shared with requirements-drafting-assistant for MCP; separate git credentials for PR creation)

**Use cases:**
- Fetch ticket data to draft specifications
- Search related tickets for context
- Read existing feature files for conventions
- Create spec branches and PRs

---

## Role Interactions

### Story Creation → Specification

```
PO creates story
    ↓
AI (requirements-analyst) reads ticket via MCP
    ↓
AI reads existing .feature files via MCP
    ↓
AI drafts spec, creates PR
    ↓
BO reviews spec (reads .feature via MCP)
    ↓
BO approves → Spec merges with @pending tags
```

### Specification → Implementation

```
Developer creates impl branch
    ↓
Developer collaborates with developer-implementation agent
    ↓
Developer implements using TDD
    ↓
Developer runs standards-compliance
    ↓
Developer opens PR
    ↓
Tech Lead reviews for architecture/standards
    ↓
CI validates @pending tags removed
    ↓
Implementation merges
```

### Escalation → Context Improvement

```
AI escalates with question
    ↓
Expert (BO/Tech Lead/QA Lead) answers
    ↓
Expert updates appropriate context file
    ↓
Context file owner approves via CODEOWNERS
    ↓
Future similar stories don't escalate
```

---

## Access Control Matrix

| System | PO | BO | Developer | QA Lead | Tech Lead | AI Agent |
|--------|----|----|-----------|---------|-----------|----------|
| **Ticketing** |
| Read tickets | ✅ Full | ✅ MCP | ✅ Full | ✅ Full | ✅ Full | ✅ MCP (service) |
| Write tickets | ✅ | ✅ Comment | ✅ | ✅ | ✅ | ❌ |
| **Source Control** |
| Read .feature files | ✅ Web/MCP | ✅ MCP | ✅ Full git | ✅ Full git | ✅ Full git | ✅ MCP (read-only) |
| Write .feature files | ❌ | ✅ PR approval | ✅ PR/commit | ✅ PR/commit | ✅ PR/commit | ✅ PR creation (git API) |
| Read code files | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ |
| Write code files | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ |
| Read context files | ✅ Web/MCP | ✅ MCP | ✅ Full git | ✅ Full git | ✅ Full git | ✅ MCP |
| Write context files | ❌ | ✅ business.md (CODEOWNERS) | ✅ PR | ✅ testing.md (CODEOWNERS) | ✅ tech/arch.md (CODEOWNERS) | ❌ |
| **CI/CD** |
| View results | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Trigger builds | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ (via PR) |
| Configure pipelines | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |

**Key Insight:** Business Owners need read access to specifications and tickets, but should not have full repository access (security: no code, credentials, secrets).

---

## CODEOWNERS Configuration

**Enforce approval requirements for context files:**

```
# Context files - specialized owners
business.md @business-owner
testing.md @qa-lead
tech_standards.md @tech-lead
architecture.md @tech-lead

# Feature files - BO must approve all changes
features/**/*.feature @business-owner

# Code files - standard review process
internal/**/*.go @tech-lead
features/step_definitions/**/*.go @tech-lead
```

**Why CODEOWNERS:**
- Enforces BO approval for all specification changes
- Ensures context file changes are reviewed by experts
- Provides audit trail of who approved what
- Prevents accidental spec changes without BO knowledge

---

## Time Investment Comparison

### Traditional Process (per story)
- **PO:** 30 min (write story) + 1 hour (Three Amigos meeting) = 1.5 hours
- **BO:** 1 hour (Three Amigos meeting) + 1 hour (manual Gherkin writing) = 2 hours
- **Developer:** 1 hour (Three Amigos meeting) + 2 hours (implementation) = 3 hours
- **QA:** 1 hour (Three Amigos meeting) = 1 hour
- **Total:** 7.5 hours

### AI-Augmented Process (per story)
- **PO:** 15 min (write story)
- **BO:** 10 min (async PR review)
- **Developer:** 2 hours (implementation with AI collaboration)
- **QA Lead:** 5 min (pattern maintenance amortized across stories)
- **Tech Lead:** 5 min (review standards-compliance report)
- **Total:** 2.5 hours

**Savings:** 5 hours per story (67% reduction)

**Key difference:** Async PR review replaces synchronous meetings. AI drafts specs, humans verify.

---

## Success Criteria

**Roles are working effectively when:**
- ✅ PO can write stories without needing to attend meetings
- ✅ BO can approve specs asynchronously in <15 minutes
- ✅ Developers implement without waiting for clarification
- ✅ QA patterns are reused across all stories
- ✅ Tech standards are consistently followed
- ✅ AI escalation rate decreases over time (context files improve)
- ✅ No one has access they don't need (security principle of least privilege)
