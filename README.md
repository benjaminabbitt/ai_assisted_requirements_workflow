# AI-Augmented Requirements Workflow

**An enterprise-grade system for accelerating requirements using AI while maintaining human oversight.**

---

## Overview

This system eliminates the requirements bottleneck that emerges when AI-assisted development accelerates implementation velocity. Traditional requirements processes—Three Amigos meetings, manual Gherkin authoring, calendar coordination—can't keep pace with AI-accelerated development.

**The solution:** Apply the same AI acceleration to requirements that we apply to development.

### Key Principle

- **AI generates deterministic artifacts** (specifications, test cases, code skeletons) based on well-defined rules and context
- **Humans make decisions and apply engineering** (business approval, security review, performance optimization, architecture)
- **Result:** AI eliminates toil; humans retain control over judgment

### Speed Targets

| Traditional Process | AI-Augmented Process |
|---------------------|----------------------|
| 30-60 min Three Amigos meeting | 5-15 min async PR review |
| + manual Gherkin writing | AI drafts specifications |
| + calendar coordination | Immediate processing |
| **Total: 2-4 hours per story** | **Total: 5-30 minutes per story** |

---

## Why BDD? Executable Specifications in a Dynamic World

**BDD (Behavior-Driven Development) isn't just testing—it's executable specifications that ensure correctness as the world changes.**

### The Problem: Code Changes, Requirements Drift, AI Generates

In modern development:
- **Developers refactor** code constantly (changing implementation details)
- **AI generates** code that may not match intent (pattern-following ≠ correctness)
- **Requirements evolve** as business needs change
- **Teams grow** and knowledge fragments

**Documentation fails:** It goes stale. No one knows if code matches docs.

**Unit tests fail:** They test implementation details, not business behavior. Refactor the code → tests break even if behavior is correct.

### The Solution: Automated Acceptance Tests as Living Contract

**Gherkin scenarios are executable specifications:**

```gherkin
@story-PROJ-1234
Scenario: User requests password reset successfully
  Given a user exists with email "user@example.com"
  When I request a password reset for "user@example.com"
  Then I should receive a reset email
  And the password reset attempt should be logged
```

**This scenario:**
1. **Defines expected behavior** in business language (not code)
2. **Runs automatically** in CI/CD (validates behavior every commit)
3. **Survives refactoring** (implementation can change, behavior stays same)
4. **Validates AI output** (AI-generated code must pass business requirements)
5. **Documents current behavior** (if scenario passes, this is what the system does)

### BDD Value in This Workflow

#### 1. **AI Generates Code → BDD Validates Correctness**

**Without BDD:**
```
AI generates code → Human reviews → "Looks good" → Ships → Broken in production
```

**With BDD:**
```
AI generates code → BDD scenarios run → Fail → AI refactors → Pass → Behavior verified
```

**Example:**
- AI generates password reset service
- Scenario runs: "Then I should receive a reset email"
- Test fails: Email not sent
- AI fixes implementation
- Scenario passes: Behavior confirmed

**Benefit:** AI can generate code fast, but BDD ensures it actually does what the business requires.

#### 2. **Developers Refactor → BDD Ensures Behavior Preserved**

**Scenario:**
```gherkin
Scenario: Calculate order total with tax
  Given an order with items totaling $100.00
  When I calculate the total for California
  Then the total should be $108.75
```

**Developer refactors:**
- Extracts tax calculator into separate service
- Changes internal data structures
- Optimizes algorithm for performance

**BDD scenario still passes → Behavior preserved.**

Without BDD, developer might unknowingly break edge cases during refactoring.

#### 3. **Requirements Evolve → BDD Tracks What System Actually Does**

**Reality:** Requirements change mid-sprint, stories get amended, edge cases emerge during implementation.

**BDD provides traceability:**
- `@story-PROJ-1234` tag links scenario to originating requirement
- Scenario text shows exactly what behavior was implemented
- If scenario passes, this is what the system does now (not what docs said 3 months ago)

**When business asks:** "Does the system handle expired tokens?"
**Answer:** Run scenarios tagged `@password-reset`. If scenario exists and passes → yes. If not → no.

**Truth lives in passing tests, not stale documentation.**

#### 4. **Multiple Developers + AI → BDD Prevents Fragmentation**

**Problem:** Different developers (or AI sessions) implement similar logic differently.

**Example:**
- API endpoint validates email format one way
- UI validates email format differently
- Background job validates email yet another way

**User tries same email → Different results depending on code path.**

**BDD catches this:**
```gherkin
Scenario: Validate email format
  Given I am using the API
  When I submit email "user@example.com"
  Then it should be accepted

Scenario: Validate email format via UI
  Given I am using the web interface
  When I submit email "user@example.com"
  Then it should be accepted
```

If validation logic fragments, one scenario fails → Developer must unify logic.

#### 5. **Business Validates Actual Behavior (Not Code Review)**

**Traditional:** Business Owner reviews code PR → "I don't understand this code, looks fine I guess?"

**BDD:** Business Owner reviews Gherkin scenarios → "Yes, this matches what I asked for."

