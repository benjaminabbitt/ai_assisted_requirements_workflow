# Using This System with Claude

This document provides meta-instructions for using Claude to implement the AI-augmented requirements workflow described in `docs/README.md`.

---

## Core Principles

### AI Usage: Deterministic Artifacts, Not Autonomous Decisions

**Good use of AI: Generate deterministic artifacts** (code, specifications, processes) based on well-defined rules and context. AI excels at mechanical work—reading documentation, parsing APIs, applying patterns, drafting specifications.

**Poor use of AI: Making decisions with loose rulesets without supervision.** AI should not make business decisions, security tradeoffs, or compliance interpretations autonomously. These require human judgment.

**This system implements this principle:**
- ✅ **AI generates artifacts:** Gherkin specifications, step definition skeletons, boundary test cases, compliance reports, unit tests
- ✅ **Humans make decisions:** Business owners approve requirements, developers implement logic, security teams validate approaches
- ✅ **Rules are explicit:** Context files, API specs, and patterns define what AI should do
- ✅ **Supervision is built-in:** Every AI-drafted spec requires human review before merge
- ✅ **Escalation is mandatory:** When AI encounters ambiguity or missing rules, it stops and asks humans
- ✅ **Audit trails required:** All AI executions must be logged (see below)

**The result:** AI eliminates toil while humans retain control over judgment.

### AI-Generated Code Requires Human Software Engineering

**AI generates code and tests** (following patterns, conventions, TDD cycles), **but humans must apply software engineering principles to both:**

#### For Tests:
- **Maintainability:** Refactor test code for clarity, eliminate duplication, improve readability
- **Consistency:** Ensure test style matches team conventions, naming is coherent across test suite
- **Performance:** Optimize slow tests, parallelize where appropriate, reduce test fixture overhead
- **Completeness:** Review edge cases, add missing scenarios, validate business logic coverage
- **Quality:** Ensure assertions are meaningful, error messages are clear, setup/teardown is proper

#### For Implementation Code:
- **Security:** Review for vulnerabilities (injection attacks, authentication flaws, data exposure)
- **Performance:** Identify bottlenecks, optimize algorithms, choose correct data structures, reduce unnecessary allocations
  - **Use correct data structures:** Maps for lookups (O(1)), not slice iteration (O(n))
  - **Consider time complexity:** Is this O(n²) when it could be O(n)?
  - **Avoid repeated work:** Multiple passes over data, redundant computations
  - **Profile hot paths:** Identify actual bottlenecks with measurements, not assumptions
- **Consistency & Reuse:** Ensure same logic used throughout application, not fragmented
  - **Extract and reuse:** Common validation, business rules, transformations should exist once
  - **Avoid fragmentation:** Same email validation everywhere (not different rules in different paths)
  - **Prevent divergence:** Logic shouldn't vary by endpoint, UI path, or reimplementation
  - **Refactor duplicates:** AI may generate similar-but-different logic across sessions (context window exhaustion)
  - **Monitor for duplication:** Regularly scan codebase for repeated patterns, similar implementations
  - **Unified behavior:** User shouldn't get different results resubmitting same data through different code paths
- **Architecture & Abstractions:** Design interfaces and boundaries that will serve the system as it evolves
  - **Create appropriate abstractions:** Identify common patterns and extract reusable interfaces
  - **Guide system design:** Determine where layers should exist, what responsibilities belong where
  - **Evaluate architectural decisions:** Consider coupling, cohesion, testability, maintainability
  - **Design for extensibility:** Build abstractions that accommodate future requirements
- **Dependency Management:** Evaluate, select, and maintain external dependencies
  - **Assess dependency quality:** Security track record, maintenance activity, community health
  - **Evaluate licensing:** Ensure compatibility with project requirements
  - **Monitor dependency health:** Watch for security advisories, breaking changes, deprecations
  - **Minimize dependency footprint:** Avoid adding dependencies for trivial functionality
  - **Version management:** Understand transitive dependencies, version conflicts, upgrade paths
- **Code Monitoring & Evolution:** Continuously improve codebase health over time
  - **Watch for technical debt:** Identify accumulating complexity, outdated patterns, architectural drift
  - **Refactor proactively:** Don't wait for crisis—improve code structure as you work
  - **Track metrics:** Monitor build times, test performance, code coverage trends
  - **Review patterns:** Ensure new code follows established conventions, update standards when needed
