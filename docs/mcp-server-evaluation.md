# MCP Server Evaluation: GitHub, GitLab, and Jira

**Date:** 2025-12-04
**Purpose:** Evaluate existing MCP servers against requirements for AI-augmented requirements workflow

This document evaluates existing MCP (Model Context Protocol) servers for GitHub, GitLab, and Jira/Atlassian against the requirements outlined in [mcp-integration-requirements.md](mcp-integration-requirements.md).

---

## Executive Summary

✅ **GitHub MCP Server:** Meets all source control requirements
✅ **GitLab MCP Server:** Meets all source control requirements
✅ **Atlassian/Jira MCP Servers:** Meet all ticketing system requirements

**Recommendation:** Use existing official and community MCP servers. No custom development required for MVP.

---

## 1. GitHub MCP Server

### Official Server
**Repository:** https://github.com/github/github-mcp-server
**Status:** Generally Available (GA) as of September 2025
**Type:** Remote MCP server (no local installation required)
**Authentication:** OAuth 2.1 + PKCE

### Capabilities vs. Requirements

| Requirement | GitHub MCP Server | Status |
|-------------|-------------------|--------|
| **Read .feature files** | ✅ Read file contents from repositories | ✅ Met |
| **Search files by pattern** | ✅ Search files, browse repository structure | ✅ Met |
| **Search file contents** | ✅ Search code, analyze commits | ✅ Met |
| **List directory contents** | ✅ Browse repository structure | ✅ Met |
| **Read from specific branches** | ✅ Access any repository/branch | ✅ Met |
| **File history** | ✅ Analyze commits, track changes | ✅ Met |
| **Path restrictions** | ⚠️ Application-level enforcement needed | ⚠️ Custom |
| **Read-only for non-developers** | ✅ OAuth scopes control access | ✅ Met |

### Available Operations

**Repository & Code:**
- Browse and query code
- Search files across repositories
- Read file contents
- Analyze commits and understand project structure

**Issues & Pull Requests:**
- Create, update, and manage issues
- Create, update, and manage pull requests
- Triage bugs, review code changes
- Maintain project boards

**GitHub Actions:**
- Monitor workflow runs
- Analyze build failures
- View CI/CD pipeline status
- Re-run failed jobs

**Security:**
- Surface code scanning alerts
- Review Dependabot alerts
- Examine security findings

**Collaboration:**
- Access discussions
- Manage notifications
- Analyze team activity

### Configuration

**Toolsets available:**
- `repos` - Repository operations
- `issues` - Issue management
- `pull_requests` - PR management
- `actions` - GitHub Actions/CI-CD
- `code_security` - Security insights

**Authentication:**
- OAuth 2.1 with PKCE
- Supports personal access tokens (PAT)
- Fine-grained permissions control
- Integrated with GitHub Copilot IDEs (VS Code, Visual Studio, JetBrains, Eclipse, Xcode) and third-party tools (Claude Desktop, Cursor)

### Gap Analysis

✅ **All core requirements met**

**Custom implementation needed:**
1. **Path restrictions:** GitHub MCP server doesn't enforce path-based access restrictions. Need application-layer filtering to block sensitive files (`.env`, secrets, credentials).
2. **Audit logging:** GitHub provides access logs via API, but we need to implement our own audit trail for MCP usage.

**Mitigation:**
- Implement wrapper service that filters file read requests based on allowed/blocked path patterns
- Use GitHub API to query access logs for audit trail
- Configure OAuth scopes to minimum necessary permissions (read-only for non-developers)

---

## 2. GitLab MCP Server

### Official Server
**Documentation:** https://docs.gitlab.com/user/gitlab_duo/model_context_protocol/mcp_server/
**Status:** Beta (as of GitLab 18.6, December 2024)
**Type:** Remote MCP server (HTTP transport) or stdio with mcp-remote
**Authentication:** OAuth 2.0 Dynamic Client Registration

### Capabilities vs. Requirements

