# FOSS Components Evaluation for MCP Integration

**Date:** 2025-12-04
**Purpose:** Evaluate existing Free and Open Source Software (FOSS) solutions for MCP integration components

This document evaluates FOSS alternatives to custom development for the components identified in [mcp-server-evaluation.md](mcp-server-evaluation.md).

---

## Executive Summary

✅ **Excellent news: FOSS solutions exist for ALL identified components**

**Recommendation:** Use existing FOSS MCP gateways instead of building custom wrapper services.

**Best option:** **IBM ContextForge** or **Lasso MCP Gateway** provide comprehensive features out-of-the-box.

**Estimated effort reduction:** 3-4 weeks custom development → **1-2 days configuration**

---

## Components Identified for Custom Development

From [mcp-server-evaluation.md](mcp-server-evaluation.md), we identified:

1. **Path restriction wrapper service** (Critical)
2. Request caching layer (Optional)
3. Rate limiting and retry logic (Optional)
4. Secrets management integration (Optional)
5. Health checks and monitoring (Optional)

---

## 1. MCP-Specific Gateway Solutions

### Option 1: IBM ContextForge MCP Gateway ⭐ RECOMMENDED

**Repository:** https://github.com/IBM/mcp-context-forge
**Status:** Open Source (no official IBM support)
**License:** Apache 2.0 (assumed, typical for IBM open source)
**Language:** Python

#### Features

**✅ All requirements met:**

| Feature | Support | Details |
|---------|---------|---------|
| **Path filtering** | ✅ Yes | URI-based access controls, resource visibility |
| **Authentication** | ✅ Yes | Email auth, SSO, multi-tenant, RBAC |
| **Authorization** | ✅ Yes | Teams, role-based access control (RBAC) |
| **Rate limiting** | ✅ Yes | Built-in rate-limiting with user-scoped OAuth tokens |
| **Caching** | ✅ Yes | Redis-backed federation caching |
| **Retry logic** | ✅ Yes | Built-in retries |
| **Monitoring** | ✅ Yes | OpenTelemetry, Phoenix, Jaeger, Zipkin |
| **Secrets management** | ✅ Yes | OAuth token management, credential injection |

#### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Requirements Agents (both use same ContextForge)           │
│  - requirements-drafting-assistant                          │
│  - requirements-analyst                                     │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│  IBM ContextForge MCP Gateway                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Security & Access Control                           │   │
│  │ - Multi-tenant RBAC                                 │   │
│  │ - OAuth 2.0 with user-scoped tokens                 │   │
│  │ - URI-based resource filtering                      │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Performance                                         │   │
│  │ - Redis caching                                     │   │
│  │ - Built-in rate limiting                           │   │
│  │ - Automatic retries                                │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Observability                                       │   │
│  │ - OpenTelemetry tracing                            │   │
│  │ - LLM metrics (tokens, costs, performance)         │   │
│  │ - Admin UI with real-time logs                     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────┬───────────────────────────────────────────────────┘
          │
          ├──────────────┬──────────────┬──────────────┐
          ▼              ▼              ▼              ▼
    GitHub MCP     GitLab MCP     Jira MCP      REST APIs
    (wrapped)      (wrapped)      (wrapped)    (converted)
```

#### Key Advantages

**Protocol Flexibility:**
- Supports HTTP, JSON-RPC, WebSocket, SSE, stdio, streamable-HTTP
- Can wrap non-MCP REST APIs as MCP servers
- Federation across multiple MCP and REST services

**Enterprise Features:**
- 400+ tests for reliability
- Multi-tenant with email auth, teams, RBAC
- SSO configuration
- Dynamic client registration

**Observability:**
- OpenTelemetry integration
- Distributed tracing across federated gateways
- LLM-specific metrics (token usage, costs, performance)
- Admin UI with real-time log viewing, filtering, search

**Resource Management:**
- Redis-backed caching
- Resource visibility controls (URI-based access)
- SSE updates for resource streaming

#### Configuration Example

```yaml
# contextforge-config.yaml
gateway:
  port: 8080