- **Error Handling:** Validate error paths are comprehensive and user-friendly
- **Maintainability:** Simplify complex logic, extract functions, improve naming
- **Correctness:** Verify business logic matches requirements, edge cases handled properly

**Example workflow (collaborative, interactive):**
1. Developer guides AI to generate implementation following patterns
2. AI generates code; developer reviews interactively
3. **Developer collaborates with AI on engineering concerns:**
   - Security: "Does this expose sensitive data? AI, add input validation here"
   - Performance: "Is this N+1 query? AI, refactor to batch this"
   - Consistency: "Does this match our error handling pattern? AI, use our standard error wrapper"
   - Architecture: "Is this the right abstraction? AI, extract this interface"
4. **Iterative refinement:** Developer reviews → guides AI to refactor → reviews again
5. **For critical sections:** Developer writes code personally when needed
6. Developer validates the solution is not just correct, but well-engineered

**This is collaborative engineering, not pull request review.** Developer and AI work together iteratively, not sequentially.

**This is software engineering, not brick-laying.** AI can generate syntactically correct code that follows patterns, but humans ensure that code is:
- **Secure** (no vulnerabilities)
- **Performant** (efficient algorithms, proper resource management)
- **Consistent** (unified style, coherent architecture across entire codebase)
- **Maintainable** (clear intent, good abstractions, manageable complexity)
- **Correct** (actually solves the business problem)

**The work is real engineering:** Making tradeoff decisions, choosing appropriate abstractions, designing for change, considering edge cases AI missed, identifying security implications, optimizing hot paths, and ensuring the codebase remains coherent as it grows.

---

## Execution Examples for Documentation

### Requirement

**Generate execution examples when prompts change** to document agent behavior and expected outputs.

### What to Capture

**For conversational agents** (requirements-drafting-assistant):
- Full conversation transcript (every question and answer)
- Context file versions referenced
- API contracts consulted
- Edge cases discovered during conversation
- Final structured output
- Conversation duration and outcome

**For one-shot agents** (requirements-analyst, bo-review, standards-compliance, developer-implementation):
- Prompt input (story, ticket, code to review)
- Agent reasoning process (analysis steps, inferences made)
- Context file versions used
- API specs referenced and validation results
- Decisions made and rationale
- Escalation triggers (if any)
- Final output (spec draft, review report, compliance report)
- Execution time and confidence level

### Storage

**Example outputs:** `docs/prompts/{agent-name}/example-output/`
- Used for documentation and examples
- Demonstrates actual agent execution
- Shows expected output format
- Updated when prompt changes

### Purpose

1. **Documentation:** Show what agent outputs look like
   - "What does a HIGH confidence spec draft look like?"
   - "What questions does requirements-drafting-assistant ask?"
   - "What does a standards-compliance report contain?"

2. **Learning:** Identify patterns in AI escalations to improve context files
   - "AI escalates frequently on auth stories → architecture.md missing auth service docs"
   - "AI never uses XYZ step → testing.md step library incomplete"

3. **Improvement:** Refine prompts based on actual execution patterns
   - "AI generates too many boundary cases → adjust testing.md patterns"
   - "AI misinterprets this business rule → clarify business.md wording"

4. **Transparency:** Humans can review AI reasoning process
   - BO can see what assumptions AI made
   - Developers can understand why spec was drafted this way

---

## Available AI Agents

This system provides 5 specialized AI agents. Each has a specific purpose and operates in a specific mode.

### 1. requirements-drafting-assistant (Conversational)

**Purpose:** Help Business Owners articulate requirements through conversational exploration.

**Mode:** CONVERSATIONAL - Back-and-forth dialogue, iterative refinement

**When to use:**
- Requirement is vague or exploratory
- Complex feature with many edge cases
- Need to validate against API constraints and business rules
- BO wants help thinking through scenarios

**Input:**
- Initial requirement idea from BO
- Context files (business.md, architecture.md, testing.md, tech_standards.md)
- Existing feature files (*.feature - to understand conventions, reuse steps, identify dependencies)
- Ticketing system data (related tickets, comments/threads, dependencies)

