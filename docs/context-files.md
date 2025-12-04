# Context Files

Context files are the foundation of AI autonomy—and therefore the foundation of requirements speed. **Investment in context files pays compound returns:** every hour spent documenting your systems saves many hours of future requirements toil.

---

## Overview

Context files are your team's shared knowledge written down. They enable AI to do what Dev and QA would do in a discussion—if it has access to the same information.

**Start with these four:**
- `business.md` - Domain knowledge, business rules, personas, compliance
- `architecture.md` - External dependencies, system constraints, third-party APIs
- `testing.md` - Step library, boundary patterns, edge case patterns
- `tech_standards.md` - Language conventions, coding patterns, directory structure

**Note:** Context files can grow and shrink based on your needs. Split them when they get too large (e.g., split business.md into business-rules.md and compliance.md), merge them if too fragmented, or add specialized files as your domain requires (like security.md, deployment.md, or integrations.md). These are a reasonable starting point, not a rigid requirement.

**Important:** When you change context file structure, update agent prompts to reference the new files. For example, if you split business.md into business-rules.md and compliance.md, modify agent prompts to read both files instead of just business.md.

---

## The Four Files

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

---

## Shared and Organization-Level Context

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

---

## Why This Matters: Scaling Expertise, Not Replacing People

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

---

## Context File Governance

### Ownership via CODEOWNERS

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

### Freshness Monitoring

Implement automated staleness detection:

- **Warning threshold:** 30 days since last update
- **Alert threshold:** 60 days since last update
- **Mechanism:** Weekly scheduled CI job checks last-modified dates
- **Action:** Creates ticket for owner to review and update or confirm current

### Version Coupling

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

---

## API Interface Requirements

**Internal dependencies:** AI reads your package management files (`go.mod`, `package.json`, `pom.xml`, `requirements.txt`, etc.) to discover internal service dependencies and their interfaces. No need to duplicate this in `architecture.md`.

**External dependencies:** Third-party APIs and external services must be documented in `architecture.md` with links to their interface specs. AI can't discover these from package files.

### What Goes Where

| Dependency Type | Source of Truth | Example |
|-----------------|-----------------|---------|
| Internal services | Package management + linked specs | `go.mod` references `github.com/company/auth-service` |
| Internal APIs | OpenAPI/GraphQL specs in repo | `/api/specs/auth-api.yaml` |
| External APIs | `architecture.md` + linked specs | Stripe, Twilio, AWS services |
| Shared libraries | Package management | Internal packages, SDKs |

### Supported Interface Formats

AI reads machine-readable interface specs to detect missing parameters, validate types, and identify constraints:

| Format | File Types | AI Capabilities |
|--------|------------|-----------------|
| OpenAPI/Swagger | `.yaml`, `.json` | Required/optional params, types, enums, constraints |
| GraphQL | `.graphql`, `.gql` | Query/mutation fields, required args, types |
| Protocol Buffers | `.proto` | Message fields, required/optional, types |
| JSON Schema | `.json` | Property constraints, required fields, patterns |
| TypeScript Interfaces | `.ts`, `.d.ts` | Type definitions, optional properties |

### architecture.md Structure

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

### What AI Detects From Interfaces

When AI reads your API interfaces, it can identify:

- **Missing required parameters** in acceptance criteria
- **Invalid enum values** that don't match the schema
- **Type mismatches** (expecting string, spec says integer)
- **Undocumented endpoints** referenced in stories
- **Deprecated fields** that shouldn't be used
- **Constraint violations** (min/max, patterns, formats)

This catches issues that would otherwise surface during implementation—or worse, in production.

---

## Keeping Files Current

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

## Context File Templates

### business.md Template

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

### architecture.md Template

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

### testing.md Template

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

### tech_standards.md Template

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

## Next Steps

- **For workflow integration:** See [Workflow Guide](workflow.md)
- **For governance details:** See [Roles & Responsibilities](roles.md)
- **For implementation:** See [Implementation Summary](implementation-summary.md)
