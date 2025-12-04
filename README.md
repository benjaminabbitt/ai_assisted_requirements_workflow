# AI-Augmented Requirements Workflow

**Behavior-Driven Development accelerated with AI, while humans retain control.**

---

## What This Is

A workflow system that applies AI acceleration to requirements—the same way teams already use AI for development.

**The problem:** AI-assisted development increases implementation velocity. Requirements processes (Three Amigos meetings, manual Gherkin writing, calendar coordination) can't keep up. Requirements become the bottleneck.

**The solution:** Humans make decisions (what to build, how to build it). AI implements those decisions (drafts specs, generates tests, writes code). Humans verify the implementation is correct (approval, engineering review).

**The result:** Requirements velocity matches development velocity. 2-4 hours per story → 5-30 minutes.

---

## Why BDD?

**When AI writes code fast, how do you know it's right?**

Traditional approaches fail:
- **Documentation** gets outdated immediately
- **Unit tests** break when you refactor (even if behavior is correct)
- **Code review** catches syntax errors but misses behavioral bugs

**BDD solves this:** Business-readable scenarios that run as tests.

```gherkin
@story-PROJ-1234
Scenario: User requests password reset successfully
  Given a user exists with email "user@example.com"
  When I request a password reset for "user@example.com"
  Then I should receive a reset email
  And the password reset attempt should be logged
```

**What this gives you:**
- Business Owner reads it and knows what you're building
- Runs automatically → validates behavior every commit
- AI-generated code must make it pass → guarantees correctness
- Refactor all you want → as long as scenarios pass, behavior is preserved

**For AI-augmented development:**
1. BO approves: "Yes, this describes what I want"
2. AI generates code until scenarios pass
3. Scenarios pass = behavior is correct (not just compiles)

Think of it as a contract between business intent and code reality. The contract runs and tells you if it's honored.

[Deep dive: 5 ways BDD prevents problems in AI-augmented development →](docs/bdd-value.md)

---

## Core Benefits

### For Business Owners
- **Async review instead of synchronous meetings** - Review specifications at your convenience
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
   - Helps Business Owners articulate requirements through guided dialogue
   - Reads existing feature files, pulls ticketing system data (related tickets, comments/threads, dependencies)
   - Explores edge cases, validates against constraints

2. **requirements-analyst** (One-Shot)
   - Analyzes story tickets and drafts Gherkin specifications
   - Validates against API contracts, generates boundary conditions

3. **bo-review** (One-Shot)
   - Guides Business Owners in reviewing AI-drafted specifications
   - Checks business logic, identifies missing scenarios

4. **developer-implementation** (Conversational)
   - Collaborates with developers on TDD implementation using subagents
   - Plan → Review → Execute model with context introspection
   - Multi-perspective review: Expert developer, neophyte, security, performance, etc.

5. **standards-compliance** (One-Shot)
   - Automated code review for technical standards compliance
   - Runs automatically on PRs before human review

### Context Files: Your Team's Common Knowledge, Externalized

Your team's shared knowledge written down in four files:
- `business.md` - Domain knowledge, business rules, personas, compliance
- `architecture.md` - External dependencies, system constraints, third-party APIs
- `testing.md` - Step library, boundary patterns, edge case patterns
- `tech_standards.md` - Language conventions, coding patterns, directory structure

**Why:** AI (and new team members) read these to understand how your team works. Knowledge that was in people's heads is now written down and reusable.

[Read more about context files →](docs/context-files.md)

---

## Quick Start

1. **Read the workflow** → [docs/workflow.md](docs/workflow.md)
2. **Review example context files** → [sample-project/](sample-project/)
3. **Set up your context files** → [docs/context-files.md](docs/context-files.md)
4. **Configure AI agents** → [docs/agents.md](docs/agents.md)
5. **Integrate with your systems** → [docs/integration.md](docs/integration.md)

---

## Documentation

### Core Concepts

**[Workflow Overview](docs/workflow.md)**
Complete end-to-end process: story creation → AI drafts spec → BO approval → developer implementation → CI validation. Includes stage-by-stage flow, decision points, and escalation paths.

