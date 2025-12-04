# Requirements Drafting Assistant

## Purpose

Help Business Owners articulate requirements through conversational exploration, considering API contracts, business rules, and technical constraints.

## How It's Different

Unlike other prompts in this system (which are one-shot executions), this is a **conversational agent** that:
- Asks clarifying questions iteratively
- Explores edge cases through dialogue
- Validates against existing context (business.md, architecture.md)
- Builds requirements incrementally
- Produces structured output when BO confirms completeness

## When to Use

**Use this agent when:**
- Starting with vague requirement ("we need password reset")
- Complex feature with many edge cases
- Need to explore API constraints
- Want to validate against existing business rules
- BO not sure what all the requirements are yet

**Skip this agent when:**
- Requirement is already well-documented
- Simple change to existing feature
- Just need Gherkin spec from existing acceptance criteria

## Integration with Workflow

```
BO has idea/problem →
  requirements-drafting-assistant (conversational) →
    Structured requirement document →
      Create ticket (PROJ-XXX) →
        requirements-analyst (draft Gherkin spec) →
          bo-review (approve spec) →
            developer-implementation (code it)
```

This agent sits at the **beginning** of the workflow, before requirements-analyst.

## Key Capabilities

### 1. Context-Aware Questions

Agent reads context files and asks questions based on:
- **business.md:** Existing business rules, domain model, ubiquitous language
- **architecture.md:** API contracts, external dependencies, system boundaries
- **testing.md:** Existing step library, common patterns
- **tech_standards.md:** Technical constraints (Go, gRPC, etc.)

**Example:**
```
BO: "We need password reset"
Agent: "I see BR-003 requires 12+ character passwords for admins.
        Should reset enforce the same complexity rules?"
```

### 2. Edge Case Exploration

Agent proactively asks "what if" questions:
- Boundary conditions (empty, null, max values)
- Error scenarios (API down, timeout, invalid input)
- Concurrent access (race conditions)
- Security implications (enumeration, brute force)

**Example:**
```
Agent: "What if someone requests reset for email that doesn't exist?
        Should we reveal that, or return generic message to prevent
        email enumeration attacks?"
```

### 3. API Constraint Validation

Agent checks architecture.md for:
- Available API endpoints and contracts
- External service limitations (rate limits, timeouts)
- Authentication/authorization patterns
- Data format requirements

**Example:**
```
Agent: "I see the email service has rate limiting (100/hour).
        Should we queue emails, or fail fast and show error to user?"
```

### 4. Business Rule Consistency

Agent identifies:
- Conflicts with existing business rules
- Missing business rules that should exist
- Opportunities to reuse existing rules

**Example:**
```
Agent: "This requirement seems to conflict with BR-007 (data access controls).
        Should we update BR-007, or adjust this requirement?"
```

### 5. Structured Output

When conversation is complete, agent produces:
- User story (As a... I want... So that...)
- Acceptance criteria (testable behaviors)
- Business rules involved (existing + new)
- Edge cases discussed with expected behaviors
- API dependencies
- Open questions for tech team
- Next steps (ticket creation, routing)

## Sample Conversation

See `sample-conversation.md` for complete example of:
- BO starts with vague requirement
- Agent asks ~15 clarifying questions
- Explores edge cases, security, compliance
- Validates against context files
- Produces complete structured requirement

**Topic:** API Key Management for Developer Portal
**Duration:** ~20 minute conversation
**Result:** Complete requirement ready for Gherkin spec drafting

## Context Files Required

Place these in the same directory or parent directory:

- `../../business.md` - Business rules, domain model
- `../../architecture.md` - API contracts, system architecture
- `../../testing.md` - Test patterns, step library
- `../../tech_standards.md` - Technical constraints

Agent reads these first to inform questions.

## What Success Looks Like

**Good conversation:**
- BO starts with 1-2 sentence idea
- Agent asks 10-20 clarifying questions
- Edge cases identified and resolved
- API constraints checked early
- Output is detailed requirement (2-3 pages)
- BO says "I wouldn't have thought of those edge cases"

**Poor conversation:**
- Agent asks too few questions (output has gaps)
- Agent asks irrelevant questions (wastes BO time)
- Agent doesn't check context files (misses constraints)
- Output still has ambiguities or conflicts

## AI vs. Human Responsibilities

**AI (this agent) handles:**
- Asking systematic clarifying questions
- Checking context files for constraints
- Identifying common edge cases
- Structuring output in consistent format
- Flagging conflicts with existing rules

**Human (BO) handles:**
- Providing business context and goals
- Making business decisions
- Prioritizing requirements
- Approving final output
- Creating ticket and routing to next phase

## Tips for Business Owners

**Do:**
- Start with the problem you're solving, not the solution
- Answer questions based on business needs, not technical constraints
- Say "I don't know" when uncertain - agent will flag as open question
- Think about your users' goals and frustrations

**Don't:**
- Worry about technical feasibility (agent will check)
- Skip edge cases ("that probably won't happen") - agent will explore them
- Assume agent knows context you haven't shared - be explicit
- Rush to finish - thorough conversation saves rework later

## Example Execution

**See:** `example-output/CONVERSATION-LOG.md`

**Summary:** Real execution with Claude subagent helping BO draft API key management requirements. Conversation covered:
- Scope clarification (external partners, not internal users)
- Security constraints (key display, revocation, rate limiting)
- Edge cases (duplicate names, account deletion, concurrent revocations)
- API design (header format, validation location)
- Compliance (data retention, audit logging)
- UX considerations (error messages, confirmation flows)

**Result:** Complete requirement document ready for ticket creation and Gherkin spec drafting.

## Success Metrics

- **Coverage:** 90%+ of edge cases identified before implementation
- **Clarity:** Developers ask < 2 clarifying questions during implementation
- **Time saved:** Reduces rework from ambiguous requirements by 60%+
- **BO satisfaction:** "This helped me think through scenarios I would have missed"

## When to Escalate

**Escalate to Product Owner if:**
- Requirement conflicts with product roadmap
- Significant new feature outside BO authority
- Budget/timeline implications need approval

**Escalate to Tech Lead if:**
- Multiple API constraints conflict (agent can't resolve)
- Feasibility concerns (agent flags but can't determine)
- Architectural decision needed (new service, database, etc.)

**Escalate to Compliance if:**
- New data retention requirements
- PII handling questions
- Regulatory constraints unclear

## Integration Notes

**Input to this agent:**
- Vague requirement or problem statement
- Story ID (if exists)
- User persona (if known)

**Output from this agent:**
- Structured requirement document
- Ready for ticket creation (PROJ-XXX)
- Then routes to requirements-analyst prompt

**Time investment:**
- 15-30 minutes for simple feature
- 30-60 minutes for complex feature
- Worth it: prevents days of rework from unclear requirements
