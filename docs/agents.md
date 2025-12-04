# AI Agents

This workflow provides 5 specialized AI agents. Each has a specific purpose, operates in a specific mode, and requires specific inputs.

For detailed meta-instructions on using these agents with Claude, see [CLAUDE.md](../CLAUDE.md).

---

## Overview

| Agent | Mode | Purpose | When to Use |
|-------|------|---------|-------------|
| **requirements-drafting-assistant** | Conversational | Help BO articulate complex requirements | Vague/complex features needing exploration |
| **requirements-analyst** | One-Shot | Analyze tickets, draft Gherkin specs | Every new story (automatic or manual) |
| **bo-review** | One-Shot | Guide BO in reviewing AI-drafted specs | When reviewing Gherkin/Cucumber PRs |
| **developer-implementation** | Conversational | Collaborate on TDD implementation | During step definition implementation |
| **standards-compliance** | One-Shot | Review code for technical standards | Automated on implementation PRs |

---

## 1. requirements-drafting-assistant

**Mode:** CONVERSATIONAL - Back-and-forth dialogue, iterative refinement

**Purpose:** Help Business Owners articulate requirements through conversational exploration.

**When to use:**
- Requirement is vague or exploratory ("we need better security")
- Complex feature with many edge cases
- Need to validate against API constraints and business rules
- BO wants help thinking through scenarios

**Input:**
- Initial requirement idea from BO
- Context files (business.md, architecture.md, testing.md, tech_standards.md)
- Access to tickets (via MCP) for historical context
- Access to .feature files (via MCP) for existing conventions

**Process:**
1. Agent asks clarifying questions (references context files, API contracts)
2. BO answers, agent explores further
3. Agent identifies edge cases, security implications, technical constraints
4. Conversation builds up complete requirement incrementally
5. Agent produces structured requirement document

**Output:**
- Structured requirement with user story, acceptance criteria, edge cases, API dependencies
- Ready for ticket creation (PROJ-XXX)
- Feeds into requirements-analyst

**Example usage:**
```bash
claude --agent requirements-drafting-assistant
> "We need two-factor authentication for admin users"

# Agent begins asking questions...
```

**See:** [Prompt and Examples](prompts/requirements-drafting-assistant/)

---

## 2. requirements-analyst

**Mode:** ONE-SHOT - Reads ticket, analyzes, produces draft spec in single execution

**Purpose:** Analyze story tickets and draft Gherkin specifications automatically.

**Trigger:** Automatic (webhook on ticket labeled "ready-for-spec") or manual

**Input:**
- User story and acceptance criteria from ticket
- Ticket ID (extracted from URL for branch naming)
- Context files (business.md, architecture.md, testing.md, tech_standards.md)
- Package management files (for internal dependencies)
- API interface specs (OpenAPI, GraphQL, protobuf, etc.)
- Access to tickets (via MCP) for related context
- Access to .feature files (via MCP) for step reuse

**Process:**
1. Extract ticket ID from URL
2. Read package management to discover internal dependencies
3. Read architecture.md for external dependencies
4. Read API interface specs to validate parameters
5. Analyze story against all context files
6. Generate boundary conditions from API constraints + testing.md patterns
7. Assess confidence level
8. Either draft Gherkin or escalate

**Output:**
- Branch: `spec/{ticket-id}`
- Feature file with scenarios tagged `@pending @story-{id}`
- Skeleton step definitions
- PR ready for BO review with confidence notes
- Creates Gherkin/Cucumber PR (triggers bo-review)

**Example usage:**
```bash
claude --agent requirements-analyst \
       --input "Ticket: https://jira.company.com/PROJ-1234"
```

**See:** [Prompt and Examples](prompts/requirements-analyst/)

---

## 3. bo-review

**Mode:** ONE-SHOT - Reads draft spec, performs review, produces report

**Purpose:** Guide Business Owners in reviewing AI-drafted specifications for business correctness.

**Trigger:** Runs on Gherkin/Cucumber PRs (spec/ branches containing .feature files)

**Input:**
- AI-drafted specification (feature file in spec/ branch)
- Original story/ticket
- Context files (business.md for business rules validation)
- API contracts (to validate technical feasibility)

**Process:**
1. Check if scenarios match story intent
2. Verify business rules are applied correctly
3. Validate error messages are appropriate
4. Identify missing critical scenarios
5. Flag security/compliance concerns
6. Assess AI assumptions and inferences

**Output:**
- Review report with APPROVED or CHANGES REQUESTED decision
- List of issues found (critical, warning)
- Questions for clarification
- Recommendation for next steps

**Example usage:**
```bash
claude --agent bo-review \
       --input "PR #123" \
       --context business.md
```

**See:** [Prompt and Examples](prompts/bo-review/)

---

## 4. developer-implementation

**Mode:** CONVERSATIONAL - Interactive, iterative implementation with developer guidance

