# AI-Augmented Requirements Workflow - Documentation Index

**Behavior-Driven Development accelerated with AI, while humans retain control.**

---

## Quick Links

**New to this workflow?** Start here:
1. [Why BDD?](bdd-value.md) - Understand the value proposition
2. [Implementation Summary](implementation-summary.md) ⭐ - Get started in 1-2 days
3. [Workflow Guide](workflow.md) - Learn the complete process
4. [Roles & Responsibilities](roles.md) - Understand who does what

**Ready to implement?**
- [Implementation Summary](implementation-summary.md) - Complete setup guide (1-2 days)
- [Integration Architecture](integration.md) - MCP gateway setup
- [CI Configuration](ci-configuration.md) - Automate validation

---

## Executive Summary

### The Problem: Requirements Can't Keep Up

**AI-assisted development has broken the requirements bottleneck in the wrong direction.**

When developers use AI correctly—with architects and senior devs reviewing output, enforcing patterns, and bringing code into compliance—implementation velocity increases dramatically. A competent team with AI assistance can implement faster than traditional requirements processes can specify.

The result: requirements become the bottleneck. Teams either:

- **Slow down development** to wait for specifications (waste)
- **Skip specifications** and accumulate technical/business debt (risk)
- **Write specifications after implementation** which defeats the purpose (theater)

None of these are acceptable. The solution is to apply the same AI acceleration to requirements that we apply to development.

### The Solution: AI-Augmented Requirements

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

### AI Usage Principles: Deterministic Artifacts, Not Autonomous Decisions

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

### How It Works: Context Files Enable AI Autonomy

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

### Speed Targets

| Path | Time to PR Ready | Human Time Required |
|------|------------------|---------------------|
| AI Confident | Same day | ~5 min review |
| AI Confident + Inferences | Same day | ~15 min review |
| AI Escalates | 2-3 days | Discussion required |

**Goal:** Most stories need only a quick PR review—no meetings, no back-and-forth, no waiting for schedules to align. The proportion depends on how complete your context files and API specs are.

**Compare to traditional:** A Three Amigos session takes 30-60 minutes per story, requires synchronous availability of 3+ people, and still needs someone to write the Gherkin afterward. This workflow reduces human time per story from hours to minutes.

---

## Core Documentation

### Getting Started

**[Implementation Summary](implementation-summary.md)** ⭐ **START HERE**

Complete setup guide using IBM ContextForge MCP Gateway (open source). Covers installation, configuration, testing, and CI integration. **1-2 days total** vs. 3-4 weeks custom development.

**[BDD Value Proposition](bdd-value.md)**

Deep dive into why automated acceptance tests matter. How BDD prevents fragmentation, validates AI output, survives refactoring, and provides business-readable specifications. Includes 5 critical values and real-world examples.

### Understanding the Workflow

**[Workflow Guide](workflow.md)**

Complete end-to-end process from story creation through implementation. Includes:
- All 5 workflow stages (including optional Stage 0)
- Step definition lifecycle
- AI confidence & escalation process
- Workflow diagram and speed targets

**[Roles & Responsibilities](roles.md)**

Who does what, required access, training requirements, and time investment comparison. Includes:
- Product Owner, Business Owner, Developer, QA Lead, Tech Lead responsibilities
- MCP access requirements per role
- CODEOWNERS configuration
- Traditional vs AI-augmented time comparison

### Key Concepts

**[Context Files](context-files.md)**

Your team's common knowledge, externalized. How to structure business.md, architecture.md, testing.md, and tech_standards.md. Includes:
- Why context files enable AI autonomy
- Shared and organization-level context
- API interface requirements
- Context file governance
- Complete templates for all four files

**[AI Agents](agents.md)**

5 specialized agents: when to use each, what inputs they need, what outputs they produce. Includes:
- requirements-drafting-assistant (conversational exploration)
- requirements-analyst (auto-draft specs)
- bo-review (guide BO review)
- developer-implementation (TDD collaboration)
- standards-compliance (automated code review)

### Implementation Details

**[Integration Architecture](integration.md)**

