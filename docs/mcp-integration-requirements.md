# MCP Integration Requirements

**Status:** Needs Development/Discovery

This document outlines the MCP (Model Context Protocol) integrations required for the AI-augmented requirements workflow. Both requirements-drafting-assistant and requirements-analyst agents need these integrations to function.

---

## Overview

**Problem:** Requirements agents need to access ticketing systems and source code repositories to synthesize requirements and draft specifications.

**Solution:** Use MCP (Model Context Protocol) to provide consistent, authenticated access to external systems.

**Key Requirement:** Credentials configured once, shared across both agents.

---

## 1. Ticketing System Access (MCP Server)

### Requirements

**Functional:**
- Read ticket data (description, acceptance criteria, labels, status, assignees)
- Search tickets by keywords, labels, components
- Read comments and discussion threads chronologically
- Fetch related tickets (linked issues, blocking/blocked by, parent/child relationships)
- Access ticket history (status changes, field updates, timeline)
- Query for recently updated tickets in specific projects

**Non-Functional:**
- Authentication: API token or OAuth
- Rate limiting: Respect ticketing system API limits
- Caching: Cache ticket data for session duration to reduce API calls
- Error handling: Graceful failures when tickets not found or access denied

### Supported Ticketing Systems

**Priority 1 (Must Support):**
- Jira (Cloud and Server)
- GitHub Issues
- Linear

**Priority 2 (Should Support):**
- Azure DevOps
- GitLab Issues
- Shortcut (formerly Clubhouse)

### Discovery Questions

**Authentication:**
- [ ] What authentication methods does each ticketing system support? (API tokens, OAuth, PAT)
- [ ] How do we securely store and rotate credentials?
- [ ] Do we need different permissions for different agents or can they share read-only access?
- [ ] How do we handle multiple ticketing systems in the same organization?

**API Capabilities:**
- [ ] What search syntax does each system support? (JQL for Jira, query language for Linear, etc.)
- [ ] What fields are available via API? (custom fields, attachments, watchers, etc.)
- [ ] How are comments/threads structured? (flat list, threaded, chronological)
- [ ] How are relationships represented? (links, references, parent/child)

**Rate Limiting:**
- [ ] What are the rate limits for each system?
- [ ] How should we handle rate limit errors? (retry with backoff, queue requests, fail fast)
- [ ] Should we implement request batching?

**Data Model:**
- [ ] How do we normalize ticket data across different systems?
- [ ] What common schema should MCP expose?
- [ ] How do we handle system-specific fields?

### MCP Server Interface

**Operations to implement:**

```typescript
interface TicketingMCPServer {
  // Fetch single ticket
  getTicket(ticketId: string): Promise<Ticket>

  // Search tickets
  searchTickets(query: string, filters?: TicketFilters): Promise<Ticket[]>

  // Get comments/discussion
  getTicketComments(ticketId: string): Promise<Comment[]>

  // Get related tickets
  getRelatedTickets(ticketId: string, relationshipType?: string): Promise<Ticket[]>

  // Get ticket history/timeline
  getTicketHistory(ticketId: string): Promise<HistoryEvent[]>
}

interface Ticket {
  id: string
  url: string
  title: string
  description: string
  status: string
  labels: string[]
  assignee?: string
  reporter: string
  createdAt: Date
  updatedAt: Date
  customFields: Record<string, any>
}

interface Comment {
  id: string
  author: string
  body: string
  createdAt: Date
  updatedAt: Date
}

interface TicketFilters {
  project?: string
  status?: string[]
  labels?: string[]
  assignee?: string
  createdAfter?: Date
  updatedAfter?: Date
}
```

### Development Tasks

- [ ] Implement Jira MCP server (Cloud API)
- [ ] Implement GitHub Issues MCP server
- [ ] Implement Linear MCP server
- [ ] Create unified ticket data model
- [ ] Implement authentication/credential management
- [ ] Implement rate limiting and retry logic
- [ ] Add request caching layer
- [ ] Write integration tests for each ticketing system
- [ ] Document configuration and deployment

---

## 2. Source Control Access for Non-Developers (MCP Server)