authentication:
  method: oauth
  sso:
    enabled: true
    provider: azure-ad # or google, okta, etc.

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
        - docs/**/*.md
        - sample-project/context/**/*
      blocked_paths:
        - .env*
        - secrets/**/*
        - "**/*.key"
        - "**/*.pem"

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
  ttl: 600 # 10 minutes

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
    endpoint: http://jaeger:4318
  metrics:
    llm_tracking: true
    cost_tracking: true
```

#### Deployment

**Via pip:**
```bash
pip install mcp-contextforge-gateway
contextforge serve --config contextforge-config.yaml
```

**Via Docker:**
```bash
docker run -p 8080:8080 \
  -v ./contextforge-config.yaml:/config.yaml \
  -e GITHUB_CLIENT_ID=${GITHUB_CLIENT_ID} \
  -e GITHUB_CLIENT_SECRET=${GITHUB_CLIENT_SECRET} \
  ibm/mcp-contextforge
```

#### Pros
- ✅ **Complete solution** - all features needed out-of-the-box
- ✅ **Battle-tested** - 400+ tests
- ✅ **Observability** - OpenTelemetry, admin UI, real-time logs
- ✅ **Flexible** - wraps REST APIs as MCP, multiple protocols
- ✅ **Enterprise-ready** - multi-tenant, RBAC, SSO

#### Cons
- ⚠️ No official IBM support (community-maintained)
- ⚠️ Python dependency (if team prefers TypeScript/Go)
- ⚠️ Redis required for caching (additional infrastructure)

---

### Option 2: Lasso MCP Gateway 🛡️ SECURITY-FOCUSED

**Repository:** https://github.com/lasso-security/mcp-gateway
**Status:** Open Source
**License:** Not specified in search results
**Language:** Python

#### Features

**✅ Security-first approach:**

| Feature | Support | Details |
|---------|---------|---------|
| **Path filtering** | ✅ Yes | Security scanner analyzes server reputation |
| **Authentication** | ✅ Yes | Token validation, JWT support |
| **Authorization** | ⚠️ Plugin | Via custom plugins |
| **Rate limiting** | ⚠️ Unknown | Not mentioned in docs |
| **Caching** | ⚠️ Unknown | Not mentioned in docs |
| **Security scanning** | ✅ Yes | Reputation analysis, hidden instruction detection |
| **Monitoring** | ✅ Yes | Xetrack plugin for tracking, logging |
| **Guardrails** | ✅ Yes | Prompt injection, PII detection, token masking |

#### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Requirements Agents                                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│  Lasso MCP Gateway (Security-Focused)                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Security Scanner                                    │   │
│  │ - Server reputation analysis                       │   │
│  │ - Hidden instruction detection                     │   │
│  │ - Automatic blocking (threshold: 30)               │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Guardrail Plugins                                   │   │
│  │ - Basic: Token/secret masking                      │   │
│  │ - Presidio: PII detection (credit cards, SSN, etc) │   │
│  │ - Lasso: Advanced AI safety (custom policies)     │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Monitoring (Xetrack Plugin)                        │   │
│  │ - SQLite/log-based experiment tracking             │   │
│  │ - Structured data monitoring                       │   │
│  │ - CLI/Python queryable                            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────┬───────────────────────────────────────────────────┘
          │
          └─────► MCP Servers (GitHub, GitLab, Jira)
```

#### Configuration Example

```json
{
  "servers": {
    "github": {
      "command": "github-mcp-server",
      "env": {
        "GITHUB_TOKEN": "${GITHUB_TOKEN}"
      }
    },
    "jira": {
      "command": "jira-mcp-server",
      "env": {
        "JIRA_TOKEN": "${JIRA_TOKEN}"
      }
    }
  }
}
```

**Launch with plugins:**
```bash
mcp-gateway -p basic -p presidio --config mcp.json
```