**Process:**
1. Read existing feature files to understand conventions and related features
2. Pull ticketing system data (related work, comments, dependencies) to understand context
3. Agent asks clarifying questions (references context files, API contracts, existing features, past tickets)
4. BO answers, agent explores further
5. Agent identifies edge cases, security implications, technical constraints
6. Conversation builds up complete requirement incrementally
7. Agent produces structured requirement document

**Output:**
- Structured requirement with user story, acceptance criteria, edge cases, API dependencies
- Ready for ticket creation (PROJ-XXX)
- Feeds into next stage (requirements-analyst)

**Execution:**
```bash
# Using Claude
claude --prompt docs/prompts/requirements-drafting-assistant/prompt.md \
       --context sample-project/context/business.md,sample-project/context/architecture.md,sample-project/context/testing.md,sample-project/context/tech_standards.md \
       --mode conversational

# Initial message: "We need to add two-factor authentication for admin users."
# Agent will begin asking questions...
```

**Log output to:** `docs/prompts/requirements-drafting-assistant/example-output/CONVERSATION-LOG-{date}.md`

**See example:** `docs/prompts/requirements-drafting-assistant/example-output/CONVERSATION-LOG.md`

---

### 2. requirements-analyst (One-Shot)

**Purpose:** Analyze story tickets and draft Gherkin specifications.

**Mode:** ONE-SHOT - Reads ticket, analyzes, produces draft spec in single execution

**Trigger:** Runs when ticket is labeled "ready-for-spec" (webhook) or manually invoked. Creates Gherkin/Cucumber PR for BO review.

**Input:**
- User story and acceptance criteria from ticket
- Ticket ID (extracted from URL for branch naming)
- Context files (business.md, architecture.md, testing.md, tech_standards.md)
- Package management files (for internal dependencies)
- API interface specs (OpenAPI, GraphQL, etc.)

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
- Feature file with scenarios tagged `@story-{id}`
- Skeleton step definitions
- PR ready for BO review with confidence notes

**Execution:**
```bash
# Using Claude
claude --prompt docs/prompts/requirements-analyst/prompt.md \
       --context sample-project/context/business.md,sample-project/context/architecture.md,sample-project/context/testing.md,sample-project/context/tech_standards.md \
       --input "Ticket URL: https://jira.company.com/PROJ-1234"

# Agent reads ticket, analyzes, drafts spec
```

**Log output to:** `docs/prompts/requirements-analyst/example-output/OUTPUT-{ticket-id}.md`

**See example:** `docs/prompts/requirements-analyst/example-output/OUTPUT.md`

---

### 3. bo-review (One-Shot)

**Purpose:** Guide Business Owners in reviewing AI-drafted specifications for business correctness.

**Mode:** ONE-SHOT - Reads draft spec, performs review, produces report

**Trigger:** Runs on Gherkin/Cucumber PRs (spec/ branches) to assist BO in reviewing specifications before approval.

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

**Execution:**
```bash
# Using Claude
claude --prompt docs/prompts/bo-review/prompt.md \
       --context sample-project/context/business.md \
       --input "PR: spec/PROJ-1234 (draft spec for password reset)"

# Agent reviews spec, produces report
```

**Log output to:** `docs/prompts/bo-review/example-output/REVIEW-{ticket-id}.md`

**See example:** `docs/prompts/bo-review/example-output/REVIEW-OUTPUT.md`

---

### 4. developer-implementation (Conversational)

**Purpose:** Collaborate with developers in implementing step definitions for approved Gherkin specifications following TDD and IoC patterns.

**Mode:** CONVERSATIONAL - Interactive, iterative implementation with developer guidance using specialized subagents

**Input:**
- Approved spec (merged to main with `@pending @story-{id}` tags)
- Skeleton step definitions (auto-generated)
- tech_standards.md (TDD cycle, IoC patterns)
- architecture.md (external dependencies, APIs)

**Process:**
1. **Plan:** Introspect context files, dependencies, existing code
2. **Review:** Developer validates plan, adjusts approach
3. **Execute:** TDD implementation with multi-perspective review
   - Expert {Language} Developer subagent (idioms, standards)
   - Neophyte Developer subagent (clarity, documentation)
   - Specialty subagents (security, performance, concurrency, etc.)
