# Implementation Summary

**Quick start guide for implementing the AI-augmented requirements workflow**

---

## What You Need

This workflow requires:
1. **MCP Gateway** to connect AI agents to ticketing systems and source control
2. **Official MCP Servers** for GitHub, GitLab, and Jira (already exist)
3. **Context Files** defining your team's domain, architecture, testing patterns, and standards

**Good news:** NO CUSTOM DEVELOPMENT NEEDED! Use existing open source solutions.

---

## Recommended Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  AI Requirements Agents                                     │
│  ┌──────────────────────┐  ┌──────────────────────┐       │
│  │ requirements-        │  │ requirements-        │       │
│  │ drafting-assistant   │  │ analyst              │       │
│  └──────────┬───────────┘  └──────────┬───────────┘       │
│             └───────────┬───────────────┘                   │
└─────────────────────────┼───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│  IBM ContextForge MCP Gateway                               │
│  - Path filtering (block .env, secrets, keys)              │
│  - Authentication (OAuth, SSO)                             │
│  - Authorization (RBAC for Business Owners, Developers)    │
│  - Rate limiting (60 req/min default)                      │
│  - Caching (Redis, 10min TTL)                              │
│  - Monitoring (OpenTelemetry, Admin UI)                    │
└─────────┬───────────────────────────────────────────────────┘
          │
          ├──────────────┬──────────────┬──────────────┐
          ▼              ▼              ▼              ▼
    GitHub MCP     GitLab MCP     Jira MCP      Linear MCP
    (official)     (official)     (official)    (future)
```

---

## Step 1: Choose Your MCP Gateway

### Option A: IBM ContextForge (Recommended)

**Best for:** Most teams - provides all features out-of-the-box

**Features:**
- ✅ Path filtering (allowed/blocked lists)
- ✅ Authentication (OAuth, SSO, multi-tenant)
- ✅ Authorization (RBAC with teams and roles)
- ✅ Rate limiting (built-in, user-scoped)
- ✅ Caching (Redis-backed, configurable TTL)
- ✅ Retry logic (automatic)
- ✅ Monitoring (OpenTelemetry, admin UI, real-time logs)
- ✅ Protocol flexibility (wraps REST APIs as MCP)

**Setup time:** 1-2 days
**Repository:** https://github.com/IBM/mcp-context-forge

### Option B: Lasso MCP Gateway

**Best for:** High-security environments

**Features:**
- ✅ Security scanner (reputation analysis)
- ✅ Guardrails (prompt injection, PII detection)
- ✅ Plugin architecture (Basic, Presidio, Lasso)
- ✅ Experiment tracking (Xetrack)

**Setup time:** 1 day
**Repository:** https://github.com/lasso-security/mcp-gateway

### Option C: Docker MCP Gateway

**Best for:** Teams already using Docker infrastructure

**Features:**
- ✅ Container isolation (each MCP server isolated)
- ✅ Built into Docker Desktop
- ✅ OAuth flows, credential management
- ⚠️ Add Traefik + Redis for rate limiting and caching

**Setup time:** 2-3 days
**Repository:** https://github.com/docker/mcp-gateway

**Comparison:** See [FOSS Components Evaluation](foss-components-evaluation.md) for detailed comparison.

---

## Step 2: Install and Configure (IBM ContextForge Example)

### Installation

```bash
# Install ContextForge
pip install mcp-contextforge-gateway

# Start Redis (for caching)
docker run -d -p 6379:6379 redis:alpine

# Optional: Start Jaeger (for tracing)
docker run -d -p 16686:16686 -p 4318:4318 jaegertracing/all-in-one
```

### Configuration

Create `contextforge-config.yaml`:

```yaml
gateway:
  port: 8080

authentication:
  method: oauth
  sso:
    enabled: true
    provider: azure-ad  # or google, okta, etc.