#### Available Plugins

**Basic Guardrails:**
- GitHub tokens, AWS access tokens, JWT tokens
- Hugging Face tokens
- No PII detection

**Presidio (Microsoft):**
- Credit cards, IP addresses
- Email addresses, phone numbers
- Social Security Numbers (SSN)

**Lasso (Enterprise):**
- Full MCP interaction visibility
- Always-on monitoring
- Prompt injection protection
- Sensitive data leakage prevention
- Customizable natural language policies

**Xetrack (Monitoring):**
- Experiment tracking
- Tool call monitoring
- SQLite database or log files
- CLI and Python query interface

#### Security Scanner

**Analyzes:**
- Server reputation (marketplace + GitHub data)
- Tool descriptions for hidden instructions
- Automatic blocking based on reputation score (threshold: 30)
- Updates configuration with scan results

#### Pros
- ✅ **Security-focused** - purpose-built for AI safety
- ✅ **Plugin architecture** - extensible, modular
- ✅ **PII detection** - Presidio integration
- ✅ **Reputation scanning** - analyzes server trustworthiness
- ✅ **Lightweight** - minimal dependencies

#### Cons
- ⚠️ Limited documentation on rate limiting/caching
- ⚠️ Fewer enterprise features than ContextForge
- ⚠️ No mention of multi-tenancy or RBAC
- ⚠️ Security focus may be overkill for internal tools

---

### Option 3: Docker MCP Gateway 🐳 CONTAINER-NATIVE

**Repository:** https://github.com/docker/mcp-gateway
**Status:** Official Docker project
**License:** Apache 2.0
**Language:** Go (assumed)

#### Features

| Feature | Support | Details |
|---------|---------|---------|
| **Path filtering** | ⚠️ Via policy | Policy management via CLI |
| **Authentication** | ✅ Yes | Credential management, OAuth flows |
| **Authorization** | ⚠️ Via policy | Policy CLI available |
| **Rate limiting** | ❌ No | Not mentioned |
| **Caching** | ❌ No | Not mentioned |
| **Container isolation** | ✅ Yes | Runs MCP servers in isolated containers |
| **Monitoring** | ✅ Yes | Logging, call tracing |
| **Lifecycle management** | ✅ Yes | On-demand container startup |

#### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Requirements Agents                                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│  Docker MCP Gateway                                         │
│  - Credential injection                                     │
│  - Policy enforcement (docker mcp policy)                   │
│  - Logging and call tracing                                 │
└─────────┬───────────────────────────────────────────────────┘
          │
          ├──────────────┬──────────────┬──────────────┐
          ▼              ▼              ▼              ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
    │ GitHub   │   │ GitLab   │   │ Jira     │   │ Other    │
    │ MCP      │   │ MCP      │   │ MCP      │   │ MCP      │
    │ Container│   │ Container│   │ Container│   │ Container│
    └──────────┘   └──────────┘   └──────────┘   └──────────┘
    (isolated)     (isolated)     (isolated)     (isolated)
