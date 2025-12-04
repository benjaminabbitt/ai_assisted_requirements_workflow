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

### Context Files: Your Team's Common Knowledge, Externalized

**What they are:** Your team's shared knowledge written down in context files.

**Start with these four:**
- `business.md` - Domain knowledge, business rules, personas, compliance
- `architecture.md` - External dependencies, system constraints, third-party APIs
- `testing.md` - Step library, boundary patterns, edge case patterns
- `tech_standards.md` - Language conventions, coding patterns, directory structure

**Why:** AI (and new team members) read these files to understand how your team works. Knowledge that was in people's heads is now written down and reusable.

**Flexibility:** Context files can grow and shrink. Split when too large, merge when too fragmented, add specialized files (like security.md, deployment.md, or integrations.md) as needed. These are a starting point, not a rigid structure. **When context files change, update agent prompts to reference the new files.**

**Governance:** Version them, assign owners, keep them current.

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