tenants:
  - name: requirements-team
    email_domain: company.com
    roles:
      - name: business-owner
        permissions:
          - read:tickets
          - read:features
          - approve:specs
      - name: requirements-analyst
        permissions:
          - read:tickets
          - read:features
          - write:specs

servers:
  # GitHub - for reading .feature files
  - name: github
    type: mcp
    url: https://github-mcp.api.github.com
    auth:
      type: oauth
      clientId: ${GITHUB_CLIENT_ID}
      clientSecret: ${GITHUB_CLIENT_SECRET}
    filters:
      allowed_paths:
        - features/**/*.feature
        - features/step_definitions/**/*
        - docs/**/*.md
        - sample-project/context/**/*
        - README.md
        - CLAUDE.md
      blocked_paths:
        - .env*
        - secrets/**/*
        - "**/*.key"
        - "**/*.pem"
        - "**/*.p12"
        - "**/id_rsa*"

  # Jira - for reading tickets
  - name: jira
    type: mcp
    url: https://company.atlassian.net
    auth:
      type: oauth
      clientId: ${JIRA_CLIENT_ID}
      clientSecret: ${JIRA_CLIENT_SECRET}

caching:
  enabled: true
  backend: redis
  redis:
    host: localhost
    port: 6379
  ttl: 600  # 10 minutes

rate_limiting:
  enabled: true
  default:
    requests_per_minute: 60
  per_user:
    business-owner: 100
    requirements-analyst: 200

observability:
  opentelemetry:
    enabled: true
    endpoint: http://localhost:4318
  metrics:
    llm_tracking: true
    cost_tracking: true
```

### Environment Variables

```bash
# GitHub OAuth
export GITHUB_CLIENT_ID=your_github_client_id
export GITHUB_CLIENT_SECRET=your_github_client_secret

# Jira OAuth
export JIRA_CLIENT_ID=your_jira_client_id
export JIRA_CLIENT_SECRET=your_jira_client_secret
```

### Start Gateway

```bash
contextforge serve --config contextforge-config.yaml
```

**Gateway running at:** http://localhost:8080
**Admin UI:** http://localhost:8080/admin
**Jaeger UI:** http://localhost:16686

---

## Step 3: Set Up Context Files

Create context files defining your team's knowledge. **Start with these four:**

**Note:** Context files can grow and shrink based on your needs. Split them when they get too large (e.g., split business.md into business-rules.md and compliance.md), merge them if too fragmented, or add specialized files as your domain requires. These are a reasonable starting point, not a rigid requirement.

**Important:** When you change context files, update agent prompts to reference the new structure. For example, if you split business.md into business-rules.md and compliance.md, update prompts to read both files instead of just business.md.

### business.md
```markdown
# Business Context

## Domain Glossary
- **Story:** User-facing requirement with acceptance criteria
- **Spec:** Gherkin scenarios defining expected behavior

## User Personas
- **Business Owner:** Approves specifications
- **Developer:** Implements features

## Business Rules
- BR-001: All user input must be validated
- BR-002: Passwords must be 12+ characters

## Compliance
- GDPR: Personal data handling requirements
```

### architecture.md
```markdown
# Architecture Context

## System Overview
- Language: Go 1.21
- Framework: Godog (BDD testing)
- Database: PostgreSQL

## API Contracts
- Internal: /api/specs/internal-api.yaml
- External: /api/specs/external-api.yaml

## External Dependencies
- Email: SendGrid API
- Authentication: Auth0
```

### testing.md
```markdown
# Testing Context

## Step Library
- `Given a user exists with email {string}` - Create test user
- `When I request a password reset for {string}` - Trigger reset
- `Then I should receive a reset email` - Verify email sent

## Boundary Patterns
- **String:** empty, whitespace, max length (255), unicode, special chars
- **Email:** valid, invalid, max length, special domains
- **Numeric:** zero, negative, min, max, overflow