4. Iterative refinement until production-ready
5. Remove `@pending` tags for this story (keep `@story-{id}`)
6. CI blocks implementation PR if `@pending` tags present for this specific story

**Tag lifecycle:**
- After spec merge: `@pending @story-PROJ-1234` (two separate tags)
- After implementation merge: `@story-PROJ-1234` only (@pending removed for that story)
- **CI/CD behavior:**
  - **General CI runs:** Skip all scenarios tagged `@pending` (incomplete implementations)
  - **Development:** Developer removes `@pending` for their story, enabling them to run those specific A/C scenarios during implementation
  - **Implementation PR blocking:** CI blocks merge if `@pending` tags present for that specific story ID

**Output:**
- Implementation guidance with examples
- Unit test examples (using primary constructors + mocks)
- Service implementation examples
- Step definition patterns

**Execution:**
```bash
# Using Claude
claude --prompt docs/prompts/developer-implementation/prompt.md \
       --context sample-project/context/tech_standards.md,sample-project/context/architecture.md,sample-project/context/testing.md \
       --input "Implement: features/auth/password_reset.feature (@pending)"

# Agent provides implementation guidance
```

**See example:** `docs/prompts/developer-implementation/prompt.md` (contains full conversational TDD example with subagents)

---

### 5. standards-compliance (One-Shot)

**Purpose:** Automated code review for compliance with project technical standards, especially IoC patterns.

**Mode:** ONE-SHOT - Runs automatically on PRs, reads code, checks standards, produces compliance report

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

**Execution:**
```bash
# Using Claude (typically automated in CI/CD)
claude --prompt docs/prompts/standards-compliance/prompt.md \
       --context sample-project/context/tech_standards.md,sample-project/context/architecture.md \
       --input "Review: internal/services/*.go"

# Agent analyzes code, produces compliance report
```

**Log output to:** `docs/prompts/standards-compliance/example-output/COMPLIANCE-{hash}.md`

**See example:** `docs/prompts/standards-compliance/example-output/COMPLIANCE-REPORT.md`

---

## Workflow Integration

### Complete Flow

**See `docs/README.md` for the canonical workflow description.**

Quick summary:
1. Story created → AI analyzes
2. AI drafts spec → Opens PR (spec/ branch)
3. BO reviews and approves → PR merges with `@pending @story-{id}` tags
4. Developer implements → Removes `@pending` for story → Tests pass → Merges
5. **CI/CD test behavior:**
   - General CI runs skip all `@pending` scenarios (incomplete implementations)
   - Implementation PRs are blocked if `@pending` tags present for that specific story
   - Developer runs A/C locally (without `@pending`) during implementation
6. Demo/release → Business validates (tracked in deployment)

### When to Use Each Agent

**requirements-drafting-assistant:**
- Before formal ticket creation
- When requirement needs exploration
- Optional but helpful for complex features

**requirements-analyst:**
- After ticket created
- Triggered automatically (webhook) or manually
- Required for every new feature

**bo-review:**
- After requirements-analyst produces draft
- Before spec merges to main
- Required - BO approval mandatory

**developer-implementation:**
- After spec approved and merged
- During implementation process
- Conversational guidance with developer throughout coding

**standards-compliance:**
- During implementation PR review (automated)
- Before code merges to main
- Runs in CI/CD to validate technical standards

---

## Context Files Required

All agents depend on well-maintained context files:

### business.md
- Domain glossary
- User personas
- Business rules
- Compliance requirements

**Used by:** requirements-drafting-assistant, requirements-analyst, bo-review

### architecture.md
- System overview
- Internal API specs (paths to OpenAPI/GraphQL files)
- External dependencies (third-party APIs)
- Technical constraints

**Used by:** requirements-drafting-assistant, requirements-analyst, standards-compliance, developer-implementation

### testing.md
- Step library (existing steps to reuse)
- Boundary condition patterns
- Fuzzing library references
- Edge case patterns

**Used by:** requirements-drafting-assistant, requirements-analyst, developer-implementation

### tech_standards.md
- Language & framework (Go, Godog, etc.)
- IoC patterns (primary constructor + production factory)
- Directory structure
- Coverage strategy
- Naming conventions

