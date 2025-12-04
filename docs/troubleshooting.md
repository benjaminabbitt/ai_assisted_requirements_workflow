# Troubleshooting Guide

Common issues and solutions for the AI-augmented requirements workflow.

---

## High AI Escalation Rate

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

---

## Incorrect Specs Generated

**Symptom:** AI drafts specs that don't match intent

**Root Cause:** Ambiguous context files or outdated API specs

**Fix:**

1. Review the execution log - what did AI reference?
2. Check if business rules in business.md are clear
3. Verify API specs match actual API implementation
4. Add examples to context files for similar scenarios
5. Ensure ticket acceptance criteria are specific and testable

**Examples of unclear vs clear:**

**Unclear:**
```markdown
## Acceptance Criteria
- User can reset password
```

**Clear:**
```markdown
## Acceptance Criteria
- [ ] User receives email with reset link within 2 minutes
- [ ] Reset link expires after 24 hours
- [ ] Password must be 12+ characters with special char
- [ ] Successful reset logs security event
```

---

## Tests Fail After Implementation

**Symptom:** Scenarios pass but unit tests show violations

**Root Cause:** Developer didn't follow IoC patterns or TDD guidance

**Fix:**

1. Run standards-compliance agent on code
2. Review developer-implementation guidance
3. Ensure tests use primary constructors + mocks
4. Verify production factories have `// coverage:ignore` and no business logic
5. Check that step definitions follow tech_standards.md patterns

---

## Agent Can't Read Tickets

**Symptom:** Agent fails with "Unable to fetch ticket" or "Access denied"

**Root Cause:** MCP gateway configuration or credentials issue

**Check:**

1. Is MCP gateway running?
   ```bash
   curl http://localhost:8080/health
   # Expected: 200 OK
   ```

2. Are credentials correct?
   ```bash
   # Check environment variables
   echo $JIRA_CLIENT_ID
   echo $JIRA_CLIENT_SECRET
   ```

3. Is MCP endpoint configured correctly in agent?
   - Check agent config points to correct gateway URL
   - Verify gateway URL is accessible from agent

4. Test MCP access directly:
   ```bash
   curl http://localhost:8080/jira/ticket/PROJ-1234
   # Expected: Ticket data in JSON
   ```

**Fix:**

- Restart MCP gateway: `contextforge serve --config contextforge-config.yaml`
- Verify OAuth tokens are valid and not expired
- Check firewall/network allows agent → gateway communication
- Review gateway logs for detailed error messages

---

## Path Filtering Not Working

**Symptom:** Agent can read sensitive files like `.env` or `secrets/`

**Root Cause:** Path restrictions not configured or not enforced by MCP gateway

**Check:**

1. Test directly with blocked path:
   ```bash
   curl http://localhost:8080/github/read?path=.env
   # Expected: 403 Forbidden
   ```

2. Test with allowed path:
   ```bash
   curl http://localhost:8080/github/read?path=features/auth/login.feature
   # Expected: 200 OK with file contents
   ```

3. Review MCP gateway config:
   ```yaml
   filters:
     allowed_paths:
       - features/**/*.feature
       - docs/**/*.md
     blocked_paths:
       - .env*
       - secrets/**/*
       - "**/*.key"
   ```

**Fix:**

- Use glob patterns correctly: `features/**/*.feature` not `features/*.feature`
- Ensure blocked_paths includes all sensitive file patterns
- Restart gateway after config changes
- Test each path pattern individually

---

## Rate Limiting Too Aggressive

**Symptom:** Agent frequently gets 429 Too Many Requests errors

**Root Cause:** Rate limits set too low for agent usage patterns

**Check:**

```bash
# Send multiple requests quickly
for i in {1..65}; do
  curl http://localhost:8080/jira/ticket/PROJ-1234
done
# Count how many succeed vs 429 errors
```

**Fix:**

Adjust in MCP gateway config:

```yaml
rate_limiting:
  enabled: true
  default:
    requests_per_minute: 120  # Increase from 60
  per_user:
    requirements-analyst: 300  # Increase for AI agents
```

Restart gateway after changes.

**Alternative:** Implement request batching in agent to reduce API calls.

---

## Caching Not Working

**Symptom:** Every request is slow, no performance improvement from caching

**Root Cause:** Redis not running or caching not enabled

**Check:**

1. Is Redis running?
   ```bash
   redis-cli ping
   # Expected: PONG
   ```

2. Is caching enabled in config?
   ```yaml
   caching:
     enabled: true
     backend: redis
     ttl: 600
   ```

3. Test cache hit:
   ```bash
   # First request (cache miss)
   time curl http://localhost:8080/jira/ticket/PROJ-1234
   # Expected: ~2 seconds

   # Second request (cache hit)
   time curl http://localhost:8080/jira/ticket/PROJ-1234
   # Expected: ~100ms (20x faster)
   ```

**Fix:**

- Start Redis: `docker run -d -p 6379:6379 redis:alpine`
- Verify MCP gateway can connect to Redis
- Check gateway logs for cache-related errors
- Verify cache TTL is reasonable (not too short)
- Check cache hit rate in admin UI: `http://localhost:8080/admin`

---

## Spec PRs Not Triggering bo-review

**Symptom:** Spec PRs created but bo-review agent doesn't run

**Root Cause:** GitHub Actions trigger not configured correctly

**Check:**

1. Is workflow file present?
   ```bash
   cat .github/workflows/gherkin-validation.yml
   ```

2. Does trigger match branch pattern?
   ```yaml
   on:
     pull_request:
       paths:
         - 'features/**/*.feature'  # Must match your structure
   ```