## Edge Case Patterns
- Authentication: expired tokens, revoked tokens, invalid signatures
- Rate limiting: burst patterns, sustained load
```

### tech_standards.md
```markdown
# Technical Standards

## Language & Framework
- Go 1.21
- Godog for BDD testing
- PostgreSQL for persistence

## IoC Patterns
- Primary constructor (for testing with mocks)
- Production factory (for real dependencies)
- Coverage: Ignore production factories

## Directory Structure
```
internal/
  domain/
    services/    # Business logic
  repositories/  # Data access
features/
  *.feature     # Gherkin scenarios
  step_definitions/  # Test implementations
```

## Naming Conventions
- Interfaces: `UserRepository`
- Implementations: `PostgresUserRepository`
- Tests: `*_test.go`
```

**See:** [Context Files Guide](context-files.md) for detailed examples.

---

## Step 4: Configure AI Agents

Both agents use the same MCP gateway (shared credentials):

### requirements-drafting-assistant

**Purpose:** Help Business Owners articulate requirements through conversation

**Configuration:**
- MCP endpoint: http://localhost:8080
- Access: Jira (read tickets, comments, related tickets)
- Access: GitHub (read .feature files for conventions)

**Usage:**
```bash
# BO has vague requirement, needs help articulating it
claude --agent requirements-drafting-assistant
> "We need password reset functionality"

# Agent asks clarifying questions, pulls historical context
# Outputs structured requirement ready for ticketing
```

### requirements-analyst

**Purpose:** Analyze tickets and draft Gherkin specifications automatically

**Configuration:**
- MCP endpoint: http://localhost:8080
- Access: Jira (read ticket, comments, related tickets)
- Access: GitHub (read .feature files, create branches/PRs)
- Context: business.md, architecture.md, testing.md, tech_standards.md

**Trigger:** Automatic (ticket labeled "ready-for-spec" via webhook) or manual

**Workflow:**
1. Reads ticket via MCP
2. Searches related tickets for context
3. Reads existing .feature files for conventions
4. Drafts Gherkin scenarios
5. Creates `spec/PROJ-1234` branch
6. **Opens Gherkin/Cucumber PR** for BO review (triggers bo-review)

**Output:** Gherkin/Cucumber PR with .feature files

### bo-review

**Purpose:** Assist Business Owners in reviewing AI-drafted Gherkin specifications

**Configuration:**
- MCP endpoint: http://localhost:8080
- Access: GitHub (read .feature files in PR)
- Access: Jira (read original ticket for context)
- Context: business.md (for business rules validation)

**Trigger:** Runs on Gherkin/Cucumber PRs (spec/ branches containing .feature files)

**Workflow:**
1. Reads AI-drafted .feature files from PR
2. Compares against original ticket requirements
3. Validates business rules are correctly applied
4. Identifies missing critical scenarios
5. Flags security/compliance concerns
6. Produces review report with APPROVED or CHANGES REQUESTED

**Output:** Review report posted as PR comment, guiding BO in approval decision

**See:** [AI Agents Guide](agents.md) for detailed documentation.

---

## Step 5: Test End-to-End

### Test requirements-drafting-assistant

```bash
# Verify agent can:
# 1. Read Jira tickets
# 2. Search related tickets
# 3. Read existing .feature files
# 4. Generate structured requirement

# Expected: Agent asks clarifying questions, references past tickets
```

### Test requirements-analyst

```bash
# Create test ticket in Jira
# Label it "ready-for-spec"
# Webhook triggers agent

# Expected:
# - Agent reads ticket via MCP
# - Agent reads existing .feature files for conventions
# - Agent creates spec/PROJ-1234 branch
# - Agent opens PR with Gherkin scenarios
# - PR assigned to Business Owner for review
```

### Verify Security

```bash
# Try to read blocked file
curl http://localhost:8080/github/read?path=.env
# Expected: 403 Forbidden