### Requirements

**Functional:**
- Read `.feature` files (Gherkin specifications) from repository
- Search `.feature` files by keywords, tags, file path patterns
- Read step definition files (`*_steps.go`, `*_steps.py`, etc.)
- List files in specific directories (e.g., `features/`, `features/step_definitions/`)
- Access file history (who changed what, when)
- Read files from specific branches (main, spec/*, impl/*)

**Access Control:**
- **Key Requirement:** Business Owners and non-developers need read-only access
- No write permissions (cannot commit, push, create branches)
- No access to sensitive files (credentials, secrets, private keys)
- Scoped to specific paths (features/, docs/, context files)

**Non-Functional:**
- Authentication: Git provider API (GitHub, GitLab, Bitbucket, Azure DevOps)
- Permissions: Read-only, scoped to repository
- Performance: Cache file contents for session duration
- Error handling: Graceful failures for missing files or access denied

### Supported Git Providers

**Priority 1 (Must Support):**
- GitHub (API v3/v4)
- GitLab (API v4)

**Priority 2 (Should Support):**
- Azure DevOps Repos
- Bitbucket Cloud
- Bitbucket Server

### Discovery Questions

**Authentication:**
- [ ] What authentication methods do we support? (Personal Access Tokens, OAuth, GitHub Apps)
- [ ] How do we handle different permission scopes?
- [ ] Should Business Owners use their own credentials or a shared read-only bot account?
- [ ] How do we revoke access when someone leaves the team?

**File Access:**
- [ ] Do we need to read files from multiple branches simultaneously?
- [ ] Should we support reading historical versions (git blame, file at specific commit)?
- [ ] How do we handle large repositories? (shallow clones, sparse checkout)
- [ ] Should we support searching across branches?

**Path Restrictions:**
- [ ] What paths should be accessible? (`features/`, `docs/`, `sample-project/context/`)
- [ ] What paths should be blocked? (`.env`, `secrets/`, `credentials/`, `*.key`, `*.pem`)
- [ ] How do we enforce path restrictions?
- [ ] Should restrictions be configurable per repository?

**Performance:**
- [ ] Should we clone the repository locally or use API calls?
- [ ] How long should we cache file contents?
- [ ] Should we implement incremental updates (git pull) or fetch fresh each time?

### MCP Server Interface

**Operations to implement:**

```typescript
interface SourceControlMCPServer {
  // Read single file
  readFile(path: string, branch?: string): Promise<FileContents>

  // Search files by pattern
  searchFiles(pattern: string, branch?: string): Promise<FilePath[]>

  // List directory contents
  listDirectory(path: string, branch?: string): Promise<DirectoryEntry[]>

  // Search file contents
  searchContents(query: string, pathPattern?: string, branch?: string): Promise<SearchResult[]>

  // Get file history
  getFileHistory(path: string, limit?: number): Promise<FileHistoryEntry[]>
}

interface FileContents {
  path: string
  branch: string
  content: string
  lastModified: Date
  lastAuthor: string
}

interface FilePath {
  path: string
  type: 'file' | 'directory'
}

interface DirectoryEntry {
  name: string
  path: string
  type: 'file' | 'directory'
  size: number
  lastModified: Date
}

interface SearchResult {
  path: string
  matches: CodeMatch[]
}

interface CodeMatch {
  lineNumber: number
  lineContent: string
  matchContext: string // surrounding lines
}

interface FileHistoryEntry {
  commit: string
  author: string
  date: Date
  message: string
  changes: 'added' | 'modified' | 'deleted'
}
```

### Security Considerations

**Path Restrictions:**
```typescript
// Allowed paths (read-only)
const ALLOWED_PATHS = [
  'features/**/*.feature',
  'features/step_definitions/**/*',
  'docs/**/*.md',
  'sample-project/context/**/*',
  'README.md',
  'CLAUDE.md'
]