```

#### Key Advantages

**Container Isolation:**
- Each MCP server runs in isolated Docker container
- Restricted privileges, network access, resource usage
- Security through containerization

**Lifecycle Management:**
- On-demand container startup when AI applications request tools
- Automatic cleanup when no longer needed

**Integration:**
- Built into Docker Desktop
- Seamless OAuth flows
- Credential management via Docker secrets

#### Configuration

Uses 4 configuration files in `~/.docker/mcp/`:
1. **docker-mcp.yaml** - Server catalog
2. **registry.yaml** - Enabled servers
3. **config.yaml** - Per-server config
4. **tools.yaml** - Per-server tool enablement

**Policy management:**
```bash
docker mcp policy --help
```

#### Pros
- ✅ **Container isolation** - strong security boundary
- ✅ **Official Docker project** - reliable, maintained
- ✅ **Integrated with Docker Desktop** - easy setup
- ✅ **Resource management** - container limits

#### Cons
- ❌ **No rate limiting** - would need separate solution
- ❌ **No caching** - would need separate solution
- ⚠️ **Docker dependency** - requires Docker runtime
- ⚠️ **Limited documentation** - policy features unclear

---

### Option 4: Traefik Hub MCP Gateway 🔀 OAUTH-FOCUSED

**Documentation:** https://doc.traefik.io/traefik-hub/mcp-gateway/
**Status:** Commercial with free tier
**License:** Proprietary (Traefik Hub), but Traefik Proxy is open source
**Language:** Go

#### Features

| Feature | Support | Details |
|---------|---------|---------|
| **Path filtering** | ❌ No | Operates on MCP methods, not HTTP paths |
| **Authentication** | ✅ Yes | OAuth 2.1/2.0, JWT integration |
| **Authorization** | ✅ Yes | Task-Based Access Control (TBAC), policies |
| **Rate limiting** | ❌ No | Not mentioned for MCP Gateway |
| **Caching** | ❌ No | Not mentioned for MCP Gateway |
| **List filtering** | ✅ Yes | Controls tool/resource discovery |
| **Session affinity** | ✅ Yes | HRW load balancing |

#### Key Advantages

**OAuth 2.1/2.0 Compliant:**
- Resource Server specification implementation
- JWT authentication middleware
- Centralized governance across MCP servers

**Fine-Grained Authorization:**
- Expression-based policy matching
- Runtime access policies (allow/deny)
- List filtering for tool/resource discovery

**Resource Management:**
- Automatic resource metadata discovery (OAuth well-known endpoints)
- Session affinity using Highest Random Weight (HRW) load balancing

#### Pros
- ✅ **OAuth-native** - proper spec compliance
- ✅ **Policy-based** - flexible authorization
- ✅ **List filtering** - control discovery

#### Cons
- ❌ **No path filtering** - doesn't solve our main requirement
- ❌ **No rate limiting** - missing key feature
- ❌ **No caching** - missing key feature
- ⚠️ **Commercial** - requires Traefik Hub (not fully open source)

---

### Option 5: Microsoft MCP Gateway ☁️ KUBERNETES-NATIVE

**Repository:** https://github.com/microsoft/mcp-gateway
**Status:** Open Source
**License:** MIT (typical for Microsoft OSS)
**Language:** Likely C#/Go (not confirmed)

#### Features

| Feature | Support | Details |
|---------|---------|---------|
| **Path filtering** | ❌ No | Not mentioned |
| **Authentication** | ✅ Yes | Azure Entra ID, bearer token |
| **Authorization** | ✅ Yes | RBAC (mcp.admin, mcp.engineer roles) |
| **Rate limiting** | ❌ No | Not mentioned |
| **Caching** | ❌ No | Not mentioned |
| **Session routing** | ✅ Yes | Session-aware stateful routing |
| **Kubernetes** | ✅ Yes | Lifecycle management in K8s |

#### Key Advantages

**Kubernetes-Native:**
- Designed for K8s environments
- Scalable, session-aware stateful routing
- Lifecycle management of MCP servers

**Access Control:**
- Azure Entra ID integration
- Role-based access control (RBAC)
- Granular read/write permissions

**RESTful Management:**
- Control plane APIs for deployment, updates, lifecycle
- Status monitoring and log access

#### Pros
- ✅ **Kubernetes-ready** - built for K8s
- ✅ **Azure integration** - natural fit for Azure users
- ✅ **Session routing** - stateful session handling

#### Cons
- ❌ **No path filtering** - missing our main requirement
- ❌ **No rate limiting/caching** - missing key features
- ⚠️ **Azure-centric** - best for Azure environments
- ⚠️ **Overkill** - if not using Kubernetes

---

## 2. Middleware Libraries (Build Your Own Gateway)

If none of the gateways fit, we can build using middleware libraries.

### FastMCP Middleware (Python)

**Repository:** https://github.com/jlowin/fastmcp
**Status:** Open Source
**License:** MIT (assumed)
**Language:** Python

#### Features

**MCP-Native Middleware:**
- Context-aware logic for authentication, authorization, caching
- Interceptor pattern for requests and responses
- Built-in templates for common patterns

**Available Middleware:**

**Authentication:**
- Token validation from HTTP headers
- JWT authentication
- Custom authentication logic

**Authorization:**
- **permit-fastmcp**: Permit.io integration (ABAC policies)
- **cerbos-fastmcp**: Cerbos policy engine integration
- Attribute-Based Access Control (ABAC) - evaluate tool arguments as attributes

**Built-in Patterns:**
- Logging
- Error handling
- Rate limiting templates
- Timing for performance monitoring

#### Example Usage

```python
from fastmcp import FastMCP
from fastmcp.middleware import authenticate, authorize, cache