| Requirement | GitLab MCP Server | Status |
|-------------|-------------------|--------|
| **Read .feature files** | ✅ Access project files | ✅ Met |
| **Search files by pattern** | ✅ Global and project-scoped search | ✅ Met |
| **Search file contents** | ✅ Semantic code search (vector embeddings on GitLab.com) | ✅ Met |
| **List directory contents** | ✅ Browse project structure | ✅ Met |
| **Read from specific branches** | ✅ Access any branch | ✅ Met |
| **File history** | ✅ Access commits and diffs | ✅ Met |
| **Path restrictions** | ⚠️ Application-level enforcement needed | ⚠️ Custom |
| **Read-only for non-developers** | ✅ OAuth permissions control | ✅ Met |

### Available Operations

**Issue Management:**
- Create new issues (title, description, assignees, milestones, labels)
- Retrieve detailed issue information
- Add comments/notes to issues
- Support internal notes (visible only to Reporters+)

**Merge Request Operations:**
- Create merge requests between branches
- Fetch merge request details, commits, and diffs
- Access pipeline information for merge requests
- Add comments to merge requests

**Pipeline & CI/CD:**
- Retrieve jobs from specific CI/CD pipelines
- View pipeline status and details
- Access build logs and results

**Search Capabilities:**
- Global search across issues, merge requests, and projects
- Group and project-scoped searching
- Filter by state and confidentiality
- **Semantic code search** using vector embeddings (GitLab.com only) - finds relevant code snippets

**Project Information:**
- Access GitLab project information
- Browse project structure and files

### Configuration

**Transport types:**
- **HTTP transport (recommended):** Direct connection without additional dependencies
- **stdio transport with mcp-remote:** Connection through proxy (requires Node.js)

**Authentication:**
- OAuth 2.0 Dynamic Client Registration
- Respects GitLab user permissions
- Fine-grained access control

### Gap Analysis

✅ **All core requirements met**

**Custom implementation needed:**
1. **Path restrictions:** GitLab MCP server doesn't enforce path-based access restrictions. Need application-layer filtering.
2. **Audit logging:** Need to implement audit trail for MCP usage.

**Mitigation:**
- Same as GitHub: implement wrapper service for path filtering
- Use GitLab audit events API for access logging
- Configure OAuth scopes for minimum necessary permissions

**Bonus capability:**
- Semantic code search on GitLab.com provides advanced search beyond simple text matching

---

## 3. Atlassian/Jira MCP Server

### Official Server
**Repository:** https://github.com/atlassian/atlassian-mcp-server
**Status:** Beta (open to all Atlassian Cloud customers)
**Type:** Remote MCP server (cloud-based bridge)
**Authentication:** OAuth 2.0

### Alternative Community Servers
- **aashari/mcp-server-atlassian-jira:** Generic HTTP tools for any Jira API endpoint
- **sooperset/mcp-atlassian:** Supports both Cloud and Server/Data Center
- **xuanxt/atlassian-mcp:** 51 tools for Jira/Confluence management

### Capabilities vs. Requirements

| Requirement | Atlassian Jira MCP | Status |
|-------------|-------------------|--------|
| **Read ticket data** | ✅ Get full issue details | ✅ Met |
| **Search tickets** | ✅ JQL queries, keyword search | ✅ Met |
| **Read comments/threads** | ✅ List and add comments | ✅ Met |
| **Get related tickets** | ✅ Access linked issues, blocks/blocked-by | ✅ Met |
| **Ticket history/timeline** | ✅ Access changelog, worklogs | ✅ Met |
| **Search by filters** | ✅ JQL with project, status, labels, assignee filters | ✅ Met |
| **Read-only access** | ✅ OAuth scopes control write permissions | ✅ Met |

### Available Operations (Official Atlassian Server)

**Search & Discovery:**
- Search across Jira and Confluence
- Find issues by project, status, assignee
- JQL query support

**Issue Management:**
- Create new issues
- Retrieve full issue details
- Update issue fields
- Bulk ticket generation from notes