// Blocked paths (never readable via MCP)
const BLOCKED_PATHS = [
  '.env',
  '.env.*',
  'secrets/**/*',
  'credentials/**/*',
  '**/*.key',
  '**/*.pem',
  '**/*.p12',
  '**/id_rsa',
  '**/id_ed25519',
  '.git/config', // contains remote URLs which may have embedded credentials
  '**/*secret*',
  '**/*password*'
]
```

**Access Control:**
- [ ] Enforce path restrictions (allowed/blocked lists)
- [ ] Return clear error messages when access denied
- [ ] Validate permissions before file operations

### Development Tasks

- [ ] Implement GitHub source control MCP server
- [ ] Implement GitLab source control MCP server
- [ ] Create path restriction enforcement layer
- [ ] Implement file search and content search
- [ ] Add caching layer for file contents
- [ ] Implement authentication/credential management
- [ ] Write integration tests for each git provider
- [ ] Document configuration and deployment
- [ ] Create security review checklist for path restrictions

---

## 3. Shared MCP Configuration

### Requirements

**Single Configuration for Both Agents:**
- Credentials stored once, referenced by both requirements-drafting-assistant and requirements-analyst
- MCP server endpoints configured once
- Connection pooling shared across agent sessions
- Consistent error handling and logging

### Configuration Schema

```yaml
mcp:
  ticketing:
    provider: "jira" | "github-issues" | "linear"
    baseUrl: "https://company.atlassian.net"
    authentication:
      type: "api-token" | "oauth" | "pat"
      credentials:
        token: "${TICKETING_API_TOKEN}" # from environment variable
    options:
      timeout: 30s
      maxRetries: 3
      cacheExpiration: 5m

  sourceControl:
    provider: "github" | "gitlab" | "azure-devops"
    repository: "org/repo"
    authentication:
      type: "pat" | "oauth" | "github-app"
      credentials:
        token: "${GIT_READ_TOKEN}" # from environment variable
    options:
      timeout: 30s
      defaultBranch: "main"
      cacheExpiration: 10m
      allowedPaths:
        - "features/**/*"
        - "docs/**/*.md"
        - "sample-project/context/**/*"
      blockedPaths:
        - ".env*"
        - "secrets/**/*"
        - "**/*.key"
        - "**/*.pem"