mcp = FastMCP("requirements-gateway")

# Authentication middleware
@mcp.middleware
async def auth_middleware(request, call_next):
    token = request.headers.get("Authorization")
    if not validate_token(token):
        raise Unauthorized("Invalid token")
    return await call_next(request)

# Authorization middleware (Permit.io)
from permit_fastmcp import PermitMiddleware

mcp.use(PermitMiddleware(
    api_key="permit_key_...",
    pdp_url="https://cloudpdp.api.permit.io"
))

# Caching middleware
@mcp.middleware
async def cache_middleware(request, call_next):
    cache_key = f"{request.method}:{request.path}"
    if cached := await redis.get(cache_key):
        return cached

    response = await call_next(request)
    await redis.set(cache_key, response, ex=600)
    return response

# Rate limiting middleware
@mcp.middleware
async def rate_limit_middleware(request, call_next):
    user = request.user
    if await is_rate_limited(user):
        raise TooManyRequests("Rate limit exceeded")
    return await call_next(request)
```

#### Pros
- ✅ **Flexible** - build exactly what you need
- ✅ **MCP-native** - designed for MCP protocol
- ✅ **Python** - easy to integrate with Python ecosystem
- ✅ **Middleware ecosystem** - Permit.io, Cerbos integrations

#### Cons
- ⚠️ **DIY** - requires building the gateway yourself
- ⚠️ **Less comprehensive** - need to add all features manually
- ⚠️ **Maintenance** - you own the codebase

---

## 3. General-Purpose API Gateways

For completeness, here are general-purpose API gateways (though MCP-specific solutions are better).

### Kong Gateway

**Features:** Authentication, rate limiting, caching, request/response transformation
**Pros:** Battle-tested, massive plugin ecosystem (70+)
**Cons:** Not MCP-aware, requires significant configuration, PostgreSQL/Cassandra dependency

### Apache APISIX

**Features:** Cloud-native, load balancing, auth, rate limiting, API security
**Pros:** High performance, Kubernetes-native, plugin support
**Cons:** Not MCP-aware, learning curve

### Tyk

**Features:** Authentication, quotas, rate limiting, versioning, GUI dashboard
**Pros:** Full API management suite, graphical interface
**Cons:** Not MCP-aware, complex setup

**Verdict:** ❌ **Not recommended** - MCP-specific gateways are better fit

---

## 4. Supporting Libraries

### Caching

**Redis:**
- **Python:** redis-py, aioredis-py
- **TypeScript:** ioredis, node-redis
- **Features:** High performance, distributed caching, TTL support

**In-Memory (simpler):**
- **Python:** functools.lru_cache, cachetools
- **TypeScript:** node-cache, memory-cache

### Rate Limiting

**Redis-Based:**
- **TypeScript:** rate-limit-redis (npm)
- **Python:** Redis + custom implementation
- **Algorithms:** Token bucket, sliding window

**Alternatives:**
- **Python:** slowapi, limits
- **TypeScript:** express-rate-limit

---

## Recommendations

### Scenario 1: Best Overall Solution ⭐

**Use IBM ContextForge MCP Gateway**

**Why:**
- ✅ All features needed (path filtering, auth, rate limiting, caching)
- ✅ Battle-tested (400+ tests)
- ✅ Observability out-of-the-box (OpenTelemetry, admin UI)
- ✅ Multi-tenant RBAC
- ✅ Can wrap REST APIs as MCP

**Setup time:** 1-2 days
**Maintenance:** Low (community-maintained)

**Configuration:**
```yaml
# Simple contextforge-config.yaml
servers:
  - github (with path filters)
  - gitlab (with path filters)
  - jira