**[BDD Value Proposition](docs/bdd-value.md)**
Deep dive into why automated acceptance tests matter. How BDD prevents fragmentation, validates AI output, survives refactoring, and provides business-readable specifications. Includes 5 critical values and real-world examples.

**[Context Files](docs/context-files.md)**
Your team's common knowledge, externalized. How to structure business.md, architecture.md, testing.md, and tech_standards.md. Versioning strategy, governance model, and freshness monitoring.

**[AI Agents](docs/agents.md)**
5 specialized agents: when to use each, what inputs they need, what outputs they produce. Includes conversational vs one-shot modes, confidence levels, and escalation triggers.

### Implementation

**[Integration Architecture](docs/integration.md)**
Connecting AI agents to your infrastructure. Ticket system webhooks, source control APIs, CI/CD pipelines. Authentication, rate limiting, error handling, and retry logic.

**[Source Control Strategy](docs/source-control.md)**
Two-branch model (spec/ and impl/), tag taxonomy (@pending, @story-{id}), branch protection rules, merge requirements. Why specs merge before implementation, and how CI enforces completeness.

**[CI Configuration](docs/ci-configuration.md)**
Automated validation rules, blocking conditions, compliance checks. Running scenarios with @pending skipped, blocking implementation PRs, standards-compliance integration, and test result reporting.

**[Roles & Responsibilities](docs/roles.md)**
Who does what: Product Owners write stories, Business Owners approve specs, Developers implement, QA Leads maintain patterns, Tech Leads govern standards. Context file ownership via CODEOWNERS.

### Guides

**[Using with Claude](CLAUDE.md)**
Meta-instructions for implementing this workflow with Claude Code. AI usage principles (deterministic artifacts not autonomous decisions), agent execution patterns, logging requirements, and best practices.

**[Scenario Design](docs/scenario-design.md)**
Writing good Gherkin scenarios. Level of detail guidelines, when to split scenarios, handling edge cases, reusing steps, and maintaining readability. Examples of good vs problematic scenarios.

**[Metrics & SLAs](docs/metrics.md)**
Measuring and maintaining requirements velocity. Target SLAs (5-30 min per story), tracking escalation rates, context file freshness, and spec approval times. Dashboard examples and alerts.

**[Troubleshooting](docs/troubleshooting.md)**
Common issues and solutions. High AI escalation rates, incorrect specs, implementation blockers, CI failures. Includes diagnostic steps and resolution patterns.

### Reference

**[Sample Project](sample-project/)**
Example context files demonstrating the patterns. Real-world business.md, architecture.md, testing.md, and tech_standards.md. Use as templates for your own project.

**[Agent Prompts](docs/prompts/)**
Complete prompts for each AI agent with execution examples. Requirements-drafting-assistant, requirements-analyst, bo-review, developer-implementation, and standards-compliance. Includes example conversations and outputs.

---

## Key Principle

**Humans make decisions. AI implements them. Humans verify the implementation is correct.**

- ✅ **Humans make decisions:**
  - What to build (requirements, features, priorities)
  - How to build it (architecture, patterns, tradeoffs)
  - What "correct" means (acceptance criteria, quality standards)

- ✅ **AI implements decisions:**
  - Generates specifications from requirements
  - Generates test cases from acceptance criteria
  - Generates code following established patterns

- ✅ **Humans verify implementation and apply engineering:**
  - Does the spec match our requirements? (BO approval)
  - Does the code match the spec? (Developer review)
  - Is the engineering sound? (Security, performance, architecture, maintainability)
  - Does this fit the bigger picture? (System design, business strategy, team knowledge)

**Not:** AI makes decisions or generates unverified implementations
**Instead:** AI eliminates implementation toil; humans decide and verify

---

## Speed Comparison

| Traditional Process | AI-Augmented Process |
|---------------------|----------------------|
| 30-60 min Three Amigos meeting | 5-15 min async PR review |
| + manual Gherkin writing | AI drafts specifications |
| + calendar coordination | Immediate processing |
| **Total: 2-4 hours per story** | **Total: 5-30 minutes per story** |

---

## License

[Specify your license here]

## Contributing

[Specify contribution guidelines here]

---

**Remember:** AI eliminates toil. Humans retain judgment. Requirements keep pace with development.
