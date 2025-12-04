# Enterprise Requirements Workflow

## Gherkin/Cucumber with AI Augmentation

**Version 2.0**

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Roles & Responsibilities](#2-roles--responsibilities)
3. [Integration Architecture](#3-integration-architecture)
4. [Source Control Strategy](#4-source-control-strategy)
5. [Context Files](#5-context-files)
6. [Workflow Stages](#6-workflow-stages)
7. [Step Definition Lifecycle](#7-step-definition-lifecycle)
8. [AI Confidence & Escalation](#8-ai-confidence--escalation)
9. [Scenario Design Guidance](#9-scenario-design-guidance)
10. [SLA Enforcement](#10-sla-enforcement)
11. [Rollback & Revert Procedures](#11-rollback--revert-procedures)
12. [Tag Taxonomy](#12-tag-taxonomy)
13. [Hotfix Path](#13-hotfix-path)
14. [Metrics](#14-metrics)
15. [Naming Conventions](#15-naming-conventions)
16. [CI Configuration](#16-ci-configuration)
17. [Appendix: AI Agent Prompt](#17-appendix-ai-agent-prompt)
18. [Appendix: Context File Templates](#18-appendix-context-file-templates)

---

## 1. Executive Summary

### 1.0 The Problem: Requirements Can't Keep Up

**AI-assisted development has broken the requirements bottleneck in the wrong direction.**

When developers use AI correctly—with architects and senior devs reviewing output, enforcing patterns, and bringing code into compliance—implementation velocity increases dramatically. A competent team with AI assistance can implement faster than traditional requirements processes can specify.

The result: requirements become the bottleneck. Teams either:

- **Slow down development** to wait for specifications (waste)
- **Skip specifications** and accumulate technical/business debt (risk)
- **Write specifications after implementation** which defeats the purpose (theater)

None of these are acceptable. The solution is to apply the same AI acceleration to requirements that we apply to development.

### 1.1 The Solution: AI-Augmented Requirements

This workflow eliminates requirements toil. Product Owners write stories; AI handles the labor-intensive analysis, drafting, and formatting. Business Owners approve. Developers implement.

**What AI handles (the toil):**

- Reading package management to discover internal dependencies
- Reading context files to understand external dependencies and system constraints
- Parsing API interfaces to validate parameters and types
- Checking for existing step definitions to reuse
- Generating boundary conditions from API constraints + QA patterns
- Applying fuzzing patterns for edge case coverage
- Drafting Gherkin scenarios (including boundary test cases)
- Creating skeleton step definitions
- Opening PRs with proper tags and links

**What humans handle (the judgment):**

- Writing the initial story and acceptance criteria
- Approving that AI's interpretation matches intent
- Resolving ambiguity when AI escalates
- Implementing the actual test code
- Final review before merge

**The ratio (target):** AI handles the mechanical work. Humans provide high-value judgment. The goal is that most stories need only a quick review—exact ratios depend on context file completeness and domain complexity.

### 1.2 AI Usage Principles: Deterministic Artifacts, Not Autonomous Decisions

**Good use of AI: Generate deterministic artifacts** (code, specifications, processes) based on well-defined rules and context. AI excels at mechanical work—reading documentation, parsing APIs, applying patterns, drafting specifications.

**Poor use of AI: Making decisions with loose rulesets without supervision.** AI should not make business decisions, security tradeoffs, or compliance interpretations autonomously. These require human judgment.

**This workflow follows this principle:**
- **AI generates artifacts:** Gherkin specifications, step definition skeletons, boundary test cases, unit tests, implementation code
- **Humans make decisions and apply engineering:**
  - **Business owners** work with conversational agents to specify requirements, then approve AI-drafted specs
  - **Developers** review AI-generated code for security, performance, consistency, and architecture (engineering, not brick-laying)
  - **Security teams** validate approaches and review security implications
- **Rules are explicit:** Context files, API specs, and patterns define what AI should do
- **Supervision is built-in:** Every AI-drafted spec requires BO approval; every implementation requires developer engineering review
- **Escalation is mandatory:** When AI encounters ambiguity or missing rules, it stops and asks humans

**The result:** AI eliminates toil (reading docs, formatting, drafting, generating boilerplate) while humans retain control over judgment (approval, priorities, tradeoffs) and apply software engineering (security, performance, architecture, maintainability).

### 1.3 How It Works: Context Files Enable AI Autonomy

**The key insight: AI can do what Dev and QA would do in a discussion—if it has access to the same information.**

Traditional BDD requires a "Three Amigos" conversation: PO explains the requirement, Dev assesses feasibility, QA identifies edge cases. This is valuable but slow. It doesn't scale when implementation is fast.

This workflow preserves the value while eliminating the synchronous meeting:

| Traditional (slow) | AI-Augmented (fast) |
|-------------------|---------------------|
| PO explains domain terms | AI reads `business.md` |
| Dev checks internal API constraints | AI reads package management + internal API specs |
| Dev checks external API constraints | AI reads `architecture.md` + external API specs |
| Dev spots missing parameters | AI parses interface specs, detects required fields |
| QA suggests edge cases | AI reads `testing.md` patterns, applies to story |
| QA recalls boundary conditions | AI generates boundaries from API constraints + fuzzing patterns |
| Team discusses, someone writes Gherkin | AI drafts Gherkin immediately |
| Async review via PR comments | Same—but draft already exists |

**Context files are your team's knowledge, externalized.** When `architecture.md` documents external dependencies, `testing.md` captures your step library and boundary patterns, and `business.md` defines your domain—AI has everything it needs to draft specifications without waiting for humans to be available.

**Internal dependencies come from code.** AI reads package management files and linked API specs to understand internal services. No need to document these separately.

**External dependencies need documentation.** Third-party APIs can't be discovered from code—document them in `architecture.md` with links to their interface specs.

**QA expertise scales through patterns.** QA defines boundary condition patterns and references fuzzing libraries once in `testing.md`. AI applies these patterns to every story automatically—generating edge cases for each parameter based on its type and constraints.

**Humans review output, not process.** Instead of attending meetings to produce specifications, humans review AI-drafted PRs. The quality gate remains. The toil disappears.

### 1.4 Speed Targets

| Path | Time to PR Ready | Human Time Required |
|------|------------------|---------------------|
| AI Confident | Same day | ~5 min review |
| AI Confident + Inferences | Same day | ~15 min review |
| AI Escalates | 2-3 days | Discussion required |

**Goal:** Most stories need only a quick PR review—no meetings, no back-and-forth, no waiting for schedules to align. The proportion depends on how complete your context files and API specs are.

**Compare to traditional:** A Three Amigos session takes 30-60 minutes per story, requires synchronous availability of 3+ people, and still needs someone to write the Gherkin afterward. This workflow reduces human time per story from hours to minutes.

### 1.5 Enforced Rules

| Rule | Enforcement | Rationale |
|------|-------------|-----------|
| `@story-{id}` on all scenarios | CI blocks merge | Traceability to business request |
| `@pending` blocks impl PR for that story | CI blocks merge | Prevents incomplete implementation |
| BO approval for .feature files | CODEOWNERS + branch protection | Business alignment guaranteed |
| Context files have CODEOWNERS | Branch protection | Accountable ownership |

### 1.6 Source Control as Single Source of Truth

All specifications live in the same repository as the code, as Gherkin/Cucumber `.feature` files:

- **Stories link to feature files** via `@story-{id}` tags
- **Tags track status:** `@pending` (awaiting implementation, removed when implementation merges)
- **Specs merge to main before implementation** so developers always work from approved specs
- **Branch protection enforces BO approval** for all feature file changes
- **CI validates completeness** (no orphan drafts, all scenarios tagged)
- **History is auditable** (who approved what, when, why)

---

## 2. Roles & Responsibilities

| Role | Primary Responsibility | Context File Ownership |
|------|------------------------|----------------------|
| **Product Owner** | Write stories with acceptance criteria | business.md |
| **Business Owner** | Approve all spec changes (enforced via CODEOWNERS) | business.md (review) |
| **Developer** | Implement step definitions, ensure code quality | (consumer) |
| **QA Lead** | Maintain test patterns, boundary conditions | testing.md |
| **Tech Lead** | Maintain technical standards, IoC patterns | tech_standards.md, architecture.md |

**Training Requirements:**
- **PO:** 2 hours (story writing, acceptance criteria)
- **BO:** 1 hour (PR review, Gherkin basics)
- **Developer:** 4 hours (IoC patterns, Godog, TDD) + reference to tech_standards.md
- **QA Lead:** 4 hours (boundary patterns, testing.md structure)
- **Tech Lead:** Minimal (defines the standards)

**Traditional vs. AI-Augmented:**
- Traditional Three Amigos: 2-4 hours per story (PO + Dev + QA meeting, manual Gherkin writing, async reviews)
- AI-Augmented: 15-30 minutes per story (PO writes story → AI drafts → BO reviews PR asynchronously)

---

## 3. Integration Architecture

**Speed requires automation.** Manual triggers, copy-paste workflows, and human-mediated handoffs reintroduce the toil this workflow eliminates. This section defines the integration points that enable fully automated story-to-specification flow.

### 4.1 System Components

| Component | Role | Examples |
|-----------|------|----------|
| Ticket System | Story intake, status tracking, escalation threads | Jira, Azure DevOps, Linear, GitHub Issues |
| AI Agent | Analyzes stories, drafts Gherkin, creates PRs | Custom agent, GitHub Copilot Workspace, Claude |
| Source Control | Stores specifications, context files, step definitions | GitHub, GitLab, Azure Repos, Bitbucket |
| CI/CD | Validates specs, runs tests, enforces rules | GitHub Actions, GitLab CI, Azure Pipelines |

### 2.2 Trigger Mechanism

AI review begins when a ticket meets trigger criteria.

#### 2.2.1 Webhook Trigger (Recommended)

1. Ticket system fires webhook on status change to "Ready for Specification"
2. Webhook payload includes ticket ID (used as branch name), title, description, acceptance criteria
3. AI agent service receives webhook, fetches context files from repo
4. AI performs analysis and either creates branch + PR or posts escalation to ticket

#### 2.2.2 Manual Invocation (Fallback)

For teams not ready for webhook automation: user tags AI agent in ticket comment or runs CLI command with ticket ID.

### 2.3 PR Creation Flow

When AI drafts specifications, it creates a branch (using ticket ID from URL) and PR:

| PR Element | Value | Purpose |
|------------|-------|---------|
| Branch | `spec/{ticket-id}` (extracted from ticket URL) | Isolates spec work |
| Target Branch | `main` | Specifications merge to trunk |
| Title | `[SPEC] {ticket-id}: {ticket-title}` | Identifies spec PRs in queue |
| Body | Links to ticket, lists scenarios, notes confidence level | Context for reviewers |
| Labels | `specification`, `ai-generated`, `needs-bo-approval` | Filtering and routing |
| Reviewers | Auto-assigns BO from CODEOWNERS | BO must approve feature files |

**After BO approval:** Spec PR merges to main with `@pending` tags. Developers then create a separate `impl/{ticket-id}` branch for implementation.

### 2.4 Status Tracking via Tags

**Tags are the source of truth for specification status—not ticket or branch status.** Ticket systems can track their own workflow, but the authoritative state lives in the feature files themselves.

| Tag | Meaning | Where It Lives |
|-----|---------|----------------|
| `@story-{id}` | Links to originating ticket | All scenarios |
| `@pending` | Awaiting implementation (removed after implementation merges) | `main` branch after spec PR merges, removed after impl PR merges |

**@pending lifecycle:** Added when spec merges, removed when implementation for that story merges. CI blocks merge if implementation PR contains `@pending` scenarios for that specific story.

**Ticket status synchronization is optional.** If your team wants ticket status to reflect spec progress, configure webhooks on PR merge events. But the tags in the feature files are what CI enforces and what matters for the workflow.

### 2.5 Required API Integrations

Your AI agent needs API access to:

- **Ticket System:** Read tickets (including extracting ID from URL), update status, post comments
- **Source Control:** Read context files, create branches, commit files, open PRs
- **CI System (optional):** Query validation results for pre-merge checks

---

## 4. Source Control Strategy

Feature files live in the same repository as the code. This section defines the directory structure, branching model, merge rules, and conflict resolution.

### 4.1 Repository Structure

**Feature files live in the same repository as the code they specify.** This ensures specifications and implementation stay in sync, enables atomic commits that include both spec and code changes, and simplifies CI configuration.

#### 3.1.1 Directory Structure

The example below shows one common layout. **Adapt directory structure to your language and framework conventions**—this is a point of flexibility. What matters is that context files, API specs, feature files, and step definitions all live in the same repo.

```
/
├── context/
│   ├── business.md
│   ├── architecture.md
│   ├── testing.md
│   └── tech_standards.md
├── api/
│   └── specs/
│       ├── auth-api.yaml
│       └── users-api.yaml
├── features/                       # Or gherkin/, specs/, test/features/, etc.
│   ├── auth/
│   │   └── auth_login.feature
│   └── editor/
│       └── editor_create.feature
├── step_definitions/               # Location varies by framework
│   ├── auth_steps.{ext}
│   └── editor_steps.{ext}
├── src/
└── .github/
```

**Why same repo:**
- Spec changes and code changes in the same commit = atomic, traceable
- Branch protection applies to both specs and code
- No cross-repo coordination or version drift
- CI runs specs against the code in the same PR

### 4.2 Branching Model

Use trunk-based development with short-lived branches. Spec branches merge after BO approval; implementation branches are separate.

#### 3.2.1 Branch Types

| Branch Pattern | Purpose | Lifecycle | Merge Target |
|----------------|---------|-----------|--------------|
| `main` | Production-ready specifications and code | Permanent | n/a |
| `spec/{ticket-id}` | AI-drafted Gherkin specifications | Merges after BO approval | `main` |
| `impl/{ticket-id}` | Step definition implementation | Merges after tests pass | `main` |
| `amend/{ticket-id}` | Changes to approved specs | Merges after BO approval | `main` |
| `hotfix-spec/{ticket-id}` | Emergency spec additions | < 1 day | `main` |

#### 3.2.2 Two-Branch Workflow

Specifications and implementation use separate branches:

1. AI extracts ticket ID from URL (e.g., `PROJ-1234`)
2. AI creates branch `spec/PROJ-1234`, drafts Gherkin with `@pending` tags, opens PR
3. BO reviews and approves PR
4. **Spec PR merges to main** — scenarios tagged `@pending @story-PROJ-1234`
5. Developer creates NEW branch `impl/PROJ-1234` from main
6. Developer implements step definitions, runs tests
7. Remove `@pending` tags from scenarios (keep `@story-PROJ-1234`)
8. Implementation PR merges to main (CI blocks if `@pending` present for this story)
9. Demo/release to business for validation (tracked outside tags)

**Why separate branches:** Spec work is complete when BO approves. Merging it immediately means main always has the latest approved specifications. Developers start implementation from a clean main branch with the spec already in place.

#### 3.2.3 Branch Protection Rules

Configure these rules on the main branch:

- **Require pull request reviews:** Minimum 1 approval (2 for post-escalation)
- **Require BO approval for feature files:** All `.feature` file changes require Business Owner sign-off (enforced via CODEOWNERS)
- **Require status checks:** CI validation must pass
- **Require linear history:** Squash or rebase merges only
- **Restrict force pushes:** No force push to main
- **Require CODEOWNERS review:** Context file changes need owner approval

### 4.3 Merge Strategy

#### 3.3.1 For Story PRs

Use squash merge. The commit message should include the ticket ID and a summary of scenarios added plus implementation notes. This keeps main history clean and traceable.

### 4.4 Conflict Resolution

When multiple stories touch the same feature file:

#### 3.4.1 Prevention

- AI checks for open PRs on same feature file before creating new PR
- If conflict detected, AI adds scenarios to existing PR instead of creating new one
- Use feature file size limits (max 15 scenarios) to reduce collision probability

#### 3.4.2 Resolution

When conflicts occur despite prevention:

1. First PR to merge wins
2. Conflicting PR rebases onto main
3. AI re-validates scenarios still make sense with merged changes
4. If AI detects logical conflict (not just textual), escalate to humans

---

## 5. Context Files

Context files are the foundation of AI autonomy—and therefore the foundation of requirements speed. **Investment in context files pays compound returns:** every hour spent documenting your systems saves many hours of future requirements toil.

### 4.1 The Four Files

| File | Purpose | AI Uses It To... | Owner |
|------|---------|------------------|-------|
| `business.md` | Domain knowledge | Validate requirements, check compliance | Product Owner |
| `architecture.md` | External dependencies + system context | Document third-party APIs, constraints AI can't discover from code | Dev/Architect |
| `testing.md` | Test knowledge + boundary patterns + fuzzing refs | Reuse steps, generate boundary cases, identify edge cases | QA Lead |
| `tech_standards.md` | Language, framework, directory structure, patterns | Generate step definitions correctly, put files in right places | Tech Lead |

**Internal dependencies:** AI reads package management files (`go.mod`, `package.json`, etc.) and linked API specs to understand internal services. No need to duplicate in context files.

**External dependencies:** Third-party APIs (Stripe, Twilio, AWS, etc.) must be documented in `architecture.md`. AI can't discover these from code.

**Critical: `testing.md` must include boundary condition patterns.** QA defines patterns for each data type (strings, integers, dates, etc.) once. AI applies these patterns to every story automatically, generating boundary test cases without QA needing to attend each story discussion.

**Critical: `tech_standards.md` must specify directory structure.** AI needs to know where feature files and step definitions go for your language/framework. Document your conventions once; AI follows them for every story.

### 4.2 Shared and Organization-Level Context

Context files can be shared across projects. Higher-level context typically lives in separate repositories.

| Level | Location | What to Share | Example |
|-------|----------|---------------|---------|
| Organization | Dedicated repo (`org-context`) | Company-wide compliance, security standards, shared personas | `compliance.md`, `security.md` |
| Domain | Dedicated repo (`domain-{name}-context`) | Domain glossary, business rules for a product area | `payments-context/business.md` |
| Team | Dedicated repo or shared location | Shared step libraries, common API specs | `team-qa-context/testing.md` |
| Project | Same repo as code | Project-specific architecture, local conventions | `context/architecture.md` |

**Layered context:** AI reads context files from multiple levels, with project-level overriding domain-level overriding org-level. This enables:

- **Consistency:** Compliance rules defined once, applied everywhere
- **Reuse:** Shared step libraries reduce duplication across projects
- **Flexibility:** Projects can override org defaults when needed

**Versioning with semver:** All shared context repos use semantic versioning. Projects declare which versions they depend on:

```yaml
# context-dependencies.yaml (in project repo)
org-context: "^2.0.0"      # Compatible with 2.x
payments-context: "~1.3.0"  # Compatible with 1.3.x
team-qa-context: "^3.1.0"   # Compatible with 3.x
```

**Version semantics for context files:**
- **Major (2.0.0 → 3.0.0):** Breaking changes—removed terms, renamed concepts, changed rules that require spec updates
- **Minor (2.0.0 → 2.1.0):** Additions—new terms, new patterns, new steps that don't break existing specs
- **Patch (2.0.0 → 2.0.1):** Fixes—typos, clarifications, no behavioral change

**Implementation:** Use git submodules, package manager, or AI agent config to fetch specific versions of shared context repos.

### 4.3 Why This Matters: Scaling Expertise, Not Replacing People

**Context files don't replace engineers and QA—they amplify their expertise.** Context files handle routine questions and encode established structures so experts can focus on difficult problems.

When AI drafts Gherkin for a login feature:

1. It reads package management to discover internal service dependencies
2. It reads the auth API interface spec to identify required parameters, response schemas, and error codes
3. It checks `architecture.md` for any external dependencies (third-party auth providers, etc.)
4. It checks `testing.md` for existing authentication steps to reuse
5. It checks `business.md` for security rules and compliance requirements
6. It checks `tech_standards.md` for how to structure the step definitions

**If all this information exists, AI doesn't need to ask about routine concerns.** It already has the patterns experts established. No meeting required. No waiting for availability. No calendar Tetris.

**AI catches what humans miss.** By reading actual API interfaces, AI detects missing required parameters, invalid enum values, and constraint violations before the story is even implemented. A developer reviewing AI's draft sees these issues flagged upfront—not discovered during implementation.

**If information is missing or the problem is complex, AI escalates with specific questions.** Humans answer asynchronously—applying judgment, making tradeoffs, resolving ambiguity. Then AI proceeds. Still faster than scheduling a meeting.

**Context files are expertise, codified.** When an engineer answers an escalation about "how to handle third-party API rate limits," that answer belongs in `architecture.md`. When QA identifies a new edge case pattern, it belongs in `testing.md`. The file becomes a living record of decisions—maintained by experts, applied by AI.

**The math:** Well-maintained context files (with links to actual API interfaces) enable AI to handle most routine stories autonomously. Complex stories that require judgment get escalated with specific questions—engineers and QA apply their minds to difficult problems, not routine ones.

### 4.4 Context File Governance

#### 4.4.1 Ownership via CODEOWNERS

Create a CODEOWNERS file to enforce accountability:

```
# Feature file ownership - BO approval required
/gherkin/features/        @business-owner-team

# Context file ownership
/context/business.md      @product-owner-team
/context/architecture.md  @architecture-team
/context/testing.md       @qa-lead-team
/context/tech_standards.md @tech-lead-team
```

**Critical:** All `.feature` file changes require Business Owner approval before merge. This ensures business alignment for all acceptance criteria.

#### 4.4.2 Freshness Monitoring

Implement automated staleness detection:

- **Warning threshold:** 30 days since last update
- **Alert threshold:** 60 days since last update
- **Mechanism:** Weekly scheduled CI job checks last-modified dates
- **Action:** Creates ticket for owner to review and update or confirm current

#### 4.4.3 Version Coupling

Context files use semantic versioning. Each file includes a version header:

```yaml
---
version: 2.1.0
last-reviewed: 2024-01-15
reviewed-by: @jane-doe
changelog: Added MFA requirements for admin roles
---
```

When AI drafts specifications, it records which context versions it used in the PR description:

```markdown
## Context Versions Used
- org-context: 2.0.3
- payments-context: 1.3.1
- project context/business.md: 1.2.0
- project context/testing.md: 3.0.0
```

This enables auditing if issues arise—you can trace back to exactly what context AI was working with.

### 4.5 API Interface Requirements

**Internal dependencies:** AI reads your package management files (`go.mod`, `package.json`, `pom.xml`, `requirements.txt`, etc.) to discover internal service dependencies and their interfaces. No need to duplicate this in `architecture.md`.

**External dependencies:** Third-party APIs and external services must be documented in `architecture.md` with links to their interface specs. AI can't discover these from package files.

#### 4.5.1 What Goes Where

| Dependency Type | Source of Truth | Example |
|-----------------|-----------------|---------|
| Internal services | Package management + linked specs | `go.mod` references `github.com/company/auth-service` |
| Internal APIs | OpenAPI/GraphQL specs in repo | `/api/specs/auth-api.yaml` |
| External APIs | `architecture.md` + linked specs | Stripe, Twilio, AWS services |
| Shared libraries | Package management | Internal packages, SDKs |

#### 4.5.2 Supported Interface Formats

AI reads machine-readable interface specs to detect missing parameters, validate types, and identify constraints:

| Format | File Types | AI Capabilities |
|--------|------------|-----------------|
| OpenAPI/Swagger | `.yaml`, `.json` | Required/optional params, types, enums, constraints |
| GraphQL | `.graphql`, `.gql` | Query/mutation fields, required args, types |
| Protocol Buffers | `.proto` | Message fields, required/optional, types |
| JSON Schema | `.json` | Property constraints, required fields, patterns |
| TypeScript Interfaces | `.ts`, `.d.ts` | Type definitions, optional properties |

#### 4.5.3 architecture.md Structure

`architecture.md` focuses on external dependencies and system context that AI can't discover from code:

```markdown
## API: Authentication Service

### Overview
Handles user authentication and session management.

### Interface Contracts
- **OpenAPI Spec:** `/api/specs/auth-api.yaml`
- **Base URL:** `https://api.example.com/v1/auth`

### Key Constraints (human context)
- Rate limit: 5 attempts per minute per IP
- Sessions expire after 24 hours
- MFA required for admin roles
```

**The key:** Link to actual specs. AI reads `/api/specs/auth-api.yaml` directly and extracts required parameters, response schemas, and error codes. The human-written notes provide context that isn't in the spec (rate limits, business rules, gotchas).

#### 4.5.4 What AI Detects From Interfaces

When AI reads your API interfaces, it can identify:

- **Missing required parameters** in acceptance criteria
- **Invalid enum values** that don't match the schema
- **Type mismatches** (expecting string, spec says integer)
- **Undocumented endpoints** referenced in stories
- **Deprecated fields** that shouldn't be used
- **Constraint violations** (min/max, patterns, formats)

This catches issues that would otherwise surface during implementation—or worse, in production.

### 4.6 Keeping Files Current

Add to sprint ceremonies as a checklist:

- [ ] New domain terms this sprint? → Update `business.md`, bump minor version
- [ ] New integrations or APIs? → Update `architecture.md`, bump minor version
- [ ] API interface changed? → Ensure spec files are updated (OpenAPI, GraphQL, etc.)
- [ ] New step patterns? → Update `testing.md`, bump minor version
- [ ] New coding conventions? → Update `tech_standards.md`, bump minor version
- [ ] Breaking changes to any context? → Bump major version, communicate to dependent projects
- [ ] **AI escalations this sprint?** → Review answers, codify into appropriate context file

**If AI escalates frequently, your context files are incomplete.** Every escalation means a story that couldn't be processed autonomously—a story that required human discussion, calendar coordination, and wait time. Fix the source (incomplete documentation), not the symptom (slow requirements).

**Escalations are learning opportunities.** When AI escalates with missing information and humans answer, those answers should typically be codified into context files:
- Missing business rule → Update `business.md`
- External API behavior → Update `architecture.md`
- New edge case pattern → Update `testing.md`
- Technical standard clarification → Update `tech_standards.md`

This creates a feedback loop: escalations → human answers → documented → fewer future escalations. **Experts apply their minds to difficult problems, then capture the decisions for future reuse.**

**API interfaces must stay in sync with implementations.** If your OpenAPI spec doesn't match your actual API, AI will flag false issues or miss real ones. Treat spec files as first-class artifacts—update them when the API changes.

**Context file maintenance is not overhead—it's investment in speed.** A 30-minute update to `architecture.md` (or the linked API spec) after adding a new endpoint can save hours of escalation and discussion across all future stories that touch that API.

**Bump versions when you make changes.** Semantic versioning enables projects to declare compatible versions and makes it easy to trace which context AI used when drafting specifications.

---

## 6. Workflow Stages

### 5.0 Stage 0 (Optional): Requirements Drafting Assistance (BO: ~15-30 minutes, conversational)

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

### 5.1 Stage 1: Story Creation (PO: ~10 minutes)

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

**Trigger:** Ticket status changes to "Ready for Specification" → AI creates branch and begins review immediately.

### 5.2 Stage 2: AI Review & Drafting (AI: seconds to minutes)

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

#### 5.2.1 If Confident

AI creates:

- `.feature` file with scenarios tagged `@pending @story-{id}` (ready for implementation after merge)
- Skeleton step definitions (method signatures, no implementation)
- PR ready for BO review with confidence notes

#### 5.2.2 If Uncertain

AI posts to ticket:

- What it's uncertain about
- Specific questions for humans
- What context is missing

### 5.3 Stage 3: BO Review (BO: ~5-15 minutes, async)

**Agent:** `bo-review` (One-Shot) - Optional guide for BO reviewers
**Prompt:** [`docs/prompts/bo-review/`](../prompts/bo-review/)

The PR is where Business Owner verifies AI's work. This is the quality gate—but it's a **review** gate, not a **creation** gate. The hard work is done.

**Compare to traditional:** BO would attend a meeting, discuss requirements, then wait for someone to write Gherkin, then review. Here, BO reviews a complete draft at their convenience.

#### 5.3.1 PR Contents

- Feature file with scenarios
- Skeleton step definitions
- Context versions used (for traceability)
- AI confidence notes (if any inferences were made)
- API validation results
- Link to original ticket

#### 5.3.2 Reviewer Checklist

| Check | Looking For |
|-------|-------------|
| Intent match | Do scenarios capture what the story asks for? |
| Inference review | Are AI's assumptions acceptable? |
| Gap check | Any obvious missing scenarios? |
| Tag check | Is `@story-{id}` present? |

#### 5.3.3 Approval Requirements

| AI Confidence | Approvals Needed | Rationale |
|---------------|------------------|-----------|
| High (no flags) | 1 (must include BO) | Routine work, BO validates business intent |
| Medium (flagged inferences) | 1 (must include BO) | Verify assumptions are acceptable |
| After escalation | 2 (must include BO) | Higher scrutiny for complex cases |

**BO approval is mandatory.** Business Owner must approve all feature file changes to ensure acceptance criteria align with business requirements.

**On BO approval:** PR merges to main. Scenarios remain `@pending`. Spec work is complete. Developers can now begin implementation.

**Note:** For meta-instructions on running these agents with Claude and logging requirements, see `../CLAUDE.md`.

### 5.4 Stage 4: Implementation

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

## 7. Step Definition Lifecycle

This section addresses how step definitions evolve from AI-generated skeletons through implementation and maintenance.

### 6.1 Skeleton Generation

When AI creates a PR, it generates skeleton step definitions based on `tech_standards.md` patterns:

```go
// Generated by AI - PROJ-1234
// TODO: Implement
func (s *AuthSteps) theUserEntersValidCredentials() error {
    panic("not implemented")
}
```

### 6.2 Implementation

Developer implements the skeleton, potentially adjusting the signature:

| Change Type | Action Required | Who Does It |
|-------------|-----------------|-------------|
| Parameter added | Update feature file step text | Developer in same PR |
| Step renamed | Update feature file step text | Developer in same PR |
| Step split into multiple | Update feature file, may need new scenarios | Developer, BO re-approves |
| Step consolidated | Update all affected feature files | Developer, BO re-approves |

### 6.3 Signature Change Protocol

When implementation requires changing a step signature:

- **Same PR rule:** Step definition changes and feature file updates must be in the same PR
- **BO re-approval:** If feature file changes, BO must approve again
- **Backward compatibility:** If step is used in multiple features, keep old signature as deprecated alias
- **Documentation:** Update `testing.md` with new step signature

### 6.4 Shared Step Libraries

For enterprises with multiple teams sharing steps:

#### 6.4.1 Structure

```
/shared-steps/
├── common/           # Cross-domain steps (login, navigation)
├── domain-a/         # Domain-specific shared steps
└── domain-b/
```

#### 6.4.2 Governance

- **Shared steps require 2 approvals** from different teams
- **Breaking changes require deprecation period** (2 sprints minimum)
- **Shared step changes trigger CI across all consuming repos**

---

## 8. AI Confidence & Escalation

### 7.1 Confidence Levels

| Level | Symbol | Meaning | AI Action |
|-------|--------|---------|-----------|
| High | ✓ | All info in context files, existing patterns | Drafts PR directly |
| Medium | ⚠ | Made reasonable inferences from context | Drafts PR with flagged assumptions |
| Low | ✗ | Missing context or high uncertainty | Escalates, does not draft |

### 7.2 Escalation Triggers

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

### 7.3 Escalation Format

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

### 7.4 Conflict Resolution Process

When business rules conflict with technical constraints:

- **AI documents both constraints** in escalation
- **Product Owner and Tech Lead** must both respond
- **Resolution recorded** in ticket and relevant context file
- **AI re-analyzes** with updated context

**After resolution:** The conflict resolution decision should be documented in the appropriate context file so AI (and future humans) understand how to handle similar situations.

---

## 9. Scenario Design Guidance

This section provides rules for well-structured Gherkin specifications.

### 8.1 Feature File Organization

#### 8.1.1 One Capability Per Feature

Each `.feature` file should cover a single user capability. If you find yourself using "and" to describe the feature, split it.

- **Good:** `auth_login.feature`, `auth_password_reset.feature`
- **Bad:** `auth_login_and_password_reset.feature`

#### 8.1.2 Size Limits

| Metric | Limit | Rationale |
|--------|-------|-----------|
| Scenarios per feature | 15 max | Reduces merge conflicts, improves readability |
| Steps per scenario | 10 max | Keeps scenarios focused |
| Examples per outline | 10 max | Prevents data-driven bloat |

### 8.2 When to Use Scenario Outline

Use Scenario Outline when:

- Same behavior with different data (validation rules, boundary conditions)
- Data combinations matter (role × action permission matrix)
- Error messages vary by input

Do NOT use Scenario Outline when:

- Workflow differs between cases (use separate scenarios)
- Only 2-3 examples exist (inline is clearer)
- Examples require different setup (Background won't work)

### 8.3 Background Usage

Use Background for:

- Common preconditions shared by ALL scenarios in the feature
- User authentication when all scenarios need a logged-in user
- Data setup that doesn't vary between scenarios

Do NOT use Background when:

- Some scenarios need different setup
- Background would exceed 5 steps
- Setup is complex enough to obscure scenario intent

### 8.4 Step Writing Rules

#### 8.4.1 Given-When-Then Discipline

| Keyword | Purpose | Verb Tense |
|---------|---------|------------|
| Given | Establish preconditions | Past/present state |
| When | Action under test | Present tense |
| Then | Assert outcomes | Present/future state |
| And/But | Continue previous keyword | Match previous |

#### 8.4.2 Step Reuse Priority

1. **First:** Reuse existing step from `testing.md` exactly as written
2. **Second:** Parameterize existing step if minor variation needed
3. **Third:** Create new step only if no existing step fits

New steps must be documented in `testing.md` before PR merge.

---

## 10. SLA Enforcement

This section defines review timelines and enforcement mechanisms.

### 9.1 Review Timeline Targets

| AI Confidence | Review SLA | Escalation After |
|---------------|------------|------------------|
| High | 4 hours | 8 hours |
| Medium | Same business day | 24 hours |
| Post-escalation | 1 business day | 48 hours |

### 9.2 Enforcement Mechanisms

#### 9.2.1 Automated Reminders

- **At 50% of SLA:** Slack/Teams notification to assigned reviewers
- **At 75% of SLA:** Escalation to team lead
- **At 100% of SLA:** Visible in team dashboard, blocks new spec PRs for that reviewer

#### 9.2.2 Auto-Merge Option

For HIGH confidence PRs only:

- If CI passes and no review within 24 hours, auto-merge with notification
- Team can disable auto-merge if preferred
- Auto-merged PRs flagged for retrospective review

### 9.3 Escalation Response SLAs

| Escalation Type | Response SLA | Resolver |
|-----------------|--------------|----------|
| Missing architecture info | 1 business day | Architect/Tech Lead |
| Business rule clarification | 4 hours | Product Owner |
| Compliance question | 2 business days | Compliance team |
| Security concern | 1 business day | Security team |

---

## 11. Rollback & Revert Procedures

This section covers how to handle specifications that need to be reversed after merge.

### 10.1 When to Revert

- **Incorrect scenarios:** Scenarios don't match intended behavior
- **Blocking development:** Scenarios prevent valid implementation
- **Business change:** Requirements changed after spec approval
- **Security issue:** Scenarios expose security vulnerability in test data

### 10.2 Revert Process

#### 10.2.1 For @pending Scenarios (Not Yet Implemented)

- Create `amend/{ticket-id}-revert` branch
- Remove or modify scenarios
- PR requires 1 approval (original author cannot approve)
- BO must approve if feature files change
- Link to original ticket with reason

#### 10.2.2 For Implemented Scenarios

- Create ticket for revert with business justification
- Revert PR must include both scenario removal AND step cleanup
- Requires 2 approvals (including BO)
- Update context files if revert reveals gap

### 10.3 Amendment vs. Revert

| Situation | Action | Branch Pattern |
|-----------|--------|----------------|
| Minor correction to scenario | Amend | `amend/{ticket-id}-fix` |
| Add missing edge case | Amend | `amend/{ticket-id}-edge` |
| Completely wrong approach | Revert + new spec | `amend/{ticket-id}-revert` |
| Feature cancelled | Revert | `amend/{ticket-id}-cancel` |

---

## 12. Tag Taxonomy

**Tags are the authoritative status of specifications.** Not ticket status, not branch status—the tags in the feature files themselves.

### 11.1 Lifecycle Tags

| Tag | Meaning | Location | In CI? |
|-----|---------|----------|--------|
| `@pending` | Awaiting implementation (temporary) | `main` branch between spec merge and impl merge | No (skipped) |

**Tag lifecycle:**
- Spec PR merges: Scenarios tagged `@pending @story-{id}`
- Implementation completes: Developer removes `@pending` tags for that story
- Implementation PR: **CI blocks merge if `@pending` present for that story**
- After merge: Scenarios have only `@story-{id}` (no lifecycle tag)
- Business validation tracked in deployment/release process

### 11.2 Required Tags

| Tag | Format | Purpose |
|-----|--------|---------|
| `@story-{id}` | `@story-PROJ-1234` | Links to originating ticket |

### 11.3 Optional Tags

| Tag | Purpose |
|-----|---------|
| `@smoke` | Include in fast feedback suite |
| `@regression` | Include in full regression suite |
| `@hotfix` | Emergency addition, expedited review |
| `@wip` | Work in progress, skip in CI |

---

## 13. Hotfix Path

For production incidents requiring immediate test coverage.

### 12.1 Standard vs. Hotfix

| Aspect | Standard Process | Hotfix Process |
|--------|------------------|----------------|
| AI review | Required | Skip |
| BO review | Full review | Security + functionality only |
| Approvals | 1-2 depending on confidence | 1 approval (BO if available) |
| Tag | `@pending` | `@hotfix @pending` |

### 12.2 Hotfix Process

1. Tag scenarios `@hotfix @pending @story-{id}`
2. CI allows `@hotfix` to merge with expedited review
3. Within 5 business days: Remove `@hotfix`, complete standard BO review
4. If not completed in 5 days: Escalate to team lead

---

## 14. Metrics

**The goal is speed parity:** Requirements should not be the bottleneck when development is AI-assisted.

### 13.1 Primary: AI Confidence Rate (Measures Context File Quality)

| Metric | Target | If Below Target |
|--------|--------|-----------------|
| AI confident (High + Medium) | Majority of stories | Context files incomplete—every miss is a story requiring human discussion |
| AI escalates (Low) | Minority of stories | Document missing systems/patterns |

**Why this is primary:** AI confidence directly determines how many stories can be processed without human wait time. Low confidence = slow requirements = development bottleneck. Track this metric over time—it should improve as context files mature.

### 13.2 Secondary: Speed (Measures Process Efficiency)

| Metric | Target | Why It Matters |
|--------|--------|----------------|
| Story → PR ready (confident) | Same day | Requirements shouldn't wait overnight |
| PR → BO Approved | < 24 hours | Review backlog = development blocked |
| Escalation resolution | 2-3 days | Complex cases still need to be fast |

### 13.3 Quality (Measures Accuracy)

| Metric | Target | If Above Target |
|--------|--------|-----------------|
| BO overrides AI | < 10% | AI missing context or miscalibrated |
| Post-merge reverts | < 5% | Review process not catching issues |
| Implementation rejects spec | < 3% | Feasibility assessment weak |

**Quality vs. Speed tradeoff:** This workflow optimizes for speed without sacrificing quality. If quality metrics slip, the problem is usually incomplete context files, not the workflow itself.

---

## 15. Naming Conventions

### 14.1 Branches

| Purpose | Pattern | Example |
|---------|---------|---------|
| Specification | `spec/{ticket-id}` | `spec/PROJ-1234` |
| Implementation | `impl/{ticket-id}` | `impl/PROJ-1234` |
| Amendment | `amend/{ticket-id}` | `amend/PROJ-1234` |
| Hotfix | `hotfix-spec/{ticket-id}` | `hotfix-spec/PROJ-9999` |

**Branch name derived from ticket URL.** AI extracts the ticket ID automatically—no manual entry required.

### 14.2 Files

Follow your framework's conventions. Common patterns:

| Type | Pattern | Notes |
|------|---------|-------|
| Features | `{domain}_{capability}.feature` | Or `{Domain}{Capability}.feature` for PascalCase conventions |
| Steps | `{domain}_steps.{ext}` | Location and naming per framework (e.g., `*_test.go`, `*Steps.java`) |

---

## 16. CI Configuration

Example GitHub Actions workflow for Gherkin validation. **Adjust paths to match your directory structure.**

```yaml
name: Gherkin Validation
on:
  pull_request:
    paths: ['features/**', 'context/**']  # Adjust to your structure

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Require @story tag
        run: |
          # Adjust path to your features directory
          for f in features/**/*.feature; do
            grep -q "@story-" "$f" || { echo "Missing @story: $f"; exit 1; }
          done

      - name: Block impl PR with @pending for this story
        if: startsWith(github.head_ref, 'impl/')
        run: |
          # Extract story ID from branch name (e.g., impl/PROJ-1234 -> PROJ-1234)
          STORY_ID=$(echo "${{ github.head_ref }}" | sed 's/^impl\///')

          # Check if any scenarios have @pending for this story
          if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
            echo "ERROR: Implementation PR contains @pending scenarios for $STORY_ID"
            echo "Remove @pending tags before merging (keep @story-$STORY_ID)"
            exit 1
          fi

      - name: Validate context file freshness
        run: |
          for f in context/*.md; do
            days=$(( ($(date +%s) - $(stat -c %Y "$f")) / 86400 ))
            if [ $days -gt 60 ]; then
              echo "WARNING: $f not updated in $days days"
            fi
          done

      - name: Run implemented scenarios (not @pending)
        run: cucumber --tags "not @pending"
```

---

## 17. Appendix: AI Agent Prompt

Use this prompt template for your AI agent:

```
You are an autonomous requirements analyst. Your job is to review user stories
and either draft Gherkin specifications or escalate to humans.

INPUTS:
- User story and acceptance criteria from ticket
- Ticket ID extracted from URL (use for branch naming: spec/{ticket-id})
- business.md (domain terms, personas, business rules, compliance)
- architecture.md (external dependencies, system context, third-party API specs)
- Package management files (go.mod, package.json, etc.) for internal dependencies
- API interface specs linked from code or architecture.md (OpenAPI, GraphQL, protobuf, etc.)
- testing.md (step library, boundary condition patterns, fuzzing libraries, edge cases)
- tech_standards.md (language, framework, directory structure, step definition patterns)

PROCESS:
1. EXTRACT ticket ID from URL for branch naming
2. READ package management to discover internal dependencies
3. READ architecture.md for external dependencies
4. READ API interface specs (internal and external) to validate parameters
5. ANALYZE the story against all context files
6. GENERATE boundary conditions from API constraints + testing.md patterns
7. ASSESS your confidence level
8. Either DRAFT Gherkin or ESCALATE

API INTERFACE ANALYSIS:
When reading API specs (OpenAPI, GraphQL, etc.), check for:
- Required parameters not mentioned in acceptance criteria → FLAG
- Invalid enum values in acceptance criteria → FLAG
- Type mismatches (string vs int, etc.) → FLAG
- Deprecated endpoints or fields → FLAG
- Missing error handling for documented error codes → FLAG
- Constraint violations (min/max, patterns) → FLAG

BOUNDARY CONDITION GENERATION:
For each parameter in the story, apply patterns from testing.md:
- String fields: empty, whitespace, max length, unicode, special chars
- Numeric fields: zero, negative, min, max, overflow
- Dates: epoch, far future, leap year, timezone edges
- Arrays: empty, single, max items, duplicates
- Reference fuzzing libraries for additional edge cases
Include boundary scenarios in Scenario Outlines where appropriate.

CONFIDENCE ASSESSMENT:
✓ HIGH - Draft directly:
  - Internal dependencies discoverable from package management
  - External dependencies documented in architecture.md
  - API interface specs validate the acceptance criteria
  - Boundary patterns available in testing.md
  - >80% of needed steps exist in testing.md
  - Business rules are clear in business.md
  - No security/auth/compliance ambiguity

⚠ MEDIUM - Draft with notes:
  - Made reasonable inferences from context
  - Flag specific assumptions for BO reviewer
  - Note any API constraints that affect scenarios
  - Generated boundary cases that may need review

✗ LOW - Escalate:
  - External dependency not documented in architecture.md
  - API interface spec missing or inaccessible
  - Required parameters unclear or conflicting
  - Security or auth changes implied
  - No similar patterns in testing.md
  - Conflicting or ambiguous requirements
  - Compliance implications unclear

DRAFTING RULES:
1. Create branch: spec/{ticket-id}
2. Reuse existing steps from testing.md exactly
3. Only create new steps when no existing step fits
4. Follow tech_standards.md for step definitions
5. Tag all scenarios @pending @story-{ticket-id}
6. If MEDIUM confidence, include inference notes for BO
7. Include API validation findings in PR description
8. Generate boundary condition scenarios using testing.md patterns
9. Use Scenario Outlines for boundary testing where multiple values apply

OUTPUT:
- Branch: spec/{ticket-id}
- PR targeting main
- Reviewer: BO (from CODEOWNERS)
- Labels: specification, ai-generated, needs-bo-approval
- PR body includes:
  - Context versions used (org-context, domain-context, project context files)
  - API validation results
  - Boundary cases generated
  - Flagged issues
  - Confidence notes
- On BO approval: PR merges, scenarios stay @pending, developers create impl/{ticket-id} branch

TRUST CALIBRATION:
- >70% confident = draft with notes
- <70% confident = escalate
- Security/auth/compliance = extra scrutiny, lower threshold
- Missing API interface = escalate (can't validate feasibility)
```

---

## 18. Appendix: Context File Templates

### 17.1 business.md Template

```markdown
---
version: 1.0.0
last-reviewed: YYYY-MM-DD
reviewed-by: @owner
changelog: Initial version
---

# Business Context

## Domain Glossary

| Term | Definition |
|------|------------|
| [term] | [definition] |

## User Personas

| Persona | Description | Access Level |
|---------|-------------|--------------|
| [name] | [description] | [access] |

## Business Rules

| ID | Rule | Exceptions |
|----|------|------------|
| BR-001 | [rule] | [exceptions] |

## Compliance Requirements

| Requirement | Standard | Affected Features |
|-------------|----------|-------------------|
| [requirement] | [standard] | [features] |
```

### 17.2 architecture.md Template

```markdown
---
version: 1.0.0
last-reviewed: YYYY-MM-DD
reviewed-by: @owner
changelog: Initial version
---

# Architecture Context

## System Overview

[Mermaid diagram or description of how services connect]

## Internal API Specs

AI discovers internal dependencies from package management. List paths to specs here:

| Service | Spec Location | Format |
|---------|---------------|--------|
| Auth Service | `/api/specs/auth-api.yaml` | OpenAPI 3.0 |
| User Service | `/api/specs/users.graphql` | GraphQL SDL |

## External Dependencies

These cannot be discovered from code—document them here:

### [External Service, e.g., Stripe]

| Attribute | Value |
|-----------|-------|
| Description | [what we use it for] |
| API Docs | [link to their docs] |
| Spec Location | `/specs/external/stripe.yaml` (local copy or link) |
| Auth | [API key, OAuth, etc.] |
| Rate Limits | [limits that affect our usage] |
| Sandbox/Test | [test environment details] |

### [Another External Service]

...

## Technical Constraints

| Constraint | Rationale | Affected Scenarios |
|------------|-----------|-------------------|
| [constraint] | [why] | [what it affects] |
```

**Key principle:** AI reads package management for internal dependencies. `architecture.md` documents external dependencies and system-level context that can't be discovered from code.

### 17.3 testing.md Template

```markdown
---
version: 1.0.0
last-reviewed: YYYY-MM-DD
reviewed-by: @owner
changelog: Initial version
---

# Testing Context

## Step Library

### [Domain] Steps

| Step | Parameters | Notes |
|------|------------|-------|
| `Given [step text]` | [params] | [notes] |
| `When [step text]` | [params] | [notes] |
| `Then [step text]` | [params] | [notes] |

## Boundary Condition Patterns

AI uses these patterns to generate boundary test cases for each data type:

| Data Type | Boundary Conditions | Example Values |
|-----------|--------------------| ---------------|
| String | empty, whitespace-only, max length, unicode, special chars | `""`, `"   "`, `"a"*255`, `"日本語"`, `"<script>"` |
| Integer | zero, negative, min, max, overflow | `0`, `-1`, `-2147483648`, `2147483647` |
| Email | valid format, invalid format, max length, special domains | `"a@b.co"`, `"invalid"`, `"test+tag@example.com"` |
| Date | epoch, far future, leap year, timezone edge | `1970-01-01`, `2099-12-31`, `2024-02-29` |
| Array | empty, single item, max items, duplicate items | `[]`, `[1]`, `[1,1,1]` |
| Currency | zero, negative, max precision, rounding edge | `0.00`, `-0.01`, `999999.99`, `0.005` |

## Fuzzing Libraries

Reference these libraries for edge case generation in step definitions:

| Library | Language | Use For |
|---------|----------|---------|
| [go-fuzz](https://github.com/dvyukov/go-fuzz) | Go | Input fuzzing |
| [Hypothesis](https://hypothesis.readthedocs.io/) | Python | Property-based testing |
| [fast-check](https://github.com/dubzzz/fast-check) | TypeScript/JS | Property-based testing |
| [QuickCheck](https://hackage.haskell.org/package/QuickCheck) | Haskell/ports | Property-based testing |
| [jqwik](https://jqwik.net/) | Java | Property-based testing |

AI references these when generating edge case scenarios. Step implementations can use these libraries for thorough input validation testing.

## Edge Case Patterns

| Pattern | When to Apply | Example |
|---------|---------------|---------|
| [pattern] | [when] | [example] |
| Concurrent access | Multi-user features | Two users edit same record |
| Timeout handling | External API calls | Service unavailable |
| Partial failure | Batch operations | 3 of 5 items fail |
| Rate limiting | High-frequency actions | Exceed API limits |
| Session expiry | Long-running workflows | Token expires mid-flow |

## Test Data

| Data Type | Strategy |
|-----------|----------|
| [type] | [how managed] |
```

**Key principle:** QA defines the patterns and boundary conditions once. AI applies them to every story automatically. This scales QA expertise across all specifications without requiring QA presence at every story discussion.

### 17.4 tech_standards.md Template

```markdown
---
version: 1.0.0
last-reviewed: YYYY-MM-DD
reviewed-by: @owner
changelog: Initial version
---

# Technical Standards

## Language & Framework

- **Language:** [language and version]
- **Test Framework:** [Cucumber-JVM, Godog, Behave, etc.]
- **Feature Location:** [path to feature files, e.g., `src/test/resources/features/`]
- **Step Definitions:** [path and naming, e.g., `src/test/java/steps/*Steps.java`]

## Step Definition Pattern

```[language]
// Template for step definitions in your language/framework
```

## Error Handling

- [pattern]

## Naming Conventions

- [conventions for your language]

## Inversion of Control

- Use constructor-based dependency injection
- Prefer GORM for database access
```

**This file tells AI where things are and how to write them.** Language-specific conventions (directory structure, file naming, step definition syntax) belong here.

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