Connecting AI agents to your infrastructure. Includes:
- System components and trigger mechanisms
- PR creation flow
- Status tracking via tags
- MCP gateway architecture
- Security considerations
- Webhook configuration examples

**[Source Control Strategy](source-control.md)**

Two-branch model (spec/ and impl/), tag taxonomy (@pending, @story-{id}), branch protection rules, merge requirements. Includes:
- Repository structure
- Branching model and workflow
- Tag lifecycle and taxonomy
- Rollback & revert procedures
- Hotfix path
- Naming conventions

**[CI Configuration](ci-configuration.md)**

Automated validation rules, blocking conditions, compliance checks. Includes:
- GitHub Actions, GitLab CI, Azure Pipelines examples
- Validation rules (required tags, @pending blocking)
- Running scenarios
- Standards compliance check
- Branch protection configuration

### Guides

**[Scenario Design](scenario-design.md)**

Writing good Gherkin scenarios. Includes:
- Feature file organization and size limits
- When to use Scenario Outline vs separate scenarios
- Background usage patterns
- Step writing rules (Given-When-Then discipline)
- Examples of well-written scenarios
- Anti-patterns to avoid

**[Metrics & SLAs](metrics.md)**

Measuring and maintaining requirements velocity. Includes:
- Primary metric: AI confidence rate
- Secondary metrics: Speed and quality
- Review timeline targets (SLAs)
- Enforcement mechanisms
- Dashboard examples
- Traditional vs AI-augmented comparison

**[Troubleshooting](troubleshooting.md)**

Common issues and solutions. Includes:
- High AI escalation rate
- Incorrect specs generated
- Agent connectivity issues
- Path filtering problems
- Rate limiting and caching issues
- Workflow performance problems

---

## Technical Reference

**[FOSS Components Evaluation](foss-components-evaluation.md)**

Comprehensive evaluation of open source MCP gateways. Includes:
- IBM ContextForge (recommended)
- Lasso MCP Gateway (security-focused)
- Docker MCP Gateway (container-native)
- Feature comparison matrix
- Cost-benefit analysis

**[MCP Server Evaluation](mcp-server-evaluation.md)**

Evaluation of official GitHub, GitLab, and Jira MCP servers against workflow requirements. Result: All requirements met by existing servers.

**[MCP Integration Requirements](mcp-integration-requirements.md)**

Technical requirements for MCP integration. Includes:
- Ticketing system access requirements
- Source control access for non-developers
- Shared credentials
- Security controls
- Role-specific access requirements

---

## Meta-Instructions

**[Using with Claude](../CLAUDE.md)**

Meta-instructions for implementing this workflow with Claude Code. Includes:
- AI usage principles (deterministic artifacts not autonomous decisions)
- Agent execution patterns
- Logging requirements (all executions must be logged)
- Best practices for each role
- Troubleshooting tips

---

## Sample Project

**[Sample Project](../sample-project/)**

Example context files demonstrating the patterns. Real-world business.md, architecture.md, testing.md, and tech_standards.md. Use as templates for your own project.

**[Agent Prompts](prompts/)**

Complete prompts for each AI agent with execution examples. Includes:
- requirements-drafting-assistant
- requirements-analyst
- bo-review
- developer-implementation
- standards-compliance

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

## Enforced Rules

| Rule | Enforcement | Rationale |
|------|-------------|-----------|
| `@story-{id}` on all scenarios | CI blocks merge | Traceability to business request |
| `@pending` blocks impl PR for that story | CI blocks merge | Prevents incomplete implementation |
| BO approval for .feature files | CODEOWNERS + branch protection | Business alignment guaranteed |
| Context files have CODEOWNERS | Branch protection | Accountable ownership |

---

## Speed Comparison

| Traditional Process | AI-Augmented Process |
|---------------------|----------------------|
| 30-60 min Three Amigos meeting | 5-15 min async PR review |
| + manual Gherkin writing | AI drafts specifications |
| + calendar coordination | Immediate processing |
| **Total: 2-4 hours per story** | **Total: 5-30 minutes per story** |

---

**Remember:** AI eliminates toil. Humans retain judgment. Requirements keep pace with development.