features:
  - Redis caching (TTL: 10min)
  - Rate limiting (60 req/min default)
  - OAuth authentication
  - RBAC (business-owner, requirements-analyst)
  - OpenTelemetry tracing
```

---

### Scenario 2: Security-First ⛑️

**Use Lasso MCP Gateway**

**Why:**
- ✅ Security scanner (reputation analysis)
- ✅ Guardrails (prompt injection, PII detection)
- ✅ Plugin architecture
- ✅ Lightweight

**Setup time:** 1 day
**Maintenance:** Low

**Best for:** High-security environments, AI safety requirements

---

### Scenario 3: Already Using Docker 🐳

**Use Docker MCP Gateway**

**Why:**
- ✅ Container isolation (strong security)
- ✅ Integrated with Docker Desktop
- ✅ Official Docker project

**Add:**
- Traefik Proxy (open source) for rate limiting
- Redis for caching

**Setup time:** 2-3 days
**Maintenance:** Medium (separate components)

---

### Scenario 4: Kubernetes Environment ☸️

**Use Microsoft MCP Gateway**

**Why:**
- ✅ Kubernetes-native
- ✅ Session-aware stateful routing
- ✅ Azure Entra ID integration

**Add:**
- Traefik or NGINX Ingress for path filtering
- Redis for caching

**Setup time:** 2-3 days
**Maintenance:** Medium

---

### Scenario 5: Maximum Control 🛠️

**Build with FastMCP + Redis**

**Why:**
- ✅ Complete control
- ✅ MCP-native middleware
- ✅ Python ecosystem

**Components:**
- FastMCP (middleware framework)
- permit-fastmcp or cerbos-fastmcp (authorization)
- Redis (caching + rate limiting)
- Custom path filtering middleware

**Setup time:** 3-5 days
**Maintenance:** High (you own the code)

---

## Feature Comparison Matrix

| Feature | ContextForge | Lasso | Docker | Traefik | Microsoft | FastMCP DIY |
|---------|--------------|-------|--------|---------|-----------|-------------|
| **Path Filtering** | ✅ URI-based | ✅ Scanner | ⚠️ Policy | ❌ No | ❌ No | ✅ Custom |
| **Authentication** | ✅ OAuth/SSO | ✅ Token | ✅ OAuth | ✅ OAuth | ✅ Azure | ✅ Custom |
| **Authorization** | ✅ RBAC | ⚠️ Plugin | ⚠️ Policy | ✅ TBAC | ✅ RBAC | ✅ Custom |
| **Rate Limiting** | ✅ Built-in | ❌ Unknown | ❌ No | ❌ No | ❌ No | ✅ Custom |
| **Caching** | ✅ Redis | ❌ Unknown | ❌ No | ❌ No | ❌ No | ✅ Custom |
| **Retry Logic** | ✅ Built-in | ❌ Unknown | ❌ No | ❌ No | ❌ No | ✅ Custom |
| **Monitoring** | ✅ OTel+UI | ✅ Xetrack | ✅ Logs | ❌ No | ✅ Logs | ✅ Custom |
| **Multi-Tenant** | ✅ Yes | ❌ No | ❌ No | ❌ No | ⚠️ Azure | ✅ Custom |
| **MCP-Native** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **Setup Time** | 1-2 days | 1 day | 2-3 days | 2-3 days | 2-3 days | 3-5 days |
| **Maintenance** | Low | Low | Medium | Medium | Medium | High |
| **Cost** | Free | Free | Free | $$ (Hub) | Free | Free |

---

## Implementation Plan with IBM ContextForge

### Phase 1: Setup and Configuration (Day 1)

**Tasks:**
1. Install IBM ContextForge
   ```bash
   pip install mcp-contextforge-gateway
   ```

2. Create configuration file
   ```yaml
   # contextforge-config.yaml
   gateway:
     port: 8080

   servers:
     - name: github
       type: mcp
       url: https://github-mcp.api.github.com
       auth: { type: oauth, ... }
       filters:
         allowed_paths: [features/**, docs/**]
         blocked_paths: [.env*, secrets/**]

     - name: jira
       type: mcp
       url: https://company.atlassian.net
       auth: { type: oauth, ... }

   caching:
     enabled: true
     backend: redis
     ttl: 600

   rate_limiting:
     enabled: true
     default: 60
   ```

3. Set up Redis (if using caching)
   ```bash
   docker run -d -p 6379:6379 redis:alpine
   ```

4. Configure environment variables
   ```bash
   export GITHUB_CLIENT_ID=xxx
   export GITHUB_CLIENT_SECRET=xxx
   export JIRA_CLIENT_ID=xxx
   export JIRA_CLIENT_SECRET=xxx
   ```

### Phase 2: Testing (Day 2)

**Tasks:**
1. Start gateway
   ```bash
   contextforge serve --config contextforge-config.yaml
   ```

2. Test with requirements-drafting-assistant
   - Read tickets from Jira
   - Read feature files from GitHub
   - Verify path restrictions work

3. Test with requirements-analyst
   - Same tests
   - Verify shared credentials work

4. Performance testing
   - Verify caching works (second request faster)
   - Verify rate limiting works (429 after limit)

### Phase 3: Monitoring Setup (Optional)

**Tasks:**
1. Set up Jaeger for tracing
   ```bash
   docker run -d -p 16686:16686 -p 4318:4318 jaegertracing/all-in-one
   ```

2. Enable OpenTelemetry in config
   ```yaml
   observability:
     opentelemetry:
       enabled: true
       endpoint: http://localhost:4318
   ```

3. Access admin UI
   - Navigate to http://localhost:8080/admin
   - View real-time logs, traces, metrics

---

## Cost-Benefit Analysis

### Custom Development (Original Plan)

**Effort:** 3-4 weeks
**Components to build:**
- Path restriction service
- Request caching layer
- Rate limiting + retry logic
- Health checks
- Monitoring/logging

**Ongoing maintenance:** High (you own all code)

### IBM ContextForge (Recommended)

**Effort:** 1-2 days configuration
**Components needed:**
- Install ContextForge
- Configure YAML file
- Set up Redis (optional, for caching)
- Set up Jaeger (optional, for tracing)

**Ongoing maintenance:** Low (community-maintained)

**Savings:** **2.5-3.5 weeks of development time**

---

## Conclusion

**Use IBM ContextForge MCP Gateway** - it provides all needed features out-of-the-box:

✅ Path filtering (URI-based resource visibility)
✅ Authentication (OAuth, SSO, multi-tenant)
✅ Authorization (RBAC with teams and roles)
✅ Rate limiting (built-in, user-scoped)
✅ Caching (Redis-backed)
✅ Retry logic (built-in)
✅ Monitoring (OpenTelemetry, admin UI, real-time logs)
✅ Protocol flexibility (wraps REST APIs as MCP)

**Alternative:** Lasso MCP Gateway if security is top priority (guardrails, PII detection, reputation scanning)

**Avoid:** Building custom solution - waste of 3-4 weeks when FOSS solutions exist

**Next steps:**
1. Install IBM ContextForge
2. Create configuration with path filters
3. Test with both requirements agents
4. Deploy to production

**Total implementation time:** 1-2 days vs. 3-4 weeks custom development