**Used by:** requirements-analyst, standards-compliance, developer-implementation

**Critical:** Keep context files current. Every outdated context file increases AI escalation rate and slows requirements process.

**Flexibility:** Context files can grow and shrink based on your needs. Split when too large, merge when too fragmented, add specialized files as needed. **When you change context file structure, update agent prompts to reference the new files.** For example, if you split business.md into business-rules.md and compliance.md, modify agent prompts to read both files.

---

## MCP Integration

Both requirements-drafting-assistant and requirements-analyst access ticketing systems and feature files via **MCP (Model Context Protocol)** with shared credentials.

**See:** [MCP Integration Requirements](docs/mcp-integration-requirements.md) for:
- Ticketing system access (read tickets, comments, related tickets, search)
- Source control access for non-developers (read `.feature` files, context files)
- Shared credential configuration
- Security and access control
- Development roadmap for MCP servers

**Key principle:** Credentials configured once, shared across both agents. Non-developers (Business Owners) get read-only access to specifications without full repository access (security: no code, credentials, secrets).

---

## Best Practices

### For Business Owners

1. **Use requirements-drafting-assistant for complex features**
   - Don't skip this step for vague requirements
   - Conversation will surface edge cases early

2. **Review bo-review reports thoroughly**
   - Don't just approve - understand the issues
   - Ask questions if rationale isn't clear

3. **Codify escalation answers into business.md**
   - When AI escalates with business rule questions, answer them
   - **Then add the answer to business.md** so AI (and future humans) have it
   - Not all answers need codification (one-off exceptions), but most do (general rules)
   - This creates feedback loop: escalations → answers → documented → fewer future escalations

4. **Maintain business.md actively**
   - AI escalation on business logic = missing context in business.md
   - Document new business rules immediately after they're decided
   - Context files are living documents of your expert decisions

### For Developers

1. **Guide AI code generation, then review and refactor**

   **Humans guide AI, stepping in to write critical sections:**
   - Provide examples for complex patterns AI should follow
   - Write critical security or performance-sensitive code personally
   - Direct AI with specific implementation strategies
   - Review and refactor AI-generated code for engineering quality

   **For tests:**
   - AI generates tests following TDD pattern with mocks
   - Human reviews test, verifies it fails for right reason
   - **Human refactors test code:** Improve naming, extract helpers, consolidate setup
   - **Human reviews coverage:** Ensure meaningful coverage, not just line count

   **For implementation:**
   - AI generates code following patterns (guided by human)
   - **Review for security:** Input validation, injection attacks, data exposure
   - **Review for performance:**
     - N+1 queries, unnecessary allocations
     - **Data structure choice:** Map for lookups, not slice iteration
     - Algorithm efficiency and time complexity
     - Profile hot paths to identify actual bottlenecks
   - **Review for consistency & reuse:**
     - Extract common logic (validation, business rules) into reusable functions
     - Ensure same logic used throughout (not reimplemented differently in each place)
     - Check for fragmented behavior (email validates one way in API, differently in UI)
     - Monitor codebase for duplication across multiple features
   - **Review for architecture & abstractions:**
     - Create appropriate interfaces for common patterns
     - Evaluate if new code belongs in existing abstractions or needs new ones
     - Consider coupling, cohesion, and long-term maintainability
   - **Review dependencies:**
     - Evaluate security, maintenance, and licensing of new dependencies
     - Monitor existing dependencies for advisories and deprecations
     - Minimize dependency footprint
   - **Monitor code health:**
     - Watch for technical debt accumulation
     - Track build and test performance over time
     - Refactor proactively as complexity grows
   - **Refactor:** Simplify complex logic, extract functions, improve clarity
   - **Validate correctness:** Edge cases, error paths, business logic

2. **Follow developer-implementation TDD guidance**
   - Guide AI to write test first (using primary constructor + mocks)
   - Guide AI to implement minimal code to pass
   - Review and guide AI to refactor, write critical sections personally
   - **Review for engineering quality, not just correctness**
   - **Remove `@pending` tags for your story during implementation** to enable running A/C scenarios locally
   - CI will block your PR if any `@pending` tags remain for your story ID