# Try to read allowed file
curl http://localhost:8080/github/read?path=features/auth/login.feature
# Expected: 200 OK with file contents
```

### Verify Caching

```bash
# First request (cache miss)
time curl http://localhost:8080/jira/ticket/PROJ-1234
# Expected: ~2 seconds

# Second request (cache hit)
time curl http://localhost:8080/jira/ticket/PROJ-1234
# Expected: ~100ms (20x faster)
```

### Verify Rate Limiting

```bash
# Send 65 requests in 1 minute (limit: 60)
for i in {1..65}; do curl http://localhost:8080/jira/ticket/PROJ-1234; done
# Expected: First 60 succeed, last 5 return 429 Too Many Requests
```

---

## Step 6: Integrate with CI/CD

### GitHub Actions Example

```yaml
name: Gherkin/Cucumber PR Automation

on:
  issues:
    types: [labeled]
  pull_request:
    paths:
      - 'features/**/*.feature'

jobs:
  # requirements-analyst: Runs when ticket labeled "ready-for-spec"
  draft-spec:
    if: github.event.label.name == 'ready-for-spec'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run requirements-analyst
        run: |
          claude --agent requirements-analyst \
                 --input "Ticket: ${{ github.event.issue.html_url }}"

      - name: Create Gherkin/Cucumber PR
        run: |
          # Agent creates spec/PROJ-1234 branch with .feature files
          # Opens PR for BO review (triggers bo-review below)

  # bo-review: Runs on Gherkin/Cucumber PRs (spec/ branches)
  bo-review-assistant:
    if: startsWith(github.head_ref, 'spec/')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run bo-review agent
        run: |
          claude --agent bo-review \
                 --input "PR #${{ github.event.pull_request.number }}" \
                 --context business.md

      - name: Post review report as PR comment
        run: |
          gh pr comment ${{ github.event.pull_request.number }} \
             --body-file review-report.md

  block-pending-impl:
    runs-on: ubuntu-latest
    if: startsWith(github.head_ref, 'impl/')
    steps:
      - uses: actions/checkout@v4

      - name: Check for @pending tags
        run: |
          STORY_ID=$(echo "${{ github.head_ref }}" | sed 's/^impl\///')
          if grep -r "@pending.*@story-$STORY_ID\|@story-$STORY_ID.*@pending" features/; then
            echo "ERROR: Implementation PR contains @pending for $STORY_ID"
            exit 1
          fi
```

**See:** [CI Configuration Guide](ci-configuration.md) for detailed setup.

---

## Timeline

| Phase | Duration | Effort |
|-------|----------|--------|
| **Phase 1: Gateway Setup** | Day 1 | Install ContextForge, configure, start services |
| **Phase 2: Context Files** | Day 2 | Create business.md, architecture.md, testing.md, tech_standards.md |
| **Phase 3: Agent Config** | Day 3 | Configure both agents to use MCP gateway |
| **Phase 4: Testing** | Day 4 | End-to-end testing, verify all features work |
| **Phase 5: CI Integration** | Day 5 | Add GitHub Actions, webhooks, automation |
| **Total** | **1 week** | **vs. 3-4 weeks custom development** |

---

## Cost Comparison

### Custom Development (Original Plan)
- **Path restriction service:** 1-2 days
- **Authentication/authorization:** 2-3 days
- **Rate limiting:** 1-2 days
- **Caching layer:** 1-2 days
- **Monitoring/health checks:** 2-3 days
- **Integration testing:** 2-3 days
- **Documentation:** 1-2 days

**Total:** 3-4 weeks + ongoing maintenance

### FOSS Solution (Recommended)
- **Install ContextForge:** 1 hour
- **Configuration:** 4-6 hours
- **Testing:** 4-6 hours
- **Documentation:** 2-3 hours

**Total:** 1-2 days + minimal ongoing maintenance

**Savings:** 2.5-3.5 weeks of development time

---

## Success Metrics

After implementation, you should achieve:

**Speed:**
- ✅ Spec review: < 15 minutes (async, no meetings)
- ✅ Ticket read: < 2 seconds
- ✅ File read: < 2 seconds (< 100ms cached)

**Security:**
- ✅ Business Owners have read-only access (no code access)
- ✅ Sensitive files blocked (.env, secrets, keys)
- ✅ Path restrictions enforced (403 on blocked paths)

**Reliability:**
- ✅ Rate limits respected (no 429 errors under normal load)
- ✅ Caching working (80%+ cache hit rate)
- ✅ Shared credentials between agents (no duplication)

**Workflow:**
- ✅ Story → Spec draft: < 30 minutes (automatic)
- ✅ BO approval: < 15 minutes (async PR review)
- ✅ AI escalation rate: < 10% (context files are complete)

---

## Troubleshooting

### Agent can't read tickets

**Check:**
- Is ContextForge running? `curl http://localhost:8080/health`
- Are Jira credentials correct? Check `JIRA_CLIENT_ID`, `JIRA_CLIENT_SECRET`
- Is MCP endpoint configured correctly in agent?