**Comments & Collaboration:**
- View comments and discussion threads
- Add comments to issues
- Access worklogs
- Track time spent

**Related Information:**
- Access linked issues (blocks, blocked-by, relates to)
- View development info (commits, PRs, branches)
- Sprint and board information

**Security:**
- Respects Jira and Confluence user permissions
- OAuth 2.0 secure authentication
- HTTPS with TLS 1.2+ encryption

### Available Operations (Community Servers)

**aashari/mcp-server-atlassian-jira provides 5 generic HTTP tools:**
- `GET` - Read any Jira API endpoint (projects, issues, comments, worklogs)
- `POST` - Create resources (issues, comments)
- `PUT` - Update existing resources
- `PATCH` - Partial updates
- `DELETE` - Remove resources

**Specific capabilities:**
- Search issues using JQL: `assignee=currentUser() AND status='In Progress'`
- List/add comments
- Get issue details with full context
- Browse projects
- Optional JMESPath filtering for extracting specific fields

**sooperset/mcp-atlassian:**
- Supports both Jira Cloud and Server/Data Center
- Confluence and Jira integration
- Community-maintained

### Configuration

**Official Atlassian MCP Server:**
```
Required: Atlassian Cloud account
Authentication: OAuth 2.0 via web flow
No special sign-up required (beta open to all customers)
```

**Community Servers (e.g., aashari):**
```yaml
Environment variables:
- ATLASSIAN_EMAIL: Your Atlassian account email
- ATLASSIAN_TOKEN: API token from Atlassian Cloud
- ATLASSIAN_DOMAIN: Subdomain (e.g., "yourcompany" from yourcompany.atlassian.net)
```

### Gap Analysis

✅ **All core requirements met**

**No custom implementation needed for MVP**

**Recommendations:**
1. **Official Atlassian server** for Jira Cloud customers (simplest, OAuth-based)
2. **aashari/mcp-server-atlassian-jira** for flexibility (generic HTTP tools access any endpoint)
3. **sooperset/mcp-atlassian** for Server/Data Center deployments

---

## 4. Cross-Cutting Concerns

### Authentication & Credentials

| Platform | Official Server Auth | Community Server Auth | Shared Config Possible |
|----------|---------------------|----------------------|----------------------|
| GitHub | OAuth 2.1 + PKCE | Personal Access Token | ✅ Yes (OAuth app) |
| GitLab | OAuth 2.0 Dynamic | Personal Access Token | ✅ Yes (OAuth app) |
| Jira | OAuth 2.0 | API Token + Email | ✅ Yes (per platform) |

**Recommendation:**
- Use OAuth for official servers (better security, token refresh, scoped permissions)
- Store credentials in environment variables or secrets manager
- Configure once per platform, reference from both requirements agents

### Path Restrictions (GitHub/GitLab)

**Problem:** MCP servers don't enforce path-based access restrictions for source control.

**Solution:** Implement application-layer wrapper service:

```typescript
class SecureSourceControlMCP {
  private allowedPaths = [
    'features/**/*.feature',
    'features/step_definitions/**/*',
    'docs/**/*.md',
    'sample-project/context/**/*',
    'README.md',
    'CLAUDE.md'
  ];

  private blockedPaths = [
    '.env*',
    'secrets/**/*',
    '**/*.key',
    '**/*.pem',
    '**/*.p12',
    '**/id_rsa',
    '**/id_ed25519'
  ];

  async readFile(path: string): Promise<FileContents> {
    if (!this.isPathAllowed(path)) {
      throw new Error(`Access denied: ${path} is not in allowed paths`);
    }

    if (this.isPathBlocked(path)) {
      throw new Error(`Access denied: ${path} is blocked for security`);
    }

    return await this.mcpServer.readFile(path);
  }

  // Similar wrappers for searchFiles, listDirectory, etc.
}
```

**Deployment:**
- Wrapper service sits between agents and MCP servers
- Enforces allow/block lists
- Logs all access attempts (allowed and denied)
- Can be configured per repository