**After implementation:**
```
Developer: "Password reset feature is done"
BO: "Show me the scenarios"
Developer: [Runs scenarios with @story-PROJ-1234 tag]
BO: [Watches scenarios execute in demo environment]
BO: "Perfect - the email is sent, audit log captures it. Approved."
```

**BO validates behavior, not code.** This is what they understand.

### Why This Matters for AI-Augmented Development

**AI accelerates code generation, but:**
- AI doesn't understand business intent (only patterns)
- AI may generate syntactically correct but behaviorally wrong code
- AI may introduce subtle bugs humans miss in code review

**BDD provides the safety net:**

1. **AI drafts specification** (Gherkin scenarios) → BO approves
2. **AI generates implementation** → BDD scenarios validate behavior
3. **Scenarios fail** → AI refactors until they pass
4. **Scenarios pass** → Behavior matches business intent

**The result:** Fast AI-generated code that's verified to be correct.

### Key Insight: Specifications That Execute

**This workflow treats Gherkin scenarios as the source of truth:**
- Not documentation (that goes stale)
- Not code (that's implementation detail)
- **Executable specifications** (that validate behavior continuously)

**In a dynamic world where:**
- Code changes constantly (refactoring, AI generation, optimization)
- Requirements evolve (amendments, edge cases, clarifications)
- Teams scale (knowledge fragments, new developers join)
- AI generates code (fast but needs validation)

**Automated acceptance tests ensure:** The system does what the business expects, validated automatically, every commit.

---

## Core Benefits

### For Business Owners

- **Async review instead of synchronous meetings** - Review specifications at your convenience, not in scheduled time blocks
- **AI catches issues early** - Missing parameters, invalid constraints, compliance concerns surfaced before development
- **Specifications stay current** - Living documentation in source control, not stale wiki pages

### For Developers

- **Approved specs before implementation** - No ambiguity, no waiting for clarification
- **Test scenarios ready to implement** - Boundary conditions and edge cases already defined
- **TDD guidance from AI** - Collaborative implementation with pattern enforcement

### For QA

- **Expertise scales without meetings** - Define patterns once in `testing.md`, AI applies them to every story
- **Consistent edge case coverage** - Boundary conditions automatically generated from API specs and QA patterns
- **Fuzzing patterns reused** - Reference fuzzing libraries once, AI uses them everywhere

### For Tech Leads

- **Standards enforced automatically** - AI checks compliance before human review
- **Architecture constraints validated** - AI reads API specs and flags violations before implementation
- **Context files capture decisions** - Expertise codified and reused, not lost in chat history

---

## What This System Provides

### 5 Specialized AI Agents

1. **requirements-drafting-assistant** (Conversational)
   - Helps Business Owners articulate complex requirements through guided dialogue
   - Explores edge cases, validates against constraints, surfaces security implications

2. **requirements-analyst** (One-Shot)
   - Analyzes story tickets and drafts Gherkin specifications
   - Validates against API contracts, generates boundary conditions, creates test scenarios

3. **bo-review** (One-Shot)
   - Guides Business Owners in reviewing AI-drafted specifications
   - Checks business logic, identifies missing scenarios, flags assumptions

4. **standards-compliance** (One-Shot)
   - Reviews code for compliance with technical standards
   - Validates IoC patterns, checks conventions, calculates compliance scores

5. **developer-implementation** (Conversational)
   - Collaborates with developers on TDD implementation
   - Guides test-first development, ensures pattern compliance, supports refactoring

### Complete Workflow Documentation

See **[docs/README.md](docs/README.md)** for comprehensive system documentation including:

- **Integration Architecture** - How to connect AI agents to your ticket system, source control, and CI/CD
- **Context Files** - The foundation of AI autonomy (business rules, architecture, testing patterns, technical standards)
- **Workflow Stages** - Step-by-step process from story creation through implementation
- **AI Confidence & Escalation** - How AI determines when to draft vs. when to ask humans
- **Source Control Strategy** - Branch protection, merge rules, tag taxonomy
- **CI Configuration** - Automated validation, blocking rules, compliance checks
- **Metrics & SLAs** - How to measure and maintain requirements velocity

### Sample Project

See **[sample-project/](sample-project/)** for example context files demonstrating:

- `business.md` - Domain glossary, personas, business rules, compliance requirements
- `architecture.md` - System overview, API specs, external dependencies, constraints
- `testing.md` - Step library, boundary patterns, fuzzing references, edge cases
- `tech_standards.md` - Language conventions, IoC patterns, directory structure

Use these as templates for your own project.

---

## Quick Start

### 1. Read the Documentation

Start with **[docs/README.md](docs/README.md)** to understand:
- How the workflow eliminates requirements toil
- Why context files enable AI autonomy
- How specifications and implementation flow through source control
- What humans approve vs. what AI handles mechanically

### 2. Review the Sample Project

Explore **[sample-project/](sample-project/)** to see:
- What level of detail AI agents need in context files
- How to structure business rules, architecture docs, and test patterns
- How context files scale expert knowledge across all stories

### 3. Set Up Your Context Files

Create your own context files using the sample project as a template:

```
your-project/
├── context/
│   ├── business.md           # Your domain knowledge
│   ├── architecture.md       # Your system architecture
│   ├── testing.md            # Your test patterns
│   └── tech_standards.md     # Your technical standards
├── api/
│   └── specs/                # Your API interface specs
├── features/                 # Your Gherkin specifications
└── step_definitions/         # Your test implementations
```

### 4. Configure AI Agents

See **[docs/prompts/](docs/prompts/)** for each agent's:
- Purpose and mode (conversational vs. one-shot)
- Input requirements
- Expected outputs
- Example usage

### 5. Integrate with Your Systems

Connect AI agents to:
- **Ticket System** (Jira, Azure DevOps, Linear, GitHub Issues)
- **Source Control** (GitHub, GitLab, Bitbucket)
- **CI/CD** (GitHub Actions, GitLab CI, Azure Pipelines)

See [docs/README.md § Integration Architecture](docs/README.md#3-integration-architecture) for webhook setup and API requirements.

---

## Using This System with Claude

See **[CLAUDE.md](CLAUDE.md)** for meta-instructions on using Claude to implement this workflow, including:

- AI usage principles (deterministic artifacts, not autonomous decisions)
- Available AI agents and when to use each
- Execution logs and audit trail requirements
- Context file governance
- Best practices for Business Owners, Developers, and Tech Leads

**Important:** All AI agent executions must be logged for audit, learning, and improvement purposes.

---

## Key Features

### AI-Powered Analysis

- **Reads package management** to discover internal dependencies (`go.mod`, `package.json`, etc.)
- **Validates against API specs** (OpenAPI, GraphQL, protobuf) to catch missing parameters and type mismatches
- **Generates boundary conditions** from API constraints + QA patterns in `testing.md`
- **Applies fuzzing patterns** for comprehensive edge case coverage
- **Escalates with specific questions** when information is missing or ambiguous

### Source Control Integration

- **Specifications live with code** - Same repo, atomic commits, no version drift
- **Branch protection enforces approval** - BO must approve all `.feature` file changes
- **Tags track status** - `@pending` for incomplete, `@story-{id}` for traceability
- **CI blocks incomplete work** - Implementation PRs can't merge with `@pending` scenarios for that story

### Context File System

- **Four files enable AI autonomy:**
  - `business.md` - Domain knowledge, business rules, personas, compliance
  - `architecture.md` - External dependencies, system constraints, third-party APIs
  - `testing.md` - Step library, boundary patterns, fuzzing references, edge cases
  - `tech_standards.md` - Language conventions, IoC patterns, directory structure

- **Versioned and governed** - Semantic versioning, CODEOWNERS, freshness monitoring
- **Shared across projects** - Organization-level, domain-level, team-level, project-level layering

### Confidence-Based Workflow

| AI Confidence | Action | Human Time |
|---------------|--------|------------|
| **High** | AI drafts PR directly | ~5 min review |
| **Medium** | AI drafts with flagged assumptions | ~15 min review |
| **Low** | AI escalates with specific questions | Discussion required |

**Goal:** Most stories need only quick PR review—no meetings, no waiting for schedules to align.

---

## Architecture Highlights

### Two-Branch Model

1. **Spec branch** (`spec/{ticket-id}`): AI drafts Gherkin → BO approves → Merges with `@pending` tags
2. **Implementation branch** (`impl/{ticket-id}`): Developer implements → Removes `@pending` → CI validates → Merges

**Why separate:** Spec work completes when BO approves. Implementation can begin immediately from clean main branch.

### Tag Lifecycle

- **After spec merge:** `@pending @story-PROJ-1234` (two separate tags)
- **After implementation merge:** `@story-PROJ-1234` only (`@pending` removed)
- **CI enforcement:** Blocks implementation PR if `@pending` present for that story

### Escalation & Learning

When AI escalates due to missing information:
1. Experts answer (applying judgment to difficult problem)
2. **Answer gets codified into appropriate context file**
3. Future similar stories don't escalate—AI has the information

**This creates a feedback loop:** Escalations → human answers → documented → fewer future escalations.

---

## Repository Structure

```
.
├── README.md                          # This file - overview and navigation
├── CLAUDE.md                          # Meta-instructions for using Claude
├── docs/
│   ├── README.md                      # Complete workflow documentation
│   └── prompts/                       # AI agent prompts and examples
│       ├── requirements-drafting-assistant/
│       ├── requirements-analyst/
│       ├── bo-review/
│       ├── standards-compliance/
│       └── developer-implementation/
└── sample-project/
    ├── README.md                      # Sample project overview
    └── context/                       # Example context files
        ├── architecture.md
        ├── business.md
        ├── testing.md
        └── tech_standards.md
```

---

## Getting Help

- **For workflow questions:** See [docs/README.md](docs/README.md)
- **For context file examples:** See [sample-project/](sample-project/)
- **For using with Claude:** See [CLAUDE.md](CLAUDE.md)
- **For AI agent prompts:** See [docs/prompts/](docs/prompts/)

---

## License

[Specify your license here]

---

## Contributing

[Specify contribution guidelines here]

---

**Remember:** AI eliminates toil. Humans retain judgment. Requirements keep pace with development.