### Path filtering not working

**Check:**
- Are paths in `allowed_paths` and `blocked_paths` correct?
- Use glob patterns: `features/**/*.feature` not `features/*.feature`
- Test directly: `curl http://localhost:8080/github/read?path=.env`

### Rate limiting too aggressive

**Adjust in config:**
```yaml
rate_limiting:
  default:
    requests_per_minute: 120  # Increase from 60
```

### Caching not working

**Check:**
- Is Redis running? `redis-cli ping` (should return `PONG`)
- Is caching enabled in config? `caching.enabled: true`
- Check cache hit rate in admin UI: http://localhost:8080/admin

**See:** [Troubleshooting Guide](troubleshooting.md) for more solutions.

---

## Next Steps

1. **Read detailed guides:**
   - [Context Files](context-files.md) - How to structure your team's knowledge
   - [AI Agents](agents.md) - How each agent works
   - [FOSS Evaluation](foss-components-evaluation.md) - Full comparison of gateway options
   - [MCP Integration Requirements](mcp-integration-requirements.md) - Technical details

2. **Set up your environment:**
   - Install IBM ContextForge (or chosen gateway)
   - Create context files
   - Configure OAuth apps for GitHub/Jira

3. **Test with sample project:**
   - Use example context files from `sample-project/`
   - Create test ticket
   - Run requirements-analyst
   - Verify spec generation works

4. **Deploy to production:**
   - Move credentials to secrets manager
   - Add monitoring and alerts
   - Train team on workflow
   - Document your specific setup

---

## Resources

**Documentation:**
- [Workflow Overview](workflow.md) - Complete end-to-end process
- [BDD Value Proposition](bdd-value.md) - Why automated acceptance tests matter
- [Roles & Responsibilities](roles.md) - Who does what

**Technical:**
- [MCP Server Evaluation](mcp-server-evaluation.md) - Evaluation of GitHub/GitLab/Jira MCP servers
- [Source Control Strategy](source-control.md) - Branch model, merge rules
- [Integration Architecture](integration.md) - Webhooks, APIs, CI/CD

**Examples:**
- [Sample Project](../sample-project/) - Example context files
- [Agent Prompts](prompts/) - Complete agent prompts with examples

---

## Support

**Questions about this workflow:**
- Review documentation in `docs/`
- Check [Troubleshooting Guide](troubleshooting.md)
- Review example conversations in `docs/prompts/*/example-output/`

**Questions about IBM ContextForge:**
- GitHub Issues: https://github.com/IBM/mcp-context-forge/issues
- Documentation: https://ibm.github.io/mcp-context-forge/

**Questions about MCP:**
- MCP Documentation: https://modelcontextprotocol.io/
- MCP Servers: https://github.com/modelcontextprotocol/servers

---

**Remember:** Humans make decisions, AI implements them, humans verify. This system accelerates requirements without removing human judgment.