### Audit Logging

**Problem:** Need comprehensive audit trail for compliance and security.

**Solution:** Implement centralized logging service:

```typescript
interface MCPAccessLog {
  timestamp: Date;
  agent: 'requirements-drafting-assistant' | 'requirements-analyst';
  user: string; // Business Owner or service account
  platform: 'github' | 'gitlab' | 'jira';
  operation: string; // 'read_file', 'search_issues', etc.
  resource: string; // file path, issue ID, etc.
  allowed: boolean;
  reason?: string; // if denied
}
```

**Use platform-native audit logs where available:**
- GitHub: Audit log API
- GitLab: Audit events API
- Jira: Audit log (available in Premium/Enterprise)

### Rate Limiting

**Problem:** API rate limits vary by platform and plan.

**Platform limits (typical for free/standard plans):**
- **GitHub:** 5,000 requests/hour (authenticated), 60/hour (unauthenticated)
- **GitLab.com:** 2,000 requests/minute (authenticated)
- **Jira Cloud:** 10 requests/second per user

**Solution:**
- Implement request caching (5-10 minute expiration for read operations)
- Use exponential backoff on rate limit errors
- Queue requests if approaching limits
- Monitor usage and alert when approaching limits

### Error Handling

**Common failure modes:**
1. **Authentication failure:** OAuth token expired or revoked
2. **Permission denied:** User doesn't have access to resource
3. **Resource not found:** Ticket/file doesn't exist
4. **Rate limit exceeded:** Too many requests
5. **Network error:** Service unavailable

**Handling strategy:**
- Graceful degradation (return partial results when possible)
- Clear error messages for users
- Automatic retry with exponential backoff for transient failures
- Escalate to human when cannot proceed

---

## 5. Deployment Architecture

### Recommended Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Requirements Agents                                        │
│  ┌──────────────────────┐  ┌──────────────────────┐       │
│  │ requirements-        │  │ requirements-        │       │
│  │ drafting-assistant   │  │ analyst              │       │
│  └──────────┬───────────┘  └──────────┬───────────┘       │
│             │                           │                   │
│             └───────────┬───────────────┘                   │
│                         │                                   │
└─────────────────────────┼───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│  Secure MCP Wrapper Service                                 │
│  - Path restriction enforcement                             │
│  - Audit logging                                            │
│  - Request caching                                          │
│  - Rate limiting                                            │
└─────────┬───────────────────────────────────────────────────┘
          │
          ├──────────────┬──────────────┬──────────────┐
          ▼              ▼              ▼              ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  GitHub     │  │  GitLab     │  │  Jira       │  │  Linear     │
│  MCP Server │  │  MCP Server │  │  MCP Server │  │  (future)   │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```

### Configuration Management

**Environment variables (development):**
```bash
# GitHub
GITHUB_TOKEN=ghp_xxx
GITHUB_OAUTH_CLIENT_ID=xxx
GITHUB_OAUTH_CLIENT_SECRET=xxx

# GitLab
GITLAB_TOKEN=glpat-xxx
GITLAB_OAUTH_CLIENT_ID=xxx
GITLAB_OAUTH_CLIENT_SECRET=xxx

