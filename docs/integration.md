# Integration Architecture

**Speed requires automation.** Manual triggers, copy-paste workflows, and human-mediated handoffs reintroduce the toil this workflow eliminates. This document defines the integration points that enable fully automated story-to-specification flow.

---

## System Components

| Component | Role | Examples |
|-----------|------|----------|
| Ticket System | Story intake, status tracking, escalation threads | Jira, Azure DevOps, Linear, GitHub Issues |
| AI Agent | Analyzes stories, drafts Gherkin, creates PRs | Custom agent, GitHub Copilot Workspace, Claude |
| Source Control | Stores specifications, context files, step definitions | GitHub, GitLab, Azure Repos, Bitbucket |
| CI/CD | Validates specs, runs tests, enforces rules | GitHub Actions, GitLab CI, Azure Pipelines |

---

## Trigger Mechanism

AI review begins when a ticket meets trigger criteria.

### Webhook Trigger (Recommended)

1. Ticket system fires webhook on ticket labeled "ready-for-spec"
2. Webhook payload includes ticket ID (used as branch name), title, description, acceptance criteria
3. AI agent service receives webhook, fetches context files from repo
4. AI performs analysis and either creates branch + PR or posts escalation to ticket

### Manual Invocation (Fallback)

For teams not ready for webhook automation: user tags AI agent in ticket comment or runs CLI command with ticket ID.

---

## PR Creation Flow

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

---

## Status Tracking via Tags

**Tags are the source of truth for specification status—not ticket or branch status.** Ticket systems can track their own workflow, but the authoritative state lives in the feature files themselves.

| Tag | Meaning | Where It Lives |
|-----|---------|----------------|
| `@story-{id}` | Links to originating ticket | All scenarios |
| `@pending` | Awaiting implementation (removed after implementation merges) | `main` branch after spec PR merges, removed after impl PR merges |

**@pending lifecycle:** Added when spec merges, removed when implementation for that story merges. CI blocks merge if implementation PR contains `@pending` scenarios for that specific story.

**Ticket status synchronization is optional.** If your team wants ticket status to reflect spec progress, configure webhooks on PR merge events. But the tags in the feature files are what CI enforces and what matters for the workflow.

---

## Required API Integrations

Your AI agent needs API access to:

- **Ticket System:** Read tickets (including extracting ID from URL), update status, post comments
- **Source Control:** Read context files, create branches, commit files, open PRs
- **CI System (optional):** Query validation results for pre-merge checks

---

## MCP Gateway Architecture

**Recommended:** Use an MCP (Model Context Protocol) gateway to provide AI agents with consistent, authenticated access to external systems.

**See:**
- [Implementation Summary](implementation-summary.md) - Quick setup with IBM ContextForge or Lasso
- [FOSS Components Evaluation](foss-components-evaluation.md) - Detailed comparison of MCP gateways
- [MCP Integration Requirements](mcp-integration-requirements.md) - Technical requirements
- [MCP Server Evaluation](mcp-server-evaluation.md) - Evaluation of GitHub, GitLab, Jira MCP servers

**Key features needed:**
- Path filtering (block .env, secrets, keys)
- Authentication (OAuth, SSO)
- Authorization (RBAC for Business Owners vs Developers)
- Rate limiting
- Caching
- Monitoring

---

## Security Considerations

### Access Control

- **Business Owners:** Read-only access to .feature files, context files, and tickets (no code access)
- **AI Agents:** Read access to tickets and .feature files via MCP; write access to create branches/PRs via git API
- **Developers:** Full repository access

### Path Restrictions

Block sensitive files from MCP access:
- `.env*` files
- `secrets/**/*`
- `**/*.key`, `**/*.pem`, `**/*.p12`
- `**/id_rsa*`, `**/id_ed25519*`
- `.git/config` (may contain embedded credentials)

**See:** [MCP Integration Requirements](mcp-integration-requirements.md#security-considerations) for complete list.

---

## Webhook Configuration Examples

### GitHub Actions

```yaml
name: Requirements Automation

on:
  issues:
    types: [labeled]

jobs:
  draft-spec:
    if: github.event.label.name == 'ready-for-spec'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run requirements-analyst
        run: |
          claude --agent requirements-analyst \
                 --input "Ticket: ${{ github.event.issue.html_url }}"
```

### Jira Webhook

Configure webhook in Jira:
- **URL:** `https://your-agent-service/webhook/jira`
- **Trigger:** Issue updated
- **Filter:** `status = "Ready for Specification"`
- **Payload:** Include issue details, comments, related issues

---

## Rate Limiting & Retry Logic

### Ticket System API

- **Default limit:** Usually 60-100 requests/minute
- **Strategy:** Exponential backoff on 429 errors
- **Caching:** Cache ticket data for 5-10 minutes to reduce API calls

### Source Control API

- **Default limit:** 5000 requests/hour (GitHub), varies by provider
- **Strategy:** Batch operations when possible
- **Caching:** Cache file contents for session duration

### CI/CD API

- **Default limit:** Varies by provider
- **Strategy:** Poll with exponential backoff
- **Avoid:** Excessive status checks on PRs

---

## Error Handling

### Ticket Read Failures

- Log error with ticket ID
- Post comment to ticket: "Unable to analyze. Please check ticket format."
- Alert team if repeated failures

### Spec Creation Failures

- Post escalation to ticket with specific error
- Tag appropriate team member for manual review
- Do not create PR if confidence is too low

### PR Creation Failures

- Log error with full context
- Retry once after 30 seconds
- If still failing, post to ticket asking for manual intervention

---

## Monitoring & Observability

### Key Metrics

- **Ticket-to-PR time:** Target < 5 minutes for confident stories
- **AI confidence distribution:** Track High/Medium/Low ratios over time
- **Escalation rate:** Target < 10% (indicates context file completeness)
- **PR merge time:** Target < 15 minutes for BO approval
- **API error rates:** Alert on > 1% failure rate

### Logging

- Log all AI agent executions (see [CLAUDE.md](../CLAUDE.md) for requirements)
- Include ticket ID, confidence level, context versions used
- Store conversation transcripts for conversational agents
- Enable traceability from requirement to implementation

### Alerts

- Alert on high escalation rate (> 20% over 1 week)
- Alert on API failures (> 5% over 1 hour)
- Alert on stuck PRs (no review after 24 hours)
- Alert on CI failures for spec validation

---

## Next Steps

- **For quick setup:** See [Implementation Summary](implementation-summary.md)
- **For MCP gateway details:** See [FOSS Components Evaluation](foss-components-evaluation.md)
- **For workflow details:** See [Workflow Guide](workflow.md)
- **For CI configuration:** See [CI Configuration Guide](ci-configuration.md)