3. **Watch for context window exhaustion issues**
   - AI may generate inconsistent code across multiple sessions
   - Look for: Different naming conventions, duplicate logic, conflicting patterns, fragmented validation
   - **Critical:** Same logic should be reused, not reimplemented
   - **Example problem:** Email validation succeeds via API but fails via UI (different implementations)
   - Refactor to maintain consistency and reuse across the entire codebase

4. **Run standards-compliance before opening PR**
   - Fix violations before review
   - Don't ignore warnings

5. **Update testing.md when you create new steps**
   - Document step signature and purpose
   - Help AI reuse your steps in future specs

**Remember:** AI generates code that follows patterns and passes tests. You ensure code is secure, performant, consistent, and maintainable. This is engineering, not brick-laying.

### For Tech Leads

1. **Maintain tech_standards.md rigorously**
   - AI follows this exactly
   - Outdated standards = incorrect code generation

2. **Review standards-compliance reports in PRs**
   - Enforce the patterns
   - Don't let violations merge

3. **Monitor AI escalation rate and codify answers**
   - High escalation rate = incomplete context files
   - **When AI escalates: Answer the question, then update the context file**
   - Escalations are learning opportunities—capture expert decisions for future reuse
   - Track which escalations lead to context updates vs. one-off answers
   - Goal: Reduce escalation rate over time as context files mature

4. **Codify escalation resolutions into context**
   - Technical constraint clarifications → `tech_standards.md`
   - API behavior patterns → `architecture.md`
   - Security/performance patterns → relevant context file
   - This creates feedback loop: fewer future escalations on similar issues

---

## Logging Examples

### Conversational Agent Log Format

```markdown
# Conversation Log: {Feature Name}

**Date:** YYYY-MM-DD
**Agent:** requirements-drafting-assistant
**Participant:** Business Owner
**Duration:** ~20 minutes
**Outcome:** Complete requirement ready for ticketing

## Conversation Transcript

**BO:** "We need two-factor authentication for admin users."

**Agent:** "I can help with that. Let me ask some questions...
1. What prompted this? Security incident or proactive?
2. Which admin users? All elevated privileges?
3. When required? Every login or only sensitive operations?"

**BO:** "Proactive after security review..."

[Full conversation with ~15-20 exchanges]

## Final Output

[Structured requirement document with user story, acceptance criteria, edge cases]

## Metadata

- Context files referenced: business.md v1.2.0, architecture.md v1.3.0
- Business rules validated: BR-002, BR-003
- API contracts checked: /api/specs/auth-api.yaml
- Edge cases identified: 12 scenarios
- Open questions: 4 (flagged for tech team)
```

### One-Shot Agent Log Format

```markdown
# Execution Log: {Agent Name}

**Date:** YYYY-MM-DD
**Agent:** {agent-name}
**Input:** {ticket-id or file path}
**Confidence:** HIGH/MEDIUM/LOW

## Analysis Process

1. Read ticket PROJ-1234
2. Extracted acceptance criteria: [list]
3. Read business.md v1.2.0 for business rules
4. Read architecture.md v1.3.0 for auth API
5. Validated parameters against /api/specs/auth-api.yaml
6. Generated boundary conditions: empty email, invalid format, max length
7. Applied testing.md patterns: string boundaries, error scenarios
8. Assessed confidence: MEDIUM (2 assumptions flagged)

## Output

[Draft spec or review report]

## Metadata

- Execution time: 2.3 seconds
- Context versions: business.md v1.2.0, architecture.md v1.3.0
- APIs validated: auth-api.yaml, users-api.yaml
- Scenarios generated: 11 (8 core, 3 boundary)
- Assumptions flagged: 2 (rate limiting behavior, email service async)
- Escalation: No
```

---

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: AI-Augmented Requirements

on:
  issues:
    types: [labeled]
  pull_request:
    types: [opened, synchronize, reopened]
  push:
    branches: [main]