# Jira
ATLASSIAN_EMAIL=user@company.com
ATLASSIAN_TOKEN=xxx
ATLASSIAN_DOMAIN=company
```

**Secrets manager (production):**
- AWS Secrets Manager
- HashiCorp Vault
- Azure Key Vault
- 1Password Secrets Automation

---

## 6. Comparison: Official vs. Community Servers

### GitHub

**Official (github/github-mcp-server):**
- ✅ Maintained by GitHub
- ✅ OAuth 2.1 with PKCE
- ✅ Remote server (no local install)
- ✅ Integrated with Copilot IDEs
- ✅ Regular updates and new features
- ❌ Requires GitHub account

**Community alternatives:**
- Various community servers available but official server is recommended

**Recommendation:** Use official GitHub MCP server

### GitLab

**Official (GitLab Duo MCP Server):**
- ✅ Maintained by GitLab
- ✅ OAuth 2.0 Dynamic Client Registration
- ✅ Beta but actively developed
- ✅ Semantic code search on GitLab.com
- ✅ HTTP or stdio transport
- ❌ Beta (not GA yet)

**Community (zereight/gitlab-mcp, LuisCusihuaman/gitlab-mcp-server):**
- ⚠️ Community-maintained
- ✅ May have additional features
- ❌ Less reliable support

**Recommendation:** Use official GitLab MCP server (despite beta status)

### Jira/Atlassian

**Official (atlassian/atlassian-mcp-server):**
- ✅ Maintained by Atlassian
- ✅ OAuth 2.0
- ✅ Remote server (cloud-based)
- ✅ Both Jira and Confluence
- ✅ Respects user permissions
- ❌ Cloud only (no Server/Data Center)
- ❌ Beta status

**Community (aashari/mcp-server-atlassian-jira):**
- ✅ Generic HTTP tools (flexible)
- ✅ Direct API access
- ✅ JMESPath filtering
- ❌ Requires manual credential management
- ❌ Community-maintained

**Community (sooperset/mcp-atlassian):**
- ✅ Supports Server/Data Center
- ✅ Both Jira and Confluence
- ❌ Community-maintained

**Recommendation:**
- **Jira Cloud:** Use official Atlassian MCP server
- **Jira Server/Data Center:** Use sooperset/mcp-atlassian
- **Need flexibility:** Use aashari/mcp-server-atlassian-jira

---

## 7. Implementation Roadmap

### Phase 1: MVP (Week 1-2)

**Objective:** Get basic MCP integration working with existing servers

**Tasks:**
1. ✅ Set up GitHub MCP server (official)
   - Create OAuth app
   - Configure credentials
   - Test file read operations
2. ✅ Set up GitLab MCP server (official)
   - Configure OAuth
   - Test issue and MR operations
3. ✅ Set up Jira MCP server (official or aashari)
   - Configure authentication
   - Test ticket read and search
4. ⚠️ Basic path restriction wrapper (simple allow/block lists)
5. ⚠️ Environment variable configuration
6. ✅ Test both agents with MCP access

**Success criteria:**
- Both agents can read tickets via Jira MCP
- Both agents can read .feature files via GitHub/GitLab MCP
- Basic security (blocked paths work)
- Shared credentials between agents

### Phase 2: Production Hardening (Week 3-4)

**Objective:** Add security, audit logging, and monitoring

**Tasks:**
1. Implement comprehensive path restriction service
2. Add audit logging for all MCP operations
3. Implement request caching layer
4. Add rate limiting and retry logic
5. Move credentials to secrets manager
6. Add health checks for MCP servers
7. Create monitoring dashboard

**Success criteria:**
- All access logged for audit
- Sensitive files blocked
- Rate limits respected
- Credentials secure

### Phase 3: Enhancement (Week 5+)

**Objective:** Optimize performance and user experience

**Tasks:**
1. Implement intelligent caching strategies
2. Add support for additional git providers (Azure DevOps, Bitbucket)
3. Add support for additional ticketing systems (Linear, Shortcut)
4. Create admin dashboard for MCP usage monitoring
5. Implement multi-repository support
6. Add advanced search features

---

## 8. Gaps and Custom Development Required

### Critical (Required for MVP)

**✅ Path Restriction Wrapper Service**
- **Why:** MCP servers don't enforce path-based access control
- **Complexity:** Low (1-2 days)
- **Implementation:** Wrapper service with allow/block list matching

**✅ Audit Logging Service**
- **Why:** Need compliance and security audit trail
- **Complexity:** Low (1-2 days)
- **Implementation:** Centralized logging service, leverage platform audit APIs

### Important (Required for Production)

**Request Caching Layer**
- **Why:** Reduce API calls, improve performance, respect rate limits
- **Complexity:** Medium (2-3 days)
- **Implementation:** In-memory cache with TTL, Redis for distributed caching

**Secrets Management Integration**
- **Why:** Production security requirements
- **Complexity:** Low (1 day per secrets manager)
- **Implementation:** Integrate with AWS Secrets Manager, Vault, or Azure Key Vault

**Health Checks and Monitoring**
- **Why:** Detect MCP server outages, credential expiration
- **Complexity:** Medium (2-3 days)
- **Implementation:** Periodic health check, alerting on failures

### Nice-to-Have (Future Enhancement)

**Advanced Search Indexing**
- **Why:** Faster search across large repositories
- **Complexity:** High (1-2 weeks)
- **Implementation:** Elasticsearch or similar

**Multi-Tenancy Support**
- **Why:** Support multiple organizations/repositories
- **Complexity:** High (2-3 weeks)
- **Implementation:** Organization-scoped configuration, credential isolation

---

## 9. Final Recommendations

### Use Existing MCP Servers ✅

**No custom MCP server development needed.** Existing official and community servers meet all core requirements.

**Recommended servers:**
1. **GitHub:** Official `github/github-mcp-server` (remote, OAuth 2.1)
2. **GitLab:** Official GitLab Duo MCP server (OAuth 2.0)
3. **Jira Cloud:** Official `atlassian/atlassian-mcp-server` (OAuth 2.0)
4. **Jira Server/DC:** Community `sooperset/mcp-atlassian`

### Build Wrapper Services ⚠️

**Custom development required for:**
1. Path restriction enforcement (GitHub/GitLab)
2. Audit logging (all platforms)
3. Request caching (all platforms)
4. Secrets management integration

**Estimated effort:** 1-2 weeks for MVP, 2-3 weeks for production-ready

### Configuration Strategy

**Shared credentials between agents:**
- Single OAuth app per platform
- Environment variables or secrets manager
- Both agents reference same configuration

**Example shared config:**
```yaml
mcp:
  github:
    server: "https://github-mcp.api.github.com" # or local if self-hosted
    auth:
      type: oauth
      clientId: ${GITHUB_OAUTH_CLIENT_ID}
      clientSecret: ${GITHUB_OAUTH_CLIENT_SECRET}
    cache:
      ttl: 600 # 10 minutes

  jira:
    server: "https://company.atlassian.net"
    auth:
      type: oauth # or api-token for community servers
      email: ${ATLASSIAN_EMAIL}
      token: ${ATLASSIAN_TOKEN}
    cache:
      ttl: 300 # 5 minutes