3. Is PR branch named correctly?
   - Expected: `spec/PROJ-1234`
   - Check: Does it start with `spec/`?

4. Check GitHub Actions tab for errors

**Fix:**

- Ensure workflow file is in `.github/workflows/`
- Verify paths match your feature file location
- Check if `startsWith(github.head_ref, 'spec/')` condition is correct
- Review GitHub Actions logs for detailed errors
- Test workflow manually: GitHub → Actions → Select workflow → Run workflow

---

## Implementation PR Blocked by @pending Check

**Symptom:** CI blocks implementation PR even though @pending tags removed

**Root Cause:** Check looking for wrong story ID or tags still present for different story

**Check:**

1. What story ID is the branch?
   ```bash
   echo "impl/PROJ-1234" | sed 's/^impl\///'
   # Expected: PROJ-1234
   ```

2. Search for @pending tags for this story:
   ```bash
   grep -r "@pending.*@story-PROJ-1234\|@story-PROJ-1234.*@pending" features/
   ```

3. Are there @pending tags for OTHER stories?
   ```bash
   grep -r "@pending" features/
   # These are OK - CI should only block for THIS story
   ```

**Fix:**

- Ensure @pending tags removed ONLY for the story being implemented
- Other stories can still have @pending tags
- Verify branch name matches story ID in tags
- Check CI script extracts story ID correctly from branch name

---

## Context File Freshness Alerts

**Symptom:** CI warns that context files are stale

**Root Cause:** Context files haven't been reviewed in 30-60 days

**Fix:**

1. **Review the context file** - is it still accurate?
2. **If accurate:** Update version header with current date and bump patch version:
   ```yaml
   ---
   version: 1.2.1  # Bump patch
   last-reviewed: 2025-01-15  # Update date
   reviewed-by: @your-username
   changelog: Reviewed - no changes needed
   ---
   ```

3. **If outdated:** Make necessary updates, bump minor/major version:
   ```yaml
   ---
   version: 1.3.0  # Minor bump for additions
   last-reviewed: 2025-01-15
   reviewed-by: @your-username
   changelog: Added new email validation patterns
   ---
   ```

4. Commit and push changes

**Prevention:** Schedule quarterly reviews of context files during sprint planning.

---

## Scenarios Passing Locally But Failing in CI

**Symptom:** `godog run` passes on developer machine but fails in CI

**Root Cause:** Environment differences (database, test data, services)

**Check:**

1. Are test dependencies available in CI?
   - Database (PostgreSQL, MySQL)
   - External services (mocked or stubbed)
   - Test data fixtures

2. Are environment variables set correctly in CI?
   ```yaml
   env:
     DATABASE_URL: postgres://...
     TEST_MODE: true
   ```

3. Are tests running in isolation?
   - Check for shared state between scenarios
   - Verify database cleanup between tests

**Fix:**

- Add database setup to CI workflow
- Use docker-compose for services in CI
- Ensure test fixtures are committed
- Add test data seeding to CI steps
- Verify cleanup happens between scenario runs

---

## BO Can't Access Feature Files

**Symptom:** Business Owner gets 403 Forbidden when trying to read .feature files via MCP

**Root Cause:** Permissions not configured or path not in allowed list

**Check:**

1. Is BO's role configured in MCP gateway?
   ```yaml
   roles:
     - name: business-owner
       permissions:
         - read:tickets
         - read:features
   ```

2. Are .feature files in allowed_paths?
   ```yaml
   allowed_paths:
     - features/**/*.feature
   ```

3. Test BO credentials:
   ```bash
   curl -H "Authorization: Bearer $BO_TOKEN" \
        http://localhost:8080/github/read?path=features/auth/login.feature
   ```

**Fix:**

- Add BO role to MCP gateway config
- Ensure .feature file paths in allowed_paths
- Verify BO has correct OAuth token
- Check MCP gateway logs for authorization errors

---

## Workflow Too Slow

**Symptom:** Story → PR taking hours instead of minutes

**Root Cause:** Multiple possible issues

**Diagnose:**

1. **Check each stage:**
   - Ticket labeled → How long until AI creates PR? (Target: < 5 min)
   - AI analysis → How long for confident stories? (Target: < 1 min)
   - PR created → How long until BO reviews? (Target: < 4 hours)

2. **Identify bottleneck:**
   - If AI slow: Check MCP gateway performance, API rate limits
   - If BO review slow: Increase SLA enforcement, enable auto-merge for high confidence
   - If escalation slow: Improve context files to reduce escalations

**Fix:**

- **Slow AI:** Increase caching TTL, upgrade MCP gateway resources
- **Slow review:** Implement SLA reminders, enable auto-merge for HIGH confidence PRs
- **High escalation rate:** Invest in context file completeness (see "High AI Escalation Rate" above)

---

## Need Help?

### Documentation

- [Workflow Guide](workflow.md) - Complete end-to-end process
- [Integration Guide](integration.md) - MCP setup and API integration
- [CI Configuration](ci-configuration.md) - CI/CD setup
- [Metrics Guide](metrics.md) - Measuring success

### Community

- **Questions about this workflow:** Review documentation in `docs/`
- **Questions about IBM ContextForge:** https://github.com/IBM/mcp-context-forge/issues
- **Questions about MCP:** https://modelcontextprotocol.io/

### Common Log Locations

- **MCP Gateway:** `contextforge.log` or stdout
- **GitHub Actions:** Actions tab → Select workflow → View logs
- **Agent Executions:** See [CLAUDE.md](../CLAUDE.md) for logging requirements