jobs:
  run-acceptance-tests:
    # General CI runs skip @pending scenarios (incomplete implementations)
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run acceptance tests (skip @pending)
        run: |
          # Run Godog/Cucumber with tag filter to skip @pending scenarios
          # Syntax may vary by framework (Godog shown here)
          go test ./... -v -tags=godog --godog.tags="~@pending"
          # This skips all scenarios tagged @pending
          # Implemented scenarios (@story-PROJ-1234 without @pending) will run

  draft-spec:
    if: github.event.label.name == 'ready-for-spec'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run requirements-analyst
        run: |
          claude --prompt docs/prompts/requirements-analyst/prompt.md \
                 --context docs/*.md \
                 --input "${{ github.event.issue.html_url }}" \
                 --output spec-draft.md

      - name: Save execution log
        run: |
          mkdir -p logs/requirements-analyst
          cp execution.log logs/requirements-analyst/${{ github.event.issue.number }}.md

      - name: Create spec PR
        run: |
          git checkout -b spec/${{ github.event.issue.number }}
          # ... create feature file from spec-draft.md
          # ... commit and push
          # ... create PR with gh cli

  compliance-check:
    runs-on: ubuntu-latest
    on:
      pull_request:
        paths: ['**/*.go']
    steps:
      - name: Run standards-compliance
        run: |
          claude --prompt docs/prompts/standards-compliance/prompt.md \
                 --context sample-project/context/tech_standards.md \
                 --input "$(git diff --name-only origin/main... | grep '\.go$')" \
                 --output compliance-report.md

      - name: Comment on PR if violations
        run: |
          # ... post compliance-report.md as PR comment if score < 100%

  block-pending-impl:
    runs-on: ubuntu-latest
    if: startsWith(github.head_ref, 'impl/')
    on:
      pull_request:
        types: [opened, synchronize, reopened]
    steps:
      - uses: actions/checkout@v4

      - name: Block impl PR with @pending for this story
        run: |
          # Extract story ID from branch name (e.g., impl/PROJ-1234 -> PROJ-1234)
          STORY_ID=$(echo "${{ github.head_ref }}" | sed 's/^impl\///')

          # Check if any scenarios have both @pending and @story-{id} tags
          if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
            echo "ERROR: Implementation PR contains @pending scenarios for story $STORY_ID"
            echo ""
            echo "The following scenarios still have @pending tags:"
            grep -rn "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/
            echo ""
            echo "Remove @pending tags for this story before merging (keep @story-$STORY_ID)"
            echo ""
            echo "Note: @pending and @story-{id} are two separate tags."
            echo "After implementation, only @story-$STORY_ID should remain."
            exit 1
          fi

          echo "✓ No @pending tags found for story $STORY_ID"
```

---

## Troubleshooting

### High AI Escalation Rate

**Symptom:** AI frequently escalates with "missing context" or "unclear requirement"

**Root Cause:** Incomplete context files or escalation answers not being codified

**Fix:**
1. **Review recent escalations** - what's missing?
2. **For each escalation, identify pattern:**
   - One-off question (story-specific) → Answer inline, don't codify
   - General pattern/rule → **Answer AND codify into context file**
3. **Codify answers into appropriate context file:**
   - Business rules → `business.md`
   - External API behavior → `architecture.md`
   - Edge case patterns → `testing.md`
   - Technical standards → `tech_standards.md`
4. **Bump context file version** when you add content
5. **Re-run escalated stories** - AI should now have the information

**Key principle:** Escalations are learning opportunities. Experts apply judgment to difficult problems, then capture decisions in context files for future reuse. This creates feedback loop that reduces escalation rate over time.

### Incorrect Specs Generated

**Symptom:** AI drafts specs that don't match intent

**Root Cause:** Ambiguous context files or outdated API specs

**Fix:**
1. Review the execution log - what did AI reference?
2. Check if business rules in business.md are clear
3. Verify API specs match actual API implementation
4. Add examples to context files for similar scenarios

### Tests Fail After Implementation

**Symptom:** Scenarios pass but unit tests show violations

**Root Cause:** Developer didn't follow IoC patterns or TDD guidance

**Fix:**
1. Run standards-compliance agent on code
2. Review developer-implementation guidance
3. Ensure tests use primary constructors + mocks
4. Verify production factories have `// coverage:ignore` and no business logic

---

## Version History

- **v1.0.0** (2025-12-04): Initial version
  - 5 AI agents defined
  - Execution logging requirements specified
  - AI usage principles documented
  - Integration guidance provided