```

### Discovery Questions

- [ ] Where should MCP configuration live? (environment variables, config file, secrets manager)
- [ ] How do we handle multiple repositories or ticketing projects?
- [ ] Should we support hot-reloading configuration changes?
- [ ] How do we test MCP connections before agents use them?

### Development Tasks

- [ ] Design configuration schema
- [ ] Implement configuration loader with validation
- [ ] Implement credential management (env vars, secrets manager integration)
- [ ] Create MCP server registry (map provider names to server implementations)
- [ ] Add health check endpoints for each MCP server
- [ ] Document configuration options and examples
- [ ] Create troubleshooting guide for connection issues

---

## 4. Testing Strategy

### Unit Tests
- [ ] Mock MCP servers for agent testing
- [ ] Test credential handling and authentication
- [ ] Test path restriction enforcement
- [ ] Test error handling (network errors, auth failures, rate limits)

### Integration Tests
- [ ] Test against real ticketing systems (sandbox environments)
- [ ] Test against real git providers (test repositories)
- [ ] Test credential rotation
- [ ] Test rate limiting and retry logic

### End-to-End Tests
- [ ] Test full agent workflow with MCP integration
- [ ] Test requirements-drafting-assistant pulling ticket data and feature files
- [ ] Test requirements-analyst reading tickets and drafting specs
- [ ] Test credential sharing across both agents

---

## 5. Documentation Needs

### For Administrators
- [ ] Installation and configuration guide
- [ ] Credential setup instructions for each provider
- [ ] Troubleshooting common issues
- [ ] Security best practices

### For Agent Developers
- [ ] MCP server API reference
- [ ] How to add a new ticketing system provider
- [ ] How to add a new git provider
- [ ] Error handling patterns

### For End Users (Business Owners)
- [ ] What data agents can access and why
- [ ] How to grant/revoke access
- [ ] Privacy and security considerations
- [ ] What to do if access is denied

---

## 6. Roles and Access Requirements

### Business Owner Role

**Needs access to:**
- Ticketing system: Read tickets, comments, related tickets (no write access)
- Source control: Read `.feature` files, step definitions, context files (no write access to code)
- Source control: Write access to `.feature` files only (via PR approval, not direct MCP write)

**Authentication:**
- Personal credentials for ticketing system (for audit trail)
- Shared read-only bot account for source control OR personal read-only access

**Use cases:**
- Review AI-drafted specifications in feature files
- Access ticket history and comments to understand context
- See related tickets to understand dependencies
- Read existing feature files to maintain consistency

### Requirements Agent (AI) Role

**Needs access to:**
- Ticketing system: Read tickets, comments, search, related tickets (no write access)
- Source control: Read `.feature` files, step definitions, context files (no write access via MCP)
- Source control: Write access to create spec branches and PRs (via git API, not MCP)

**Authentication:**
- Service account credentials for ticketing system
- Service account credentials for source control (with PR creation permissions)

**Use cases:**
- Fetch ticket data to draft specifications
- Search for related tickets to understand context
- Read existing feature files to reuse steps and maintain conventions
- Create spec branches and PRs (via git API, separate from MCP read access)

### Developer Role

**Needs access to:**
- Full repository access (read and write code)
- Ticketing system: Read and update tickets
- CI/CD: Trigger builds, view test results

**Authentication:**
- Personal credentials for both systems

**Use cases:**
- Implement step definitions and service code
- Run and debug BDD scenarios locally
- Review and merge PRs

### Differentiation

| Capability | Business Owner | AI Agent | Developer |
|------------|---------------|----------|-----------|
| Read tickets | ✅ Personal | ✅ Service | ✅ Personal |
| Write/update tickets | ❌ | ❌ | ✅ |
| Read `.feature` files | ✅ MCP (read-only) | ✅ MCP (read-only) | ✅ Full git |
| Write `.feature` files | ✅ Via PR approval | ✅ Via PR creation | ✅ Via PR/commit |
| Read code files | ❌ | ❌ | ✅ |
| Write code files | ❌ | ❌ | ✅ |
| Read context files | ✅ MCP | ✅ MCP | ✅ Full git |
| Write context files | ✅ Via PR (CODEOWNERS) | ❌ | ✅ Via PR |

**Key Insight:** Business Owners need read access to specifications and tickets, but should not have full repository access (security concern - no access to code, credentials, secrets).

---

## Priority and Phasing

### Phase 1: MVP (Required for workflow to function)
1. Ticketing system MCP server for primary provider (Jira or GitHub Issues)
2. Source control MCP server for primary git provider (GitHub)
3. Basic authentication (API tokens)
4. Path restrictions for source control
5. Shared configuration for both agents

### Phase 2: Production-Ready
1. Additional ticketing system providers
2. Additional git providers
3. Caching layer
4. Rate limiting and retry logic
5. Health checks and monitoring

### Phase 3: Enhanced
1. OAuth authentication
2. Multiple repositories/projects
3. Advanced search capabilities
4. Performance optimizations
5. Dashboard for monitoring MCP usage

---

## Open Questions

**Architecture:**
- [ ] Should MCP servers run as separate processes or embedded in agent runtime?
- [ ] How do we scale MCP servers for multiple concurrent agent sessions?
- [ ] Should we use connection pooling or create connections per request?

**Deployment:**
- [ ] Where do MCP servers run? (same host as agents, separate service, cloud-hosted)
- [ ] How do we handle updates and versioning of MCP servers?
- [ ] What monitoring and observability do we need?

**Security:**
- [ ] How do we ensure credentials never leak into logs or agent outputs?
- [ ] Should we encrypt cached data?
- [ ] How often should we rotate service account credentials?

---

## Success Criteria

**MCP integration is successful when:**
- ✅ Both requirements-drafting-assistant and requirements-analyst can read tickets and feature files
- ✅ Business Owners can access specifications without full repository access
- ✅ Credentials are configured once and shared across both agents
- ✅ Path restrictions prevent access to sensitive files
- ✅ Agents can search and retrieve relevant historical context
- ✅ System handles authentication failures gracefully
- ✅ Performance is acceptable (< 2s for typical ticket or file read)