**Purpose:** Collaborate with developers in implementing step definitions for approved Gherkin specifications following TDD and IoC patterns.

**When to use:**
- After spec approved and merged (scenarios tagged `@pending @story-{id}`)
- Developer needs guidance on TDD implementation
- Developer wants AI collaboration on code generation

**Input:**
- Approved spec (merged to main with `@pending @story-{id}` tags)
- Skeleton step definitions (auto-generated)
- tech_standards.md (TDD cycle, IoC patterns)
- architecture.md (external dependencies, APIs)

**Process:**
1. Review spec scenarios (tagged `@pending @story-{id}`)
2. Guide TDD implementation (write test first, implement, test passes, refactor)
3. Ensure services use primary constructors
4. Verify production factories have `// coverage:ignore` and no business logic
5. Run scenarios until passing
6. Remove `@pending` tags for this story (keep `@story-{id}`)

**Output:**
- Implementation guidance with examples
- Unit test examples (using primary constructors + mocks)
- Service implementation examples
- Step definition patterns

**Tag lifecycle:**
- After spec merge: `@pending @story-PROJ-1234` (two separate tags)
- After implementation merge: `@story-PROJ-1234` only (@pending removed for that story)
- CI enforcement: Story-specific blocking (not global)

**Example usage:**
```bash
claude --agent developer-implementation \
       --input "Implement: features/auth/password_reset.feature (@pending)"
```

**See:** [Prompt and Examples](prompts/developer-implementation/)

---

## 5. standards-compliance

**Mode:** ONE-SHOT - Reads code, checks standards, produces compliance report

**Purpose:** Review Go code for compliance with project technical standards, especially IoC patterns.

**Trigger:** Runs automatically on implementation PRs (impl/ branches)

**Input:**
- Go code files to review
- tech_standards.md (IoC patterns, conventions)
- architecture.md (for architecture constraints)

**Process:**
1. Check for missing primary constructors
2. Detect business logic in production factories
3. Verify `// coverage:ignore` markers on factories
4. Check tests use primary constructors (not production factories)
5. Validate godoc comments on exported items
6. Calculate compliance score

**Output:**
- Compliance report with score (0-100%)
- Critical violations (must fix)
- Warnings (should fix)
- Specific fixes with code examples
- Per-service compliance checklist

**Example usage:**
```bash
claude --agent standards-compliance \
       --input "Review: internal/services/*.go" \
       --context tech_standards.md,architecture.md
```

**See:** [Prompt and Examples](prompts/standards-compliance/)

---

## Agent Execution & Logging

**IMPORTANT:** All AI agent executions MUST be logged for audit, learning, and improvement purposes.

**For conversational agents:**
- Full conversation transcript
- Context file versions referenced
- API contracts consulted
- Edge cases discovered
- Final structured output
- Conversation duration and outcome

**For one-shot agents:**
- Prompt input
- Agent reasoning process
- Context file versions used
- API specs referenced
- Decisions made and rationale
- Escalation triggers (if any)
- Final output
- Execution time and confidence level

**See:** [CLAUDE.md - Execution Logs](../CLAUDE.md#execution-logs-and-audit-trail) for complete requirements.

---

## MCP Access Requirements

All agents require MCP (Model Context Protocol) access to ticketing systems and source control:

**requirements-drafting-assistant & requirements-analyst:**
- Ticketing system: Read tickets, comments, related tickets, search
- Source control: Read .feature files, context files (no code access)

**bo-review:**
- Source control: Read .feature files in PRs, read context files
- Ticketing system: Read original ticket for context

**developer-implementation:**
- Source control: Read .feature files, context files, code for context

**standards-compliance:**
- Source control: Read code files for review

**See:**
- [Implementation Summary](implementation-summary.md) - Quick setup with MCP gateway
- [MCP Integration Requirements](mcp-integration-requirements.md) - Technical details
- [FOSS Components Evaluation](foss-components-evaluation.md) - Gateway options

---

## Prompt Templates

Complete prompts for each agent are available in the `prompts/` directory:

- [requirements-drafting-assistant/prompt.md](prompts/requirements-drafting-assistant/prompt.md)
- [requirements-analyst/prompt.md](prompts/requirements-analyst/prompt.md)
- [bo-review/prompt.md](prompts/bo-review/prompt.md)
- [developer-implementation/prompt.md](prompts/developer-implementation/prompt.md)
- [standards-compliance/prompt.md](prompts/standards-compliance/prompt.md)

Each prompt includes:
- Detailed instructions
- Input/output specifications
- Confidence assessment criteria
- Escalation triggers
- Example executions

---

## Next Steps

- **For workflow integration:** See [Workflow Guide](workflow.md)
- **For MCP setup:** See [Implementation Summary](implementation-summary.md)
- **For meta-instructions:** See [CLAUDE.md](../CLAUDE.md)