```

---

## 10. Success Criteria

**MCP integration is successful when:**

✅ Both requirements-drafting-assistant and requirements-analyst can:
- Read tickets from Jira with full context (description, comments, related tickets)
- Search tickets using keywords and JQL
- Read `.feature` files from GitHub/GitLab repositories
- Search feature files for existing steps and patterns
- Access file history and commits

✅ Security requirements met:
- Business Owners have read-only access (no code access)
- Sensitive files are blocked (.env, secrets, keys)
- All access is logged for audit

✅ Performance acceptable:
- Ticket read: < 2 seconds
- File read: < 2 seconds
- Search operations: < 5 seconds
- Cache hit rate: > 80%

✅ Operational requirements met:
- Credentials configured once, shared between agents
- Graceful error handling (clear messages, automatic retry)
- Health checks detect MCP server issues
- Rate limits respected (no 429 errors)

---

## Conclusion

**Existing MCP servers for GitHub, GitLab, and Jira fully meet our requirements.** No custom MCP server development is needed.

**Focus custom development on:**
1. Security wrapper (path restrictions, audit logging)
2. Performance layer (caching, rate limiting)
3. Operational tooling (health checks, monitoring)

**Total estimated effort:** 3-4 weeks from prototype to production-ready system.

**Next steps:**
1. Set up accounts and OAuth apps for each platform
2. Deploy official MCP servers (GitHub, GitLab, Jira)
3. Build and test path restriction wrapper service
4. Implement audit logging
5. Test with both requirements agents
6. Iterate based on feedback
