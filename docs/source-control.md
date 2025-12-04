# Source Control Strategy

Feature files live in the same repository as the code. This document defines the directory structure, branching model, merge rules, tag taxonomy, and special procedures.

---

## Repository Structure

**Feature files live in the same repository as the code they specify.** This ensures specifications and implementation stay in sync, enables atomic commits that include both spec and code changes, and simplifies CI configuration.

### Directory Structure

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

---

## Branching Model

Use trunk-based development with short-lived branches. Spec branches merge after BO approval; implementation branches are separate.

### Branch Types

| Branch Pattern | Purpose | Lifecycle | Merge Target |
|----------------|---------|-----------|--------------|
| `main` | Production-ready specifications and code | Permanent | n/a |
| `spec/{ticket-id}` | AI-drafted Gherkin specifications | Merges after BO approval | `main` |
| `impl/{ticket-id}` | Step definition implementation | Merges after tests pass | `main` |
| `amend/{ticket-id}` | Changes to approved specs | Merges after BO approval | `main` |
| `hotfix-spec/{ticket-id}` | Emergency spec additions | < 1 day | `main` |

### Two-Branch Workflow

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

### Branch Protection Rules

Configure these rules on the main branch:

- **Require pull request reviews:** Minimum 1 approval (2 for post-escalation)
- **Require BO approval for feature files:** All `.feature` file changes require Business Owner sign-off (enforced via CODEOWNERS)
- **Require status checks:** CI validation must pass
- **Require linear history:** Squash or rebase merges only
- **Restrict force pushes:** No force push to main
- **Require CODEOWNERS review:** Context file changes need owner approval

---

## Merge Strategy

### For Story PRs

Use squash merge. The commit message should include the ticket ID and a summary of scenarios added plus implementation notes. This keeps main history clean and traceable.

---

## Conflict Resolution

When multiple stories touch the same feature file:

### Prevention

- AI checks for open PRs on same feature file before creating new PR
- If conflict detected, AI adds scenarios to existing PR instead of creating new one
- Use feature file size limits (max 15 scenarios) to reduce collision probability

### Resolution

When conflicts occur despite prevention:

1. First PR to merge wins
2. Conflicting PR rebases onto main
3. AI re-validates scenarios still make sense with merged changes
4. If AI detects logical conflict (not just textual), escalate to humans

---

## Tag Taxonomy

**Tags are the authoritative status of specifications.** Not ticket status, not branch status—the tags in the feature files themselves.

### Lifecycle Tags

| Tag | Meaning | Location | In CI? |
|-----|---------|----------|--------|
| `@pending` | Awaiting implementation (temporary) | `main` branch between spec merge and impl merge | No (skipped) |

**Tag lifecycle:**
- Spec PR merges: Scenarios tagged `@pending @story-{id}`
- Implementation completes: Developer removes `@pending` tags for that story
- Implementation PR: **CI blocks merge if `@pending` present for that story**
- After merge: Scenarios have only `@story-{id}` (no lifecycle tag)
- Business validation tracked in deployment/release process

### Required Tags

| Tag | Format | Purpose |
|-----|--------|---------|
| `@story-{id}` | `@story-PROJ-1234` | Links to originating ticket |

### Optional Tags

| Tag | Purpose |
|-----|---------|
| `@smoke` | Include in fast feedback suite |
| `@regression` | Include in full regression suite |
| `@hotfix` | Emergency addition, expedited review |
| `@wip` | Work in progress, skip in CI |

---

## Rollback & Revert Procedures

This section covers how to handle specifications that need to be reversed after merge.

### When to Revert

- **Incorrect scenarios:** Scenarios don't match intended behavior
- **Blocking development:** Scenarios prevent valid implementation
- **Business change:** Requirements changed after spec approval
- **Security issue:** Scenarios expose security vulnerability in test data

### Revert Process

#### For @pending Scenarios (Not Yet Implemented)

- Create `amend/{ticket-id}-revert` branch
- Remove or modify scenarios
- PR requires 1 approval (original author cannot approve)
- BO must approve if feature files change
- Link to original ticket with reason

#### For Implemented Scenarios

- Create ticket for revert with business justification
- Revert PR must include both scenario removal AND step cleanup
- Requires 2 approvals (including BO)
- Update context files if revert reveals gap

### Amendment vs. Revert

| Situation | Action | Branch Pattern |
|-----------|--------|----------------|
| Minor correction to scenario | Amend | `amend/{ticket-id}-fix` |
| Add missing edge case | Amend | `amend/{ticket-id}-edge` |
| Completely wrong approach | Revert + new spec | `amend/{ticket-id}-revert` |
| Feature cancelled | Revert | `amend/{ticket-id}-cancel` |

---

## Hotfix Path

For production incidents requiring immediate test coverage.

### Standard vs. Hotfix

| Aspect | Standard Process | Hotfix Process |
|--------|------------------|----------------|
| AI review | Required | Skip |
| BO review | Full review | Security + functionality only |
| Approvals | 1-2 depending on confidence | 1 approval (BO if available) |
| Tag | `@pending` | `@hotfix @pending` |

### Hotfix Process

1. Tag scenarios `@hotfix @pending @story-{id}`
2. CI allows `@hotfix` to merge with expedited review
3. Within 5 business days: Remove `@hotfix`, complete standard BO review
4. If not completed in 5 days: Escalate to team lead

---

## Naming Conventions

### Branches

| Purpose | Pattern | Example |
|---------|---------|---------|
| Specification | `spec/{ticket-id}` | `spec/PROJ-1234` |
| Implementation | `impl/{ticket-id}` | `impl/PROJ-1234` |
| Amendment | `amend/{ticket-id}` | `amend/PROJ-1234` |
| Hotfix | `hotfix-spec/{ticket-id}` | `hotfix-spec/PROJ-9999` |

**Branch name derived from ticket URL.** AI extracts the ticket ID automatically—no manual entry required.

### Files

Follow your framework's conventions. Common patterns:

| Type | Pattern | Notes |
|------|---------|-------|
| Features | `{domain}_{capability}.feature` | Or `{Domain}{Capability}.feature` for PascalCase conventions |
| Steps | `{domain}_steps.{ext}` | Location and naming per framework (e.g., `*_test.go`, `*Steps.java`) |

---

## Enforced Rules

| Rule | Enforcement | Rationale |
|------|-------------|-----------|
| `@story-{id}` on all scenarios | CI blocks merge | Traceability to business request |
| `@pending` blocks impl PR for that story | CI blocks merge | Prevents incomplete implementation |
| BO approval for .feature files | CODEOWNERS + branch protection | Business alignment guaranteed |
| Context files have CODEOWNERS | Branch protection | Accountable ownership |

---

## Source Control as Single Source of Truth

All specifications live in the same repository as the code, as Gherkin/Cucumber `.feature` files:

- **Stories link to feature files** via `@story-{id}` tags
- **Tags track status:** `@pending` (awaiting implementation, removed when implementation merges)
- **Specs merge to main before implementation** so developers always work from approved specs
- **Branch protection enforces BO approval** for all feature file changes
- **CI validates completeness** (no orphan drafts, all scenarios tagged)
- **History is auditable** (who approved what, when, why)

---

## Next Steps

- **For workflow details:** See [Workflow Guide](workflow.md)
- **For CI configuration:** See [CI Configuration Guide](ci-configuration.md)
- **For context file governance:** See [Context Files Guide](context-files.md)
